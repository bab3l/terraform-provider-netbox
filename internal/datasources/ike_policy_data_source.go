// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &IKEPolicyDataSource{}

	_ datasource.DataSourceWithConfigure = &IKEPolicyDataSource{}
)

// NewIKEPolicyDataSource returns a new IKEPolicy data source.

func NewIKEPolicyDataSource() datasource.DataSource {

	return &IKEPolicyDataSource{}

}

// IKEPolicyDataSource defines the data source implementation.

type IKEPolicyDataSource struct {
	client *netbox.APIClient
}

// IKEPolicyDataSourceModel describes the data source data model.

type IKEPolicyDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Description types.String `tfsdk:"description"`

	Version types.Int64 `tfsdk:"version"`

	Mode types.String `tfsdk:"mode"`

	Proposals types.List `tfsdk:"proposals"`

	Comments types.String `tfsdk:"comments"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.List `tfsdk:"tags"`
}

// Metadata returns the data source type name.

func (d *IKEPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_ike_policy"

}

// Schema defines the schema for the data source.

func (d *IKEPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about an IKE (Internet Key Exchange) Policy in Netbox. IKE policies group together IKE proposals and define the IKE version and mode for IPSec VPN connections.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The ID of the IKE policy. Either `id` or `name` must be specified.",

				Optional: true,

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the IKE policy. Either `id` or `name` must be specified.",

				Optional: true,

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "The description of the IKE policy.",

				Computed: true,
			},

			"version": schema.Int64Attribute{

				MarkdownDescription: "The IKE version. Values: `1` (IKEv1), `2` (IKEv2).",

				Computed: true,
			},

			"mode": schema.StringAttribute{

				MarkdownDescription: "The IKE negotiation mode. Values: `aggressive`, `main`. Only applicable for IKEv1.",

				Computed: true,
			},

			"proposals": schema.ListAttribute{

				MarkdownDescription: "The list of IKE proposal IDs associated with this policy.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Comments about the IKE policy.",

				Computed: true,
			},

			"display_name": schema.StringAttribute{

				MarkdownDescription: "The display name of the IKE policy.",

				Computed: true,
			},

			"tags": schema.ListAttribute{

				MarkdownDescription: "The tags assigned to this IKE policy.",

				Computed: true,

				ElementType: types.StringType,
			},
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *IKEPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *IKEPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data IKEPolicyDataSourceModel

	// Read Terraform configuration data into the model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var ike *netbox.IKEPolicy

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

		tflog.Debug(ctx, "Reading IKEPolicy by ID", map[string]interface{}{

			"id": id,
		})

		result, httpResp, err := d.client.VpnAPI.VpnIkePoliciesRetrieve(ctx, id).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading IKEPolicy",

				utils.FormatAPIError(fmt.Sprintf("retrieve IKE policy ID %d", id), err, httpResp),
			)

			return

		}

		ike = result

	case utils.IsSet(data.Name):

		// Looking up by name

		tflog.Debug(ctx, "Reading IKEPolicy by name", map[string]interface{}{

			"name": data.Name.ValueString(),
		})

		listReq := d.client.VpnAPI.VpnIkePoliciesList(ctx)

		listReq = listReq.Name([]string{data.Name.ValueString()})

		results, httpResp, err := listReq.Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error listing IKEPolicies",

				utils.FormatAPIError(fmt.Sprintf("list IKE policies with name %q", data.Name.ValueString()), err, httpResp),
			)

			return

		}

		if results.Count == 0 {

			resp.Diagnostics.AddError(

				"IKEPolicy not found",

				fmt.Sprintf("No IKE policy found with name %q", data.Name.ValueString()),
			)

			return

		}

		if results.Count > 1 {

			resp.Diagnostics.AddError(

				"Multiple IKEPolicies found",

				fmt.Sprintf("Found %d IKE policies with name %q. Please use id instead.", results.Count, data.Name.ValueString()),
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

	d.mapIKEPolicyToState(ike, &data)

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapIKEPolicyToState maps an IKEPolicy API response to the Terraform state model.

func (d *IKEPolicyDataSource) mapIKEPolicyToState(ike *netbox.IKEPolicy, data *IKEPolicyDataSourceModel) {

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

	// Version

	if ike.Version.Value != nil {

		data.Version = types.Int64Value(int64(*ike.Version.Value))

	} else {

		data.Version = types.Int64Null()

	}

	// Mode

	if ike.Mode != nil && ike.Mode.Value != nil && *ike.Mode.Value != "" {

		data.Mode = types.StringValue(string(*ike.Mode.Value))

	} else {

		data.Mode = types.StringNull()

	}

	// Proposals

	if len(ike.Proposals) > 0 {

		proposalIDs := make([]int64, len(ike.Proposals))

		for i, proposal := range ike.Proposals {

			proposalIDs[i] = int64(proposal.Id)

		}

		proposalsValue, _ := types.ListValueFrom(context.Background(), types.Int64Type, proposalIDs)

		data.Proposals = proposalsValue

	} else {

		data.Proposals = types.ListNull(types.Int64Type)

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

}
