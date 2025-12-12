// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DeviceBayTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceBayTemplateDataSource{}
)

// NewDeviceBayTemplateDataSource returns a new DeviceBayTemplate data source.
func NewDeviceBayTemplateDataSource() datasource.DataSource {
	return &DeviceBayTemplateDataSource{}
}

// DeviceBayTemplateDataSource defines the data source implementation.
type DeviceBayTemplateDataSource struct {
	client *netbox.APIClient
}

// DeviceBayTemplateDataSourceModel describes the data source data model.
type DeviceBayTemplateDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	DeviceType     types.String `tfsdk:"device_type"`
	DeviceTypeName types.String `tfsdk:"device_type_name"`
	Name           types.String `tfsdk:"name"`
	Label          types.String `tfsdk:"label"`
	Description    types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *DeviceBayTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_bay_template"
}

// Schema defines the schema for the data source.
func (d *DeviceBayTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Device Bay Template in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the device bay template. Either `id` or `name` (with `device_type`) must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: "The ID or slug of the device type this template belongs to. Required when looking up by name.",
				Optional:            true,
				Computed:            true,
			},
			"device_type_name": schema.StringAttribute{
				MarkdownDescription: "The model name of the device type.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the device bay template.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label for the device bay.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the device bay template.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DeviceBayTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DeviceBayTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeviceBayTemplateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var template *netbox.DeviceBayTemplate

	if utils.IsSet(data.ID) {
		// Looking up by ID
		id, err := strconv.Atoi(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing DeviceBayTemplate ID",
				fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
			)
			return
		}

		tflog.Debug(ctx, "Reading DeviceBayTemplate by ID", map[string]interface{}{
			"id": id,
		})

		id32, err := utils.SafeInt32(int64(id))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}
		t, httpResp, err := d.client.DcimAPI.DcimDeviceBayTemplatesRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading DeviceBayTemplate",
				utils.FormatAPIError(fmt.Sprintf("read device bay template %d", id), err, httpResp),
			)
			return
		}
		template = t
	} else if utils.IsSet(data.Name) {
		// Looking up by name (requires device_type)
		tflog.Debug(ctx, "Reading DeviceBayTemplate by name", map[string]interface{}{
			"name":        data.Name.ValueString(),
			"device_type": data.DeviceType.ValueString(),
		})

		listReq := d.client.DcimAPI.DcimDeviceBayTemplatesList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		// If device_type is specified, filter by it
		if utils.IsSet(data.DeviceType) {
			// Try to parse as ID first
			if deviceTypeID, err := strconv.Atoi(data.DeviceType.ValueString()); err == nil {
				deviceTypeID32, convErr := utils.SafeInt32(int64(deviceTypeID))
				if convErr != nil {
					resp.Diagnostics.AddError("Invalid device_type ID", fmt.Sprintf("device_type ID overflow: %s", convErr))
					return
				}
				listReq = listReq.DeviceTypeId([]int32{deviceTypeID32})
			}
		}

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing DeviceBayTemplates",
				utils.FormatAPIError(fmt.Sprintf("list device bay templates with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"DeviceBayTemplate not found",
				fmt.Sprintf("No device bay template found with name %q", data.Name.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple DeviceBayTemplates found",
				fmt.Sprintf("Found %d device bay templates with name %q. Please specify 'device_type' to narrow results or use 'id'.", results.Count, data.Name.ValueString()),
			)
			return
		}

		template = &results.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"Either 'id' or 'name' must be specified to look up a device bay template.",
		)
		return
	}

	// Map response to model
	d.mapTemplateToDataSourceModel(template, &data)

	tflog.Debug(ctx, "Read DeviceBayTemplate", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapTemplateToDataSourceModel maps a Netbox DeviceBayTemplate to the Terraform data source model.
func (d *DeviceBayTemplateDataSource) mapTemplateToDataSourceModel(template *netbox.DeviceBayTemplate, data *DeviceBayTemplateDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", template.Id))
	data.DeviceType = types.StringValue(fmt.Sprintf("%d", template.DeviceType.GetId()))
	data.DeviceTypeName = types.StringValue(template.DeviceType.GetModel())
	data.Name = types.StringValue(template.Name)

	// Label
	if template.Label != nil && *template.Label != "" {
		data.Label = types.StringValue(*template.Label)
	} else {
		data.Label = types.StringNull()
	}

	// Description
	if template.Description != nil && *template.Description != "" {
		data.Description = types.StringValue(*template.Description)
	} else {
		data.Description = types.StringNull()
	}
}
