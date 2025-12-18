// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ModuleResource{}

	_ resource.ResourceWithConfigure = &ModuleResource{}

	_ resource.ResourceWithImportState = &ModuleResource{}
)

// NewModuleResource returns a new resource implementing the module resource.

func NewModuleResource() resource.Resource {

	return &ModuleResource{}

}

// ModuleResource defines the resource implementation.

type ModuleResource struct {
	client *netbox.APIClient
}

// ModuleResourceModel describes the resource data model.

type ModuleResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	ModuleBay types.Int32 `tfsdk:"module_bay"`

	ModuleType types.String `tfsdk:"module_type"`

	Status types.String `tfsdk:"status"`

	Serial types.String `tfsdk:"serial"`

	AssetTag types.String `tfsdk:"asset_tag"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ModuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_module"

}

// Schema defines the schema for the resource.

func (r *ModuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a module in NetBox. Modules are hardware components installed in module bays within devices. This is the recommended replacement for inventory items (deprecated in NetBox v4.3).",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the module.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The device this module is installed in (ID or name).",

				Required: true,
			},

			"module_bay": schema.Int32Attribute{

				MarkdownDescription: "The module bay ID where this module is installed.",

				Required: true,
			},

			"module_type": schema.StringAttribute{

				MarkdownDescription: "The module type (ID or model name).",

				Required: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "Operational status. Valid values: `offline`, `active`, `planned`, `staged`, `failed`, `decommissioning`.",

				Optional: true,

				Computed: true,
			},

			"serial": schema.StringAttribute{

				MarkdownDescription: "Serial number of the module.",

				Optional: true,
			},

			"asset_tag": schema.StringAttribute{

				MarkdownDescription: "A unique tag used to identify this module.",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the module.",

				Optional: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *ModuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ModuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ModuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Lookup device

	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Lookup module type

	moduleType, diags := lookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritableModuleRequest(*device, data.ModuleBay.ValueInt32(), *moduleType)

	// Set optional fields

	if !data.Status.IsNull() && !data.Status.IsUnknown() {

		status := netbox.ModuleStatusValue(data.Status.ValueString())

		apiReq.SetStatus(status)

	}

	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {

		apiReq.SetSerial(data.Serial.ValueString())

	}

	if !data.AssetTag.IsNull() && !data.AssetTag.IsUnknown() {

		apiReq.SetAssetTag(data.AssetTag.ValueString())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

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

	tflog.Debug(ctx, "Creating module", map[string]interface{}{

		"device": data.Device.ValueString(),

		"module_bay": data.ModuleBay.ValueInt32(),

		"module_type": data.ModuleType.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimModulesCreate(ctx).WritableModuleRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating module",

			utils.FormatAPIError("create module", err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created module", map[string]interface{}{

		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *ModuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ModuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	moduleID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Module ID",

			fmt.Sprintf("Module ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading module", map[string]interface{}{

		"id": moduleID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimModulesRetrieve(ctx, moduleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading module",

			utils.FormatAPIError(fmt.Sprintf("read module ID %d", moduleID), err, httpResp),
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

func (r *ModuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ModuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	moduleID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Module ID",

			fmt.Sprintf("Module ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	// Lookup device

	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Lookup module type

	moduleType, diags := lookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritableModuleRequest(*device, data.ModuleBay.ValueInt32(), *moduleType)

	// Set optional fields

	if !data.Status.IsNull() && !data.Status.IsUnknown() {

		status := netbox.ModuleStatusValue(data.Status.ValueString())

		apiReq.SetStatus(status)

	}

	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {

		apiReq.SetSerial(data.Serial.ValueString())

	}

	if !data.AssetTag.IsNull() && !data.AssetTag.IsUnknown() {

		apiReq.SetAssetTag(data.AssetTag.ValueString())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

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

	tflog.Debug(ctx, "Updating module", map[string]interface{}{

		"id": moduleID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimModulesUpdate(ctx, moduleID).WritableModuleRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating module",

			utils.FormatAPIError(fmt.Sprintf("update module ID %d", moduleID), err, httpResp),
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

func (r *ModuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ModuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	moduleID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Module ID",

			fmt.Sprintf("Module ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting module", map[string]interface{}{

		"id": moduleID,
	})

	httpResp, err := r.client.DcimAPI.DcimModulesDestroy(ctx, moduleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting module",

			utils.FormatAPIError(fmt.Sprintf("delete module ID %d", moduleID), err, httpResp),
		)

		return

	}

}

// ImportState imports an existing resource.

func (r *ModuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	moduleID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Module ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimModulesRetrieve(ctx, moduleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing module",

			utils.FormatAPIError(fmt.Sprintf("import module ID %d", moduleID), err, httpResp),
		)

		return

	}

	var data ModuleResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *ModuleResource) mapResponseToModel(ctx context.Context, module *netbox.Module, data *ModuleResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", module.GetId()))

	// Map device - preserve user's input format

	if device := module.GetDevice(); device.Id != 0 {

		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())

	}

	// Map module_bay

	moduleBay := module.GetModuleBay()

	data.ModuleBay = types.Int32Value(moduleBay.Id)

	// Map module_type - preserve user's input format

	if mt := module.GetModuleType(); mt.Id != 0 {

		data.ModuleType = utils.UpdateReferenceAttribute(data.ModuleType, mt.GetModel(), "", mt.GetId())

	}

	// Map status

	if module.Status != nil {

		data.Status = types.StringValue(string(module.Status.GetValue()))

	} else {

		data.Status = types.StringNull()

	}

	// Map serial

	if serial, ok := module.GetSerialOk(); ok && serial != nil && *serial != "" {

		data.Serial = types.StringValue(*serial)

	} else {

		data.Serial = types.StringNull()

	}

	// Map asset_tag

	if module.AssetTag.IsSet() && module.AssetTag.Get() != nil && *module.AssetTag.Get() != "" {

		data.AssetTag = types.StringValue(*module.AssetTag.Get())

	} else {

		data.AssetTag = types.StringNull()

	}

	// Map description

	if desc, ok := module.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map comments

	if comments, ok := module.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if module.HasTags() && len(module.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(module.GetTags())

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

	if module.HasCustomFields() {

		apiCustomFields := module.GetCustomFields()

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
