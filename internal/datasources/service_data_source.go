// Package datasources contains Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &ServiceDataSource{}
)

// NewServiceDataSource returns a new data source implementing the Service data source.
func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

// ServiceDataSource defines the data source implementation.
type ServiceDataSource struct {
	client *netbox.APIClient
}

// ServiceDataSourceModel describes the data source data model.
type ServiceDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Device         types.String `tfsdk:"device"`
	VirtualMachine types.String `tfsdk:"virtual_machine"`
	Name           types.String `tfsdk:"name"`
	Protocol       types.String `tfsdk:"protocol"`
	Ports          types.List   `tfsdk:"ports"`
	IPAddresses    types.List   `tfsdk:"ipaddresses"`
	Description    types.String `tfsdk:"description"`
	Comments       types.String `tfsdk:"comments"`
	Tags           types.Set    `tfsdk:"tags"`
	CustomFields   types.Set    `tfsdk:"custom_fields"`
	DisplayName    types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.
func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

// Schema defines the schema for the data source.
func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a network service in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the service. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "The device this service runs on (ID). Use with name for lookup.",
				Optional:            true,
				Computed:            true,
			},
			"virtual_machine": schema.StringAttribute{
				MarkdownDescription: "The virtual machine this service runs on (ID). Use with name for lookup.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service. Use with device or virtual_machine for lookup.",
				Optional:            true,
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol used by the service.",
				Computed:            true,
			},
			"ports": schema.ListAttribute{
				MarkdownDescription: "List of port numbers the service listens on.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"ipaddresses": schema.ListAttribute{
				MarkdownDescription: "List of IP address IDs associated with this service.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the service.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments.",
				Computed:            true,
			},
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the service."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the data source data.
func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var svc *netbox.Service

	// Look up by ID or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		svcID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Service ID",
				fmt.Sprintf("Service ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}
		tflog.Debug(ctx, "Reading service by ID", map[string]interface{}{
			"id": svcID,
		})
		result, httpResp, err := d.client.IpamAPI.IpamServicesRetrieve(ctx, svcID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading service",
				utils.FormatAPIError(fmt.Sprintf("read service ID %d", svcID), err, httpResp),
			)
			return
		}
		svc = result

	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Look up by name
		tflog.Debug(ctx, "Reading service by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		listReq := d.client.IpamAPI.IpamServicesList(ctx).Name([]string{data.Name.ValueString()})

		// Filter by device if provided
		if !data.Device.IsNull() && !data.Device.IsUnknown() {
			deviceID, err := utils.ParseID(data.Device.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid Device ID",
					fmt.Sprintf("Device ID must be a number, got: %s", data.Device.ValueString()),
				)
				return
			}
			listReq = listReq.DeviceId([]*int32{&deviceID})
		}

		// Filter by virtual_machine if provided
		if !data.VirtualMachine.IsNull() && !data.VirtualMachine.IsUnknown() {
			vmID, err := utils.ParseID(data.VirtualMachine.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid Virtual Machine ID",
					fmt.Sprintf("Virtual machine ID must be a number, got: %s", data.VirtualMachine.ValueString()),
				)
				return
			}
			listReq = listReq.VirtualMachineId([]*int32{&vmID})
		}
		listResp, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading service",
				utils.FormatAPIError(fmt.Sprintf("read service by name %s", data.Name.ValueString()), err, httpResp),
			)
			return
		}
		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Service not found",
				fmt.Sprintf("No service found with name: %s", data.Name.ValueString()),
			)
			return
		}
		if listResp.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple services found",
				fmt.Sprintf("Found %d services with name: %s. Please specify device or virtual_machine to narrow results.", listResp.GetCount(), data.Name.ValueString()),
			)
			return
		}
		svc = &listResp.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to look up a service.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, svc, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ServiceDataSource) mapResponseToModel(ctx context.Context, svc *netbox.Service, data *ServiceDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", svc.GetId()))
	data.Name = types.StringValue(svc.GetName())

	// Map device
	if svc.Device.IsSet() && svc.Device.Get() != nil {
		data.Device = types.StringValue(fmt.Sprintf("%d", svc.Device.Get().GetId()))
	} else {
		data.Device = types.StringNull()
	}

	// Map virtual_machine
	if svc.VirtualMachine.IsSet() && svc.VirtualMachine.Get() != nil {
		data.VirtualMachine = types.StringValue(fmt.Sprintf("%d", svc.VirtualMachine.Get().GetId()))
	} else {
		data.VirtualMachine = types.StringNull()
	}

	// Map protocol
	if protocol, ok := svc.GetProtocolOk(); ok && protocol != nil {
		data.Protocol = types.StringValue(string(protocol.GetValue()))
	} else {
		data.Protocol = types.StringNull()
	}

	// Map ports
	ports := svc.GetPorts()
	portsInt64 := make([]int64, len(ports))
	for i, p := range ports {
		portsInt64[i] = int64(p)
	}
	portsValue, portsDiags := types.ListValueFrom(ctx, types.Int64Type, portsInt64)
	diags.Append(portsDiags...)
	if diags.HasError() {
		return
	}
	data.Ports = portsValue

	// Map IP addresses
	if svc.HasIpaddresses() && len(svc.GetIpaddresses()) > 0 {
		ipAddrs := svc.GetIpaddresses()
		ipIDs := make([]int64, len(ipAddrs))
		for i, ip := range ipAddrs {
			ipIDs[i] = int64(ip.GetId())
		}
		ipValue, ipDiags := types.ListValueFrom(ctx, types.Int64Type, ipIDs)
		diags.Append(ipDiags...)
		if diags.HasError() {
			return
		}
		data.IPAddresses = ipValue
	} else {
		data.IPAddresses = types.ListNull(types.Int64Type)
	}

	// Map description
	if desc, ok := svc.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := svc.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Map display_name
	if svc.GetDisplay() != "" {
		data.DisplayName = types.StringValue(svc.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle tags
	if svc.HasTags() && len(svc.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(svc.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields - datasources return ALL fields
	if svc.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(svc.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
