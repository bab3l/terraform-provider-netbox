// Package datasources contains Terraform data source implementations for the Netbox provider.
//

// This package integrates with the go-netbox OpenAPI client to provide
// read-only access to Netbox resources via Terraform data sources.

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
var _ datasource.DataSource = &DeviceTypeDataSource{}

func NewDeviceTypeDataSource() datasource.DataSource {
	return &DeviceTypeDataSource{}
}

// DeviceTypeDataSource defines the data source implementation.
type DeviceTypeDataSource struct {
	client *netbox.APIClient
}

// DeviceTypeDataSourceModel describes the data source data model.
type DeviceTypeDataSourceModel struct {
	ID                     types.String  `tfsdk:"id"`
	Manufacturer           types.String  `tfsdk:"manufacturer"`
	Model                  types.String  `tfsdk:"model"`
	Slug                   types.String  `tfsdk:"slug"`
	DefaultPlatform        types.String  `tfsdk:"default_platform"`
	PartNumber             types.String  `tfsdk:"part_number"`
	UHeight                types.Float64 `tfsdk:"u_height"`
	ExcludeFromUtilization types.Bool    `tfsdk:"exclude_from_utilization"`
	IsFullDepth            types.Bool    `tfsdk:"is_full_depth"`
	SubdeviceRole          types.String  `tfsdk:"subdevice_role"`
	Airflow                types.String  `tfsdk:"airflow"`
	Weight                 types.Float64 `tfsdk:"weight"`
	WeightUnit             types.String  `tfsdk:"weight_unit"`
	Description            types.String  `tfsdk:"description"`
	Comments               types.String  `tfsdk:"comments"`
	DisplayName            types.String  `tfsdk:"display_name"`
	Tags                   types.Set     `tfsdk:"tags"`
	CustomFields           types.Set     `tfsdk:"custom_fields"`
	DeviceCount            types.Int64   `tfsdk:"device_count"`
}

func (d *DeviceTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_type"
}

func (d *DeviceTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a device type in Netbox. Device types represent a particular make and model of device. You can identify the device type using `id`, `slug`, or `model`.",
		Attributes: map[string]schema.Attribute{
			"id":           nbschema.DSIDAttribute("device type"),
			"manufacturer": nbschema.DSComputedStringAttribute("The manufacturer of this device type. Returns the manufacturer slug."),
			"model": schema.StringAttribute{
				MarkdownDescription: "Model name of the device type. Can be used to identify the device type instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
			},
			"slug":                     nbschema.DSSlugAttribute("device type"),
			"default_platform":         nbschema.DSComputedStringAttribute("Default platform for devices of this type. Returns the platform slug."),
			"part_number":              nbschema.DSComputedStringAttribute("Discrete part number for this device type."),
			"u_height":                 nbschema.DSComputedFloat64Attribute("Height of the device type in rack units."),
			"exclude_from_utilization": nbschema.DSComputedBoolAttribute("Whether devices of this type are excluded when calculating rack utilization."),
			"is_full_depth":            nbschema.DSComputedBoolAttribute("Whether the device type consumes both front and rear rack faces."),
			"subdevice_role":           nbschema.DSComputedStringAttribute("Subdevice role (parent or child)."),
			"airflow":                  nbschema.DSComputedStringAttribute("Airflow direction for the device type."),
			"weight":                   nbschema.DSComputedFloat64Attribute("Weight of the device type."),
			"weight_unit":              nbschema.DSComputedStringAttribute("Unit of weight (kg, g, lb, oz)."),
			"description":              nbschema.DSComputedStringAttribute("Detailed description of the device type."),
			"comments":                 nbschema.DSComputedStringAttribute("Comments about the device type (supports Markdown)."),
			"display_name":             nbschema.DSComputedStringAttribute("The display name of the device type."),
			"tags":                     nbschema.DSTagsAttribute(),
			"custom_fields":            nbschema.DSCustomFieldsAttribute(),
			"device_count":             nbschema.DSComputedInt64Attribute("Number of devices of this type."),
		},
	}
}

