// Package resources contains Terraform resource implementations for the Netbox provider.

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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AggregateResource{}
	_ resource.ResourceWithConfigure   = &AggregateResource{}
	_ resource.ResourceWithImportState = &AggregateResource{}
	_ resource.ResourceWithIdentity    = &AggregateResource{}
)

// NewAggregateResource returns a new Aggregate resource.
func NewAggregateResource() resource.Resource {
	return &AggregateResource{}
}

// AggregateResource defines the resource implementation.
type AggregateResource struct {
	client *netbox.APIClient
}

// AggregateResourceModel describes the resource data model.
type AggregateResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Prefix       types.String `tfsdk:"prefix"`
	RIR          types.String `tfsdk:"rir"`
	Tenant       types.String `tfsdk:"tenant"`
	DateAdded    types.String `tfsdk:"date_added"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *AggregateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate"
}

// Schema defines the schema for the resource.
func (r *AggregateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an aggregate in Netbox. Aggregates are top-level IP address blocks that represent the entire address space available for allocation by an organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the aggregate.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prefix": nbschema.PrefixAttribute("The IP prefix in CIDR notation (e.g., 10.0.0.0/8, 2001:db8::/32)."),
			"rir":    nbschema.RequiredReferenceAttributeWithDiffSuppress("RIR", "ID, name, or slug of the Regional Internet Registry (RIR) this aggregate belongs to. Required."),

			"tenant": nbschema.ReferenceAttributeWithDiffSuppress("tenant", "ID or slug of the tenant this aggregate is assigned to."),
			"date_added": schema.StringAttribute{
				MarkdownDescription: "The date this aggregate was added (YYYY-MM-DD format).",
				Optional:            true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("aggregate"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *AggregateResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure sets the client for the resource.
func (r *AggregateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new aggregate resource.
func (r *AggregateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AggregateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request (pass nil state since this is a new resource)
	createReq, diags := r.buildCreateRequest(ctx, &data, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating aggregate", map[string]interface{}{
		"prefix": data.Prefix.ValueString(),
	})

	// Call API to create aggregate
	aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesCreate(ctx).WritableAggregateRequest(*createReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating aggregate",
			fmt.Sprintf("Could not create aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, aggregate, &data)
	tflog.Debug(ctx, "Created aggregate", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the aggregate resource.
func (r *AggregateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AggregateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Reading aggregate", map[string]interface{}{
		"id": id,
	})

	// Call API to read aggregate
	aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if utils.HandleNotFound(httpResp, func() {
			tflog.Debug(ctx, "Aggregate not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error reading aggregate",
			fmt.Sprintf("Could not read aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read aggregate", httpResp, http.StatusOK) {
		return
	}

	// Preserve the custom_fields plan/state if it's null or empty
	var planSet types.Set
	if data.CustomFields.IsNull() || len(data.CustomFields.Elements()) == 0 {
		planSet = data.CustomFields
	}

	// Map response to model
	r.mapResponseToModel(ctx, aggregate, &data)

	// Restore null/empty custom_fields if it was null/empty before
	if !planSet.IsNull() || (planSet.IsNull() && data.CustomFields.IsNull()) {
		data.CustomFields = planSet
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the aggregate resource.
func (r *AggregateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AggregateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state for merge-aware custom fields
	var state AggregateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}

	// Build the update request with state for merge-aware custom fields
	updateReq, diags := r.buildCreateRequest(ctx, &data, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating aggregate", map[string]interface{}{
		"id": id,
	})

	// Call API to update aggregate
	aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesUpdate(ctx, id).WritableAggregateRequest(*updateReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating aggregate",
			fmt.Sprintf("Could not update aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, aggregate, &data)
	tflog.Debug(ctx, "Updated aggregate", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the aggregate resource.
func (r *AggregateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AggregateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting aggregate", map[string]interface{}{
		"id": id,
	})

	// Call API to delete aggregate
	httpResp, err := r.client.IpamAPI.IpamAggregatesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() {
			tflog.Debug(ctx, "Aggregate already deleted", map[string]interface{}{
				"id": id,
			})
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting aggregate",
			fmt.Sprintf("Could not delete aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "delete aggregate", httpResp, http.StatusNoContent) {
		return
	}
	tflog.Debug(ctx, "Deleted aggregate", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing aggregate.
func (r *AggregateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
				"Invalid aggregate ID",
				fmt.Sprintf("Aggregate ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing aggregate",
				fmt.Sprintf("Could not read aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
			)
			return
		}

		var data AggregateResourceModel
		if rir := aggregate.GetRir(); rir.Id != 0 {
			data.RIR = types.StringValue(fmt.Sprintf("%d", rir.Id))
		} else {
			data.RIR = types.StringNull()
		}
		if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
			data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		} else {
			data.Tenant = types.StringNull()
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, aggregate.HasTags(), aggregate.GetTags(), data.Tags)
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

		r.mapResponseToModel(ctx, aggregate, &data)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, aggregate.CustomFields, &resp.Diagnostics)
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

// buildCreateRequest builds a WritableAggregateRequest from the model.
// state is optional and only provided during updates for merge-aware custom fields.
func (r *AggregateResource) buildCreateRequest(ctx context.Context, data *AggregateResourceModel, state *AggregateResourceModel) (*netbox.WritableAggregateRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Look up RIR (required)
	rir, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
	diags.Append(rirDiags...)
	if diags.HasError() {
		return nil, diags
	}
	createReq := netbox.NewWritableAggregateRequest(data.Prefix.ValueString(), *rir)

	// Handle tenant (optional)
	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return nil, diags
		}
		createReq.SetTenant(*tenant)
	} else if data.Tenant.IsNull() {
		createReq.SetTenantNil()
	}

	// Handle date_added (optional)
	if !data.DateAdded.IsNull() && !data.DateAdded.IsUnknown() {
		createReq.SetDateAdded(data.DateAdded.ValueString())
	} else if data.DateAdded.IsNull() {
		createReq.SetDateAddedNil()
	}

	// Apply description and comments
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		createReq.SetDescription(desc)
	} else if data.Description.IsNull() {
		createReq.SetDescription("")
	}
	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		createReq.SetComments(comments)
	} else if data.Comments.IsNull() {
		createReq.SetComments("")
	}

	// Apply tags (slug list)
	utils.ApplyTagsFromSlugs(ctx, r.client, createReq, data.Tags, &diags)
	if diags.HasError() {
		return nil, diags
	}

	// Apply custom fields with merge awareness
	if state != nil {
		// Update: use merge-aware helper
		utils.ApplyCustomFieldsWithMerge(ctx, createReq, data.CustomFields, state.CustomFields, &diags)
	} else {
		// Create: apply custom fields directly
		utils.ApplyCustomFields(ctx, createReq, data.CustomFields, &diags)
	}
	if diags.HasError() {
		return nil, diags
	}

	return createReq, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *AggregateResource) mapResponseToModel(ctx context.Context, aggregate *netbox.Aggregate, data *AggregateResourceModel) {
	var diags diag.Diagnostics

	data.ID = types.StringValue(fmt.Sprintf("%d", aggregate.GetId()))
	data.Prefix = types.StringValue(aggregate.GetPrefix())

	// Map RIR (store ID to avoid import drift)
	if rir := aggregate.GetRir(); rir.Id != 0 {
		data.RIR = types.StringValue(fmt.Sprintf("%d", rir.Id))
	} else {
		data.RIR = types.StringNull()
	}

	// Map tenant (store ID to avoid import drift)
	if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
	}

	// Map date_added
	if dateAdded := aggregate.GetDateAdded(); dateAdded != "" {
		data.DateAdded = types.StringValue(dateAdded)
	} else {
		data.DateAdded = types.StringNull()
	}

	// Map description
	if description, ok := aggregate.GetDescriptionOk(); ok && description != nil && *description != "" {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := aggregate.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags (slug list)
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, aggregate.HasTags(), aggregate.GetTags(), data.Tags)
	if diags.HasError() {
		return
	}

	// Custom Fields - filter to owned fields only
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, aggregate.CustomFields, &diags)
}
