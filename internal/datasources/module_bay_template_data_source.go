// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &ModuleBayTemplateDataSource{}

	_ datasource.DataSourceWithConfigure = &ModuleBayTemplateDataSource{}
)

// NewModuleBayTemplateDataSource returns a new data source implementing the module bay template data source.

func NewModuleBayTemplateDataSource() datasource.DataSource {

	return &ModuleBayTemplateDataSource{}

}

// ModuleBayTemplateDataSource defines the data source implementation.

type ModuleBayTemplateDataSource struct {
	client *netbox.APIClient
}

// ModuleBayTemplateDataSourceModel describes the data source data model.

type ModuleBayTemplateDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	DeviceTypeID types.String `tfsdk:"device_type_id"`

	ModuleType types.String `tfsdk:"module_type"`

	ModuleTypeID types.String `tfsdk:"module_type_id"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Position types.String `tfsdk:"position"`

	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.

func (d *ModuleBayTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_module_bay_template"

}

// Schema defines the schema for the data source.

func (d *ModuleBayTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a module bay template in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the module bay template.",

				Required: true,
			},

			"device_type": schema.StringAttribute{

				MarkdownDescription: "The model name of the device type this module bay template belongs to.",

				Computed: true,
			},

			"device_type_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the device type this module bay template belongs to.",

				Computed: true,
			},

			"module_type": schema.StringAttribute{

				MarkdownDescription: "The model name of the module type this module bay template belongs to.",

				Computed: true,
			},

			"module_type_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the module type this module bay template belongs to.",

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the module bay template.",

				Computed: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the module bay template.",

				Computed: true,
			},

			"position": schema.StringAttribute{

				MarkdownDescription: "Identifier to reference when renaming installed components.",

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the module bay template.",

				Computed: true,
			},
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *ModuleBayTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read reads the data source.

func (d *ModuleBayTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data ModuleBayTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	tflog.Debug(ctx, "Reading module bay template", map[string]interface{}{"id": id})

	// Read from API

	result, httpResp, err := d.client.DcimAPI.DcimModuleBayTemplatesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading module bay template",

			utils.FormatAPIError(fmt.Sprintf("read module bay template ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to state

	d.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapToState maps the API response to the Terraform state.

func (d *ModuleBayTemplateDataSource) mapToState(ctx context.Context, result *netbox.ModuleBayTemplate, data *ModuleBayTemplateDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	// Map device type

	if result.HasDeviceType() && result.GetDeviceType().Id != 0 {

		deviceType := result.GetDeviceType()

		data.DeviceType = types.StringValue(deviceType.GetModel())

		data.DeviceTypeID = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))

	} else {

		data.DeviceType = types.StringNull()

		data.DeviceTypeID = types.StringNull()

	}

	// Map module type

	if result.HasModuleType() && result.GetModuleType().Id != 0 {

		moduleType := result.GetModuleType()

		data.ModuleType = types.StringValue(moduleType.GetModel())

		data.ModuleTypeID = types.StringValue(fmt.Sprintf("%d", moduleType.GetId()))

	} else {

		data.ModuleType = types.StringNull()

		data.ModuleTypeID = types.StringNull()

	}

	// Map label

	if result.HasLabel() && result.GetLabel() != "" {

		data.Label = types.StringValue(result.GetLabel())

	} else {

		data.Label = types.StringNull()

	}

	// Map position

	if result.HasPosition() && result.GetPosition() != "" {

		data.Position = types.StringValue(result.GetPosition())

	} else {

		data.Position = types.StringNull()

	}

	// Map description

	if result.HasDescription() && result.GetDescription() != "" {

		data.Description = types.StringValue(result.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

}
