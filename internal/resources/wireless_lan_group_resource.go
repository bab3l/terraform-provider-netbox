// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &WirelessLANGroupResource{}

	_ resource.ResourceWithConfigure = &WirelessLANGroupResource{}

	_ resource.ResourceWithImportState = &WirelessLANGroupResource{}
)

// NewWirelessLANGroupResource returns a new resource implementing the wireless LAN group resource.

func NewWirelessLANGroupResource() resource.Resource {
	return &WirelessLANGroupResource{}
}

// WirelessLANGroupResource defines the resource implementation.

type WirelessLANGroupResource struct {
	client *netbox.APIClient
}

// WirelessLANGroupResourceModel describes the resource data model.

type WirelessLANGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	Parent types.String `tfsdk:"parent"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *WirelessLANGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_lan_group"
}

// Schema defines the schema for the resource.

func (r *WirelessLANGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a wireless LAN group in NetBox. Wireless LAN groups organize wireless networks hierarchically.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the wireless LAN group.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the wireless LAN group.",

				Required: true,
			},

			"slug": schema.StringAttribute{
				MarkdownDescription: "A unique slug identifier for the wireless LAN group.",

				Required: true,
			},

			"parent": schema.StringAttribute{
				MarkdownDescription: "Parent wireless LAN group (ID or slug) for hierarchical organization.",

				Optional: true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("wireless LAN group"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

// Configure adds the provider configured client to the resource.

func (r *WirelessLANGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.

func (r *WirelessLANGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build request - use WritableWirelessLANGroupRequest

	apiReq := netbox.NewWritableWirelessLANGroupRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID, err := utils.ParseID(data.Parent.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(

				"Invalid Parent ID",

				fmt.Sprintf("Parent must be a numeric ID, got: %s", data.Parent.ValueString()),
			)

			return
		}

		apiReq.SetParent(parentID)
	}

	// Handle description and metadata fields
	utils.ApplyDescription(apiReq, data.Description)

	// Store plan values for filter-to-owned population later
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Apply tags and custom fields
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating wireless LAN group", map[string]interface{}{
		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsCreate(ctx).WritableWirelessLANGroupRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating wireless LAN group",

			utils.FormatAPIError(fmt.Sprintf("create wireless LAN group %s", data.Name.ValueString()), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case response.HasTags() && len(response.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)

	tflog.Trace(ctx, "Created wireless LAN group", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.

func (r *WirelessLANGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Wireless LAN Group ID",

			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", data.ID.ValueString()),
		)

		return
	}

	tflog.Debug(ctx, "Reading wireless LAN group", map[string]interface{}{
		"id": groupID,
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsRetrieve(ctx, groupID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading wireless LAN group",

			utils.FormatAPIError(fmt.Sprintf("read wireless LAN group ID %d", groupID), err, httpResp),
		)

		return
	}

	// Store state values for filter-to-owned (preserve null vs empty set distinction)
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only (preserves null/empty state)
	wasExplicitlyEmpty := !stateTags.IsNull() && !stateTags.IsUnknown() && len(stateTags.Elements()) == 0
	switch {
	case response.HasTags() && len(response.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, response.GetCustomFields(), &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.

func (r *WirelessLANGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Wireless LAN Group ID",

			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", plan.ID.ValueString()),
		)

		return
	}

	// Build request - use WritableWirelessLANGroupRequest

	apiReq := netbox.NewWritableWirelessLANGroupRequest(plan.Name.ValueString(), plan.Slug.ValueString())

	// Set optional fields

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.SetDescription(plan.Description.ValueString())
	}

	switch {
	case plan.Parent.IsUnknown():
		// Leave unchanged
	case plan.Parent.IsNull():
		// NetBox PATCH doesn't clear omitted optional fields; clear explicitly.
		apiReq.SetParentNil()
	default:
		parentID, err := utils.ParseID(plan.Parent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent must be a numeric ID, got: %s", plan.Parent.ValueString()),
			)
			return
		}
		apiReq.SetParent(parentID)
	}

	// Handle description and metadata fields
	utils.ApplyDescription(apiReq, plan.Description)

	// Apply tags and custom fields with merge-aware helpers
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating wireless LAN group", map[string]interface{}{
		"id": groupID,
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsUpdate(ctx, groupID).WritableWirelessLANGroupRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating wireless LAN group",

			utils.FormatAPIError(fmt.Sprintf("update wireless LAN group ID %d", groupID), err, httpResp),
		)

		return
	}

	// Map response to plan model

	// Save the plan's custom fields/tags before mapping (for filter-to-owned pattern)
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case response.HasTags() && len(response.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		plan.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		plan.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		plan.Tags = types.SetNull(types.StringType)
	}
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.

func (r *WirelessLANGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WirelessLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Wireless LAN Group ID",

			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", data.ID.ValueString()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting wireless LAN group", map[string]interface{}{
		"id": groupID,
	})

	httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsDestroy(ctx, groupID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting wireless LAN group",

			utils.FormatAPIError(fmt.Sprintf("delete wireless LAN group ID %d", groupID), err, httpResp),
		)

		return
	}
}

// ImportState imports an existing resource.

func (r *WirelessLANGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupID, err := utils.ParseID(req.ID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", req.ID),
		)

		return
	}

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLanGroupsRetrieve(ctx, groupID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error importing wireless LAN group",

			utils.FormatAPIError(fmt.Sprintf("import wireless LAN group ID %d", groupID), err, httpResp),
		)

		return
	}

	var data WirelessLANGroupResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *WirelessLANGroupResource) mapResponseToModel(ctx context.Context, group *netbox.WirelessLANGroup, data *WirelessLANGroupResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", group.GetId()))

	data.Name = types.StringValue(group.GetName())

	data.Slug = types.StringValue(group.GetSlug())

	// Map description

	if desc, ok := group.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map parent

	if group.Parent.IsSet() && group.Parent.Get() != nil {
		parent := group.Parent.Get()

		userParent := data.Parent.ValueString()

		if userParent == parent.GetName() || userParent == parent.GetSlug() || userParent == parent.GetDisplay() || userParent == fmt.Sprintf("%d", parent.GetId()) {
			// Keep user's original value
		} else {
			data.Parent = types.StringValue(parent.GetName())
		}
	} else {
		data.Parent = types.StringNull()
	}

	// Handle tags (slug list) with empty-set preservation
	wasExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
	switch {
	case group.HasTags() && len(group.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(group.GetTags()))
		for _, tag := range group.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, group.GetCustomFields(), diags)
}
