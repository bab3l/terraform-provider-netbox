// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InterfaceDataSource{}

func NewInterfaceDataSource() datasource.DataSource {
	return &InterfaceDataSource{}
}

// InterfaceDataSource defines the data source implementation.
type InterfaceDataSource struct {
	client *netbox.APIClient
}

// InterfaceDataSourceModel describes the data source data model.
type InterfaceDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Parent        types.String `tfsdk:"parent"`
	Bridge        types.String `tfsdk:"bridge"`
	Lag           types.String `tfsdk:"lag"`
	Mtu           types.Int64  `tfsdk:"mtu"`
	MacAddress    types.String `tfsdk:"mac_address"`
	Speed         types.Int64  `tfsdk:"speed"`
	Duplex        types.String `tfsdk:"duplex"`
	Wwn           types.String `tfsdk:"wwn"`
	MgmtOnly      types.Bool   `tfsdk:"mgmt_only"`
	Description   types.String `tfsdk:"description"`
	Mode          types.String `tfsdk:"mode"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

func (d *InterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface"
}

func (d *InterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an interface in Netbox. Interfaces represent physical or virtual network interfaces on devices. You can identify the interface using `id` or by combining `device` and `name`.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.DSIDAttribute("interface"),
			"device": schema.StringAttribute{
				MarkdownDescription: "ID or name of the device this interface belongs to. Required when looking up by `name`.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the interface (e.g., 'eth0', 'GigabitEthernet0/0'). Must be combined with `device` for lookup.",
				Optional:            true,
				Computed:            true,
			},
			"label":          nbschema.DSComputedStringAttribute("Physical label on the interface."),
			"type":           nbschema.DSComputedStringAttribute("Type of interface (e.g., 'virtual', '1000base-t', '10gbase-x-sfpp')."),
			"enabled":        nbschema.DSComputedBoolAttribute("Whether the interface is enabled."),
			"parent":         nbschema.DSComputedStringAttribute("ID of the parent interface (for sub-interfaces)."),
			"bridge":         nbschema.DSComputedStringAttribute("ID of the bridge interface this interface belongs to."),
			"lag":            nbschema.DSComputedStringAttribute("ID of the LAG this interface is a member of."),
			"mtu":            nbschema.DSComputedInt64Attribute("Maximum transmission unit (MTU) size."),
			"mac_address":    nbschema.DSComputedStringAttribute("MAC address of the interface."),
			"speed":          nbschema.DSComputedInt64Attribute("Interface speed in Kbps."),
			"duplex":         nbschema.DSComputedStringAttribute("Duplex mode (half, full, auto)."),
			"wwn":            nbschema.DSComputedStringAttribute("World Wide Name for Fibre Channel interfaces."),
			"mgmt_only":      nbschema.DSComputedBoolAttribute("Interface is used only for out-of-band management."),
			"description":    nbschema.DSComputedStringAttribute("Brief description of the interface."),
			"mode":           nbschema.DSComputedStringAttribute("802.1Q mode (access, tagged, tagged-all)."),
			"mark_connected": nbschema.DSComputedBoolAttribute("Treat as if a cable is connected."),
			"tags":           nbschema.DSTagsAttribute(),
			"custom_fields":  nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *InterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var iface *netbox.Interface
	var httpResp *http.Response
	var err error

	// Look up interface by ID or by device+name
	switch {
	case !data.ID.IsNull() && data.ID.ValueString() != "":
		// Look up by ID
		var id int32
		if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &id); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Interface ID",
				fmt.Sprintf("Interface ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Looking up interface by ID", map[string]interface{}{
			"id": id,
		})

		iface, httpResp, err = d.client.DcimAPI.DcimInterfacesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading interface",
				utils.FormatAPIError(fmt.Sprintf("read interface ID %d", id), err, httpResp),
			)
			return
		}
	case !data.Device.IsNull() && data.Device.ValueString() != "" && !data.Name.IsNull() && data.Name.ValueString() != "":
		// Look up by device + name
		tflog.Debug(ctx, "Looking up interface by device and name", map[string]interface{}{
			"device": data.Device.ValueString(),
			"name":   data.Name.ValueString(),
		})

		// Try to resolve device identifier to ID
		deviceValue := data.Device.ValueString()
		var deviceID int32
		if _, parseErr := fmt.Sscanf(deviceValue, "%d", &deviceID); parseErr != nil {
			// Not a number, try to look up by name
			deviceList, deviceResp, deviceErr := d.client.DcimAPI.DcimDevicesList(ctx).Name([]string{deviceValue}).Execute()
			defer utils.CloseResponseBody(deviceResp)
			if deviceErr != nil {
				resp.Diagnostics.AddError(
					"Error looking up device",
					utils.FormatAPIError(fmt.Sprintf("look up device '%s'", deviceValue), deviceErr, deviceResp),
				)
				return
			}
			if len(deviceList.Results) == 0 {
				resp.Diagnostics.AddError(
					"Device not found",
					fmt.Sprintf("No device found with name '%s'", deviceValue),
				)
				return
			}
			if len(deviceList.Results) > 1 {
				resp.Diagnostics.AddError(
					"Multiple devices found",
					fmt.Sprintf("Multiple devices found with name '%s'. Please use device ID.", deviceValue),
				)
				return
			}
			deviceID = deviceList.Results[0].GetId()
		}

		// Look up interface by device ID and name
		list, listResp, listErr := d.client.DcimAPI.DcimInterfacesList(ctx).
			DeviceId([]int32{deviceID}).
			Name([]string{data.Name.ValueString()}).
			Execute()
		defer utils.CloseResponseBody(listResp)
		if listErr != nil {
			resp.Diagnostics.AddError(
				"Error reading interface",
				utils.FormatAPIError(fmt.Sprintf("list interfaces on device %d with name '%s'", deviceID, data.Name.ValueString()), listErr, listResp),
			)
			return
		}

		if len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"Interface not found",
				fmt.Sprintf("No interface found on device '%s' with name '%s'", deviceValue, data.Name.ValueString()),
			)
			return
		}

		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple interfaces found",
				fmt.Sprintf("Multiple interfaces found on device '%s' with name '%s'. Please use ID to identify uniquely.", deviceValue, data.Name.ValueString()),
			)
			return
		}

		iface = &list.Results[0]
	default:
		resp.Diagnostics.AddError(
			"Missing identifier",
			"You must specify either `id` or both `device` and `name`",
		)
		return
	}

	// Map the interface to the data source model
	d.mapInterfaceToDataSource(ctx, iface, &data, resp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapInterfaceToDataSource maps an Interface from the API to the data source model.
func (d *InterfaceDataSource) mapInterfaceToDataSource(ctx context.Context, iface *netbox.Interface, data *InterfaceDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", iface.GetId()))
	data.Name = types.StringValue(iface.GetName())

	// Device
	device := iface.GetDevice()
	data.Device = types.StringValue(fmt.Sprintf("%d", device.GetId()))

	// Type
	ifaceType := iface.GetType()
	if value, ok := ifaceType.GetValueOk(); ok && value != nil {
		data.Type = types.StringValue(string(*value))
	} else {
		data.Type = types.StringNull()
	}

	// Label
	if label, ok := iface.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Enabled
	if enabled, ok := iface.GetEnabledOk(); ok && enabled != nil {
		data.Enabled = types.BoolValue(*enabled)
	} else {
		data.Enabled = types.BoolValue(true)
	}

	// Parent
	if iface.HasParent() {
		parent := iface.GetParent()
		if parent.GetId() != 0 {
			data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		} else {
			data.Parent = types.StringNull()
		}
	} else {
		data.Parent = types.StringNull()
	}

	// Bridge
	if iface.HasBridge() {
		bridge := iface.GetBridge()
		if bridge.GetId() != 0 {
			data.Bridge = types.StringValue(fmt.Sprintf("%d", bridge.GetId()))
		} else {
			data.Bridge = types.StringNull()
		}
	} else {
		data.Bridge = types.StringNull()
	}

	// LAG
	if iface.HasLag() {
		lag := iface.GetLag()
		if lag.GetId() != 0 {
			data.Lag = types.StringValue(fmt.Sprintf("%d", lag.GetId()))
		} else {
			data.Lag = types.StringNull()
		}
	} else {
		data.Lag = types.StringNull()
	}

	// MTU
	if mtu, ok := iface.GetMtuOk(); ok && mtu != nil {
		data.Mtu = types.Int64Value(int64(*mtu))
	} else {
		data.Mtu = types.Int64Null()
	}

	// MAC Address
	if macAddr, ok := iface.GetMacAddressOk(); ok && macAddr != nil && *macAddr != "" {
		data.MacAddress = types.StringValue(*macAddr)
	} else {
		data.MacAddress = types.StringNull()
	}

	// Speed
	if speed, ok := iface.GetSpeedOk(); ok && speed != nil {
		data.Speed = types.Int64Value(int64(*speed))
	} else {
		data.Speed = types.Int64Null()
	}

	// Duplex
	if iface.HasDuplex() {
		duplex := iface.GetDuplex()
		if value, ok := duplex.GetValueOk(); ok && value != nil {
			data.Duplex = types.StringValue(string(*value))
		} else {
			data.Duplex = types.StringNull()
		}
	} else {
		data.Duplex = types.StringNull()
	}

	// WWN
	if wwn, ok := iface.GetWwnOk(); ok && wwn != nil && *wwn != "" {
		data.Wwn = types.StringValue(*wwn)
	} else {
		data.Wwn = types.StringNull()
	}

	// MgmtOnly
	if mgmtOnly, ok := iface.GetMgmtOnlyOk(); ok && mgmtOnly != nil {
		data.MgmtOnly = types.BoolValue(*mgmtOnly)
	} else {
		data.MgmtOnly = types.BoolValue(false)
	}

	// Description
	if desc, ok := iface.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Mode
	if iface.HasMode() {
		mode := iface.GetMode()
		if value, ok := mode.GetValueOk(); ok && value != nil {
			data.Mode = types.StringValue(string(*value))
		} else {
			data.Mode = types.StringNull()
		}
	} else {
		data.Mode = types.StringNull()
	}

	// MarkConnected
	if markConnected, ok := iface.GetMarkConnectedOk(); ok && markConnected != nil {
		data.MarkConnected = types.BoolValue(*markConnected)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Tags
	if iface.HasTags() && len(iface.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(iface.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields
	if iface.HasCustomFields() && len(iface.GetCustomFields()) > 0 {
		customFields := utils.MapToCustomFieldModels(iface.GetCustomFields(), nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(cfDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
