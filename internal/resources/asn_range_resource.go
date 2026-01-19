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
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ASNRangeResource{}
	_ resource.ResourceWithConfigure   = &ASNRangeResource{}
	_ resource.ResourceWithImportState = &ASNRangeResource{}
	_ resource.ResourceWithIdentity    = &ASNRangeResource{}
)

// NewASNRangeResource returns a new ASNRange resource.
func NewASNRangeResource() resource.Resource {
	return &ASNRangeResource{}
}

// ASNRangeResource defines the resource implementation.
type ASNRangeResource struct {
	client *netbox.APIClient
}

// ASNRangeResourceModel describes the resource data model.
type ASNRangeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	RIR          types.String `tfsdk:"rir"`
	Start        types.String `tfsdk:"start"`
	End          types.String `tfsdk:"end"`
	Tenant       types.String `tfsdk:"tenant"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ASNRangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asn_range"
}

// Schema defines the schema for the resource.
func (r *ASNRangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an ASN Range in Netbox. ASN ranges define a contiguous range of Autonomous System Numbers that can be allocated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the ASN range.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": nbschema.NameAttribute("ASN range", 100),
			"slug": nbschema.SlugAttribute("ASN range"),
			"rir":  nbschema.RequiredReferenceAttribute("RIR", "ID or slug of the Regional Internet Registry (RIR) responsible for this ASN range. Required."),
			"start": schema.StringAttribute{
				MarkdownDescription: "The starting ASN in this range. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer",
					),
				},
			},
			"end": schema.StringAttribute{
				MarkdownDescription: "The ending ASN in this range. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer",
					),
				},
			},
			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this ASN range."),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("ASN range"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ASNRangeResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *ASNRangeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ASNRangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ASNRangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup RIR
	rirRef, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
	resp.Diagnostics.Append(rirDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse start and end to int64
	var start, end int64
	if _, err := fmt.Sscanf(data.Start.ValueString(), "%d", &start); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Start",
			fmt.Sprintf("Unable to parse start %q: %s", data.Start.ValueString(), err.Error()),
		)
		return
	}
	if _, err := fmt.Sscanf(data.End.ValueString(), "%d", &end); err != nil {
		resp.Diagnostics.AddError(
			"Invalid End",
			fmt.Sprintf("Unable to parse end %q: %s", data.End.ValueString(), err.Error()),
		)
		return
	}

	// Create the ASNRange request
	asnRangeRequest := netbox.NewASNRangeRequest(
		data.Name.ValueString(),
		data.Slug.ValueString(),
		*rirRef,
		start,
		end,
	)

	// Set optional fields (pass nil for state since this is Create)
	r.setOptionalFields(ctx, asnRangeRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating ASNRange", map[string]interface{}{
		"name":  data.Name.ValueString(),
		"start": start,
		"end":   end,
	})

	// Create the ASNRange
	asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesCreate(ctx).ASNRangeRequest(*asnRangeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ASNRange",
			utils.FormatAPIError("create ASNRange", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ASNRangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ASNRangeResourceModel

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
	tflog.Debug(ctx, "Reading ASNRange", map[string]interface{}{
		"id": id,
	})

	// Get the ASNRange from Netbox
	asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading ASNRange",
			utils.FormatAPIError(fmt.Sprintf("read ASNRange ID %d", id), err, httpResp),
		)
		return
	}

	// Save original custom_fields state before mapping
	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Read ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Preserve original custom_fields state if it was null or empty
	// This prevents unmanaged/cleared fields from reappearing in state
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ASNRangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ASNRangeResourceModel
	var state ASNRangeResourceModel

	// Read both plan and state for merge-aware custom fields handling
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan as the data source
	data := plan

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Lookup RIR
	rirRef, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
	resp.Diagnostics.Append(rirDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse start and end to int64
	var start, end int64
	if _, err := fmt.Sscanf(data.Start.ValueString(), "%d", &start); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Start",
			fmt.Sprintf("Unable to parse start %q: %s", data.Start.ValueString(), err.Error()),
		)
		return
	}
	if _, err := fmt.Sscanf(data.End.ValueString(), "%d", &end); err != nil {
		resp.Diagnostics.AddError(
			"Invalid End",
			fmt.Sprintf("Unable to parse end %q: %s", data.End.ValueString(), err.Error()),
		)
		return
	}

	// Create the ASNRange request
	asnRangeRequest := netbox.NewASNRangeRequest(
		data.Name.ValueString(),
		data.Slug.ValueString(),
		*rirRef,
		start,
		end,
	)

	// Set optional fields (pass state for merge-aware custom fields handling)
	r.setOptionalFields(ctx, asnRangeRequest, &data, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating ASNRange", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Update the ASNRange
	asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesUpdate(ctx, id).ASNRangeRequest(*asnRangeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ASNRange",
			utils.FormatAPIError(fmt.Sprintf("update ASNRange ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ASNRangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ASNRangeResourceModel

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
	tflog.Debug(ctx, "Deleting ASNRange", map[string]interface{}{
		"id": id,
	})

	// Delete the ASNRange
	httpResp, err := r.client.IpamAPI.IpamAsnRangesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting ASNRange",
			utils.FormatAPIError(fmt.Sprintf("delete ASNRange ID %d", id), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted ASNRange", map[string]interface{}{
		"id": id,
	})
}

func (r *ASNRangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", parsed.ID, err.Error()),
			)
			return
		}

		asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing ASNRange",
				utils.FormatAPIError(fmt.Sprintf("read ASNRange ID %d", id), err, httpResp),
			)
			return
		}

		var data ASNRangeResourceModel
		if asnRange.Rir.GetSlug() != "" {
			data.RIR = types.StringValue(asnRange.Rir.GetSlug())
		} else {
			data.RIR = types.StringValue(fmt.Sprintf("%d", asnRange.Rir.GetId()))
		}
		if asnRange.HasTenant() && asnRange.Tenant.Get() != nil {
			tenant := asnRange.Tenant.Get()
			if tenant.GetSlug() != "" {
				data.Tenant = types.StringValue(tenant.GetSlug())
			} else {
				data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
			}
		}
		if len(asnRange.GetTags()) > 0 {
			tagSlugs := make([]string, 0, len(asnRange.GetTags()))
			for _, tag := range asnRange.GetTags() {
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

		r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, asnRange.CustomFields, &resp.Diagnostics)
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the ASNRange request from the resource model.
func (r *ASNRangeResource) setOptionalFields(ctx context.Context, asnRangeRequest *netbox.ASNRangeRequest, data *ASNRangeResourceModel, state *ASNRangeResourceModel, diags *diag.Diagnostics) {
	// Tenant
	if utils.IsSet(data.Tenant) {
		tenantRef, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return
		}
		asnRangeRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	} else if data.Tenant.IsNull() {
		// Explicitly set to null to clear the field
		asnRangeRequest.SetTenantNil()
	}

	// Description
	utils.ApplyDescription(asnRangeRequest, data.Description)

	// Tags
	utils.ApplyTagsFromSlugs(ctx, r.client, asnRangeRequest, data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Custom fields with merge-aware handling
	if state != nil {
		// Update operation - merge custom fields to preserve unmanaged fields
		utils.ApplyCustomFieldsWithMerge(ctx, asnRangeRequest, data.CustomFields, state.CustomFields, diags)
	} else {
		// Create operation - apply custom fields directly
		utils.ApplyCustomFields(ctx, asnRangeRequest, data.CustomFields, diags)
	}
	if diags.HasError() {
		return
	}
}

// mapASNRangeToState maps a Netbox ASNRange to the Terraform state model.
func (r *ASNRangeResource) mapASNRangeToState(ctx context.Context, asnRange *netbox.ASNRange, data *ASNRangeResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", asnRange.Id))
	data.Name = types.StringValue(asnRange.Name)
	data.Slug = types.StringValue(asnRange.Slug)

	// RIR - preserve user's input format
	rir := asnRange.Rir
	data.RIR = utils.UpdateReferenceAttribute(data.RIR, rir.GetName(), rir.GetSlug(), rir.GetId())
	data.Start = types.StringValue(fmt.Sprintf("%d", asnRange.Start))
	data.End = types.StringValue(fmt.Sprintf("%d", asnRange.End))

	// Tenant - preserve user's input format
	if asnRange.HasTenant() && asnRange.Tenant.Get() != nil {
		tenant := asnRange.Tenant.Get()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Description
	if asnRange.Description != nil && *asnRange.Description != "" {
		data.Description = types.StringValue(*asnRange.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags
	var tagSlugs []string
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	case len(asnRange.GetTags()) > 0:
		for _, tag := range asnRange.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if diags.HasError() {
		return
	}

	// Custom Fields - filter to owned fields only
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, asnRange.CustomFields, diags)
}
