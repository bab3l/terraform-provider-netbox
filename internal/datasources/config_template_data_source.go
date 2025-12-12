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
var _ datasource.DataSource = &ConfigTemplateDataSource{}

// NewConfigTemplateDataSource returns a new data source implementing the config template data source.
func NewConfigTemplateDataSource() datasource.DataSource {
	return &ConfigTemplateDataSource{}
}

// ConfigTemplateDataSource defines the data source implementation.
type ConfigTemplateDataSource struct {
	client *netbox.APIClient
}

// ConfigTemplateDataSourceModel describes the data source data model.
type ConfigTemplateDataSourceModel struct {
	ID           types.Int32  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	TemplateCode types.String `tfsdk:"template_code"`
	DataSource   types.Int32  `tfsdk:"data_source"`
	DataPath     types.String `tfsdk:"data_path"`
}

// Metadata returns the data source type name.
func (d *ConfigTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_template"
}

// Schema defines the schema for the data source.
func (d *ConfigTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a config template in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the config template to retrieve. If specified, other filter attributes are ignored.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Filter by config template name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the config template.",
				Computed:            true,
			},
			"template_code": schema.StringAttribute{
				MarkdownDescription: "Jinja2 template code.",
				Computed:            true,
			},
			"data_source": schema.Int32Attribute{
				MarkdownDescription: "The ID of the data source the template is synced from.",
				Computed:            true,
			},
			"data_path": schema.StringAttribute{
				MarkdownDescription: "Path to remote file (relative to data source root).",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ConfigTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ConfigTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConfigTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var template *netbox.ConfigTemplate

	// If ID is provided, look up by ID directly
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		templateID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading config template by ID", map[string]interface{}{
			"id": templateID,
		})

		result, httpResp, err := d.client.ExtrasAPI.ExtrasConfigTemplatesRetrieve(ctx, templateID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading config template",
				utils.FormatAPIError(fmt.Sprintf("read config template ID %d", templateID), err, httpResp),
			)
			return
		}
		template = result
	} else {
		// Build search request with filters
		listReq := d.client.ExtrasAPI.ExtrasConfigTemplatesList(ctx)

		if !data.Name.IsNull() && !data.Name.IsUnknown() {
			listReq = listReq.Name([]string{data.Name.ValueString()})
		}

		tflog.Debug(ctx, "Searching for config template")

		result, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error searching for config template",
				utils.FormatAPIError("search for config template", err, httpResp),
			)
			return
		}

		if result.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"No config template found",
				"No config template matching the specified criteria was found.",
			)
			return
		}

		if result.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple config templates found",
				fmt.Sprintf("Found %d config templates matching the specified criteria. Please refine your search.", result.GetCount()),
			)
			return
		}

		template = &result.GetResults()[0]
	}

	// Map response to model
	d.mapResponseToModel(template, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ConfigTemplateDataSource) mapResponseToModel(template *netbox.ConfigTemplate, data *ConfigTemplateDataSourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())
	data.TemplateCode = types.StringValue(template.GetTemplateCode())

	// Map description
	if desc, ok := template.GetDescriptionOk(); ok && desc != nil {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringValue("")
	}

	// Map data source
	if dataSource, ok := template.GetDataSourceOk(); ok && dataSource != nil {
		data.DataSource = types.Int32Value(dataSource.GetId())
	} else {
		data.DataSource = types.Int32Null()
	}

	// Map data path
	data.DataPath = types.StringValue(template.GetDataPath())
}
