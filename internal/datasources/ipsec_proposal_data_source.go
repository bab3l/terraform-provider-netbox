// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &IPSecProposalDataSource{}
	_ datasource.DataSourceWithConfigure = &IPSecProposalDataSource{}
)

// NewIPSecProposalDataSource returns a new IPSecProposal data source.
func NewIPSecProposalDataSource() datasource.DataSource {
	return &IPSecProposalDataSource{}
}

// IPSecProposalDataSource defines the data source implementation.
type IPSecProposalDataSource struct {
	client *netbox.APIClient
}

// IPSecProposalDataSourceModel describes the data source data model.
type IPSecProposalDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	DisplayName             types.String `tfsdk:"display_name"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	EncryptionAlgorithm     types.String `tfsdk:"encryption_algorithm"`
	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`
	SALifetimeSeconds       types.Int64  `tfsdk:"sa_lifetime_seconds"`
	SALifetimeData          types.Int64  `tfsdk:"sa_lifetime_data"`
	Comments                types.String `tfsdk:"comments"`
	Tags                    types.List   `tfsdk:"tags"`
	CustomFields            types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *IPSecProposalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_proposal"
}

// Schema defines the schema for the data source.
func (d *IPSecProposalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an IPSec Proposal in Netbox. IPSec proposals define the security parameters for the IPSec phase 2 (ESP/AH) negotiation in VPN connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the IPSec proposal. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec proposal. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the IPSec proposal.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the IPSec proposal.",
				Computed:            true,
			},
			"encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The encryption algorithm for the IPSec proposal. Values: `aes-128-cbc`, `aes-128-gcm`, `aes-192-cbc`, `aes-192-gcm`, `aes-256-cbc`, `aes-256-gcm`, `3des-cbc`, `des-cbc`.",
				Computed:            true,
			},
			"authentication_algorithm": schema.StringAttribute{
				MarkdownDescription: "The authentication algorithm (hash) for the IPSec proposal. Values: `hmac-sha1`, `hmac-sha256`, `hmac-sha384`, `hmac-sha512`, `hmac-md5`.",
				Computed:            true,
			},
			"sa_lifetime_seconds": schema.Int64Attribute{
				MarkdownDescription: "Security association lifetime in seconds.",
				Computed:            true,
			},
			"sa_lifetime_data": schema.Int64Attribute{
				MarkdownDescription: "Security association lifetime in kilobytes.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments about the IPSec proposal.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this IPSec proposal.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *IPSecProposalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *IPSecProposalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPSecProposalDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var ipsec *netbox.IPSecProposal

	// Check if we're looking up by ID
	switch {
	case utils.IsSet(data.ID):
		id, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}
		tflog.Debug(ctx, "Reading IPSecProposal by ID", map[string]interface{}{
			"id": id,
		})
		result, httpResp, err := d.client.VpnAPI.VpnIpsecProposalsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading IPSecProposal",
				utils.FormatAPIError(fmt.Sprintf("retrieve IPSec proposal ID %d", id), err, httpResp),
			)
			return
		}
		ipsec = result

	case utils.IsSet(data.Name):
		// Looking up by name
		tflog.Debug(ctx, "Reading IPSecProposal by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		listReq := d.client.VpnAPI.VpnIpsecProposalsList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})
		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing IPSecProposals",
				utils.FormatAPIError(fmt.Sprintf("list IPSec proposals with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}
		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"IPSecProposal not found",
				fmt.Sprintf("No IPSec proposal found with name %q", data.Name.ValueString()),
			)
			return
		}
		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple IPSecProposals found",
				fmt.Sprintf("Found %d IPSec proposals with name %q. Please use id instead.", results.Count, data.Name.ValueString()),
			)
			return
		}
		ipsec = &results.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either id or name must be specified.",
		)
		return
	}

	// Map the result to state
	d.mapIPSecProposalToState(ipsec, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapIPSecProposalToState maps an IPSecProposal API response to the Terraform state model.
func (d *IPSecProposalDataSource) mapIPSecProposalToState(ipsec *netbox.IPSecProposal, data *IPSecProposalDataSourceModel) {
	// ID
	data.ID = types.StringValue(fmt.Sprintf("%d", ipsec.Id))

	// Display Name
	if ipsec.GetDisplay() != "" {
		data.DisplayName = types.StringValue(ipsec.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Name
	data.Name = types.StringValue(ipsec.Name)

	// Description
	if ipsec.Description != nil && *ipsec.Description != "" {
		data.Description = types.StringValue(*ipsec.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Encryption Algorithm
	if ipsec.EncryptionAlgorithm != nil && ipsec.EncryptionAlgorithm.Value != nil {
		data.EncryptionAlgorithm = types.StringValue(string(*ipsec.EncryptionAlgorithm.Value))
	} else {
		data.EncryptionAlgorithm = types.StringNull()
	}

	// Authentication Algorithm
	if ipsec.AuthenticationAlgorithm != nil && ipsec.AuthenticationAlgorithm.Value != nil {
		data.AuthenticationAlgorithm = types.StringValue(string(*ipsec.AuthenticationAlgorithm.Value))
	} else {
		data.AuthenticationAlgorithm = types.StringNull()
	}

	// SA Lifetime Seconds
	if ipsec.SaLifetimeSeconds.IsSet() && ipsec.SaLifetimeSeconds.Get() != nil {
		data.SALifetimeSeconds = types.Int64Value(int64(*ipsec.SaLifetimeSeconds.Get()))
	} else {
		data.SALifetimeSeconds = types.Int64Null()
	}

	// SA Lifetime Data
	if ipsec.SaLifetimeData.IsSet() && ipsec.SaLifetimeData.Get() != nil {
		data.SALifetimeData = types.Int64Value(int64(*ipsec.SaLifetimeData.Get()))
	} else {
		data.SALifetimeData = types.Int64Null()
	}

	// Comments
	if ipsec.Comments != nil && *ipsec.Comments != "" {
		data.Comments = types.StringValue(*ipsec.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags
	if len(ipsec.Tags) > 0 {
		tags := make([]string, len(ipsec.Tags))

		for i, tag := range ipsec.Tags {
			tags[i] = tag.Name
		}
		tagsValue, _ := types.ListValueFrom(context.Background(), types.StringType, tags)
		data.Tags = tagsValue
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Custom fields - datasources return ALL fields
	data.CustomFields = utils.CustomFieldsSetFromAPI(context.Background(), ipsec.HasCustomFields(), ipsec.GetCustomFields(), nil)
}
