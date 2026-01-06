// Package datasources contains Terraform data source implementations for the Netbox provider.
//

// This package integrates with the go-netbox OpenAPI client to provide
// read-only access to Netbox resources via Terraform data sources.

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
var _ datasource.DataSource = &TunnelDataSource{}

func NewTunnelDataSource() datasource.DataSource {
	return &TunnelDataSource{}
}

// TunnelDataSource defines the data source implementation.
type TunnelDataSource struct {
	client *netbox.APIClient
}

// TunnelDataSourceModel describes the data source data model.
type TunnelDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Status         types.String `tfsdk:"status"`
	Group          types.String `tfsdk:"group"`
	GroupID        types.String `tfsdk:"group_id"`
	Encapsulation  types.String `tfsdk:"encapsulation"`
	IPSecProfile   types.String `tfsdk:"ipsec_profile"`
	IPSecProfileID types.String `tfsdk:"ipsec_profile_id"`
	Tenant         types.String `tfsdk:"tenant"`
	TenantID       types.String `tfsdk:"tenant_id"`
	TunnelID       types.Int64  `tfsdk:"tunnel_id"`
	Description    types.String `tfsdk:"description"`
	Comments       types.String `tfsdk:"comments"`
	DisplayName    types.String `tfsdk:"display_name"`
	Tags           types.Set    `tfsdk:"tags"`
	CustomFields   types.Set    `tfsdk:"custom_fields"`
}

func (d *TunnelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunnel"
}

func (d *TunnelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a VPN tunnel in Netbox. Tunnels represent secure connections between network endpoints. You can identify the tunnel using `id` or `name`.",
		Attributes: map[string]schema.Attribute{
			"id":               nbschema.DSIDAttribute("tunnel"),
			"name":             nbschema.DSNameAttribute("tunnel"),
			"status":           nbschema.DSComputedStringAttribute("Operational status of the tunnel (planned, active, disabled)."),
			"group":            nbschema.DSComputedStringAttribute("Name of the tunnel group."),
			"group_id":         nbschema.DSComputedStringAttribute("ID of the tunnel group."),
			"encapsulation":    nbschema.DSComputedStringAttribute("Encapsulation protocol for the tunnel."),
			"ipsec_profile":    nbschema.DSComputedStringAttribute("Name of the IPSec profile."),
			"ipsec_profile_id": nbschema.DSComputedStringAttribute("ID of the IPSec profile."),
			"tenant":           nbschema.DSComputedStringAttribute("Name of the tenant."),
			"tenant_id":        nbschema.DSComputedStringAttribute("ID of the tenant."),
			"tunnel_id":        nbschema.DSComputedInt64Attribute("Tunnel identifier (numeric ID used by the tunnel protocol)."),
			"description":      nbschema.DSComputedStringAttribute("Detailed description of the tunnel."),
			"comments":         nbschema.DSComputedStringAttribute("Additional comments about the tunnel."),
			"display_name":     nbschema.DSComputedStringAttribute("Display name of the tunnel."),
			"tags":             nbschema.DSTagsAttribute(),
			"custom_fields":    nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *TunnelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *TunnelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TunnelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var tunnel *netbox.Tunnel
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID or name
	switch {
	case !data.ID.IsNull():
		// Search by ID
		tunnelID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading tunnel by ID", map[string]interface{}{
			"id": tunnelID,
		})

		// Parse the tunnel ID to int32 for the API call
		var tunnelIDInt int32
		if _, parseErr := fmt.Sscanf(tunnelID, "%d", &tunnelIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Tunnel ID",
				fmt.Sprintf("Tunnel ID must be a number, got: %s", tunnelID),
			)
			return
		}

		// Retrieve the tunnel via API
		tunnel, httpResp, err = d.client.VpnAPI.VpnTunnelsRetrieve(ctx, tunnelIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)

	case !data.Name.IsNull():
		// Search by name
		tunnelName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading tunnel by name", map[string]interface{}{
			"name": tunnelName,
		})

		// List tunnels with name filter
		var tunnels *netbox.PaginatedTunnelList
		tunnels, httpResp, err = d.client.VpnAPI.VpnTunnelsList(ctx).Name([]string{tunnelName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading tunnel",
				utils.FormatAPIError("read tunnel by name", err, httpResp),
			)
			return
		}
		if len(tunnels.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Tunnel Not Found",
				fmt.Sprintf("No tunnel found with name: %s", tunnelName),
			)
			return
		}
		if len(tunnels.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Tunnels Found",
				fmt.Sprintf("Multiple tunnels found with name: %s. Tunnel names may not be unique in Netbox.", tunnelName),
			)
			return
		}
		tunnel = &tunnels.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Tunnel Identifier",
			"Either 'id' or 'name' must be specified to identify the tunnel.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tunnel",
			utils.FormatAPIError("read tunnel", err, httpResp),
		)
		return
	}
	if httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Tunnel Not Found",
			"The specified tunnel was not found in Netbox.",
		)
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error reading tunnel",
			fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusOK, httpResp.StatusCode),
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
		data.Group = types.StringValue(tunnel.Group.Get().GetName())
		data.GroupID = types.StringValue(fmt.Sprintf("%d", tunnel.Group.Get().GetId()))
	} else {
		data.Group = types.StringNull()
		data.GroupID = types.StringNull()
	}

	// Handle ipsec_profile reference
	if tunnel.HasIpsecProfile() && tunnel.IpsecProfile.IsSet() && tunnel.IpsecProfile.Get() != nil {
		data.IPSecProfile = types.StringValue(tunnel.IpsecProfile.Get().GetName())
		data.IPSecProfileID = types.StringValue(fmt.Sprintf("%d", tunnel.IpsecProfile.Get().GetId()))
	} else {
		data.IPSecProfile = types.StringNull()
		data.IPSecProfileID = types.StringNull()
	}

	// Handle tenant reference
	if tunnel.HasTenant() && tunnel.Tenant.IsSet() && tunnel.Tenant.Get() != nil {
		data.Tenant = types.StringValue(tunnel.Tenant.Get().GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tunnel.Tenant.Get().GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
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

	// Handle display_name
	if tunnel.GetDisplay() != "" {
		data.DisplayName = types.StringValue(tunnel.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
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
		// For data sources, we extract all available custom fields
		customFields := utils.MapToCustomFieldModels(tunnel.GetCustomFields(), nil)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
