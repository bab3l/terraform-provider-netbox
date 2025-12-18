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

var _ resource.Resource = &ContactGroupResource{}

var _ resource.ResourceWithImportState = &ContactGroupResource{}

func NewContactGroupResource() resource.Resource {

	return &ContactGroupResource{}

}

type ContactGroupResource struct {
	client *netbox.APIClient
}

// GetClient returns the API client for testing purposes.

func (r *ContactGroupResource) GetClient() *netbox.APIClient {

	return r.client

}

type ContactGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Parent types.String `tfsdk:"parent"`

	ParentID types.String `tfsdk:"parent_id"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *ContactGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_contact_group"

}

func (r *ContactGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a contact group in Netbox. Contact groups provide a hierarchical way to organize contacts for better management.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("contact group"),

			"name": nbschema.NameAttribute("contact group", 100),

			"slug": nbschema.SlugAttribute("contact group"),

			"parent": nbschema.ReferenceAttribute("contact group", "ID or slug of the parent contact group. Leave empty for top-level groups."),

			"parent_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the parent contact group.",
			},

			"description": nbschema.DescriptionAttribute("contact group"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *ContactGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ContactGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ContactGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating contact group", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Build the request

	contactGroupRequest := netbox.WritableContactGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),

		Description: utils.StringPtr(data.Description),
	}

	// Set parent if provided

	if utils.IsSet(data.Parent) {

		parentID, parentDiags := netboxlookup.LookupContactGroupID(ctx, r.client, data.Parent.ValueString())

		resp.Diagnostics.Append(parentDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		contactGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		contactGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		contactGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Create via API

	contactGroup, httpResp, err := r.client.TenancyAPI.TenancyContactGroupsCreate(ctx).WritableContactGroupRequest(contactGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		handler := utils.CreateErrorHandler{

			ResourceType: "netbox_contact_group",

			ResourceName: "this.contact_group",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {

				list, _, lookupErr := r.client.TenancyAPI.TenancyContactGroupsList(lookupCtx).Slug([]string{slug}).Execute()

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

		resp.Diagnostics.AddError("Error creating contact group", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))

		return

	}

	r.mapContactGroupToState(ctx, contactGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "created a contact group resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ContactGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ContactGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	contactGroupID := data.ID.ValueString()

	contactGroupIDInt := utils.ParseInt32FromString(contactGroupID)

	if contactGroupIDInt == 0 {

		resp.Diagnostics.AddError("Invalid Contact Group ID", fmt.Sprintf("Contact Group ID must be a number, got: %s", contactGroupID))

		return

	}

	contactGroup, httpResp, err := r.client.TenancyAPI.TenancyContactGroupsRetrieve(ctx, contactGroupIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error reading contact group", utils.FormatAPIError(fmt.Sprintf("read contact group ID %s", contactGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode == 404 {

		resp.State.RemoveResource(ctx)

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error reading contact group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	r.mapContactGroupToState(ctx, contactGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ContactGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ContactGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	contactGroupID := data.ID.ValueString()

	contactGroupIDInt := utils.ParseInt32FromString(contactGroupID)

	if contactGroupIDInt == 0 {

		resp.Diagnostics.AddError("Invalid Contact Group ID", fmt.Sprintf("Contact Group ID must be a number, got: %s", contactGroupID))

		return

	}

	tflog.Debug(ctx, "Updating contact group", map[string]interface{}{

		"id": contactGroupID,

		"name": data.Name.ValueString(),
	})

	// Build the request

	contactGroupRequest := netbox.WritableContactGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),

		Description: utils.StringPtr(data.Description),
	}

	// Set parent if provided

	if utils.IsSet(data.Parent) {

		parentID, parentDiags := netboxlookup.LookupContactGroupID(ctx, r.client, data.Parent.ValueString())

		resp.Diagnostics.Append(parentDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		contactGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		contactGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		contactGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Update via API

	contactGroup, httpResp, err := r.client.TenancyAPI.TenancyContactGroupsUpdate(ctx, contactGroupIDInt).WritableContactGroupRequest(contactGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error updating contact group", utils.FormatAPIError(fmt.Sprintf("update contact group ID %s", contactGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error updating contact group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	r.mapContactGroupToState(ctx, contactGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ContactGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ContactGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	contactGroupID := data.ID.ValueString()

	contactGroupIDInt := utils.ParseInt32FromString(contactGroupID)

	if contactGroupIDInt == 0 {

		resp.Diagnostics.AddError("Invalid Contact Group ID", fmt.Sprintf("Contact Group ID must be a number, got: %s", contactGroupID))

		return

	}

	tflog.Debug(ctx, "Deleting contact group", map[string]interface{}{"id": contactGroupID})

	httpResp, err := r.client.TenancyAPI.TenancyContactGroupsDestroy(ctx, contactGroupIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error deleting contact group", utils.FormatAPIError(fmt.Sprintf("delete contact group ID %s", contactGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 204 {

		resp.Diagnostics.AddError("Error deleting contact group", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return

	}

	tflog.Trace(ctx, "deleted a contact group resource")

}

func (r *ContactGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapContactGroupToState maps API response to Terraform state.

func (r *ContactGroupResource) mapContactGroupToState(ctx context.Context, contactGroup *netbox.ContactGroup, data *ContactGroupResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", contactGroup.GetId()))

	data.Name = types.StringValue(contactGroup.GetName())

	data.Slug = types.StringValue(contactGroup.GetSlug())

	data.Description = utils.StringFromAPI(contactGroup.HasDescription(), contactGroup.GetDescription, data.Description)

	// Handle parent reference

	if contactGroup.HasParent() {

		parent := contactGroup.GetParent()

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

	if contactGroup.HasTags() {

		tags := utils.NestedTagsToTagModels(contactGroup.GetTags())

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

	case contactGroup.HasCustomFields() && !data.CustomFields.IsNull():

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if !diags.HasError() {

			customFields := utils.MapToCustomFieldModels(contactGroup.GetCustomFields(), stateCustomFields)

			customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

			diags.Append(cfValueDiags...)

			if !diags.HasError() {

				data.CustomFields = customFieldsValue

			}

		}

	case contactGroup.HasCustomFields():

		customFields := utils.MapToCustomFieldModels(contactGroup.GetCustomFields(), []utils.CustomFieldModel{})

		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfValueDiags...)

		if !diags.HasError() {

			data.CustomFields = customFieldsValue

		}

	default:

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
