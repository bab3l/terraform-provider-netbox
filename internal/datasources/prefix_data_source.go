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
	_ datasource.DataSource              = &PrefixDataSource{}
	_ datasource.DataSourceWithConfigure = &PrefixDataSource{}
)

// NewPrefixDataSource returns a new Prefix data source.
func NewPrefixDataSource() datasource.DataSource {
	return &PrefixDataSource{}
}

// PrefixDataSource defines the data source implementation.
type PrefixDataSource struct {
	client *netbox.APIClient
}

// PrefixDataSourceModel describes the data source data model.
type PrefixDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Prefix       types.String `tfsdk:"prefix"`
	Site         types.String `tfsdk:"site"`
	SiteID       types.Int64  `tfsdk:"site_id"`
	VRF          types.String `tfsdk:"vrf"`
	VRFID        types.Int64  `tfsdk:"vrf_id"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantID     types.Int64  `tfsdk:"tenant_id"`
	VLAN         types.String `tfsdk:"vlan"`
	VLANID       types.Int64  `tfsdk:"vlan_id"`
	Status       types.String `tfsdk:"status"`
	Role         types.String `tfsdk:"role"`
	RoleID       types.Int64  `tfsdk:"role_id"`
	IsPool       types.Bool   `tfsdk:"is_pool"`
	MarkUtilized types.Bool   `tfsdk:"mark_utilized"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *PrefixDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prefix"
}

// Schema defines the schema for the data source.
func (d *PrefixDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a prefix in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the prefix. Either `id` or `prefix` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "The IP prefix in CIDR notation (e.g., 192.168.1.0/24).",
				Optional:            true,
				Computed:            true,
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "The name of the site this prefix is assigned to.",
				Computed:            true,
			},
			"site_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the site this prefix is assigned to.",
				Computed:            true,
			},
			"vrf": schema.StringAttribute{
				MarkdownDescription: "The name of the VRF this prefix is assigned to.",
				Computed:            true,
			},
			"vrf_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VRF this prefix is assigned to.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant this prefix is assigned to.",
				Computed:            true,
			},
			"tenant_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the tenant this prefix is assigned to.",
				Computed:            true,
			},
			"vlan": schema.StringAttribute{
				MarkdownDescription: "The display name of the VLAN this prefix is assigned to.",
				Computed:            true,
			},
			"vlan_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VLAN this prefix is assigned to.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the prefix (container, active, reserved, deprecated).",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The name of the role for this prefix.",
				Computed:            true,
			},
			"role_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the role for this prefix.",
				Computed:            true,
			},
			"is_pool": schema.BoolAttribute{
				MarkdownDescription: "If true, all IP addresses within this prefix are considered usable.",
				Computed:            true,
			},
			"mark_utilized": schema.BoolAttribute{
				MarkdownDescription: "If true, the prefix is treated as fully utilized.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the prefix.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments for the prefix.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this prefix.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *PrefixDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PrefixDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrefixDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var prefix *netbox.Prefix

	// Check if we're looking up by ID
	if utils.IsSet(data.ID) {
		var idInt int
		_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Reading prefix by ID", map[string]interface{}{
			"id": idInt,
		})

		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}

		result, httpResp, err := d.client.IpamAPI.IpamPrefixesRetrieve(ctx, id32).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading prefix",
				utils.FormatAPIError(fmt.Sprintf("retrieve prefix ID %d", idInt), err, httpResp),
			)
			return
		}
		prefix = result
	} else if utils.IsSet(data.Prefix) {
		// Looking up by prefix CIDR
		tflog.Debug(ctx, "Reading prefix by CIDR", map[string]interface{}{
			"prefix": data.Prefix.ValueString(),
		})

		listReq := d.client.IpamAPI.IpamPrefixesList(ctx)
		listReq = listReq.Prefix([]string{data.Prefix.ValueString()})

		results, httpResp, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing prefixes",
				utils.FormatAPIError(fmt.Sprintf("list prefixes with CIDR %q", data.Prefix.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"Prefix not found",
				fmt.Sprintf("No prefix found with CIDR %q", data.Prefix.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple prefixes found",
				fmt.Sprintf("Found %d prefixes with CIDR %q. Please specify the ID to uniquely identify the prefix.", results.Count, data.Prefix.ValueString()),
			)
			return
		}

		prefix = &results.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Either 'id' or 'prefix' must be specified to look up a prefix.",
		)
		return
	}

	// Map the prefix to our model
	d.mapPrefixToState(ctx, prefix, &data)

	tflog.Debug(ctx, "Read prefix", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"prefix": data.Prefix.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapPrefixToState maps a Netbox Prefix to the Terraform state model.
func (d *PrefixDataSource) mapPrefixToState(ctx context.Context, prefix *netbox.Prefix, data *PrefixDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", prefix.Id))
	data.Prefix = types.StringValue(prefix.Prefix)

	// Site
	if prefix.Site.IsSet() && prefix.Site.Get() != nil {
		data.Site = types.StringValue(prefix.Site.Get().Name)
		data.SiteID = types.Int64Value(int64(prefix.Site.Get().Id))
	} else {
		data.Site = types.StringNull()
		data.SiteID = types.Int64Null()
	}

	// VRF
	if prefix.Vrf.IsSet() && prefix.Vrf.Get() != nil {
		data.VRF = types.StringValue(prefix.Vrf.Get().Name)
		data.VRFID = types.Int64Value(int64(prefix.Vrf.Get().Id))
	} else {
		data.VRF = types.StringNull()
		data.VRFID = types.Int64Null()
	}

	// Tenant
	if prefix.Tenant.IsSet() && prefix.Tenant.Get() != nil {
		data.Tenant = types.StringValue(prefix.Tenant.Get().Name)
		data.TenantID = types.Int64Value(int64(prefix.Tenant.Get().Id))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.Int64Null()
	}

	// VLAN
	if prefix.Vlan.IsSet() && prefix.Vlan.Get() != nil {
		data.VLAN = types.StringValue(prefix.Vlan.Get().Display)
		data.VLANID = types.Int64Value(int64(prefix.Vlan.Get().Id))
	} else {
		data.VLAN = types.StringNull()
		data.VLANID = types.Int64Null()
	}

	// Status
	if prefix.Status != nil {
		data.Status = types.StringValue(string(prefix.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Role
	if prefix.Role.IsSet() && prefix.Role.Get() != nil {
		data.Role = types.StringValue(prefix.Role.Get().Name)
		data.RoleID = types.Int64Value(int64(prefix.Role.Get().Id))
	} else {
		data.Role = types.StringNull()
		data.RoleID = types.Int64Null()
	}

	// IsPool
	if prefix.IsPool != nil {
		data.IsPool = types.BoolValue(*prefix.IsPool)
	} else {
		data.IsPool = types.BoolNull()
	}

	// MarkUtilized
	if prefix.MarkUtilized != nil {
		data.MarkUtilized = types.BoolValue(*prefix.MarkUtilized)
	} else {
		data.MarkUtilized = types.BoolNull()
	}

	// Description
	if prefix.Description != nil && *prefix.Description != "" {
		data.Description = types.StringValue(*prefix.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if prefix.Comments != nil && *prefix.Comments != "" {
		data.Comments = types.StringValue(*prefix.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags
	if len(prefix.Tags) > 0 {
		tagNames := make([]string, 0, len(prefix.Tags))
		for _, tag := range prefix.Tags {
			tagNames = append(tagNames, tag.Name)
		}
		tagList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}
}
