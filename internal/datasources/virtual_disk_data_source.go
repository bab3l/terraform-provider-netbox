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
	_ datasource.DataSource              = &VirtualDiskDataSource{}
	_ datasource.DataSourceWithConfigure = &VirtualDiskDataSource{}
)

// NewVirtualDiskDataSource returns a new VirtualDisk data source.
func NewVirtualDiskDataSource() datasource.DataSource {
	return &VirtualDiskDataSource{}
}

// VirtualDiskDataSource defines the data source implementation.
type VirtualDiskDataSource struct {
	client *netbox.APIClient
}

// VirtualDiskDataSourceModel describes the data source data model.
type VirtualDiskDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	VirtualMachine     types.String `tfsdk:"virtual_machine"`
	VirtualMachineName types.String `tfsdk:"virtual_machine_name"`
	Size               types.String `tfsdk:"size"`
	Description        types.String `tfsdk:"description"`
	Tags               types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *VirtualDiskDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_disk"
}

// Schema defines the schema for the data source.
func (d *VirtualDiskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a virtual disk in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the virtual disk. Either `id` or `name` (with `virtual_machine`) must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the virtual disk.",
				Optional:            true,
				Computed:            true,
			},
			"virtual_machine": schema.StringAttribute{
				MarkdownDescription: "The ID or name of the virtual machine. Required when looking up by name.",
				Optional:            true,
				Computed:            true,
			},
			"virtual_machine_name": schema.StringAttribute{
				MarkdownDescription: "The name of the virtual machine this disk belongs to.",
				Computed:            true,
			},
			"size": schema.StringAttribute{
				MarkdownDescription: "The size of the virtual disk in GB.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the virtual disk.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this virtual disk.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *VirtualDiskDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *VirtualDiskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtualDiskDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var vd *netbox.VirtualDisk

	// Check if we're looking up by ID or name
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

		tflog.Debug(ctx, "Reading VirtualDisk by ID", map[string]interface{}{
			"id": idInt,
		})

		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}

		result, httpResp, err := d.client.VirtualizationAPI.VirtualizationVirtualDisksRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VirtualDisk",
				utils.FormatAPIError(fmt.Sprintf("retrieve VirtualDisk ID %d", idInt), err, httpResp),
			)
			return
		}
		vd = result
	case utils.IsSet(data.Name):
		// Looking up by name - requires virtual_machine
		if !utils.IsSet(data.VirtualMachine) {
			resp.Diagnostics.AddError(
				"Missing virtual_machine",
				"When looking up a virtual disk by name, the 'virtual_machine' attribute must be specified.",
			)
			return
		}

		tflog.Debug(ctx, "Reading VirtualDisk by name", map[string]interface{}{
			"name":            data.Name.ValueString(),
			"virtual_machine": data.VirtualMachine.ValueString(),
		})

		listReq := d.client.VirtualizationAPI.VirtualizationVirtualDisksList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		// Try to parse virtual_machine as ID first
		var vmID int
		if _, err := fmt.Sscanf(data.VirtualMachine.ValueString(), "%d", &vmID); err == nil {
			vmID32, err := utils.SafeInt32(int64(vmID))
			if err != nil {
				resp.Diagnostics.AddError("Invalid Virtual Machine ID", fmt.Sprintf("Virtual Machine ID value overflow: %s", err))
				return
			}
			listReq = listReq.VirtualMachineId([]int32{vmID32})
		} else {
			// Otherwise, use as name
			listReq = listReq.VirtualMachine([]string{data.VirtualMachine.ValueString()})
		}

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing VirtualDisks",
				utils.FormatAPIError(fmt.Sprintf("list VirtualDisks with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"VirtualDisk not found",
				fmt.Sprintf("No VirtualDisk found with name %q for virtual machine %q", data.Name.ValueString(), data.VirtualMachine.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple VirtualDisks found",
				fmt.Sprintf("Found %d VirtualDisks with name %q. Please use 'id' to specify the exact VirtualDisk.", results.Count, data.Name.ValueString()),
			)
			return
		}

		vd = &results.Results[0]
	default:
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"Either 'id' or 'name' (with 'virtual_machine') must be specified to look up a VirtualDisk.",
		)
		return
	}

	// Map response to model
	d.mapVirtualDiskToDataSourceModel(ctx, vd, &data)

	tflog.Debug(ctx, "Read VirtualDisk", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapVirtualDiskToDataSourceModel maps a Netbox VirtualDisk to the Terraform data source model.
func (d *VirtualDiskDataSource) mapVirtualDiskToDataSourceModel(ctx context.Context, vd *netbox.VirtualDisk, data *VirtualDiskDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vd.Id))
	data.Name = types.StringValue(vd.Name)
	data.VirtualMachine = types.StringValue(fmt.Sprintf("%d", vd.VirtualMachine.GetId()))
	data.VirtualMachineName = types.StringValue(vd.VirtualMachine.GetName())
	data.Size = types.StringValue(fmt.Sprintf("%d", vd.Size))

	// Description
	if vd.Description != nil && *vd.Description != "" {
		data.Description = types.StringValue(*vd.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags - convert to list of strings (tag names)
	if len(vd.Tags) > 0 {
		tagNames := make([]string, len(vd.Tags))
		for i, tag := range vd.Tags {
			tagNames[i] = tag.Name
		}
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}
}
