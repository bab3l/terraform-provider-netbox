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
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
var _ resource.Resource = &TunnelTerminationResource{}
var _ resource.ResourceWithImportState = &TunnelTerminationResource{}

func NewTunnelTerminationResource() resource.Resource {
	return &TunnelTerminationResource{}
}

// TunnelTerminationResource defines the resource implementation.
type TunnelTerminationResource struct {
	client *netbox.APIClient
}

// TunnelTerminationResourceModel describes the resource data model.
type TunnelTerminationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Tunnel          types.String `tfsdk:"tunnel"`
	Role            types.String `tfsdk:"role"`
	TerminationType types.String `tfsdk:"termination_type"`
	TerminationID   types.Int64  `tfsdk:"termination_id"`
	OutsideIP       types.String `tfsdk:"outside_ip"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

func (r *TunnelTerminationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunnel_termination"
}

func (r *TunnelTerminationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a VPN tunnel termination in Netbox. Tunnel terminations represent the endpoints of a tunnel, typically devices or virtual machines.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("tunnel termination"),
			"tunnel": schema.StringAttribute{
				MarkdownDescription: "ID of the tunnel this termination belongs to.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role of this tunnel termination. Valid values: `peer`, `hub`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"peer",
						"hub",
					),
				},
			},
			"termination_type": schema.StringAttribute{
				MarkdownDescription: "Content type of the termination object. Valid values: `dcim.device`, `virtualization.virtualmachine`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"dcim.device",
						"virtualization.virtualmachine",
					),
				},
			},
			"termination_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the termination object (device or virtual machine).",
				Optional:            true,
			},
			"outside_ip": schema.StringAttribute{
				MarkdownDescription: "ID of the outside IP address for this tunnel termination.",
				Optional:            true,
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *TunnelTerminationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TunnelTerminationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TunnelTerminationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse tunnel ID
	tunnelID, err := utils.ParseID(data.Tunnel.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tunnel ID",
			fmt.Sprintf("Unable to parse tunnel ID: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Creating tunnel termination", map[string]interface{}{
		"tunnel_id":        tunnelID,
		"termination_type": data.TerminationType.ValueString(),
	})

	// Build the API request - need to use BriefTunnelRequest for tunnel
	briefTunnel := *netbox.NewBriefTunnelRequest("")
	// Use additional properties to set the tunnel ID
	briefTunnel.AdditionalProperties = map[string]interface{}{
		"id": int(tunnelID),
	}

	tunnelTerminationRequest := netbox.NewWritableTunnelTerminationRequest(
		briefTunnel,
		data.TerminationType.ValueString(),
	)

	// Actually, we need to use AdditionalProperties to pass the tunnel ID
	tunnelTerminationRequest.AdditionalProperties = make(map[string]interface{})
	tunnelTerminationRequest.AdditionalProperties["tunnel"] = int(tunnelID)

	// Set role if provided
	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		role := netbox.PatchedWritableTunnelTerminationRequestRole(data.Role.ValueString())
		tunnelTerminationRequest.Role = &role
	}

	// Set termination_id if provided
	if !data.TerminationID.IsNull() && !data.TerminationID.IsUnknown() {
		tunnelTerminationRequest.TerminationId = *netbox.NewNullableInt64(netbox.PtrInt64(data.TerminationID.ValueInt64()))
	}

	// Set outside_ip if provided
	if !data.OutsideIP.IsNull() && !data.OutsideIP.IsUnknown() {
		outsideIPID, err := utils.ParseID(data.OutsideIP.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Outside IP ID",
				fmt.Sprintf("Unable to parse outside IP ID: %s", err),
			)
			return
		}
		tunnelTerminationRequest.AdditionalProperties["outside_ip"] = int(outsideIPID)
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tagModels []utils.TagModel
		diags := data.Tags.ElementsAs(ctx, &tagModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelTerminationRequest.Tags = utils.TagsToNestedTagRequests(tagModels)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() {
		var customFieldModels []utils.CustomFieldModel
		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelTerminationRequest.CustomFields = utils.CustomFieldsToMap(customFieldModels)
	}

	// Create the tunnel termination via API
	tunnelTermination, httpResp, err := r.client.VpnAPI.VpnTunnelTerminationsCreate(ctx).WritableTunnelTerminationRequest(*tunnelTerminationRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tunnel termination",
			utils.FormatAPIError("create tunnel termination", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapTunnelTerminationToState(ctx, tunnelTermination, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created tunnel termination", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunnelTerminationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TunnelTerminationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse tunnel termination ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Reading tunnel termination", map[string]interface{}{
		"id": id,
	})

	tunnelTermination, httpResp, err := r.client.VpnAPI.VpnTunnelTerminationsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading tunnel termination",
			utils.FormatAPIError(fmt.Sprintf("read tunnel termination ID %d", id), err, httpResp),
		)
		return
	}

	r.mapTunnelTerminationToState(ctx, tunnelTermination, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunnelTerminationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TunnelTerminationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse tunnel termination ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	// Parse tunnel ID
	tunnelID, err := utils.ParseID(data.Tunnel.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tunnel ID",
			fmt.Sprintf("Unable to parse tunnel ID: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Updating tunnel termination", map[string]interface{}{
		"id":        id,
		"tunnel_id": tunnelID,
	})

	// Build the API request
	briefTunnel := *netbox.NewBriefTunnelRequest("")
	briefTunnel.AdditionalProperties = map[string]interface{}{
		"id": int(tunnelID),
	}

	tunnelTerminationRequest := netbox.NewWritableTunnelTerminationRequest(
		briefTunnel,
		data.TerminationType.ValueString(),
	)

	tunnelTerminationRequest.AdditionalProperties = make(map[string]interface{})
	tunnelTerminationRequest.AdditionalProperties["tunnel"] = int(tunnelID)

	// Set role if provided
	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		role := netbox.PatchedWritableTunnelTerminationRequestRole(data.Role.ValueString())
		tunnelTerminationRequest.Role = &role
	}

	// Set termination_id if provided (or null to clear)
	if !data.TerminationID.IsNull() && !data.TerminationID.IsUnknown() {
		tunnelTerminationRequest.TerminationId = *netbox.NewNullableInt64(netbox.PtrInt64(data.TerminationID.ValueInt64()))
	} else {
		tunnelTerminationRequest.TerminationId = *netbox.NewNullableInt64(nil)
	}

	// Set outside_ip if provided (or null to clear)
	if !data.OutsideIP.IsNull() && !data.OutsideIP.IsUnknown() {
		outsideIPID, err := utils.ParseID(data.OutsideIP.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Outside IP ID",
				fmt.Sprintf("Unable to parse outside IP ID: %s", err),
			)
			return
		}
		tunnelTerminationRequest.AdditionalProperties["outside_ip"] = int(outsideIPID)
	} else {
		tunnelTerminationRequest.AdditionalProperties["outside_ip"] = nil
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tagModels []utils.TagModel
		diags := data.Tags.ElementsAs(ctx, &tagModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelTerminationRequest.Tags = utils.TagsToNestedTagRequests(tagModels)
	} else {
		tunnelTerminationRequest.Tags = []netbox.NestedTagRequest{}
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() {
		var customFieldModels []utils.CustomFieldModel
		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tunnelTerminationRequest.CustomFields = utils.CustomFieldsToMap(customFieldModels)
	} else {
		tunnelTerminationRequest.CustomFields = map[string]interface{}{}
	}

	// Update the tunnel termination via API
	tunnelTermination, httpResp, err := r.client.VpnAPI.VpnTunnelTerminationsUpdate(ctx, id).WritableTunnelTerminationRequest(*tunnelTerminationRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tunnel termination",
			utils.FormatAPIError(fmt.Sprintf("update tunnel termination ID %d", id), err, httpResp),
		)
		return
	}

	r.mapTunnelTerminationToState(ctx, tunnelTermination, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunnelTerminationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TunnelTerminationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse tunnel termination ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Deleting tunnel termination", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.VpnAPI.VpnTunnelTerminationsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting tunnel termination",
			utils.FormatAPIError(fmt.Sprintf("delete tunnel termination ID %d", id), err, httpResp),
		)
		return
	}
}

func (r *TunnelTerminationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapTunnelTerminationToState maps a TunnelTermination from the API to the Terraform state model.
func (r *TunnelTerminationResource) mapTunnelTerminationToState(ctx context.Context, tunnelTermination *netbox.TunnelTermination, data *TunnelTerminationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tunnelTermination.GetId()))
	data.Tunnel = types.StringValue(fmt.Sprintf("%d", tunnelTermination.Tunnel.GetId()))
	data.TerminationType = types.StringValue(tunnelTermination.GetTerminationType())

	// Handle role - check if value is set before accessing
	if tunnelTermination.Role.HasValue() {
		data.Role = types.StringValue(string(tunnelTermination.Role.GetValue()))
	} else {
		data.Role = types.StringNull()
	}

	// Handle termination_id
	if tunnelTermination.HasTerminationId() && tunnelTermination.TerminationId.IsSet() {
		if val := tunnelTermination.TerminationId.Get(); val != nil {
			data.TerminationID = types.Int64Value(*val)
		} else {
			data.TerminationID = types.Int64Null()
		}
	} else {
		data.TerminationID = types.Int64Null()
	}

	// Handle outside_ip reference
	if tunnelTermination.HasOutsideIp() && tunnelTermination.OutsideIp.IsSet() && tunnelTermination.OutsideIp.Get() != nil {
		data.OutsideIP = types.StringValue(fmt.Sprintf("%d", tunnelTermination.OutsideIp.Get().GetId()))
	} else {
		data.OutsideIP = types.StringNull()
	}

	// Handle tags
	if tunnelTermination.HasTags() && len(tunnelTermination.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(tunnelTermination.GetTags())
		tagsValue, tagsDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagsDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if tunnelTermination.HasCustomFields() && len(tunnelTermination.GetCustomFields()) > 0 {
		var existingCustomFields []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			cfDiags := data.CustomFields.ElementsAs(ctx, &existingCustomFields, false)
			diags.Append(cfDiags...)
			if diags.HasError() {
				return
			}
		}
		customFields := utils.MapToCustomFieldModels(tunnelTermination.GetCustomFields(), existingCustomFields)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
