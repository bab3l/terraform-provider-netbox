package resources

import (
	"context"
	"fmt"

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

var _ resource.Resource = &TenantGroupResource{}

var _ resource.ResourceWithImportState = &TenantGroupResource{}

func NewTenantGroupResource() resource.Resource {

	return &TenantGroupResource{}

}

type TenantGroupResource struct {
	client *netbox.APIClient
}

// GetClient returns the API client for testing purposes.

func (r *TenantGroupResource) GetClient() *netbox.APIClient {

	return r.client

}

type TenantGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Parent types.String `tfsdk:"parent"`

	ParentID types.String `tfsdk:"parent_id"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *TenantGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_tenant_group"

}

func (r *TenantGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a tenant group in Netbox. Tenant groups provide a hierarchical way to organize tenants for multi-tenancy scenarios.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("tenant group"),

			"name": nbschema.NameAttribute("tenant group", 100),

			"slug": nbschema.SlugAttribute("tenant group"),

			"parent": nbschema.ReferenceAttribute("tenant group", "ID or slug of the parent tenant group. Leave empty for top-level groups."),

			"parent_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the parent tenant group.",
			},

			"description": nbschema.DescriptionAttribute("tenant group"),

			"display_name": nbschema.DisplayNameAttribute("tenant group"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *TenantGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *TenantGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data TenantGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating tenant group", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Build the request

	tenantGroupRequest := netbox.WritableTenantGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),

		Description: utils.StringPtr(data.Description),
	}

	// Set parent if provided

	if utils.IsSet(data.Parent) {

		parentID, parentDiags := netboxlookup.LookupTenantGroupID(ctx, r.client, data.Parent.ValueString())

		resp.Diagnostics.Append(parentDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		tenantGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		tenantGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		tenantGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Create via API

	tenantGroup, httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsCreate(ctx).WritableTenantGroupRequest(tenantGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		handler := utils.CreateErrorHandler{

			ResourceType: "netbox_tenant_group",

			ResourceName: "this.tenant_group",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {

				list, _, lookupErr := r.client.TenancyAPI.TenancyTenantGroupsList(lookupCtx).Slug([]string{slug}).Execute()

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

		resp.Diagnostics.AddError("Error creating tenant group", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))

		return

	}

	r.mapTenantGroupToState(ctx, tenantGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "created a tenant group resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TenantGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data TenantGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tenantGroupID := data.ID.ValueString()

	tenantGroupIDInt := utils.ParseInt32FromString(tenantGroupID)

	if tenantGroupIDInt == 0 {

		resp.Diagnostics.AddError("Invalid Tenant Group ID", fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID))

		return

	}

	tenantGroup, httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, tenantGroupIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error reading tenant group", utils.FormatAPIError(fmt.Sprintf("read tenant group ID %s", tenantGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode == 404 {

		resp.State.RemoveResource(ctx)

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error reading tenant group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	r.mapTenantGroupToState(ctx, tenantGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TenantGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data TenantGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tenantGroupID := data.ID.ValueString()

	tenantGroupIDInt := utils.ParseInt32FromString(tenantGroupID)

	if tenantGroupIDInt == 0 {

		resp.Diagnostics.AddError("Invalid Tenant Group ID", fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID))

		return

	}

	tflog.Debug(ctx, "Updating tenant group", map[string]interface{}{

		"id": tenantGroupID,

		"name": data.Name.ValueString(),
	})

	// Build the request

	tenantGroupRequest := netbox.WritableTenantGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),

		Description: utils.StringPtr(data.Description),
	}

	// Set parent if provided

	if utils.IsSet(data.Parent) {

		parentID, parentDiags := netboxlookup.LookupTenantGroupID(ctx, r.client, data.Parent.ValueString())

		resp.Diagnostics.Append(parentDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		tenantGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		tenantGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		tenantGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Update via API

	tenantGroup, httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsUpdate(ctx, tenantGroupIDInt).WritableTenantGroupRequest(tenantGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error updating tenant group", utils.FormatAPIError(fmt.Sprintf("update tenant group ID %s", tenantGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error updating tenant group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	r.mapTenantGroupToState(ctx, tenantGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TenantGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data TenantGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tenantGroupID := data.ID.ValueString()

	tenantGroupIDInt := utils.ParseInt32FromString(tenantGroupID)

	if tenantGroupIDInt == 0 {

		resp.Diagnostics.AddError("Invalid Tenant Group ID", fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID))

		return

	}

	tflog.Debug(ctx, "Deleting tenant group", map[string]interface{}{"id": tenantGroupID})

	httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsDestroy(ctx, tenantGroupIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error deleting tenant group", utils.FormatAPIError(fmt.Sprintf("delete tenant group ID %s", tenantGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 204 {

		resp.Diagnostics.AddError("Error deleting tenant group", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return

	}

	tflog.Trace(ctx, "deleted a tenant group resource")

}

func (r *TenantGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapTenantGroupToState maps API response to Terraform state.

func (r *TenantGroupResource) mapTenantGroupToState(ctx context.Context, tenantGroup *netbox.TenantGroup, data *TenantGroupResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", tenantGroup.GetId()))

	data.Name = types.StringValue(tenantGroup.GetName())

	data.Slug = types.StringValue(tenantGroup.GetSlug())

	data.Description = utils.StringFromAPI(tenantGroup.HasDescription(), tenantGroup.GetDescription, data.Description)

	// Handle display_name

	if tenantGroup.GetDisplay() != "" {

		data.DisplayName = types.StringValue(tenantGroup.GetDisplay())

	} else {

		data.DisplayName = types.StringNull()

	}

	// Handle parent reference

	if tenantGroup.HasParent() {

		parent := tenantGroup.GetParent()

		if parent.GetId() != 0 {

			data.ParentID = types.StringValue(fmt.Sprintf("%d", parent.GetId()))

			userParent := data.Parent.ValueString()

			if userParent == parent.GetName() || userParent == parent.GetSlug() || userParent == parent.GetDisplay() || userParent == fmt.Sprintf("%d", parent.GetId()) {

				// Keep user's original value

			} else {

				data.Parent = types.StringValue(parent.GetName())

			}

		} else {

			data.Parent = types.StringNull()
			data.ParentID = types.StringNull()

		}

	} else {

		data.Parent = types.StringNull()
		data.ParentID = types.StringNull()

	}

	// Handle tags

	if tenantGroup.HasTags() {

		tags := utils.NestedTagsToTagModels(tenantGroup.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if !diags.HasError() {

			data.Tags = tagsValue

		}

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	switch {

	case tenantGroup.HasCustomFields() && !data.CustomFields.IsNull():

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if !diags.HasError() {

			customFields := utils.MapToCustomFieldModels(tenantGroup.GetCustomFields(), stateCustomFields)

			customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

			diags.Append(cfValueDiags...)

			if !diags.HasError() {

				data.CustomFields = customFieldsValue

			}

		}

	case tenantGroup.HasCustomFields():

		customFields := utils.MapToCustomFieldModels(tenantGroup.GetCustomFields(), []utils.CustomFieldModel{})

		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfValueDiags...)

		if !diags.HasError() {

			data.CustomFields = customFieldsValue

		}

	default:

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
