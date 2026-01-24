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
	_ resource.Resource                = &DevicePrimaryIPResource{}
	_ resource.ResourceWithConfigure   = &DevicePrimaryIPResource{}
	_ resource.ResourceWithImportState = &DevicePrimaryIPResource{}
)

// NewDevicePrimaryIPResource returns a new resource implementing the device primary IP resource.
func NewDevicePrimaryIPResource() resource.Resource {
	return &DevicePrimaryIPResource{}
}

// DevicePrimaryIPResource defines the resource implementation.
type DevicePrimaryIPResource struct {
	client *netbox.APIClient
}

// DevicePrimaryIPResourceModel describes the resource data model.
type DevicePrimaryIPResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Device     types.String `tfsdk:"device"`
	PrimaryIP4 types.String `tfsdk:"primary_ip4"`
	PrimaryIP6 types.String `tfsdk:"primary_ip6"`
	OobIP      types.String `tfsdk:"oob_ip"`
}

// Metadata returns the resource type name.
func (r *DevicePrimaryIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_primary_ip"
}

// Schema defines the schema for the resource.
func (r *DevicePrimaryIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages primary IP assignments for a device in NetBox. This resource is intended to avoid circular dependencies with interfaces and IP addresses.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the device.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress("device", "ID or name of the device to update. Required."),
			"primary_ip4": schema.StringAttribute{
				MarkdownDescription: "Primary IPv4 address assigned to this device (ID or address).",
				Optional:            true,
			},
			"primary_ip6": schema.StringAttribute{
				MarkdownDescription: "Primary IPv6 address assigned to this device (ID or address).",
				Optional:            true,
			},
			"oob_ip": schema.StringAttribute{
				MarkdownDescription: "Out-of-band management IP address assigned to this device (ID or address).",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DevicePrimaryIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create sets primary IP assignments for a device.
func (r *DevicePrimaryIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DevicePrimaryIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.IsSet(data.PrimaryIP4) && !utils.IsSet(data.PrimaryIP6) && !utils.IsSet(data.OobIP) {
		resp.Diagnostics.AddError(
			"Missing primary IP assignment",
			"At least one of primary_ip4, primary_ip6, or oob_ip must be set.",
		)
		return
	}

	deviceID, diags := resolveDeviceID(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patch := netbox.NewPatchedWritableDeviceWithConfigContextRequest()
	applyDevicePrimaryIPPatch(ctx, r.client, patch, &data, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := r.client.DcimAPI.DcimDevicesPartialUpdate(ctx, deviceID).
		PatchedWritableDeviceWithConfigContextRequest(*patch).
		Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting device primary IPs",
			utils.FormatAPIError(fmt.Sprintf("update device ID %d", deviceID), err, httpResp),
		)
		return
	}

	r.mapToState(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *DevicePrimaryIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DevicePrimaryIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device ID",
			fmt.Sprintf("Device ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	device, httpResp, err := r.client.DcimAPI.DcimDevicesRetrieve(ctx, deviceID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading device",
			utils.FormatAPIError(fmt.Sprintf("read device ID %d", deviceID), err, httpResp),
		)
		return
	}

	r.mapToState(device, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates primary IP assignments for a device.
func (r *DevicePrimaryIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DevicePrimaryIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID, diags := resolveDeviceID(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patch := netbox.NewPatchedWritableDeviceWithConfigContextRequest()
	applyDevicePrimaryIPPatch(ctx, r.client, patch, &data, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := r.client.DcimAPI.DcimDevicesPartialUpdate(ctx, deviceID).
		PatchedWritableDeviceWithConfigContextRequest(*patch).
		Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device primary IPs",
			utils.FormatAPIError(fmt.Sprintf("update device ID %d", deviceID), err, httpResp),
		)
		return
	}

	r.mapToState(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete clears the primary IP assignments for a device.
func (r *DevicePrimaryIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DevicePrimaryIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device ID",
			fmt.Sprintf("Device ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	patch := netbox.NewPatchedWritableDeviceWithConfigContextRequest()
	patch.SetPrimaryIp4Nil()
	patch.SetPrimaryIp6Nil()
	patch.SetOobIpNil()

	_, httpResp, err := r.client.DcimAPI.DcimDevicesPartialUpdate(ctx, deviceID).
		PatchedWritableDeviceWithConfigContextRequest(*patch).
		Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error clearing device primary IPs",
			utils.FormatAPIError(fmt.Sprintf("update device ID %d", deviceID), err, httpResp),
		)
	}
}

// ImportState imports an existing device primary IP assignment into Terraform.
func (r *DevicePrimaryIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DevicePrimaryIPResource) mapToState(device *netbox.Device, data *DevicePrimaryIPResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", device.GetId()))
	data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())

	if device.HasPrimaryIp4() && device.GetPrimaryIp4().Id != 0 {
		ip := device.GetPrimaryIp4()
		data.PrimaryIP4 = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.PrimaryIP4 = types.StringNull()
	}

	if device.HasPrimaryIp6() && device.GetPrimaryIp6().Id != 0 {
		ip := device.GetPrimaryIp6()
		data.PrimaryIP6 = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.PrimaryIP6 = types.StringNull()
	}

	if device.HasOobIp() && device.GetOobIp().Id != 0 {
		ip := device.GetOobIp()
		data.OobIP = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.OobIP = types.StringNull()
	}
}

func applyDevicePrimaryIPPatch(ctx context.Context, client *netbox.APIClient, patch *netbox.PatchedWritableDeviceWithConfigContextRequest, data *DevicePrimaryIPResourceModel, includeNull bool, diags *diag.Diagnostics) {
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

	if utils.IsSet(data.OobIP) {
		ipAddr, ipDiags := netboxlookup.LookupIPAddress(ctx, client, data.OobIP.ValueString())
		diags.Append(ipDiags...)
		if diags.HasError() {
			return
		}
		patch.SetOobIp(*ipAddr)
	} else if includeNull && data.OobIP.IsNull() {
		patch.SetOobIpNil()
	}
}

func resolveDeviceID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {
	if id, err := utils.ParseID(value); err == nil {
		return id, nil
	}
	return netboxlookup.GenericLookupID(ctx, value, netboxlookup.DeviceLookupConfig(client), func(d *netbox.Device) int32 {
		return d.GetId()
	})
}
