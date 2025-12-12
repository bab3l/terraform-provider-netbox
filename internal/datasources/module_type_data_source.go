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
var _ datasource.DataSource = &ModuleTypeDataSource{}

// NewModuleTypeDataSource returns a new data source implementing the module type data source.
func NewModuleTypeDataSource() datasource.DataSource {
	return &ModuleTypeDataSource{}
}

// ModuleTypeDataSource defines the data source implementation.
type ModuleTypeDataSource struct {
	client *netbox.APIClient
}

// ModuleTypeDataSourceModel describes the data source data model.
type ModuleTypeDataSourceModel struct {
	ID             types.Int32   `tfsdk:"id"`
	Model          types.String  `tfsdk:"model"`
	ManufacturerID types.Int32   `tfsdk:"manufacturer_id"`
	Manufacturer   types.String  `tfsdk:"manufacturer"`
	PartNumber     types.String  `tfsdk:"part_number"`
	Airflow        types.String  `tfsdk:"airflow"`
	Weight         types.Float64 `tfsdk:"weight"`
	WeightUnit     types.String  `tfsdk:"weight_unit"`
	Description    types.String  `tfsdk:"description"`
	Comments       types.String  `tfsdk:"comments"`
}

// Metadata returns the data source type name.
func (d *ModuleTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_type"
}

// Schema defines the schema for the data source.
func (d *ModuleTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a module type in NetBox. Module types define hardware module specifications (model, manufacturer, etc.) that can be instantiated as modules within devices.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the module type.",
				Optional:            true,
				Computed:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model name/number of the module type. Used for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"manufacturer_id": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the manufacturer. Used with model for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"manufacturer": schema.StringAttribute{
				MarkdownDescription: "The name of the manufacturer.",
				Computed:            true,
			},
			"part_number": schema.StringAttribute{
				MarkdownDescription: "Discrete part number (optional).",
				Computed:            true,
			},
			"airflow": schema.StringAttribute{
				MarkdownDescription: "Airflow direction.",
				Computed:            true,
			},
			"weight": schema.Float64Attribute{
				MarkdownDescription: "Weight of the module.",
				Computed:            true,
			},
			"weight_unit": schema.StringAttribute{
				MarkdownDescription: "Unit for weight measurement.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the module type.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ModuleTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ModuleTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModuleTypeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var moduleType *netbox.ModuleType

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		// Lookup by ID
		typeID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading module type by ID", map[string]interface{}{
			"id": typeID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimModuleTypesRetrieve(ctx, typeID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading module type",
				utils.FormatAPIError(fmt.Sprintf("read module type ID %d", typeID), err, httpResp),
			)
			return
		}
		moduleType = response
	} else if !data.Model.IsNull() && !data.Model.IsUnknown() {
		// Lookup by model (and optionally manufacturer_id)
		model := data.Model.ValueString()

		tflog.Debug(ctx, "Reading module type by model", map[string]interface{}{
			"model": model,
		})

		listReq := d.client.DcimAPI.DcimModuleTypesList(ctx).Model([]string{model})

		if !data.ManufacturerID.IsNull() && !data.ManufacturerID.IsUnknown() {
			listReq = listReq.ManufacturerId([]int32{data.ManufacturerID.ValueInt32()})
		}

		response, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading module type",
				utils.FormatAPIError(fmt.Sprintf("read module type by model %s", model), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Module Type Not Found",
				fmt.Sprintf("No module type found with model: %s", model),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Module Types Found",
				fmt.Sprintf("Found %d module types with model %s. Please provide manufacturer_id or use ID to select a specific one.", count, model),
			)
			return
		}

		moduleType = &response.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'model' must be specified to lookup a module type.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(moduleType, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ModuleTypeDataSource) mapResponseToModel(moduleType *netbox.ModuleType, data *ModuleTypeDataSourceModel) {
	data.ID = types.Int32Value(moduleType.GetId())
	data.Model = types.StringValue(moduleType.GetModel())

	// Map manufacturer
	if mfr := moduleType.GetManufacturer(); mfr.Id != 0 {
		data.ManufacturerID = types.Int32Value(mfr.GetId())
		data.Manufacturer = types.StringValue(mfr.GetName())
	}

	// Map part_number
	if partNum, ok := moduleType.GetPartNumberOk(); ok && partNum != nil && *partNum != "" {
		data.PartNumber = types.StringValue(*partNum)
	} else {
		data.PartNumber = types.StringNull()
	}

	// Map airflow
	if moduleType.Airflow.IsSet() && moduleType.Airflow.Get() != nil {
		data.Airflow = types.StringValue(string(moduleType.Airflow.Get().GetValue()))
	} else {
		data.Airflow = types.StringNull()
	}

	// Map weight
	if moduleType.Weight.IsSet() && moduleType.Weight.Get() != nil {
		data.Weight = types.Float64Value(*moduleType.Weight.Get())
	} else {
		data.Weight = types.Float64Null()
	}

	// Map weight_unit
	if moduleType.WeightUnit.IsSet() && moduleType.WeightUnit.Get() != nil {
		data.WeightUnit = types.StringValue(string(moduleType.WeightUnit.Get().GetValue()))
	} else {
		data.WeightUnit = types.StringNull()
	}

	// Map description
	if desc, ok := moduleType.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := moduleType.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}
}
