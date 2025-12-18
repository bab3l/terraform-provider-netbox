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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &FrontPortResource{}

	_ resource.ResourceWithConfigure = &FrontPortResource{}

	_ resource.ResourceWithImportState = &FrontPortResource{}
)

// NewFrontPortResource returns a new resource implementing the front port resource.

func NewFrontPortResource() resource.Resource {

	return &FrontPortResource{}

}

// FrontPortResource defines the resource implementation.

type FrontPortResource struct {
	client *netbox.APIClient
}

// FrontPortResourceModel describes the resource data model.

type FrontPortResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	Color types.String `tfsdk:"color"`

	RearPort types.String `tfsdk:"rear_port"`

	RearPortPosition types.Int32 `tfsdk:"rear_port_position"`

	Description types.String `tfsdk:"description"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *FrontPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_front_port"

}

// Schema defines the schema for the resource.

func (r *FrontPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a front port in NetBox. Front ports represent physical ports on the front of a device, typically used for patch panels and fiber distribution. They are mapped to rear ports.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the front port.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The device this front port belongs to (ID or name).",

				Required: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the front port.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the front port.",

				Optional: true,
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "The type of front port (e.g., `8p8c`, `8p6c`, `110-punch`, `bnc`, `f`, `n`, `mrj21`, `fc`, `lc`, `lc-pc`, `lc-upc`, `lc-apc`, `lsh`, `mpo`, `mtrj`, `sc`, `sc-pc`, `sc-upc`, `sc-apc`, `st`, `cs`, `sn`, `splice`, `other`).",

				Required: true,
			},

			"color": schema.StringAttribute{

				MarkdownDescription: "Color of the front port in hex format (e.g., `aa1409`).",

				Optional: true,
			},

			"rear_port": schema.StringAttribute{

				MarkdownDescription: "The rear port that this front port maps to (ID).",

				Required: true,
			},

			"rear_port_position": schema.Int32Attribute{

				MarkdownDescription: "Position on the rear port (1-1024). Default is 1.",

				Optional: true,

				Computed: true,

				Default: int32default.StaticInt32(1),
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the front port.",

				Optional: true,
			},

			"mark_connected": schema.BoolAttribute{

				MarkdownDescription: "Treat as if a cable is connected.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(false),
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *FrontPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *FrontPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data FrontPortResourceModel

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

	// Parse rear port ID

	rearPortID, err := utils.ParseID(data.RearPort.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Rear Port ID",

			fmt.Sprintf("Could not parse rear port ID %q: %s", data.RearPort.ValueString(), err),
		)

		return

	}

	// Build request

	apiReq := netbox.NewWritableFrontPortRequest(*device, data.Name.ValueString(), netbox.FrontPortTypeValue(data.Type.ValueString()), rearPortID)

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Color.IsNull() && !data.Color.IsUnknown() {

		apiReq.SetColor(data.Color.ValueString())

	}

	if !data.RearPortPosition.IsNull() && !data.RearPortPosition.IsUnknown() {

		apiReq.SetRearPortPosition(data.RearPortPosition.ValueInt32())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {

		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())

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

	tflog.Debug(ctx, "Creating front port", map[string]interface{}{

		"device": data.Device.ValueString(),

		"name": data.Name.ValueString(),

		"rear_port": rearPortID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsCreate(ctx).WritableFrontPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating front port",

			utils.FormatAPIError("create front port", err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created front port", map[string]interface{}{

		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read retrieves the resource.

func (r *FrontPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data FrontPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)

		return

	}

	tflog.Debug(ctx, "Reading front port", map[string]interface{}{

		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Front port not found, removing from state", map[string]interface{}{

				"id": portID,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading front port",

			utils.FormatAPIError(fmt.Sprintf("read front port ID %d", portID), err, httpResp),
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

func (r *FrontPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data FrontPortResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)

		return

	}

	// Lookup device

	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse rear port ID

	rearPortID, err := utils.ParseID(data.RearPort.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Rear Port ID",

			fmt.Sprintf("Could not parse rear port ID %q: %s", data.RearPort.ValueString(), err),
		)

		return

	}

	// Build request

	apiReq := netbox.NewWritableFrontPortRequest(*device, data.Name.ValueString(), netbox.FrontPortTypeValue(data.Type.ValueString()), rearPortID)

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Color.IsNull() && !data.Color.IsUnknown() {

		apiReq.SetColor(data.Color.ValueString())

	}

	if !data.RearPortPosition.IsNull() && !data.RearPortPosition.IsUnknown() {

		apiReq.SetRearPortPosition(data.RearPortPosition.ValueInt32())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {

		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())

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

	tflog.Debug(ctx, "Updating front port", map[string]interface{}{

		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsUpdate(ctx, portID).WritableFrontPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating front port",

			utils.FormatAPIError(fmt.Sprintf("update front port ID %d", portID), err, httpResp),
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

// Delete removes the resource.

func (r *FrontPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data FrontPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)

		return

	}

	tflog.Debug(ctx, "Deleting front port", map[string]interface{}{

		"id": portID,
	})

	httpResp, err := r.client.DcimAPI.DcimFrontPortsDestroy(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Resource already deleted

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting front port",

			utils.FormatAPIError(fmt.Sprintf("delete front port ID %d", portID), err, httpResp),
		)

		return

	}

}

// ImportState imports the resource.

func (r *FrontPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	portID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Front Port ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing front port",

			utils.FormatAPIError(fmt.Sprintf("import front port ID %d", portID), err, httpResp),
		)

		return

	}

	var data FrontPortResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *FrontPortResource) mapResponseToModel(ctx context.Context, port *netbox.FrontPort, data *FrontPortResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", port.GetId()))

	data.Name = types.StringValue(port.GetName())

	// Map device - preserve user's input format

	if device := port.GetDevice(); device.Id != 0 {

		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())

	}

	// Map type

	data.Type = types.StringValue(string(port.Type.GetValue()))

	// Map label

	if label, ok := port.GetLabelOk(); ok && label != nil && *label != "" {

		data.Label = types.StringValue(*label)

	} else {

		data.Label = types.StringNull()

	}

	// Map color

	if color, ok := port.GetColorOk(); ok && color != nil && *color != "" {

		data.Color = types.StringValue(*color)

	} else {

		data.Color = types.StringNull()

	}

	// Map rear port - preserve user's input format

	if rearPort := port.GetRearPort(); rearPort.Id != 0 {

		data.RearPort = utils.UpdateReferenceAttribute(data.RearPort, rearPort.GetName(), "", rearPort.GetId())

	}

	// Map rear port position

	if rearPortPos, ok := port.GetRearPortPositionOk(); ok && rearPortPos != nil {

		data.RearPortPosition = types.Int32Value(*rearPortPos)

	} else {

		data.RearPortPosition = types.Int32Value(1)

	}

	// Map description

	if desc, ok := port.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map mark_connected

	if markConnected, ok := port.GetMarkConnectedOk(); ok && markConnected != nil {

		data.MarkConnected = types.BoolValue(*markConnected)

	} else {

		data.MarkConnected = types.BoolValue(false)

	}

	// Handle tags

	if port.HasTags() && len(port.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(port.GetTags())

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

	if port.HasCustomFields() {

		apiCustomFields := port.GetCustomFields()

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
