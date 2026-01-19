// Package resources contains Terraform resource implementations for the Netbox provider.

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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &ClusterTypeResource{}
	_ resource.ResourceWithConfigure   = &ClusterTypeResource{}
	_ resource.ResourceWithImportState = &ClusterTypeResource{}
	_ resource.ResourceWithIdentity    = &ClusterTypeResource{}
)

// NewClusterTypeResource returns a new Cluster Type resource.
func NewClusterTypeResource() resource.Resource {
	return &ClusterTypeResource{}
}

// ClusterTypeResource defines the resource implementation.
type ClusterTypeResource struct {
	client *netbox.APIClient
}

// ClusterTypeResourceModel describes the resource data model.
type ClusterTypeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ClusterTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_type"
}

// Schema defines the schema for the resource.
func (r *ClusterTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a cluster type in Netbox. Cluster types define the technology or platform used for virtualization clusters (e.g., 'VMware vSphere', 'Proxmox', 'Kubernetes').",
		Attributes: map[string]schema.Attribute{
			"id":   nbschema.IDAttribute("cluster type"),
			"name": nbschema.NameAttribute("cluster type", 100),
			"slug": nbschema.SlugAttribute("cluster type"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("cluster type"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ClusterTypeResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure sets up the resource with the provider client.
func (r *ClusterTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// mapClusterTypeToState maps a ClusterType from the API to the Terraform state model.
func (r *ClusterTypeResource) mapClusterTypeToState(ctx context.Context, clusterType *netbox.ClusterType, data *ClusterTypeResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", clusterType.GetId()))
	data.Name = types.StringValue(clusterType.GetName())
	data.Slug = types.StringValue(clusterType.GetSlug())

	// Handle description
	data.Description = utils.NullableStringFromAPI(
		clusterType.HasDescription() && clusterType.GetDescription() != "",
		clusterType.GetDescription,
		data.Description,
	)

	// Tags and custom fields are now handled in Create/Read/Update with filter-to-owned pattern
}

// Create creates a new cluster type resource.
func (r *ClusterTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating cluster type", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the cluster type request
	clusterTypeRequest := netbox.ClusterTypeRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	utils.ApplyDescription(&clusterTypeRequest, data.Description)

	// Handle tags and custom_fields
	utils.ApplyTagsFromSlugs(ctx, r.client, &clusterTypeRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, &clusterTypeRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	clusterType, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterTypesCreate(ctx).ClusterTypeRequest(clusterTypeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster type",
			utils.FormatAPIError("create cluster type", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Created cluster type", map[string]interface{}{
		"id":   clusterType.GetId(),
		"name": clusterType.GetName(),
	})

	// Save plan state for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Map response to state
	r.mapClusterTypeToState(ctx, clusterType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case clusterType.HasTags() && len(clusterType.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(clusterType.GetTags()))
		for _, tag := range clusterType.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, clusterType.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ClusterTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	clusterTypeID := data.ID.ValueString()
	var clusterTypeIDInt int32
	clusterTypeIDInt, err := utils.ParseID(clusterTypeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cluster Type ID",
			fmt.Sprintf("Cluster Type ID must be a number, got: %s", clusterTypeID),
		)
		return
	}
	tflog.Debug(ctx, "Reading cluster type", map[string]interface{}{
		"id": clusterTypeID,
	})

	// Call the API
	clusterType, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterTypesRetrieve(ctx, clusterTypeIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Cluster type not found, removing from state", map[string]interface{}{
				"id": clusterTypeID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster type",
			utils.FormatAPIError(fmt.Sprintf("read cluster type ID %s", clusterTypeID), err, httpResp),
		)
		return
	}

	// Save state for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Map response to state
	r.mapClusterTypeToState(ctx, clusterType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve null/empty state values for tags and custom_fields
	wasExplicitlyEmpty := !stateTags.IsNull() && !stateTags.IsUnknown() && len(stateTags.Elements()) == 0
	switch {
	case clusterType.HasTags() && len(clusterType.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(clusterType.GetTags()))
		for _, tag := range clusterType.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, clusterType.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ClusterTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read both state and plan for merge-aware custom fields
	var state, plan ClusterTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	clusterTypeID := plan.ID.ValueString()
	var clusterTypeIDInt int32
	clusterTypeIDInt, err := utils.ParseID(clusterTypeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cluster Type ID",
			fmt.Sprintf("Cluster Type ID must be a number, got: %s", clusterTypeID),
		)
		return
	}
	tflog.Debug(ctx, "Updating cluster type", map[string]interface{}{
		"id":   clusterTypeID,
		"name": plan.Name.ValueString(),
		"slug": plan.Slug.ValueString(),
	})

	// Build the cluster type request
	clusterTypeRequest := netbox.ClusterTypeRequest{
		Name: plan.Name.ValueString(),
		Slug: plan.Slug.ValueString(),
	}

	// Set optional fields if provided
	utils.ApplyDescription(&clusterTypeRequest, plan.Description)

	// Handle tags and custom_fields with merge-aware helpers
	utils.ApplyTagsFromSlugs(ctx, r.client, &clusterTypeRequest, plan.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFieldsWithMerge(ctx, &clusterTypeRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)

	// Call the API
	clusterType, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterTypesUpdate(ctx, clusterTypeIDInt).ClusterTypeRequest(clusterTypeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster type",
			utils.FormatAPIError(fmt.Sprintf("update cluster type ID %s", clusterTypeID), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Updated cluster type", map[string]interface{}{
		"id":   clusterType.GetId(),
		"name": clusterType.GetName(),
	})

	// Save plan state for filter-to-owned pattern
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Map response to state
	r.mapClusterTypeToState(ctx, clusterType, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case clusterType.HasTags() && len(clusterType.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(clusterType.GetTags()))
		for _, tag := range clusterType.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		plan.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		plan.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		plan.Tags = types.SetNull(types.StringType)
	}
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, clusterType.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *ClusterTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	clusterTypeID := data.ID.ValueString()
	var clusterTypeIDInt int32
	clusterTypeIDInt, err := utils.ParseID(clusterTypeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cluster Type ID",
			fmt.Sprintf("Cluster Type ID must be a number, got: %s", clusterTypeID),
		)
		return
	}
	tflog.Debug(ctx, "Deleting cluster type", map[string]interface{}{
		"id": clusterTypeID,
	})

	// Call the API
	httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterTypesDestroy(ctx, clusterTypeIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted, consider success
			tflog.Debug(ctx, "Cluster type already deleted", map[string]interface{}{
				"id": clusterTypeID,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting cluster type",
			utils.FormatAPIError(fmt.Sprintf("delete cluster type ID %s", clusterTypeID), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted cluster type", map[string]interface{}{
		"id": clusterTypeID,
	})
}

// ImportState imports an existing resource into Terraform.
func (r *ClusterTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		id, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid cluster type ID",
				fmt.Sprintf("Cluster type ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		clusterType, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterTypesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing cluster type",
				utils.FormatAPIError("read cluster type", err, httpResp),
			)
			return
		}

		var data ClusterTypeResourceModel
		if clusterType.HasTags() {
			var tagSlugs []string
			for _, tag := range clusterType.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		}
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

		r.mapClusterTypeToState(ctx, clusterType, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, clusterType.GetCustomFields(), &resp.Diagnostics)
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
