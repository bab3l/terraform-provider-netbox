// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &VMInterfaceResource{}

	_ resource.ResourceWithConfigure = &VMInterfaceResource{}

	_ resource.ResourceWithImportState = &VMInterfaceResource{}

	_ resource.ResourceWithIdentity = &VMInterfaceResource{}
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

	Parent types.String `tfsdk:"parent"`

	Bridge types.String `tfsdk:"bridge"`

	UntaggedVLAN types.String `tfsdk:"untagged_vlan"`

	TaggedVLANs types.Set `tfsdk:"tagged_vlans"`

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
				Validators: []validator.String{
					validators.ValidMACAddress(),
				},
			},

			"mode": schema.StringAttribute{
				MarkdownDescription: "The 802.1Q mode of the interface. Valid values are: `access`, `tagged`, `tagged-all`.",

				Optional: true,
			},

			"parent": nbschema.ReferenceAttributeWithDiffSuppress(
				"parent interface",
				"Name or ID of the parent interface (for sub-interfaces).",
			),

			"bridge": nbschema.ReferenceAttributeWithDiffSuppress(
				"bridge interface",
				"Name or ID of the bridge interface this interface belongs to.",
			),

			"untagged_vlan": nbschema.ReferenceAttributeWithDiffSuppress(
				"vlan",
				"The name or ID of the untagged VLAN (for access or tagged mode).",
			),

			"tagged_vlans": schema.SetAttribute{
				MarkdownDescription: "Set of VLAN names or IDs to tag on this interface. Can only be set when mode is `tagged` or `tagged-all`.",
				Optional:            true,
				ElementType:         types.StringType,
			},

			"vrf": nbschema.ReferenceAttributeWithDiffSuppress(
				"vrf",
				"The name or ID of the VRF assigned to this interface.",
			),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("VM interface"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *VMInterfaceResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
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

	// Parent

	if iface.Parent.IsSet() && iface.Parent.Get() != nil {
		parent := iface.Parent.Get()
		data.Parent = utils.UpdateReferenceAttribute(data.Parent, parent.GetName(), "", parent.GetId())
	} else {
		data.Parent = types.StringNull()
	}

	// Bridge

	if iface.Bridge.IsSet() && iface.Bridge.Get() != nil {
		bridge := iface.Bridge.Get()
		data.Bridge = utils.UpdateReferenceAttribute(data.Bridge, bridge.GetName(), "", bridge.GetId())
	} else {
		data.Bridge = types.StringNull()
	}

	// Untagged VLAN

	if iface.UntaggedVlan.IsSet() && iface.UntaggedVlan.Get() != nil {
		vlan := iface.UntaggedVlan.Get()

		data.UntaggedVLAN = utils.UpdateReferenceAttribute(data.UntaggedVLAN, vlan.GetName(), "", vlan.GetId())
	} else {
		data.UntaggedVLAN = types.StringNull()
	}

	// Tagged VLANs

	data.TaggedVLANs = updateVMInterfaceTaggedVLANs(ctx, data.TaggedVLANs, iface.GetTaggedVlans(), diags)

	// VRF

	if iface.Vrf.IsSet() && iface.Vrf.Get() != nil {
		vrf := iface.Vrf.Get()

		data.VRF = utils.UpdateReferenceAttribute(data.VRF, vrf.GetName(), "", vrf.GetId())
	} else {
		data.VRF = types.StringNull()
	}

	// Handle tags (slug list) with empty-set preservation
	wasExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
	switch {
	case iface.HasTags() && len(iface.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(iface.GetTags()))
		for _, tag := range iface.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, iface.GetCustomFields(), diags)
}

// buildVMInterfaceRequest builds a WritableVMInterfaceRequest from the resource model.

