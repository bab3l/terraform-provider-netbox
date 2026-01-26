// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &ClusterResource{}
	_ resource.ResourceWithConfigure   = &ClusterResource{}
	_ resource.ResourceWithImportState = &ClusterResource{}
	_ resource.ResourceWithIdentity    = &ClusterResource{}
)

// NewClusterResource returns a new Cluster resource.
func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the resource implementation.
type ClusterResource struct {
	client *netbox.APIClient
}

// ClusterResourceModel describes the resource data model.
type ClusterResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Group        types.String `tfsdk:"group"`
	Status       types.String `tfsdk:"status"`
	Tenant       types.String `tfsdk:"tenant"`
	Site         types.String `tfsdk:"site"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the resource.
func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtualization cluster in Netbox. Clusters represent a pool of physical resources (such as compute, storage, and networking) that can be used to run virtual machines.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the cluster.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": nbschema.NameAttribute("cluster", 100),
			"type": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the cluster type (e.g., 'VMware vSphere', 'Proxmox').",
				Required:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the cluster group this cluster belongs to.",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the cluster. Valid values are: `planned`, `staging`, `active`, `decommissioning`, `offline`. Defaults to `active`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the tenant this cluster is assigned to.",
				Optional:            true,
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the site where this cluster is located.",
				Optional:            true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("cluster"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ClusterResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure sets up the resource with the provider client.
func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// mapClusterToState maps a Cluster from the API to the Terraform state model.
func (r *ClusterResource) mapClusterToState(cluster *netbox.Cluster, data *ClusterResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", cluster.GetId()))
	data.Name = types.StringValue(cluster.GetName())

	// Type (required field) - preserve user's input format (ID, name, or slug)
	data.Type = utils.PreserveReferenceFormat(data.Type, cluster.Type.GetId(), cluster.Type.GetSlug(), "")

	// Group (optional)
	if cluster.Group.IsSet() && cluster.Group.Get() != nil {
		group := cluster.Group.Get()
		data.Group = utils.PreserveOptionalReferenceFormat(data.Group, true, group.GetId(), group.GetSlug(), "")
	} else {
		data.Group = types.StringNull()
	}

	// Status
	if cluster.HasStatus() {
		data.Status = types.StringValue(string(cluster.Status.GetValue()))
	} else {
		data.Status = types.StringValue("active")
	}

	// Tenant (optional)
	if cluster.Tenant.IsSet() && cluster.Tenant.Get() != nil {
		tenant := cluster.Tenant.Get()
		data.Tenant = utils.PreserveOptionalReferenceFormat(data.Tenant, true, tenant.GetId(), tenant.GetName(), tenant.GetSlug())
	} else {
		data.Tenant = types.StringNull()
	}

	// Site (optional)
	if cluster.Site.IsSet() && cluster.Site.Get() != nil {
		site := cluster.Site.Get()
		data.Site = utils.PreserveOptionalReferenceFormat(data.Site, true, site.GetId(), site.GetName(), site.GetSlug())
	} else {
		data.Site = types.StringNull()
	}

	// Description
	data.Description = utils.NullableStringFromAPI(
		cluster.HasDescription() && cluster.GetDescription() != "",
		cluster.GetDescription,
		data.Description,
	)

	// Comments
	data.Comments = utils.NullableStringFromAPI(
		cluster.HasComments() && cluster.GetComments() != "",
		cluster.GetComments,
		data.Comments,
	)

	// Tags and custom fields are now handled in Create/Read/Update with filter-to-owned pattern
}

// buildClusterRequest builds a WritableClusterRequest from the resource model.
func (r *ClusterResource) buildClusterRequest(ctx context.Context, data *ClusterResourceModel, diags *diag.Diagnostics) *netbox.WritableClusterRequest {
	// Lookup cluster type (required)
	clusterType := utils.ResolveRequiredReference(ctx, r.client, data.Type, netboxlookup.LookupClusterType, diags)
	if diags.HasError() {
		return nil
	}
	clusterRequest := &netbox.WritableClusterRequest{
		Name: data.Name.ValueString(),
		Type: *clusterType,
	}

	// Group
	if group := utils.ResolveOptionalReference(ctx, r.client, data.Group, netboxlookup.LookupClusterGroup, diags); group != nil {
		clusterRequest.Group = *netbox.NewNullableBriefClusterGroupRequest(group)
	} else if data.Group.IsNull() {
		clusterRequest.SetGroupNil()
	}

	// Status
	if utils.IsSet(data.Status) {
		status := netbox.ClusterStatusValue(data.Status.ValueString())
		clusterRequest.Status = &status
	}

	// Tenant
	if tenant := utils.ResolveOptionalReference(ctx, r.client, data.Tenant, netboxlookup.LookupTenant, diags); tenant != nil {
		clusterRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	} else if data.Tenant.IsNull() {
		clusterRequest.SetTenantNil()
	}

	// Site
	if site := utils.ResolveOptionalReference(ctx, r.client, data.Site, netboxlookup.LookupSite, diags); site != nil {
		clusterRequest.Site = *netbox.NewNullableBriefSiteRequest(site)
	} else if data.Site.IsNull() {
		clusterRequest.SetSiteNil()
	}

	// Apply description and comments
	utils.ApplyDescriptiveFields(clusterRequest, data.Description, data.Comments)

	// Apply tags
	utils.ApplyTagsFromSlugs(ctx, r.client, clusterRequest, data.Tags, diags)
	if diags.HasError() {
		return nil
	}

	// Apply custom fields
	utils.ApplyCustomFields(ctx, clusterRequest, data.CustomFields, diags)
	return clusterRequest
}

// buildClusterRequestWithState builds a WritableClusterRequest with merge-aware custom fields.
func (r *ClusterResource) buildClusterRequestWithState(ctx context.Context, plan *ClusterResourceModel, state *ClusterResourceModel, diags *diag.Diagnostics) *netbox.WritableClusterRequest {
	// Lookup cluster type (required)
	clusterType := utils.ResolveRequiredReference(ctx, r.client, plan.Type, netboxlookup.LookupClusterType, diags)
	if diags.HasError() {
		return nil
	}
	clusterRequest := &netbox.WritableClusterRequest{
		Name: plan.Name.ValueString(),
		Type: *clusterType,
	}

	// Group
	if group := utils.ResolveOptionalReference(ctx, r.client, plan.Group, netboxlookup.LookupClusterGroup, diags); group != nil {
		clusterRequest.Group = *netbox.NewNullableBriefClusterGroupRequest(group)
	} else if plan.Group.IsNull() {
		clusterRequest.SetGroupNil()
	}

	// Status
	if utils.IsSet(plan.Status) {
		status := netbox.ClusterStatusValue(plan.Status.ValueString())
		clusterRequest.Status = &status
	}

	// Tenant
	if tenant := utils.ResolveOptionalReference(ctx, r.client, plan.Tenant, netboxlookup.LookupTenant, diags); tenant != nil {
		clusterRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	} else if plan.Tenant.IsNull() {
		clusterRequest.SetTenantNil()
	}

	// Site
	if site := utils.ResolveOptionalReference(ctx, r.client, plan.Site, netboxlookup.LookupSite, diags); site != nil {
		clusterRequest.Site = *netbox.NewNullableBriefSiteRequest(site)
	} else if plan.Site.IsNull() {
		clusterRequest.SetSiteNil()
	}

	// Apply description and comments
	utils.ApplyDescriptiveFields(clusterRequest, plan.Description, plan.Comments)

	// Apply tags
	utils.ApplyTagsFromSlugs(ctx, r.client, clusterRequest, plan.Tags, diags)
	if diags.HasError() {
		return nil
	}

	// Apply custom fields with merge (merge-aware)
	utils.ApplyCustomFieldsWithMerge(ctx, clusterRequest, plan.CustomFields, state.CustomFields, diags)
	return clusterRequest
}

// Create creates a new cluster resource.
func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating cluster", map[string]interface{}{
		"name": data.Name.ValueString(),
		"type": data.Type.ValueString(),
	})

	// Build the cluster request
	clusterRequest := r.buildClusterRequest(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	cluster, httpResp, err := r.client.VirtualizationAPI.VirtualizationClustersCreate(ctx).WritableClusterRequest(*clusterRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			utils.FormatAPIError("create cluster", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Created cluster", map[string]interface{}{
		"id":   cluster.GetId(),
		"name": cluster.GetName(),
	})

	// Save plan state for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Map response to state
	r.mapClusterToState(cluster, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, cluster.HasTags(), cluster.GetTags(), planTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, cluster.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	clusterID := data.ID.ValueString()
	var clusterIDInt int32
	clusterIDInt, err := utils.ParseID(clusterID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cluster ID",
			fmt.Sprintf("Cluster ID must be a number, got: %s", clusterID),
		)
		return
	}
	tflog.Debug(ctx, "Reading cluster", map[string]interface{}{
		"id": clusterID,
	})

	// Call the API
	cluster, httpResp, err := r.client.VirtualizationAPI.VirtualizationClustersRetrieve(ctx, clusterIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Cluster not found, removing from state", map[string]interface{}{
				"id": clusterID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster",
			utils.FormatAPIError(fmt.Sprintf("read cluster ID %s", clusterID), err, httpResp),
		)
		return
	}

	// Save state for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Map response to state
	r.mapClusterToState(cluster, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve null/empty state values for tags and custom_fields
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, cluster.HasTags(), cluster.GetTags(), stateTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, cluster.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read both state and plan for merge-aware custom fields
	var state, plan ClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	clusterID := plan.ID.ValueString()
	var clusterIDInt int32
	clusterIDInt, err := utils.ParseID(clusterID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cluster ID",
			fmt.Sprintf("Cluster ID must be a number, got: %s", clusterID),
		)
		return
	}
	tflog.Debug(ctx, "Updating cluster", map[string]interface{}{
		"id":   clusterID,
		"name": plan.Name.ValueString(),
	})

	// Build the cluster request with state for merge-aware custom fields
	clusterRequest := r.buildClusterRequestWithState(ctx, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	cluster, httpResp, err := r.client.VirtualizationAPI.VirtualizationClustersUpdate(ctx, clusterIDInt).WritableClusterRequest(*clusterRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster",
			utils.FormatAPIError(fmt.Sprintf("update cluster ID %s", clusterID), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Updated cluster", map[string]interface{}{
		"id":   cluster.GetId(),
		"name": cluster.GetName(),
	})

	// Save plan state for filter-to-owned pattern
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Map response to state
	r.mapClusterToState(cluster, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	plan.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, cluster.HasTags(), cluster.GetTags(), planTags)
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, cluster.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	clusterID := data.ID.ValueString()
	var clusterIDInt int32
	clusterIDInt, err := utils.ParseID(clusterID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cluster ID",
			fmt.Sprintf("Cluster ID must be a number, got: %s", clusterID),
		)
		return
	}
	tflog.Debug(ctx, "Deleting cluster", map[string]interface{}{
		"id": clusterID,
	})

	// Call the API
	httpResp, err := r.client.VirtualizationAPI.VirtualizationClustersDestroy(ctx, clusterIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted, consider success
			tflog.Debug(ctx, "Cluster already deleted", map[string]interface{}{
				"id": clusterID,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			utils.FormatAPIError(fmt.Sprintf("delete cluster ID %s", clusterID), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted cluster", map[string]interface{}{
		"id": clusterID,
	})
}

// ImportState imports an existing resource into Terraform.
func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		clusterIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Cluster ID",
				fmt.Sprintf("Cluster ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		cluster, httpResp, err := r.client.VirtualizationAPI.VirtualizationClustersRetrieve(ctx, clusterIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing cluster",
				utils.FormatAPIError("read cluster", err, httpResp),
			)
			return
		}

		var data ClusterResourceModel
		data.Type = types.StringValue(fmt.Sprintf("%d", cluster.Type.GetId()))
		if cluster.Group.IsSet() && cluster.Group.Get() != nil {
			data.Group = types.StringValue(fmt.Sprintf("%d", cluster.Group.Get().GetId()))
		}
		if cluster.Tenant.IsSet() && cluster.Tenant.Get() != nil {
			data.Tenant = types.StringValue(fmt.Sprintf("%d", cluster.Tenant.Get().GetId()))
		}
		if cluster.Site.IsSet() && cluster.Site.Get() != nil {
			data.Site = types.StringValue(fmt.Sprintf("%d", cluster.Site.Get().GetId()))
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, cluster.HasTags(), cluster.GetTags(), data.Tags)
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

		r.mapClusterToState(cluster, &data)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, cluster.GetCustomFields(), &resp.Diagnostics)
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
