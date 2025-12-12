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
var _ datasource.DataSource = &RearPortDataSource{}

// NewRearPortDataSource returns a new data source implementing the rear port data source.
func NewRearPortDataSource() datasource.DataSource {
	return &RearPortDataSource{}
}

// RearPortDataSource defines the data source implementation.
type RearPortDataSource struct {
	client *netbox.APIClient
}

// RearPortDataSourceModel describes the data source data model.
type RearPortDataSourceModel struct {
	ID            types.Int32  `tfsdk:"id"`
	DeviceID      types.Int32  `tfsdk:"device_id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	Color         types.String `tfsdk:"color"`
	Positions     types.Int32  `tfsdk:"positions"`
	Description   types.String `tfsdk:"description"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
}

// Metadata returns the data source type name.
func (d *RearPortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rear_port"
}

// Schema defines the schema for the data source.
func (d *RearPortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a rear port in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the rear port.",
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
				MarkdownDescription: "The name of the rear port. Used with device_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the rear port.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of rear port.",
				Computed:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the rear port in hex format.",
				Computed:            true,
			},
			"positions": schema.Int32Attribute{
				MarkdownDescription: "Number of front ports that may be mapped to this rear port.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the rear port.",
				Computed:            true,
			},
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Whether the port is marked as connected.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RearPortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *RearPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RearPortDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var port *netbox.RearPort

	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Lookup by ID
		portID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading rear port by ID", map[string]interface{}{
			"id": portID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimRearPortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading rear port",
				utils.FormatAPIError(fmt.Sprintf("read rear port ID %d", portID), err, httpResp),
			)
			return
		}
		port = response
	case !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by device_id and name
		deviceID := data.DeviceID.ValueInt32()
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading rear port by device and name", map[string]interface{}{
			"device_id": deviceID,
			"name":      name,
		})

		response, httpResp, err := d.client.DcimAPI.DcimRearPortsList(ctx).DeviceId([]int32{deviceID}).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading rear port",
				utils.FormatAPIError(fmt.Sprintf("read rear port by device ID %d and name %s", deviceID, name), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Rear Port Not Found",
				fmt.Sprintf("No rear port found for device ID %d with name: %s", deviceID, name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Rear Ports Found",
				fmt.Sprintf("Found %d rear ports for device ID %d with name %s. Please use ID to select a specific one.", count, deviceID, name),
			)
			return
		}

		port = &response.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or both 'device_id' and 'name' must be specified to lookup a rear port.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(port, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *RearPortDataSource) mapResponseToModel(port *netbox.RearPort, data *RearPortDataSourceModel) {
	data.ID = types.Int32Value(port.GetId())
	data.Name = types.StringValue(port.GetName())

	// Map device
	if device := port.GetDevice(); device.Id != 0 {
		data.DeviceID = types.Int32Value(device.Id)
		data.Device = types.StringValue(device.GetName())
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

	// Map positions
	if positions, ok := port.GetPositionsOk(); ok && positions != nil {
		data.Positions = types.Int32Value(*positions)
	} else {
		data.Positions = types.Int32Value(1)
	}

	// Map description
	if desc, ok := port.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if mc, ok := port.GetMarkConnectedOk(); ok && mc != nil {
		data.MarkConnected = types.BoolValue(*mc)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}
}
