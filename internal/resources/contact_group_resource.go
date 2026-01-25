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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ContactGroupResource{}
var _ resource.ResourceWithImportState = &ContactGroupResource{}
var _ resource.ResourceWithIdentity = &ContactGroupResource{}

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
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	ParentID     types.String `tfsdk:"parent_id"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *ContactGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact_group"
}

func (r *ContactGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a contact group in Netbox. Contact groups provide a hierarchical way to organize contacts for better management.",
		Attributes: map[string]schema.Attribute{
			"id":     nbschema.IDAttribute("contact group"),
			"name":   nbschema.NameAttribute("contact group", 100),
			"slug":   nbschema.SlugAttribute("contact group"),
			"parent": nbschema.ReferenceAttribute("contact group", "ID or slug of the parent contact group. Leave empty for top-level groups."),
			"parent_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the parent contact group.",
			},
			"tags":          nbschema.TagsSlugAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("contact group"))

	// Tags and custom fields are defined directly in the schema above.
}

func (r *ContactGroupResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
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
	}

	// Apply description
	utils.ApplyDescription(&contactGroupRequest, data.Description)

	// Set parent if provided
	if utils.IsSet(data.Parent) {
		parentID, parentDiags := netboxlookup.LookupContactGroupID(ctx, r.client, data.Parent.ValueString())
		resp.Diagnostics.Append(parentDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contactGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)
	}

	// Apply tags and custom_fields
	utils.ApplyTagsFromSlugs(ctx, r.client, &contactGroupRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, &contactGroupRequest, data.CustomFields, &resp.Diagnostics)

	// Store plan values for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Create via API
	contactGroup, httpResp, err := r.client.TenancyAPI.TenancyContactGroupsCreate(ctx).WritableContactGroupRequest(contactGroupRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_contact_group",
			ResourceName: "this.contact_group",
			SlugValue:    data.Slug.ValueString(),
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

	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating contact group", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}

	r.mapContactGroupToState(contactGroup, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, contactGroup.HasTags(), contactGroup.GetTags(), planTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, contactGroup.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

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
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading contact group", utils.FormatAPIError(fmt.Sprintf("read contact group ID %s", contactGroupID), err, httpResp))
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error reading contact group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Store state values for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	r.mapContactGroupToState(contactGroup, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, contactGroup.HasTags(), contactGroup.GetTags(), stateTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, contactGroup.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ContactGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactGroupID := plan.ID.ValueString()
	contactGroupIDInt := utils.ParseInt32FromString(contactGroupID)
	if contactGroupIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Contact Group ID", fmt.Sprintf("Contact Group ID must be a number, got: %s", contactGroupID))
		return
	}
	tflog.Debug(ctx, "Updating contact group", map[string]interface{}{
		"id":   contactGroupID,
		"name": plan.Name.ValueString(),
	})

	// Build the request
	contactGroupRequest := netbox.WritableContactGroupRequest{
		Name: plan.Name.ValueString(),
		Slug: plan.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&contactGroupRequest, plan.Description)

	// Set parent if provided
	if utils.IsSet(plan.Parent) {
		parentID, parentDiags := netboxlookup.LookupContactGroupID(ctx, r.client, plan.Parent.ValueString())
		resp.Diagnostics.Append(parentDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contactGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)
	} else if plan.Parent.IsNull() {
		contactGroupRequest.SetParentNil()
	}

	// Apply tags and custom_fields with merge-aware helpers
	utils.ApplyTagsFromSlugs(ctx, r.client, &contactGroupRequest, plan.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFieldsWithMerge(ctx, &contactGroupRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)

	// Update via API
	contactGroup, httpResp, err := r.client.TenancyAPI.TenancyContactGroupsUpdate(ctx, contactGroupIDInt).WritableContactGroupRequest(contactGroupRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating contact group", utils.FormatAPIError(fmt.Sprintf("update contact group ID %s", contactGroupID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating contact group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Map response to model
	r.mapContactGroupToState(contactGroup, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	plan.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, contactGroup.HasTags(), contactGroup.GetTags(), plan.Tags)
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, plan.CustomFields, contactGroup.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting contact group", utils.FormatAPIError(fmt.Sprintf("delete contact group ID %s", contactGroupID), err, httpResp))
		return
	}

	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError("Error deleting contact group", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
	tflog.Trace(ctx, "deleted a contact group resource")
}

func (r *ContactGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		contactGroupIDInt := utils.ParseInt32FromString(parsed.ID)
		if contactGroupIDInt == 0 {
			resp.Diagnostics.AddError("Invalid Contact Group ID", fmt.Sprintf("Contact Group ID must be a number, got: %s", parsed.ID))
			return
		}

		contactGroup, httpResp, err := r.client.TenancyAPI.TenancyContactGroupsRetrieve(ctx, contactGroupIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing contact group", utils.FormatAPIError(fmt.Sprintf("read contact group ID %s", parsed.ID), err, httpResp))
			return
		}
		if httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError("Error importing contact group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
			return
		}

		var data ContactGroupResourceModel
		if contactGroup.HasParent() {
			parent := contactGroup.GetParent()
			if parent.GetId() != 0 {
				data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
			}
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

		r.mapContactGroupToState(contactGroup, &data)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, contactGroup.HasTags(), contactGroup.GetTags(), data.Tags)

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, contactGroup.GetCustomFields(), &resp.Diagnostics)
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

// mapContactGroupToState maps API response to Terraform state.
func (r *ContactGroupResource) mapContactGroupToState(contactGroup *netbox.ContactGroup, data *ContactGroupResourceModel) {
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

	// Tags and custom fields are now handled in Create/Read/Update methods using filter-to-owned pattern
}
