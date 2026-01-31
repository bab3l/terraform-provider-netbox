// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &LocationResource{}
	_ resource.ResourceWithImportState = &LocationResource{}
	_ resource.ResourceWithIdentity    = &LocationResource{}
)

func NewLocationResource() resource.Resource {
	return &LocationResource{}
}

// LocationResource defines the resource implementation.
type LocationResource struct {
	client *netbox.APIClient
}

// LocationResourceModel describes the resource data model.
type LocationResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Site         types.String `tfsdk:"site"`
	Parent       types.String `tfsdk:"parent"`
	Status       types.String `tfsdk:"status"`
	Tenant       types.String `tfsdk:"tenant"`
	Facility     types.String `tfsdk:"facility"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *LocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (r *LocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a location in Netbox. Locations represent physical areas within a site, such as buildings, floors, or rooms. Locations can be nested hierarchically.",
		Attributes: map[string]schema.Attribute{
			"id":       nbschema.IDAttribute("location"),
			"name":     nbschema.NameAttribute("location", 100),
			"slug":     nbschema.SlugAttribute("location"),
			"site":     nbschema.RequiredReferenceAttributeWithDiffSuppress("site", "ID or slug of the site this location belongs to. Required."),
			"parent":   nbschema.ReferenceAttributeWithDiffSuppress("parent location", "ID or slug of the parent location. Leave empty for top-level locations within the site."),
			"status":   nbschema.StatusAttribute([]string{"planned", "staging", "active", "decommissioning", "retired"}, "Operational status of the location. Defaults to `active`."),
			"tenant":   nbschema.ReferenceAttributeWithDiffSuppress("tenant", "ID or slug of the tenant that owns this location."),
			"facility": nbschema.FacilityAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("location"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *LocationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *LocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LocationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating location", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
		"site": data.Site.ValueString(),
	})

	// Lookup site
	siteRef, diags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	locationRequest := netbox.NewWritableLocationRequest(data.Name.ValueString(), data.Slug.ValueString(), *siteRef)

	// Set optional parent
	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID, parentDiags := netboxlookup.LookupLocationID(ctx, r.client, data.Parent.ValueString())
		resp.Diagnostics.Append(parentDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		locationRequest.Parent = *netbox.NewNullableInt32(&parentID)
	} else if data.Parent.IsNull() {
		locationRequest.SetParentNil()
	}

	// Set optional status
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.LocationStatusValue(data.Status.ValueString())
		locationRequest.Status = &status
	}

	// Set optional tenant
	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		locationRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	} else if data.Tenant.IsNull() {
		locationRequest.SetTenantNil()
	}

	// Set optional facility
	if !data.Facility.IsNull() && !data.Facility.IsUnknown() {
		facility := data.Facility.ValueString()
		locationRequest.Facility = &facility
	}

	// Set optional description
	utils.ApplyDescription(locationRequest, data.Description)

	// Handle tags
	utils.ApplyTagsFromSlugs(ctx, r.client, locationRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, locationRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	location, httpResp, err := r.client.DcimAPI.DcimLocationsCreate(ctx).WritableLocationRequest(*locationRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating location",
			utils.FormatAPIError("create location", err, httpResp),
		)
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error creating location",
			fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusCreated, httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	r.mapLocationToState(ctx, location, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "created a location resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LocationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	locationID := data.ID.ValueString()
	tflog.Debug(ctx, "Reading location", map[string]interface{}{
		"id": locationID,
	})
	var locationIDInt int32
	locationIDInt, err := utils.ParseID(locationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Location ID",
			fmt.Sprintf("Location ID must be a number, got: %s", locationID),
		)
		return
	}
	location, httpResp, err := r.client.DcimAPI.DcimLocationsRetrieve(ctx, locationIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading location",
			utils.FormatAPIError(fmt.Sprintf("read location ID %s", locationID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error reading location",
			fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusOK, httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	r.mapLocationToState(ctx, location, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan LocationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	locationID := plan.ID.ValueString()
	tflog.Debug(ctx, "Updating location", map[string]interface{}{
		"id": locationID,
	})
	var locationIDInt int32
	locationIDInt, err := utils.ParseID(locationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Location ID",
			fmt.Sprintf("Location ID must be a number, got: %s", locationID),
		)
		return
	}

	// Lookup site
	siteRef, diags := netboxlookup.LookupSite(ctx, r.client, plan.Site.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	locationRequest := netbox.NewWritableLocationRequest(plan.Name.ValueString(), plan.Slug.ValueString(), *siteRef)

	// Set optional parent
	if !plan.Parent.IsNull() && !plan.Parent.IsUnknown() {
		parentID, parentDiags := netboxlookup.LookupLocationID(ctx, r.client, plan.Parent.ValueString())
		resp.Diagnostics.Append(parentDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		locationRequest.Parent = *netbox.NewNullableInt32(&parentID)
	} else if plan.Parent.IsNull() {
		locationRequest.SetParentNil()
	}

	// Set optional status
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		status := netbox.LocationStatusValue(plan.Status.ValueString())
		locationRequest.Status = &status
	}

	// Set optional tenant
	if !plan.Tenant.IsNull() && !plan.Tenant.IsUnknown() {
		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, plan.Tenant.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		locationRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	} else if plan.Tenant.IsNull() {
		locationRequest.SetTenantNil()
	}

	// Set optional facility
	if !plan.Facility.IsNull() && !plan.Facility.IsUnknown() {
		facility := plan.Facility.ValueString()
		locationRequest.Facility = &facility
	}

	// Set optional description
	utils.ApplyDescription(locationRequest, plan.Description)

	// Handle tags (prefer plan, fallback to state)
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, locationRequest, plan.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, locationRequest, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle custom fields with merge-aware logic
	utils.ApplyCustomFieldsWithMerge(ctx, locationRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	location, httpResp, err := r.client.DcimAPI.DcimLocationsUpdate(ctx, locationIDInt).WritableLocationRequest(*locationRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating location",
			utils.FormatAPIError(fmt.Sprintf("update location ID %s", locationID), err, httpResp),
		)
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error updating location",
			fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusOK, httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	r.mapLocationToState(ctx, location, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "updated a location resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LocationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	locationID := data.ID.ValueString()
	tflog.Debug(ctx, "Deleting location", map[string]interface{}{
		"id": locationID,
	})
	var locationIDInt int32
	locationIDInt, err := utils.ParseID(locationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Location ID",
			fmt.Sprintf("Location ID must be a number, got: %s", locationID),
		)
		return
	}
	httpResp, err := r.client.DcimAPI.DcimLocationsDestroy(ctx, locationIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting location",
			utils.FormatAPIError(fmt.Sprintf("delete location ID %s", locationID), err, httpResp),
		)
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting location",
			fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusNoContent, httpResp.StatusCode),
		)
		return
	}
	tflog.Trace(ctx, "deleted a location resource")
}

func (r *LocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		locationIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Location ID",
				fmt.Sprintf("Location ID must be a number, got: %s", parsed.ID),
			)
			return
		}
		location, httpResp, err := r.client.DcimAPI.DcimLocationsRetrieve(ctx, locationIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing location",
				utils.FormatAPIError(fmt.Sprintf("read location ID %s", parsed.ID), err, httpResp),
			)
			return
		}

		var data LocationResourceModel
		site := location.GetSite()
		if site.Id != 0 {
			data.Site = types.StringValue(fmt.Sprintf("%d", site.GetId()))
		}
		if location.HasParent() && location.GetParent().Id != 0 {
			parent := location.GetParent()
			data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		}
		if location.HasTenant() && location.GetTenant().Id != 0 {
			tenant := location.GetTenant()
			data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, location.HasTags(), location.GetTags(), data.Tags)
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

		r.mapLocationToState(ctx, location, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, location.GetCustomFields(), &resp.Diagnostics)
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

// mapLocationToState maps a Location API response to the Terraform state model.
func (r *LocationResource) mapLocationToState(ctx context.Context, location *netbox.Location, data *LocationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", location.GetId()))
	data.Name = types.StringValue(location.GetName())
	data.Slug = types.StringValue(location.GetSlug())

	// Site - preserve the user's configured value (ID, slug, or name)
	site := location.GetSite()
	data.Site = utils.UpdateReferenceAttribute(data.Site, site.GetName(), site.GetSlug(), site.Id)

	// Parent - preserve user's input format
	if location.HasParent() && location.GetParent().Id != 0 {
		parent := location.GetParent()
		data.Parent = utils.UpdateReferenceAttribute(data.Parent, parent.GetName(), parent.GetSlug(), parent.GetId())
	} else {
		data.Parent = types.StringNull()
	}

	// Status
	if location.HasStatus() {
		status := location.GetStatus()
		if status.Value != nil {
			data.Status = types.StringValue(string(*status.Value))
		} else {
			data.Status = types.StringValue("active")
		}
	} else {
		data.Status = types.StringValue("active")
	}

	// Tenant - preserve the user's configured value (ID, slug, or name)
	if location.HasTenant() && location.GetTenant().Id != 0 {
		tenant := location.GetTenant()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.Id)
	} else {
		data.Tenant = types.StringNull()
	}

	// Facility
	if location.HasFacility() {
		facility := location.GetFacility()
		if facility == "" {
			data.Facility = types.StringNull()
		} else {
			data.Facility = types.StringValue(facility)
		}
	} else {
		data.Facility = types.StringNull()
	}

	// Description
	if location.HasDescription() {
		desc := location.GetDescription()
		if desc == "" {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags using filter-to-owned approach
	planTags := data.Tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, location.HasTags(), location.GetTags(), planTags)

	// Handle custom fields - only populate fields that are in plan (owned by this resource)
	if location.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, location.GetCustomFields(), diags)
	}
}
