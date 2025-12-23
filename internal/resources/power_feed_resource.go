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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &PowerFeedResource{}

	_ resource.ResourceWithConfigure = &PowerFeedResource{}

	_ resource.ResourceWithImportState = &PowerFeedResource{}
)

// NewPowerFeedResource returns a new resource implementing the power feed resource.

func NewPowerFeedResource() resource.Resource {

	return &PowerFeedResource{}

}

// PowerFeedResource defines the resource implementation.

type PowerFeedResource struct {
	client *netbox.APIClient
}

// PowerFeedResourceModel describes the resource data model.

type PowerFeedResourceModel struct {
	ID types.String `tfsdk:"id"`

	PowerPanel types.String `tfsdk:"power_panel"`

	Rack types.String `tfsdk:"rack"`

	Name types.String `tfsdk:"name"`

	Status types.String `tfsdk:"status"`

	Type types.String `tfsdk:"type"`

	Supply types.String `tfsdk:"supply"`

	Phase types.String `tfsdk:"phase"`

	Voltage types.Int64 `tfsdk:"voltage"`

	Amperage types.Int64 `tfsdk:"amperage"`

	MaxUtilization types.Int64 `tfsdk:"max_utilization"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Description types.String `tfsdk:"description"`

	Tenant types.String `tfsdk:"tenant"`

	Comments types.String `tfsdk:"comments"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *PowerFeedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_power_feed"

}

// Schema defines the schema for the resource.

