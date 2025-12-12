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
var _ datasource.DataSource = &InterfaceTemplateDataSource{}

// NewInterfaceTemplateDataSource returns a new data source implementing the interface template data source.
func NewInterfaceTemplateDataSource() datasource.DataSource {
	return &InterfaceTemplateDataSource{}
}

// InterfaceTemplateDataSource defines the data source implementation.
type InterfaceTemplateDataSource struct {
	client *netbox.APIClient
}

// InterfaceTemplateDataSourceModel describes the data source data model.
type InterfaceTemplateDataSourceModel struct {
	ID          types.Int32  `tfsdk:"id"`
	DeviceType  types.Int32  `tfsdk:"device_type"`
	ModuleType  types.Int32  `tfsdk:"module_type"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Type        types.String `tfsdk:"type"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	MgmtOnly    types.Bool   `tfsdk:"mgmt_only"`
	Description types.String `tfsdk:"description"`
	Bridge      types.Int32  `tfsdk:"bridge"`
	PoeMode     types.String `tfsdk:"poe_mode"`
	PoeType     types.String `tfsdk:"poe_type"`
	RfRole      types.String `tfsdk:"rf_role"`
}

// Metadata returns the data source type name.
func (d *InterfaceTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_template"
}

// Schema defines the schema for the data source.
func (d *InterfaceTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an interface template in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the interface template to retrieve. If specified, other filter attributes are ignored.",
				Optional:            true,
				Computed:            true,
			},
			"device_type": schema.Int32Attribute{
				MarkdownDescription: "Filter by device type ID.",
				Optional:            true,
				Computed:            true,
			},
			"module_type": schema.Int32Attribute{
				MarkdownDescription: "Filter by module type ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Filter by interface template name.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "The physical label of the interface template.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the interface.",
				Optional:            true,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is enabled by default.",
				Computed:            true,
			},
			"mgmt_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is for management only.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the interface template.",
				Computed:            true,
			},
			"bridge": schema.Int32Attribute{
				MarkdownDescription: "The ID of the bridge interface template this interface belongs to.",
				Computed:            true,
			},
			"poe_mode": schema.StringAttribute{
				MarkdownDescription: "PoE mode (pd or pse).",
				Computed:            true,
			},
			"poe_type": schema.StringAttribute{
				MarkdownDescription: "PoE type.",
				Computed:            true,
			},
			"rf_role": schema.StringAttribute{
				MarkdownDescription: "Wireless role (ap or station).",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *InterfaceTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *InterfaceTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var template *netbox.InterfaceTemplate

	// If ID is provided, look up by ID directly
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		templateID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading interface template by ID", map[string]interface{}{
			"id": templateID,
		})

		result, httpResp, err := d.client.DcimAPI.DcimInterfaceTemplatesRetrieve(ctx, templateID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading interface template",
				utils.FormatAPIError(fmt.Sprintf("read interface template ID %d", templateID), err, httpResp),
			)
			return
		}
		template = result
	} else {
		// Build search request with filters
		listReq := d.client.DcimAPI.DcimInterfaceTemplatesList(ctx)

		if !data.Name.IsNull() && !data.Name.IsUnknown() {
			listReq = listReq.Name([]string{data.Name.ValueString()})
		}
		if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
			deviceTypeID := data.DeviceType.ValueInt32()
			listReq = listReq.DeviceTypeId([]*int32{&deviceTypeID})
		}
		if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {
			moduleTypeID := data.ModuleType.ValueInt32()
			listReq = listReq.ModuleTypeId([]*int32{&moduleTypeID})
		}
		if !data.Type.IsNull() && !data.Type.IsUnknown() {
			listReq = listReq.Type_([]string{data.Type.ValueString()})
		}

		tflog.Debug(ctx, "Searching for interface template")

		result, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error searching for interface template",
				utils.FormatAPIError("search for interface template", err, httpResp),
			)
			return
		}

		if result.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"No interface template found",
				"No interface template matching the specified criteria was found.",
			)
			return
		}

		if result.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple interface templates found",
				fmt.Sprintf("Found %d interface templates matching the specified criteria. Please refine your search.", result.GetCount()),
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
func (d *InterfaceTemplateDataSource) mapResponseToModel(template *netbox.InterfaceTemplate, data *InterfaceTemplateDataSourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())
	data.Type = types.StringValue(string(template.Type.GetValue()))

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

	// Map label
	if label, ok := template.GetLabelOk(); ok && label != nil {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringValue("")
	}

	// Map enabled
	if enabled, ok := template.GetEnabledOk(); ok && enabled != nil {
		data.Enabled = types.BoolValue(*enabled)
	} else {
		data.Enabled = types.BoolValue(true)
	}

	// Map mgmt_only
	if mgmtOnly, ok := template.GetMgmtOnlyOk(); ok && mgmtOnly != nil {
		data.MgmtOnly = types.BoolValue(*mgmtOnly)
	} else {
		data.MgmtOnly = types.BoolValue(false)
	}

	// Map description
	if desc, ok := template.GetDescriptionOk(); ok && desc != nil {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringValue("")
	}

	// Map bridge
	if template.Bridge.IsSet() && template.Bridge.Get() != nil {
		data.Bridge = types.Int32Value(template.Bridge.Get().Id)
	} else {
		data.Bridge = types.Int32Null()
	}

	// Map poe_mode
	if poeMode, ok := template.GetPoeModeOk(); ok && poeMode != nil && poeMode.Value != nil {
		data.PoeMode = types.StringValue(string(*poeMode.Value))
	} else {
		data.PoeMode = types.StringNull()
	}

	// Map poe_type
	if poeType, ok := template.GetPoeTypeOk(); ok && poeType != nil && poeType.Value != nil {
		data.PoeType = types.StringValue(string(*poeType.Value))
	} else {
		data.PoeType = types.StringNull()
	}

	// Map rf_role
	if rfRole, ok := template.GetRfRoleOk(); ok && rfRole != nil && rfRole.Value != nil {
		data.RfRole = types.StringValue(string(*rfRole.Value))
	} else {
		data.RfRole = types.StringNull()
	}
}
