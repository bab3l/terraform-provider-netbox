// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
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

var _ resource.Resource = &RegionResource{}

var _ resource.ResourceWithImportState = &RegionResource{}

func NewRegionResource() resource.Resource {

	return &RegionResource{}

}

// RegionResource defines the resource implementation.

type RegionResource struct {
	client *netbox.APIClient
}

// RegionResourceModel describes the resource data model.

type RegionResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	DisplayName types.String `tfsdk:"display_name"`

	Parent types.String `tfsdk:"parent"`

	ParentID types.String `tfsdk:"parent_id"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *RegionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_region"

}

func (r *RegionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a region in Netbox. Regions provide a hierarchical way to organize sites geographically, such as continents, countries, states, or cities.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("region"),

			"name": nbschema.NameAttribute("region", 100),

			"slug": nbschema.SlugAttribute("region"),

			"display_name": nbschema.DisplayNameAttribute("region"),

			"parent": nbschema.ReferenceAttribute("parent region", "ID or slug of the parent region. Leave empty for top-level regions. This enables hierarchical organization of geographic areas."),

			"parent_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the parent region.",
			},

			"description": nbschema.DescriptionAttribute("region"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *RegionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *RegionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data RegionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating region", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Prepare the region request

	regionRequest := netbox.WritableRegionRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string field

	regionRequest.Description = utils.StringPtr(data.Description)

	// Set optional parent

	if utils.IsSet(data.Parent) {

		parentID, parentDiags := netboxlookup.LookupRegionID(ctx, r.client, data.Parent.ValueString())

		resp.Diagnostics.Append(parentDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		regionRequest.Parent = *netbox.NewNullableInt32(&parentID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		regionRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFieldModels []utils.CustomFieldModel

		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		regionRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)

	}

	// Call the API

	region, httpResp, err := r.client.DcimAPI.DcimRegionsCreate(ctx).WritableRegionRequest(regionRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error creating region", utils.FormatAPIError("create region", err, httpResp))

		return

	}

	if httpResp.StatusCode != 201 {

		resp.Diagnostics.AddError("Error creating region", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state

	r.mapRegionToState(ctx, region, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "created a region resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RegionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data RegionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	regionID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading region", map[string]interface{}{"id": regionID})

	var regionIDInt int32

	regionIDInt, err := utils.ParseID(regionID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Region ID", fmt.Sprintf("Region ID must be a number, got: %s", regionID))

		return

	}

	region, httpResp, err := r.client.DcimAPI.DcimRegionsRetrieve(ctx, regionIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error reading region", utils.FormatAPIError(fmt.Sprintf("read region ID %s", regionID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error reading region", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state

	r.mapRegionToState(ctx, region, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RegionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data RegionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	regionID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating region", map[string]interface{}{"id": regionID})

	var regionIDInt int32

	regionIDInt, err := utils.ParseID(regionID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Region ID", fmt.Sprintf("Region ID must be a number, got: %s", regionID))

		return

	}

	// Prepare the region request

	regionRequest := netbox.WritableRegionRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string field

	regionRequest.Description = utils.StringPtr(data.Description)

	// Set optional parent

	if utils.IsSet(data.Parent) {

		parentID, parentDiags := netboxlookup.LookupRegionID(ctx, r.client, data.Parent.ValueString())

		resp.Diagnostics.Append(parentDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		regionRequest.Parent = *netbox.NewNullableInt32(&parentID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		regionRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFieldModels []utils.CustomFieldModel

		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		regionRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)

	}

	// Call the API

	region, httpResp, err := r.client.DcimAPI.DcimRegionsUpdate(ctx, regionIDInt).WritableRegionRequest(regionRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error updating region", utils.FormatAPIError(fmt.Sprintf("update region ID %s", regionID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error updating region", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state

	r.mapRegionToState(ctx, region, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "updated a region resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RegionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data RegionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	regionID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting region", map[string]interface{}{"id": regionID})

	var regionIDInt int32

	regionIDInt, err := utils.ParseID(regionID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Region ID", fmt.Sprintf("Region ID must be a number, got: %s", regionID))

		return

	}

	httpResp, err := r.client.DcimAPI.DcimRegionsDestroy(ctx, regionIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error deleting region", utils.FormatAPIError(fmt.Sprintf("delete region ID %s", regionID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 204 {

		resp.Diagnostics.AddError("Error deleting region", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return

	}

	tflog.Trace(ctx, "deleted a region resource")

}

func (r *RegionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapRegionToState maps API response to Terraform state.

func (r *RegionResource) mapRegionToState(ctx context.Context, region *netbox.Region, data *RegionResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", region.GetId()))

	data.Name = types.StringValue(region.GetName())

	data.Slug = types.StringValue(region.GetSlug())

	data.DisplayName = types.StringValue(region.GetDisplay())

	// Handle parent

	if region.HasParent() && region.GetParent().Id != 0 {

		parent := region.GetParent()

		data.ParentID = types.StringValue(fmt.Sprintf("%d", parent.GetId()))

		userParent := data.Parent.ValueString()

		if userParent == parent.GetName() || userParent == parent.GetSlug() || userParent == parent.GetDisplay() || userParent == fmt.Sprintf("%d", parent.GetId()) {

			// Keep user's original value

		} else {

			data.Parent = types.StringValue(parent.GetName())

		}

	} else {

		data.Parent = types.StringNull()
		data.ParentID = types.StringNull()

	}

	// Handle description - use StringFromAPI to treat empty string as null

	data.Description = utils.StringFromAPI(region.HasDescription(), region.GetDescription, data.Description)

	// Handle tags

	if region.HasTags() {

		tags := utils.NestedTagsToTagModels(region.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if !diags.HasError() {

			data.Tags = tagsValue

		}

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	if region.HasCustomFields() {

		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {

			cfDiags := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(cfDiags...)

		}

		if !diags.HasError() {

			customFields := utils.MapToCustomFieldModels(region.GetCustomFields(), existingModels)

			customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

			diags.Append(cfDiags...)

			if !cfDiags.HasError() {

				data.CustomFields = customFieldsValue

			}

		}

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
