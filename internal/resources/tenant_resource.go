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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &TenantResource{}

var _ resource.ResourceWithImportState = &TenantResource{}
var _ resource.ResourceWithIdentity = &TenantResource{}

func NewTenantResource() resource.Resource {
	return &TenantResource{}
}

// TenantResource defines the resource implementation.

type TenantResource struct {
	client *netbox.APIClient
}

// TenantResourceModel describes the resource data model.

type TenantResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Group types.String `tfsdk:"group"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *TenantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (r *TenantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tenant in Netbox. Tenants represent individual customers or organizational units in multi-tenancy scenarios, allowing you to organize and track resources by client or department.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("tenant"),

			"name": nbschema.NameAttribute("tenant", 100),

			"slug": nbschema.SlugAttribute("tenant"),

			"group": nbschema.ReferenceAttributeWithDiffSuppress("tenant group", "ID or slug of the tenant group that this tenant belongs to."),
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("tenant"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *TenantResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *TenantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TenantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating tenant", map[string]interface{}{
		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Prepare the tenant request

	tenantRequest := netbox.TenantRequest{
		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Apply common descriptive fields (description, comments)
	utils.ApplyDescriptiveFields(&tenantRequest, data.Description, data.Comments)

	// Handle group relationship - lookup the group details by ID

	if groupRef := utils.ResolveOptionalReference(ctx, r.client, data.Group, netboxlookup.LookupTenantGroup, &resp.Diagnostics); groupRef != nil {
		tenantRequest.Group = *netbox.NewNullableBriefTenantGroupRequest(groupRef)
	} else if data.Group.IsNull() {
		tenantRequest.SetGroupNil()
	}

	// Apply common metadata fields (tags, custom_fields)
	utils.ApplyTagsFromSlugs(ctx, r.client, &tenantRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, &tenantRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the tenant via API

	tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsCreate(ctx).TenantRequest(tenantRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		// Use enhanced error handler that detects duplicates and provides import hints

		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_tenant",

			ResourceName: "this.tenant",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				list, _, lookupErr := r.client.TenancyAPI.TenancyTenantsList(lookupCtx).Slug([]string{slug}).Execute()

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
		resp.Diagnostics.AddError("Error creating tenant", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))

		return
	}

	if tenant == nil {
		resp.Diagnostics.AddError("Tenant API returned nil", "No tenant object returned from Netbox API.")

		return
	}

	// Map response to state using helper
	planTags := data.Tags
	planCustomFields := data.CustomFields

	r.mapTenantToState(ctx, tenant, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	var tagSlugs []string
	switch {
	case planTags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(planTags.Elements()) == 0:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	case tenant.HasTags():
		for _, tag := range tenant.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, tenant.GetCustomFields(), &resp.Diagnostics)

	tflog.Debug(ctx, "Created tenant", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TenantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := data.ID.ValueString()

	var tenantIDInt int32

	tenantIDInt, err := utils.ParseID(tenantID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID))

		return
	}

	tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading tenant", utils.FormatAPIError(fmt.Sprintf("read tenant ID %s", tenantID), err, httpResp))

		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading tenant", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return
	}

	// Map response to state using helper
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	r.mapTenantToState(ctx, tenant, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve null/empty state values for tags and custom_fields
	var stateTagSlugs []string
	switch {
	case stateTags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(stateTags.Elements()) == 0:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	case tenant.HasTags():
		for _, tag := range tenant.GetTags() {
			stateTagSlugs = append(stateTagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, stateTagSlugs)
	default:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, tenant.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read both state and plan for merge-aware custom fields
	var state, plan TenantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := plan.ID.ValueString()

	var tenantIDInt int32

	tenantIDInt, err := utils.ParseID(tenantID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID))

		return
	}

	// Prepare the tenant update request

	tenantRequest := netbox.TenantRequest{
		Name: plan.Name.ValueString(),

		Slug: plan.Slug.ValueString(),
	}

	// Apply common descriptive fields (description, comments)
	utils.ApplyDescriptiveFields(&tenantRequest, plan.Description, plan.Comments)

	// Handle group relationship

	if groupRef := utils.ResolveOptionalReference(ctx, r.client, plan.Group, netboxlookup.LookupTenantGroup, &resp.Diagnostics); groupRef != nil {
		tenantRequest.Group = *netbox.NewNullableBriefTenantGroupRequest(groupRef)
	} else if plan.Group.IsNull() {
		tenantRequest.SetGroupNil()
	}

	// Apply common metadata fields (tags, custom_fields) with merge-aware helpers
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &tenantRequest, plan.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &tenantRequest, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFieldsWithMerge(ctx, &tenantRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsUpdate(ctx, tenantIDInt).TenantRequest(tenantRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant", utils.FormatAPIError(fmt.Sprintf("update tenant ID %s", tenantID), err, httpResp))

		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating tenant", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return
	}

	// Save plan state for filter-to-owned pattern
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Map response to state using helper
	r.mapTenantToState(ctx, tenant, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	var updateTagSlugs []string
	switch {
	case planTags.IsNull():
		plan.Tags = types.SetNull(types.StringType)
	case len(planTags.Elements()) == 0:
		plan.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	case tenant.HasTags():
		for _, tag := range tenant.GetTags() {
			updateTagSlugs = append(updateTagSlugs, tag.GetSlug())
		}
		plan.Tags = utils.TagsSlugToSet(ctx, updateTagSlugs)
	default:
		plan.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, tenant.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TenantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := data.ID.ValueString()

	var tenantIDInt int32

	tenantIDInt, err := utils.ParseID(tenantID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID))

		return
	}

	httpResp, err := r.client.TenancyAPI.TenancyTenantsDestroy(ctx, tenantIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError("Error deleting tenant", utils.FormatAPIError(fmt.Sprintf("delete tenant ID %s", tenantID), err, httpResp))

		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error deleting tenant", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return
	}
}

func (r *TenantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		tenantIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Tenant ID must be a number, got: %s", parsed.ID))
			return
		}
		tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing tenant", utils.FormatAPIError(fmt.Sprintf("read tenant ID %s", parsed.ID), err, httpResp))
			return
		}

		var data TenantResourceModel
		if tenant.HasGroup() && tenant.GetGroup().Id != 0 {
			group := tenant.GetGroup()
			data.Group = types.StringValue(fmt.Sprintf("%d", group.Id))
		}
		if tenant.HasTags() {
			tagSlugs := make([]string, 0, len(tenant.GetTags()))
			for _, tag := range tenant.GetTags() {
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

		r.mapTenantToState(ctx, tenant, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, tenant.GetCustomFields(), &resp.Diagnostics)
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapTenantToState maps API response to Terraform state using state helpers.

func (r *TenantResource) mapTenantToState(ctx context.Context, tenant *netbox.Tenant, data *TenantResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))

	data.Name = types.StringValue(tenant.GetName())

	data.Slug = types.StringValue(tenant.GetSlug())

	// Handle group reference
	groupRef := utils.PreserveOptionalReferenceWithID(
		data.Group,
		tenant.HasGroup() && tenant.GetGroup().Id != 0,
		tenant.GetGroup().Id,
		tenant.GetGroup().Name,
		tenant.GetGroup().Slug,
	)
	data.Group = groupRef.Reference

	// Handle optional string fields using helpers
	data.Description = utils.StringFromAPI(tenant.HasDescription(), tenant.GetDescription, data.Description)
	data.Comments = utils.StringFromAPI(tenant.HasComments(), tenant.GetComments, data.Comments)

	// Tags and custom fields are now handled in Create/Read/Update with filter-to-owned pattern
}
