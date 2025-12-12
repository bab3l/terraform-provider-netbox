// Package resources provides Terraform resource implementations for NetBox objects.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &WirelessLANGroupResource{}
	_ resource.ResourceWithConfigure   = &WirelessLANGroupResource{}
	_ resource.ResourceWithImportState = &WirelessLANGroupResource{}
)

// NewWirelessLANGroupResource returns a new resource implementing the wireless LAN group resource.
func NewWirelessLANGroupResource() resource.Resource {
	return &WirelessLANGroupResource{}
}

// WirelessLANGroupResource defines the resource implementation.
type WirelessLANGroupResource struct {
	client *netbox.APIClient
}

// WirelessLANGroupResourceModel describes the resource data model.
type WirelessLANGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Parent       types.String `tfsdk:"parent"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *WirelessLANGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_lan_group"
}

// Schema defines the schema for the resource.
func (r *WirelessLANGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a wireless LAN group in NetBox. Wireless LAN groups organize wireless networks hierarchically.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the wireless LAN group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the wireless LAN group.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "A unique slug identifier for the wireless LAN group.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the wireless LAN group.",
				Optional:            true,
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "Parent wireless LAN group (ID or slug) for hierarchical organization.",
				Optional:            true,
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *WirelessLANGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.
func (r *WirelessLANGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request - use WritableWirelessLANGroupRequest
	apiReq := netbox.NewWritableWirelessLANGroupRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID, err := utils.ParseID(data.Parent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent must be a numeric ID, got: %s", data.Parent.ValueString()),
			)
			return
		}
		apiReq.SetParent(parentID)
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTags(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var cfModels []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))
	}

	tflog.Debug(ctx, "Creating wireless LAN group", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsCreate(ctx).WritableWirelessLANGroupRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating wireless LAN group",
			utils.FormatAPIError(fmt.Sprintf("create wireless LAN group %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created wireless LAN group", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *WirelessLANGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Wireless LAN Group ID",
			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Reading wireless LAN group", map[string]interface{}{
		"id": groupID,
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsRetrieve(ctx, groupID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading wireless LAN group",
			utils.FormatAPIError(fmt.Sprintf("read wireless LAN group ID %d", groupID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *WirelessLANGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Wireless LAN Group ID",
			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	// Build request - use WritableWirelessLANGroupRequest
	apiReq := netbox.NewWritableWirelessLANGroupRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID, err := utils.ParseID(data.Parent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent must be a numeric ID, got: %s", data.Parent.ValueString()),
			)
			return
		}
		apiReq.SetParent(parentID)
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTags(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var cfModels []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))
	}

	tflog.Debug(ctx, "Updating wireless LAN group", map[string]interface{}{
		"id": groupID,
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsUpdate(ctx, groupID).WritableWirelessLANGroupRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating wireless LAN group",
			utils.FormatAPIError(fmt.Sprintf("update wireless LAN group ID %d", groupID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.
func (r *WirelessLANGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Wireless LAN Group ID",
			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Deleting wireless LAN group", map[string]interface{}{
		"id": groupID,
	})

	httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsDestroy(ctx, groupID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting wireless LAN group",
			utils.FormatAPIError(fmt.Sprintf("delete wireless LAN group ID %d", groupID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *WirelessLANGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupID, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", req.ID),
		)
		return
	}

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsRetrieve(ctx, groupID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing wireless LAN group",
			utils.FormatAPIError(fmt.Sprintf("import wireless LAN group ID %d", groupID), err, httpResp),
		)
		return
	}

	var data WirelessLANGroupResourceModel
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *WirelessLANGroupResource) mapResponseToModel(ctx context.Context, group *netbox.WirelessLANGroup, data *WirelessLANGroupResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	data.Name = types.StringValue(group.GetName())
	data.Slug = types.StringValue(group.GetSlug())

	// Map description
	if desc, ok := group.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map parent
	if group.Parent.IsSet() && group.Parent.Get() != nil {
		data.Parent = types.StringValue(fmt.Sprintf("%d", group.Parent.Get().GetId()))
	} else {
		data.Parent = types.StringNull()
	}

	// Handle tags
	if group.HasTags() && len(group.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(group.GetTags())
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
	if group.HasCustomFields() {
		apiCustomFields := group.GetCustomFields()
		var stateCustomFieldModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
			data.CustomFields.ElementsAs(ctx, &stateCustomFieldModels, false)
		}
		customFields := utils.MapToCustomFieldModels(apiCustomFields, stateCustomFieldModels)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