func (r *PowerFeedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a power feed in NetBox. Power feeds represent connections from power panels to racks.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the power feed.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"power_panel": schema.StringAttribute{

				MarkdownDescription: "The power panel this feed originates from (ID or name).",

				Required: true,
			},

			"rack": schema.StringAttribute{

				MarkdownDescription: "The rack this feed connects to (ID or name).",

				Optional: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the power feed.",

				Required: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "Status of the power feed. Valid values: `offline`, `active`, `planned`, `failed`. Default: `active`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("active"),
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "Type of the power feed. Valid values: `primary`, `redundant`. Default: `primary`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("primary"),
			},

			"supply": schema.StringAttribute{

				MarkdownDescription: "Supply type. Valid values: `ac`, `dc`. Default: `ac`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("ac"),
			},

			"phase": schema.StringAttribute{

				MarkdownDescription: "Phase type. Valid values: `single-phase`, `three-phase`. Default: `single-phase`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("single-phase"),
			},

			"voltage": schema.Int64Attribute{

				MarkdownDescription: "Voltage in volts. Default: 120.",

				Optional: true,

				Computed: true,

				Default: int64default.StaticInt64(120),
			},

			"amperage": schema.Int64Attribute{

				MarkdownDescription: "Amperage in amps. Default: 20.",

				Optional: true,

				Computed: true,

				Default: int64default.StaticInt64(20),
			},

			"max_utilization": schema.Int64Attribute{

				MarkdownDescription: "Maximum utilization percentage (1-100). Default: 80.",

				Optional: true,

				Computed: true,

				Default: int64default.StaticInt64(80),
			},

			"mark_connected": schema.BoolAttribute{

				MarkdownDescription: "Treat as if a cable is connected. Default: false.",

				Optional: true,

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the power feed.",

				Optional: true,
			},

			"tenant": schema.StringAttribute{

				MarkdownDescription: "The tenant this power feed belongs to (ID or slug).",

				Optional: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about the power feed.",

				Optional: true,
			},

			"display_name": nbschema.DisplayNameAttribute("power feed"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *PowerFeedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *PowerFeedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data PowerFeedResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Lookup power panel

	powerPanel, diags := lookup.LookupPowerPanel(ctx, r.client, data.PowerPanel.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritablePowerFeedRequest(*powerPanel, data.Name.ValueString())

	// Set optional fields

	if !data.Rack.IsNull() && !data.Rack.IsUnknown() {

		rack, diags := lookup.LookupRack(ctx, r.client, data.Rack.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetRack(*rack)

	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {

		status := netbox.PatchedWritablePowerFeedRequestStatus(data.Status.ValueString())

		apiReq.SetStatus(status)

	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {

		feedType := netbox.PatchedWritablePowerFeedRequestType(data.Type.ValueString())

		apiReq.SetType(feedType)

	}

	if !data.Supply.IsNull() && !data.Supply.IsUnknown() {

		supply := netbox.PatchedWritablePowerFeedRequestSupply(data.Supply.ValueString())

		apiReq.SetSupply(supply)

	}

	if !data.Phase.IsNull() && !data.Phase.IsUnknown() {

		phase := netbox.PatchedWritablePowerFeedRequestPhase(data.Phase.ValueString())

		apiReq.SetPhase(phase)

	}

	if !data.Voltage.IsNull() && !data.Voltage.IsUnknown() {

		voltage, err := utils.SafeInt32FromValue(data.Voltage)

		if err != nil {

			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Voltage value overflow: %s", err))

			return

		}

		apiReq.SetVoltage(voltage)

	}

	if !data.Amperage.IsNull() && !data.Amperage.IsUnknown() {

		amperage, err := utils.SafeInt32FromValue(data.Amperage)

		if err != nil {

			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Amperage value overflow: %s", err))

			return

		}

		apiReq.SetAmperage(amperage)

	}

	if !data.MaxUtilization.IsNull() && !data.MaxUtilization.IsUnknown() {

		maxUtilization, err := utils.SafeInt32FromValue(data.MaxUtilization)

		if err != nil {

			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("MaxUtilization value overflow: %s", err))

			return

		}

		apiReq.SetMaxUtilization(maxUtilization)

	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {

		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {

		tenant, diags := lookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTenant(*tenant)

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

	tflog.Debug(ctx, "Creating power feed", map[string]interface{}{

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerFeedsCreate(ctx).WritablePowerFeedRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating power feed",

			utils.FormatAPIError(fmt.Sprintf("create power feed %s", data.Name.ValueString()), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created power feed", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *PowerFeedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data PowerFeedResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	pfID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Power Feed ID",

			fmt.Sprintf("Power feed ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading power feed", map[string]interface{}{

		"id": pfID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerFeedsRetrieve(ctx, pfID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading power feed",

			utils.FormatAPIError(fmt.Sprintf("read power feed ID %d", pfID), err, httpResp),
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

func (r *PowerFeedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data PowerFeedResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	pfID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Power Feed ID",

			fmt.Sprintf("Power feed ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	// Lookup power panel

	powerPanel, diags := lookup.LookupPowerPanel(ctx, r.client, data.PowerPanel.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritablePowerFeedRequest(*powerPanel, data.Name.ValueString())

	// Set optional fields

	if !data.Rack.IsNull() && !data.Rack.IsUnknown() {

		rack, diags := lookup.LookupRack(ctx, r.client, data.Rack.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetRack(*rack)

	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {

		status := netbox.PatchedWritablePowerFeedRequestStatus(data.Status.ValueString())

		apiReq.SetStatus(status)

	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {

		feedType := netbox.PatchedWritablePowerFeedRequestType(data.Type.ValueString())

		apiReq.SetType(feedType)

	}

	if !data.Supply.IsNull() && !data.Supply.IsUnknown() {

		supply := netbox.PatchedWritablePowerFeedRequestSupply(data.Supply.ValueString())

		apiReq.SetSupply(supply)

	}

	if !data.Phase.IsNull() && !data.Phase.IsUnknown() {

		phase := netbox.PatchedWritablePowerFeedRequestPhase(data.Phase.ValueString())

		apiReq.SetPhase(phase)

	}

	if !data.Voltage.IsNull() && !data.Voltage.IsUnknown() {

		voltage, err := utils.SafeInt32FromValue(data.Voltage)

		if err != nil {

			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Voltage value overflow: %s", err))

			return

		}

		apiReq.SetVoltage(voltage)

	}

	if !data.Amperage.IsNull() && !data.Amperage.IsUnknown() {

		amperage, err := utils.SafeInt32FromValue(data.Amperage)

		if err != nil {

			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Amperage value overflow: %s", err))

			return

		}

		apiReq.SetAmperage(amperage)

	}

	if !data.MaxUtilization.IsNull() && !data.MaxUtilization.IsUnknown() {

		maxUtilization, err := utils.SafeInt32FromValue(data.MaxUtilization)

		if err != nil {

			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("MaxUtilization value overflow: %s", err))

			return

		}

		apiReq.SetMaxUtilization(maxUtilization)

	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {

		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {

		tenant, diags := lookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTenant(*tenant)

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

	tflog.Debug(ctx, "Updating power feed", map[string]interface{}{

		"id": pfID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerFeedsUpdate(ctx, pfID).WritablePowerFeedRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating power feed",

			utils.FormatAPIError(fmt.Sprintf("update power feed ID %d", pfID), err, httpResp),
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

func (r *PowerFeedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data PowerFeedResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	pfID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Power Feed ID",

			fmt.Sprintf("Power feed ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting power feed", map[string]interface{}{

		"id": pfID,
	})

	httpResp, err := r.client.DcimAPI.DcimPowerFeedsDestroy(ctx, pfID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting power feed",

			utils.FormatAPIError(fmt.Sprintf("delete power feed ID %d", pfID), err, httpResp),
		)

		return

	}

}

// ImportState imports an existing resource.

func (r *PowerFeedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	pfID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Power feed ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimPowerFeedsRetrieve(ctx, pfID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing power feed",

			utils.FormatAPIError(fmt.Sprintf("import power feed ID %d", pfID), err, httpResp),
		)

		return

	}

	var data PowerFeedResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *PowerFeedResource) mapResponseToModel(ctx context.Context, pf *netbox.PowerFeed, data *PowerFeedResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", pf.GetId()))

	data.Name = types.StringValue(pf.GetName())

	// DisplayName
	if pf.Display != "" {
		data.DisplayName = types.StringValue(pf.Display)
	} else {
		data.DisplayName = types.StringNull()
	}

	// Map power panel

	data.PowerPanel = types.StringValue(fmt.Sprintf("%d", pf.PowerPanel.GetId()))

	// Map rack - preserve user's input format

	if pf.Rack.IsSet() && pf.Rack.Get() != nil {

		rack := pf.Rack.Get()

		data.Rack = utils.UpdateReferenceAttribute(data.Rack, rack.GetName(), "", rack.GetId())

	} else {

		data.Rack = types.StringNull()

	}

	// Map status

	if status, ok := pf.GetStatusOk(); ok && status != nil {

		data.Status = types.StringValue(string(status.GetValue()))

	} else {

		data.Status = types.StringNull()

	}

	// Map type

	if feedType, ok := pf.GetTypeOk(); ok && feedType != nil {

		data.Type = types.StringValue(string(feedType.GetValue()))

	} else {

		data.Type = types.StringNull()

	}

	// Map supply

	if supply, ok := pf.GetSupplyOk(); ok && supply != nil {

		data.Supply = types.StringValue(string(supply.GetValue()))

	} else {

		data.Supply = types.StringNull()

	}

	// Map phase

	if phase, ok := pf.GetPhaseOk(); ok && phase != nil {

		data.Phase = types.StringValue(string(phase.GetValue()))

	} else {

		data.Phase = types.StringNull()

	}

	// Map voltage

	if voltage, ok := pf.GetVoltageOk(); ok && voltage != nil {

		data.Voltage = types.Int64Value(int64(*voltage))

	} else {

		data.Voltage = types.Int64Value(120)

	}

	// Map amperage

	if amperage, ok := pf.GetAmperageOk(); ok && amperage != nil {

		data.Amperage = types.Int64Value(int64(*amperage))

	} else {

		data.Amperage = types.Int64Value(20)

	}

	// Map max_utilization

	if maxUtil, ok := pf.GetMaxUtilizationOk(); ok && maxUtil != nil {

		data.MaxUtilization = types.Int64Value(int64(*maxUtil))

	} else {

		data.MaxUtilization = types.Int64Value(80)

	}

	// Map mark_connected

	if mc, ok := pf.GetMarkConnectedOk(); ok && mc != nil {

		data.MarkConnected = types.BoolValue(*mc)

	} else {

		data.MarkConnected = types.BoolNull()

	}

	// Map description

	if desc, ok := pf.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map tenant - preserve user's input format

	if pf.Tenant.IsSet() && pf.Tenant.Get() != nil {

		tenant := pf.Tenant.Get()

		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())

	} else {

		data.Tenant = types.StringNull()

	}

	// Map comments

	if comments, ok := pf.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if pf.HasTags() && len(pf.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(pf.GetTags())

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

	if pf.HasCustomFields() {

		apiCustomFields := pf.GetCustomFields()

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
