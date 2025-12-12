// Package resources contains Terraform resource implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to create, read, update,
// and delete resources in Netbox via Terraform.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TunnelResource{}
var _ resource.ResourceWithImportState = &TunnelResource{}

func NewTunnelResource() resource.Resource {
	return &TunnelResource{}
}

// TunnelResource defines the resource implementation.
type TunnelResource struct {
	client *netbox.APIClient
}

// TunnelResourceModel describes the resource data model.
type TunnelResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Status        types.String `tfsdk:"status"`
	Group         types.String `tfsdk:"group"`
	Encapsulation types.String `tfsdk:"encapsulation"`
	IPSecProfile  types.String `tfsdk:"ipsec_profile"`
	Tenant        types.String `tfsdk:"tenant"`
	TunnelID      types.Int64  `tfsdk:"tunnel_id"`
	Description   types.String `tfsdk:"description"`
	Comments      types.String `tfsdk:"comments"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

func (r *TunnelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunnel"
}

func (r *TunnelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a VPN tunnel in Netbox. Tunnels represent secure connections between network endpoints using various encapsulation protocols like IPSec, GRE, WireGuard, etc.",

		Attributes: map[string]schema.Attribute{
			"id":   nbschema.IDAttribute("tunnel"),
			"name": nbschema.NameAttribute("tunnel", 100),
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the tunnel. Valid values: `planned`, `active`, `disabled`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"planned",
						"active",
						"disabled",
					),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "ID of the tunnel group this tunnel belongs to.",
				Optional:            true,
			},
			"encapsulation": schema.StringAttribute{
				MarkdownDescription: "Encapsulation protocol for the tunnel. Valid values: `ipsec-transport`, `ipsec-tunnel`, `ip-ip`, `gre`, `wireguard`, `openvpn`, `l2tp`, `pptp`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ipsec-transport",
						"ipsec-tunnel",
						"ip-ip",
						"gre",
						"wireguard",
						"openvpn",
						"l2tp",
						"pptp",
					),
				},
			},
			"ipsec_profile": schema.StringAttribute{
				MarkdownDescription: "ID of the IPSec profile for this tunnel (required for IPSec encapsulation types).",
				Optional:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant this tunnel belongs to.",
				Optional:            true,
			},
			"tunnel_id": schema.Int64Attribute{
				MarkdownDescription: "Tunnel identifier (numeric ID used by the tunnel protocol).",
				Optional:            true,
			},
			"description":   nbschema.DescriptionAttribute("tunnel"),
			"comments":      nbschema.CommentsAttribute("tunnel"),
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *TunnelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *TunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TunnelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request using WritableTunnelRequest
	tunnelRequest := netbox.NewWritableTunnelRequest(
		data.Name.ValueString(),
		netbox.PatchedWritableTunnelRequestEncapsulation(data.Encapsulation.ValueString()),
	)

	// Set status - default to "active" if not provided (Netbox requires status)
	statusValue := "active"
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		statusValue = data.Status.ValueString()
	}
	status := netbox.PatchedWritableTunnelRequestStatus(statusValue)
	tunnelRequest.Status = &status

	// Initialize additional properties for integer ID references
	tunnelRequest.AdditionalProperties = make(map[string]interface{})

	// Set optional fields
	if !data.Group.IsNull() && data.Group.ValueString() != "" {
		groupID, err := utils.ParseID(data.Group.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Tunnel Group ID",
				fmt.Sprintf("Unable to parse tunnel group ID: %s", err),
			)
			return
		}
		tunnelRequest.AdditionalProperties["group"] = int(groupID)
	}

	if !data.IPSecProfile.IsNull() && data.IPSecProfile.ValueString() != "" {
		ipsecProfileID, err := utils.ParseID(data.IPSecProfile.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid IPSec Profile ID",
				fmt.Sprintf("Unable to parse IPSec profile ID: %s", err),
			)
			return
		}
		tunnelRequest.AdditionalProperties["ipsec_profile"] = int(ipsecProfileID)
	}

	if !data.Tenant.IsNull() && data.Tenant.ValueString() != "" {
		tenantID, err := utils.ParseID(data.Tenant.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Tenant ID",
				fmt.Sprintf("Unable to parse tenant ID: %s", err),
			)
			return
		}
		tunnelRequest.AdditionalProperties["tenant"] = int(tenantID)
	}

	if !data.TunnelID.IsNull() {
		tunnelRequest.TunnelId = *netbox.NewNullableInt64(netbox.PtrInt64(data.TunnelID.ValueInt64()))
	}

	if !data.Description.IsNull() && data.Description.ValueString() != "" {
		tunnelRequest.Description = netbox.PtrString(data.Description.ValueString())
	}

	if !data.Comments.IsNull() && data.Comments.ValueString() != "" {
		tunnelRequest.Comments = netbox.PtrString(data.Comments.ValueString())
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tagModels []utils.TagModel
		diags := data.Tags.ElementsAs(ctx, &tagModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelRequest.Tags = utils.TagsToNestedTagRequests(tagModels)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() {
		var customFieldModels []utils.CustomFieldModel
		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelRequest.CustomFields = utils.CustomFieldsToMap(customFieldModels)
	}

	tflog.Debug(ctx, "Creating tunnel", map[string]interface{}{
		"name":          data.Name.ValueString(),
		"encapsulation": data.Encapsulation.ValueString(),
	})

	// Create the tunnel via API
	tunnel, httpResp, err := r.client.VpnAPI.VpnTunnelsCreate(ctx).WritableTunnelRequest(*tunnelRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tunnel",
			utils.FormatAPIError("create tunnel", err, httpResp),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", tunnel.GetId()))
	data.Status = types.StringValue(string(tunnel.Status.GetValue()))

	// Handle group reference from response
	if tunnel.HasGroup() && tunnel.Group.IsSet() && tunnel.Group.Get() != nil {
		data.Group = types.StringValue(fmt.Sprintf("%d", tunnel.Group.Get().GetId()))
	}

	// Handle ipsec_profile reference from response
	if tunnel.HasIpsecProfile() && tunnel.IpsecProfile.IsSet() && tunnel.IpsecProfile.Get() != nil {
		data.IPSecProfile = types.StringValue(fmt.Sprintf("%d", tunnel.IpsecProfile.Get().GetId()))
	}

	// Handle tenant reference from response
	if tunnel.HasTenant() && tunnel.Tenant.IsSet() && tunnel.Tenant.Get() != nil {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tunnel.Tenant.Get().GetId()))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TunnelResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	tunnelID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tunnel ID",
			fmt.Sprintf("Unable to parse tunnel ID: %s", err),
		)
		return
	}

	// Get the tunnel from the API
	tunnel, httpResp, err := r.client.VpnAPI.VpnTunnelsRetrieve(ctx, tunnelID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Tunnel not found, removing from state", map[string]interface{}{
				"id": data.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading tunnel",
			utils.FormatAPIError("read tunnel", err, httpResp),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", tunnel.GetId()))
	data.Name = types.StringValue(tunnel.GetName())
	data.Status = types.StringValue(string(tunnel.Status.GetValue()))
	data.Encapsulation = types.StringValue(string(tunnel.Encapsulation.GetValue()))

	// Handle group reference
	if tunnel.HasGroup() && tunnel.Group.IsSet() && tunnel.Group.Get() != nil {
		data.Group = types.StringValue(fmt.Sprintf("%d", tunnel.Group.Get().GetId()))
	} else {
		data.Group = types.StringNull()
	}

	// Handle ipsec_profile reference
	if tunnel.HasIpsecProfile() && tunnel.IpsecProfile.IsSet() && tunnel.IpsecProfile.Get() != nil {
		data.IPSecProfile = types.StringValue(fmt.Sprintf("%d", tunnel.IpsecProfile.Get().GetId()))
	} else {
		data.IPSecProfile = types.StringNull()
	}

	// Handle tenant reference
	if tunnel.HasTenant() && tunnel.Tenant.IsSet() && tunnel.Tenant.Get() != nil {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tunnel.Tenant.Get().GetId()))
	} else {
		data.Tenant = types.StringNull()
	}

	// Handle tunnel_id
	if tunnel.HasTunnelId() && tunnel.TunnelId.IsSet() {
		if val := tunnel.TunnelId.Get(); val != nil {
			data.TunnelID = types.Int64Value(*val)
		} else {
			data.TunnelID = types.Int64Null()
		}
	} else {
		data.TunnelID = types.Int64Null()
	}

	// Handle description
	if tunnel.HasDescription() && tunnel.GetDescription() != "" {
		data.Description = types.StringValue(tunnel.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments
	if tunnel.HasComments() && tunnel.GetComments() != "" {
		data.Comments = types.StringValue(tunnel.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if tunnel.HasTags() {
		tags := utils.NestedTagsToTagModels(tunnel.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if tunnel.HasCustomFields() {
		// Get existing custom field models to preserve type information
		var existingCustomFields []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			diags := data.CustomFields.ElementsAs(ctx, &existingCustomFields, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		customFields := utils.MapToCustomFieldModels(tunnel.GetCustomFields(), existingCustomFields)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunnelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TunnelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	tunnelID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tunnel ID",
			fmt.Sprintf("Unable to parse tunnel ID: %s", err),
		)
		return
	}

	// Build the API request using WritableTunnelRequest
	tunnelRequest := netbox.NewWritableTunnelRequest(
		data.Name.ValueString(),
		netbox.PatchedWritableTunnelRequestEncapsulation(data.Encapsulation.ValueString()),
	)

	// Set status - default to "active" if not provided (Netbox requires status)
	statusValue := "active"
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		statusValue = data.Status.ValueString()
	}
	status := netbox.PatchedWritableTunnelRequestStatus(statusValue)
	tunnelRequest.Status = &status

	// Initialize additional properties for integer ID references
	tunnelRequest.AdditionalProperties = make(map[string]interface{})

	if !data.Group.IsNull() && data.Group.ValueString() != "" {
		groupID, err := utils.ParseID(data.Group.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Tunnel Group ID",
				fmt.Sprintf("Unable to parse tunnel group ID: %s", err),
			)
			return
		}
		tunnelRequest.AdditionalProperties["group"] = int(groupID)
	} else {
		tunnelRequest.AdditionalProperties["group"] = nil
	}

	if !data.IPSecProfile.IsNull() && data.IPSecProfile.ValueString() != "" {
		ipsecProfileID, err := utils.ParseID(data.IPSecProfile.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid IPSec Profile ID",
				fmt.Sprintf("Unable to parse IPSec profile ID: %s", err),
			)
			return
		}
		tunnelRequest.AdditionalProperties["ipsec_profile"] = int(ipsecProfileID)
	} else {
		tunnelRequest.AdditionalProperties["ipsec_profile"] = nil
	}

	if !data.Tenant.IsNull() && data.Tenant.ValueString() != "" {
		tenantID, err := utils.ParseID(data.Tenant.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Tenant ID",
				fmt.Sprintf("Unable to parse tenant ID: %s", err),
			)
			return
		}
		tunnelRequest.AdditionalProperties["tenant"] = int(tenantID)
	} else {
		tunnelRequest.AdditionalProperties["tenant"] = nil
	}

	if !data.TunnelID.IsNull() {
		tunnelRequest.TunnelId = *netbox.NewNullableInt64(netbox.PtrInt64(data.TunnelID.ValueInt64()))
	} else {
		tunnelRequest.TunnelId = *netbox.NewNullableInt64(nil)
	}

	if !data.Description.IsNull() && data.Description.ValueString() != "" {
		tunnelRequest.Description = netbox.PtrString(data.Description.ValueString())
	} else {
		tunnelRequest.Description = netbox.PtrString("")
	}

	if !data.Comments.IsNull() && data.Comments.ValueString() != "" {
		tunnelRequest.Comments = netbox.PtrString(data.Comments.ValueString())
	} else {
		tunnelRequest.Comments = netbox.PtrString("")
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tagModels []utils.TagModel
		diags := data.Tags.ElementsAs(ctx, &tagModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelRequest.Tags = utils.TagsToNestedTagRequests(tagModels)
	} else {
		tunnelRequest.Tags = []netbox.NestedTagRequest{}
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() {
		var customFieldModels []utils.CustomFieldModel
		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelRequest.CustomFields = utils.CustomFieldsToMap(customFieldModels)
	} else {
		tunnelRequest.CustomFields = map[string]interface{}{}
	}

	tflog.Debug(ctx, "Updating tunnel", map[string]interface{}{
		"id":            data.ID.ValueString(),
		"name":          data.Name.ValueString(),
		"encapsulation": data.Encapsulation.ValueString(),
	})

	// Update the tunnel via API
	tunnel, httpResp, err := r.client.VpnAPI.VpnTunnelsUpdate(ctx, tunnelID).WritableTunnelRequest(*tunnelRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tunnel",
			utils.FormatAPIError("update tunnel", err, httpResp),
		)
		return
	}

	// Update status from response
	data.Status = types.StringValue(string(tunnel.Status.GetValue()))

	// Handle group reference from response
	if tunnel.HasGroup() && tunnel.Group.IsSet() && tunnel.Group.Get() != nil {
		data.Group = types.StringValue(fmt.Sprintf("%d", tunnel.Group.Get().GetId()))
	}

	// Handle ipsec_profile reference from response
	if tunnel.HasIpsecProfile() && tunnel.IpsecProfile.IsSet() && tunnel.IpsecProfile.Get() != nil {
		data.IPSecProfile = types.StringValue(fmt.Sprintf("%d", tunnel.IpsecProfile.Get().GetId()))
	}

	// Handle tenant reference from response
	if tunnel.HasTenant() && tunnel.Tenant.IsSet() && tunnel.Tenant.Get() != nil {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tunnel.Tenant.Get().GetId()))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunnelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TunnelResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	tunnelID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tunnel ID",
			fmt.Sprintf("Unable to parse tunnel ID: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Deleting tunnel", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Delete the tunnel via API
	httpResp, err := r.client.VpnAPI.VpnTunnelsDestroy(ctx, tunnelID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting tunnel",
			utils.FormatAPIError("delete tunnel", err, httpResp),
		)
		return
	}
}

func (r *TunnelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
