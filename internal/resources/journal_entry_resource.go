// Package resources contains Terraform resource implementations for the Netbox provider.

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &JournalEntryResource{}
	_ resource.ResourceWithImportState = &JournalEntryResource{}
	_ resource.ResourceWithIdentity    = &JournalEntryResource{}
)

func NewJournalEntryResource() resource.Resource {
	return &JournalEntryResource{}
}

// JournalEntryResource defines the resource implementation.
type JournalEntryResource struct {
	client *netbox.APIClient
}

// JournalEntryResourceModel describes the resource data model.
type JournalEntryResourceModel struct {
	ID                 types.Int32  `tfsdk:"id"`
	AssignedObjectType types.String `tfsdk:"assigned_object_type"`
	AssignedObjectID   types.Int64  `tfsdk:"assigned_object_id"`
	Kind               types.String `tfsdk:"kind"`
	Comments           types.String `tfsdk:"comments"`
	Tags               types.Set    `tfsdk:"tags"`
	CustomFields       types.Set    `tfsdk:"custom_fields"`
}

func (r *JournalEntryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_journal_entry"
}

func (r *JournalEntryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a journal entry in NetBox. Journal entries allow you to record notes, comments, and documentation against any object in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the journal entry.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"assigned_object_type": schema.StringAttribute{
				MarkdownDescription: "The content type of the assigned object (e.g., `dcim.device`, `dcim.site`, `ipam.ipaddress`).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"assigned_object_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the assigned object.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "The kind/severity of the journal entry. Valid values: `info`, `success`, `warning`, `danger`. Defaults to `info`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("info"),
			},
		},
	}

	// Add comments attribute
	maps.Copy(resp.Schema.Attributes, map[string]schema.Attribute{
		"comments": nbschema.CommentsAttribute("journal entry"),
	})

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *JournalEntryResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *JournalEntryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JournalEntryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JournalEntryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating Journal Entry", map[string]interface{}{
		"assigned_object_type": data.AssignedObjectType.ValueString(),
		"assigned_object_id":   data.AssignedObjectID.ValueInt64(),
	})

	// Prepare the Journal Entry request
	journalEntryRequest := netbox.WritableJournalEntryRequest{
		AssignedObjectType: data.AssignedObjectType.ValueString(),
		AssignedObjectId:   data.AssignedObjectID.ValueInt64(),
		Comments:           data.Comments.ValueString(),
	}

	// Set optional fields
	r.setOptionalFields(ctx, &journalEntryRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Journal Entry via API
	journalEntry, httpResp, err := r.client.ExtrasAPI.ExtrasJournalEntriesCreate(ctx).WritableJournalEntryRequest(journalEntryRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Journal Entry",
			utils.FormatAPIError("creating journal entry", err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "create journal entry", httpResp, http.StatusCreated) {
		return
	}
	tflog.Debug(ctx, "Created Journal Entry", map[string]interface{}{
		"id": journalEntry.GetId(),
	})

	// Map response back to state
	r.mapJournalEntryToState(ctx, journalEntry, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(fmt.Sprintf("%d", data.ID.ValueInt32())), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JournalEntryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JournalEntryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields from state for potential restoration
	originalCustomFields := data.CustomFields

	// Get the Journal Entry via API
	id := data.ID.ValueInt32()
	journalEntry, httpResp, err := r.client.ExtrasAPI.ExtrasJournalEntriesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() {
			tflog.Debug(ctx, "Journal Entry not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Journal Entry",
			utils.FormatAPIError("reading journal entry", err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read journal entry", httpResp, http.StatusOK) {
		return
	}

	// Map response to state
	r.mapJournalEntryToState(ctx, journalEntry, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was null or empty before (not managed or explicitly cleared),
	// restore that state after mapping.
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		tflog.Debug(ctx, "Custom fields unmanaged/cleared, preserving original state during Read", map[string]interface{}{
			"was_null":  originalCustomFields.IsNull(),
			"was_empty": !originalCustomFields.IsNull() && len(originalCustomFields.Elements()) == 0,
		})
		data.CustomFields = originalCustomFields
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(fmt.Sprintf("%d", data.ID.ValueInt32())), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JournalEntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read BOTH state and plan for merge-aware custom fields
	var state, plan JournalEntryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := plan.ID.ValueInt32()
	tflog.Debug(ctx, "Updating Journal Entry", map[string]interface{}{
		"id": id,
	})

	// Prepare the Journal Entry request
	journalEntryRequest := netbox.WritableJournalEntryRequest{
		AssignedObjectType: plan.AssignedObjectType.ValueString(),
		AssignedObjectId:   plan.AssignedObjectID.ValueInt64(),
		Comments:           plan.Comments.ValueString(),
	}

	// Set optional fields with state for merge
	r.setOptionalFields(ctx, &journalEntryRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Store plan tags/customfields for filter-to-owned population
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Update the Journal Entry via API
	journalEntry, httpResp, err := r.client.ExtrasAPI.ExtrasJournalEntriesUpdate(ctx, id).WritableJournalEntryRequest(journalEntryRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Journal Entry",
			utils.FormatAPIError("updating journal entry", err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "update journal entry", httpResp, http.StatusOK) {
		return
	}
	tflog.Debug(ctx, "Updated Journal Entry", map[string]interface{}{
		"id": journalEntry.GetId(),
	})

	// Map response back to state
	plan.Tags = planTags
	plan.CustomFields = planCustomFields
	r.mapJournalEntryToState(ctx, journalEntry, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(fmt.Sprintf("%d", plan.ID.ValueInt32())), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *JournalEntryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JournalEntryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueInt32()
	tflog.Debug(ctx, "Deleting Journal Entry", map[string]interface{}{
		"id": id,
	})

	// Delete the Journal Entry via API
	httpResp, err := r.client.ExtrasAPI.ExtrasJournalEntriesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, nil) {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Journal Entry",
			utils.FormatAPIError("deleting journal entry", err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "delete journal entry", httpResp, http.StatusNoContent) {
		return
	}
	tflog.Debug(ctx, "Deleted Journal Entry", map[string]interface{}{
		"id": id,
	})
}

func (r *JournalEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			resp.Diagnostics.AddError(
				"Error Importing Journal Entry",
				fmt.Sprintf("Could not parse ID %q: %s", parsed.ID, err),
			)
			return
		}

		journalEntry, httpResp, err := r.client.ExtrasAPI.ExtrasJournalEntriesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing journal entry", utils.FormatAPIError("reading journal entry", err, httpResp))
			return
		}
		if !utils.ValidateStatusCode(&resp.Diagnostics, "import journal entry", httpResp, http.StatusOK) {
			return
		}

		var data JournalEntryResourceModel
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

		r.mapJournalEntryToState(ctx, journalEntry, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, journalEntry.CustomFields, &resp.Diagnostics)
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

	id, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Journal Entry",
			fmt.Sprintf("Could not parse ID %q: %s", req.ID, err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// setOptionalFields sets optional fields on the Journal Entry request.
func (r *JournalEntryResource) setOptionalFields(ctx context.Context, journalEntryRequest *netbox.WritableJournalEntryRequest, plan *JournalEntryResourceModel, state *JournalEntryResourceModel, diags *diag.Diagnostics) {
	// Kind
	if utils.IsSet(plan.Kind) {
		kind, err := netbox.NewJournalEntryKindValueFromValue(plan.Kind.ValueString())
		if err != nil {
			diags.AddError("Invalid Kind", fmt.Sprintf("Invalid journal entry kind value: %s", plan.Kind.ValueString()))
			return
		}
		journalEntryRequest.Kind = kind
	}

	// Apply Tags and CustomFields with merge-aware pattern
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, journalEntryRequest, plan.Tags, diags)
	} else if state != nil && utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, journalEntryRequest, state.Tags, diags)
	}
	if diags.HasError() {
		return
	}
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, journalEntryRequest, plan.CustomFields, state.CustomFields, diags)
	} else {
		utils.ApplyCustomFields(ctx, journalEntryRequest, plan.CustomFields, diags)
	}
	if diags.HasError() {
		return
	}
}

// mapJournalEntryToState maps a Journal Entry API response to the Terraform state model.
func (r *JournalEntryResource) mapJournalEntryToState(ctx context.Context, journalEntry *netbox.JournalEntry, data *JournalEntryResourceModel, diags *diag.Diagnostics) {
	data.ID = types.Int32Value(journalEntry.GetId())
	data.AssignedObjectType = types.StringValue(journalEntry.GetAssignedObjectType())
	data.AssignedObjectID = types.Int64Value(journalEntry.GetAssignedObjectId())
	data.Comments = types.StringValue(journalEntry.GetComments())

	// Kind - always set since it's computed (defaults to \"info\")
	if journalEntry.Kind != nil && journalEntry.Kind.Value != nil {
		data.Kind = types.StringValue(string(*journalEntry.Kind.Value))
	} else {
		data.Kind = types.StringValue("info") // Default value
	}

	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, len(journalEntry.Tags) > 0, journalEntry.Tags, data.Tags)
	if diags.HasError() {
		return
	}

	// Handle custom fields using filter-to-owned helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, journalEntry.CustomFields, diags)
}
