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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SiteDataSource{}

func NewSiteDataSource() datasource.DataSource {
	return &SiteDataSource{}
}

// SiteDataSource defines the data source implementation.
type SiteDataSource struct {
	client *netbox.APIClient
}

// SiteDataSourceModel describes the data source data model.
type SiteDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Status       types.String `tfsdk:"status"`
	Region       types.String `tfsdk:"region"`
	Group        types.String `tfsdk:"group"`
	Tenant       types.String `tfsdk:"tenant"`
	Facility     types.String `tfsdk:"facility"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *SiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (d *SiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a site in Netbox. Sites represent physical locations such as data centers, offices, or other facilities where network infrastructure is deployed. You can identify the site using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the site. Specify `id`, `slug`, or `name` to identify the site.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the site. Can be used to identify the site instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the site. Specify `id`, `slug`, or `name` to identify the site.",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the site (e.g., `planned`, `staging`, `active`, `decommissioning`, `retired`).",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the region where this site is located.",
				Computed:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the site group.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the tenant that owns this site.",
				Computed:            true,
			},
			"facility": schema.StringAttribute{
				MarkdownDescription: "Local facility identifier or description.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the site.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the site.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this site.",
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
				MarkdownDescription: "Custom fields assigned to this site.",
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

func (d *SiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SiteDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var site *netbox.Site
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		// Search by ID
		siteID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading site by ID", map[string]interface{}{
			"id": siteID,
		})

		// Parse the site ID to int32 for the API call
		var siteIDInt int32
		if _, parseErr := fmt.Sscanf(siteID, "%d", &siteIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Site ID",
				fmt.Sprintf("Site ID must be a number, got: %s", siteID),
			)
			return
		}

		// Retrieve the site via API
		site, httpResp, err = d.client.DcimAPI.DcimSitesRetrieve(ctx, siteIDInt).Execute()
	} else if !data.Slug.IsNull() {
		// Search by slug
		siteSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading site by slug", map[string]interface{}{
			"slug": siteSlug,
		})

		// List sites with slug filter
		sites, httpResp, err := d.client.DcimAPI.DcimSitesList(ctx).Slug([]string{siteSlug}).Execute()
		if err == nil && httpResp.StatusCode == 200 {
			if len(sites.GetResults()) == 0 {
				resp.Diagnostics.AddError(
					"Site Not Found",
					fmt.Sprintf("No site found with slug: %s", siteSlug),
				)
				return
			}
			if len(sites.GetResults()) > 1 {
				resp.Diagnostics.AddError(
					"Multiple Sites Found",
					fmt.Sprintf("Multiple sites found with slug: %s. This should not happen as slugs should be unique.", siteSlug),
				)
				return
			}
			site = &sites.GetResults()[0]
		}
	} else if !data.Name.IsNull() {
		// Search by name
		siteName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading site by name", map[string]interface{}{
			"name": siteName,
		})

		// List sites with name filter
		sites, httpResp, err := d.client.DcimAPI.DcimSitesList(ctx).Name([]string{siteName}).Execute()
		if err == nil && httpResp.StatusCode == 200 {
			if len(sites.GetResults()) == 0 {
				resp.Diagnostics.AddError(
					"Site Not Found",
					fmt.Sprintf("No site found with name: %s", siteName),
				)
				return
			}
			if len(sites.GetResults()) > 1 {
				resp.Diagnostics.AddError(
					"Multiple Sites Found",
					fmt.Sprintf("Multiple sites found with name: %s. Site names may not be unique in Netbox.", siteName),
				)
				return
			}
			site = &sites.GetResults()[0]
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing Site Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the site.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading site",
			fmt.Sprintf("Could not read site: %s", err),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError(
			"Site Not Found",
			"The specified site was not found in Netbox.",
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading site",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", site.GetId()))
	data.Name = types.StringValue(site.GetName())
	data.Slug = types.StringValue(site.GetSlug())

	if site.HasStatus() {
		status := site.GetStatus()
		if status.HasValue() {
			statusValue, _ := status.GetValueOk()
			data.Status = types.StringValue(string(*statusValue))
		}
	} else {
		data.Status = types.StringNull()
	}

	// Handle optional relationships
	if site.HasRegion() {
		region := site.GetRegion()
		data.Region = types.StringValue(region.GetName())
	} else {
		data.Region = types.StringNull()
	}

	if site.HasGroup() {
		group := site.GetGroup()
		data.Group = types.StringValue(group.GetName())
	} else {
		data.Group = types.StringNull()
	}

	if site.HasTenant() {
		tenant := site.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
	} else {
		data.Tenant = types.StringNull()
	}

	if site.HasDescription() {
		data.Description = types.StringValue(site.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	if site.HasComments() {
		data.Comments = types.StringValue(site.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	if site.HasFacility() {
		data.Facility = types.StringValue(site.GetFacility())
	} else {
		data.Facility = types.StringNull()
	}

	// Handle tags
	if site.HasTags() {
		tags := utils.NestedTagsToTagModels(site.GetTags())
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
	if site.HasCustomFields() {
		// For data sources, we extract all available custom fields
		customFields := utils.MapToCustomFieldModels(site.GetCustomFields(), nil)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
