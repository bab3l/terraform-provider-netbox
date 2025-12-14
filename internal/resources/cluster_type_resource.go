// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
			"id":            nbschema.IDAttribute("cluster type"),
			"name":          nbschema.NameAttribute("cluster type", 100),
			"slug":          nbschema.SlugAttribute("cluster type"),
			"description":   nbschema.DescriptionAttribute("cluster type"),
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
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
	if clusterType.HasDescription() && clusterType.GetDescription() != "" {
		data.Description = types.StringValue(clusterType.GetDescription())
	} else if !data.Description.IsNull() {
		data.Description = types.StringNull()
	}

	// Handle tags
	if clusterType.HasTags() {
		tags := utils.NestedTagsToTagModels(clusterType.GetTags())
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
	if clusterType.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(clusterType.GetCustomFields(), stateCustomFields)
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
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		clusterTypeRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		clusterTypeRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		clusterTypeRequest.CustomFields = utils.CustomFieldsToMap(customFields)
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

	// Map response to state
	r.mapClusterTypeToState(ctx, clusterType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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
		if httpResp != nil && httpResp.StatusCode == 404 {
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

	// Map response to state
	r.mapClusterTypeToState(ctx, clusterType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ClusterTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
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

	tflog.Debug(ctx, "Updating cluster type", map[string]interface{}{
		"id":   clusterTypeID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the cluster type request
	clusterTypeRequest := netbox.ClusterTypeRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		clusterTypeRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		clusterTypeRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		clusterTypeRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

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

	// Map response to state
	r.mapClusterTypeToState(ctx, clusterType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		if httpResp != nil && httpResp.StatusCode == 404 {
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
