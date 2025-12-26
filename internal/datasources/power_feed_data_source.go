// Package datasources contains Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &PowerFeedDataSource{}

	_ datasource.DataSourceWithConfigure = &PowerFeedDataSource{}
)

// NewPowerFeedDataSource returns a new data source implementing the PowerFeed data source.

func NewPowerFeedDataSource() datasource.DataSource {

	return &PowerFeedDataSource{}

}

// PowerFeedDataSource defines the data source implementation.

type PowerFeedDataSource struct {
	client *netbox.APIClient
}

// PowerFeedDataSourceModel describes the data source data model.

type PowerFeedDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	DisplayName types.String `tfsdk:"display_name"`

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

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *PowerFeedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_power_feed"

}

// Schema defines the schema for the data source.

func (d *PowerFeedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Retrieves information about a power feed in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the power feed. Use this to look up by ID.",

				Optional: true,

				Computed: true,
			},

			"display_name": schema.StringAttribute{

				MarkdownDescription: "The display name of the power feed.",

				Computed: true,
			},

			"power_panel": schema.StringAttribute{

				MarkdownDescription: "The power panel this feed originates from (ID).",

				Optional: true,

				Computed: true,
			},

			"rack": schema.StringAttribute{

				MarkdownDescription: "The rack this feed connects to (ID).",

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the power feed. Use with power_panel for lookup.",

				Optional: true,

				Computed: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "Status of the power feed.",

				Computed: true,
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "Type of the power feed.",

				Computed: true,
			},

			"supply": schema.StringAttribute{

				MarkdownDescription: "Supply type.",

				Computed: true,
			},

			"phase": schema.StringAttribute{

				MarkdownDescription: "Phase type.",

				Computed: true,
			},

			"voltage": schema.Int64Attribute{

				MarkdownDescription: "Voltage in volts.",

				Computed: true,
			},

			"amperage": schema.Int64Attribute{

				MarkdownDescription: "Amperage in amps.",

				Computed: true,
			},

			"max_utilization": schema.Int64Attribute{

				MarkdownDescription: "Maximum utilization percentage.",

				Computed: true,
			},

			"mark_connected": schema.BoolAttribute{

				MarkdownDescription: "Whether the power feed is treated as connected.",

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the power feed.",

				Computed: true,
			},

			"tenant": schema.StringAttribute{

				MarkdownDescription: "The tenant this power feed belongs to (ID).",

				Computed: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments.",

				Computed: true,
			},

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *PowerFeedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {

		return

	}

	client, ok := req.ProviderData.(*netbox.APIClient)

	if !ok {

		resp.Diagnostics.AddError(

			"Unexpected Data Source Configure Type",

			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return

	}

	d.client = client

}

// Read refreshes the data source data.

func (d *PowerFeedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data PowerFeedDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var pf *netbox.PowerFeed

	// Look up by ID if provided

	switch {

	case !data.ID.IsNull() && !data.ID.IsUnknown():

		pfID, err := utils.ParseID(data.ID.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid Power Feed ID",

				fmt.Sprintf("Power feed ID must be a number, got: %s", data.ID.ValueString()),
			)

			return

		}

		tflog.Debug(ctx, "Reading power feed by ID", map[string]interface{}{

			"id": pfID,
		})

		result, httpResp, err := d.client.DcimAPI.DcimPowerFeedsRetrieve(ctx, pfID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading power feed",

				utils.FormatAPIError(fmt.Sprintf("read power feed ID %d", pfID), err, httpResp),
			)

			return

		}

		pf = result

	case !data.Name.IsNull() && !data.Name.IsUnknown():

		// Look up by name

		tflog.Debug(ctx, "Reading power feed by name", map[string]interface{}{

			"name": data.Name.ValueString(),
		})

		listReq := d.client.DcimAPI.DcimPowerFeedsList(ctx).Name([]string{data.Name.ValueString()})

		// Filter by power_panel if provided

		if !data.PowerPanel.IsNull() && !data.PowerPanel.IsUnknown() {

			ppID, err := utils.ParseID(data.PowerPanel.ValueString())

			if err != nil {

				resp.Diagnostics.AddError(

					"Invalid Power Panel ID",

					fmt.Sprintf("Power panel ID must be a number, got: %s", data.PowerPanel.ValueString()),
				)

				return

			}

			listReq = listReq.PowerPanelId([]int32{ppID})

		}

		listResp, httpResp, err := listReq.Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading power feed",

				utils.FormatAPIError(fmt.Sprintf("read power feed by name %s", data.Name.ValueString()), err, httpResp),
			)

			return

		}

		if listResp.GetCount() == 0 {

			resp.Diagnostics.AddError(

				"Power feed not found",

				fmt.Sprintf("No power feed found with name: %s", data.Name.ValueString()),
			)

			return

		}

		if listResp.GetCount() > 1 {

			resp.Diagnostics.AddError(

				"Multiple power feeds found",

				fmt.Sprintf("Found %d power feeds with name: %s. Please specify the power_panel to narrow results.", listResp.GetCount(), data.Name.ValueString()),
			)

			return

		}

		pf = &listResp.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either 'id' or 'name' must be specified to look up a power feed.",
		)

		return

	}

	// Map response to model

	d.mapResponseToModel(ctx, pf, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (d *PowerFeedDataSource) mapResponseToModel(ctx context.Context, pf *netbox.PowerFeed, data *PowerFeedDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", pf.GetId()))

	// Display Name

	if pf.GetDisplay() != "" {

		data.DisplayName = types.StringValue(pf.GetDisplay())

	} else {

		data.DisplayName = types.StringNull()

	}

	data.Name = types.StringValue(pf.GetName())

	// Map power panel

	data.PowerPanel = types.StringValue(fmt.Sprintf("%d", pf.PowerPanel.GetId()))

	// Map rack

	if pf.Rack.IsSet() && pf.Rack.Get() != nil {

		data.Rack = types.StringValue(fmt.Sprintf("%d", pf.Rack.Get().GetId()))

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

		data.Voltage = types.Int64Null()

	}

	// Map amperage

	if amperage, ok := pf.GetAmperageOk(); ok && amperage != nil {

		data.Amperage = types.Int64Value(int64(*amperage))

	} else {

		data.Amperage = types.Int64Null()

	}

	// Map max_utilization

	if maxUtil, ok := pf.GetMaxUtilizationOk(); ok && maxUtil != nil {

		data.MaxUtilization = types.Int64Value(int64(*maxUtil))

	} else {

		data.MaxUtilization = types.Int64Null()

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

	// Map tenant

	if pf.Tenant.IsSet() && pf.Tenant.Get() != nil {

		data.Tenant = types.StringValue(fmt.Sprintf("%d", pf.Tenant.Get().GetId()))

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

		customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)

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
