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
	_ datasource.DataSource              = &IPRangeDataSource{}
	_ datasource.DataSourceWithConfigure = &IPRangeDataSource{}
)

// NewIPRangeDataSource returns a new IP Range data source.
func NewIPRangeDataSource() datasource.DataSource {
	return &IPRangeDataSource{}
}

// IPRangeDataSource defines the data source implementation.
type IPRangeDataSource struct {
	client *netbox.APIClient
}

// IPRangeDataSourceModel describes the data source data model.
type IPRangeDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	DisplayName  types.String `tfsdk:"display_name"`
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
	Size         types.Int64  `tfsdk:"size"`
	VRF          types.String `tfsdk:"vrf"`
	VRFID        types.Int64  `tfsdk:"vrf_id"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantID     types.Int64  `tfsdk:"tenant_id"`
	Status       types.String `tfsdk:"status"`
	Role         types.String `tfsdk:"role"`
	RoleID       types.Int64  `tfsdk:"role_id"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	MarkUtilized types.Bool   `tfsdk:"mark_utilized"`
	Tags         types.List   `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *IPRangeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_range"
}

// Schema defines the schema for the data source.
func (d *IPRangeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an IP address range in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the IP range. Must be specified to look up a specific range.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the IP range.",
				Computed:            true,
			},
			"start_address": schema.StringAttribute{
				MarkdownDescription: "The starting IP address of the range. Can be used with `end_address` to look up a range.",
				Optional:            true,
				Computed:            true,
			},
			"end_address": schema.StringAttribute{
				MarkdownDescription: "The ending IP address of the range. Can be used with `start_address` to look up a range.",
				Optional:            true,
				Computed:            true,
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "The number of IP addresses in the range.",
				Computed:            true,
			},
			"vrf": schema.StringAttribute{
				MarkdownDescription: "The name of the VRF this IP range is assigned to.",
				Computed:            true,
			},
			"vrf_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VRF this IP range is assigned to.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant this IP range is assigned to.",
				Computed:            true,
			},
			"tenant_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the tenant this IP range is assigned to.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the IP range.",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The name of the IPAM role for this IP range.",
				Computed:            true,
			},
			"role_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the IPAM role for this IP range.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the IP range.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments for the IP range.",
				Computed:            true,
			},
			"mark_utilized": schema.BoolAttribute{
				MarkdownDescription: "Whether this range is treated as fully utilized.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this IP range.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *IPRangeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *IPRangeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPRangeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var ipRange *netbox.IPRange

	// Check if we're looking up by ID
	switch {
	case utils.IsSet(data.ID):
		var idInt int
		_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}
		tflog.Debug(ctx, "Reading IP range by ID", map[string]interface{}{
			"id": idInt,
		})
		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}
		result, httpResp, err := d.client.IpamAPI.IpamIpRangesRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading IP range",
				utils.FormatAPIError(fmt.Sprintf("retrieve IP range ID %d", idInt), err, httpResp),
			)
			return
		}
		ipRange = result

	case utils.IsSet(data.StartAddress) && utils.IsSet(data.EndAddress):
		// Looking up by start and end address

		tflog.Debug(ctx, "Reading IP range by addresses", map[string]interface{}{
			"start_address": data.StartAddress.ValueString(),
			"end_address":   data.EndAddress.ValueString(),
		})
		listReq := d.client.IpamAPI.IpamIpRangesList(ctx)
		listReq = listReq.StartAddress([]string{data.StartAddress.ValueString()})
		listReq = listReq.EndAddress([]string{data.EndAddress.ValueString()})
		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing IP ranges",
				utils.FormatAPIError(fmt.Sprintf("list IP ranges with start %q and end %q", data.StartAddress.ValueString(), data.EndAddress.ValueString()), err, httpResp),
			)
			return
		}
		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"IP range not found",
				fmt.Sprintf("No IP range found with start address %q and end address %q", data.StartAddress.ValueString(), data.EndAddress.ValueString()),
			)
			return
		}
		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple IP ranges found",
				fmt.Sprintf("Found %d IP ranges with start address %q and end address %q. Please use 'id' to specify the exact range.", results.Count, data.StartAddress.ValueString(), data.EndAddress.ValueString()),
			)
			return
		}
		ipRange = &results.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"Either 'id' or both 'start_address' and 'end_address' must be specified to look up an IP range.",
		)
		return
	}

	// Map response to model
	d.mapIPRangeToDataSourceModel(ctx, ipRange, &data)
	tflog.Debug(ctx, "Read IP range", map[string]interface{}{
		"id":            data.ID.ValueString(),
		"start_address": data.StartAddress.ValueString(),
		"end_address":   data.EndAddress.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapIPRangeToDataSourceModel maps a Netbox IPRange to the Terraform data source model.
func (d *IPRangeDataSource) mapIPRangeToDataSourceModel(ctx context.Context, ipRange *netbox.IPRange, data *IPRangeDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", ipRange.Id))

	// Display Name
	if ipRange.GetDisplay() != "" {
		data.DisplayName = types.StringValue(ipRange.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	data.StartAddress = types.StringValue(ipRange.StartAddress)
	data.EndAddress = types.StringValue(ipRange.EndAddress)
	data.Size = types.Int64Value(int64(ipRange.Size))

	// VRF
	if ipRange.Vrf.IsSet() && ipRange.Vrf.Get() != nil {
		vrfObj := ipRange.Vrf.Get()
		data.VRF = types.StringValue(vrfObj.Name)
		data.VRFID = types.Int64Value(int64(vrfObj.Id))
	} else {
		data.VRF = types.StringNull()
		data.VRFID = types.Int64Null()
	}

	// Tenant
	if ipRange.Tenant.IsSet() && ipRange.Tenant.Get() != nil {
		tenantObj := ipRange.Tenant.Get()
		data.Tenant = types.StringValue(tenantObj.Name)
		data.TenantID = types.Int64Value(int64(tenantObj.Id))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.Int64Null()
	}

	// Status
	if ipRange.Status != nil {
		data.Status = types.StringValue(string(ipRange.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Role
	if ipRange.Role.IsSet() && ipRange.Role.Get() != nil {
		roleObj := ipRange.Role.Get()
		data.Role = types.StringValue(roleObj.Name)
		data.RoleID = types.Int64Value(int64(roleObj.Id))
	} else {
		data.Role = types.StringNull()
		data.RoleID = types.Int64Null()
	}

	// Description
	if ipRange.Description != nil && *ipRange.Description != "" {
		data.Description = types.StringValue(*ipRange.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if ipRange.Comments != nil && *ipRange.Comments != "" {
		data.Comments = types.StringValue(*ipRange.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Mark Utilized
	if ipRange.MarkUtilized != nil {
		data.MarkUtilized = types.BoolValue(*ipRange.MarkUtilized)
	} else {
		data.MarkUtilized = types.BoolValue(false)
	}

	// Tags - convert to list of strings (tag names)
	if len(ipRange.Tags) > 0 {
		tagNames := make([]string, len(ipRange.Tags))
		for i, tag := range ipRange.Tags {
			tagNames[i] = tag.Name
		}
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Handle custom fields - datasources return ALL fields
	if ipRange.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(ipRange.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
