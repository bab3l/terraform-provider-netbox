package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &InterfacesDataSource{}
	_ datasource.DataSourceWithConfigure = &InterfacesDataSource{}
)

func NewInterfacesDataSource() datasource.DataSource {
	return &InterfacesDataSource{}
}

type InterfacesDataSource struct {
	client *netbox.APIClient
}

type InterfacesDataSourceModel struct {
	Filter     []utils.QueryFilterModel `tfsdk:"filter"`
	IDs        types.List               `tfsdk:"ids"`
	Names      types.List               `tfsdk:"names"`
	Interfaces types.List               `tfsdk:"interfaces"`
}

type interfaceQueryResultModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func interfaceQueryResultObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
	}
}

func (d *InterfacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interfaces"
}

func (d *InterfacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Query device interfaces in NetBox (DCIM) using AWS-style filter blocks. Multiple `filter` blocks are ANDed; values within a filter are ORed.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				MarkdownDescription: "List of interface IDs that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"names": schema.ListAttribute{
				MarkdownDescription: "List of interface names that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"interfaces": schema.ListNestedAttribute{
				MarkdownDescription: "List of matching interfaces as objects containing `id` and `name`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Interface ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Interface name.",
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
							MarkdownDescription: "Filter key name (e.g. `name`, `name__ic`, `device`, `device_id`, `site`, `site_id`, `type`, `enabled`, `tag`, `q`, `custom_field`, `custom_field_value`).",
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

func (d *InterfacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func parseInt32SliceFromStrings(values []string) ([]int32, error) {
	ids := make([]int32, 0, len(values))
	for _, s := range values {
		id, err := utils.ParseID(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseBoolFilterValue(values []string) (bool, error) {
	if len(values) != 1 {
		return false, fmt.Errorf("expected exactly one value")
	}
	s := strings.TrimSpace(strings.ToLower(values[0]))
	switch s {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean %q", values[0])
	}
}

func (d *InterfacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfacesDataSourceModel

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
			"At least one `filter` block must be provided to avoid accidentally listing all interfaces.",
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

	listReq := d.client.DcimAPI.DcimInterfacesList(ctx)

	for name, values := range filters {
		switch name {
		case filterKeyName:
			listReq = listReq.Name(values)
		case filterKeyNameIc:
			listReq = listReq.NameIc(values)
		case "device":
			listReq = listReq.Device(stringPointers(values))
		case "device_id":
			ids, err := parseInt32SliceFromStrings(values)
			if err != nil {
				resp.Diagnostics.AddError("Invalid filter values", fmt.Sprintf("Filter device_id must be numeric IDs: %s", err))
				return
			}
			listReq = listReq.DeviceId(ids)
		case filterKeySite:
			listReq = listReq.Site(values)
		case "site_id":
			ids, err := parseInt32SliceFromStrings(values)
			if err != nil {
				resp.Diagnostics.AddError("Invalid filter values", fmt.Sprintf("Filter site_id must be numeric IDs: %s", err))
				return
			}
			listReq = listReq.SiteId(ids)
		case "type":
			listReq = listReq.Type_(values)
		case "enabled":
			b, err := parseBoolFilterValue(values)
			if err != nil {
				resp.Diagnostics.AddError("Invalid filter values", "Filter `enabled` requires exactly one boolean value (true/false).")
				return
			}
			listReq = listReq.Enabled(b)
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
				fmt.Sprintf("Unsupported filter name %q for netbox_interfaces.", name),
			)
			return
		}
	}

	tflog.Debug(ctx, "Querying interfaces", map[string]interface{}{
		"filters": filters,
	})

	const pageLimit int32 = 100
	var (
		offset  int32
		results []netbox.Interface
	)

	for {
		pageReq := listReq.Limit(pageLimit).Offset(offset)
		page, httpResp, err := pageReq.Execute()
		utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying interfaces",
				utils.FormatAPIError("list interfaces", err, httpResp),
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
		for _, iface := range results {
			customFields := iface.GetCustomFields()
			if utils.MatchesCustomFieldFilters(customFields, customFieldExists, customFieldValueFilters) {
				filtered = append(filtered, iface)
			}
		}
		results = filtered
	}

	ids := make([]string, 0, len(results))
	names := make([]string, 0, len(results))
	items := make([]interfaceQueryResultModel, 0, len(results))

	for _, iface := range results {
		id := fmt.Sprintf("%d", iface.GetId())
		name := iface.GetName()
		ids = append(ids, id)
		names = append(names, name)
		items = append(items, interfaceQueryResultModel{ID: types.StringValue(id), Name: types.StringValue(name)})
	}

	idsValue, idDiags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(idDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	namesValue, nameDiags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(nameDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	itemsValue, itemsDiags := types.ListValueFrom(ctx, interfaceQueryResultObjectType(), items)
	resp.Diagnostics.Append(itemsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.IDs = idsValue
	data.Names = namesValue
	data.Interfaces = itemsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
