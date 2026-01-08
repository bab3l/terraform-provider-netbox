// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"strconv"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
	_ resource.Resource                = &IPSecProfileResource{}
	_ resource.ResourceWithConfigure   = &IPSecProfileResource{}
	_ resource.ResourceWithImportState = &IPSecProfileResource{}
)

// NewIPSecProfileResource returns a new IPSecProfile resource.
func NewIPSecProfileResource() resource.Resource {
	return &IPSecProfileResource{}
}

// IPSecProfileResource defines the resource implementation.
type IPSecProfileResource struct {
	client *netbox.APIClient
}

// IPSecProfileResourceModel describes the resource data model.
type IPSecProfileResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Mode         types.String `tfsdk:"mode"`
	IKEPolicy    types.String `tfsdk:"ike_policy"`
	IPSecPolicy  types.String `tfsdk:"ipsec_policy"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IPSecProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_profile"
}

// Schema defines the schema for the resource.
func (r *IPSecProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IPSec Profile in Netbox. IPSec profiles combine IKE and IPSec policies to define complete VPN configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IPSec profile.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec profile. Required.",
				Required:            true,
			},
			"description": nbschema.DescriptionAttribute("IPSec profile"),
			"mode": schema.StringAttribute{
				MarkdownDescription: "The IPSec mode. Required. Valid values: `esp` (Encapsulating Security Payload), `ah` (Authentication Header).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("esp", "ah"),
				},
			},
			"ike_policy": schema.StringAttribute{
				MarkdownDescription: "The name of the IKE policy to use. Required.",
				Required:            true,
			},
			"ipsec_policy": schema.StringAttribute{
				MarkdownDescription: "The name of the IPSec policy to use. Required.",
				Required:            true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("IPSec profile"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *IPSecProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IPSecProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IPSecProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IPSecProfile request
	mode := netbox.IPSecProfileModeValue(data.Mode.ValueString())

	// Parse policy IDs
	ikePolicyID, _ := strconv.Atoi(data.IKEPolicy.ValueString())
	ipsecPolicyID, _ := strconv.Atoi(data.IPSecPolicy.ValueString())

	// Create dummy policy objects (will be overridden by AdditionalProperties)
	ikePolicy := netbox.NewBriefIKEPolicyRequest("placeholder")
	ipsecPolicy := netbox.NewBriefIPSecPolicyRequest("placeholder")
	ipsecRequest := netbox.NewWritableIPSecProfileRequest(
		data.Name.ValueString(),
		mode,
		*ikePolicy,
		*ipsecPolicy,
	)

	// Override the policy fields with integer IDs using AdditionalProperties
	// This replaces the nested objects with simple integer IDs that Netbox API accepts
	ipsecRequest.AdditionalProperties = map[string]interface{}{
		"ike_policy":   ikePolicyID,
		"ipsec_policy": ipsecPolicyID,
	}

	// Set optional fields
	r.setOptionalFields(ctx, ipsecRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating IPSecProfile", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Create the IPSecProfile
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProfilesCreate(ctx).WritableIPSecProfileRequest(*ipsecRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IPSecProfile",
			utils.FormatAPIError("create IPSec profile", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIPSecProfileToState(ctx, ipsec, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created IPSecProfile", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *IPSecProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IPSecProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec profile ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading IPSecProfile", map[string]interface{}{
		"id": id,
	})

	// Read the IPSecProfile
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProfilesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "IPSecProfile not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IPSecProfile",
			utils.FormatAPIError("read IPSec profile", err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state

	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapIPSecProfileToState(ctx, ipsec, &data, &resp.Diagnostics)

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
func (r *IPSecProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan IPSecProfileResourceModel

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
			fmt.Sprintf("Could not parse IPSec profile ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Create the IPSecProfile request
	mode := netbox.IPSecProfileModeValue(plan.Mode.ValueString())

	// Parse policy IDs
	ikePolicyID, _ := strconv.Atoi(plan.IKEPolicy.ValueString())
	ipsecPolicyID, _ := strconv.Atoi(plan.IPSecPolicy.ValueString())

	// Create dummy policy objects (will be overridden by AdditionalProperties)
	ikePolicy := netbox.NewBriefIKEPolicyRequest("placeholder")
	ipsecPolicy := netbox.NewBriefIPSecPolicyRequest("placeholder")
	ipsecRequest := netbox.NewWritableIPSecProfileRequest(
		plan.Name.ValueString(),
		mode,
		*ikePolicy,
		*ipsecPolicy,
	)

	// Override the policy fields with integer IDs using AdditionalProperties
	// This replaces the nested objects with simple integer IDs that Netbox API accepts
	ipsecRequest.AdditionalProperties = map[string]interface{}{
		"ike_policy":   ikePolicyID,
		"ipsec_policy": ipsecPolicyID,
	}

	// Set optional fields with state for merge-aware custom fields
	r.setOptionalFields(ctx, ipsecRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IPSecProfile", map[string]interface{}{
		"id":   id,
		"name": plan.Name.ValueString(),
	})

	// Update the IPSecProfile
	ipsec, httpResp, err := r.client.VpnAPI.VpnIpsecProfilesUpdate(ctx, id).WritableIPSecProfileRequest(*ipsecRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IPSecProfile",
			utils.FormatAPIError("update IPSec profile", err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapIPSecProfileToState(ctx, ipsec, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, ipsec.GetCustomFields(), &resp.Diagnostics)

	tflog.Debug(ctx, "Updated IPSecProfile", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IPSecProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IPSecProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing ID",
			fmt.Sprintf("Could not parse IPSec profile ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting IPSecProfile", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Delete the IPSecProfile
	httpResp, err := r.client.VpnAPI.VpnIpsecProfilesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting IPSecProfile",
			utils.FormatAPIError("delete IPSec profile", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted IPSecProfile", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state from an existing resource.
func (r *IPSecProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the WritableIPSecProfileRequest.
func (r *IPSecProfileResource) setOptionalFields(ctx context.Context, ipsecRequest *netbox.WritableIPSecProfileRequest, plan *IPSecProfileResourceModel, state *IPSecProfileResourceModel, diags *diag.Diagnostics) {
	// Set description
	utils.ApplyDescription(ipsecRequest, plan.Description)

	// Set comments, tags, and custom fields with merge-aware helpers
	utils.ApplyComments(ipsecRequest, plan.Comments)
	utils.ApplyTags(ctx, ipsecRequest, plan.Tags, diags)
	// Apply custom fields with merge logic to preserve unmanaged fields
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, ipsecRequest, plan.CustomFields, state.CustomFields, diags)
	} else {
		// During Create, no state exists yet
		utils.ApplyCustomFields(ctx, ipsecRequest, plan.CustomFields, diags)
	}
}

// mapIPSecProfileToState maps an IPSecProfile API response to the Terraform state model.
func (r *IPSecProfileResource) mapIPSecProfileToState(ctx context.Context, ipsec *netbox.IPSecProfile, data *IPSecProfileResourceModel, diags *diag.Diagnostics) {
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

	// Mode
	if ipsec.Mode.Value != nil {
		data.Mode = types.StringValue(string(*ipsec.Mode.Value))
	}

	// IKE Policy (store ID as string)
	data.IKEPolicy = types.StringValue(fmt.Sprintf("%d", ipsec.IkePolicy.Id))

	// IPSec Policy (store ID as string)
	data.IPSecPolicy = types.StringValue(fmt.Sprintf("%d", ipsec.IpsecPolicy.Id))

	// Comments
	if ipsec.Comments != nil && *ipsec.Comments != "" {
		data.Comments = types.StringValue(*ipsec.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsFromAPI(ctx, ipsec.HasTags(), ipsec.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, ipsec.HasCustomFields(), ipsec.GetCustomFields(), data.CustomFields, diags)
}
