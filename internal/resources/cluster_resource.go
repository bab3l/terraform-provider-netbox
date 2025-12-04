// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ClusterResource{}
	_ resource.ResourceWithConfigure   = &ClusterResource{}
	_ resource.ResourceWithImportState = &ClusterResource{}
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
			"description": nbschema.DescriptionAttribute("cluster"),
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the cluster.",
				Optional:            true,
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
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
func (r *ClusterResource) mapClusterToState(ctx context.Context, cluster *netbox.Cluster, data *ClusterResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", cluster.GetId()))
	data.Name = types.StringValue(cluster.GetName())

	// Type (always present - required field)
	// Preserve the user's input format (name or slug) to avoid state drift
	clusterTypeName := cluster.Type.GetName()
	clusterTypeSlug := cluster.Type.GetSlug()
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		configuredValue := data.Type.ValueString()
		if configuredValue == clusterTypeSlug {
			data.Type = types.StringValue(clusterTypeSlug)
		} else {
			data.Type = types.StringValue(clusterTypeName)
		}
	} else {
		data.Type = types.StringValue(clusterTypeName)
	}

	// Group
	if cluster.Group.IsSet() && cluster.Group.Get() != nil {
		data.Group = types.StringValue(cluster.Group.Get().GetName())
	} else {
		data.Group = types.StringNull()
	}

	// Status
	if cluster.HasStatus() {
		data.Status = types.StringValue(string(cluster.Status.GetValue()))
	} else {
		data.Status = types.StringValue("active")
	}

	// Tenant
	if cluster.Tenant.IsSet() && cluster.Tenant.Get() != nil {
		data.Tenant = types.StringValue(cluster.Tenant.Get().GetName())
	} else {
		data.Tenant = types.StringNull()
	}

	// Site
	if cluster.Site.IsSet() && cluster.Site.Get() != nil {
		data.Site = types.StringValue(cluster.Site.Get().GetName())
	} else {
		data.Site = types.StringNull()
	}

	// Description
	if cluster.HasDescription() && cluster.GetDescription() != "" {
		data.Description = types.StringValue(cluster.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if cluster.HasComments() && cluster.GetComments() != "" {
		data.Comments = types.StringValue(cluster.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if cluster.HasTags() {
		tags := utils.NestedTagsToTagModels(cluster.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if cluster.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(cluster.GetCustomFields(), stateCustomFields)
		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfValueDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}

// buildClusterRequest builds a WritableClusterRequest from the resource model.
func (r *ClusterResource) buildClusterRequest(ctx context.Context, data *ClusterResourceModel, diags *diag.Diagnostics) *netbox.WritableClusterRequest {
	// Lookup cluster type (required)
	clusterType, typeDiags := netboxlookup.LookupClusterType(ctx, r.client, data.Type.ValueString())
	diags.Append(typeDiags...)
	if diags.HasError() {
		return nil
	}

	clusterRequest := &netbox.WritableClusterRequest{
		Name: data.Name.ValueString(),
		Type: *clusterType,
	}

	// Group
	if utils.IsSet(data.Group) {
		group, groupDiags := netboxlookup.LookupClusterGroup(ctx, r.client, data.Group.ValueString())
		diags.Append(groupDiags...)
		if diags.HasError() {
			return nil
		}
		clusterRequest.Group = *netbox.NewNullableBriefClusterGroupRequest(group)
	}

	// Status
	if utils.IsSet(data.Status) {
		status := netbox.ClusterStatusValue(data.Status.ValueString())
		clusterRequest.Status = &status
	}

	// Tenant
	if utils.IsSet(data.Tenant) {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return nil
		}
		clusterRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	}

	// Site
	if utils.IsSet(data.Site) {
		site, siteDiags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())
		diags.Append(siteDiags...)
		if diags.HasError() {
			return nil
		}
		clusterRequest.Site = *netbox.NewNullableBriefSiteRequest(site)
	}

	// Description
	if utils.IsSet(data.Description) {
		description := data.Description.ValueString()
		clusterRequest.Description = &description
	}

	// Comments
	if utils.IsSet(data.Comments) {
		comments := data.Comments.ValueString()
		clusterRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		diags.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if diags.HasError() {
			return nil
		}
		clusterRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if diags.HasError() {
			return nil
		}
		clusterRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

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

	// Map response to state
	r.mapClusterToState(ctx, cluster, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
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

	// Map response to state
	r.mapClusterToState(ctx, cluster, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
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

	tflog.Debug(ctx, "Updating cluster", map[string]interface{}{
		"id":   clusterID,
		"name": data.Name.ValueString(),
	})

	// Build the cluster request
	clusterRequest := r.buildClusterRequest(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	cluster, httpResp, err := r.client.VirtualizationAPI.VirtualizationClustersUpdate(ctx, clusterIDInt).WritableClusterRequest(*clusterRequest).Execute()
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

	// Map response to state
	r.mapClusterToState(ctx, cluster, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
