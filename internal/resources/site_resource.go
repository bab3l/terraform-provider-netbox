// Package resources contains Terraform resource implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// CRUD operations for Netbox resources via Terraform.
package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SiteResource{}
var _ resource.ResourceWithImportState = &SiteResource{}

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

// SiteResource defines the resource implementation.
type SiteResource struct {
	client *netbox.APIClient
}

// SiteResourceModel describes the resource data model.
type SiteResourceModel struct {
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

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a site in Netbox. Sites represent physical locations such as data centers, offices, or other facilities where network infrastructure is deployed.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the site (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the site. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the site. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the site. Valid values include: `planned`, `staging`, `active`, `decommissioning`, `retired`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"planned",
						"staging",
						"active",
						"decommissioning",
						"retired",
					),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the region where this site is located. Regions help organize sites geographically.",
				Optional:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the site group. Site groups provide an additional level of organization.",
				Optional:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the tenant that owns this site. Used for multi-tenancy scenarios.",
				Optional:            true,
			},
			"facility": schema.StringAttribute{
				MarkdownDescription: "Local facility identifier or description (e.g., building name, floor, room number).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the site, its purpose, or other relevant information.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the site. Supports Markdown formatting.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this site. Tags provide a way to categorize and organize resources.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
							},
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
								validators.ValidSlug(),
							},
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this site. Custom fields allow you to store additional structured data.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 50),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`),
									"must start with a letter and contain only letters, numbers, and underscores",
								),
							},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject).",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"text",
									"longtext",
									"integer",
									"boolean",
									"date",
									"url",
									"json",
									"select",
									"multiselect",
									"object",
									"multiobject",
									"multiple",  // legacy
									"selection", // legacy
								),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(1000),
							},
						},
					},
				},
			},
		},
	}
}

func (r *SiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netbox.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create site using go-netbox client
	tflog.Debug(ctx, "Creating site", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Prepare the site request
	siteRequest := netbox.WritableSiteRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Handle tenant relationship
	if !data.Tenant.IsNull() {
		var tenantIDInt int32
		if _, err := fmt.Sscanf(data.Tenant.ValueString(), "%d", &tenantIDInt); err == nil {
			tenantRef, diags := netboxlookup.LookupTenantBrief(ctx, r.client, fmt.Sprintf("%d", tenantIDInt))
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			siteRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
		}
	}
	// Handle region relationship
	if !data.Region.IsNull() {
		var regionIDInt int32
		if _, err := fmt.Sscanf(data.Region.ValueString(), "%d", &regionIDInt); err == nil {
			regionRef, diags := netboxlookup.LookupRegionBrief(ctx, r.client, fmt.Sprintf("%d", regionIDInt))
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			siteRequest.Region = *netbox.NewNullableBriefRegionRequest(regionRef)
		}
	}
	// Handle group relationship
	if !data.Group.IsNull() {
		var groupIDInt int32
		if _, err := fmt.Sscanf(data.Group.ValueString(), "%d", &groupIDInt); err == nil {
			groupRef, diags := netboxlookup.LookupSiteGroupBrief(ctx, r.client, fmt.Sprintf("%d", groupIDInt))
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			siteRequest.Group = *netbox.NewNullableBriefSiteGroupRequest(groupRef)
		}
	}

	// Set optional fields if provided
	if !data.Status.IsNull() {
		statusValue := netbox.LocationStatusValue(data.Status.ValueString())
		siteRequest.Status = &statusValue
	}
	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		siteRequest.Description = &description
	}
	if !data.Comments.IsNull() {
		comments := data.Comments.ValueString()
		siteRequest.Comments = &comments
	}
	if !data.Facility.IsNull() {
		facility := data.Facility.ValueString()
		siteRequest.Facility = &facility
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Create the site via API
	site, httpResp, err := r.client.DcimAPI.DcimSitesCreate(ctx).WritableSiteRequest(siteRequest).Execute()
	if err != nil {
		// Use enhanced error handler that detects duplicates and provides import hints
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_site",
			ResourceName: "this.site", // Terraform resource name placeholder
			SlugValue:    data.Slug.ValueString(),
			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				// Try to look up existing site by slug
				list, _, lookupErr := r.client.DcimAPI.DcimSitesList(lookupCtx).Slug([]string{slug}).Execute()
				if lookupErr != nil {
					return "", lookupErr
				}
				if list != nil && len(list.Results) > 0 {
					return fmt.Sprintf("%d", list.Results[0].GetId()), nil
				}
				return "", nil
			},
		}
		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating site",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
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
		data.Status = types.StringValue("active") // default status
	}

	if site.HasDescription() {
		desc := site.GetDescription()
		// Keep null if original was null and API returns empty string
		if desc == "" && data.Description.IsNull() {
			// Keep as null
		} else {
			data.Description = types.StringValue(desc)
		}
	}

	if site.HasComments() {
		comments := site.GetComments()
		// Keep null if original was null and API returns empty string
		if comments == "" && data.Comments.IsNull() {
			// Keep as null
		} else {
			data.Comments = types.StringValue(comments)
		}
	}

	if site.HasFacility() {
		facility := site.GetFacility()
		// Keep null if original was null and API returns empty string
		if facility == "" && data.Facility.IsNull() {
			// Keep as null
		} else {
			data.Facility = types.StringValue(facility)
		}
	}

	tflog.Trace(ctx, "created a site resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site ID from state
	siteID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading site", map[string]interface{}{
		"id": siteID,
	})

	// Parse the site ID to int32 for the API call
	var siteIDInt int32
	if _, err := fmt.Sscanf(siteID, "%d", &siteIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site ID",
			fmt.Sprintf("Site ID must be a number, got: %s", siteID),
		)
		return
	}

	// Retrieve the site via API
	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading site",
			utils.FormatAPIError(fmt.Sprintf("read site ID %s", siteID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		// Site no longer exists, remove from state
		resp.State.RemoveResource(ctx)
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
	}

	if site.HasDescription() {
		desc := site.GetDescription()
		// Keep null if original was null and API returns empty string
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	if site.HasComments() {
		comments := site.GetComments()
		// Keep null if original was null and API returns empty string
		if comments == "" && data.Comments.IsNull() {
			data.Comments = types.StringNull()
		} else {
			data.Comments = types.StringValue(comments)
		}
	} else {
		data.Comments = types.StringNull()
	}

	if site.HasFacility() {
		facility := site.GetFacility()
		// Keep null if original was null and API returns empty string
		if facility == "" && data.Facility.IsNull() {
			data.Facility = types.StringNull()
		} else {
			data.Facility = types.StringValue(facility)
		}
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

	// Handle custom fields - we need to preserve the state structure
	if site.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(site.GetCustomFields(), stateCustomFields)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site ID from state
	siteID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating site", map[string]interface{}{
		"id":   siteID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Parse the site ID to int32 for the API call
	var siteIDInt int32
	if _, err := fmt.Sscanf(siteID, "%d", &siteIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site ID",
			fmt.Sprintf("Site ID must be a number, got: %s", siteID),
		)
		return
	}

	// Prepare the site update request
	siteRequest := netbox.WritableSiteRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set status if provided
	if !data.Status.IsNull() {
		statusValue := netbox.LocationStatusValue(data.Status.ValueString())
		siteRequest.Status = &statusValue
	}

	// Handle tenant relationship
	if !data.Tenant.IsNull() {
		var tenantIDInt int32
		if _, err := fmt.Sscanf(data.Tenant.ValueString(), "%d", &tenantIDInt); err == nil {
			tenantRef, diags := netboxlookup.LookupTenantBrief(ctx, r.client, fmt.Sprintf("%d", tenantIDInt))
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			siteRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
		}
	}
	// Handle region relationship
	if !data.Region.IsNull() {
		var regionIDInt int32
		if _, err := fmt.Sscanf(data.Region.ValueString(), "%d", &regionIDInt); err == nil {
			regionRef, diags := netboxlookup.LookupRegionBrief(ctx, r.client, fmt.Sprintf("%d", regionIDInt))
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			siteRequest.Region = *netbox.NewNullableBriefRegionRequest(regionRef)
		}
	}
	// Handle group relationship
	if !data.Group.IsNull() {
		var groupIDInt int32
		if _, err := fmt.Sscanf(data.Group.ValueString(), "%d", &groupIDInt); err == nil {
			groupRef, diags := netboxlookup.LookupSiteGroupBrief(ctx, r.client, fmt.Sprintf("%d", groupIDInt))
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			siteRequest.Group = *netbox.NewNullableBriefSiteGroupRequest(groupRef)
		}
	}

	// Set optional fields if provided
	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		siteRequest.Description = &description
	}
	if !data.Comments.IsNull() {
		comments := data.Comments.ValueString()
		siteRequest.Comments = &comments
	}
	if !data.Facility.IsNull() {
		facility := data.Facility.ValueString()
		siteRequest.Facility = &facility
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Update the site via API
	site, httpResp, err := r.client.DcimAPI.DcimSitesUpdate(ctx, siteIDInt).WritableSiteRequest(siteRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating site",
			utils.FormatAPIError(fmt.Sprintf("update site ID %s", siteID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error updating site",
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
	}

	if site.HasDescription() {
		desc := site.GetDescription()
		if desc == "" && data.Description.IsNull() {
			// Keep null if originally null and API returns empty
		} else {
			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	if site.HasComments() {
		comments := site.GetComments()
		if comments == "" && data.Comments.IsNull() {
			// Keep null if originally null and API returns empty
		} else {
			data.Comments = types.StringValue(comments)
		}
	} else {
		data.Comments = types.StringNull()
	}

	if site.HasFacility() {
		facility := site.GetFacility()
		if facility == "" && data.Facility.IsNull() {
			// Keep null if originally null and API returns empty
		} else {
			data.Facility = types.StringValue(facility)
		}
	} else {
		data.Facility = types.StringNull()
	}

	// Handle tags in response
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

	// Handle custom fields in response
	if site.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(site.GetCustomFields(), stateCustomFields)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site ID from state
	siteID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting site", map[string]interface{}{
		"id": siteID,
	})

	// Parse the site ID to int32 for the API call
	var siteIDInt int32
	if _, err := fmt.Sscanf(siteID, "%d", &siteIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site ID",
			fmt.Sprintf("Site ID must be a number, got: %s", siteID),
		)
		return
	}

	// Delete the site via API
	httpResp, err := r.client.DcimAPI.DcimSitesDestroy(ctx, siteIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting site",
			utils.FormatAPIError(fmt.Sprintf("delete site ID %s", siteID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error deleting site",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted a site resource")
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
