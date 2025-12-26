// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &PowerPortDataSource{}

// NewPowerPortDataSource returns a new data source implementing the power port data source.

func NewPowerPortDataSource() datasource.DataSource {

	return &PowerPortDataSource{}

}

// PowerPortDataSource defines the data source implementation.

type PowerPortDataSource struct {
	client *netbox.APIClient
}

// PowerPortDataSourceModel describes the data source data model.

type PowerPortDataSourceModel struct {
	ID types.Int32 `tfsdk:"id"`

	DeviceID types.Int32 `tfsdk:"device_id"`

	Device types.String `tfsdk:"device"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	MaximumDraw types.Int32 `tfsdk:"maximum_draw"`

	AllocatedDraw types.Int32 `tfsdk:"allocated_draw"`

	Description types.String `tfsdk:"description"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	DisplayName types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.

func (d *PowerPortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_power_port"

}

// Schema defines the schema for the data source.

func (d *PowerPortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Retrieves information about a power port in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.Int32Attribute{

				MarkdownDescription: "The unique numeric ID of the power port.",

				Optional: true,

				Computed: true,
			},

			"device_id": schema.Int32Attribute{

				MarkdownDescription: "The numeric ID of the device. Used with name for lookup when ID is not provided.",

				Optional: true,

				Computed: true,
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The name of the device.",

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the power port. Used with device_id for lookup when ID is not provided.",

				Optional: true,

				Computed: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the power port.",

				Computed: true,
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "Power port type.",

				Computed: true,
			},

			"maximum_draw": schema.Int32Attribute{

				MarkdownDescription: "Maximum power draw in watts.",

				Computed: true,
			},

			"allocated_draw": schema.Int32Attribute{

				MarkdownDescription: "Allocated power draw in watts.",

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the power port.",

				Computed: true,
			},

			"mark_connected": schema.BoolAttribute{

				MarkdownDescription: "Treat as if a cable is connected.",

				Computed: true,
			},

			"display_name": schema.StringAttribute{

				MarkdownDescription: "The display name of the power port.",

				Computed: true,
			},
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *PowerPortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read retrieves the data source data.

func (d *PowerPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data PowerPortDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var powerPort *netbox.PowerPort

	switch {

	case !data.ID.IsNull() && !data.ID.IsUnknown():

		// Lookup by ID

		portID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading power port by ID", map[string]interface{}{

			"id": portID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimPowerPortsRetrieve(ctx, portID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading power port",

				utils.FormatAPIError(fmt.Sprintf("read power port ID %d", portID), err, httpResp),
			)

			return

		}

		powerPort = response

	case !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown():

		// Lookup by device_id and name

		deviceID := data.DeviceID.ValueInt32()

		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading power port by device and name", map[string]interface{}{

			"device_id": deviceID,

			"name": name,
		})

		response, httpResp, err := d.client.DcimAPI.DcimPowerPortsList(ctx).DeviceId([]int32{deviceID}).Name([]string{name}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading power port",

				utils.FormatAPIError(fmt.Sprintf("read power port by device ID %d and name %s", deviceID, name), err, httpResp),
			)

			return

		}

		count := int(response.GetCount())

		if count == 0 {

			resp.Diagnostics.AddError(

				"Power Port Not Found",

				fmt.Sprintf("No power port found for device ID %d with name: %s", deviceID, name),
			)

			return

		}

		if count > 1 {

			resp.Diagnostics.AddError(

				"Multiple Power Ports Found",

				fmt.Sprintf("Found %d power ports for device ID %d with name %s. Please use ID to select a specific one.", count, deviceID, name),
			)

			return

		}

		powerPort = &response.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either 'id' or both 'device_id' and 'name' must be specified to lookup a power port.",
		)

		return

	}

	// Map response to model

	d.mapResponseToModel(powerPort, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (d *PowerPortDataSource) mapResponseToModel(powerPort *netbox.PowerPort, data *PowerPortDataSourceModel) {

	data.ID = types.Int32Value(powerPort.GetId())

	data.Name = types.StringValue(powerPort.GetName())

	// Map device

	if device := powerPort.GetDevice(); device.Id != 0 {

		data.DeviceID = types.Int32Value(device.Id)

		data.Device = types.StringValue(device.GetName())

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
	// Map display_name

	if powerPort.GetDisplay() != "" {

		data.DisplayName = types.StringValue(powerPort.GetDisplay())

	} else {

		data.DisplayName = types.StringNull()

	}
}
