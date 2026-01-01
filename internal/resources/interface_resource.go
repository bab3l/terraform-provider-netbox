// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InterfaceResource{}
var _ resource.ResourceWithImportState = &InterfaceResource{}

func NewInterfaceResource() resource.Resource {
	return &InterfaceResource{}
}

// InterfaceResource defines the resource implementation.
type InterfaceResource struct {
	client *netbox.APIClient
}

// InterfaceResourceModel describes the resource data model.
type InterfaceResourceModel struct {
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

func (r *InterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface"
}

func (r *InterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an interface on a device in Netbox. Interfaces represent physical or virtual network interfaces on devices, including Ethernet ports, LAG interfaces, virtual interfaces, and more.",
		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("interface"),
			"device": schema.StringAttribute{
				MarkdownDescription: "ID or name of the device this interface belongs to. Required.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the interface (e.g., 'eth0', 'GigabitEthernet0/0'). Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label on the interface.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of interface. Common values: `virtual`, `bridge`, `lag`, `1000base-t`, `10gbase-t`, `10gbase-x-sfpp`, `25gbase-x-sfp28`, `40gbase-x-qsfpp`, `100gbase-x-qsfp28`. Required.",
				Required:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is enabled. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "ID of the parent interface (for sub-interfaces).",
				Optional:            true,
			},
			"bridge": schema.StringAttribute{
				MarkdownDescription: "ID of the bridge interface this interface belongs to.",
				Optional:            true,
			},
			"lag": schema.StringAttribute{
				MarkdownDescription: "ID of the LAG (Link Aggregation Group) this interface is a member of.",
				Optional:            true,
			},
			"mtu": schema.Int64Attribute{
				MarkdownDescription: "Maximum transmission unit (MTU) size. Common values: 1500 (Ethernet), 9000 (Jumbo frames).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65536),
				},
			},
			"mac_address": schema.StringAttribute{
				MarkdownDescription: "MAC address of the interface in format `AA:BB:CC:DD:EE:FF`.",
				Optional:            true,
			},
			"speed": schema.Int64Attribute{
				MarkdownDescription: "Interface speed in Kbps (e.g., 1000000 for 1Gbps, 10000000 for 10Gbps).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"duplex": schema.StringAttribute{
				MarkdownDescription: "Duplex mode. Valid values: `half`, `full`, `auto`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("half", "full", "auto", ""),
				},
			},
			"wwn": schema.StringAttribute{
				MarkdownDescription: "World Wide Name (WWN) for Fibre Channel interfaces.",
				Optional:            true,
			},
			"mgmt_only": schema.BoolAttribute{
				MarkdownDescription: "This interface is used only for out-of-band management.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "802.1Q mode. Valid values: `access`, `tagged`, `tagged-all`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("access", "tagged", "tagged-all", ""),
				},
			},
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected, even if no cable is attached.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("interface"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.
func (r *InterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new interface in Netbox.
func (r *InterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup the device
	deviceRef, diags := netboxlookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the interface request
	interfaceType := netbox.InterfaceTypeValue(data.Type.ValueString())
	interfaceReq := netbox.NewWritableInterfaceRequest(*deviceRef, data.Name.ValueString(), interfaceType)

	// Set optional fields
	r.setOptionalFields(ctx, interfaceReq, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply metadata fields (tags, custom_fields)
	utils.ApplyMetadataFields(ctx, interfaceReq, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating interface", map[string]interface{}{
		"device": data.Device.ValueString(),
		"name":   data.Name.ValueString(),
		"type":   data.Type.ValueString(),
	})
	iface, httpResp, err := r.client.DcimAPI.DcimInterfacesCreate(ctx).WritableInterfaceRequest(*interfaceReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating interface",
			utils.FormatAPIError("create interface", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapInterfaceToState(ctx, iface, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created interface", map[string]interface{}{
		"id":   iface.GetId(),
		"name": iface.GetName(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	interfaceID := utils.ParseInt32FromString(data.ID.ValueString())
	if interfaceID == 0 {
		resp.Diagnostics.AddError("Invalid Interface ID", "Interface ID must be a number.")
		return
	}
	tflog.Debug(ctx, "Reading interface", map[string]interface{}{
		"id": interfaceID,
	})
	iface, httpResp, err := r.client.DcimAPI.DcimInterfacesRetrieve(ctx, interfaceID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Interface not found, removing from state", map[string]interface{}{
				"id": interfaceID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading interface",
			utils.FormatAPIError("read interface", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapInterfaceToState(ctx, iface, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates an existing interface in Netbox.
func (r *InterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	interfaceID := utils.ParseInt32FromString(data.ID.ValueString())
	if interfaceID == 0 {
		resp.Diagnostics.AddError("Invalid Interface ID", "Interface ID must be a number.")
		return
	}

	// Lookup the device
	deviceRef, diags := netboxlookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the interface request
	interfaceType := netbox.InterfaceTypeValue(data.Type.ValueString())
	interfaceReq := netbox.NewWritableInterfaceRequest(*deviceRef, data.Name.ValueString(), interfaceType)

	// Set optional fields
	r.setOptionalFields(ctx, interfaceReq, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply metadata fields (tags, custom_fields)
	utils.ApplyMetadataFields(ctx, interfaceReq, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating interface", map[string]interface{}{
		"id":   interfaceID,
		"name": data.Name.ValueString(),
	})
	iface, httpResp, err := r.client.DcimAPI.DcimInterfacesUpdate(ctx, interfaceID).WritableInterfaceRequest(*interfaceReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating interface",
			utils.FormatAPIError("update interface", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapInterfaceToState(ctx, iface, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updated interface", map[string]interface{}{
		"id":   iface.GetId(),
		"name": iface.GetName(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes an interface from Netbox.
func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	interfaceID := utils.ParseInt32FromString(data.ID.ValueString())
	if interfaceID == 0 {
		resp.Diagnostics.AddError("Invalid Interface ID", "Interface ID must be a number.")
		return
	}
	tflog.Debug(ctx, "Deleting interface", map[string]interface{}{
		"id": interfaceID,
	})
	httpResp, err := r.client.DcimAPI.DcimInterfacesDestroy(ctx, interfaceID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Interface already deleted", map[string]interface{}{
				"id": interfaceID,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting interface",
			utils.FormatAPIError("delete interface", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted interface", map[string]interface{}{
		"id": interfaceID,
	})
}

// ImportState imports an existing interface into Terraform state.
func (r *InterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the interface request.
func (r *InterfaceResource) setOptionalFields(ctx context.Context, interfaceReq *netbox.WritableInterfaceRequest, data *InterfaceResourceModel, diags *diag.Diagnostics) {
	// Label
	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		label := data.Label.ValueString()
		interfaceReq.Label = &label
	}

	// Enabled
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		enabled := data.Enabled.ValueBool()
		interfaceReq.Enabled = &enabled
	}

	// Parent interface (WritableInterfaceRequest uses int32, not NestedInterfaceRequest)
	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID := utils.ParseInt32FromString(data.Parent.ValueString())
		if parentID != 0 {
			interfaceReq.Parent = *netbox.NewNullableInt32(&parentID)
		}
	}

	// Bridge interface (WritableInterfaceRequest uses int32, not NestedInterfaceRequest)
	if !data.Bridge.IsNull() && !data.Bridge.IsUnknown() {
		bridgeID := utils.ParseInt32FromString(data.Bridge.ValueString())
		if bridgeID != 0 {
			interfaceReq.Bridge = *netbox.NewNullableInt32(&bridgeID)
		}
	}

	// LAG interface (WritableInterfaceRequest uses int32, not NestedInterfaceRequest)
	if !data.Lag.IsNull() && !data.Lag.IsUnknown() {
		lagID := utils.ParseInt32FromString(data.Lag.ValueString())
		if lagID != 0 {
			interfaceReq.Lag = *netbox.NewNullableInt32(&lagID)
		}
	}

	// MTU
	if !data.Mtu.IsNull() && !data.Mtu.IsUnknown() {
		mtu, err := utils.SafeInt32FromValue(data.Mtu)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("Mtu value overflow: %s", err))
			return
		}
		interfaceReq.Mtu = *netbox.NewNullableInt32(&mtu)
	}

	// MAC Address
	if !data.MacAddress.IsNull() && !data.MacAddress.IsUnknown() {
		macAddr := data.MacAddress.ValueString()
		interfaceReq.MacAddress = *netbox.NewNullableString(&macAddr)
	}

	// Speed
	if !data.Speed.IsNull() && !data.Speed.IsUnknown() {
		speed, err := utils.SafeInt32FromValue(data.Speed)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("Speed value overflow: %s", err))
			return
		}
		interfaceReq.Speed = *netbox.NewNullableInt32(&speed)
	}

	// Duplex
	if !data.Duplex.IsNull() && !data.Duplex.IsUnknown() && data.Duplex.ValueString() != "" {
		duplex := netbox.InterfaceRequestDuplex(data.Duplex.ValueString())
		interfaceReq.Duplex = *netbox.NewNullableInterfaceRequestDuplex(&duplex)
	}

	// WWN
	if !data.Wwn.IsNull() && !data.Wwn.IsUnknown() {
		wwn := data.Wwn.ValueString()
		interfaceReq.Wwn = *netbox.NewNullableString(&wwn)
	}

	// MgmtOnly
	if !data.MgmtOnly.IsNull() && !data.MgmtOnly.IsUnknown() {
		mgmtOnly := data.MgmtOnly.ValueBool()
		interfaceReq.MgmtOnly = &mgmtOnly
	}

	// Description
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		interfaceReq.Description = &desc
	}

	// Mode (WritableInterfaceRequest uses PatchedWritableInterfaceRequestMode)
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() && data.Mode.ValueString() != "" {
		mode := netbox.PatchedWritableInterfaceRequestMode(data.Mode.ValueString())
		interfaceReq.Mode = &mode
	}

	// MarkConnected
	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		markConnected := data.MarkConnected.ValueBool()
		interfaceReq.MarkConnected = &markConnected
	}
}

// mapInterfaceToState maps a Netbox Interface to the Terraform state model.
func (r *InterfaceResource) mapInterfaceToState(ctx context.Context, iface *netbox.Interface, data *InterfaceResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", iface.GetId()))
	data.Name = types.StringValue(iface.GetName())

	// Device
	device := iface.GetDevice()
	userDevice := data.Device.ValueString()
	if userDevice == device.GetName() || userDevice == device.GetDisplay() || userDevice == fmt.Sprintf("%d", device.GetId()) {
		// Keep user's original value
	} else {
		data.Device = types.StringValue(device.GetName())
	}

	// Type
	ifaceType := iface.GetType()
	if value, ok := ifaceType.GetValueOk(); ok && value != nil {
		data.Type = types.StringValue(string(*value))
	}

	// Label
	if label, ok := iface.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Enabled - always set value since it's computed (defaults to true)
	if enabled, ok := iface.GetEnabledOk(); ok && enabled != nil {
		data.Enabled = types.BoolValue(*enabled)
	} else {
		// Set default value for computed field
		data.Enabled = types.BoolValue(true)
	}

	// Parent
	if iface.HasParent() {
		parent := iface.GetParent()
		if parent.GetId() != 0 {
			userParent := data.Parent.ValueString()
			if userParent == parent.GetName() || userParent == parent.GetDisplay() || userParent == fmt.Sprintf("%d", parent.GetId()) {
				// Keep user's original value
			} else {
				data.Parent = types.StringValue(parent.GetName())
			}
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
			userBridge := data.Bridge.ValueString()
			if userBridge == bridge.GetName() || userBridge == bridge.GetDisplay() || userBridge == fmt.Sprintf("%d", bridge.GetId()) {
				// Keep user's original value
			} else {
				data.Bridge = types.StringValue(bridge.GetName())
			}
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
			userLag := data.Lag.ValueString()
			if userLag == lag.GetName() || userLag == lag.GetDisplay() || userLag == fmt.Sprintf("%d", lag.GetId()) {
				// Keep user's original value
			} else {
				data.Lag = types.StringValue(lag.GetName())
			}
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
	if iface.HasTags() {
		tags := utils.NestedTagsToTagModels(iface.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if !tagDiags.HasError() {
			data.Tags = tagsValue
		} else {
			data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields
	if iface.HasCustomFields() {
		var existingCustomFields []utils.CustomFieldModel
		if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
			data.CustomFields.ElementsAs(ctx, &existingCustomFields, false)
		}
		customFields := utils.MapToCustomFieldModels(iface.GetCustomFields(), existingCustomFields)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
