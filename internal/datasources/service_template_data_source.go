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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServiceTemplateDataSource{}

// NewServiceTemplateDataSource returns a new data source implementing the service template data source.
func NewServiceTemplateDataSource() datasource.DataSource {
	return &ServiceTemplateDataSource{}
}

// ServiceTemplateDataSource defines the data source implementation.
type ServiceTemplateDataSource struct {
	client *netbox.APIClient
}

// ServiceTemplateDataSourceModel describes the data source data model.
type ServiceTemplateDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Protocol     types.String `tfsdk:"protocol"`
	Ports        types.List   `tfsdk:"ports"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ServiceTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_template"
}

// Schema defines the schema for the data source.
func (d *ServiceTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a service template in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id":       nbschema.DSIDAttribute("service template"),
			"name":     nbschema.DSNameAttribute("service template"),
			"protocol": nbschema.DSComputedStringAttribute("Protocol used by the service template (tcp, udp, sctp)."),
			"ports": schema.ListAttribute{
				MarkdownDescription: "List of port numbers the service template listens on.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"description": nbschema.DSComputedStringAttribute("Description of the service template."),
			"comments":    nbschema.DSComputedStringAttribute("Comments about the service template."),
			"tags":        nbschema.DSTagsAttribute(),
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this service template.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field.",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ServiceTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the service template data.
func (d *ServiceTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serviceTemplate *netbox.ServiceTemplate
	var httpResp *http.Response
	var err error

	// Lookup by ID or name
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

		tflog.Debug(ctx, "Reading service template by ID", map[string]interface{}{
			"id": id,
		})

		serviceTemplate, httpResp, err = d.client.IpamAPI.IpamServiceTemplatesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by name
		tflog.Debug(ctx, "Reading service template by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		list, listResp, listErr := d.client.IpamAPI.IpamServiceTemplatesList(ctx).
			Name([]string{data.Name.ValueString()}).
			Execute()
		httpResp = listResp
		defer utils.CloseResponseBody(httpResp)
		err = listErr

		if err == nil {
			results := list.GetResults()
			if len(results) == 0 {
				resp.Diagnostics.AddError(
					"Not Found",
					fmt.Sprintf("No service template found with name: %s", data.Name.ValueString()),
				)
				return
			}
			if len(results) > 1 {
				resp.Diagnostics.AddError(
					"Multiple Found",
					fmt.Sprintf("Multiple service templates found with name: %s. Please use id for a more specific lookup.", data.Name.ValueString()),
				)
				return
			}
			serviceTemplate = &results[0]
		}
	default:
		resp.Diagnostics.AddError(
			"Missing Identifier",
			"Either 'id' or 'name' must be specified to look up a service template.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading service template",
			utils.FormatAPIError("read service template", err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, serviceTemplate, &data, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state.
func (d *ServiceTemplateDataSource) mapResponseToState(ctx context.Context, serviceTemplate *netbox.ServiceTemplate, data *ServiceTemplateDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", serviceTemplate.GetId()))
	data.Name = types.StringValue(serviceTemplate.GetName())

	// Handle protocol
	if serviceTemplate.HasProtocol() {
		protocol := serviceTemplate.GetProtocol()
		data.Protocol = types.StringValue(string(protocol.GetValue()))
	} else {
		data.Protocol = types.StringNull()
	}

	// Handle ports
	if serviceTemplate.Ports != nil {
		var ports []int64
		for _, p := range serviceTemplate.Ports {
			ports = append(ports, int64(p))
		}
		portsList, diagErr := types.ListValueFrom(ctx, types.Int64Type, ports)
		diags.Append(diagErr...)
		data.Ports = portsList
	} else {
		data.Ports = types.ListNull(types.Int64Type)
	}

	// Handle description
	if serviceTemplate.HasDescription() && serviceTemplate.GetDescription() != "" {
		data.Description = types.StringValue(serviceTemplate.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments
	if serviceTemplate.HasComments() && serviceTemplate.GetComments() != "" {
		data.Comments = types.StringValue(serviceTemplate.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if serviceTemplate.HasTags() && len(serviceTemplate.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(serviceTemplate.GetTags())
		tagsValue, diagErr := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(diagErr...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if serviceTemplate.HasCustomFields() {
		var existingModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			diagErr := data.CustomFields.ElementsAs(ctx, &existingModels, false)
			diags.Append(diagErr...)
		}
		customFields := utils.MapToCustomFieldModels(serviceTemplate.GetCustomFields(), existingModels)
		if len(customFields) > 0 {
			customFieldsValue, diagErr := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
			diags.Append(diagErr...)
			data.CustomFields = customFieldsValue
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
