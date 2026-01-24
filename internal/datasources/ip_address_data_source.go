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
	_ datasource.DataSource              = &IPAddressDataSource{}
	_ datasource.DataSourceWithConfigure = &IPAddressDataSource{}
)

// NewIPAddressDataSource returns a new IP Address data source.
func NewIPAddressDataSource() datasource.DataSource {
	return &IPAddressDataSource{}
}

// IPAddressDataSource defines the data source implementation.
type IPAddressDataSource struct {
	client *netbox.APIClient
}

// IPAddressDataSourceModel describes the data source data model.
type IPAddressDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Address            types.String `tfsdk:"address"`
	DisplayName        types.String `tfsdk:"display_name"`
	VRF                types.String `tfsdk:"vrf"`
	VRFID              types.Int64  `tfsdk:"vrf_id"`
	Tenant             types.String `tfsdk:"tenant"`
	TenantID           types.Int64  `tfsdk:"tenant_id"`
	Status             types.String `tfsdk:"status"`
	Role               types.String `tfsdk:"role"`
	AssignedObjectType types.String `tfsdk:"assigned_object_type"`
	AssignedObjectID   types.Int64  `tfsdk:"assigned_object_id"`
	NatInside          types.String `tfsdk:"nat_inside"`
	DNSName            types.String `tfsdk:"dns_name"`
	Description        types.String `tfsdk:"description"`
	Comments           types.String `tfsdk:"comments"`
	Tags               types.List   `tfsdk:"tags"`
	CustomFields       types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *IPAddressDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_address"
}

// Schema defines the schema for the data source.
func (d *IPAddressDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an IP address in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the IP address. Either `id` or `address` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The IP address with prefix length (e.g., 192.168.1.1/24).",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the IP address.",
				Computed:            true,
			},
			"vrf": schema.StringAttribute{
				MarkdownDescription: "The name of the VRF this IP address is assigned to.",
				Computed:            true,
			},
			"vrf_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VRF this IP address is assigned to.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant this IP address is assigned to.",
				Computed:            true,
			},
			"tenant_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the tenant this IP address is assigned to.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the IP address.",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the IP address.",
				Computed:            true,
			},
			"assigned_object_type": schema.StringAttribute{
				MarkdownDescription: "The content type of the assigned object.",
				Computed:            true,
			},
			"assigned_object_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the assigned object.",
				Computed:            true,
			},
			"nat_inside": schema.StringAttribute{
				MarkdownDescription: "The ID of the inside IP address for NAT (the IP for which this address is the outside IP).",
				Computed:            true,
			},
			"dns_name": schema.StringAttribute{
				MarkdownDescription: "Hostname or FQDN.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the IP address.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments for the IP address.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this IP address.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *IPAddressDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *IPAddressDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPAddressDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var ipAddress *netbox.IPAddress

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
		tflog.Debug(ctx, "Reading IP address by ID", map[string]interface{}{
			"id": idInt,
		})
		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}
		result, httpResp, err := d.client.IpamAPI.IpamIpAddressesRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading IP address",
				utils.FormatAPIError(fmt.Sprintf("retrieve IP address ID %d", idInt), err, httpResp),
			)
			return
		}
		ipAddress = result

	case utils.IsSet(data.Address):
		// Looking up by address
		tflog.Debug(ctx, "Reading IP address by address", map[string]interface{}{
			"address": data.Address.ValueString(),
		})
		listReq := d.client.IpamAPI.IpamIpAddressesList(ctx)
		listReq = listReq.Address([]string{data.Address.ValueString()})
		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing IP addresses",
				utils.FormatAPIError(fmt.Sprintf("list IP addresses with address %q", data.Address.ValueString()), err, httpResp),
			)
			return
		}
		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"IP address not found",
				fmt.Sprintf("No IP address found with address %q", data.Address.ValueString()),
			)
			return
		}
		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple IP addresses found",
				fmt.Sprintf("Found %d IP addresses with address %q. Please specify the ID to uniquely identify the IP address.", results.Count, data.Address.ValueString()),
			)
			return
		}
		ipAddress = &results.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Either 'id' or 'address' must be specified to look up an IP address.",
		)
		return
	}

	// Map the IP address to our model
	d.mapIPAddressToState(ctx, ipAddress, &data)
	tflog.Debug(ctx, "Read IP address", map[string]interface{}{
		"id":      data.ID.ValueString(),
		"address": data.Address.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapIPAddressToState maps a Netbox IPAddress to the Terraform state model.
func (d *IPAddressDataSource) mapIPAddressToState(ctx context.Context, ipAddress *netbox.IPAddress, data *IPAddressDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", ipAddress.Id))
	data.Address = types.StringValue(ipAddress.Address)

	// Display Name
	if ipAddress.GetDisplay() != "" {
		data.DisplayName = types.StringValue(ipAddress.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// VRF
	if ipAddress.Vrf.IsSet() && ipAddress.Vrf.Get() != nil {
		data.VRF = types.StringValue(ipAddress.Vrf.Get().Name)
		data.VRFID = types.Int64Value(int64(ipAddress.Vrf.Get().Id))
	} else {
		data.VRF = types.StringNull()
		data.VRFID = types.Int64Null()
	}

	// Tenant
	if ipAddress.Tenant.IsSet() && ipAddress.Tenant.Get() != nil {
		data.Tenant = types.StringValue(ipAddress.Tenant.Get().Name)
		data.TenantID = types.Int64Value(int64(ipAddress.Tenant.Get().Id))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.Int64Null()
	}

	// Status
	if ipAddress.Status != nil {
		data.Status = types.StringValue(string(ipAddress.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Role
	if ipAddress.Role != nil {
		data.Role = types.StringValue(string(ipAddress.Role.GetValue()))
	} else {
		data.Role = types.StringNull()
	}

	// Assigned Object Type
	if ipAddress.AssignedObjectType.IsSet() && ipAddress.AssignedObjectType.Get() != nil {
		data.AssignedObjectType = types.StringValue(*ipAddress.AssignedObjectType.Get())
	} else {
		data.AssignedObjectType = types.StringNull()
	}

	// Assigned Object ID
	if ipAddress.AssignedObjectId.IsSet() && ipAddress.AssignedObjectId.Get() != nil {
		data.AssignedObjectID = types.Int64Value(*ipAddress.AssignedObjectId.Get())
	} else {
		data.AssignedObjectID = types.Int64Null()
	}

	// NAT Inside
	if ipAddress.NatInside.IsSet() && ipAddress.NatInside.Get() != nil {
		nat := ipAddress.NatInside.Get()
		data.NatInside = types.StringValue(fmt.Sprintf("%d", nat.GetId()))
	} else {
		data.NatInside = types.StringNull()
	}

	// DNS Name
	if ipAddress.DnsName != nil && *ipAddress.DnsName != "" {
		data.DNSName = types.StringValue(*ipAddress.DnsName)
	} else {
		data.DNSName = types.StringNull()
	}

	// Description
	if ipAddress.Description != nil && *ipAddress.Description != "" {
		data.Description = types.StringValue(*ipAddress.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if ipAddress.Comments != nil && *ipAddress.Comments != "" {
		data.Comments = types.StringValue(*ipAddress.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags (slug list)
	if len(ipAddress.Tags) > 0 {
		tagSlugs := make([]string, 0, len(ipAddress.Tags))
		for _, tag := range ipAddress.Tags {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		tagList, _ := types.ListValueFrom(ctx, types.StringType, tagSlugs)
		data.Tags = tagList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Handle custom fields - datasources return ALL fields
	if ipAddress.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(ipAddress.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
