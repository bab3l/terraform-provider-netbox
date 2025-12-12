// Package datasources provides Terraform data source implementations for NetBox objects.
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
var _ datasource.DataSource = &FrontPortTemplateDataSource{}

// NewFrontPortTemplateDataSource returns a new data source implementing the front port template data source.
func NewFrontPortTemplateDataSource() datasource.DataSource {
	return &FrontPortTemplateDataSource{}
}

// FrontPortTemplateDataSource defines the data source implementation.
type FrontPortTemplateDataSource struct {
	client *netbox.APIClient
}

// FrontPortTemplateDataSourceModel describes the data source data model.
type FrontPortTemplateDataSourceModel struct {
	ID               types.Int32  `tfsdk:"id"`
	DeviceType       types.Int32  `tfsdk:"device_type"`
	ModuleType       types.Int32  `tfsdk:"module_type"`
	Name             types.String `tfsdk:"name"`
	Label            types.String `tfsdk:"label"`
	Type             types.String `tfsdk:"type"`
	Color            types.String `tfsdk:"color"`
	RearPort         types.String `tfsdk:"rear_port"`
	RearPortID       types.Int32  `tfsdk:"rear_port_id"`
	RearPortPosition types.Int32  `tfsdk:"rear_port_position"`
	Description      types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *FrontPortTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_front_port_template"
}

// Schema defines the schema for the data source.
func (d *FrontPortTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a front port template in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the front port template.",
				Optional:            true,
				Computed:            true,
			},
			"device_type": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the device type. Used with name for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"module_type": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the module type. Used with name for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the front port template. Used with device_type or module_type for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the front port template.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of front port.",
				Computed:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the front port in hex format.",
				Computed:            true,
			},
			"rear_port": schema.StringAttribute{
				MarkdownDescription: "The name of the rear port template this front port maps to.",
				Computed:            true,
			},
			"rear_port_id": schema.Int32Attribute{
				MarkdownDescription: "The ID of the rear port template this front port maps to.",
				Computed:            true,
			},
			"rear_port_position": schema.Int32Attribute{
				MarkdownDescription: "Position on the rear port that this front port maps to.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the front port template.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *FrontPortTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves the data source data.
func (d *FrontPortTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FrontPortTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var template *netbox.FrontPortTemplate

	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Lookup by ID
		templateID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading front port template by ID", map[string]interface{}{
			"id": templateID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimFrontPortTemplatesRetrieve(ctx, templateID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading front port template",
				utils.FormatAPIError(fmt.Sprintf("read front port template ID %d", templateID), err, httpResp),
			)
			return
		}
		template = response
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by device_type/module_type and name
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading front port template by name", map[string]interface{}{
			"name": name,
		})

		listReq := d.client.DcimAPI.DcimFrontPortTemplatesList(ctx).Name([]string{name})

		if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
			deviceTypeID := data.DeviceType.ValueInt32()
			listReq = listReq.DeviceTypeId([]*int32{&deviceTypeID})
		}
		if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {
			moduleTypeID := data.ModuleType.ValueInt32()
			listReq = listReq.ModuleTypeId([]*int32{&moduleTypeID})
		}

		response, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading front port template",
				utils.FormatAPIError(fmt.Sprintf("read front port template by name %s", name), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Front Port Template Not Found",
				fmt.Sprintf("No front port template found with name: %s", name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Front Port Templates Found",
				fmt.Sprintf("Found %d front port templates with name %s. Please specify device_type or module_type to narrow results, or use ID.", count, name),
			)
			return
		}

		template = &response.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to lookup a front port template.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(template, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *FrontPortTemplateDataSource) mapResponseToModel(template *netbox.FrontPortTemplate, data *FrontPortTemplateDataSourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())

	// Map device type
	if template.DeviceType.IsSet() && template.DeviceType.Get() != nil {
		data.DeviceType = types.Int32Value(template.DeviceType.Get().Id)
	} else {
		data.DeviceType = types.Int32Null()
	}

	// Map module type
	if template.ModuleType.IsSet() && template.ModuleType.Get() != nil {
		data.ModuleType = types.Int32Value(template.ModuleType.Get().Id)
	} else {
		data.ModuleType = types.Int32Null()
	}

	// Map type
	data.Type = types.StringValue(string(template.Type.GetValue()))

	// Map label
	if label, ok := template.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map color
	if color, ok := template.GetColorOk(); ok && color != nil && *color != "" {
		data.Color = types.StringValue(*color)
	} else {
		data.Color = types.StringNull()
	}

	// Map rear port
	data.RearPort = types.StringValue(template.RearPort.GetName())
	data.RearPortID = types.Int32Value(template.RearPort.GetId())

	// Map rear port position
	if pos, ok := template.GetRearPortPositionOk(); ok && pos != nil {
		data.RearPortPosition = types.Int32Value(*pos)
	} else {
		data.RearPortPosition = types.Int32Null()
	}

	// Map description
	if desc, ok := template.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}
}
