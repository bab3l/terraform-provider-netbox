// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &VirtualMachineDataSource{}

// NewVirtualMachineDataSource returns a new Virtual Machine data source.

func NewVirtualMachineDataSource() datasource.DataSource {
	return &VirtualMachineDataSource{}
}

// VirtualMachineDataSource defines the data source implementation.

type VirtualMachineDataSource struct {
	client *netbox.APIClient
}

// VirtualMachineDataSourceModel describes the data source data model.

type VirtualMachineDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Status types.String `tfsdk:"status"`

	Site types.String `tfsdk:"site"`

	Cluster types.String `tfsdk:"cluster"`

	Role types.String `tfsdk:"role"`

	Tenant types.String `tfsdk:"tenant"`

	Platform types.String `tfsdk:"platform"`

	Vcpus types.Float64 `tfsdk:"vcpus"`

	Memory types.Int64 `tfsdk:"memory"`

	Disk types.Int64 `tfsdk:"disk"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *VirtualMachineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

// Schema defines the schema for the data source.

func (d *VirtualMachineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a virtual machine in Netbox. Virtual machines represent virtualized compute instances. You can identify the virtual machine using `id` or `name`.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.DSIDAttribute("virtual machine"),

			"name": nbschema.DSNameAttribute("virtual machine"),

			"status": nbschema.DSComputedStringAttribute("The status of the virtual machine (offline, active, planned, staged, failed, decommissioning)."),

			"site": nbschema.DSComputedStringAttribute("The site where this virtual machine is located."),

			"cluster": nbschema.DSComputedStringAttribute("The cluster this virtual machine belongs to."),

			"role": nbschema.DSComputedStringAttribute("The device role for this virtual machine."),

			"tenant": nbschema.DSComputedStringAttribute("The tenant this virtual machine is assigned to."),

			"platform": nbschema.DSComputedStringAttribute("The platform (operating system) running on this virtual machine."),

			"vcpus": nbschema.DSComputedFloat64Attribute("The number of virtual CPUs allocated to this virtual machine."),

			"memory": nbschema.DSComputedInt64Attribute("The amount of memory (in MB) allocated to this virtual machine."),

			"disk": nbschema.DSComputedInt64Attribute("The total disk space (in GB) allocated to this virtual machine."),

			"description": nbschema.DSComputedStringAttribute("Detailed description of the virtual machine."),

			"comments": nbschema.DSComputedStringAttribute("Additional comments or notes about the virtual machine."),

			"display_name": nbschema.DSComputedStringAttribute("Display name of the virtual machine."),

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure sets up the data source with the provider client.

func (d *VirtualMachineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves data from the API.

func (d *VirtualMachineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtualMachineDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var vm *netbox.VirtualMachineWithConfigContext

	var err error

	var httpResp *http.Response

	// Determine if we're searching by ID or name

	switch {
	case !data.ID.IsNull():

		// Search by ID

		vmID := data.ID.ValueString()

		tflog.Debug(ctx, "Reading virtual machine by ID", map[string]interface{}{
			"id": vmID,
		})

		var vmIDInt int32

		if _, parseErr := fmt.Sscanf(vmID, "%d", &vmIDInt); parseErr != nil {
			resp.Diagnostics.AddError(

				"Invalid Virtual Machine ID",

				fmt.Sprintf("Virtual Machine ID must be a number, got: %s", vmID),
			)

			return
		}

		vm, httpResp, err = d.client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx, vmIDInt).Execute()

		defer utils.CloseResponseBody(httpResp)

	case !data.Name.IsNull():

		// Search by name

		vmName := data.Name.ValueString()

		tflog.Debug(ctx, "Reading virtual machine by name", map[string]interface{}{
			"name": vmName,
		})

		var vms *netbox.PaginatedVirtualMachineWithConfigContextList

		vms, httpResp, err = d.client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{vmName}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError(

				"Error reading virtual machine",

				utils.FormatAPIError("read virtual machine by name", err, httpResp),
			)

			return
		}

		if len(vms.GetResults()) == 0 {
			resp.Diagnostics.AddError(

				"Virtual Machine Not Found",

				fmt.Sprintf("No virtual machine found with name: %s", vmName),
			)

			return
		}

		if len(vms.GetResults()) > 1 {
			resp.Diagnostics.AddError(

				"Multiple Virtual Machines Found",

				fmt.Sprintf("Multiple virtual machines found with name: %s. Virtual machine names may not be unique in Netbox.", vmName),
			)

			return
		}

		vm = &vms.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Virtual Machine Identifier",

			"Either 'id' or 'name' must be specified to identify the virtual machine.",
		)

		return
	}

	if err != nil {
		resp.Diagnostics.AddError(

			"Error reading virtual machine",

			utils.FormatAPIError("read virtual machine", err, httpResp),
		)

		return
	}

	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(

			"Virtual Machine Not Found",

			fmt.Sprintf("No virtual machine found with ID: %s", data.ID.ValueString()),
		)

		return
	}

	// Map response to state

	data.ID = types.StringValue(fmt.Sprintf("%d", vm.GetId()))

	data.Name = types.StringValue(vm.GetName())

	// Status

	if vm.HasStatus() {
		data.Status = types.StringValue(string(vm.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Site

	if vm.Site.IsSet() && vm.Site.Get() != nil {
		data.Site = types.StringValue(vm.Site.Get().GetName())
	} else {
		data.Site = types.StringNull()
	}

	// Cluster

	if vm.Cluster.IsSet() && vm.Cluster.Get() != nil {
		data.Cluster = types.StringValue(vm.Cluster.Get().GetName())
	} else {
		data.Cluster = types.StringNull()
	}

	// Role

	if vm.Role.IsSet() && vm.Role.Get() != nil {
		data.Role = types.StringValue(vm.Role.Get().GetName())
	} else {
		data.Role = types.StringNull()
	}

	// Tenant

	if vm.Tenant.IsSet() && vm.Tenant.Get() != nil {
		data.Tenant = types.StringValue(vm.Tenant.Get().GetName())
	} else {
		data.Tenant = types.StringNull()
	}

	// Platform

	if vm.Platform.IsSet() && vm.Platform.Get() != nil {
		data.Platform = types.StringValue(vm.Platform.Get().GetName())
	} else {
		data.Platform = types.StringNull()
	}

	// Vcpus

	if vm.Vcpus.IsSet() && vm.Vcpus.Get() != nil {
		data.Vcpus = types.Float64Value(*vm.Vcpus.Get())
	} else {
		data.Vcpus = types.Float64Null()
	}

	// Memory

	if vm.Memory.IsSet() && vm.Memory.Get() != nil {
		data.Memory = types.Int64Value(int64(*vm.Memory.Get()))
	} else {
		data.Memory = types.Int64Null()
	}

	// Disk

	if vm.Disk.IsSet() && vm.Disk.Get() != nil {
		data.Disk = types.Int64Value(int64(*vm.Disk.Get()))
	} else {
		data.Disk = types.Int64Null()
	}

	// Description

	if vm.HasDescription() && vm.GetDescription() != "" {
		data.Description = types.StringValue(vm.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments

	if vm.HasComments() && vm.GetComments() != "" {
		data.Comments = types.StringValue(vm.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle display_name

	if vm.GetDisplay() != "" {
		data.DisplayName = types.StringValue(vm.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle tags

	if vm.HasTags() && len(vm.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(vm.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {
			return
		}

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields

	if vm.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(vm.GetCustomFields(), nil)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		resp.Diagnostics.Append(cfDiags...)

		if resp.Diagnostics.HasError() {
			return
		}

		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Debug(ctx, "Read virtual machine", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
