// Package datasources provides Terraform data source implementations for Netbox resources.
//

// Data sources allow Terraform configurations to read existing Netbox resources
// by their unique identifiers or names.

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
var _ datasource.DataSource = &TunnelTerminationDataSource{}

func NewTunnelTerminationDataSource() datasource.DataSource {
	return &TunnelTerminationDataSource{}
}

// TunnelTerminationDataSource defines the data source implementation.
type TunnelTerminationDataSource struct {
	client *netbox.APIClient
}

// TunnelTerminationDataSourceModel describes the data source data model.
type TunnelTerminationDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Tunnel          types.String `tfsdk:"tunnel"`
	TunnelName      types.String `tfsdk:"tunnel_name"`
	Role            types.String `tfsdk:"role"`
	TerminationType types.String `tfsdk:"termination_type"`
	TerminationID   types.Int64  `tfsdk:"termination_id"`
	OutsideIP       types.String `tfsdk:"outside_ip"`
	DisplayName     types.String `tfsdk:"display_name"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

func (d *TunnelTerminationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunnel_termination"
}

func (d *TunnelTerminationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a VPN tunnel termination in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id": nbschema.DSIDAttribute("tunnel termination"),
			"tunnel": schema.StringAttribute{
				MarkdownDescription: "ID of the tunnel. Use this to filter tunnel terminations by tunnel.",
				Optional:            true,
				Computed:            true,
			},
			"tunnel_name": schema.StringAttribute{
				MarkdownDescription: "Name of the tunnel. Use this to filter tunnel terminations by tunnel name.",
				Optional:            true,
				Computed:            true,
			},
			"role": nbschema.DSComputedStringAttribute("Role of the tunnel termination (peer, hub)."),
			"termination_type": schema.StringAttribute{
				MarkdownDescription: "Content type of the termination object.",
				Computed:            true,
			},
			"termination_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the termination object (device or virtual machine).",
				Computed:            true,
			},
			"outside_ip":    nbschema.DSComputedStringAttribute("ID of the outside IP address."),
			"display_name":  nbschema.DSComputedStringAttribute("Display name of the tunnel termination."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *TunnelTerminationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TunnelTerminationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TunnelTerminationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var tunnelTermination *netbox.TunnelTermination
	var err error

	// Look up by ID or tunnel reference
	switch {
	case !data.ID.IsNull() && data.ID.ValueString() != "":
		id, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError(
				"Error parsing ID",
				fmt.Sprintf("Could not parse tunnel termination ID %s: %s", data.ID.ValueString(), parseErr),
			)
			return
		}
		tflog.Debug(ctx, "Reading tunnel termination by ID", map[string]interface{}{
			"id": id,
		})
		var httpResp *http.Response
		tunnelTermination, httpResp, err = d.client.VpnAPI.VpnTunnelTerminationsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading tunnel termination",
				fmt.Sprintf("Unable to read tunnel termination ID %d: %s", id, err.Error()),
			)
			return
		}

	case !data.Tunnel.IsNull() || !data.TunnelName.IsNull():
		// Look up by tunnel reference
		var tunnelID int32
		if !data.Tunnel.IsNull() && data.Tunnel.ValueString() != "" {
			id, parseErr := utils.ParseID(data.Tunnel.ValueString())
			if parseErr != nil {
				resp.Diagnostics.AddError(
					"Error parsing Tunnel ID",
					fmt.Sprintf("Could not parse tunnel ID %s: %s", data.Tunnel.ValueString(), parseErr),
				)
				return
			}
			tunnelID = id
		} else if !data.TunnelName.IsNull() && data.TunnelName.ValueString() != "" {
			// Look up tunnel by name
			tunnelList, httpResp, lookupErr := d.client.VpnAPI.VpnTunnelsList(ctx).Name([]string{data.TunnelName.ValueString()}).Execute()
			defer utils.CloseResponseBody(httpResp)
			if lookupErr != nil {
				resp.Diagnostics.AddError(
					"Error looking up tunnel",
					utils.FormatAPIError(fmt.Sprintf("look up tunnel named '%s'", data.TunnelName.ValueString()), lookupErr, httpResp),
				)
				return
			}
			if len(tunnelList.Results) == 0 {
				resp.Diagnostics.AddError(
					"Tunnel not found",
					fmt.Sprintf("No tunnel found with name '%s'", data.TunnelName.ValueString()),
				)
				return
			}
			tunnelID = tunnelList.Results[0].GetId()
		}
		tflog.Debug(ctx, "Reading tunnel termination by tunnel ID", map[string]interface{}{
			"tunnel_id": tunnelID,
		})

		// List tunnel terminations for this tunnel
		list, httpResp, listErr := d.client.VpnAPI.VpnTunnelTerminationsList(ctx).TunnelId([]int32{tunnelID}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if listErr != nil {
			resp.Diagnostics.AddError(
				"Error listing tunnel terminations",
				utils.FormatAPIError(fmt.Sprintf("list tunnel terminations for tunnel ID %d", tunnelID), listErr, httpResp),
			)
			return
		}
		if len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"Tunnel termination not found",
				fmt.Sprintf("No tunnel termination found for tunnel ID %d", tunnelID),
			)
			return
		}

		// Return the first termination found
		tunnelTermination = &list.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Either 'id' or 'tunnel'/'tunnel_name' must be specified.",
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", tunnelTermination.GetId()))
	data.Tunnel = types.StringValue(fmt.Sprintf("%d", tunnelTermination.Tunnel.GetId()))
	data.TunnelName = types.StringValue(tunnelTermination.Tunnel.GetName())
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
		data.OutsideIP = types.StringValue(tunnelTermination.OutsideIp.Get().GetDisplay())
	} else {
		data.OutsideIP = types.StringNull()
	}

	// Handle display_name
	if tunnelTermination.GetDisplay() != "" {
		data.DisplayName = types.StringValue(tunnelTermination.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle tags
	if tunnelTermination.HasTags() && len(tunnelTermination.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(tunnelTermination.GetTags())
		tagsValue, tagsDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(tagsDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if tunnelTermination.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(tunnelTermination.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(cfDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
