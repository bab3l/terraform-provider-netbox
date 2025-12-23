// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &VMInterfaceDataSource{}

// NewVMInterfaceDataSource returns a new VM Interface data source.

func NewVMInterfaceDataSource() datasource.DataSource {

	return &VMInterfaceDataSource{}

}

// VMInterfaceDataSource defines the data source implementation.

type VMInterfaceDataSource struct {
	client *netbox.APIClient
}

// VMInterfaceDataSourceModel describes the data source data model.

type VMInterfaceDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	VirtualMachine types.String `tfsdk:"virtual_machine"`

	Name types.String `tfsdk:"name"`

	Enabled types.Bool `tfsdk:"enabled"`

	MTU types.Int64 `tfsdk:"mtu"`

	MACAddress types.String `tfsdk:"mac_address"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	Mode types.String `tfsdk:"mode"`

	UntaggedVLAN types.String `tfsdk:"untagged_vlan"`

	VRF types.String `tfsdk:"vrf"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *VMInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_vm_interface"

}

// Schema defines the schema for the data source.

func (d *VMInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a virtual machine interface in Netbox. You can identify the interface using `id`, or by specifying both `name` and `virtual_machine`.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.DSIDAttribute("VM interface"),

			"virtual_machine": schema.StringAttribute{

				MarkdownDescription: "The name of the virtual machine. Required when looking up by name.",

				Optional: true,

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the interface. Required when looking up by name (along with virtual_machine).",

				Optional: true,

				Computed: true,
			},

			"enabled": nbschema.DSComputedBoolAttribute("Whether the interface is enabled."),

			"mtu": nbschema.DSComputedInt64Attribute("The Maximum Transmission Unit (MTU) size for the interface."),

			"mac_address": nbschema.DSComputedStringAttribute("The MAC address of the interface."),

			"description": nbschema.DSComputedStringAttribute("Detailed description of the VM interface."),

			"display_name": nbschema.DSComputedStringAttribute("Display name for the VM interface."),

			"mode": nbschema.DSComputedStringAttribute("The 802.1Q mode of the interface (access, tagged, tagged-all)."),

			"untagged_vlan": nbschema.DSComputedStringAttribute("The untagged VLAN assigned to this interface."),

			"vrf": nbschema.DSComputedStringAttribute("The VRF assigned to this interface."),

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

// Configure sets up the data source with the provider client.

func (d *VMInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read retrieves data from the API.

func (d *VMInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data VMInterfaceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var iface *netbox.VMInterface

	var err error

	var httpResp *http.Response

	// Determine if we're searching by ID or by name+virtual_machine

	switch {

	case !data.ID.IsNull():

		// Search by ID

		ifaceID := data.ID.ValueString()

		tflog.Debug(ctx, "Reading VM interface by ID", map[string]interface{}{

			"id": ifaceID,
		})

		var ifaceIDInt int32

		if _, parseErr := fmt.Sscanf(ifaceID, "%d", &ifaceIDInt); parseErr != nil {

			resp.Diagnostics.AddError(

				"Invalid VM Interface ID",

				fmt.Sprintf("VM Interface ID must be a number, got: %s", ifaceID),
			)

			return

		}

		iface, httpResp, err = d.client.VirtualizationAPI.VirtualizationInterfacesRetrieve(ctx, ifaceIDInt).Execute()

		defer utils.CloseResponseBody(httpResp)

	case !data.Name.IsNull() && !data.VirtualMachine.IsNull():

		// Search by name and virtual machine

		ifaceName := data.Name.ValueString()

		vmName := data.VirtualMachine.ValueString()

		tflog.Debug(ctx, "Reading VM interface by name and virtual machine", map[string]interface{}{

			"name": ifaceName,

			"virtual_machine": vmName,
		})

		var ifaces *netbox.PaginatedVMInterfaceList

		ifaces, httpResp, err = d.client.VirtualizationAPI.VirtualizationInterfacesList(ctx).Name([]string{ifaceName}).VirtualMachine([]string{vmName}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading VM interface",

				utils.FormatAPIError("read VM interface by name", err, httpResp),
			)

			return

		}

		if len(ifaces.GetResults()) == 0 {

			resp.Diagnostics.AddError(

				"VM Interface Not Found",

				fmt.Sprintf("No VM interface found with name '%s' on virtual machine '%s'", ifaceName, vmName),
			)

			return

		}

		if len(ifaces.GetResults()) > 1 {

			resp.Diagnostics.AddError(

				"Multiple VM Interfaces Found",

				fmt.Sprintf("Multiple VM interfaces found with name '%s' on virtual machine '%s'.", ifaceName, vmName),
			)

			return

		}

		iface = &ifaces.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing VM Interface Identifier",

			"Either 'id' or both 'name' and 'virtual_machine' must be specified to identify the VM interface.",
		)

		return

	}

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading VM interface",

			utils.FormatAPIError("read VM interface", err, httpResp),
		)

		return

	}

	if httpResp != nil && httpResp.StatusCode == 404 {

		resp.Diagnostics.AddError(

			"VM Interface Not Found",

			fmt.Sprintf("No VM interface found with ID: %s", data.ID.ValueString()),
		)

		return

	}

	// Map response to state

	data.ID = types.StringValue(fmt.Sprintf("%d", iface.GetId()))

	data.Name = types.StringValue(iface.GetName())

	// Virtual Machine (always present - required field)

	data.VirtualMachine = types.StringValue(iface.VirtualMachine.GetName())

	// Enabled

	if iface.HasEnabled() {

		data.Enabled = types.BoolValue(iface.GetEnabled())

	} else {

		data.Enabled = types.BoolNull()

	}

	// MTU

	if iface.Mtu.IsSet() && iface.Mtu.Get() != nil {

		data.MTU = types.Int64Value(int64(*iface.Mtu.Get()))

	} else {

		data.MTU = types.Int64Null()

	}

	// MAC Address

	if iface.MacAddress.IsSet() && iface.MacAddress.Get() != nil && *iface.MacAddress.Get() != "" {

		data.MACAddress = types.StringValue(*iface.MacAddress.Get())

	} else {

		data.MACAddress = types.StringNull()

	}

	// Description

	if iface.HasDescription() && iface.GetDescription() != "" {

		data.Description = types.StringValue(iface.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Mode

	if iface.HasMode() {

		data.Mode = types.StringValue(string(iface.Mode.GetValue()))

	} else {

		data.Mode = types.StringNull()

	}

	// Untagged VLAN

	if iface.UntaggedVlan.IsSet() && iface.UntaggedVlan.Get() != nil {

		data.UntaggedVLAN = types.StringValue(iface.UntaggedVlan.Get().GetName())

	} else {

		data.UntaggedVLAN = types.StringNull()

	}

	// VRF

	if iface.Vrf.IsSet() && iface.Vrf.Get() != nil {

		data.VRF = types.StringValue(iface.Vrf.Get().GetName())

	} else {

		data.VRF = types.StringNull()

	}

	// Display name

	if displayName := iface.GetDisplay(); displayName != "" {

		data.DisplayName = types.StringValue(displayName)

	} else {

		data.DisplayName = types.StringNull()

	}

	// Handle tags

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

	// Handle custom fields

	if iface.HasCustomFields() {

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

	tflog.Debug(ctx, "Read VM interface", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
