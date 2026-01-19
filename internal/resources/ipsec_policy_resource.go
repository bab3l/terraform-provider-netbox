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
	_ resource.Resource                = &IPSecPolicyResource{}
	_ resource.ResourceWithConfigure   = &IPSecPolicyResource{}
	_ resource.ResourceWithImportState = &IPSecPolicyResource{}
	_ resource.ResourceWithIdentity    = &IPSecPolicyResource{}
)

// NewIPSecPolicyResource returns a new IPSecPolicy resource.
func NewIPSecPolicyResource() resource.Resource {
	return &IPSecPolicyResource{}
}

// IPSecPolicyResource defines the resource implementation.
type IPSecPolicyResource struct {
	client *netbox.APIClient
}

// IPSecPolicyResourceModel describes the resource data model.
type IPSecPolicyResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Proposals    types.Set    `tfsdk:"proposals"`
	PFSGroup     types.Int64  `tfsdk:"pfs_group"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IPSecPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_policy"
}

// Schema defines the schema for the resource.
func (r *IPSecPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IPSec Policy in Netbox. IPSec policies group together IPSec proposals and define the PFS (Perfect Forward Secrecy) group for VPN connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IPSec policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec policy. Required.",
				Required:            true,
			},
			"description": nbschema.DescriptionAttribute("IPSec policy"),
			"proposals": schema.SetAttribute{
				MarkdownDescription: "A set of IPSec proposal IDs to associate with this policy.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"pfs_group": schema.Int64Attribute{
				MarkdownDescription: "The Diffie-Hellman group for Perfect Forward Secrecy. Optional. Valid values: 1, 2, 5, 14-34.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2, 5, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34),
				},
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("IPSec policy"))

	// Add tags and custom_fields
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *IPSecPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *IPSecPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IPSecPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IPSecPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IPSecPolicy request
	ipsecRequest := netbox.NewWritableIPSecPolicyRequest(data.Name.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, ipsecRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating IPSecPolicy", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Create the IPSecPolicy
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecPoliciesCreate(ctx).WritableIPSecPolicyRequest(*ipsecRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IPSecPolicy",
			utils.FormatAPIError("create IPSec policy", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIPSecPolicyToState(ctx, ipsec, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created IPSecPolicy", map[string]interface{}{
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
func (r *IPSecPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IPSecPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec policy ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading IPSecPolicy", map[string]interface{}{
		"id": id,
	})

	// Read the IPSecPolicy
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecPoliciesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "IPSecPolicy not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IPSecPolicy",
			utils.FormatAPIError("read IPSec policy", err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state

	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapIPSecPolicyToState(ctx, ipsec, &data, &resp.Diagnostics)

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
func (r *IPSecPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan IPSecPolicyResourceModel

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
			fmt.Sprintf("Could not parse IPSec policy ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Create the IPSecPolicy request
	ipsecRequest := netbox.NewWritableIPSecPolicyRequest(plan.Name.ValueString())

	// Set optional fields with state for merge-aware custom fields
	r.setOptionalFields(ctx, ipsecRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IPSecPolicy", map[string]interface{}{
		"id":   id,
		"name": plan.Name.ValueString(),
	})

	// Update the IPSecPolicy
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecPoliciesUpdate(ctx, id).WritableIPSecPolicyRequest(*ipsecRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IPSecPolicy",
			utils.FormatAPIError("update IPSec policy", err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapIPSecPolicyToState(ctx, ipsec, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, ipsec.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated IPSecPolicy", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IPSecPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IPSecPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec policy ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting IPSecPolicy", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Delete the IPSecPolicy
	httpResp, err := r.client.VpnAPI.VpnIpsecPoliciesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting IPSecPolicy",
			utils.FormatAPIError("delete IPSec policy", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted IPSecPolicy", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state from an existing resource.
func (r *IPSecPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("IPSec policy ID must be a number, got: %s", parsed.ID))
			return
		}

		ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecPoliciesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing IPSecPolicy", utils.FormatAPIError("read IPSec policy", err, httpResp))
			return
		}

		var data IPSecPolicyResourceModel
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

		r.mapIPSecPolicyToState(ctx, ipsec, &data, &resp.Diagnostics)
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the WritableIPSecPolicyRequest.
func (r *IPSecPolicyResource) setOptionalFields(ctx context.Context, ipsecRequest *netbox.WritableIPSecPolicyRequest, plan *IPSecPolicyResourceModel, state *IPSecPolicyResourceModel, diags *diag.Diagnostics) {
	// Proposals
	if utils.IsSet(plan.Proposals) {
		var proposalIDs []int64
		diags.Append(plan.Proposals.ElementsAs(ctx, &proposalIDs, false)...)
		if diags.HasError() {
			return
		}
		proposals := make([]int32, len(proposalIDs))
		for i, id := range proposalIDs {
			val, err := utils.SafeInt32(id)
			if err != nil {
				diags.AddError("Invalid value", fmt.Sprintf("Proposal ID value overflow: %s", err))
				return
			}
			proposals[i] = val
		}
		ipsecRequest.Proposals = proposals
	} else if state != nil && plan.Proposals.IsNull() && utils.IsSet(state.Proposals) {
		// Explicitly clear on Update when removed from config.
		// Note: go-netbox treats a non-nil empty slice as set, so it will be serialized.
		ipsecRequest.Proposals = []int32{}
	}

	// PFS Group
	if utils.IsSet(plan.PFSGroup) {
		pfsGroupVal, err := utils.SafeInt32FromValue(plan.PFSGroup)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("PFSGroup value overflow: %s", err))
			return
		}
		pfsGroup := netbox.PatchedWritableIPSecPolicyRequestPfsGroup(pfsGroupVal)
		ipsecRequest.PfsGroup = *netbox.NewNullablePatchedWritableIPSecPolicyRequestPfsGroup(&pfsGroup)
	}
	// Note: NetBox rejects clearing pfs_group (400 "may not be blank"), so omission should
	// retain the prior value (Optional+Computed).

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

// mapIPSecPolicyToState maps an IPSecPolicy API response to the Terraform state model.
func (r *IPSecPolicyResource) mapIPSecPolicyToState(ctx context.Context, ipsec *netbox.IPSecPolicy, data *IPSecPolicyResourceModel, diags *diag.Diagnostics) {
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

	// Proposals
	if len(ipsec.Proposals) > 0 {
		proposalIDs := make([]int64, len(ipsec.Proposals))
		for i, proposal := range ipsec.Proposals {
			proposalIDs[i] = int64(proposal.Id)
		}
		proposalsValue, _ := types.SetValueFrom(ctx, types.Int64Type, proposalIDs)
		data.Proposals = proposalsValue
	} else {
		data.Proposals = types.SetNull(types.Int64Type)
	}

	// PFS Group
	if ipsec.PfsGroup != nil && ipsec.PfsGroup.Value != nil {
		data.PFSGroup = types.Int64Value(int64(*ipsec.PfsGroup.Value))
	} else {
		data.PFSGroup = types.Int64Null()
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
