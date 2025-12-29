package resources

import (
	"context"
	"fmt"
	"maps"

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

var _ resource.Resource = &ClusterGroupResource{}
var _ resource.ResourceWithImportState = &ClusterGroupResource{}

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
			"id":           nbschema.IDAttribute("cluster group"),
			"name":         nbschema.NameAttribute("cluster group", 100),
			"slug":         nbschema.SlugAttribute("cluster group"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("cluster group"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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

	// Handle tags and custom_fields
	utils.ApplyMetadataFields(ctx, &clusterGroupRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

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

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError("Error creating cluster group", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}

	r.mapClusterGroupToState(ctx, clusterGroup, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
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
		resp.Diagnostics.AddError("Error reading cluster group", utils.FormatAPIError(fmt.Sprintf("read cluster group ID %s", clusterGroupID), err, httpResp))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading cluster group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	r.mapClusterGroupToState(ctx, clusterGroup, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterGroupID := data.ID.ValueString()

	clusterGroupIDInt := utils.ParseInt32FromString(clusterGroupID)

	if clusterGroupIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Cluster Group ID", fmt.Sprintf("Cluster Group ID must be a number, got: %s", clusterGroupID))

		return
	}

	tflog.Debug(ctx, "Updating cluster group", map[string]interface{}{
		"id": clusterGroupID,

		"name": data.Name.ValueString(),
	})

	// Build the request

	clusterGroupRequest := netbox.ClusterGroupRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&clusterGroupRequest, data.Description)

	// Handle tags and custom_fields
	utils.ApplyMetadataFields(ctx, &clusterGroupRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	// Update via API

	clusterGroup, httpResp, err := r.client.VirtualizationAPI.VirtualizationClusterGroupsUpdate(ctx, clusterGroupIDInt).ClusterGroupRequest(clusterGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError("Error updating cluster group", utils.FormatAPIError(fmt.Sprintf("update cluster group ID %s", clusterGroupID), err, httpResp))

		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating cluster group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return
	}

	r.mapClusterGroupToState(ctx, clusterGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		if httpResp != nil && httpResp.StatusCode == 404 {
			return // Already deleted
		}

		resp.Diagnostics.AddError("Error deleting cluster group", utils.FormatAPIError(fmt.Sprintf("delete cluster group ID %s", clusterGroupID), err, httpResp))

		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error deleting cluster group", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return
	}

	tflog.Trace(ctx, "deleted a cluster group resource")
}

func (r *ClusterGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapClusterGroupToState maps API response to Terraform state.

func (r *ClusterGroupResource) mapClusterGroupToState(ctx context.Context, clusterGroup *netbox.ClusterGroup, data *ClusterGroupResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", clusterGroup.GetId()))

	data.Name = types.StringValue(clusterGroup.GetName())

	// DisplayName
	if clusterGroup.Display != "" {
	} else {
	}

	data.Slug = types.StringValue(clusterGroup.GetSlug())

	data.Description = utils.StringFromAPI(clusterGroup.HasDescription(), clusterGroup.GetDescription, data.Description)

	// Handle tags
	data.Tags = utils.PopulateTagsFromNestedTags(ctx, clusterGroup.HasTags(), clusterGroup.GetTags(), diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields
	data.CustomFields = utils.PopulateCustomFieldsFromMap(ctx, clusterGroup.HasCustomFields(), clusterGroup.GetCustomFields(), data.CustomFields, diags)
}
