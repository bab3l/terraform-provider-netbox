// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

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
	_ resource.Resource = &PowerPortResource{}

	_ resource.ResourceWithConfigure = &PowerPortResource{}

	_ resource.ResourceWithImportState = &PowerPortResource{}
)

// NewPowerPortResource returns a new resource implementing the power port resource.

func NewPowerPortResource() resource.Resource {
	return &PowerPortResource{}
}

// PowerPortResource defines the resource implementation.

type PowerPortResource struct {
	client *netbox.APIClient
}

// PowerPortResourceModel describes the resource data model.

type PowerPortResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	MaximumDraw types.Int32 `tfsdk:"maximum_draw"`

	AllocatedDraw types.Int32 `tfsdk:"allocated_draw"`

	Description types.String `tfsdk:"description"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *PowerPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_port"
}

// Schema defines the schema for the resource.

func (r *PowerPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a power port in NetBox. Power ports represent power supply connections on devices.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the power port.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": schema.StringAttribute{
				MarkdownDescription: "The device this power port belongs to (ID or name).",

				Required: true,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the power port.",

				Required: true,
			},

			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the power port.",

				Optional: true,
			},

			"type": schema.StringAttribute{
				MarkdownDescription: "Power port type (e.g., `iec-60320-c14`, `nema-5-15p`, etc.).",

				Optional: true,
			},

			"maximum_draw": schema.Int32Attribute{
				MarkdownDescription: "Maximum power draw in watts.",

				Optional: true,
			},

			"allocated_draw": schema.Int32Attribute{
				MarkdownDescription: "Allocated power draw in watts.",

				Optional: true,
			},

			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("power port"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.

func (r *PowerPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PowerPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PowerPortResourceModel

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

	apiReq := netbox.NewWritablePowerPortRequest(*device, data.Name.ValueString())

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		apiReq.SetLabel(data.Label.ValueString())
	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		portType := netbox.PatchedWritablePowerPortRequestType(data.Type.ValueString())

		apiReq.SetType(portType)
	}

	if !data.MaximumDraw.IsNull() && !data.MaximumDraw.IsUnknown() {
		apiReq.SetMaximumDraw(data.MaximumDraw.ValueInt32())
	}

	if !data.AllocatedDraw.IsNull() && !data.AllocatedDraw.IsUnknown() {
		apiReq.SetAllocatedDraw(data.AllocatedDraw.ValueInt32())
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields

	utils.ApplyDescription(apiReq, data.Description)

	utils.ApplyMetadataFields(ctx, apiReq, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating power port", map[string]interface{}{
		"device": data.Device.ValueString(),

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortsCreate(ctx).WritablePowerPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating power port",

			utils.FormatAPIError(fmt.Sprintf("create power port %s", data.Name.ValueString()), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created power port", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.

func (r *PowerPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PowerPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Power Port ID",

			fmt.Sprintf("Power Port ID must be a number, got: %s", data.ID.ValueString()),
		)

		return
	}

	tflog.Debug(ctx, "Reading power port", map[string]interface{}{
		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading power port",

			utils.FormatAPIError(fmt.Sprintf("read power port ID %d", portID), err, httpResp),
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

func (r *PowerPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PowerPortResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Power Port ID",

			fmt.Sprintf("Power Port ID must be a number, got: %s", data.ID.ValueString()),
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

	apiReq := netbox.NewWritablePowerPortRequest(*device, data.Name.ValueString())

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		apiReq.SetLabel(data.Label.ValueString())
	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		portType := netbox.PatchedWritablePowerPortRequestType(data.Type.ValueString())

		apiReq.SetType(portType)
	}

	if !data.MaximumDraw.IsNull() && !data.MaximumDraw.IsUnknown() {
		apiReq.SetMaximumDraw(data.MaximumDraw.ValueInt32())
	}

	if !data.AllocatedDraw.IsNull() && !data.AllocatedDraw.IsUnknown() {
		apiReq.SetAllocatedDraw(data.AllocatedDraw.ValueInt32())
	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields

	utils.ApplyDescription(apiReq, data.Description)

	utils.ApplyMetadataFields(ctx, apiReq, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating power port", map[string]interface{}{
		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortsUpdate(ctx, portID).WritablePowerPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating power port",

			utils.FormatAPIError(fmt.Sprintf("update power port ID %d", portID), err, httpResp),
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

func (r *PowerPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PowerPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Power Port ID",

			fmt.Sprintf("Power Port ID must be a number, got: %s", data.ID.ValueString()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting power port", map[string]interface{}{
		"id": portID,
	})

	httpResp, err := r.client.DcimAPI.DcimPowerPortsDestroy(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting power port",

			utils.FormatAPIError(fmt.Sprintf("delete power port ID %d", portID), err, httpResp),
		)

		return
	}
}

// ImportState imports an existing resource.

func (r *PowerPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	portID, err := utils.ParseID(req.ID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Power Port ID must be a number, got: %s", req.ID),
		)

		return
	}

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error importing power port",

			utils.FormatAPIError(fmt.Sprintf("import power port ID %d", portID), err, httpResp),
		)

		return
	}

	var data PowerPortResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *PowerPortResource) mapResponseToModel(ctx context.Context, powerPort *netbox.PowerPort, data *PowerPortResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", powerPort.GetId()))

	data.Name = types.StringValue(powerPort.GetName())

	// Map device - preserve user's input format

	if device := powerPort.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map label

	if label, ok := powerPort.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map type

	if powerPort.Type.IsSet() && powerPort.Type.Get() != nil {
		data.Type = types.StringValue(string(powerPort.Type.Get().GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map maximum_draw

	if powerPort.MaximumDraw.IsSet() && powerPort.MaximumDraw.Get() != nil {
		data.MaximumDraw = types.Int32Value(*powerPort.MaximumDraw.Get())
	} else {
		data.MaximumDraw = types.Int32Null()
	}

	// Map allocated_draw

	if powerPort.AllocatedDraw.IsSet() && powerPort.AllocatedDraw.Get() != nil {
		data.AllocatedDraw = types.Int32Value(*powerPort.AllocatedDraw.Get())
	} else {
		data.AllocatedDraw = types.Int32Null()
	}

	// Map description

	if desc, ok := powerPort.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected

	if mc, ok := powerPort.GetMarkConnectedOk(); ok && mc != nil {
		data.MarkConnected = types.BoolValue(*mc)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Handle tags
	data.Tags = utils.PopulateTagsFromAPI(ctx, powerPort.HasTags(), powerPort.GetTags(), data.Tags, diags)

	// Handle custom fields
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, powerPort.HasCustomFields(), powerPort.GetCustomFields(), data.CustomFields, diags)
}
