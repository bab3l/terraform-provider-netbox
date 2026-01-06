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
var _ datasource.DataSource = &LocationDataSource{}

func NewLocationDataSource() datasource.DataSource {
	return &LocationDataSource{}
}

// LocationDataSource defines the data source implementation.
type LocationDataSource struct {
	client *netbox.APIClient
}

// LocationDataSourceModel describes the data source data model.
type LocationDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	DisplayName  types.String `tfsdk:"display_name"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Site         types.String `tfsdk:"site"`
	SiteID       types.String `tfsdk:"site_id"`
	Parent       types.String `tfsdk:"parent"`
	ParentID     types.String `tfsdk:"parent_id"`
	Status       types.String `tfsdk:"status"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantID     types.String `tfsdk:"tenant_id"`
	Facility     types.String `tfsdk:"facility"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *LocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (d *LocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a location in Netbox. Locations represent physical areas within a site, such as buildings, floors, or rooms. Locations can be nested hierarchically to model complex site layouts. You can identify the location using `id`, `slug`, or `name`.",
		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("location"),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the location."),
			"name":          nbschema.DSNameAttribute("location"),
			"slug":          nbschema.DSSlugAttribute("location"),
			"site":          nbschema.DSComputedStringAttribute("Name of the site where this location resides."),
			"site_id":       nbschema.DSComputedStringAttribute("ID of the site where this location resides."),
			"parent":        nbschema.DSComputedStringAttribute("Name of the parent location."),
			"parent_id":     nbschema.DSComputedStringAttribute("ID of the parent location."),
			"status":        nbschema.DSComputedStringAttribute("Operational status of the location (e.g., `planned`, `staging`, `active`, `decommissioning`, `retired`)."),
			"tenant":        nbschema.DSComputedStringAttribute("Name of the tenant that owns this location."),
			"tenant_id":     nbschema.DSComputedStringAttribute("ID of the tenant that owns this location."),
			"facility":      nbschema.DSComputedStringAttribute("Local facility identifier or description."),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the location."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LocationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *LocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var location *netbox.Location
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	switch {
	case !data.ID.IsNull():
		locationID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading location by ID", map[string]interface{}{
			"id": locationID,
		})
		var locationIDInt int32
		if _, parseErr := fmt.Sscanf(locationID, "%d", &locationIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Location ID",
				fmt.Sprintf("Location ID must be a number, got: %s", locationID),
			)
			return
		}
		location, httpResp, err = d.client.DcimAPI.DcimLocationsRetrieve(ctx, locationIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)

	case !data.Slug.IsNull():
		locationSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading location by slug", map[string]interface{}{
			"slug": locationSlug,
		})
		var locations *netbox.PaginatedLocationList
		locations, httpResp, err = d.client.DcimAPI.DcimLocationsList(ctx).Slug([]string{locationSlug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading location",
				utils.FormatAPIError("read location by slug", err, httpResp),
			)
			return
		}
		if len(locations.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Location Not Found",
				fmt.Sprintf("No location found with slug: %s", locationSlug),
			)
			return
		}
		if len(locations.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Locations Found",
				fmt.Sprintf("Multiple locations found with slug: %s. This should not happen as slugs should be unique.", locationSlug),
			)
			return
		}
		location = &locations.GetResults()[0]

	case !data.Name.IsNull():
		locationName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading location by name", map[string]interface{}{
			"name": locationName,
		})
		var locations *netbox.PaginatedLocationList
		locations, httpResp, err = d.client.DcimAPI.DcimLocationsList(ctx).Name([]string{locationName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading location",
				utils.FormatAPIError("read location by name", err, httpResp),
			)
			return
		}
		if len(locations.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Location Not Found",
				fmt.Sprintf("No location found with name: %s", locationName),
			)
			return
		}
		if len(locations.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Locations Found",
				fmt.Sprintf("Multiple locations found with name: %s. Location names may not be unique in Netbox.", locationName),
			)
			return
		}
		location = &locations.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Location Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the location.",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading location",
			utils.FormatAPIError("read location", err, httpResp),
		)
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error reading location",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", location.GetId()))

	// Display Name
	if location.GetDisplay() != "" {
		data.DisplayName = types.StringValue(location.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
	data.Name = types.StringValue(location.GetName())
	data.Slug = types.StringValue(location.GetSlug())

	// Map site (required field)
	site := location.GetSite()
	data.Site = types.StringValue(site.GetName())
	data.SiteID = types.StringValue(fmt.Sprintf("%d", site.GetId()))

	// Map parent (optional, hierarchical)
	if location.HasParent() && location.GetParent().Id != 0 {
		parent := location.GetParent()
		data.Parent = types.StringValue(parent.GetName())
		data.ParentID = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
		data.ParentID = types.StringNull()
	}

	// Map status
	if location.HasStatus() {
		status := location.GetStatus()
		data.Status = types.StringValue(string(status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Map tenant
	if location.HasTenant() && location.GetTenant().Id != 0 {
		tenant := location.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Map facility
	if location.HasFacility() && location.GetFacility() != "" {
		data.Facility = types.StringValue(location.GetFacility())
	} else {
		data.Facility = types.StringNull()
	}

	// Map description
	if location.HasDescription() {
		data.Description = types.StringValue(location.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if location.HasTags() {
		tags := utils.NestedTagsToTagModels(location.GetTags())
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
	if location.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(location.GetCustomFields(), []utils.CustomFieldModel{})
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
	tflog.Debug(ctx, "Read location", map[string]interface{}{
		"id":   location.GetId(),
		"name": location.GetName(),
		"slug": location.GetSlug(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
