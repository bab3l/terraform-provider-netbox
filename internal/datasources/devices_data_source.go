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
	_ datasource.DataSource              = &DevicesDataSource{}
	_ datasource.DataSourceWithConfigure = &DevicesDataSource{}
)

func NewDevicesDataSource() datasource.DataSource {
	return &DevicesDataSource{}
}

type DevicesDataSource struct {
	client *netbox.APIClient
}

type DevicesDataSourceModel struct {
	Filter  []utils.QueryFilterModel `tfsdk:"filter"`
	IDs     types.List               `tfsdk:"ids"`
	Names   types.List               `tfsdk:"names"`
	Devices types.List               `tfsdk:"devices"`
}

type deviceQueryResultModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func deviceQueryResultObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
	}
}

func (d *DevicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *DevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Query devices in NetBox using AWS-style filter blocks. Multiple `filter` blocks are ANDed; values within a filter are ORed.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				MarkdownDescription: "List of device IDs that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"names": schema.ListAttribute{
				MarkdownDescription: "Best-effort list of device names that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"devices": schema.ListNestedAttribute{
				MarkdownDescription: "List of matching devices as objects containing `id` and `name`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Device ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Device name.",
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
							MarkdownDescription: "Filter key name (e.g. `name`, `name__ic`, `serial`, `status`, `site`, `tag`, `q`, `custom_field`, `custom_field_value`).",
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

func (d *DevicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DevicesDataSourceModel

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
			"At least one `filter` block must be provided to avoid accidentally listing all devices.",
		)
		return
	}

	customFieldExists := filters["custom_field"]
	customFieldValueRaw := filters["custom_field_value"]
	delete(filters, "custom_field")
	delete(filters, "custom_field_value")

	customFieldValueFilters, valueDiags := utils.ParseCustomFieldValueFilters(customFieldValueRaw)
	resp.Diagnostics.Append(valueDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listReq := d.client.DcimAPI.DcimDevicesList(ctx)

	// Apply supported filters
	for name, values := range filters {
		switch name {
		case "name":
			listReq = listReq.Name(values)
		case "name__ic":
			listReq = listReq.NameIc(values)
		case "serial":
			listReq = listReq.Serial(values)
		case "status":
			listReq = listReq.Status(values)
		case "site":
			listReq = listReq.Site(values)
		case "tag":
			listReq = listReq.Tag(values)
		case "q":
			if len(values) != 1 {
				resp.Diagnostics.AddError(
					"Invalid filter values",
					"Filter `q` requires exactly one value.",
				)
				return
			}
			listReq = listReq.Q(values[0])
		default:
			resp.Diagnostics.AddError(
				"Unsupported filter",
				fmt.Sprintf("Unsupported filter name %q for netbox_devices.", name),
			)
			return
		}
	}

	tflog.Debug(ctx, "Querying devices", map[string]interface{}{
		"filters": filters,
	})

	// Paginate to ensure we return all matches
	const pageLimit int32 = 100
	var (
		offset  int32
		results []netbox.Device
	)

	for {
		pageReq := listReq.Limit(pageLimit).Offset(offset)
		page, httpResp, err := pageReq.Execute()
		utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying devices",
				utils.FormatAPIError("list devices", err, httpResp),
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
		for _, device := range results {
			customFields := device.GetCustomFields()
			if utils.MatchesCustomFieldFilters(customFields, customFieldExists, customFieldValueFilters) {
				filtered = append(filtered, device)
			}
		}
		results = filtered
	}

	ids := make([]string, 0, len(results))
	names := make([]string, 0, len(results))
	devices := make([]deviceQueryResultModel, 0, len(results))

	for _, device := range results {
		id := fmt.Sprintf("%d", device.GetId())
		ids = append(ids, id)

		if device.HasName() && device.Name.Get() != nil && *device.Name.Get() != "" {
			name := *device.Name.Get()
			names = append(names, name)
			devices = append(devices, deviceQueryResultModel{ID: types.StringValue(id), Name: types.StringValue(name)})
		} else {
			names = append(names, "")
			devices = append(devices, deviceQueryResultModel{ID: types.StringValue(id), Name: types.StringValue("")})
		}
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
	devicesValue, devicesDiags := types.ListValueFrom(ctx, deviceQueryResultObjectType(), devices)
	resp.Diagnostics.Append(devicesDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.IDs = idsValue
	data.Names = namesValue
	data.Devices = devicesValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
