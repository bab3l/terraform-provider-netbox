// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &LocationResource{}

var _ resource.ResourceWithImportState = &LocationResource{}

func NewLocationResource() resource.Resource {
	return &LocationResource{}
}

// LocationResource defines the resource implementation.

type LocationResource struct {
	client *netbox.APIClient
}

// LocationResourceModel describes the resource data model.

type LocationResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Site types.String `tfsdk:"site"`

	Parent types.String `tfsdk:"parent"`

	Status types.String `tfsdk:"status"`

	Tenant types.String `tfsdk:"tenant"`

	Facility types.String `tfsdk:"facility"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *LocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (r *LocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a location in Netbox. Locations represent physical areas within a site, such as buildings, floors, or rooms. Locations can be nested hierarchically.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("location"),

			"name": nbschema.NameAttribute("location", 100),

			"slug": nbschema.SlugAttribute("location"),

			"site": nbschema.RequiredReferenceAttribute("site", "ID or slug of the site this location belongs to. Required."),

			"parent": nbschema.ReferenceAttribute("parent location", "ID or slug of the parent location. Leave empty for top-level locations within the site."),

			"status": nbschema.StatusAttribute([]string{"planned", "staging", "active", "decommissioning", "retired"}, "Operational status of the location. Defaults to `active`."),

			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this location."),

			"facility": nbschema.FacilityAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("location"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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
	}

	// Set optional facility

	if !data.Facility.IsNull() && !data.Facility.IsUnknown() {
		facility := data.Facility.ValueString()

		locationRequest.Facility = &facility
	}

	// Set optional description

	utils.ApplyDescription(locationRequest, data.Description)

	// Handle tags

	utils.ApplyMetadataFields(ctx, locationRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

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

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(

			"Error creating location",

			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)

		return
	}

	// Map response to state

	r.mapLocationToState(ctx, location, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

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
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(

			"Error reading location",

			utils.FormatAPIError(fmt.Sprintf("read location ID %s", locationID), err, httpResp),
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

	r.mapLocationToState(ctx, location, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LocationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	locationID := data.ID.ValueString()

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
	}

	// Set optional facility

	if !data.Facility.IsNull() && !data.Facility.IsUnknown() {
		facility := data.Facility.ValueString()

		locationRequest.Facility = &facility
	}

	// Set optional description

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()

		locationRequest.Description = &desc
	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		locationRequest.Tags = tags
	}

	// Apply metadata fields (tags, custom_fields)

	utils.ApplyMetadataFields(ctx, locationRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

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

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(

			"Error updating location",

			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)

		return
	}

	// Map response to state

	r.mapLocationToState(ctx, location, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updated a location resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting location",

			utils.FormatAPIError(fmt.Sprintf("delete location ID %s", locationID), err, httpResp),
		)

		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(

			"Error deleting location",

			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)

		return
	}

	tflog.Trace(ctx, "deleted a location resource")
}

func (r *LocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

		switch {
		case facility == "" && data.Facility.IsNull():

			data.Facility = types.StringNull()

		case facility == "":

			data.Facility = types.StringNull()

		default:

			data.Facility = types.StringValue(facility)
		}
	} else {
		data.Facility = types.StringNull()
	}

	// Description

	if location.HasDescription() {
		desc := location.GetDescription()

		switch {
		case desc == "" && data.Description.IsNull():

			data.Description = types.StringNull()

		case desc == "":

			data.Description = types.StringNull()

		default:

			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags

	if location.HasTags() {
		tags := utils.NestedTagsToTagModels(location.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if diags.HasError() {
			return
		}

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields

	if location.HasCustomFields() {
		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {
			cfDiags := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(cfDiags...)

			if diags.HasError() {
				return
			}
		}

		customFields := utils.MapToCustomFieldModels(location.GetCustomFields(), existingModels)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		if diags.HasError() {
			return
		}

		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
