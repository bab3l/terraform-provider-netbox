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
	_ datasource.DataSource              = &IKEProposalDataSource{}
	_ datasource.DataSourceWithConfigure = &IKEProposalDataSource{}
)

// NewIKEProposalDataSource returns a new IKEProposal data source.
func NewIKEProposalDataSource() datasource.DataSource {
	return &IKEProposalDataSource{}
}

// IKEProposalDataSource defines the data source implementation.
type IKEProposalDataSource struct {
	client *netbox.APIClient
}

// IKEProposalDataSourceModel describes the data source data model.
type IKEProposalDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	AuthenticationMethod    types.String `tfsdk:"authentication_method"`
	EncryptionAlgorithm     types.String `tfsdk:"encryption_algorithm"`
	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`
	Group                   types.Int64  `tfsdk:"group"`
	SALifetime              types.Int64  `tfsdk:"sa_lifetime"`
	Comments                types.String `tfsdk:"comments"`
	DisplayName             types.String `tfsdk:"display_name"`
	Tags                    types.List   `tfsdk:"tags"`
	CustomFields            types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *IKEProposalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_proposal"
}

// Schema defines the schema for the data source.
func (d *IKEProposalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an IKE (Internet Key Exchange) Proposal in Netbox. IKE proposals define the security parameters for the IKE phase 1 negotiation in IPSec VPN connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the IKE proposal. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IKE proposal. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the IKE proposal.",
				Computed:            true,
			},
			"authentication_method": schema.StringAttribute{
				MarkdownDescription: "The authentication method for the IKE proposal. Values: `preshared-keys`, `certificates`, `rsa-signatures`, `dsa-signatures`.",
				Computed:            true,
			},
			"encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The encryption algorithm for the IKE proposal. Values: `aes-128-cbc`, `aes-128-gcm`, `aes-192-cbc`, `aes-192-gcm`, `aes-256-cbc`, `aes-256-gcm`, `3des-cbc`, `des-cbc`.",
				Computed:            true,
			},
			"authentication_algorithm": schema.StringAttribute{
				MarkdownDescription: "The authentication algorithm (hash) for the IKE proposal. Values: `hmac-sha1`, `hmac-sha256`, `hmac-sha384`, `hmac-sha512`, `hmac-md5`.",
				Computed:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The Diffie-Hellman group for the IKE proposal.",
				Computed:            true,
			},
			"sa_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Security association lifetime in seconds.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments about the IKE proposal.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the IKE proposal.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this IKE proposal.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *IKEProposalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *IKEProposalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IKEProposalDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var ike *netbox.IKEProposal

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
		tflog.Debug(ctx, "Reading IKEProposal by ID", map[string]interface{}{
			"id": id,
		})
		result, httpResp, err := d.client.VpnAPI.VpnIkeProposalsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading IKEProposal",
				utils.FormatAPIError(fmt.Sprintf("retrieve IKE proposal ID %d", id), err, httpResp),
			)
			return
		}
		ike = result

	case utils.IsSet(data.Name):
		// Looking up by name
		tflog.Debug(ctx, "Reading IKEProposal by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		listReq := d.client.VpnAPI.VpnIkeProposalsList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})
		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing IKEProposals",
				utils.FormatAPIError(fmt.Sprintf("list IKE proposals with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}
		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"IKEProposal not found",
				fmt.Sprintf("No IKE proposal found with name %q", data.Name.ValueString()),
			)
			return
		}
		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple IKEProposals found",
				fmt.Sprintf("Found %d IKE proposals with name %q. Please use id instead.", results.Count, data.Name.ValueString()),
			)
			return
		}
		ike = &results.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either id or name must be specified.",
		)
		return
	}

	// Map the result to state
	d.mapIKEProposalToState(ike, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapIKEProposalToState maps an IKEProposal API response to the Terraform state model.
func (d *IKEProposalDataSource) mapIKEProposalToState(ike *netbox.IKEProposal, data *IKEProposalDataSourceModel) {
	// ID
	data.ID = types.StringValue(fmt.Sprintf("%d", ike.Id))

	// Name
	data.Name = types.StringValue(ike.Name)

	// Description
	if ike.Description != nil && *ike.Description != "" {
		data.Description = types.StringValue(*ike.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Authentication Method
	if ike.AuthenticationMethod.Value != nil {
		data.AuthenticationMethod = types.StringValue(string(*ike.AuthenticationMethod.Value))
	}

	// Encryption Algorithm
	if ike.EncryptionAlgorithm.Value != nil {
		data.EncryptionAlgorithm = types.StringValue(string(*ike.EncryptionAlgorithm.Value))
	}

	// Authentication Algorithm
	if ike.AuthenticationAlgorithm != nil && ike.AuthenticationAlgorithm.Value != nil {
		data.AuthenticationAlgorithm = types.StringValue(string(*ike.AuthenticationAlgorithm.Value))
	} else {
		data.AuthenticationAlgorithm = types.StringNull()
	}

	// Group
	if ike.Group.Value != nil {
		data.Group = types.Int64Value(int64(*ike.Group.Value))
	}

	// SA Lifetime
	if ike.SaLifetime.IsSet() && ike.SaLifetime.Get() != nil {
		data.SALifetime = types.Int64Value(int64(*ike.SaLifetime.Get()))
	} else {
		data.SALifetime = types.Int64Null()
	}

	// Comments
	if ike.Comments != nil && *ike.Comments != "" {
		data.Comments = types.StringValue(*ike.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Display name
	if ike.GetDisplay() != "" {
		data.DisplayName = types.StringValue(ike.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Tags
	if len(ike.Tags) > 0 {
		tags := make([]string, len(ike.Tags))
		for i, tag := range ike.Tags {
			tags[i] = tag.Name
		}
		tagsValue, _ := types.ListValueFrom(context.Background(), types.StringType, tags)
		data.Tags = tagsValue
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Custom fields - datasources return ALL fields
	data.CustomFields = utils.CustomFieldsSetFromAPI(context.Background(), ike.HasCustomFields(), ike.GetCustomFields(), nil)
}