func (r *VMInterfaceResource) buildVMInterfaceRequest(ctx context.Context, plan *VMInterfaceResourceModel, state *VMInterfaceResourceModel, diags *diag.Diagnostics) *netbox.WritableVMInterfaceRequest {
	// Lookup virtual machine (required)

	vm, vmDiags := netboxlookup.LookupVirtualMachine(ctx, r.client, plan.VirtualMachine.ValueString())

	diags.Append(vmDiags...)

	if diags.HasError() {
		return nil
	}

	ifaceRequest := &netbox.WritableVMInterfaceRequest{
		VirtualMachine: *vm,

		Name: plan.Name.ValueString(),
	}

	// Enabled

	if utils.IsSet(plan.Enabled) {
		enabled := plan.Enabled.ValueBool()

		ifaceRequest.Enabled = &enabled
	}

	// MTU

	if utils.IsSet(plan.MTU) {
		mtu, err := utils.SafeInt32FromValue(plan.MTU)

		if err != nil {
			diags.AddError("Invalid MTU value", fmt.Sprintf("MTU value overflow: %s", err))

			return nil
		}

		ifaceRequest.Mtu = *netbox.NewNullableInt32(&mtu)
	} else if plan.MTU.IsNull() && utils.IsSet(state.MTU) {
		// Only explicitly set to nil if we're clearing a previously set value
		ifaceRequest.SetMtuNil()
	}

	// MAC Address

	if utils.IsSet(plan.MACAddress) {
		macAddress := plan.MACAddress.ValueString()

		ifaceRequest.MacAddress = *netbox.NewNullableString(&macAddress)
	} else if plan.MACAddress.IsNull() && utils.IsSet(state.MACAddress) {
		// Only explicitly set to nil if we're clearing a previously set value
		ifaceRequest.SetMacAddressNil()
	}

	// Description

	utils.ApplyDescription(ifaceRequest, plan.Description)

	// Mode

	if utils.IsSet(plan.Mode) {
		mode := netbox.PatchedWritableInterfaceRequestMode(plan.Mode.ValueString())

		ifaceRequest.Mode = &mode
	} else if plan.Mode.IsNull() && utils.IsSet(state.Mode) {
		// Only explicitly set to empty string if we're clearing a previously set value
		emptyMode := netbox.PatchedWritableInterfaceRequestMode("")
		ifaceRequest.Mode = &emptyMode
	}

	// Parent

	if utils.IsSet(plan.Parent) {
		parentID, parentDiags := r.resolveVMInterfaceID(ctx, plan.Parent.ValueString())
		diags.Append(parentDiags...)
		if diags.HasError() {
			return nil
		}
		ifaceRequest.Parent = *netbox.NewNullableInt32(&parentID)
	} else if plan.Parent.IsNull() && utils.IsSet(state.Parent) {
		ifaceRequest.Parent.Set(nil)
	}

	// Bridge

	if utils.IsSet(plan.Bridge) {
		bridgeID, bridgeDiags := r.resolveVMInterfaceID(ctx, plan.Bridge.ValueString())
		diags.Append(bridgeDiags...)
		if diags.HasError() {
			return nil
		}
		ifaceRequest.Bridge = *netbox.NewNullableInt32(&bridgeID)
	} else if plan.Bridge.IsNull() && utils.IsSet(state.Bridge) {
		ifaceRequest.Bridge.Set(nil)
	}

	// Untagged VLAN

	if utils.IsSet(plan.UntaggedVLAN) {
		vlan, vlanDiags := netboxlookup.LookupVLAN(ctx, r.client, plan.UntaggedVLAN.ValueString())

		diags.Append(vlanDiags...)

		if diags.HasError() {
			return nil
		}

		ifaceRequest.UntaggedVlan = *netbox.NewNullableBriefVLANRequest(vlan)
	} else if plan.UntaggedVLAN.IsNull() && utils.IsSet(state.UntaggedVLAN) {
		// Only explicitly set to nil if we're clearing a previously set value
		ifaceRequest.SetUntaggedVlanNil()
	}

	// Tagged VLANs

	if utils.IsSet(plan.TaggedVLANs) {
		vlanIDs, vlanDiags := r.resolveTaggedVLANIDs(ctx, plan.TaggedVLANs)
		diags.Append(vlanDiags...)
		if diags.HasError() {
			return nil
		}
		ifaceRequest.TaggedVlans = vlanIDs
	} else if plan.TaggedVLANs.IsNull() && utils.IsSet(state.TaggedVLANs) {
		ifaceRequest.TaggedVlans = []int32{}
	}

	// VRF

	if utils.IsSet(plan.VRF) {
		vrf, vrfDiags := netboxlookup.LookupVRF(ctx, r.client, plan.VRF.ValueString())

		diags.Append(vrfDiags...)

		if diags.HasError() {
			return nil
		}

		ifaceRequest.Vrf = *netbox.NewNullableBriefVRFRequest(vrf)
	} else if plan.VRF.IsNull() && utils.IsSet(state.VRF) {
		// Only explicitly set to nil if we're clearing a previously set value
		ifaceRequest.SetVrfNil()
	}

	// Apply metadata fields individually with merge-aware helpers

	utils.ApplyTagsFromSlugs(ctx, r.client, ifaceRequest, plan.Tags, diags)
	utils.ApplyCustomFieldsWithMerge(ctx, ifaceRequest, plan.CustomFields, state.CustomFields, diags)

	if diags.HasError() {
		return nil
	}

	return ifaceRequest
}

