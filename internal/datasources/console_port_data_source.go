// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ConsolePortDataSource{}

// NewConsolePortDataSource returns a new data source implementing the console port data source.
func NewConsolePortDataSource() datasource.DataSource {
	return &ConsolePortDataSource{}
}

// ConsolePortDataSource defines the data source implementation.
type ConsolePortDataSource struct {
	client *netbox.APIClient
}

// ConsolePortDataSourceModel describes the data source data model.
type ConsolePortDataSourceModel struct {
	ID            types.Int32  `tfsdk:"id"`
	DeviceID      types.Int32  `tfsdk:"device_id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	Speed         types.Int32  `tfsdk:"speed"`
	Description   types.String `tfsdk:"description"`
	DisplayName   types.String `tfsdk:"display_name"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ConsolePortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_port"
}

// Schema defines the schema for the data source.
func (d *ConsolePortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a console port in NetBox. You can identify the console port using `id` or the combination of `device_id` and `name`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the console port.",
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
				MarkdownDescription: "The name of the console port. Used with device_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the console port.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Console port type.",
				Computed:            true,
			},
			"speed": schema.Int32Attribute{
				MarkdownDescription: "Console port speed in bps.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the console port.",
				Computed:            true,
			},
			"display_name": nbschema.DSComputedStringAttribute("The display name of the console port."),
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected.",
				Computed:            true,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ConsolePortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ConsolePortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConsolePortDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var consolePort *netbox.ConsolePort
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Lookup by ID
		portID := data.ID.ValueInt32()
		tflog.Debug(ctx, "Reading console port by ID", map[string]interface{}{
			"id": portID,
		})
		response, httpResp, err := d.client.DcimAPI.DcimConsolePortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading console port",
				utils.FormatAPIError(fmt.Sprintf("read console port ID %d", portID), err, httpResp),
			)
			return
		}
		consolePort = response

	case !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by device_id and name
		deviceID := data.DeviceID.ValueInt32()
		name := data.Name.ValueString()
		tflog.Debug(ctx, "Reading console port by device and name", map[string]interface{}{
			"device_id": deviceID,
			"name":      name,
		})
		response, httpResp, err := d.client.DcimAPI.DcimConsolePortsList(ctx).DeviceId([]int32{deviceID}).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading console port",
				utils.FormatAPIError(fmt.Sprintf("read console port by device ID %d and name %s", deviceID, name), err, httpResp),
			)
			return
		}
		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Console Port Not Found",
				fmt.Sprintf("No console port found for device ID %d with name: %s", deviceID, name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Console Ports Found",
				fmt.Sprintf("Found %d console ports for device ID %d with name %s. Please use ID to select a specific one.", count, deviceID, name),
			)
			return
		}
		consolePort = &response.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or both 'device_id' and 'name' must be specified to lookup a console port.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(consolePort, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ConsolePortDataSource) mapResponseToModel(consolePort *netbox.ConsolePort, data *ConsolePortDataSourceModel) {
	data.ID = types.Int32Value(consolePort.GetId())
	data.Name = types.StringValue(consolePort.GetName())

	// Map device
	if device := consolePort.GetDevice(); device.Id != 0 {
		data.DeviceID = types.Int32Value(device.Id)
		data.Device = types.StringValue(device.GetName())
	}

	// Map label
	if label, ok := consolePort.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map type
	if consolePort.Type != nil {
		data.Type = types.StringValue(string(consolePort.Type.GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map speed
	if consolePort.Speed.IsSet() && consolePort.Speed.Get() != nil {
		data.Speed = types.Int32Value(int32(consolePort.Speed.Get().GetValue()))
	} else {
		data.Speed = types.Int32Null()
	}

	// Map description
	if desc, ok := consolePort.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if mc, ok := consolePort.GetMarkConnectedOk(); ok && mc != nil {
		data.MarkConnected = types.BoolValue(*mc)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Map display name
	if consolePort.GetDisplay() != "" {
		data.DisplayName = types.StringValue(consolePort.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Map custom fields - datasources return ALL fields
	if consolePort.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(consolePort.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(context.Background(), utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
