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
	_ datasource.DataSource              = &VirtualMachinesDataSource{}
	_ datasource.DataSourceWithConfigure = &VirtualMachinesDataSource{}
)

func NewVirtualMachinesDataSource() datasource.DataSource {
	return &VirtualMachinesDataSource{}
}

type VirtualMachinesDataSource struct {
	client *netbox.APIClient
}

type VirtualMachinesDataSourceModel struct {
	Filter          []utils.QueryFilterModel `tfsdk:"filter"`
	IDs             types.List               `tfsdk:"ids"`
	Names           types.List               `tfsdk:"names"`
	VirtualMachines types.List               `tfsdk:"virtual_machines"`
}

type virtualMachineQueryResultModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func virtualMachineQueryResultObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
	}
}

func (d *VirtualMachinesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machines"
}

func (d *VirtualMachinesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Query virtual machines in NetBox using AWS-style filter blocks. Multiple `filter` blocks are ANDed; values within a filter are ORed.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				MarkdownDescription: "List of virtual machine IDs that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"names": schema.ListAttribute{
				MarkdownDescription: "Best-effort list of virtual machine names that match the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"virtual_machines": schema.ListNestedAttribute{
				MarkdownDescription: "List of matching virtual machines as objects containing `id` and `name`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Virtual machine ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Virtual machine name.",
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
							MarkdownDescription: "Filter key name (e.g. `name`, `name__ic`, `status`, `cluster`, `site`, `tag`, `q`, `custom_field`, `custom_field_value`).",
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

func (d *VirtualMachinesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VirtualMachinesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtualMachinesDataSourceModel

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
			"At least one `filter` block must be provided to avoid accidentally listing all virtual machines.",
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

	listReq := d.client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx)

	for name, values := range filters {
		switch name {
		case "name":
			listReq = listReq.Name(values)
		case "name__ic":
			listReq = listReq.NameIc(values)
		case "status":
			listReq = listReq.Status(values)
		case "cluster":
			listReq = listReq.Cluster(values)
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
				fmt.Sprintf("Unsupported filter name %q for netbox_virtual_machines.", name),
			)
			return
		}
	}

	tflog.Debug(ctx, "Querying virtual machines", map[string]interface{}{
		"filters": filters,
	})

	const pageLimit int32 = 100
	var (
		offset  int32
		results []netbox.VirtualMachineWithConfigContext
	)

	for {
		pageReq := listReq.Limit(pageLimit).Offset(offset)
		page, httpResp, err := pageReq.Execute()
		utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying virtual machines",
				utils.FormatAPIError("list virtual machines", err, httpResp),
			)
			return
		}

		pageResults := page.GetResults()
		if len(pageResults) == 0 {
			break
		}

		results = append(results, pageResults...)
		pageLen := len(pageResults)
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
		for _, vm := range results {
			customFields := vm.GetCustomFields()
			if utils.MatchesCustomFieldFilters(customFields, customFieldExists, customFieldValueFilters) {
				filtered = append(filtered, vm)
			}
		}
		results = filtered
	}

	ids := make([]string, 0, len(results))
	names := make([]string, 0, len(results))
	vms := make([]virtualMachineQueryResultModel, 0, len(results))

	for _, vm := range results {
		id := fmt.Sprintf("%d", vm.GetId())
		name := vm.GetName()
		ids = append(ids, id)
		names = append(names, name)
		vms = append(vms, virtualMachineQueryResultModel{ID: types.StringValue(id), Name: types.StringValue(name)})
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
	vmsValue, vmsDiags := types.ListValueFrom(ctx, virtualMachineQueryResultObjectType(), vms)
	resp.Diagnostics.Append(vmsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.IDs = idsValue
	data.Names = namesValue
	data.VirtualMachines = vmsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
