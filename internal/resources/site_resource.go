// Package resources contains Terraform resource implementations for the Netbox provider.

//

// This package integrates with the go-netbox OpenAPI client to provide

// CRUD operations for Netbox resources via Terraform.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Status types.String `tfsdk:"status"`

	Region types.String `tfsdk:"region"`

	RegionID types.String `tfsdk:"region_id"`

	Group types.String `tfsdk:"group"`

	GroupID types.String `tfsdk:"group_id"`

	Tenant types.String `tfsdk:"tenant"`

	TenantID types.String `tfsdk:"tenant_id"`

	Facility types.String `tfsdk:"facility"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_site"

}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a site in Netbox. Sites represent physical locations such as data centers, offices, or other facilities where network infrastructure is deployed.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("site"),

			"name": nbschema.NameAttribute("site", 100),

			"slug": nbschema.SlugAttribute("site"),

			"status": nbschema.StatusAttribute(

				[]string{"planned", "staging", "active", "decommissioning", "retired"},

				"Operational status of the site.",
			),

			"region": nbschema.ReferenceAttribute("region", "ID or slug of the region where this site is located."),

			"region_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the region.",
			},

			"group": nbschema.ReferenceAttribute("site group", "ID or slug of the site group."),

			"group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the site group.",
			},

			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this site."),

			"tenant_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the tenant.",
			},

			"facility": schema.StringAttribute{

				MarkdownDescription: "Local facility identifier or description (e.g., building name, floor, room number).",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.LengthAtMost(50),
				},
			},

			"description": nbschema.DescriptionAttribute("site"),

			"comments": nbschema.CommentsAttributeWithLimit("site", 1000),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
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

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating site", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Prepare the site request

	siteRequest := netbox.WritableSiteRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string fields

	siteRequest.Description = utils.StringPtr(data.Description)

	siteRequest.Comments = utils.StringPtr(data.Comments)

	siteRequest.Facility = utils.StringPtr(data.Facility)

	// Set status if provided

	if utils.IsSet(data.Status) {

		statusValue := netbox.LocationStatusValue(data.Status.ValueString())

		siteRequest.Status = &statusValue

	}

	// Handle tenant relationship

	if utils.IsSet(data.Tenant) {

		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)

	}

	// Handle region relationship

	if utils.IsSet(data.Region) {

		regionRef, diags := netboxlookup.LookupRegion(ctx, r.client, data.Region.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Region = *netbox.NewNullableBriefRegionRequest(regionRef)

	}

	// Handle group relationship

	if utils.IsSet(data.Group) {

		groupRef, diags := netboxlookup.LookupSiteGroup(ctx, r.client, data.Group.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Group = *netbox.NewNullableBriefSiteGroupRequest(groupRef)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Create the site via API

	site, httpResp, err := r.client.DcimAPI.DcimSitesCreate(ctx).WritableSiteRequest(siteRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		handler := utils.CreateErrorHandler{

			ResourceType: "netbox_site",

			ResourceName: "this.site",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {

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

		resp.Diagnostics.AddError("Error creating site", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))

		return

	}

	if site == nil {

		resp.Diagnostics.AddError("Site API returned nil", "No site object returned from Netbox API.")

		return

	}

	// Map response to state using helper

	r.mapSiteToState(ctx, site, &data)

	tflog.Debug(ctx, "Created site", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data SiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	siteID := data.ID.ValueString()

	var siteIDInt int32

	siteIDInt, err := utils.ParseID(siteID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", siteID))

		return

	}

	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError(fmt.Sprintf("read site ID %s", siteID), err, httpResp))

		return

	}

	if httpResp.StatusCode == 404 {

		resp.State.RemoveResource(ctx)

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error reading site", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state using helper

	r.mapSiteToState(ctx, site, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data SiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	siteID := data.ID.ValueString()

	var siteIDInt int32

	siteIDInt, err := utils.ParseID(siteID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", siteID))

		return

	}

	// Prepare the site update request

	siteRequest := netbox.WritableSiteRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string fields

	siteRequest.Description = utils.StringPtr(data.Description)

	siteRequest.Comments = utils.StringPtr(data.Comments)

	siteRequest.Facility = utils.StringPtr(data.Facility)

	// Set status if provided

	if utils.IsSet(data.Status) {

		statusValue := netbox.LocationStatusValue(data.Status.ValueString())

		siteRequest.Status = &statusValue

	}

	// Handle tenant relationship

	if utils.IsSet(data.Tenant) {

		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)

	}

	// Handle region relationship

	if utils.IsSet(data.Region) {

		regionRef, diags := netboxlookup.LookupRegion(ctx, r.client, data.Region.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Region = *netbox.NewNullableBriefRegionRequest(regionRef)

	}

	// Handle group relationship

	if utils.IsSet(data.Group) {

		groupRef, diags := netboxlookup.LookupSiteGroup(ctx, r.client, data.Group.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Group = *netbox.NewNullableBriefSiteGroupRequest(groupRef)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Update the site via API

	site, httpResp, err := r.client.DcimAPI.DcimSitesUpdate(ctx, siteIDInt).WritableSiteRequest(siteRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error updating site", utils.FormatAPIError(fmt.Sprintf("update site ID %s", siteID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error updating site", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state using helper

	r.mapSiteToState(ctx, site, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data SiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	siteID := data.ID.ValueString()

	var siteIDInt int32

	siteIDInt, err := utils.ParseID(siteID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", siteID))

		return

	}

	httpResp, err := r.client.DcimAPI.DcimSitesDestroy(ctx, siteIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error deleting site", utils.FormatAPIError(fmt.Sprintf("delete site ID %s", siteID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 204 {

		resp.Diagnostics.AddError("Error deleting site", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return

	}

}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapSiteToState maps API response to Terraform state using state helpers.

func (r *SiteResource) mapSiteToState(ctx context.Context, site *netbox.Site, data *SiteResourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", site.GetId()))

	data.Name = types.StringValue(site.GetName())

	data.Slug = types.StringValue(site.GetSlug())

	// Handle status

	if site.HasStatus() {

		status := site.GetStatus()

		if status.HasValue() {

			statusValue, _ := status.GetValueOk()

			data.Status = types.StringValue(string(*statusValue))

		}

	} else {

		data.Status = types.StringValue("active")

	}

	// Handle tenant reference

	if site.HasTenant() {

		tenant := site.GetTenant()

		if tenant.Id != 0 {

			data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))

			userTenant := data.Tenant.ValueString()

			if userTenant == tenant.GetName() || userTenant == tenant.GetSlug() || userTenant == tenant.GetDisplay() || userTenant == fmt.Sprintf("%d", tenant.GetId()) {

				// Keep user's original value

			} else {

				data.Tenant = types.StringValue(tenant.GetName())

			}

		} else {

			data.Tenant = types.StringNull()
			data.TenantID = types.StringNull()

		}

	} else {

		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()

	}

	// Handle region reference

	if site.HasRegion() {

		region := site.GetRegion()

		if region.Id != 0 {

			data.RegionID = types.StringValue(fmt.Sprintf("%d", region.GetId()))

			userRegion := data.Region.ValueString()

			if userRegion == region.GetName() || userRegion == region.GetSlug() || userRegion == region.GetDisplay() || userRegion == fmt.Sprintf("%d", region.GetId()) {

				// Keep user's original value

			} else {

				data.Region = types.StringValue(region.GetName())

			}

		} else {

			data.Region = types.StringNull()
			data.RegionID = types.StringNull()

		}

	} else {

		data.Region = types.StringNull()
		data.RegionID = types.StringNull()

	}

	// Handle group reference

	if site.HasGroup() {

		group := site.GetGroup()

		if group.Id != 0 {

			data.GroupID = types.StringValue(fmt.Sprintf("%d", group.GetId()))

			userGroup := data.Group.ValueString()

			if userGroup == group.GetName() || userGroup == group.GetSlug() || userGroup == group.GetDisplay() || userGroup == fmt.Sprintf("%d", group.GetId()) {

				// Keep user's original value

			} else {

				data.Group = types.StringValue(group.GetName())

			}

		} else {

			data.Group = types.StringNull()
			data.GroupID = types.StringNull()

		}

	} else {

		data.Group = types.StringNull()
		data.GroupID = types.StringNull()

	}

	// Handle optional string fields using helpers

	data.Facility = utils.StringFromAPI(site.HasFacility(), site.GetFacility, data.Facility)

	data.Description = utils.StringFromAPI(site.HasDescription(), site.GetDescription, data.Description)

	data.Comments = utils.StringFromAPI(site.HasComments(), site.GetComments, data.Comments)

	// Handle tags

	if site.HasTags() {

		tags := utils.NestedTagsToTagModels(site.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		if !tagDiags.HasError() {

			data.Tags = tagsValue

		}

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields - preserve state structure

	if site.HasCustomFields() && !data.CustomFields.IsNull() {

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		if !cfDiags.HasError() {

			customFields := utils.MapToCustomFieldModels(site.GetCustomFields(), stateCustomFields)

			customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

			if !cfValueDiags.HasError() {

				data.CustomFields = customFieldsValue

			}

		}

	} else if data.CustomFields.IsNull() {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