func (r *VMInterfaceResource) resolveVMInterfaceID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	ifaceID, err := utils.ParseID(value)
	if err == nil {
		return ifaceID, nil
	}

	ifaces, httpResp, listErr := r.client.VirtualizationAPI.VirtualizationInterfacesList(ctx).Name([]string{value}).Execute()
	defer utils.CloseResponseBody(httpResp)
	if listErr != nil {
		return 0, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Error resolving VM interface",
			utils.FormatAPIError(fmt.Sprintf("lookup VM interface %s", value), listErr, httpResp),
		)}
	}

	results := ifaces.GetResults()
	if len(results) == 0 {
		return 0, diag.Diagnostics{diag.NewErrorDiagnostic(
			"VM interface not found",
			fmt.Sprintf("No VM interface found with name %q", value),
		)}
	}
	if len(results) > 1 {
		return 0, diag.Diagnostics{diag.NewErrorDiagnostic(
			"VM interface name not unique",
			fmt.Sprintf("Multiple VM interfaces found with name %q; use an ID instead", value),
		)}
	}

	return results[0].GetId(), nil
}

func (r *VMInterfaceResource) resolveTaggedVLANIDs(ctx context.Context, taggedVlans types.Set) ([]int32, diag.Diagnostics) {
	var vlanRefs []string
	var diags diag.Diagnostics
	if setDiags := taggedVlans.ElementsAs(ctx, &vlanRefs, false); setDiags.HasError() {
		diags.Append(setDiags...)
		return nil, diags
	}

	ids := make([]int32, 0, len(vlanRefs))
	for _, vlanRef := range vlanRefs {
		id, vlanDiags := netboxlookup.GenericLookupID(ctx, vlanRef, netboxlookup.VLANLookupConfig(r.client), func(v *netbox.VLAN) int32 {
			return v.GetId()
		})
		diags.Append(vlanDiags...)
		if diags.HasError() {
			return nil, diags
		}
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, diags
}

func updateVMInterfaceTaggedVLANs(ctx context.Context, current types.Set, vlans []netbox.VLAN, diags *diag.Diagnostics) types.Set {
	existing := map[string]struct{}{}
	if !current.IsNull() && !current.IsUnknown() {
		var currentValues []string
		diags.Append(current.ElementsAs(ctx, &currentValues, false)...)
		for _, value := range currentValues {
			existing[value] = struct{}{}
		}
	}

	if len(vlans) == 0 {
		if !current.IsNull() && !current.IsUnknown() && len(current.Elements()) == 0 {
			return types.SetValueMust(types.StringType, []attr.Value{})
		}
		return types.SetNull(types.StringType)
	}

	values := make([]string, 0, len(vlans))
	for _, vlan := range vlans {
		idValue := fmt.Sprintf("%d", vlan.GetId())
		nameValue := vlan.GetName()
		value := idValue
		if _, ok := existing[nameValue]; ok {
			value = nameValue
		}
		values = append(values, value)
	}

	sort.Strings(values)
	setValue, setDiags := types.SetValueFrom(ctx, types.StringType, values)
	diags.Append(setDiags...)
	return setValue
}

func validateVMInterfaceModeAndTaggedVLANs(ctx context.Context, plan *VMInterfaceResourceModel, diags *diag.Diagnostics) {
	if !utils.IsSet(plan.TaggedVLANs) {
		return
	}

	var vlanRefs []string
	if setDiags := plan.TaggedVLANs.ElementsAs(ctx, &vlanRefs, false); setDiags.HasError() {
		diags.Append(setDiags...)
		return
	}
	if len(vlanRefs) == 0 {
		return
	}

	mode := ""
	if utils.IsSet(plan.Mode) {
		mode = plan.Mode.ValueString()
	}

	if mode != "tagged" && mode != "tagged-all" {
		diags.AddError(
			"Invalid tagged_vlans configuration",
			"tagged_vlans can only be set when mode is \"tagged\" or \"tagged-all\".",
		)
	}
}

// Create creates a new VM interface resource.

func (r *VMInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VMInterfaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	validateVMInterfaceModeAndTaggedVLANs(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating VM interface", map[string]interface{}{
		"name": data.Name.ValueString(),

		"virtual_machine": data.VirtualMachine.ValueString(),
	})

	// Build the interface request

	var emptyState VMInterfaceResourceModel
	ifaceRequest := r.buildVMInterfaceRequest(ctx, &data, &emptyState, &resp.Diagnostics)

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

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.

func (r *VMInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VMInterfaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Preserve original custom_fields from state for potential restoration
	originalCustomFields := data.CustomFields

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

	// If custom_fields was null or empty before (not managed or explicitly cleared),
	// restore that state after mapping.
	// This prevents Terraform from trying to manage fields that aren't in the configuration.
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		tflog.Debug(ctx, "Custom fields unmanaged/cleared, preserving original state during Read", map[string]interface{}{
			"was_null":  originalCustomFields.IsNull(),
			"was_empty": !originalCustomFields.IsNull() && len(originalCustomFields.Elements()) == 0,
		})
		data.CustomFields = originalCustomFields
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.

func (r *VMInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan VMInterfaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	validateVMInterfaceModeAndTaggedVLANs(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	ifaceID := plan.ID.ValueString()

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

		"name": plan.Name.ValueString(),
	})

	// Build the interface request

	ifaceRequest := r.buildVMInterfaceRequest(ctx, &plan, &state, &resp.Diagnostics)

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

	// Store plan's custom_fields to filter the response
	planCustomFields := plan.CustomFields

	// Map response to state

	r.mapVMInterfaceToState(ctx, iface, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Filter custom_fields to only those owned by this resource
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, iface.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		ifaceIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid VM Interface ID",
				fmt.Sprintf("VM Interface ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		iface, httpResp, err := r.client.VirtualizationAPI.VirtualizationInterfacesRetrieve(ctx, ifaceIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing VM interface", utils.FormatAPIError(fmt.Sprintf("read VM interface ID %s", parsed.ID), err, httpResp))
			return
		}

		var data VMInterfaceResourceModel
		data.Mode = types.StringUnknown()
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapVMInterfaceToState(ctx, iface, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, iface.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)

	// Mark mode as unknown during import so that Read method knows to populate it
	// This allows import to preserve mode field while regular reads without mode in config don't cause drift
	resp.State.SetAttribute(ctx, path.Root("mode"), types.StringUnknown())
}
