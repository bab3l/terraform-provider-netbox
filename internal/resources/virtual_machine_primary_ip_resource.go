// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &VirtualMachinePrimaryIPResource{}
	_ resource.ResourceWithConfigure   = &VirtualMachinePrimaryIPResource{}
	_ resource.ResourceWithImportState = &VirtualMachinePrimaryIPResource{}
)

// NewVirtualMachinePrimaryIPResource returns a new resource implementing the VM primary IP resource.
func NewVirtualMachinePrimaryIPResource() resource.Resource {
	return &VirtualMachinePrimaryIPResource{}
}

// VirtualMachinePrimaryIPResource defines the resource implementation.
type VirtualMachinePrimaryIPResource struct {
	client *netbox.APIClient
}

// VirtualMachinePrimaryIPResourceModel describes the resource data model.
type VirtualMachinePrimaryIPResourceModel struct {
	ID             types.String `tfsdk:"id"`
	VirtualMachine types.String `tfsdk:"virtual_machine"`
	PrimaryIP4     types.String `tfsdk:"primary_ip4"`
	PrimaryIP6     types.String `tfsdk:"primary_ip6"`
}

// Metadata returns the resource type name.
func (r *VirtualMachinePrimaryIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_primary_ip"
}

// Schema defines the schema for the resource.
func (r *VirtualMachinePrimaryIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages primary IP assignments for a virtual machine in NetBox. This resource is intended to avoid circular dependencies with VM interfaces and IP addresses.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the virtual machine.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_machine": nbschema.RequiredReferenceAttributeWithDiffSuppress("virtual machine", "ID or name of the virtual machine to update. Required."),
			"primary_ip4": schema.StringAttribute{
				MarkdownDescription: "Primary IPv4 address assigned to this virtual machine (ID or address).",
				Optional:            true,
			},
			"primary_ip6": schema.StringAttribute{
				MarkdownDescription: "Primary IPv6 address assigned to this virtual machine (ID or address).",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VirtualMachinePrimaryIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create sets primary IP assignments for a virtual machine.
func (r *VirtualMachinePrimaryIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualMachinePrimaryIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.IsSet(data.PrimaryIP4) && !utils.IsSet(data.PrimaryIP6) {
		resp.Diagnostics.AddError(
			"Missing primary IP assignment",
			"At least one of primary_ip4 or primary_ip6 must be set.",
		)
		return
	}

	vmID, diags := resolveVirtualMachineID(ctx, r.client, data.VirtualMachine.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patch := netbox.NewPatchedWritableVirtualMachineWithConfigContextRequest()
	applyVirtualMachinePrimaryIPPatch(ctx, r.client, patch, &data, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesPartialUpdate(ctx, vmID).
		PatchedWritableVirtualMachineWithConfigContextRequest(*patch).
		Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting virtual machine primary IPs",
			utils.FormatAPIError(fmt.Sprintf("update virtual machine ID %d", vmID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "set virtual machine primary IPs", httpResp, http.StatusOK) {
		return
	}

	r.mapToState(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *VirtualMachinePrimaryIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualMachinePrimaryIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Virtual Machine ID",
			fmt.Sprintf("Virtual Machine ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	vm, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx, vmID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() { resp.State.RemoveResource(ctx) }) {
			return
		}
		resp.Diagnostics.AddError(
			"Error reading virtual machine",
			utils.FormatAPIError(fmt.Sprintf("read virtual machine ID %d", vmID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read virtual machine", httpResp, http.StatusOK) {
		return
	}

	r.mapToState(vm, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates primary IP assignments for a virtual machine.
func (r *VirtualMachinePrimaryIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtualMachinePrimaryIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmID, diags := resolveVirtualMachineID(ctx, r.client, data.VirtualMachine.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patch := netbox.NewPatchedWritableVirtualMachineWithConfigContextRequest()
	applyVirtualMachinePrimaryIPPatch(ctx, r.client, patch, &data, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesPartialUpdate(ctx, vmID).
		PatchedWritableVirtualMachineWithConfigContextRequest(*patch).
		Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating virtual machine primary IPs",
			utils.FormatAPIError(fmt.Sprintf("update virtual machine ID %d", vmID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "update virtual machine primary IPs", httpResp, http.StatusOK) {
		return
	}

	r.mapToState(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete clears the primary IP assignments for a virtual machine.
func (r *VirtualMachinePrimaryIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualMachinePrimaryIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Virtual Machine ID",
			fmt.Sprintf("Virtual Machine ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	patch := netbox.NewPatchedWritableVirtualMachineWithConfigContextRequest()
	patch.SetPrimaryIp4Nil()
	patch.SetPrimaryIp6Nil()

	_, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesPartialUpdate(ctx, vmID).
		PatchedWritableVirtualMachineWithConfigContextRequest(*patch).
		Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, nil) {
			return
		}
		resp.Diagnostics.AddError(
			"Error clearing virtual machine primary IPs",
			utils.FormatAPIError(fmt.Sprintf("update virtual machine ID %d", vmID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "clear virtual machine primary IPs", httpResp, http.StatusOK) {
		return
	}
}

// ImportState imports an existing virtual machine primary IP assignment into Terraform.
func (r *VirtualMachinePrimaryIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

func (r *VirtualMachinePrimaryIPResource) mapToState(vm *netbox.VirtualMachineWithConfigContext, data *VirtualMachinePrimaryIPResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vm.GetId()))
	data.VirtualMachine = utils.UpdateReferenceAttribute(data.VirtualMachine, vm.GetName(), "", vm.GetId())

	if vm.HasPrimaryIp4() && vm.GetPrimaryIp4().Id != 0 {
		ip := vm.GetPrimaryIp4()
		data.PrimaryIP4 = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.PrimaryIP4 = types.StringNull()
	}

	if vm.HasPrimaryIp6() && vm.GetPrimaryIp6().Id != 0 {
		ip := vm.GetPrimaryIp6()
		data.PrimaryIP6 = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.PrimaryIP6 = types.StringNull()
	}
}

func applyVirtualMachinePrimaryIPPatch(ctx context.Context, client *netbox.APIClient, patch *netbox.PatchedWritableVirtualMachineWithConfigContextRequest, data *VirtualMachinePrimaryIPResourceModel, includeNull bool, diags *diag.Diagnostics) {
	if utils.IsSet(data.PrimaryIP4) {
		ipAddr, ipDiags := netboxlookup.LookupIPAddress(ctx, client, data.PrimaryIP4.ValueString())
		diags.Append(ipDiags...)
		if diags.HasError() {
			return
		}
		patch.SetPrimaryIp4(*ipAddr)
	} else if includeNull && data.PrimaryIP4.IsNull() {
		patch.SetPrimaryIp4Nil()
	}

	if utils.IsSet(data.PrimaryIP6) {
		ipAddr, ipDiags := netboxlookup.LookupIPAddress(ctx, client, data.PrimaryIP6.ValueString())
		diags.Append(ipDiags...)
		if diags.HasError() {
			return
		}
		patch.SetPrimaryIp6(*ipAddr)
	} else if includeNull && data.PrimaryIP6.IsNull() {
		patch.SetPrimaryIp6Nil()
	}
}

func resolveVirtualMachineID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {
	if id, err := utils.ParseID(value); err == nil {
		return id, nil
	}
	return netboxlookup.GenericLookupID(ctx, value, netboxlookup.VirtualMachineLookupConfig(client), func(vm *netbox.VirtualMachineWithConfigContext) int32 {
		return vm.GetId()
	})
}
