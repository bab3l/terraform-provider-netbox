// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
var _ resource.Resource = &JournalEntryResource{}
var _ resource.ResourceWithImportState = &JournalEntryResource{}

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

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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
	r.setOptionalFields(ctx, &journalEntryRequest, &data, &resp.Diagnostics)
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
	tflog.Debug(ctx, "Created Journal Entry", map[string]interface{}{
		"id": journalEntry.GetId(),
	})

	// Map response back to state
	r.mapJournalEntryToState(ctx, journalEntry, &data, &resp.Diagnostics)
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

	// Get the Journal Entry via API
	id := data.ID.ValueInt32()
	journalEntry, httpResp, err := r.client.ExtrasAPI.ExtrasJournalEntriesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Journal Entry not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Journal Entry",
			utils.FormatAPIError("reading journal entry", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapJournalEntryToState(ctx, journalEntry, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JournalEntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JournalEntryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueInt32()
	tflog.Debug(ctx, "Updating Journal Entry", map[string]interface{}{
		"id": id,
	})

	// Prepare the Journal Entry request
	journalEntryRequest := netbox.WritableJournalEntryRequest{
		AssignedObjectType: data.AssignedObjectType.ValueString(),
		AssignedObjectId:   data.AssignedObjectID.ValueInt64(),
		Comments:           data.Comments.ValueString(),
	}

	// Set optional fields
	r.setOptionalFields(ctx, &journalEntryRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

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
	tflog.Debug(ctx, "Updated Journal Entry", map[string]interface{}{
		"id": journalEntry.GetId(),
	})

	// Map response back to state
	r.mapJournalEntryToState(ctx, journalEntry, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Journal Entry",
			utils.FormatAPIError("deleting journal entry", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted Journal Entry", map[string]interface{}{
		"id": id,
	})
}

func (r *JournalEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *JournalEntryResource) setOptionalFields(ctx context.Context, journalEntryRequest *netbox.WritableJournalEntryRequest, data *JournalEntryResourceModel, diags *diag.Diagnostics) {
	// Kind
	if utils.IsSet(data.Kind) {
		kind, err := netbox.NewJournalEntryKindValueFromValue(data.Kind.ValueString())
		if err != nil {
			diags.AddError("Invalid Kind", fmt.Sprintf("Invalid journal entry kind value: %s", data.Kind.ValueString()))
			return
		}
		journalEntryRequest.Kind = kind
	}

	// Apply Tags and CustomFields
	utils.ApplyTags(ctx, journalEntryRequest, data.Tags, diags)
	if diags.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, journalEntryRequest, data.CustomFields, diags)
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

	// Tags
	if len(journalEntry.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(journalEntry.Tags)
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields
	if journalEntry.CustomFields != nil && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		customFields := utils.MapToCustomFieldModels(journalEntry.CustomFields, stateCustomFields)
		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfValueDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	}
}
