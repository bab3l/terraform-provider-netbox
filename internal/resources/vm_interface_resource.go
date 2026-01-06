// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &VMInterfaceResource{}

	_ resource.ResourceWithConfigure = &VMInterfaceResource{}

	_ resource.ResourceWithImportState = &VMInterfaceResource{}
)

// NewVMInterfaceResource returns a new VM Interface resource.

func NewVMInterfaceResource() resource.Resource {
	return &VMInterfaceResource{}
}

// VMInterfaceResource defines the resource implementation.

type VMInterfaceResource struct {
	client *netbox.APIClient
}

// VMInterfaceResourceModel describes the resource data model.

type VMInterfaceResourceModel struct {
	ID types.String `tfsdk:"id"`

	VirtualMachine types.String `tfsdk:"virtual_machine"`

	Name types.String `tfsdk:"name"`

	Enabled types.Bool `tfsdk:"enabled"`

	MTU types.Int64 `tfsdk:"mtu"`

	MACAddress types.String `tfsdk:"mac_address"`

	Description types.String `tfsdk:"description"`

	Mode types.String `tfsdk:"mode"`

	UntaggedVLAN types.String `tfsdk:"untagged_vlan"`

	VRF types.String `tfsdk:"vrf"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *VMInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_interface"
}

// Schema defines the schema for the resource.

func (r *VMInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual machine interface in Netbox. VM interfaces represent the virtual network interfaces attached to a virtual machine.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the VM interface.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"virtual_machine": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the virtual machine this interface belongs to.",

				Required: true,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the interface (e.g., 'eth0', 'ens192').",

				Required: true,
			},

			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is enabled. Defaults to true.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(true),
			},

			"mtu": schema.Int64Attribute{
				MarkdownDescription: "The Maximum Transmission Unit (MTU) size for the interface.",

				Optional: true,
			},

			"mac_address": schema.StringAttribute{
				MarkdownDescription: "The MAC address of the interface.",

				Optional: true,
			},

			"mode": schema.StringAttribute{
				MarkdownDescription: "The 802.1Q mode of the interface. Valid values are: `access`, `tagged`, `tagged-all`.",

				Optional: true,
			},

			"untagged_vlan": nbschema.ReferenceAttributeWithDiffSuppress(
				"vlan",
				"The name or ID of the untagged VLAN (for access or tagged mode).",
			),

			"vrf": nbschema.ReferenceAttributeWithDiffSuppress(
				"vrf",
				"The name or ID of the VRF assigned to this interface.",
			),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("VM interface"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure sets up the resource with the provider client.
func (r *VMInterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// mapVMInterfaceToState maps a VMInterface from the API to the Terraform state model.

func (r *VMInterfaceResource) mapVMInterfaceToState(ctx context.Context, iface *netbox.VMInterface, data *VMInterfaceResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", iface.GetId()))

	data.Name = types.StringValue(iface.GetName())

	// Virtual Machine (always present - required field)

	data.VirtualMachine = utils.UpdateReferenceAttribute(data.VirtualMachine, iface.VirtualMachine.GetName(), "", iface.VirtualMachine.GetId())

	// Enabled

	if iface.HasEnabled() {
		data.Enabled = types.BoolValue(iface.GetEnabled())
	} else {
		data.Enabled = types.BoolValue(true)
	}

	// MTU

	if iface.Mtu.IsSet() && iface.Mtu.Get() != nil {
		data.MTU = types.Int64Value(int64(*iface.Mtu.Get()))
	} else {
		data.MTU = types.Int64Null()
	}

	// MAC Address

	if iface.MacAddress.IsSet() && iface.MacAddress.Get() != nil && *iface.MacAddress.Get() != "" {
		apiMac := *iface.MacAddress.Get()

		if !data.MACAddress.IsNull() && !data.MACAddress.IsUnknown() {
			if strings.EqualFold(data.MACAddress.ValueString(), apiMac) {
				// Keep user's casing

				apiMac = data.MACAddress.ValueString()
			}
		}

		data.MACAddress = types.StringValue(apiMac)
	} else {
		data.MACAddress = types.StringNull()
	}

	// Description

	if iface.HasDescription() && iface.GetDescription() != "" {
		data.Description = types.StringValue(iface.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Mode - only set if user specified it in config, or during import
	// This prevents Terraform from seeing drift when the API returns a mode but config doesn't specify one
	// During import, data.ID would be unknown initially, but by the time we reach mapVMInterfaceToState,
	// the ID has been set. Instead, we check if mode is unknown (happens during import state refresh)

	if !data.Mode.IsNull() || data.Mode.IsUnknown() {
		if iface.HasMode() {
			data.Mode = types.StringValue(string(iface.Mode.GetValue()))
		} else {
			data.Mode = types.StringNull()
		}
	}

	// Untagged VLAN

	if iface.UntaggedVlan.IsSet() && iface.UntaggedVlan.Get() != nil {
		vlan := iface.UntaggedVlan.Get()

		data.UntaggedVLAN = utils.UpdateReferenceAttribute(data.UntaggedVLAN, vlan.GetName(), "", vlan.GetId())
	} else {
		data.UntaggedVLAN = types.StringNull()
	}

	// VRF

	if iface.Vrf.IsSet() && iface.Vrf.Get() != nil {
		vrf := iface.Vrf.Get()

		data.VRF = utils.UpdateReferenceAttribute(data.VRF, vrf.GetName(), "", vrf.GetId())
	} else {
		data.VRF = types.StringNull()
	}

	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsFromAPI(ctx, iface.HasTags(), iface.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, iface.HasCustomFields(), iface.GetCustomFields(), data.CustomFields, diags)
}

// buildVMInterfaceRequest builds a WritableVMInterfaceRequest from the resource model.

func (r *VMInterfaceResource) buildVMInterfaceRequest(ctx context.Context, data *VMInterfaceResourceModel, diags *diag.Diagnostics) *netbox.WritableVMInterfaceRequest {
	// Lookup virtual machine (required)

	vm, vmDiags := netboxlookup.LookupVirtualMachine(ctx, r.client, data.VirtualMachine.ValueString())

	diags.Append(vmDiags...)

	if diags.HasError() {
		return nil
	}

	ifaceRequest := &netbox.WritableVMInterfaceRequest{
		VirtualMachine: *vm,

		Name: data.Name.ValueString(),
	}

	// Enabled

	if utils.IsSet(data.Enabled) {
		enabled := data.Enabled.ValueBool()

		ifaceRequest.Enabled = &enabled
	}

	// MTU

	if utils.IsSet(data.MTU) {
		mtu, err := utils.SafeInt32FromValue(data.MTU)

		if err != nil {
			diags.AddError("Invalid MTU value", fmt.Sprintf("MTU value overflow: %s", err))

			return nil
		}

		ifaceRequest.Mtu = *netbox.NewNullableInt32(&mtu)
	}

	// MAC Address

	if utils.IsSet(data.MACAddress) {
		macAddress := data.MACAddress.ValueString()

		ifaceRequest.MacAddress = *netbox.NewNullableString(&macAddress)
	}

	// Description

	if utils.IsSet(data.Description) {
		description := data.Description.ValueString()

		ifaceRequest.Description = &description
	}

	// Mode

	if utils.IsSet(data.Mode) {
		mode := netbox.PatchedWritableInterfaceRequestMode(data.Mode.ValueString())

		ifaceRequest.Mode = &mode
	}

	// Untagged VLAN

	if utils.IsSet(data.UntaggedVLAN) {
		vlan, vlanDiags := netboxlookup.LookupVLAN(ctx, r.client, data.UntaggedVLAN.ValueString())

		diags.Append(vlanDiags...)

		if diags.HasError() {
			return nil
		}

		ifaceRequest.UntaggedVlan = *netbox.NewNullableBriefVLANRequest(vlan)
	}

	// VRF

	if utils.IsSet(data.VRF) {
		vrf, vrfDiags := netboxlookup.LookupVRF(ctx, r.client, data.VRF.ValueString())

		diags.Append(vrfDiags...)

		if diags.HasError() {
			return nil
		}

		ifaceRequest.Vrf = *netbox.NewNullableBriefVRFRequest(vrf)
	}

	// Apply metadata fields (tags, custom_fields)

	utils.ApplyMetadataFields(ctx, ifaceRequest, data.Tags, data.CustomFields, diags)

	if diags.HasError() {
		return nil
	}

	return ifaceRequest
}

// Create creates a new VM interface resource.

func (r *VMInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VMInterfaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating VM interface", map[string]interface{}{
		"name": data.Name.ValueString(),

		"virtual_machine": data.VirtualMachine.ValueString(),
	})

	// Build the interface request

	ifaceRequest := r.buildVMInterfaceRequest(ctx, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API

	iface, httpResp, err := r.client.VirtualizationAPI.VirtualizationInterfacesCreate(ctx).WritableVMInterfaceRequest(*ifaceRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating VM interface",

			utils.FormatAPIError("create VM interface", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Created VM interface", map[string]interface{}{
		"id": iface.GetId(),

		"name": iface.GetName(),
	})

	// Map response to state

	r.mapVMInterfaceToState(ctx, iface, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.

func (r *VMInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VMInterfaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	ifaceID := data.ID.ValueString()

	var ifaceIDInt int32

	ifaceIDInt, err := utils.ParseID(ifaceID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid VM Interface ID",

			fmt.Sprintf("VM Interface ID must be a number, got: %s", ifaceID),
		)

		return
	}

	tflog.Debug(ctx, "Reading VM interface", map[string]interface{}{
		"id": ifaceID,
	})

	// Call the API

	iface, httpResp, err := r.client.VirtualizationAPI.VirtualizationInterfacesRetrieve(ctx, ifaceIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "VM interface not found, removing from state", map[string]interface{}{
				"id": ifaceID,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading VM interface",

			utils.FormatAPIError(fmt.Sprintf("read VM interface ID %s", ifaceID), err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapVMInterfaceToState(ctx, iface, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.

func (r *VMInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VMInterfaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	ifaceID := data.ID.ValueString()

	var ifaceIDInt int32

	ifaceIDInt, err := utils.ParseID(ifaceID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid VM Interface ID",

			fmt.Sprintf("VM Interface ID must be a number, got: %s", ifaceID),
		)

		return
	}

	tflog.Debug(ctx, "Updating VM interface", map[string]interface{}{
		"id": ifaceID,

		"name": data.Name.ValueString(),
	})

	// Build the interface request

	ifaceRequest := r.buildVMInterfaceRequest(ctx, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API

	iface, httpResp, err := r.client.VirtualizationAPI.VirtualizationInterfacesUpdate(ctx, ifaceIDInt).WritableVMInterfaceRequest(*ifaceRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating VM interface",

			utils.FormatAPIError(fmt.Sprintf("update VM interface ID %s", ifaceID), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Updated VM interface", map[string]interface{}{
		"id": iface.GetId(),

		"name": iface.GetName(),
	})

	// Map response to state

	r.mapVMInterfaceToState(ctx, iface, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state.

func (r *VMInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VMInterfaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	ifaceID := data.ID.ValueString()

	var ifaceIDInt int32

	ifaceIDInt, err := utils.ParseID(ifaceID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid VM Interface ID",

			fmt.Sprintf("VM Interface ID must be a number, got: %s", ifaceID),
		)

		return
	}

	tflog.Debug(ctx, "Deleting VM interface", map[string]interface{}{
		"id": ifaceID,
	})

	// Call the API

	httpResp, err := r.client.VirtualizationAPI.VirtualizationInterfacesDestroy(ctx, ifaceIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted, consider success

			tflog.Debug(ctx, "VM interface already deleted", map[string]interface{}{
				"id": ifaceID,
			})

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting VM interface",

			utils.FormatAPIError(fmt.Sprintf("delete VM interface ID %s", ifaceID), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted VM interface", map[string]interface{}{
		"id": ifaceID,
	})
}

// ImportState imports an existing resource into Terraform.

func (r *VMInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// Mark mode as unknown during import so that Read method knows to populate it
	// This allows import to preserve mode field while regular reads without mode in config don't cause drift
	resp.State.SetAttribute(ctx, path.Root("mode"), types.StringUnknown())
}
