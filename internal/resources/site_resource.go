// Package resources contains Terraform resource implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// CRUD operations for Netbox resources via Terraform.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
var _ resource.ResourceWithIdentity = &SiteResource{}

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

// SiteResource defines the resource implementation.
type SiteResource struct {
	client *netbox.APIClient
}

// SiteResourceModel describes the resource data model.
type SiteResourceModel struct {
	ID              types.String  `tfsdk:"id"`
	Name            types.String  `tfsdk:"name"`
	Slug            types.String  `tfsdk:"slug"`
	Status          types.String  `tfsdk:"status"`
	Region          types.String  `tfsdk:"region"`
	RegionID        types.String  `tfsdk:"region_id"`
	Group           types.String  `tfsdk:"group"`
	GroupID         types.String  `tfsdk:"group_id"`
	Tenant          types.String  `tfsdk:"tenant"`
	TenantID        types.String  `tfsdk:"tenant_id"`
	Facility        types.String  `tfsdk:"facility"`
	TimeZone        types.String  `tfsdk:"time_zone"`
	PhysicalAddress types.String  `tfsdk:"physical_address"`
	ShippingAddress types.String  `tfsdk:"shipping_address"`
	Latitude        types.Float64 `tfsdk:"latitude"`
	Longitude       types.Float64 `tfsdk:"longitude"`
	Description     types.String  `tfsdk:"description"`
	Comments        types.String  `tfsdk:"comments"`
	Tags            types.Set     `tfsdk:"tags"`
	CustomFields    types.Set     `tfsdk:"custom_fields"`
}

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a site in Netbox. Sites represent physical locations such as data centers, offices, or other facilities where network infrastructure is deployed.",
		Attributes: map[string]schema.Attribute{
			"id":   nbschema.IDAttribute("site"),
			"name": nbschema.NameAttribute("site", 100),
			"slug": nbschema.SlugAttribute("site"),
			"status": nbschema.StatusAttribute(
				[]string{"planned", "staging", "active", "decommissioning", "retired"},
				"Operational status of the site.",
			),
			"region": nbschema.ReferenceAttributeWithDiffSuppress("region", "ID or slug of the region where this site is located."),
			"region_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the region.",
			},
			"group": nbschema.ReferenceAttributeWithDiffSuppress("site group", "ID or slug of the site group."),
			"group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the site group.",
			},
			"tenant": nbschema.ReferenceAttributeWithDiffSuppress("tenant", "ID or slug of the tenant that owns this site."),
			"tenant_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the tenant.",
			},
			"facility": schema.StringAttribute{
				MarkdownDescription: "Local facility identifier or description (e.g., building name, floor, room number).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"time_zone": schema.StringAttribute{
				MarkdownDescription: "Time zone for this site (IANA name, e.g. 'America/Los_Angeles').",
				Optional:            true,
			},
			"physical_address": schema.StringAttribute{
				MarkdownDescription: "Physical address of the site.",
				Optional:            true,
			},
			"shipping_address": schema.StringAttribute{
				MarkdownDescription: "Shipping address for the site (if different from physical address).",
				Optional:            true,
			},
			"latitude": schema.Float64Attribute{
				MarkdownDescription: "GPS latitude coordinate in decimal format (xx.yyyyyy).",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(-90, 90),
				},
			},
			"longitude": schema.Float64Attribute{
				MarkdownDescription: "GPS longitude coordinate in decimal format (xx.yyyyyy).",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(-180, 180),
				},
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("site"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *SiteResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
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
	siteRequest.Facility = utils.StringPtr(data.Facility)
	siteRequest.PhysicalAddress = utils.StringPtr(data.PhysicalAddress)
	siteRequest.ShippingAddress = utils.StringPtr(data.ShippingAddress)
	if utils.IsSet(data.TimeZone) {
		siteRequest.SetTimeZone(data.TimeZone.ValueString())
	}
	if utils.IsSet(data.Latitude) {
		latitude := data.Latitude.ValueFloat64()
		siteRequest.SetLatitude(latitude)
	}
	if utils.IsSet(data.Longitude) {
		longitude := data.Longitude.ValueFloat64()
		siteRequest.SetLongitude(longitude)
	}

	// Set status if provided
	if utils.IsSet(data.Status) {
		statusValue := netbox.LocationStatusValue(data.Status.ValueString())
		siteRequest.Status = &statusValue
	}

	// Handle tenant relationship
	if tenantRef := utils.ResolveOptionalReference(ctx, r.client, data.Tenant, netboxlookup.LookupTenant, &resp.Diagnostics); tenantRef != nil {
		siteRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	} else if data.Tenant.IsNull() {
		siteRequest.SetTenantNil()
	}

	// Handle region relationship
	if regionRef := utils.ResolveOptionalReference(ctx, r.client, data.Region, netboxlookup.LookupRegion, &resp.Diagnostics); regionRef != nil {
		siteRequest.Region = *netbox.NewNullableBriefRegionRequest(regionRef)
	} else if data.Region.IsNull() {
		siteRequest.SetRegionNil()
	}

	// Handle group relationship
	if groupRef := utils.ResolveOptionalReference(ctx, r.client, data.Group, netboxlookup.LookupSiteGroup, &resp.Diagnostics); groupRef != nil {
		siteRequest.Group = *netbox.NewNullableBriefSiteGroupRequest(groupRef)
	} else if data.Group.IsNull() {
		siteRequest.SetGroupNil()
	}

	// Apply description and comments
	utils.ApplyDescriptiveFields(&siteRequest, data.Description, data.Comments)

	// Apply tags from slugs
	utils.ApplyTagsFromSlugs(ctx, r.client, &siteRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply custom fields
	utils.ApplyCustomFields(ctx, &siteRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the site via API
	site, httpResp, err := r.client.DcimAPI.DcimSitesCreate(ctx).WritableSiteRequest(siteRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_site",
			ResourceName: "this.site",
			SlugValue:    data.Slug.ValueString(),
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
	r.mapSiteToState(ctx, site, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created site", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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

	if httpResp != nil && httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError(fmt.Sprintf("read site ID %s", siteID), err, httpResp))
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading site", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Map response to state using helper
	r.mapSiteToState(ctx, site, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan SiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := plan.ID.ValueString()
	var siteIDInt int32
	siteIDInt, err := utils.ParseID(siteID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", siteID))
		return
	}

	// Prepare the site update request
	siteRequest := netbox.WritableSiteRequest{
		Name: plan.Name.ValueString(),
		Slug: plan.Slug.ValueString(),
	}

	// Handle optional string fields
	// NetBox uses PATCH semantics: omitting a field does not clear it.
	// When the user removes an optional string from config, send an explicit empty string to clear it.
	if plan.Facility.IsNull() {
		empty := ""
		siteRequest.Facility = &empty
	} else {
		siteRequest.Facility = utils.StringPtr(plan.Facility)
	}
	if plan.PhysicalAddress.IsNull() {
		empty := ""
		siteRequest.PhysicalAddress = &empty
	} else {
		siteRequest.PhysicalAddress = utils.StringPtr(plan.PhysicalAddress)
	}
	if plan.ShippingAddress.IsNull() {
		empty := ""
		siteRequest.ShippingAddress = &empty
	} else {
		siteRequest.ShippingAddress = utils.StringPtr(plan.ShippingAddress)
	}
	if utils.IsSet(plan.TimeZone) {
		siteRequest.SetTimeZone(plan.TimeZone.ValueString())
	} else if plan.TimeZone.IsNull() {
		siteRequest.SetTimeZoneNil()
	}
	if utils.IsSet(plan.Latitude) {
		latitude := plan.Latitude.ValueFloat64()
		siteRequest.SetLatitude(latitude)
	} else if plan.Latitude.IsNull() {
		siteRequest.SetLatitudeNil()
	}
	if utils.IsSet(plan.Longitude) {
		longitude := plan.Longitude.ValueFloat64()
		siteRequest.SetLongitude(longitude)
	} else if plan.Longitude.IsNull() {
		siteRequest.SetLongitudeNil()
	}

	// Set status if provided
	if utils.IsSet(plan.Status) {
		statusValue := netbox.LocationStatusValue(plan.Status.ValueString())
		siteRequest.Status = &statusValue
	}

	// Handle tenant relationship
	if tenantRef := utils.ResolveOptionalReference(ctx, r.client, plan.Tenant, netboxlookup.LookupTenant, &resp.Diagnostics); tenantRef != nil {
		siteRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	} else if plan.Tenant.IsNull() {
		siteRequest.SetTenantNil()
	}

	// Handle region relationship
	if regionRef := utils.ResolveOptionalReference(ctx, r.client, plan.Region, netboxlookup.LookupRegion, &resp.Diagnostics); regionRef != nil {
		siteRequest.Region = *netbox.NewNullableBriefRegionRequest(regionRef)
	} else if plan.Region.IsNull() {
		siteRequest.SetRegionNil()
	}

	// Handle group relationship
	if groupRef := utils.ResolveOptionalReference(ctx, r.client, plan.Group, netboxlookup.LookupSiteGroup, &resp.Diagnostics); groupRef != nil {
		siteRequest.Group = *netbox.NewNullableBriefSiteGroupRequest(groupRef)
	} else if plan.Group.IsNull() {
		siteRequest.SetGroupNil()
	}

	// Apply description and comments
	utils.ApplyDescriptiveFields(&siteRequest, plan.Description, plan.Comments)

	// Handle tags and custom fields - merge-aware for partial management
	// If tags are in plan, use plan. If not, preserve state tags.
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &siteRequest, plan.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &siteRequest, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply custom fields with merge logic (preserves unmanaged fields from state)
	utils.ApplyCustomFieldsWithMerge(ctx, &siteRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
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
	r.mapSiteToState(ctx, site, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
		// Ignore 404 errors (resource already deleted)
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Site already deleted", map[string]interface{}{
				"id": siteID,
			})
			return
		}
		resp.Diagnostics.AddError("Error deleting site", utils.FormatAPIError(fmt.Sprintf("delete site ID %s", siteID), err, httpResp))
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error deleting site", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		siteIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", parsed.ID))
			return
		}
		site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing site", utils.FormatAPIError(fmt.Sprintf("read site ID %s", parsed.ID), err, httpResp))
			return
		}

		var data SiteResourceModel
		if site.HasRegion() && site.GetRegion().Id != 0 {
			region := site.GetRegion()
			data.Region = types.StringValue(fmt.Sprintf("%d", region.GetId()))
		}
		if site.HasGroup() && site.GetGroup().Id != 0 {
			group := site.GetGroup()
			data.Group = types.StringValue(fmt.Sprintf("%d", group.GetId()))
		}
		if site.HasTenant() && site.GetTenant().Id != 0 {
			tenant := site.GetTenant()
			data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		}
		if site.HasTags() {
			tagSlugs := make([]string, 0, len(site.GetTags()))
			for _, tag := range site.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		} else {
			data.Tags = types.SetNull(types.StringType)
		}
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapSiteToState(ctx, site, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, site.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapSiteToState maps API response to Terraform state using state helpers.
func (r *SiteResource) mapSiteToState(ctx context.Context, site *netbox.Site, data *SiteResourceModel, diags *diag.Diagnostics) {
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
	var tenantResult utils.ReferenceWithID
	if site.HasTenant() {
		tenant := site.GetTenant()
		tenantResult = utils.PreserveOptionalReferenceWithID(data.Tenant, tenant.Id != 0, tenant.GetId(), tenant.GetName(), tenant.GetSlug())
	} else {
		tenantResult = utils.PreserveOptionalReferenceWithID(data.Tenant, false, 0, "", "")
	}
	data.Tenant = tenantResult.Reference
	data.TenantID = tenantResult.ID

	// Handle region reference
	var regionResult utils.ReferenceWithID
	if site.HasRegion() {
		region := site.GetRegion()
		regionResult = utils.PreserveOptionalReferenceWithID(data.Region, region.Id != 0, region.GetId(), region.GetName(), region.GetSlug())
	} else {
		regionResult = utils.PreserveOptionalReferenceWithID(data.Region, false, 0, "", "")
	}
	data.Region = regionResult.Reference
	data.RegionID = regionResult.ID

	// Handle group reference
	var groupResult utils.ReferenceWithID
	if site.HasGroup() {
		group := site.GetGroup()
		groupResult = utils.PreserveOptionalReferenceWithID(data.Group, group.Id != 0, group.GetId(), group.GetName(), group.GetSlug())
	} else {
		groupResult = utils.PreserveOptionalReferenceWithID(data.Group, false, 0, "", "")
	}
	data.Group = groupResult.Reference
	data.GroupID = groupResult.ID

	// Handle optional string fields using helpers
	data.Facility = utils.StringFromAPI(site.HasFacility(), site.GetFacility, data.Facility)
	data.Description = utils.StringFromAPI(site.HasDescription(), site.GetDescription, data.Description)
	data.Comments = utils.StringFromAPI(site.HasComments(), site.GetComments, data.Comments)
	data.TimeZone = utils.StringFromAPI(site.HasTimeZone(), site.GetTimeZone, data.TimeZone)
	data.PhysicalAddress = utils.StringFromAPI(site.HasPhysicalAddress(), site.GetPhysicalAddress, data.PhysicalAddress)
	data.ShippingAddress = utils.StringFromAPI(site.HasShippingAddress(), site.GetShippingAddress, data.ShippingAddress)

	// Handle latitude/longitude
	if site.HasLatitude() && site.Latitude.Get() != nil {
		data.Latitude = types.Float64Value(*site.Latitude.Get())
	} else {
		data.Latitude = types.Float64Null()
	}
	if site.HasLongitude() && site.Longitude.Get() != nil {
		data.Longitude = types.Float64Value(*site.Longitude.Get())
	} else {
		data.Longitude = types.Float64Null()
	}

	// Handle tags
	var tagSlugs []string
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	case site.HasTags():
		for _, tag := range site.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	}

	// Handle custom fields - use filtered-to-owned for partial management
	if site.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, site.GetCustomFields(), diags)
	}
}
