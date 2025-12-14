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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ConsoleServerPortResource{}

	_ resource.ResourceWithConfigure = &ConsoleServerPortResource{}

	_ resource.ResourceWithImportState = &ConsoleServerPortResource{}
)

// NewConsoleServerPortResource returns a new resource implementing the console server port resource.

func NewConsoleServerPortResource() resource.Resource {

	return &ConsoleServerPortResource{}

}

// ConsoleServerPortResource defines the resource implementation.

type ConsoleServerPortResource struct {
	client *netbox.APIClient
}

// ConsoleServerPortResourceModel describes the resource data model.

type ConsoleServerPortResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	Speed types.Int32 `tfsdk:"speed"`

	Description types.String `tfsdk:"description"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ConsoleServerPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_console_server_port"

}

// Schema defines the schema for the resource.

func (r *ConsoleServerPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a console server port in NetBox. Console server ports are physical console connections on console servers that provide remote access to other devices.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the console server port.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The device this console server port belongs to (ID or name).",

				Required: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the console server port.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the console server port.",

				Optional: true,
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "Console server port type. Valid values: `de-9`, `db-25`, `rj-11`, `rj-12`, `rj-45`, `mini-din-8`, `usb-a`, `usb-b`, `usb-c`, `usb-mini-a`, `usb-mini-b`, `usb-micro-a`, `usb-micro-b`, `usb-micro-ab`, `other`.",

				Optional: true,
			},

			"speed": schema.Int32Attribute{

				MarkdownDescription: "Console server port speed in bps. Valid values: `1200`, `2400`, `4800`, `9600`, `19200`, `38400`, `57600`, `115200`.",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the console server port.",

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

func (r *ConsoleServerPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ConsoleServerPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ConsoleServerPortResourceModel

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

	// Build request

	apiReq := netbox.NewWritableConsoleServerPortRequest(*device, data.Name.ValueString())

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {

		portType := netbox.PatchedWritableConsolePortRequestType(data.Type.ValueString())

		apiReq.SetType(portType)

	}

	if !data.Speed.IsNull() && !data.Speed.IsUnknown() {

		speed := netbox.PatchedWritableConsolePortRequestSpeed(data.Speed.ValueInt32())

		apiReq.SetSpeed(speed)

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

	tflog.Debug(ctx, "Creating console server port", map[string]interface{}{

		"device": data.Device.ValueString(),

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsCreate(ctx).WritableConsoleServerPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating console server port",

			utils.FormatAPIError(fmt.Sprintf("create console server port %s", data.Name.ValueString()), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created console server port", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *ConsoleServerPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ConsoleServerPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Console Server Port ID",

			fmt.Sprintf("Console Server Port ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading console server port", map[string]interface{}{

		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading console server port",

			utils.FormatAPIError(fmt.Sprintf("read console server port ID %d", portID), err, httpResp),
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

func (r *ConsoleServerPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ConsoleServerPortResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Console Server Port ID",

			fmt.Sprintf("Console Server Port ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	// Lookup device

	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritableConsoleServerPortRequest(*device, data.Name.ValueString())

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {

		portType := netbox.PatchedWritableConsolePortRequestType(data.Type.ValueString())

		apiReq.SetType(portType)

	}

	if !data.Speed.IsNull() && !data.Speed.IsUnknown() {

		speed := netbox.PatchedWritableConsolePortRequestSpeed(data.Speed.ValueInt32())

		apiReq.SetSpeed(speed)

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

	tflog.Debug(ctx, "Updating console server port", map[string]interface{}{

		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsUpdate(ctx, portID).WritableConsoleServerPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating console server port",

			utils.FormatAPIError(fmt.Sprintf("update console server port ID %d", portID), err, httpResp),
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

func (r *ConsoleServerPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ConsoleServerPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Console Server Port ID",

			fmt.Sprintf("Console Server Port ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting console server port", map[string]interface{}{

		"id": portID,
	})

	httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsDestroy(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting console server port",

			utils.FormatAPIError(fmt.Sprintf("delete console server port ID %d", portID), err, httpResp),
		)

		return

	}

}

// ImportState imports an existing resource.

func (r *ConsoleServerPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	portID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Console Server Port ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing console server port",

			utils.FormatAPIError(fmt.Sprintf("import console server port ID %d", portID), err, httpResp),
		)

		return

	}

	var data ConsoleServerPortResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *ConsoleServerPortResource) mapResponseToModel(ctx context.Context, consoleServerPort *netbox.ConsoleServerPort, data *ConsoleServerPortResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", consoleServerPort.GetId()))

	data.Name = types.StringValue(consoleServerPort.GetName())

	// Map device

	if device := consoleServerPort.GetDevice(); device.Id != 0 {

		data.Device = types.StringValue(fmt.Sprintf("%d", device.GetId()))

	}

	// Map label

	if label, ok := consoleServerPort.GetLabelOk(); ok && label != nil && *label != "" {

		data.Label = types.StringValue(*label)

	} else {

		data.Label = types.StringNull()

	}

	// Map type

	if consoleServerPort.Type != nil {

		data.Type = types.StringValue(string(consoleServerPort.Type.GetValue()))

	} else {

		data.Type = types.StringNull()

	}

	// Map speed

	if consoleServerPort.Speed.IsSet() && consoleServerPort.Speed.Get() != nil {

		data.Speed = types.Int32Value(int32(consoleServerPort.Speed.Get().GetValue()))

	} else {

		data.Speed = types.Int32Null()

	}

	// Map description

	if desc, ok := consoleServerPort.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map mark_connected

	if mc, ok := consoleServerPort.GetMarkConnectedOk(); ok && mc != nil {

		data.MarkConnected = types.BoolValue(*mc)

	} else {

		data.MarkConnected = types.BoolValue(false)

	}

	// Handle tags

	if consoleServerPort.HasTags() && len(consoleServerPort.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(consoleServerPort.GetTags())

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

	if consoleServerPort.HasCustomFields() {

		apiCustomFields := consoleServerPort.GetCustomFields()

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
