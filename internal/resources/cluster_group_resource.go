package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &ClusterGroupResource{}
	_ resource.ResourceWithImportState = &ClusterGroupResource{}
	_ resource.ResourceWithIdentity    = &ClusterGroupResource{}
)

func NewClusterGroupResource() resource.Resource {
	return &ClusterGroupResource{}
}

type ClusterGroupResource struct {
	client *netbox.APIClient
}

// GetClient returns the API client for testing purposes.
func (r *ClusterGroupResource) GetClient() *netbox.APIClient {
	return r.client
}

type ClusterGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *ClusterGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_group"
}

func (r *ClusterGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a cluster group in Netbox. Cluster groups provide a way to organize virtualization clusters for better management (e.g., by datacenter, environment, or team).",
		Attributes: map[string]schema.Attribute{
			"id":   nbschema.IDAttribute("cluster group"),
			"name": nbschema.NameAttribute("cluster group", 100),
			"slug": nbschema.SlugAttribute("cluster group"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("cluster group"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ClusterGroupResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *ClusterGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *ClusterGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating cluster group", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the request
	clusterGroupRequest := netbox.ClusterGroupRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&clusterGroupRequest, data.Description)

	// Apply tags and custom_fields
	utils.ApplyTagsFromSlugs(ctx, r.client, &clusterGroupRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, &clusterGroupRequest, data.CustomFields, &resp.Diagnostics)

	// Store plan values for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Create via API
	clusterGroup, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterGroupsCreate(ctx).ClusterGroupRequest(clusterGroupRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_cluster_group",
			ResourceName: "this.cluster_group",
			SlugValue:    data.Slug.ValueString(),
			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				list, _, lookupErr := r.client.VirtualizationAPI.VirtualizationClusterGroupsList(lookupCtx).Slug([]string{slug}).Execute()
				if lookupErr != nil {
					return "", lookupErr
				}
				if list != nil && len(list.Results) > 0 {
					return fmt.Sprintf("%d", list.Results[0].GetId()), nil
				}
				return "", nil
			},
		}
		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "create cluster group", httpResp, http.StatusCreated) {
		return
	}

	r.mapClusterGroupToState(clusterGroup, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, clusterGroup.HasTags(), clusterGroup.GetTags(), planTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, clusterGroup.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	tflog.Trace(ctx, "created a cluster group resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterGroupID := data.ID.ValueString()
	clusterGroupIDInt := utils.ParseInt32FromString(clusterGroupID)
	if clusterGroupIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Cluster Group ID", fmt.Sprintf("Cluster Group ID must be a number, got: %s", clusterGroupID))
		return
	}
	clusterGroup, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterGroupsRetrieve(ctx, clusterGroupIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() { resp.State.RemoveResource(ctx) }) {
			return
		}
		resp.Diagnostics.AddError("Error reading cluster group", utils.FormatAPIError(fmt.Sprintf("read cluster group ID %s", clusterGroupID), err, httpResp))
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read cluster group", httpResp, http.StatusOK) {
		return
	}

	// Store state values for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	r.mapClusterGroupToState(clusterGroup, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, clusterGroup.HasTags(), clusterGroup.GetTags(), stateTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, clusterGroup.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ClusterGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	clusterGroupID := plan.ID.ValueString()
	clusterGroupIDInt := utils.ParseInt32FromString(clusterGroupID)
	if clusterGroupIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Cluster Group ID", fmt.Sprintf("Cluster Group ID must be a number, got: %s", clusterGroupID))
		return
	}
	tflog.Debug(ctx, "Updating cluster group", map[string]interface{}{
		"id":   clusterGroupID,
		"name": plan.Name.ValueString(),
	})

	// Build the request
	clusterGroupRequest := netbox.ClusterGroupRequest{
		Name: plan.Name.ValueString(),
		Slug: plan.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&clusterGroupRequest, plan.Description)

	// Apply tags and custom_fields with merge-aware helpers
	utils.ApplyTagsFromSlugs(ctx, r.client, &clusterGroupRequest, plan.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFieldsWithMerge(ctx, &clusterGroupRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)

	// Update via API
	clusterGroup, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterGroupsUpdate(ctx, clusterGroupIDInt).ClusterGroupRequest(clusterGroupRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating cluster group", utils.FormatAPIError(fmt.Sprintf("update cluster group ID %s", clusterGroupID), err, httpResp))
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "update cluster group", httpResp, http.StatusOK) {
		return
	}
	r.mapClusterGroupToState(clusterGroup, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	plan.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, clusterGroup.HasTags(), clusterGroup.GetTags(), plan.Tags)
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, plan.CustomFields, clusterGroup.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClusterGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	clusterGroupID := data.ID.ValueString()
	clusterGroupIDInt := utils.ParseInt32FromString(clusterGroupID)
	if clusterGroupIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Cluster Group ID", fmt.Sprintf("Cluster Group ID must be a number, got: %s", clusterGroupID))
		return
	}
	tflog.Debug(ctx, "Deleting cluster group", map[string]interface{}{"id": clusterGroupID})
	httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterGroupsDestroy(ctx, clusterGroupIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() {}) {
			return
		}
		resp.Diagnostics.AddError("Error deleting cluster group", utils.FormatAPIError(fmt.Sprintf("delete cluster group ID %s", clusterGroupID), err, httpResp))
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "delete cluster group", httpResp, http.StatusNoContent) {
		return
	}
	tflog.Trace(ctx, "deleted a cluster group resource")
}

func (r *ClusterGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		clusterGroupIDInt := utils.ParseInt32FromString(parsed.ID)
		if clusterGroupIDInt == 0 {
			resp.Diagnostics.AddError("Invalid Cluster Group ID", fmt.Sprintf("Cluster Group ID must be a number, got: %s", parsed.ID))
			return
		}

		clusterGroup, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterGroupsRetrieve(ctx, clusterGroupIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing cluster group", utils.FormatAPIError(fmt.Sprintf("read cluster group ID %s", parsed.ID), err, httpResp))
			return
		}
		if !utils.ValidateStatusCode(&resp.Diagnostics, "import cluster group", httpResp, http.StatusOK) {
			return
		}

		var data ClusterGroupResourceModel
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, clusterGroup.HasTags(), clusterGroup.GetTags(), data.Tags)
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapClusterGroupToState(clusterGroup, &data)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, clusterGroup.HasTags(), clusterGroup.GetTags(), data.Tags)

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, clusterGroup.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapClusterGroupToState maps API response to Terraform state.
func (r *ClusterGroupResource) mapClusterGroupToState(clusterGroup *netbox.ClusterGroup, data *ClusterGroupResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", clusterGroup.GetId()))
	data.Name = types.StringValue(clusterGroup.GetName())
	data.Slug = types.StringValue(clusterGroup.GetSlug())
	data.Description = utils.StringFromAPI(clusterGroup.HasDescription(), clusterGroup.GetDescription, data.Description)

	// Tags and custom fields are now handled in Create/Read/Update methods using filter-to-owned pattern
}
