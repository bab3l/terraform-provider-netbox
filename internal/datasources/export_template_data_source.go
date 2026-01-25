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
var _ datasource.DataSource = &ExportTemplateDataSource{}

// NewExportTemplateDataSource returns a new data source implementing the export template data source.
func NewExportTemplateDataSource() datasource.DataSource {
	return &ExportTemplateDataSource{}
}

// ExportTemplateDataSource defines the data source implementation.
type ExportTemplateDataSource struct {
	client *netbox.APIClient
}

// ExportTemplateDataSourceModel describes the data source data model.
type ExportTemplateDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ObjectTypes   types.List   `tfsdk:"object_types"`
	Description   types.String `tfsdk:"description"`
	TemplateCode  types.String `tfsdk:"template_code"`
	MimeType      types.String `tfsdk:"mime_type"`
	FileExtension types.String `tfsdk:"file_extension"`
	AsAttachment  types.Bool   `tfsdk:"as_attachment"`
}

// Metadata returns the data source type name.
func (d *ExportTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export_template"
}

// Schema defines the schema for the data source.
func (d *ExportTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an export template in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id":   nbschema.DSIDAttribute("export template"),
			"name": nbschema.DSNameAttribute("export template"),
			"object_types": schema.ListAttribute{
				MarkdownDescription: "List of object types this template applies to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description":    nbschema.DSComputedStringAttribute("Description of the export template."),
			"template_code":  nbschema.DSComputedStringAttribute("Jinja2 template code."),
			"mime_type":      nbschema.DSComputedStringAttribute("MIME type for the rendered output."),
			"file_extension": nbschema.DSComputedStringAttribute("Extension to append to the rendered filename."),
			"as_attachment":  nbschema.DSComputedBoolAttribute("Whether to download file as attachment."),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ExportTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the export template data.
func (d *ExportTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExportTemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var exportTemplate *netbox.ExportTemplate
	var httpResp *http.Response
	var err error

	// Lookup by ID
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		id, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), parseErr),
			)
			return
		}
		tflog.Debug(ctx, "Reading export template by ID", map[string]interface{}{
			"id": id,
		})
		exportTemplate, httpResp, err = d.client.ExtrasAPI.ExtrasExportTemplatesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)

	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by name
		tflog.Debug(ctx, "Reading export template by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		list, listResp, listErr := d.client.ExtrasAPI.ExtrasExportTemplatesList(ctx).
			Name([]string{data.Name.ValueString()}).
			Execute()
		defer utils.CloseResponseBody(listResp)
		httpResp = listResp
		err = listErr
		if err == nil {
			results := list.GetResults()
			exportTemplateResult, ok := utils.ExpectSingleResult(
				results,
				"Not Found",
				fmt.Sprintf("No export template found with name: %s", data.Name.ValueString()),
				"Multiple Found",
				fmt.Sprintf("Multiple export templates found with name: %s. Please use id for a more specific lookup.", data.Name.ValueString()),
				&resp.Diagnostics,
			)
			if !ok {
				return
			}
			exportTemplate = exportTemplateResult
		}

	default:
		resp.Diagnostics.AddError(
			"Missing Identifier",
			"Either 'id' or 'name' must be specified to look up an export template.",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading export template",
			utils.FormatAPIError("read export template", err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, exportTemplate, &data, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state.
func (d *ExportTemplateDataSource) mapResponseToState(ctx context.Context, exportTemplate *netbox.ExportTemplate, data *ExportTemplateDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", exportTemplate.GetId()))
	data.Name = types.StringValue(exportTemplate.GetName())
	data.TemplateCode = types.StringValue(exportTemplate.GetTemplateCode())

	// Handle object types
	if len(exportTemplate.ObjectTypes) > 0 {
		objectTypesList, diags := types.ListValueFrom(ctx, types.StringType, exportTemplate.ObjectTypes)
		resp.Diagnostics.Append(diags...)
		data.ObjectTypes = objectTypesList
	} else {
		data.ObjectTypes = types.ListNull(types.StringType)
	}

	// Handle description
	if exportTemplate.HasDescription() && exportTemplate.GetDescription() != "" {
		data.Description = types.StringValue(exportTemplate.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle mime type
	if exportTemplate.HasMimeType() && exportTemplate.GetMimeType() != "" {
		data.MimeType = types.StringValue(exportTemplate.GetMimeType())
	} else {
		data.MimeType = types.StringNull()
	}

	// Handle file extension
	if exportTemplate.HasFileExtension() && exportTemplate.GetFileExtension() != "" {
		data.FileExtension = types.StringValue(exportTemplate.GetFileExtension())
	} else {
		data.FileExtension = types.StringNull()
	}

	// Handle as_attachment
	if exportTemplate.HasAsAttachment() {
		data.AsAttachment = types.BoolValue(exportTemplate.GetAsAttachment())
	} else {
		data.AsAttachment = types.BoolValue(true) // Default
	}
}
