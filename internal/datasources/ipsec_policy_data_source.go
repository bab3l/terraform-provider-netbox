// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &IPSecPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &IPSecPolicyDataSource{}
)

// NewIPSecPolicyDataSource returns a new IPSecPolicy data source.
func NewIPSecPolicyDataSource() datasource.DataSource {
	return &IPSecPolicyDataSource{}
}

// IPSecPolicyDataSource defines the data source implementation.
type IPSecPolicyDataSource struct {
	client *netbox.APIClient
}

// IPSecPolicyDataSourceModel describes the data source data model.
type IPSecPolicyDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Proposals   types.List   `tfsdk:"proposals"`
	PFSGroup    types.Int64  `tfsdk:"pfs_group"`
	Comments    types.String `tfsdk:"comments"`
	Tags        types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *IPSecPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_policy"
}

// Schema defines the schema for the data source.
func (d *IPSecPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an IPSec Policy in Netbox. IPSec policies group together IPSec proposals and define the PFS (Perfect Forward Secrecy) group for VPN connections.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the IPSec policy. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec policy. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the IPSec policy.",
				Computed:            true,
			},
			"proposals": schema.ListAttribute{
				MarkdownDescription: "The list of IPSec proposal IDs associated with this policy.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"pfs_group": schema.Int64Attribute{
				MarkdownDescription: "The Diffie-Hellman group for Perfect Forward Secrecy.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments about the IPSec policy.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this IPSec policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *IPSecPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *IPSecPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPSecPolicyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ipsec *netbox.IPSecPolicy

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

		tflog.Debug(ctx, "Reading IPSecPolicy by ID", map[string]interface{}{
			"id": id,
		})

		result, httpResp, err := d.client.VpnAPI.VpnIpsecPoliciesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading IPSecPolicy",
				utils.FormatAPIError(fmt.Sprintf("retrieve IPSec policy ID %d", id), err, httpResp),
			)
			return
		}
		ipsec = result
	case utils.IsSet(data.Name):
		// Looking up by name
		tflog.Debug(ctx, "Reading IPSecPolicy by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listReq := d.client.VpnAPI.VpnIpsecPoliciesList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing IPSecPolicies",
				utils.FormatAPIError(fmt.Sprintf("list IPSec policies with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"IPSecPolicy not found",
				fmt.Sprintf("No IPSec policy found with name %q", data.Name.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple IPSecPolicies found",
				fmt.Sprintf("Found %d IPSec policies with name %q. Please use id instead.", results.Count, data.Name.ValueString()),
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
	d.mapIPSecPolicyToState(ipsec, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapIPSecPolicyToState maps an IPSecPolicy API response to the Terraform state model.
func (d *IPSecPolicyDataSource) mapIPSecPolicyToState(ipsec *netbox.IPSecPolicy, data *IPSecPolicyDataSourceModel) {
	// ID
	data.ID = types.StringValue(fmt.Sprintf("%d", ipsec.Id))

	// Name
	data.Name = types.StringValue(ipsec.Name)

	// Description
	if ipsec.Description != nil && *ipsec.Description != "" {
		data.Description = types.StringValue(*ipsec.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Proposals
	if len(ipsec.Proposals) > 0 {
		proposalIDs := make([]int64, len(ipsec.Proposals))
		for i, proposal := range ipsec.Proposals {
			proposalIDs[i] = int64(proposal.Id)
		}
		proposalsValue, _ := types.ListValueFrom(context.Background(), types.Int64Type, proposalIDs)
		data.Proposals = proposalsValue
	} else {
		data.Proposals = types.ListNull(types.Int64Type)
	}

	// PFS Group
	if ipsec.PfsGroup != nil && ipsec.PfsGroup.Value != nil {
		data.PFSGroup = types.Int64Value(int64(*ipsec.PfsGroup.Value))
	} else {
		data.PFSGroup = types.Int64Null()
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
}
