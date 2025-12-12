// Package datasources provides Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ConsoleServerPortDataSource{}

// NewConsoleServerPortDataSource returns a new data source implementing the console server port data source.
func NewConsoleServerPortDataSource() datasource.DataSource {
	return &ConsoleServerPortDataSource{}
}

// ConsoleServerPortDataSource defines the data source implementation.
type ConsoleServerPortDataSource struct {
	client *netbox.APIClient
}

// ConsoleServerPortDataSourceModel describes the data source data model.
type ConsoleServerPortDataSourceModel struct {
	ID            types.Int32  `tfsdk:"id"`
	DeviceID      types.Int32  `tfsdk:"device_id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	Speed         types.Int32  `tfsdk:"speed"`
	Description   types.String `tfsdk:"description"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
}

// Metadata returns the data source type name.
func (d *ConsoleServerPortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_server_port"
}

// Schema defines the schema for the data source.
func (d *ConsoleServerPortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a console server port in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the console server port.",
				Optional:            true,
				Computed:            true,
			},
			"device_id": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the device. Used with name for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "The name of the device.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the console server port. Used with device_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the console server port.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Console server port type.",
				Computed:            true,
			},
			"speed": schema.Int32Attribute{
				MarkdownDescription: "Console server port speed in bps.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the console server port.",
				Computed:            true,
			},
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ConsoleServerPortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ConsoleServerPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConsoleServerPortDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var consoleServerPort *netbox.ConsoleServerPort

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		// Lookup by ID
		portID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading console server port by ID", map[string]interface{}{
			"id": portID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimConsoleServerPortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading console server port",
				utils.FormatAPIError(fmt.Sprintf("read console server port ID %d", portID), err, httpResp),
			)
			return
		}
		consoleServerPort = response
	} else if !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown() {
		// Lookup by device_id and name
		deviceID := data.DeviceID.ValueInt32()
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading console server port by device and name", map[string]interface{}{
			"device_id": deviceID,
			"name":      name,
		})

		response, httpResp, err := d.client.DcimAPI.DcimConsoleServerPortsList(ctx).DeviceId([]int32{deviceID}).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading console server port",
				utils.FormatAPIError(fmt.Sprintf("read console server port by device ID %d and name %s", deviceID, name), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Console Server Port Not Found",
				fmt.Sprintf("No console server port found for device ID %d with name: %s", deviceID, name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Console Server Ports Found",
				fmt.Sprintf("Found %d console server ports for device ID %d with name %s. Please use ID to select a specific one.", count, deviceID, name),
			)
			return
		}

		consoleServerPort = &response.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or both 'device_id' and 'name' must be specified to lookup a console server port.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(consoleServerPort, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ConsoleServerPortDataSource) mapResponseToModel(consoleServerPort *netbox.ConsoleServerPort, data *ConsoleServerPortDataSourceModel) {
	data.ID = types.Int32Value(consoleServerPort.GetId())
	data.Name = types.StringValue(consoleServerPort.GetName())

	// Map device
	if device := consoleServerPort.GetDevice(); device.Id != 0 {
		data.DeviceID = types.Int32Value(device.Id)
		data.Device = types.StringValue(device.GetName())
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
}
