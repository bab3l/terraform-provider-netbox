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
	_ resource.Resource = &PowerOutletResource{}

	_ resource.ResourceWithConfigure = &PowerOutletResource{}

	_ resource.ResourceWithImportState = &PowerOutletResource{}
)

// NewPowerOutletResource returns a new resource implementing the power outlet resource.

func NewPowerOutletResource() resource.Resource {

	return &PowerOutletResource{}

}

// PowerOutletResource defines the resource implementation.

type PowerOutletResource struct {
	client *netbox.APIClient
}

// PowerOutletResourceModel describes the resource data model.

type PowerOutletResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	PowerPort types.Int32 `tfsdk:"power_port"`

	FeedLeg types.String `tfsdk:"feed_leg"`

	Description types.String `tfsdk:"description"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *PowerOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_power_outlet"

}

// Schema defines the schema for the resource.

func (r *PowerOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a power outlet in NetBox. Power outlets represent power distribution connections on PDUs and other power distribution devices.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the power outlet.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The device this power outlet belongs to (ID or name).",

				Required: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the power outlet.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the power outlet.",

				Optional: true,
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "Power outlet type (e.g., `iec-60320-c13`, `nema-5-15r`, etc.).",

				Optional: true,
			},

			"power_port": schema.Int32Attribute{

				MarkdownDescription: "The power port ID that feeds this outlet.",

				Optional: true,
			},

			"feed_leg": schema.StringAttribute{

				MarkdownDescription: "Phase leg for three-phase power. Valid values: `A`, `B`, `C`.",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the power outlet.",

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

func (r *PowerOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *PowerOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data PowerOutletResourceModel

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

	apiReq := netbox.NewWritablePowerOutletRequest(*device, data.Name.ValueString())

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {

		outletType := netbox.PatchedWritablePowerOutletRequestType(data.Type.ValueString())

		apiReq.SetType(outletType)

	}

	if !data.PowerPort.IsNull() && !data.PowerPort.IsUnknown() {

		powerPortReq := netbox.BriefPowerPortRequest{

			Name: fmt.Sprintf("Power Port %d", data.PowerPort.ValueInt32()),
		}

		apiReq.SetPowerPort(powerPortReq)

	}

	if !data.FeedLeg.IsNull() && !data.FeedLeg.IsUnknown() {

		feedLeg := netbox.PatchedWritablePowerOutletRequestFeedLeg(data.FeedLeg.ValueString())

		apiReq.SetFeedLeg(feedLeg)

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

	tflog.Debug(ctx, "Creating power outlet", map[string]interface{}{

		"device": data.Device.ValueString(),

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsCreate(ctx).WritablePowerOutletRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating power outlet",

			utils.FormatAPIError(fmt.Sprintf("create power outlet %s", data.Name.ValueString()), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created power outlet", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *PowerOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data PowerOutletResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	outletID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Power Outlet ID",

			fmt.Sprintf("Power Outlet ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading power outlet", map[string]interface{}{

		"id": outletID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsRetrieve(ctx, outletID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading power outlet",

			utils.FormatAPIError(fmt.Sprintf("read power outlet ID %d", outletID), err, httpResp),
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

func (r *PowerOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data PowerOutletResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	outletID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Power Outlet ID",

			fmt.Sprintf("Power Outlet ID must be a number, got: %s", data.ID.ValueString()),
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

	apiReq := netbox.NewWritablePowerOutletRequest(*device, data.Name.ValueString())

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {

		outletType := netbox.PatchedWritablePowerOutletRequestType(data.Type.ValueString())

		apiReq.SetType(outletType)

	}

	if !data.PowerPort.IsNull() && !data.PowerPort.IsUnknown() {

		powerPortReq := netbox.BriefPowerPortRequest{

			Name: fmt.Sprintf("Power Port %d", data.PowerPort.ValueInt32()),
		}

		apiReq.SetPowerPort(powerPortReq)

	}

	if !data.FeedLeg.IsNull() && !data.FeedLeg.IsUnknown() {

		feedLeg := netbox.PatchedWritablePowerOutletRequestFeedLeg(data.FeedLeg.ValueString())

		apiReq.SetFeedLeg(feedLeg)

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

	tflog.Debug(ctx, "Updating power outlet", map[string]interface{}{

		"id": outletID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsUpdate(ctx, outletID).WritablePowerOutletRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating power outlet",

			utils.FormatAPIError(fmt.Sprintf("update power outlet ID %d", outletID), err, httpResp),
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

func (r *PowerOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data PowerOutletResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	outletID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Power Outlet ID",

			fmt.Sprintf("Power Outlet ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting power outlet", map[string]interface{}{

		"id": outletID,
	})

	httpResp, err := r.client.DcimAPI.DcimPowerOutletsDestroy(ctx, outletID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting power outlet",

			utils.FormatAPIError(fmt.Sprintf("delete power outlet ID %d", outletID), err, httpResp),
		)

		return

	}

}

// ImportState imports an existing resource.

func (r *PowerOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	outletID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Power Outlet ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsRetrieve(ctx, outletID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing power outlet",

			utils.FormatAPIError(fmt.Sprintf("import power outlet ID %d", outletID), err, httpResp),
		)

		return

	}

	var data PowerOutletResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *PowerOutletResource) mapResponseToModel(ctx context.Context, powerOutlet *netbox.PowerOutlet, data *PowerOutletResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", powerOutlet.GetId()))

	data.Name = types.StringValue(powerOutlet.GetName())

	// Map device

	if device := powerOutlet.GetDevice(); device.Id != 0 {

		data.Device = types.StringValue(fmt.Sprintf("%d", device.GetId()))

	}

	// Map label

	if label, ok := powerOutlet.GetLabelOk(); ok && label != nil && *label != "" {

		data.Label = types.StringValue(*label)

	} else {

		data.Label = types.StringNull()

	}

	// Map type

	if powerOutlet.Type.IsSet() && powerOutlet.Type.Get() != nil {

		data.Type = types.StringValue(string(powerOutlet.Type.Get().GetValue()))

	} else {

		data.Type = types.StringNull()

	}

	// Map power_port

	if powerOutlet.PowerPort.IsSet() && powerOutlet.PowerPort.Get() != nil {

		data.PowerPort = types.Int32Value(powerOutlet.PowerPort.Get().Id)

	} else {

		data.PowerPort = types.Int32Null()

	}

	// Map feed_leg

	if powerOutlet.FeedLeg.IsSet() && powerOutlet.FeedLeg.Get() != nil {

		data.FeedLeg = types.StringValue(string(powerOutlet.FeedLeg.Get().GetValue()))

	} else {

		data.FeedLeg = types.StringNull()

	}

	// Map description

	if desc, ok := powerOutlet.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map mark_connected

	if mc, ok := powerOutlet.GetMarkConnectedOk(); ok && mc != nil {

		data.MarkConnected = types.BoolValue(*mc)

	} else {

		data.MarkConnected = types.BoolValue(false)

	}

	// Handle tags

	if powerOutlet.HasTags() && len(powerOutlet.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(powerOutlet.GetTags())

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

	if powerOutlet.HasCustomFields() {

		apiCustomFields := powerOutlet.GetCustomFields()

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
