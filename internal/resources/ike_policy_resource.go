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
	_ resource.Resource                = &IKEPolicyResource{}
	_ resource.ResourceWithConfigure   = &IKEPolicyResource{}
	_ resource.ResourceWithImportState = &IKEPolicyResource{}
)

// NewIKEPolicyResource returns a new IKEPolicy resource.
func NewIKEPolicyResource() resource.Resource {
	return &IKEPolicyResource{}
}

// IKEPolicyResource defines the resource implementation.
type IKEPolicyResource struct {
	client *netbox.APIClient
}

// IKEPolicyResourceModel describes the resource data model.
type IKEPolicyResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Version      types.Int64  `tfsdk:"version"`
	Mode         types.String `tfsdk:"mode"`
	Proposals    types.Set    `tfsdk:"proposals"`
	PresharedKey types.String `tfsdk:"preshared_key"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IKEPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ike_policy"
}

// Schema defines the schema for the resource.
func (r *IKEPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IKE (Internet Key Exchange) Policy in Netbox. IKE policies group together IKE proposals and define the IKE version and mode for IPSec VPN connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IKE policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IKE policy. Required.",
				Required:            true,
			},
			"description": nbschema.DescriptionAttribute("IKE policy"),
			"version": schema.Int64Attribute{
				MarkdownDescription: "The IKE version. Valid values: `1` (IKEv1), `2` (IKEv2). Defaults to 1.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The IKE negotiation mode. Valid values: `aggressive`, `main`. Only applicable for IKEv1.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("aggressive", "main"),
				},
			},
			"proposals": schema.SetAttribute{
				MarkdownDescription: "A set of IKE proposal IDs to associate with this policy.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"preshared_key": schema.StringAttribute{
				MarkdownDescription: "The pre-shared key for IKE authentication. Optional.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("IKE policy"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.
func (r *IKEPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IKEPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IKEPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IKEPolicy request
	ikeRequest := netbox.NewWritableIKEPolicyRequest(data.Name.ValueString())

	// Set optional fields (no state during Create)
	r.setOptionalFields(ctx, ikeRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating IKEPolicy", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Create the IKEPolicy
	ike, httpResp, err := r.client.VpnAPI.VpnIkePoliciesCreate(ctx).WritableIKEPolicyRequest(*ikeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IKEPolicy",
			utils.FormatAPIError("create IKE policy", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIKEPolicyToState(ctx, ike, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created IKEPolicy", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *IKEPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IKEPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IKE policy ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading IKEPolicy", map[string]interface{}{
		"id": id,
	})

	// Read the IKEPolicy
	ike, httpResp, err := r.client.VpnAPI.VpnIkePoliciesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "IKEPolicy not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IKEPolicy",
			utils.FormatAPIError("read IKE policy", err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state
	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapIKEPolicyToState(ctx, ike, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was null or empty before, restore that state
	// This prevents drift when config doesn't declare custom_fields
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *IKEPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan IKEPolicyResourceModel

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
			fmt.Sprintf("Could not parse IKE policy ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Create the IKEPolicy request
	ikeRequest := netbox.NewWritableIKEPolicyRequest(plan.Name.ValueString())

	// Set optional fields with state for merge-aware custom fields
	r.setOptionalFields(ctx, ikeRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IKEPolicy", map[string]interface{}{
		"id":   id,
		"name": plan.Name.ValueString(),
	})

	// Update the IKEPolicy
	ike, httpResp, err := r.client.VpnAPI.VpnIkePoliciesUpdate(ctx, id).WritableIKEPolicyRequest(*ikeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IKEPolicy",
			utils.FormatAPIError("update IKE policy", err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapIKEPolicyToState(ctx, ike, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, ike.GetCustomFields(), &resp.Diagnostics)

	tflog.Debug(ctx, "Updated IKEPolicy", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IKEPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IKEPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IKE policy ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting IKEPolicy", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Delete the IKEPolicy
	httpResp, err := r.client.VpnAPI.VpnIkePoliciesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting IKEPolicy",
			utils.FormatAPIError("delete IKE policy", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted IKEPolicy", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state from an existing resource.
func (r *IKEPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the WritableIKEPolicyRequest.
func (r *IKEPolicyResource) setOptionalFields(ctx context.Context, ikeRequest *netbox.WritableIKEPolicyRequest, plan *IKEPolicyResourceModel, state *IKEPolicyResourceModel, diags *diag.Diagnostics) {
	// Set description
	utils.ApplyDescription(ikeRequest, plan.Description)

	// Version
	if utils.IsSet(plan.Version) {
		versionVal, err := utils.SafeInt32FromValue(plan.Version)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("Version value overflow: %s", err))
			return
		}
		version := netbox.PatchedWritableIKEPolicyRequestVersion(versionVal)
		ikeRequest.Version = &version
	}

	// Mode
	if utils.IsSet(plan.Mode) {
		mode := netbox.PatchedWritableIKEPolicyRequestMode(plan.Mode.ValueString())
		ikeRequest.Mode = &mode
	}

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
		ikeRequest.Proposals = proposals
	}

	// Preshared Key
	if utils.IsSet(plan.PresharedKey) {
		key := plan.PresharedKey.ValueString()
		ikeRequest.PresharedKey = &key
	}

	// Set comments, tags, and custom fields with merge-aware helpers
	utils.ApplyComments(ikeRequest, plan.Comments)
	utils.ApplyTags(ctx, ikeRequest, plan.Tags, diags)
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

// mapIKEPolicyToState maps an IKEPolicy API response to the Terraform state model.
func (r *IKEPolicyResource) mapIKEPolicyToState(ctx context.Context, ike *netbox.IKEPolicy, data *IKEPolicyResourceModel, diags *diag.Diagnostics) {
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

	// Version
	if ike.Version.Value != nil {
		data.Version = types.Int64Value(int64(*ike.Version.Value))
	} else {
		data.Version = types.Int64Null()
	}

	// Mode
	if ike.Mode != nil && ike.Mode.Value != nil && *ike.Mode.Value != "" {
		data.Mode = types.StringValue(string(*ike.Mode.Value))
	} else {
		data.Mode = types.StringNull()
	}

	// Proposals
	if len(ike.Proposals) > 0 {
		proposalIDs := make([]int64, len(ike.Proposals))
		for i, proposal := range ike.Proposals {
			proposalIDs[i] = int64(proposal.Id)
		}
		proposalsValue, _ := types.SetValueFrom(ctx, types.Int64Type, proposalIDs)
		data.Proposals = proposalsValue
	} else {
		data.Proposals = types.SetNull(types.Int64Type)
	}

	// Preshared Key - don't read from API (sensitive, not returned)
	// Keep the value from state if it exists
	if data.PresharedKey.IsNull() {
		data.PresharedKey = types.StringNull()
	}

	// Comments
	if ike.Comments != nil && *ike.Comments != "" {
		data.Comments = types.StringValue(*ike.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsFromAPI(ctx, ike.HasTags(), ike.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields using filter-to-owned pattern
	// Only populate fields that are declared in the config
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, ike.GetCustomFields(), diags)
}
