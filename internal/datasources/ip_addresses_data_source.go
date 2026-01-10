package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &IPAddressesDataSource{}
	_ datasource.DataSourceWithConfigure = &IPAddressesDataSource{}
)

func NewIPAddressesDataSource() datasource.DataSource {
	return &IPAddressesDataSource{}
}

type IPAddressesDataSource struct {
	client *netbox.APIClient
}

type IPAddressesDataSourceModel struct {
	Filter      []utils.QueryFilterModel `tfsdk:"filter"`
	IDs         types.List               `tfsdk:"ids"`
	Addresses   types.List               `tfsdk:"addresses"`
	IPAddresses types.List               `tfsdk:"ip_addresses"`
}

type ipAddressQueryResultModel struct {
	ID      types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
}

func ipAddressQueryResultObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":      types.StringType,
			"address": types.StringType,
		},
	}
}

func (d *IPAddressesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_addresses"
}

func (d *IPAddressesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Query IP addresses in NetBox using AWS-style filter blocks. Multiple `filter` blocks are ANDed; values within a filter are ORed.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				MarkdownDescription: "List of IP address IDs that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"addresses": schema.ListAttribute{
				MarkdownDescription: "List of IP address strings (with prefix length) that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ip_addresses": schema.ListNestedAttribute{
				MarkdownDescription: "List of matching IP addresses as objects containing `id` and `address`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "IP address ID.",
							Computed:            true,
						},
						"address": schema.StringAttribute{
							MarkdownDescription: "IP address string with prefix length (e.g. 192.0.2.1/24).",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SetNestedBlock{
				MarkdownDescription: "Filter criteria. At least one filter must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Filter key name (e.g. `address`, `status`, `role`, `tenant`, `tenant_id`, `vrf`, `vrf_id`, `dns_name`, `description`, `tag`, `q`, `custom_field`, `custom_field_value`).",
							Required:            true,
						},
						"values": schema.ListAttribute{
							MarkdownDescription: "List of values for this filter.",
							Required:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *IPAddressesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func stringPointers(values []string) []*string {
	ptrs := make([]*string, 0, len(values))
	for _, s := range values {
		v := s
		ptrs = append(ptrs, &v)
	}
	return ptrs
}

func int32PointersFromStrings(values []string) ([]*int32, error) {
	ptrs := make([]*int32, 0, len(values))
	for _, s := range values {
		id, err := utils.ParseID(s)
		if err != nil {
			return nil, err
		}
		v := id
		ptrs = append(ptrs, &v)
	}
	return ptrs, nil
}

func (d *IPAddressesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPAddressesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters, filterDiags := utils.ExpandQueryFilters(ctx, data.Filter)
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(filters) == 0 {
		resp.Diagnostics.AddError(
			"Missing filters",
			"At least one `filter` block must be provided to avoid accidentally listing all IP addresses.",
		)
		return
	}

	customFieldExists := filters[filterKeyCustomField]
	customFieldValueRaw := filters[filterKeyCustomFieldValue]
	delete(filters, filterKeyCustomField)
	delete(filters, filterKeyCustomFieldValue)

	customFieldValueFilters, valueDiags := utils.ParseCustomFieldValueFilters(customFieldValueRaw)
	resp.Diagnostics.Append(valueDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listReq := d.client.IpamAPI.IpamIpAddressesList(ctx)

	for name, values := range filters {
		switch name {
		case "address":
			listReq = listReq.Address(values)
		case filterKeyStatus:
			listReq = listReq.Status(values)
		case "role":
			listReq = listReq.Role(values)
		case "tenant":
			listReq = listReq.Tenant(values)
		case "tenant_id":
			ptrs, err := int32PointersFromStrings(values)
			if err != nil {
				resp.Diagnostics.AddError("Invalid filter values", fmt.Sprintf("Filter tenant_id must be numeric IDs: %s", err))
				return
			}
			listReq = listReq.TenantId(ptrs)
		case "vrf":
			listReq = listReq.Vrf(stringPointers(values))
		case "vrf_id":
			ptrs, err := int32PointersFromStrings(values)
			if err != nil {
				resp.Diagnostics.AddError("Invalid filter values", fmt.Sprintf("Filter vrf_id must be numeric IDs: %s", err))
				return
			}
			listReq = listReq.VrfId(ptrs)
		case "dns_name":
			listReq = listReq.DnsName(values)
		case "description":
			listReq = listReq.Description(values)
		case filterKeyTag:
			listReq = listReq.Tag(values)
		case filterKeyQ:
			if len(values) != 1 {
				resp.Diagnostics.AddError("Invalid filter values", "Filter `q` requires exactly one value.")
				return
			}
			listReq = listReq.Q(values[0])
		default:
			resp.Diagnostics.AddError(
				"Unsupported filter",
				fmt.Sprintf("Unsupported filter name %q for netbox_ip_addresses.", name),
			)
			return
		}
	}

	tflog.Debug(ctx, "Querying IP addresses", map[string]interface{}{
		"filters": filters,
	})

	const pageLimit int32 = 100
	var (
		offset  int32
		results []netbox.IPAddress
	)

	for {
		pageReq := listReq.Limit(pageLimit).Offset(offset)
		page, httpResp, err := pageReq.Execute()
		utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying IP addresses",
				utils.FormatAPIError("list IP addresses", err, httpResp),
			)
			return
		}

		if len(page.Results) == 0 {
			break
		}

		results = append(results, page.Results...)
		pageLen := len(page.Results)
		if pageLen < 0 || pageLen > int(pageLimit) {
			resp.Diagnostics.AddError(
				"Unexpected API response",
				fmt.Sprintf("Expected page size to be between 0 and %d, got %d", pageLimit, pageLen),
			)
			return
		}
		offset += int32(pageLen)

		if offset >= page.GetCount() {
			break
		}
	}

	if len(customFieldExists) > 0 || len(customFieldValueFilters) > 0 {
		filtered := results[:0]
		for _, ip := range results {
			customFields := ip.GetCustomFields()
			if utils.MatchesCustomFieldFilters(customFields, customFieldExists, customFieldValueFilters) {
				filtered = append(filtered, ip)
			}
		}
		results = filtered
	}

	ids := make([]string, 0, len(results))
	addrs := make([]string, 0, len(results))
	items := make([]ipAddressQueryResultModel, 0, len(results))

	for _, ip := range results {
		id := fmt.Sprintf("%d", ip.GetId())
		addr := ip.Address
		ids = append(ids, id)
		addrs = append(addrs, addr)
		items = append(items, ipAddressQueryResultModel{ID: types.StringValue(id), Address: types.StringValue(addr)})
	}

	idsValue, idDiags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(idDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	addrsValue, addrDiags := types.ListValueFrom(ctx, types.StringType, addrs)
	resp.Diagnostics.Append(addrDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	itemsValue, itemsDiags := types.ListValueFrom(ctx, ipAddressQueryResultObjectType(), items)
	resp.Diagnostics.Append(itemsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.IDs = idsValue
	data.Addresses = addrsValue
	data.IPAddresses = itemsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
