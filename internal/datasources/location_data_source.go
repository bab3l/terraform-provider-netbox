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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the location. Specify `id`, `slug`, or `name` to identify the location.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the location. Can be used to identify the location instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the location. Specify `id`, `slug`, or `name` to identify the location.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "Name of the site where this location resides.",
				Computed:            true,
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "ID of the site where this location resides.",
				Computed:            true,
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "Name of the parent location.",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent location.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the location (e.g., `planned`, `staging`, `active`, `decommissioning`, `retired`).",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant that owns this location.",
				Computed:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant that owns this location.",
				Computed:            true,
			},
			"facility": schema.StringAttribute{
				MarkdownDescription: "Local facility identifier or description.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the location.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this location.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the tag.",
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the tag.",
							Computed:            true,
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this location.",
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
	if !data.ID.IsNull() {
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
	} else if !data.Slug.IsNull() {
		locationSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading location by slug", map[string]interface{}{
			"slug": locationSlug,
		})

		var locations *netbox.PaginatedLocationList
		locations, httpResp, err = d.client.DcimAPI.DcimLocationsList(ctx).Slug([]string{locationSlug}).Execute()
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
	} else if !data.Name.IsNull() {
		locationName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading location by name", map[string]interface{}{
			"name": locationName,
		})

		var locations *netbox.PaginatedLocationList
		locations, httpResp, err = d.client.DcimAPI.DcimLocationsList(ctx).Name([]string{locationName}).Execute()
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
	} else {
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

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading location",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", location.GetId()))
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
