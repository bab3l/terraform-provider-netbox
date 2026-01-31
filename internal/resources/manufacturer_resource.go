package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
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

var (
	_ resource.Resource                = &ManufacturerResource{}
	_ resource.ResourceWithImportState = &ManufacturerResource{}
	_ resource.ResourceWithIdentity    = &ManufacturerResource{}
)

func NewManufacturerResource() resource.Resource {
	return &ManufacturerResource{}
}

type ManufacturerResource struct {
	client *netbox.APIClient
}

type ManufacturerResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *ManufacturerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manufacturer"
}

func (r *ManufacturerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a manufacturer in Netbox. Manufacturers are used to group devices and platforms by vendor.",
		Attributes: map[string]schema.Attribute{
			"id":   nbschema.IDAttribute("manufacturer"),
			"name": nbschema.NameAttribute("manufacturer", 100),
			"slug": nbschema.SlugAttribute("manufacturer"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("manufacturer"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ManufacturerResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Implement Configure, Create, Read, Update, Delete, ImportState methods here.
func (r *ManufacturerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ManufacturerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerRequest := netbox.ManufacturerRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&manufacturerRequest, data.Description)

	// Handle tags
	utils.ApplyTagsFromSlugs(ctx, r.client, &manufacturerRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle custom fields (no merge needed for Create)
	utils.ApplyCustomFields(ctx, &manufacturerRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersCreate(ctx).ManufacturerRequest(manufacturerRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating manufacturer", utils.FormatAPIError("create manufacturer", err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating manufacturer", fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusCreated, httpResp.StatusCode))
		return
	}
	if manufacturer == nil {
		resp.Diagnostics.AddError("Manufacturer API returned nil", "No manufacturer object returned from Netbox API.")

		return
	}

	// Map response to state using helpers
	r.mapManufacturerToState(ctx, manufacturer, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created manufacturer", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManufacturerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerID := data.ID.ValueString()
	var manufacturerIDInt int32
	manufacturerIDInt, err := utils.ParseID(manufacturerID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))
		return
	}
	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersRetrieve(ctx, manufacturerIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading manufacturer", utils.FormatAPIError(fmt.Sprintf("read manufacturer ID %s", manufacturerID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error reading manufacturer", fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusOK, httpResp.StatusCode))
		return
	}

	// Save original custom_fields state before mapping
	originalCustomFields := data.CustomFields

	// Map response to state using helpers
	r.mapManufacturerToState(ctx, manufacturer, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields state if it was null or empty
	// This prevents unmanaged/cleared fields from reappearing in state
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManufacturerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ManufacturerResourceModel
	var state ManufacturerResourceModel

	// Read both plan and state for merge-aware custom fields handling
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan as the data source
	data := plan

	manufacturerID := data.ID.ValueString()
	var manufacturerIDInt int32
	manufacturerIDInt, err := utils.ParseID(manufacturerID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))
		return
	}
	manufacturerRequest := netbox.ManufacturerRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&manufacturerRequest, data.Description)

	// Handle tags
	utils.ApplyTagsFromSlugs(ctx, r.client, &manufacturerRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle custom fields with merge-aware logic
	utils.ApplyCustomFieldsWithMerge(ctx, &manufacturerRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersUpdate(ctx, manufacturerIDInt).ManufacturerRequest(manufacturerRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating manufacturer", utils.FormatAPIError(fmt.Sprintf("update manufacturer ID %s", manufacturerID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating manufacturer", fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusOK, httpResp.StatusCode))
		return
	}

	// Map response to state using helpers
	r.mapManufacturerToState(ctx, manufacturer, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManufacturerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerID := data.ID.ValueString()
	var manufacturerIDInt int32
	manufacturerIDInt, err := utils.ParseID(manufacturerID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))
		return
	}
	httpResp, err := r.client.DcimAPI.DcimManufacturersDestroy(ctx, manufacturerIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return // Already deleted
		}
		resp.Diagnostics.AddError("Error deleting manufacturer", utils.FormatAPIError(fmt.Sprintf("delete manufacturer ID %s", manufacturerID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError("Error deleting manufacturer", fmt.Sprintf("Expected HTTP %d, got: %d", http.StatusNoContent, httpResp.StatusCode))
		return
	}
}

func (r *ManufacturerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		manufacturerIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", parsed.ID))
			return
		}

		manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersRetrieve(ctx, manufacturerIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing manufacturer", utils.FormatAPIError(fmt.Sprintf("import manufacturer ID %s", parsed.ID), err, httpResp))
			return
		}

		var data ManufacturerResourceModel
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, manufacturer.HasTags(), manufacturer.GetTags(), data.Tags)
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

		r.mapManufacturerToState(ctx, manufacturer, &data, &resp.Diagnostics)
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, manufacturer.GetCustomFields(), &resp.Diagnostics)
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

// mapManufacturerToState maps API response to Terraform state using state helpers.
func (r *ManufacturerResource) mapManufacturerToState(ctx context.Context, manufacturer *netbox.Manufacturer, data *ManufacturerResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))
	data.Name = types.StringValue(manufacturer.GetName())
	data.Slug = types.StringValue(manufacturer.GetSlug())
	data.Description = utils.StringFromAPI(manufacturer.HasDescription(), manufacturer.GetDescription, data.Description)

	// Handle tags using filter-to-owned approach
	planTags := data.Tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, manufacturer.HasTags(), manufacturer.GetTags(), planTags)

	// Handle custom fields - filter to owned fields only
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, manufacturer.GetCustomFields(), diags)
}
