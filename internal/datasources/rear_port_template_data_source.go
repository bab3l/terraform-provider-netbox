// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &RearPortTemplateDataSource{}

// NewRearPortTemplateDataSource returns a new data source implementing the rear port template data source.

func NewRearPortTemplateDataSource() datasource.DataSource {
	return &RearPortTemplateDataSource{}
}

// RearPortTemplateDataSource defines the data source implementation.

type RearPortTemplateDataSource struct {
	client *netbox.APIClient
}

// RearPortTemplateDataSourceModel describes the data source data model.

type RearPortTemplateDataSourceModel struct {
	ID types.Int32 `tfsdk:"id"`

	DeviceType types.Int32 `tfsdk:"device_type"`

	ModuleType types.Int32 `tfsdk:"module_type"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	Color types.String `tfsdk:"color"`

	Positions types.Int32 `tfsdk:"positions"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.

func (d *RearPortTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rear_port_template"
}

// Schema defines the schema for the data source.

func (d *RearPortTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a rear port template in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the rear port template.",

				Optional: true,

				Computed: true,
			},

			"device_type": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the device type. Used with name for lookup when ID is not provided.",

				Optional: true,

				Computed: true,
			},

			"module_type": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the module type. Used with name for lookup when ID is not provided.",

				Optional: true,

				Computed: true,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the rear port template. Used with device_type or module_type for lookup when ID is not provided.",

				Optional: true,

				Computed: true,
			},

			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the rear port template.",

				Computed: true,
			},

			"type": schema.StringAttribute{
				MarkdownDescription: "The type of rear port.",

				Computed: true,
			},

			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the rear port in hex format.",

				Computed: true,
			},

			"positions": schema.Int32Attribute{
				MarkdownDescription: "Number of front ports that may be mapped to this rear port.",

				Computed: true,
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the rear port template.",

				Computed: true,
			},

			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the rear port template.",

				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.

func (d *RearPortTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RearPortTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RearPortTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var template *netbox.RearPortTemplate

	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():

		// Lookup by ID

		templateID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading rear port template by ID", map[string]interface{}{
			"id": templateID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimRearPortTemplatesRetrieve(ctx, templateID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError(

				"Error reading rear port template",

				utils.FormatAPIError(fmt.Sprintf("read rear port template ID %d", templateID), err, httpResp),
			)

			return
		}

		template = response

	case !data.Name.IsNull() && !data.Name.IsUnknown():

		// Lookup by device_type/module_type and name

		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading rear port template by name", map[string]interface{}{
			"name": name,
		})

		listReq := d.client.DcimAPI.DcimRearPortTemplatesList(ctx).Name([]string{name})

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

				"Error reading rear port template",

				utils.FormatAPIError(fmt.Sprintf("read rear port template by name %s", name), err, httpResp),
			)

			return
		}

		count := int(response.GetCount())

		if count == 0 {
			resp.Diagnostics.AddError(

				"Rear Port Template Not Found",

				fmt.Sprintf("No rear port template found with name: %s", name),
			)

			return
		}

		if count > 1 {
			resp.Diagnostics.AddError(

				"Multiple Rear Port Templates Found",

				fmt.Sprintf("Found %d rear port templates with name %s. Please specify device_type or module_type to narrow results, or use ID.", count, name),
			)

			return
		}

		template = &response.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either 'id' or 'name' must be specified to lookup a rear port template.",
		)

		return
	}

	// Map response to model

	d.mapResponseToModel(template, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.

func (d *RearPortTemplateDataSource) mapResponseToModel(template *netbox.RearPortTemplate, data *RearPortTemplateDataSourceModel) {
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

	// Map positions

	if positions, ok := template.GetPositionsOk(); ok && positions != nil {
		data.Positions = types.Int32Value(*positions)
	} else {
		data.Positions = types.Int32Null()
	}

	// Map description

	if desc, ok := template.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map display_name

	if template.GetDisplay() != "" {
		data.DisplayName = types.StringValue(template.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
