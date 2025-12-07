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
var _ datasource.DataSource = &PowerOutletDataSource{}

// NewPowerOutletDataSource returns a new data source implementing the power outlet data source.
func NewPowerOutletDataSource() datasource.DataSource {
	return &PowerOutletDataSource{}
}

// PowerOutletDataSource defines the data source implementation.
type PowerOutletDataSource struct {
	client *netbox.APIClient
}

// PowerOutletDataSourceModel describes the data source data model.
type PowerOutletDataSourceModel struct {
	ID            types.Int32  `tfsdk:"id"`
	DeviceID      types.Int32  `tfsdk:"device_id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	PowerPort     types.Int32  `tfsdk:"power_port"`
	FeedLeg       types.String `tfsdk:"feed_leg"`
	Description   types.String `tfsdk:"description"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
}

// Metadata returns the data source type name.
func (d *PowerOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_outlet"
}

// Schema defines the schema for the data source.
func (d *PowerOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a power outlet in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the power outlet.",
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
				MarkdownDescription: "The name of the power outlet. Used with device_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the power outlet.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Power outlet type.",
				Computed:            true,
			},
			"power_port": schema.Int32Attribute{
				MarkdownDescription: "The power port ID that feeds this outlet.",
				Computed:            true,
			},
			"feed_leg": schema.StringAttribute{
				MarkdownDescription: "Phase leg for three-phase power.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the power outlet.",
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
func (d *PowerOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PowerOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PowerOutletDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var powerOutlet *netbox.PowerOutlet

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		// Lookup by ID
		outletID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading power outlet by ID", map[string]interface{}{
			"id": outletID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimPowerOutletsRetrieve(ctx, outletID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading power outlet",
				utils.FormatAPIError(fmt.Sprintf("read power outlet ID %d", outletID), err, httpResp),
			)
			return
		}
		powerOutlet = response
	} else if !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown() {
		// Lookup by device_id and name
		deviceID := data.DeviceID.ValueInt32()
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading power outlet by device and name", map[string]interface{}{
			"device_id": deviceID,
			"name":      name,
		})

		response, httpResp, err := d.client.DcimAPI.DcimPowerOutletsList(ctx).DeviceId([]int32{deviceID}).Name([]string{name}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading power outlet",
				utils.FormatAPIError(fmt.Sprintf("read power outlet by device ID %d and name %s", deviceID, name), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Power Outlet Not Found",
				fmt.Sprintf("No power outlet found for device ID %d with name: %s", deviceID, name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Power Outlets Found",
				fmt.Sprintf("Found %d power outlets for device ID %d with name %s. Please use ID to select a specific one.", count, deviceID, name),
			)
			return
		}

		powerOutlet = &response.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or both 'device_id' and 'name' must be specified to lookup a power outlet.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(powerOutlet, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *PowerOutletDataSource) mapResponseToModel(powerOutlet *netbox.PowerOutlet, data *PowerOutletDataSourceModel) {
	data.ID = types.Int32Value(powerOutlet.GetId())
	data.Name = types.StringValue(powerOutlet.GetName())

	// Map device
	if device := powerOutlet.GetDevice(); device.Id != 0 {
		data.DeviceID = types.Int32Value(device.Id)
		data.Device = types.StringValue(device.GetName())
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
}
