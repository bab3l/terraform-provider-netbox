// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RackDataSource{}

func NewRackDataSource() datasource.DataSource {
	return &RackDataSource{}
}

// RackDataSource defines the data source implementation.
type RackDataSource struct {
	client *netbox.APIClient
}

// RackDataSourceModel describes the data source data model.
type RackDataSourceModel struct {
	// Lookup fields (one required)
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Computed fields from the rack
	Site          types.String `tfsdk:"site"`
	SiteID        types.String `tfsdk:"site_id"`
	Location      types.String `tfsdk:"location"`
	LocationID    types.String `tfsdk:"location_id"`
	Tenant        types.String `tfsdk:"tenant"`
	TenantID      types.String `tfsdk:"tenant_id"`
	Status        types.String `tfsdk:"status"`
	Role          types.String `tfsdk:"role"`
	RoleID        types.String `tfsdk:"role_id"`
	Serial        types.String `tfsdk:"serial"`
	AssetTag      types.String `tfsdk:"asset_tag"`
	RackType      types.String `tfsdk:"rack_type"`
	RackTypeID    types.String `tfsdk:"rack_type_id"`
	FormFactor    types.String `tfsdk:"form_factor"`
	Width         types.String `tfsdk:"width"`
	UHeight       types.String `tfsdk:"u_height"`
	StartingUnit  types.String `tfsdk:"starting_unit"`
	Weight        types.String `tfsdk:"weight"`
	MaxWeight     types.String `tfsdk:"max_weight"`
	WeightUnit    types.String `tfsdk:"weight_unit"`
	DescUnits     types.Bool   `tfsdk:"desc_units"`
	OuterWidth    types.String `tfsdk:"outer_width"`
	OuterDepth    types.String `tfsdk:"outer_depth"`
	OuterUnit     types.String `tfsdk:"outer_unit"`
	MountingDepth types.String `tfsdk:"mounting_depth"`
	Airflow       types.String `tfsdk:"airflow"`
	Description   types.String `tfsdk:"description"`
	Comments      types.String `tfsdk:"comments"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

func (d *RackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack"
}

func (d *RackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a rack in Netbox. Racks represent physical equipment enclosures used to organize network infrastructure within a site or location. You can identify the rack using `id` or `name`.",

		Attributes: map[string]schema.Attribute{
			// Lookup fields
			"id":   nbschema.DSIDAttribute("rack"),
			"name": nbschema.DSNameAttribute("rack"),

			// Reference computed fields with name/id pairs
			"site":         nbschema.DSComputedStringAttribute("Name of the site where this rack is located."),
			"site_id":      nbschema.DSComputedStringAttribute("ID of the site where this rack is located."),
			"location":     nbschema.DSComputedStringAttribute("Name of the location within the site."),
			"location_id":  nbschema.DSComputedStringAttribute("ID of the location within the site."),
			"tenant":       nbschema.DSComputedStringAttribute("Name of the tenant that owns this rack."),
			"tenant_id":    nbschema.DSComputedStringAttribute("ID of the tenant that owns this rack."),
			"role":         nbschema.DSComputedStringAttribute("Name of the functional role of the rack."),
			"role_id":      nbschema.DSComputedStringAttribute("ID of the functional role of the rack."),
			"rack_type":    nbschema.DSComputedStringAttribute("Model/name of the rack type."),
			"rack_type_id": nbschema.DSComputedStringAttribute("ID of the rack type."),

			// Status and physical attributes
			"status":         nbschema.DSComputedStringAttribute("Operational status of the rack (`reserved`, `available`, `planned`, `active`, `deprecated`)."),
			"serial":         nbschema.DSComputedStringAttribute("Serial number of the rack."),
			"asset_tag":      nbschema.DSComputedStringAttribute("Unique asset tag for the rack."),
			"form_factor":    nbschema.DSComputedStringAttribute("Physical form factor of the rack (`2-post-frame`, `4-post-frame`, `4-post-cabinet`, `wall-frame`, `wall-frame-vertical`, `wall-cabinet`, `wall-cabinet-vertical`)."),
			"width":          nbschema.DSComputedStringAttribute("Rail-to-rail width of the rack in inches (`10`, `19`, `21`, `23`)."),
			"u_height":       nbschema.DSComputedStringAttribute("Height of the rack in rack units."),
			"starting_unit":  nbschema.DSComputedStringAttribute("Starting unit number for the rack (bottom)."),
			"weight":         nbschema.DSComputedStringAttribute("Weight of the rack itself."),
			"max_weight":     nbschema.DSComputedStringAttribute("Maximum weight capacity of the rack."),
			"weight_unit":    nbschema.DSComputedStringAttribute("Unit of measurement for weight (`kg`, `g`, `lb`, `oz`)."),
			"desc_units":     nbschema.DSComputedBoolAttribute("If true, rack units are numbered in descending order (top to bottom)."),
			"outer_width":    nbschema.DSComputedStringAttribute("Outer width of the rack."),
			"outer_depth":    nbschema.DSComputedStringAttribute("Outer depth of the rack."),
			"outer_unit":     nbschema.DSComputedStringAttribute("Unit of measurement for outer dimensions (`mm`, `in`)."),
			"mounting_depth": nbschema.DSComputedStringAttribute("Maximum depth of equipment that can be installed (in mm)."),
			"airflow":        nbschema.DSComputedStringAttribute("Direction of airflow through the rack (`front-to-rear`, `rear-to-front`, `passive`, `mixed`)."),

			// Text fields
			"description": nbschema.DSComputedStringAttribute("Description of the rack."),
			"comments":    nbschema.DSComputedStringAttribute("Additional comments or notes about the rack."),

			// Tags and custom fields
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RackDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves data from the Netbox API.
func (d *RackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RackDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rack *netbox.Rack
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID or name
	switch {
	case !data.ID.IsNull():
		rackID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading rack by ID", map[string]interface{}{
			"id": rackID,
		})

		var rackIDInt int32
		if _, parseErr := fmt.Sscanf(rackID, "%d", &rackIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Rack ID",
				fmt.Sprintf("Rack ID must be a number, got: %s", rackID),
			)
			return
		}

		rack, httpResp, err = d.client.DcimAPI.DcimRacksRetrieve(ctx, rackIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Name.IsNull():
		rackName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading rack by name", map[string]interface{}{
			"name": rackName,
		})

		var racks *netbox.PaginatedRackList
		racks, httpResp, err = d.client.DcimAPI.DcimRacksList(ctx).Name([]string{rackName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading rack",
				utils.FormatAPIError("read rack by name", err, httpResp),
			)
			return
		}
		if len(racks.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Rack Not Found",
				fmt.Sprintf("No rack found with name: %s", rackName),
			)
			return
		}
		if len(racks.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Racks Found",
				fmt.Sprintf("Multiple racks found with name: %s. Rack names may not be unique across sites in Netbox. Consider using the rack ID instead.", rackName),
			)
			return
		}
		rack = &racks.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Rack Identifier",
			"Either 'id' or 'name' must be specified to identify the rack.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading rack",
			utils.FormatAPIError("read rack", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading rack",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	mapRackDataSourceToState(ctx, rack, &data)

	tflog.Debug(ctx, "Read rack", map[string]interface{}{
		"id":   rack.GetId(),
		"name": rack.GetName(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapRackDataSourceToState maps a Netbox Rack to the data source model.
func mapRackDataSourceToState(ctx context.Context, rack *netbox.Rack, data *RackDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rack.GetId()))
	data.Name = types.StringValue(rack.GetName())

	// Map site (required field)
	site := rack.GetSite()
	data.Site = types.StringValue(site.GetName())
	data.SiteID = types.StringValue(fmt.Sprintf("%d", site.GetId()))

	// Map location (optional)
	if rack.HasLocation() && rack.GetLocation().Id != 0 {
		location := rack.GetLocation()
		data.Location = types.StringValue(location.GetName())
		data.LocationID = types.StringValue(fmt.Sprintf("%d", location.GetId()))
	} else {
		data.Location = types.StringNull()
		data.LocationID = types.StringNull()
	}

	// Map tenant (optional)
	if rack.HasTenant() && rack.GetTenant().Id != 0 {
		tenant := rack.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Map status
	if rack.HasStatus() {
		status := rack.GetStatus()
		if value, ok := status.GetValueOk(); ok && value != nil {
			data.Status = types.StringValue(string(*value))
		} else {
			data.Status = types.StringNull()
		}
	} else {
		data.Status = types.StringNull()
	}

	// Map role (optional)
	if rack.HasRole() && rack.GetRole().Id != 0 {
		role := rack.GetRole()
		data.Role = types.StringValue(role.GetName())
		data.RoleID = types.StringValue(fmt.Sprintf("%d", role.GetId()))
	} else {
		data.Role = types.StringNull()
		data.RoleID = types.StringNull()
	}

	// Map serial
	if serial := rack.GetSerial(); serial != "" {
		data.Serial = types.StringValue(serial)
	} else {
		data.Serial = types.StringNull()
	}

	// Map asset_tag
	if assetTag, ok := rack.GetAssetTagOk(); ok && assetTag != nil && *assetTag != "" {
		data.AssetTag = types.StringValue(*assetTag)
	} else {
		data.AssetTag = types.StringNull()
	}

	// Map rack_type (optional)
	if rack.HasRackType() && rack.GetRackType().Id != 0 {
		rackType := rack.GetRackType()
		data.RackType = types.StringValue(rackType.GetModel())
		data.RackTypeID = types.StringValue(fmt.Sprintf("%d", rackType.GetId()))
	} else {
		data.RackType = types.StringNull()
		data.RackTypeID = types.StringNull()
	}

	// Map form_factor
	if rack.HasFormFactor() {
		formFactor := rack.GetFormFactor()
		if value, ok := formFactor.GetValueOk(); ok && value != nil {
			data.FormFactor = types.StringValue(string(*value))
		} else {
			data.FormFactor = types.StringNull()
		}
	} else {
		data.FormFactor = types.StringNull()
	}

	// Map width
	if rack.HasWidth() {
		width := rack.GetWidth()
		if value, ok := width.GetValueOk(); ok && value != nil {
			data.Width = types.StringValue(fmt.Sprintf("%d", *value))
		} else {
			data.Width = types.StringNull()
		}
	} else {
		data.Width = types.StringNull()
	}

	// Map u_height
	if uHeight, ok := rack.GetUHeightOk(); ok && uHeight != nil {
		data.UHeight = types.StringValue(fmt.Sprintf("%d", *uHeight))
	} else {
		data.UHeight = types.StringNull()
	}

	// Map starting_unit
	if startingUnit, ok := rack.GetStartingUnitOk(); ok && startingUnit != nil {
		data.StartingUnit = types.StringValue(fmt.Sprintf("%d", *startingUnit))
	} else {
		data.StartingUnit = types.StringNull()
	}

	// Map weight
	if weight, ok := rack.GetWeightOk(); ok && weight != nil {
		data.Weight = types.StringValue(fmt.Sprintf("%g", *weight))
	} else {
		data.Weight = types.StringNull()
	}

	// Map max_weight
	if maxWeight, ok := rack.GetMaxWeightOk(); ok && maxWeight != nil {
		data.MaxWeight = types.StringValue(fmt.Sprintf("%d", *maxWeight))
	} else {
		data.MaxWeight = types.StringNull()
	}

	// Map weight_unit
	if rack.HasWeightUnit() {
		weightUnit := rack.GetWeightUnit()
		if value, ok := weightUnit.GetValueOk(); ok && value != nil {
			data.WeightUnit = types.StringValue(string(*value))
		} else {
			data.WeightUnit = types.StringNull()
		}
	} else {
		data.WeightUnit = types.StringNull()
	}

	// Map desc_units
	if descUnits, ok := rack.GetDescUnitsOk(); ok && descUnits != nil {
		data.DescUnits = types.BoolValue(*descUnits)
	} else {
		data.DescUnits = types.BoolNull()
	}

	// Map outer_width
	if outerWidth, ok := rack.GetOuterWidthOk(); ok && outerWidth != nil {
		data.OuterWidth = types.StringValue(fmt.Sprintf("%d", *outerWidth))
	} else {
		data.OuterWidth = types.StringNull()
	}

	// Map outer_depth
	if outerDepth, ok := rack.GetOuterDepthOk(); ok && outerDepth != nil {
		data.OuterDepth = types.StringValue(fmt.Sprintf("%d", *outerDepth))
	} else {
		data.OuterDepth = types.StringNull()
	}

	// Map outer_unit
	if rack.HasOuterUnit() {
		outerUnit := rack.GetOuterUnit()
		if value, ok := outerUnit.GetValueOk(); ok && value != nil {
			data.OuterUnit = types.StringValue(string(*value))
		} else {
			data.OuterUnit = types.StringNull()
		}
	} else {
		data.OuterUnit = types.StringNull()
	}

	// Map mounting_depth
	if mountingDepth, ok := rack.GetMountingDepthOk(); ok && mountingDepth != nil {
		data.MountingDepth = types.StringValue(fmt.Sprintf("%d", *mountingDepth))
	} else {
		data.MountingDepth = types.StringNull()
	}

	// Map airflow
	if rack.HasAirflow() {
		airflow := rack.GetAirflow()
		if value, ok := airflow.GetValueOk(); ok && value != nil {
			data.Airflow = types.StringValue(string(*value))
		} else {
			data.Airflow = types.StringNull()
		}
	} else {
		data.Airflow = types.StringNull()
	}

	// Map description
	if description := rack.GetDescription(); description != "" {
		data.Description = types.StringValue(description)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments := rack.GetComments(); comments != "" {
		data.Comments = types.StringValue(comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if rack.HasTags() {
		tags := utils.NestedTagsToTagModels(rack.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		if !diags.HasError() {
			data.Tags = tagsValue
		} else {
			data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if rack.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(rack.GetCustomFields(), []utils.CustomFieldModel{})
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !diags.HasError() {
			data.CustomFields = customFieldsValue
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