func (d *DeviceTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *DeviceTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeviceTypeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var deviceType *netbox.DeviceType
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or model
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Search by ID
		deviceTypeID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading device type by ID", map[string]interface{}{
			"id": deviceTypeID,
		})

		// Parse the device type ID to int32 for the API call
		var deviceTypeIDInt int32
		if _, parseErr := fmt.Sscanf(deviceTypeID, "%d", &deviceTypeIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Device Type ID",
				fmt.Sprintf("Device Type ID must be a number, got: %s", deviceTypeID),
			)
			return
		}

		// Retrieve the device type via API
		deviceType, httpResp, err = d.client.DcimAPI.DcimDeviceTypesRetrieve(ctx, deviceTypeIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)

	case !data.Slug.IsNull() && !data.Slug.IsUnknown():
		// Search by slug
		deviceTypeSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading device type by slug", map[string]interface{}{
			"slug": deviceTypeSlug,
		})

		// List device types with slug filter
		var deviceTypes *netbox.PaginatedDeviceTypeList
		deviceTypes, httpResp, err = d.client.DcimAPI.DcimDeviceTypesList(ctx).Slug([]string{deviceTypeSlug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading device type",
				utils.FormatAPIError("read device type by slug", err, httpResp),
			)
			return
		}
		results := deviceTypes.GetResults()
		deviceTypeResult, ok := utils.ExpectSingleResult(
			results,
			"Device Type Not Found",
			fmt.Sprintf("No device type found with slug: %s", deviceTypeSlug),
			"Multiple Device Types Found",
			fmt.Sprintf("Multiple device types found with slug: %s. This should not happen as slugs should be unique.", deviceTypeSlug),
			&resp.Diagnostics,
		)
		if !ok {
			return
		}
		deviceType = deviceTypeResult

	case !data.Model.IsNull() && !data.Model.IsUnknown():
		// Search by model
		deviceTypeModel := data.Model.ValueString()
		tflog.Debug(ctx, "Reading device type by model", map[string]interface{}{
			"model": deviceTypeModel,
		})

		// List device types with model filter
		var deviceTypes *netbox.PaginatedDeviceTypeList
		deviceTypes, httpResp, err = d.client.DcimAPI.DcimDeviceTypesList(ctx).Model([]string{deviceTypeModel}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading device type",
				utils.FormatAPIError("read device type by model", err, httpResp),
			)
			return
		}
		results := deviceTypes.GetResults()
		deviceTypeResult, ok := utils.ExpectSingleResult(
			results,
			"Device Type Not Found",
			fmt.Sprintf("No device type found with model: %s", deviceTypeModel),
			"Multiple Device Types Found",
			fmt.Sprintf("Multiple device types found with model: %s. Specify a more precise filter.", deviceTypeModel),
			&resp.Diagnostics,
		)
		if !ok {
			return
		}
		deviceType = deviceTypeResult

	default:
		resp.Diagnostics.AddError(
			"Missing Device Type Identifier",
			"Either 'id', 'slug', or 'model' must be specified to identify the device type.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading device type",
			utils.FormatAPIError("read device type", err, httpResp),
		)
		return
	}
	if httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Device Type Not Found",
			"The specified device type was not found in Netbox.",
		)
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error reading device type",
			fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusOK, httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))
	data.Model = types.StringValue(deviceType.GetModel())
	data.Slug = types.StringValue(deviceType.GetSlug())

	// Handle manufacturer
	if manufacturer := deviceType.GetManufacturer(); manufacturer.Id != 0 {
		data.Manufacturer = types.StringValue(manufacturer.GetSlug())
	} else {
		data.Manufacturer = types.StringNull()
	}

	// Handle default_platform
	if deviceType.HasDefaultPlatform() && deviceType.GetDefaultPlatform().Id != 0 {
		platform := deviceType.GetDefaultPlatform()

		data.DefaultPlatform = types.StringValue(platform.GetSlug())
	} else {
		data.DefaultPlatform = types.StringNull()
	}

	// Handle part_number
	if deviceType.HasPartNumber() && deviceType.GetPartNumber() != "" {
		data.PartNumber = types.StringValue(deviceType.GetPartNumber())
	} else {
		data.PartNumber = types.StringNull()
	}

	// Handle u_height
	if deviceType.HasUHeight() {
		data.UHeight = types.Float64Value(deviceType.GetUHeight())
	} else {
		data.UHeight = types.Float64Value(1.0) // Default
	}

	// Handle exclude_from_utilization
	if deviceType.HasExcludeFromUtilization() {
		data.ExcludeFromUtilization = types.BoolValue(deviceType.GetExcludeFromUtilization())
	} else {
		data.ExcludeFromUtilization = types.BoolValue(false)
	}

	// Handle is_full_depth
	if deviceType.HasIsFullDepth() {
		data.IsFullDepth = types.BoolValue(deviceType.GetIsFullDepth())
	} else {
		data.IsFullDepth = types.BoolValue(true)
	}

	// Handle subdevice_role
	if deviceType.HasSubdeviceRole() && deviceType.SubdeviceRole.IsSet() {
		subdeviceRole := deviceType.GetSubdeviceRole()
		data.SubdeviceRole = types.StringValue(string(subdeviceRole.GetValue()))
	} else {
		data.SubdeviceRole = types.StringNull()
	}

	// Handle airflow
	if deviceType.HasAirflow() && deviceType.Airflow.IsSet() {
		airflow := deviceType.GetAirflow()
		data.Airflow = types.StringValue(string(airflow.GetValue()))
	} else {
		data.Airflow = types.StringNull()
	}

	// Handle weight
	if deviceType.HasWeight() && deviceType.Weight.Get() != nil {
		data.Weight = types.Float64Value(*deviceType.Weight.Get())
	} else {
		data.Weight = types.Float64Null()
	}

	// Handle weight_unit
	if deviceType.HasWeightUnit() && deviceType.WeightUnit.IsSet() {
		weightUnit := deviceType.GetWeightUnit()
		data.WeightUnit = types.StringValue(string(weightUnit.GetValue()))
	} else {
		data.WeightUnit = types.StringNull()
	}

	// Handle description
	if deviceType.HasDescription() && deviceType.GetDescription() != "" {
		data.Description = types.StringValue(deviceType.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments
	if deviceType.HasComments() && deviceType.GetComments() != "" {
		data.Comments = types.StringValue(deviceType.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if deviceType.HasTags() {
		tags := utils.NestedTagsToTagModels(deviceType.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	data.CustomFields = utils.CustomFieldsSetFromAPI(ctx, deviceType.HasCustomFields(), deviceType.GetCustomFields(), &resp.Diagnostics)

	// Handle device_count (read-only, always present)
	data.DeviceCount = types.Int64Value(deviceType.GetDeviceCount())

	// Handle display_name
	if deviceType.GetDisplay() != "" {
		data.DisplayName = types.StringValue(deviceType.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
