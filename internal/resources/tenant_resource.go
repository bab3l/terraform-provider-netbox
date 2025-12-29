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

	DisplayName types.String `tfsdk:"display_name"`

	Group types.String `tfsdk:"group"`

	GroupID types.String `tfsdk:"group_id"`

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

			"display_name": nbschema.DisplayNameAttribute("tenant"),

			"group":    nbschema.ReferenceAttribute("tenant group", "Name, Slug, or ID of the tenant group that this tenant belongs to."),
			"group_id": nbschema.ComputedIDAttribute("tenant group"),
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("tenant"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())

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

	}

	// Apply common metadata fields (tags, custom_fields)
	utils.ApplyMetadataFields(ctx, &tenantRequest, data.Tags, data.CustomFields, &resp.Diagnostics)
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
	r.mapTenantToState(ctx, tenant, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created tenant", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

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
	r.mapTenantToState(ctx, tenant, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data TenantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	// Prepare the tenant update request

	tenantRequest := netbox.TenantRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Apply common descriptive fields (description, comments)
	utils.ApplyDescriptiveFields(&tenantRequest, data.Description, data.Comments)

	// Handle group relationship

	if groupRef := utils.ResolveOptionalReference(ctx, r.client, data.Group, netboxlookup.LookupTenantGroup, &resp.Diagnostics); groupRef != nil {

		tenantRequest.Group = *netbox.NewNullableBriefTenantGroupRequest(groupRef)

	}

	// Apply common metadata fields (tags, custom_fields)
	utils.ApplyMetadataFields(ctx, &tenantRequest, data.Tags, data.CustomFields, &resp.Diagnostics)
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

	// Map response to state using helper
	r.mapTenantToState(ctx, tenant, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapTenantToState maps API response to Terraform state using state helpers.

func (r *TenantResource) mapTenantToState(ctx context.Context, tenant *netbox.Tenant, data *TenantResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))

	data.Name = types.StringValue(tenant.GetName())

	data.Slug = types.StringValue(tenant.GetSlug())

	data.DisplayName = types.StringValue(tenant.GetDisplay())

	// Handle group reference
	groupRef := utils.PreserveOptionalReferenceWithID(
		data.Group,
		tenant.HasGroup() && tenant.GetGroup().Id != 0,
		tenant.GetGroup().Id,
		tenant.GetGroup().Name,
		tenant.GetGroup().Slug,
	)
	data.Group = groupRef.Reference
	data.GroupID = groupRef.ID

	// Handle optional string fields using helpers
	data.Description = utils.StringFromAPI(tenant.HasDescription(), tenant.GetDescription, data.Description)
	data.Comments = utils.StringFromAPI(tenant.HasComments(), tenant.GetComments, data.Comments)

	// Handle tags
	data.Tags = utils.PopulateTagsFromNestedTags(ctx, tenant.HasTags(), tenant.GetTags(), diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields
	data.CustomFields = utils.PopulateCustomFieldsFromMap(ctx, tenant.HasCustomFields(), tenant.GetCustomFields(), data.CustomFields, diags)
}
