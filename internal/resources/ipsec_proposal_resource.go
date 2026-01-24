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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &IPSecProposalResource{}
	_ resource.ResourceWithConfigure   = &IPSecProposalResource{}
	_ resource.ResourceWithImportState = &IPSecProposalResource{}
	_ resource.ResourceWithIdentity    = &IPSecProposalResource{}
)

// NewIPSecProposalResource returns a new IPSecProposal resource.
func NewIPSecProposalResource() resource.Resource {
	return &IPSecProposalResource{}
}

// IPSecProposalResource defines the resource implementation.
type IPSecProposalResource struct {
	client *netbox.APIClient
}

// IPSecProposalResourceModel describes the resource data model.
type IPSecProposalResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	EncryptionAlgorithm     types.String `tfsdk:"encryption_algorithm"`
	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`
	SALifetimeSeconds       types.Int64  `tfsdk:"sa_lifetime_seconds"`
	SALifetimeData          types.Int64  `tfsdk:"sa_lifetime_data"`
	Comments                types.String `tfsdk:"comments"`
	Tags                    types.Set    `tfsdk:"tags"`
	CustomFields            types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IPSecProposalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_proposal"
}

// Schema defines the schema for the resource.
func (r *IPSecProposalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IPSec Proposal in Netbox. IPSec proposals define the security parameters for the IPSec phase 2 (ESP/AH) negotiation in VPN connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IPSec proposal.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec proposal. Required.",
				Required:            true,
			},
			"description": nbschema.DescriptionAttribute("IPSec proposal"),
			"encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The encryption algorithm for the IPSec proposal. Optional. Valid values: `aes-128-cbc`, `aes-128-gcm`, `aes-192-cbc`, `aes-192-gcm`, `aes-256-cbc`, `aes-256-gcm`, `3des-cbc`, `des-cbc`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("aes-128-cbc", "aes-128-gcm", "aes-192-cbc", "aes-192-gcm", "aes-256-cbc", "aes-256-gcm", "3des-cbc", "des-cbc"),
				},
			},

			"authentication_algorithm": schema.StringAttribute{
				MarkdownDescription: "The authentication algorithm (hash) for the IPSec proposal. Optional. Valid values: `hmac-sha1`, `hmac-sha256`, `hmac-sha384`, `hmac-sha512`, `hmac-md5`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("hmac-sha1", "hmac-sha256", "hmac-sha384", "hmac-sha512", "hmac-md5"),
				},
			},
			"sa_lifetime_seconds": schema.Int64Attribute{
				MarkdownDescription: "Security association lifetime in seconds. Optional.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sa_lifetime_data": schema.Int64Attribute{
				MarkdownDescription: "Security association lifetime in kilobytes. Optional.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("IPSec proposal"))

	// Add tags and custom_fields
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *IPSecProposalResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *IPSecProposalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IPSecProposalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IPSecProposalResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IPSecProposal request
	ipsecRequest := netbox.NewWritableIPSecProposalRequest(data.Name.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, ipsecRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating IPSecProposal", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Create the IPSecProposal
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProposalsCreate(ctx).WritableIPSecProposalRequest(*ipsecRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IPSecProposal",
			utils.FormatAPIError("create IPSec proposal", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIPSecProposalToState(ctx, ipsec, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created IPSecProposal", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *IPSecProposalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IPSecProposalResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec proposal ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading IPSecProposal", map[string]interface{}{
		"id": id,
	})

	// Read the IPSecProposal
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProposalsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "IPSecProposal not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IPSecProposal",
			utils.FormatAPIError("read IPSec proposal", err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state

	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapIPSecProposalToState(ctx, ipsec, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// If custom_fields was null or empty before, restore that state

	// This prevents drift when config doesn't declare custom_fields

	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {

		data.CustomFields = originalCustomFields

	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *IPSecProposalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan IPSecProposalResourceModel

	// Read both state and plan for merge-aware custom fields handling
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec proposal ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Create the IPSecProposal request
	ipsecRequest := netbox.NewWritableIPSecProposalRequest(plan.Name.ValueString())

	// Set optional fields with state for merge-aware custom fields
	r.setOptionalFields(ctx, ipsecRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IPSecProposal", map[string]interface{}{
		"id":   id,
		"name": plan.Name.ValueString(),
	})

	// Update the IPSecProposal
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProposalsUpdate(ctx, id).WritableIPSecProposalRequest(*ipsecRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IPSecProposal",
			utils.FormatAPIError("update IPSec proposal", err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapIPSecProposalToState(ctx, ipsec, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, ipsec.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated IPSecProposal", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IPSecProposalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IPSecProposalResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec proposal ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting IPSecProposal", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Delete the IPSecProposal
	httpResp, err := r.client.VpnAPI.VpnIpsecProposalsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting IPSecProposal",
			utils.FormatAPIError("delete IPSec proposal", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted IPSecProposal", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state from an existing resource.
func (r *IPSecProposalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("IPSec proposal ID must be a number, got: %s", parsed.ID))
			return
		}

		ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProposalsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing IPSecProposal", utils.FormatAPIError("read IPSec proposal", err, httpResp))
			return
		}

		var data IPSecProposalResourceModel
		if ipsec.HasTags() {
			tagSlugs := make([]string, 0, len(ipsec.GetTags()))
			for _, tag := range ipsec.GetTags() {
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

		r.mapIPSecProposalToState(ctx, ipsec, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, ipsec.GetCustomFields(), &resp.Diagnostics)
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

// setOptionalFields sets optional fields on the WritableIPSecProposalRequest.
func (r *IPSecProposalResource) setOptionalFields(ctx context.Context, ipsecRequest *netbox.WritableIPSecProposalRequest, plan *IPSecProposalResourceModel, state *IPSecProposalResourceModel, diags *diag.Diagnostics) {
	// Encryption Algorithm
	if utils.IsSet(plan.EncryptionAlgorithm) {
		encAlg := netbox.Encryption(plan.EncryptionAlgorithm.ValueString())
		ipsecRequest.EncryptionAlgorithm = &encAlg
	}
	// Note: encryption_algorithm cannot be cleared once set in NetBox (sticky field)

	// Authentication Algorithm
	if utils.IsSet(plan.AuthenticationAlgorithm) {
		authAlg := netbox.Authentication(plan.AuthenticationAlgorithm.ValueString())
		ipsecRequest.AuthenticationAlgorithm = &authAlg
	}
	// Note: authentication_algorithm cannot be cleared once set in NetBox (sticky field)

	// SA Lifetime Seconds
	if utils.IsSet(plan.SALifetimeSeconds) {
		lifetime, err := utils.SafeInt32FromValue(plan.SALifetimeSeconds)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("SALifetimeSeconds value overflow: %s", err))
			return
		}
		ipsecRequest.SaLifetimeSeconds = *netbox.NewNullableInt32(&lifetime)
	} else if plan.SALifetimeSeconds.IsNull() && state != nil {
		// Explicitly clear by setting to NullableInt32 with nil value
		ipsecRequest.SaLifetimeSeconds = *netbox.NewNullableInt32(nil)
	}

	// SA Lifetime Data
	if utils.IsSet(plan.SALifetimeData) {
		lifetime, err := utils.SafeInt32FromValue(plan.SALifetimeData)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("SALifetimeData value overflow: %s", err))
			return
		}
		ipsecRequest.SaLifetimeData = *netbox.NewNullableInt32(&lifetime)
	} else if plan.SALifetimeData.IsNull() && state != nil {
		// Explicitly clear by setting to NullableInt32 with nil value
		ipsecRequest.SaLifetimeData = *netbox.NewNullableInt32(nil)
	}

	// Set description
	utils.ApplyDescription(ipsecRequest, plan.Description)

	// Set comments
	utils.ApplyComments(ipsecRequest, plan.Comments)

	// Apply tags
	utils.ApplyTagsFromSlugs(ctx, r.client, ipsecRequest, plan.Tags, diags)

	// Apply custom fields with merge logic to preserve unmanaged fields
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, ipsecRequest, plan.CustomFields, state.CustomFields, diags)
	} else {
		// During Create, no state exists yet
		utils.ApplyCustomFields(ctx, ipsecRequest, plan.CustomFields, diags)
	}
}

// mapIPSecProposalToState maps an IPSecProposal API response to the Terraform state model.
func (r *IPSecProposalResource) mapIPSecProposalToState(ctx context.Context, ipsec *netbox.IPSecProposal, data *IPSecProposalResourceModel, diags *diag.Diagnostics) {
	// ID
	data.ID = types.StringValue(fmt.Sprintf("%d", ipsec.Id))

	// Name
	data.Name = types.StringValue(ipsec.Name)

	// Description
	if ipsec.Description != nil && *ipsec.Description != "" {
		data.Description = types.StringValue(*ipsec.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Encryption Algorithm
	if ipsec.EncryptionAlgorithm != nil && ipsec.EncryptionAlgorithm.Value != nil {
		data.EncryptionAlgorithm = types.StringValue(string(*ipsec.EncryptionAlgorithm.Value))
	} else {
		data.EncryptionAlgorithm = types.StringNull()
	}

	// Authentication Algorithm
	if ipsec.AuthenticationAlgorithm != nil && ipsec.AuthenticationAlgorithm.Value != nil {
		data.AuthenticationAlgorithm = types.StringValue(string(*ipsec.AuthenticationAlgorithm.Value))
	} else {
		data.AuthenticationAlgorithm = types.StringNull()
	}

	// SA Lifetime Seconds
	if ipsec.SaLifetimeSeconds.IsSet() && ipsec.SaLifetimeSeconds.Get() != nil {
		data.SALifetimeSeconds = types.Int64Value(int64(*ipsec.SaLifetimeSeconds.Get()))
	} else {
		data.SALifetimeSeconds = types.Int64Null()
	}

	// SA Lifetime Data
	if ipsec.SaLifetimeData.IsSet() && ipsec.SaLifetimeData.Get() != nil {
		data.SALifetimeData = types.Int64Value(int64(*ipsec.SaLifetimeData.Get()))
	} else {
		data.SALifetimeData = types.Int64Null()
	}

	// Comments
	if ipsec.Comments != nil && *ipsec.Comments != "" {
		data.Comments = types.StringValue(*ipsec.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags with filter-to-owned pattern
	planTags := data.Tags
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case ipsec.HasTags() && len(ipsec.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(ipsec.GetTags()))
		for _, tag := range ipsec.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	if diags.HasError() {
		return
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, ipsec.HasCustomFields(), ipsec.GetCustomFields(), data.CustomFields, diags)
}
