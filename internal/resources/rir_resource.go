// Package resources contains Terraform resource implementations for the Netbox provider.

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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &RIRResource{}
	_ resource.ResourceWithConfigure   = &RIRResource{}
	_ resource.ResourceWithImportState = &RIRResource{}
	_ resource.ResourceWithIdentity    = &RIRResource{}
)

// NewRIRResource returns a new RIR resource.
func NewRIRResource() resource.Resource {
	return &RIRResource{}
}

// RIRResource defines the resource implementation.
type RIRResource struct {
	client *netbox.APIClient
}

// RIRResourceModel describes the resource data model.
type RIRResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	IsPrivate    types.Bool   `tfsdk:"is_private"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *RIRResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rir"
}

// Schema defines the schema for the resource.
func (r *RIRResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Regional Internet Registry (RIR) in Netbox. RIRs are organizations that manage the allocation and registration of Internet number resources (IP addresses, ASNs).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the RIR.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": nbschema.NameAttribute("RIR", 100),
			"slug": nbschema.SlugAttribute("RIR"),
			"is_private": schema.BoolAttribute{
				MarkdownDescription: "Whether IP space managed by this RIR is considered private. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("RIR"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *RIRResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *RIRResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *RIRResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RIRResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the RIR request
	rirRequest := netbox.NewRIRRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, rirRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating RIR", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Create the RIR
	rir, httpResp, err := r.client.IpamAPI.IpamRirsCreate(ctx).RIRRequest(*rirRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating RIR",
			utils.FormatAPIError("create RIR", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *RIRResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RIRResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Reading RIR", map[string]interface{}{
		"id": id,
	})

	// Get the RIR from Netbox
	rir, httpResp, err := r.client.IpamAPI.IpamRirsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading RIR",
			utils.FormatAPIError(fmt.Sprintf("read RIR ID %d", id), err, httpResp),
		)
		return
	}

	// Preserve the custom_fields plan/state if it's null or empty
	var planSet types.Set
	if data.CustomFields.IsNull() || len(data.CustomFields.Elements()) == 0 {
		planSet = data.CustomFields
	}

	// Map response to model
	r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)

	// Restore null/empty custom_fields if it was null/empty before
	if !planSet.IsNull() || (planSet.IsNull() && data.CustomFields.IsNull()) {
		data.CustomFields = planSet
	}

	tflog.Debug(ctx, "Read RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *RIRResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RIRResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state for merge-aware custom fields
	var state RIRResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Create the RIR request
	rirRequest := netbox.NewRIRRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields with state for merge-aware custom fields
	r.setOptionalFields(ctx, rirRequest, &data, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating RIR", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Update the RIR
	rir, httpResp, err := r.client.IpamAPI.IpamRirsUpdate(ctx, id).RIRRequest(*rirRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating RIR",
			utils.FormatAPIError(fmt.Sprintf("update RIR ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *RIRResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RIRResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting RIR", map[string]interface{}{
		"id": id,
	})

	// Delete the RIR
	httpResp, err := r.client.IpamAPI.IpamRirsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting RIR",
			utils.FormatAPIError(fmt.Sprintf("delete RIR ID %d", id), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted RIR", map[string]interface{}{
		"id": id,
	})
}

func (r *RIRResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		id, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to parse ID %q: %s", parsed.ID, err.Error()))
			return
		}

		rir, httpResp, err := r.client.IpamAPI.IpamRirsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing RIR", utils.FormatAPIError(fmt.Sprintf("read RIR ID %d", id), err, httpResp))
			return
		}

		var data RIRResourceModel
		if len(rir.Tags) > 0 {
			tagSlugs := make([]string, 0, len(rir.Tags))
			for _, tag := range rir.Tags {
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

		r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, rir.CustomFields, &resp.Diagnostics)
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

// setOptionalFields sets optional fields on the RIR request from the resource model.
// state is optional and only provided during updates for merge-aware custom fields.
func (r *RIRResource) setOptionalFields(ctx context.Context, rirRequest *netbox.RIRRequest, data *RIRResourceModel, state *RIRResourceModel, diags *diag.Diagnostics) {
	// Is Private
	if utils.IsSet(data.IsPrivate) {
		isPrivate := data.IsPrivate.ValueBool()
		rirRequest.IsPrivate = &isPrivate
	}

	// Apply description
	utils.ApplyDescription(rirRequest, data.Description)

	// Apply tags
	utils.ApplyTagsFromSlugs(ctx, r.client, rirRequest, data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Apply custom fields with merge awareness
	if state != nil {
		// Update: use merge-aware helper
		utils.ApplyCustomFieldsWithMerge(ctx, rirRequest, data.CustomFields, state.CustomFields, diags)
	} else {
		// Create: apply custom fields directly
		utils.ApplyCustomFields(ctx, rirRequest, data.CustomFields, diags)
	}
	if diags.HasError() {
		return
	}
}

// mapRIRToState maps a Netbox RIR to the Terraform state model.
func (r *RIRResource) mapRIRToState(ctx context.Context, rir *netbox.RIR, data *RIRResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rir.Id))
	data.Name = types.StringValue(rir.Name)
	data.Slug = types.StringValue(rir.Slug)

	// Is Private
	if rir.IsPrivate != nil {
		data.IsPrivate = types.BoolValue(*rir.IsPrivate)
	} else {
		data.IsPrivate = types.BoolValue(false)
	}

	// Description
	if rir.Description != nil && *rir.Description != "" {
		data.Description = types.StringValue(*rir.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags - filter to owned slugs only
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	case len(rir.Tags) > 0:
		var tagSlugs []string
		for _, tag := range rir.Tags {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	}
	// Custom Fields - filter to owned fields only
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, rir.CustomFields, diags)
}
