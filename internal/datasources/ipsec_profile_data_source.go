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
	_ datasource.DataSource              = &IPSecProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &IPSecProfileDataSource{}
)

// NewIPSecProfileDataSource returns a new IPSecProfile data source.
func NewIPSecProfileDataSource() datasource.DataSource {
	return &IPSecProfileDataSource{}
}

// IPSecProfileDataSource defines the data source implementation.
type IPSecProfileDataSource struct {
	client *netbox.APIClient
}

// IPSecProfileDataSourceModel describes the data source data model.
type IPSecProfileDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Mode        types.String `tfsdk:"mode"`
	IKEPolicy   types.String `tfsdk:"ike_policy"`
	IPSecPolicy types.String `tfsdk:"ipsec_policy"`
	Comments    types.String `tfsdk:"comments"`
	Tags        types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *IPSecProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_profile"
}

// Schema defines the schema for the data source.
func (d *IPSecProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an IPSec Profile in Netbox. IPSec profiles combine IKE and IPSec policies to define complete VPN configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the IPSec profile. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec profile. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the IPSec profile.",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The IPSec mode. Values: `esp` (Encapsulating Security Payload), `ah` (Authentication Header).",
				Computed:            true,
			},
			"ike_policy": schema.StringAttribute{
				MarkdownDescription: "The name of the IKE policy used by this profile.",
				Computed:            true,
			},
			"ipsec_policy": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec policy used by this profile.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments about the IPSec profile.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this IPSec profile.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *IPSecProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *IPSecProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPSecProfileDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ipsec *netbox.IPSecProfile

	// Check if we're looking up by ID
	if utils.IsSet(data.ID) {
		id, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Reading IPSecProfile by ID", map[string]interface{}{
			"id": id,
		})

		result, httpResp, err := d.client.VpnAPI.VpnIpsecProfilesRetrieve(ctx, id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading IPSecProfile",
				utils.FormatAPIError(fmt.Sprintf("retrieve IPSec profile ID %d", id), err, httpResp),
			)
			return
		}
		ipsec = result
	} else if utils.IsSet(data.Name) {
		// Looking up by name
		tflog.Debug(ctx, "Reading IPSecProfile by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listReq := d.client.VpnAPI.VpnIpsecProfilesList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		results, httpResp, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing IPSecProfiles",
				utils.FormatAPIError(fmt.Sprintf("list IPSec profiles with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"IPSecProfile not found",
				fmt.Sprintf("No IPSec profile found with name %q", data.Name.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple IPSecProfiles found",
				fmt.Sprintf("Found %d IPSec profiles with name %q. Please use id instead.", results.Count, data.Name.ValueString()),
			)
			return
		}

		ipsec = &results.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either id or name must be specified.",
		)
		return
	}

	// Map the result to state
	d.mapIPSecProfileToState(ipsec, &data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapIPSecProfileToState maps an IPSecProfile API response to the Terraform state model.
func (d *IPSecProfileDataSource) mapIPSecProfileToState(ipsec *netbox.IPSecProfile, data *IPSecProfileDataSourceModel) {
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

	// Mode
	if ipsec.Mode.Value != nil {
		data.Mode = types.StringValue(string(*ipsec.Mode.Value))
	} else {
		data.Mode = types.StringNull()
	}

	// IKE Policy
	data.IKEPolicy = types.StringValue(ipsec.IkePolicy.Name)

	// IPSec Policy
	data.IPSecPolicy = types.StringValue(ipsec.IpsecPolicy.Name)

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
