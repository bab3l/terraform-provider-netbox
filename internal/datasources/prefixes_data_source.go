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
	_ datasource.DataSource              = &PrefixesDataSource{}
	_ datasource.DataSourceWithConfigure = &PrefixesDataSource{}
)

func NewPrefixesDataSource() datasource.DataSource {
	return &PrefixesDataSource{}
}

type PrefixesDataSource struct {
	client *netbox.APIClient
}

type PrefixesDataSourceModel struct {
	Filter   []utils.QueryFilterModel `tfsdk:"filter"`
	IDs      types.List               `tfsdk:"ids"`
	CIDRs    types.List               `tfsdk:"cidrs"`
	Prefixes types.List               `tfsdk:"prefixes"`
}

type prefixQueryResultModel struct {
	ID     types.String `tfsdk:"id"`
	Prefix types.String `tfsdk:"prefix"`
}

func prefixQueryResultObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":     types.StringType,
			"prefix": types.StringType,
		},
	}
}

func (d *PrefixesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prefixes"
}

func (d *PrefixesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Query prefixes in NetBox using AWS-style filter blocks. Multiple `filter` blocks are ANDed; values within a filter are ORed.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				MarkdownDescription: "List of prefix IDs that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cidrs": schema.ListAttribute{
				MarkdownDescription: "List of prefixes in CIDR notation (e.g. 192.0.2.0/24) that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"prefixes": schema.ListNestedAttribute{
				MarkdownDescription: "List of matching prefixes as objects containing `id` and `prefix`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Prefix ID.",
							Computed:            true,
						},
						"prefix": schema.StringAttribute{
							MarkdownDescription: "Prefix in CIDR notation (e.g. 192.0.2.0/24).",
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
							MarkdownDescription: "Filter key name (e.g. `prefix`, `status`, `role`, `tenant`, `tenant_id`, `vrf`, `vrf_id`, `site`, `site_id`, `description`, `tag`, `within`, `contains`, `q`, `custom_field`, `custom_field_value`).",
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

func (d *PrefixesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PrefixesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrefixesDataSourceModel

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
			"At least one `filter` block must be provided to avoid accidentally listing all prefixes.",
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

	listReq := d.client.IpamAPI.IpamPrefixesList(ctx)

	for name, values := range filters {
		switch name {
		case "prefix":
			listReq = listReq.Prefix(values)
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
		case filterKeySite:
			listReq = listReq.Site(values)
		case "site_id":
			ptrs, err := int32PointersFromStrings(values)
			if err != nil {
				resp.Diagnostics.AddError("Invalid filter values", fmt.Sprintf("Filter site_id must be numeric IDs: %s", err))
				return
			}
			listReq = listReq.SiteId(ptrs)
		case "description":
			listReq = listReq.Description(values)
		case filterKeyTag:
			listReq = listReq.Tag(values)
		case "within":
			if len(values) != 1 {
				resp.Diagnostics.AddError("Invalid filter values", "Filter `within` requires exactly one value.")
				return
			}
			listReq = listReq.Within(values[0])
		case "contains":
			if len(values) != 1 {
				resp.Diagnostics.AddError("Invalid filter values", "Filter `contains` requires exactly one value.")
				return
			}
			listReq = listReq.Contains(values[0])
		case filterKeyQ:
			if len(values) != 1 {
				resp.Diagnostics.AddError("Invalid filter values", "Filter `q` requires exactly one value.")
				return
			}
			listReq = listReq.Q(values[0])
		default:
			resp.Diagnostics.AddError(
				"Unsupported filter",
				fmt.Sprintf("Unsupported filter name %q for netbox_prefixes.", name),
			)
			return
		}
	}

	tflog.Debug(ctx, "Querying prefixes", map[string]interface{}{
		"filters": filters,
	})

	const pageLimit int32 = 100
	var (
		offset  int32
		results []netbox.Prefix
	)

	for {
		pageReq := listReq.Limit(pageLimit).Offset(offset)
		page, httpResp, err := pageReq.Execute()
		utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying prefixes",
				utils.FormatAPIError("list prefixes", err, httpResp),
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
		for _, p := range results {
			customFields := p.GetCustomFields()
			if utils.MatchesCustomFieldFilters(customFields, customFieldExists, customFieldValueFilters) {
				filtered = append(filtered, p)
			}
		}
		results = filtered
	}

	ids := make([]string, 0, len(results))
	cidrs := make([]string, 0, len(results))
	items := make([]prefixQueryResultModel, 0, len(results))

	for _, p := range results {
		id := fmt.Sprintf("%d", p.GetId())
		cidr := p.GetPrefix()
		ids = append(ids, id)
		cidrs = append(cidrs, cidr)
		items = append(items, prefixQueryResultModel{ID: types.StringValue(id), Prefix: types.StringValue(cidr)})
	}

	idsValue, idDiags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(idDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cidrsValue, cidrDiags := types.ListValueFrom(ctx, types.StringType, cidrs)
	resp.Diagnostics.Append(cidrDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	itemsValue, itemsDiags := types.ListValueFrom(ctx, prefixQueryResultObjectType(), items)
	resp.Diagnostics.Append(itemsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.IDs = idsValue
	data.CIDRs = cidrsValue
	data.Prefixes = itemsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
