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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                = &IKEProposalResource{}
	_ resource.ResourceWithConfigure   = &IKEProposalResource{}
	_ resource.ResourceWithImportState = &IKEProposalResource{}
	_ resource.ResourceWithIdentity    = &IKEProposalResource{}
)

// NewIKEProposalResource returns a new IKEProposal resource.
func NewIKEProposalResource() resource.Resource {
	return &IKEProposalResource{}
}

// IKEProposalResource defines the resource implementation.
type IKEProposalResource struct {
	client *netbox.APIClient
}

// IKEProposalResourceModel describes the resource data model.
type IKEProposalResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	AuthenticationMethod    types.String `tfsdk:"authentication_method"`
	EncryptionAlgorithm     types.String `tfsdk:"encryption_algorithm"`
	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`
	Group                   types.Int64  `tfsdk:"group"`
	SALifetime              types.Int64  `tfsdk:"sa_lifetime"`
	Comments                types.String `tfsdk:"comments"`
	Tags                    types.Set    `tfsdk:"tags"`
	CustomFields            types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IKEProposalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_proposal"
}

// Schema defines the schema for the resource.
func (r *IKEProposalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IKE (Internet Key Exchange) Proposal in Netbox. IKE proposals define the security parameters for the IKE phase 1 negotiation in IPSec VPN connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IKE proposal.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IKE proposal. Required.",
				Required:            true,
			},
			"description": nbschema.DescriptionAttribute("IKE proposal"),
			"authentication_method": schema.StringAttribute{
				MarkdownDescription: "The authentication method for the IKE proposal. Required. Valid values: `preshared-keys`, `certificates`, `rsa-signatures`, `dsa-signatures`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("preshared-keys", "certificates", "rsa-signatures", "dsa-signatures"),
				},
			},
			"encryption_algorithm": schema.StringAttribute{
				MarkdownDescription: "The encryption algorithm for the IKE proposal. Required. Valid values: `aes-128-cbc`, `aes-128-gcm`, `aes-192-cbc`, `aes-192-gcm`, `aes-256-cbc`, `aes-256-gcm`, `3des-cbc`, `des-cbc`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("aes-128-cbc", "aes-128-gcm", "aes-192-cbc", "aes-192-gcm", "aes-256-cbc", "aes-256-gcm", "3des-cbc", "des-cbc"),
				},
			},
			"authentication_algorithm": schema.StringAttribute{
				MarkdownDescription: "The authentication algorithm (hash) for the IKE proposal. Optional. Valid values: `hmac-sha1`, `hmac-sha256`, `hmac-sha384`, `hmac-sha512`, `hmac-md5`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("hmac-sha1", "hmac-sha256", "hmac-sha384", "hmac-sha512", "hmac-md5"),
				},
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The Diffie-Hellman group for the IKE proposal. Required. Valid values: 1, 2, 5, 14-34.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2, 5, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34),
				},
			},
			"sa_lifetime": schema.Int64Attribute{
				MarkdownDescription: "Security association lifetime in seconds. Optional.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("IKE proposal"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *IKEProposalResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *IKEProposalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IKEProposalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IKEProposalResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IKEProposal request
	authMethod := netbox.IKEProposalAuthenticationMethodValue(data.AuthenticationMethod.ValueString())
	encAlg := netbox.IKEProposalEncryptionAlgorithmValue(data.EncryptionAlgorithm.ValueString())
	groupVal, err := utils.SafeInt32FromValue(data.Group)
	if err != nil {
		resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Group value overflow: %s", err))
		return
	}
	group := netbox.PatchedWritableIKEProposalRequestGroup(groupVal)
	ikeRequest := netbox.NewWritableIKEProposalRequest(
		data.Name.ValueString(),
		authMethod,
		encAlg,
		group,
	)

	// Set optional fields
	r.setOptionalFields(ctx, ikeRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating IKEProposal", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Create the IKEProposal
	ike, httpResp, err := r.client.VpnAPI.VpnIkeProposalsCreate(ctx).WritableIKEProposalRequest(*ikeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IKEProposal",
			utils.FormatAPIError("create IKE proposal", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIKEProposalToState(ctx, ike, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created IKEProposal", map[string]interface{}{
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
func (r *IKEProposalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IKEProposalResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IKE proposal ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading IKEProposal", map[string]interface{}{
		"id": id,
	})

	// Read the IKEProposal
	ike, httpResp, err := r.client.VpnAPI.VpnIkeProposalsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "IKEProposal not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IKEProposal",
			utils.FormatAPIError("read IKE proposal", err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state

	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapIKEProposalToState(ctx, ike, &data, &resp.Diagnostics)

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
func (r *IKEProposalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan IKEProposalResourceModel

	// Read current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IKE proposal ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Create the IKEProposal request
	authMethod := netbox.IKEProposalAuthenticationMethodValue(plan.AuthenticationMethod.ValueString())
	encAlg := netbox.IKEProposalEncryptionAlgorithmValue(plan.EncryptionAlgorithm.ValueString())
	groupVal, err := utils.SafeInt32FromValue(plan.Group)
	if err != nil {
		resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Group value overflow: %s", err))
		return
	}
	group := netbox.PatchedWritableIKEProposalRequestGroup(groupVal)
	ikeRequest := netbox.NewWritableIKEProposalRequest(
		plan.Name.ValueString(),
		authMethod,
		encAlg,
		group,
	)

	// Set optional fields
	r.setOptionalFields(ctx, ikeRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IKEProposal", map[string]interface{}{
		"id":   id,
		"name": plan.Name.ValueString(),
	})

	// Update the IKEProposal
	ike, httpResp, err := r.client.VpnAPI.VpnIkeProposalsUpdate(ctx, id).WritableIKEProposalRequest(*ikeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IKEProposal",
			utils.FormatAPIError("update IKE proposal", err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapIKEProposalToState(ctx, ike, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, ike.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated IKEProposal", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IKEProposalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IKEProposalResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IKE proposal ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting IKEProposal", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Delete the IKEProposal
	httpResp, err := r.client.VpnAPI.VpnIkeProposalsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting IKEProposal",
			utils.FormatAPIError("delete IKE proposal", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted IKEProposal", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state from an existing resource.
func (r *IKEProposalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("IKE proposal ID must be a number, got: %s", parsed.ID))
			return
		}

		ike, httpResp, err := r.client.VpnAPI.VpnIkeProposalsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing IKEProposal", utils.FormatAPIError("read IKE proposal", err, httpResp))
			return
		}

		var data IKEProposalResourceModel
		if ike.HasTags() {
			tagSlugs := make([]string, 0, len(ike.GetTags()))
			for _, tag := range ike.GetTags() {
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

		r.mapIKEProposalToState(ctx, ike, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, ike.GetCustomFields(), &resp.Diagnostics)
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

// setOptionalFields sets optional fields on the WritableIKEProposalRequest.
func (r *IKEProposalResource) setOptionalFields(ctx context.Context, ikeRequest *netbox.WritableIKEProposalRequest, plan *IKEProposalResourceModel, state *IKEProposalResourceModel, diags *diag.Diagnostics) {
	// Set description
	utils.ApplyDescription(ikeRequest, plan.Description)

	// Authentication Algorithm
	if utils.IsSet(plan.AuthenticationAlgorithm) {
		authAlg := netbox.PatchedWritableIKEProposalRequestAuthenticationAlgorithm(plan.AuthenticationAlgorithm.ValueString())
		ikeRequest.AuthenticationAlgorithm = &authAlg
	}

	// SA Lifetime
	if utils.IsSet(plan.SALifetime) {
		lifetime, err := utils.SafeInt32FromValue(plan.SALifetime)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("SALifetime value overflow: %s", err))
			return
		}
		ikeRequest.SaLifetime = *netbox.NewNullableInt32(&lifetime)
	} else if state != nil && utils.IsSet(state.SALifetime) {
		// Explicitly clear removed optional nullable int field (NetBox PATCH semantics).
		ikeRequest.SaLifetime = *netbox.NewNullableInt32(nil)
	}

	// Set comments, tags, and custom fields with merge-aware helpers
	utils.ApplyComments(ikeRequest, plan.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, ikeRequest, plan.Tags, diags)

	// Apply custom fields with merge logic to preserve unmanaged fields
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, ikeRequest, plan.CustomFields, state.CustomFields, diags)
	} else {
		// During Create, no state exists yet
		utils.ApplyCustomFields(ctx, ikeRequest, plan.CustomFields, diags)
	}
	if diags.HasError() {
		return
	}
}

// mapIKEProposalToState maps an IKEProposal API response to the Terraform state model.
func (r *IKEProposalResource) mapIKEProposalToState(ctx context.Context, ike *netbox.IKEProposal, data *IKEProposalResourceModel, diags *diag.Diagnostics) {
	// ID
	data.ID = types.StringValue(fmt.Sprintf("%d", ike.Id))

	// Name
	data.Name = types.StringValue(ike.Name)

	// Description
	if ike.Description != nil && *ike.Description != "" {
		data.Description = types.StringValue(*ike.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Authentication Method
	if ike.AuthenticationMethod.Value != nil {
		data.AuthenticationMethod = types.StringValue(string(*ike.AuthenticationMethod.Value))
	}

	// Encryption Algorithm
	if ike.EncryptionAlgorithm.Value != nil {
		data.EncryptionAlgorithm = types.StringValue(string(*ike.EncryptionAlgorithm.Value))
	}

	// Authentication Algorithm
	if ike.AuthenticationAlgorithm != nil && ike.AuthenticationAlgorithm.Value != nil && string(*ike.AuthenticationAlgorithm.Value) != "" {
		data.AuthenticationAlgorithm = types.StringValue(string(*ike.AuthenticationAlgorithm.Value))
	} else {
		data.AuthenticationAlgorithm = types.StringNull()
	}

	// Group
	if ike.Group.Value != nil {
		data.Group = types.Int64Value(int64(*ike.Group.Value))
	}

	// SA Lifetime
	if ike.SaLifetime.IsSet() && ike.SaLifetime.Get() != nil {
		data.SALifetime = types.Int64Value(int64(*ike.SaLifetime.Get()))
	} else {
		data.SALifetime = types.Int64Null()
	}

	// Comments
	if ike.Comments != nil && *ike.Comments != "" {
		data.Comments = types.StringValue(*ike.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags with filter-to-owned pattern
	planTags := data.Tags
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case ike.HasTags() && len(ike.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(ike.GetTags()))
		for _, tag := range ike.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, ike.HasCustomFields(), ike.GetCustomFields(), data.CustomFields, diags)
}
