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
var _ datasource.DataSource = &PowerOutletTemplateDataSource{}

// NewPowerOutletTemplateDataSource returns a new data source implementing the power outlet template data source.
func NewPowerOutletTemplateDataSource() datasource.DataSource {
	return &PowerOutletTemplateDataSource{}
}

// PowerOutletTemplateDataSource defines the data source implementation.
type PowerOutletTemplateDataSource struct {
	client *netbox.APIClient
}

// PowerOutletTemplateDataSourceModel describes the data source data model.
type PowerOutletTemplateDataSourceModel struct {
	ID          types.Int32  `tfsdk:"id"`
	DeviceType  types.Int32  `tfsdk:"device_type"`
	ModuleType  types.Int32  `tfsdk:"module_type"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Type        types.String `tfsdk:"type"`
	PowerPort   types.Int32  `tfsdk:"power_port"`
	FeedLeg     types.String `tfsdk:"feed_leg"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *PowerOutletTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_outlet_template"
}

// Schema defines the schema for the data source.
func (d *PowerOutletTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a power outlet template in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the power outlet template.",
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
				MarkdownDescription: "The name of the power outlet template. Used with device_type or module_type for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the power outlet template.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of power outlet.",
				Computed:            true,
			},
			"power_port": schema.Int32Attribute{
				MarkdownDescription: "The power port template ID that feeds this outlet.",
				Computed:            true,
			},
			"feed_leg": schema.StringAttribute{
				MarkdownDescription: "Phase leg for three-phase power.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the power outlet template.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *PowerOutletTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PowerOutletTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PowerOutletTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var template *netbox.PowerOutletTemplate

	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Lookup by ID
		templateID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading power outlet template by ID", map[string]interface{}{
			"id": templateID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimPowerOutletTemplatesRetrieve(ctx, templateID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading power outlet template",
				utils.FormatAPIError(fmt.Sprintf("read power outlet template ID %d", templateID), err, httpResp),
			)
			return
		}
		template = response
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by device_type/module_type and name
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading power outlet template by name", map[string]interface{}{
			"name": name,
		})

		listReq := d.client.DcimAPI.DcimPowerOutletTemplatesList(ctx).Name([]string{name})

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
				"Error reading power outlet template",
				utils.FormatAPIError(fmt.Sprintf("read power outlet template by name %s", name), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Power Outlet Template Not Found",
				fmt.Sprintf("No power outlet template found with name: %s", name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Power Outlet Templates Found",
				fmt.Sprintf("Found %d power outlet templates with name %s. Please specify device_type or module_type to narrow results, or use ID.", count, name),
			)
			return
		}

		template = &response.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to lookup a power outlet template.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(template, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *PowerOutletTemplateDataSource) mapResponseToModel(template *netbox.PowerOutletTemplate, data *PowerOutletTemplateDataSourceModel) {
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

	// Map label
	if label, ok := template.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map type
	if template.Type.IsSet() && template.Type.Get() != nil {
		data.Type = types.StringValue(string(template.Type.Get().GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map power_port
	if template.PowerPort.IsSet() && template.PowerPort.Get() != nil {
		data.PowerPort = types.Int32Value(template.PowerPort.Get().Id)
	} else {
		data.PowerPort = types.Int32Null()
	}

	// Map feed_leg
	if template.FeedLeg.IsSet() && template.FeedLeg.Get() != nil {
		data.FeedLeg = types.StringValue(string(template.FeedLeg.Get().GetValue()))
	} else {
		data.FeedLeg = types.StringNull()
	}

	// Map description
	if desc, ok := template.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}
}
