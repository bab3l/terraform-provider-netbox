// Package resources contains Terraform resource implementations for the Netbox provider.

//

// This package integrates with the go-netbox OpenAPI client to provide

// CRUD operations for Netbox resources via Terraform.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &SiteGroupResource{}

var _ resource.ResourceWithImportState = &SiteGroupResource{}

func NewSiteGroupResource() resource.Resource {

	return &SiteGroupResource{}

}

// SiteGroupResource defines the resource implementation.

type SiteGroupResource struct {
	client *netbox.APIClient
}

// SiteGroupResourceModel describes the resource data model.

type SiteGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Parent types.String `tfsdk:"parent"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *SiteGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_site_group"

}

func (r *SiteGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a site group in Netbox. Site groups provide a hierarchical way to organize sites, allowing you to create nested organizational structures for better management and reporting of your physical locations.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("site group"),

			"name": nbschema.NameAttribute("site group", 100),

			"slug": nbschema.SlugAttribute("site group"),

			"parent": nbschema.IDOnlyReferenceAttribute("parent site group", "ID of the parent site group. Leave empty for top-level site groups."),

			"description": nbschema.DescriptionAttribute("site group"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *SiteGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *SiteGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data SiteGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating site group", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Prepare the site group request

	siteGroupRequest := netbox.WritableSiteGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string field

	siteGroupRequest.Description = utils.StringPtr(data.Description)

	// Handle parent reference

	if utils.IsSet(data.Parent) {

		parentIDInt, err := utils.ParseID(data.Parent.ValueString())

		if err != nil {

			resp.Diagnostics.AddError("Invalid Parent ID", fmt.Sprintf("Parent ID must be a number, got: %s", data.Parent.ValueString()))

			return

		}

		siteGroupRequest.Parent = *netbox.NewNullableInt32(&parentIDInt)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteGroupRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	// Create the site group via API

	siteGroup, httpResp, err := r.client.DcimAPI.DcimSiteGroupsCreate(ctx).WritableSiteGroupRequest(siteGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		handler := utils.CreateErrorHandler{

			ResourceType: "netbox_site_group",

			ResourceName: "this.site_group",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {

				list, _, lookupErr := r.client.DcimAPI.DcimSiteGroupsList(lookupCtx).Slug([]string{slug}).Execute()

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

		resp.Diagnostics.AddError("Error creating site group", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))

		return

	}

	if siteGroup == nil {

		resp.Diagnostics.AddError("SiteGroup API returned nil", "No site group object returned from Netbox API.")

		return

	}

	// Map response to state using helper

	r.mapSiteGroupToState(ctx, siteGroup, &data)

	tflog.Debug(ctx, "Created site group", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SiteGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data SiteGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	siteGroupID := data.ID.ValueString()

	var siteGroupIDInt int32

	siteGroupIDInt, err := utils.ParseID(siteGroupID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Site Group ID", fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID))

		return

	}

	siteGroup, httpResp, err := r.client.DcimAPI.DcimSiteGroupsRetrieve(ctx, siteGroupIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error reading site group", utils.FormatAPIError(fmt.Sprintf("read site group ID %s", siteGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode == 404 {

		resp.State.RemoveResource(ctx)

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error reading site group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state using helper

	r.mapSiteGroupToState(ctx, siteGroup, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SiteGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data SiteGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	siteGroupID := data.ID.ValueString()

	var siteGroupIDInt int32

	siteGroupIDInt, err := utils.ParseID(siteGroupID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Site Group ID", fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID))

		return

	}

	// Prepare the site group update request

	siteGroupRequest := netbox.WritableSiteGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string field

	siteGroupRequest.Description = utils.StringPtr(data.Description)

	// Handle parent reference

	if utils.IsSet(data.Parent) {

		parentIDInt, err := utils.ParseID(data.Parent.ValueString())

		if err != nil {

			resp.Diagnostics.AddError("Invalid Parent ID", fmt.Sprintf("Parent ID must be a number, got: %s", data.Parent.ValueString()))

			return

		}

		siteGroupRequest.Parent = *netbox.NewNullableInt32(&parentIDInt)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteGroupRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		siteGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	siteGroup, httpResp, err := r.client.DcimAPI.DcimSiteGroupsUpdate(ctx, siteGroupIDInt).WritableSiteGroupRequest(siteGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error updating site group", utils.FormatAPIError(fmt.Sprintf("update site group ID %s", siteGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error updating site group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state using helper

	r.mapSiteGroupToState(ctx, siteGroup, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SiteGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data SiteGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	siteGroupID := data.ID.ValueString()

	var siteGroupIDInt int32

	siteGroupIDInt, err := utils.ParseID(siteGroupID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Site Group ID", fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID))

		return

	}

	httpResp, err := r.client.DcimAPI.DcimSiteGroupsDestroy(ctx, siteGroupIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error deleting site group", utils.FormatAPIError(fmt.Sprintf("delete site group ID %s", siteGroupID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 204 {

		resp.Diagnostics.AddError("Error deleting site group", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))

		return

	}

}

func (r *SiteGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapSiteGroupToState maps API response to Terraform state using state helpers.

func (r *SiteGroupResource) mapSiteGroupToState(ctx context.Context, siteGroup *netbox.SiteGroup, data *SiteGroupResourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", siteGroup.GetId()))

	data.Name = types.StringValue(siteGroup.GetName())

	data.Slug = types.StringValue(siteGroup.GetSlug())

	// Handle parent reference

	if siteGroup.HasParent() {

		parent := siteGroup.GetParent()

		if parent.GetId() != 0 {

			data.Parent = utils.ReferenceIDFromAPI(true, parent.GetId, data.Parent)

		} else {

			data.Parent = types.StringNull()

		}

	} else {

		data.Parent = types.StringNull()

	}

	// Handle optional string fields using helpers

	data.Description = utils.StringFromAPI(siteGroup.HasDescription(), siteGroup.GetDescription, data.Description)

	// Handle tags

	if siteGroup.HasTags() {

		tags := utils.NestedTagsToTagModels(siteGroup.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		if !tagDiags.HasError() {

			data.Tags = tagsValue

		}

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields - preserve state structure

	if siteGroup.HasCustomFields() && !data.CustomFields.IsNull() {

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		if !cfDiags.HasError() {

			customFields := utils.MapToCustomFieldModels(siteGroup.GetCustomFields(), stateCustomFields)

			customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

			if !cfValueDiags.HasError() {

				data.CustomFields = customFieldsValue

			}

		}

	} else if data.CustomFields.IsNull() {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
