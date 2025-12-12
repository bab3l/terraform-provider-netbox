// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IPSecProposalResource{}
	_ resource.ResourceWithConfigure   = &IPSecProposalResource{}
	_ resource.ResourceWithImportState = &IPSecProposalResource{}
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
			"comments":      nbschema.CommentsAttribute("IPSec proposal"),
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the resource.
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
	r.setOptionalFields(ctx, ipsecRequest, &data, &resp.Diagnostics)
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
		if httpResp != nil && httpResp.StatusCode == 404 {
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

	// Map response to model
	r.mapIPSecProposalToState(ctx, ipsec, &data, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *IPSecProposalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IPSecProposalResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
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

	// Create the IPSecProposal request
	ipsecRequest := netbox.NewWritableIPSecProposalRequest(data.Name.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, ipsecRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating IPSecProposal", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
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

	// Map response to model
	r.mapIPSecProposalToState(ctx, ipsec, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated IPSecProposal", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		if httpResp != nil && httpResp.StatusCode == 404 {
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the WritableIPSecProposalRequest.
func (r *IPSecProposalResource) setOptionalFields(ctx context.Context, ipsecRequest *netbox.WritableIPSecProposalRequest, data *IPSecProposalResourceModel, diags *diag.Diagnostics) {
	// Description
	ipsecRequest.Description = utils.StringPtr(data.Description)

	// Encryption Algorithm
	if utils.IsSet(data.EncryptionAlgorithm) {
		encAlg := netbox.Encryption(data.EncryptionAlgorithm.ValueString())
		ipsecRequest.EncryptionAlgorithm = &encAlg
	}

	// Authentication Algorithm
	if utils.IsSet(data.AuthenticationAlgorithm) {
		authAlg := netbox.Authentication(data.AuthenticationAlgorithm.ValueString())
		ipsecRequest.AuthenticationAlgorithm = &authAlg
	}

	// SA Lifetime Seconds
	if utils.IsSet(data.SALifetimeSeconds) {
		lifetime, err := utils.SafeInt32FromValue(data.SALifetimeSeconds)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("SALifetimeSeconds value overflow: %s", err))
			return
		}
		ipsecRequest.SaLifetimeSeconds = *netbox.NewNullableInt32(&lifetime)
	}

	// SA Lifetime Data
	if utils.IsSet(data.SALifetimeData) {
		lifetime, err := utils.SafeInt32FromValue(data.SALifetimeData)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("SALifetimeData value overflow: %s", err))
			return
		}
		ipsecRequest.SaLifetimeData = *netbox.NewNullableInt32(&lifetime)
	}

	// Comments
	ipsecRequest.Comments = utils.StringPtr(data.Comments)

	// Tags
	if utils.IsSet(data.Tags) {
		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		ipsecRequest.Tags = tags
	}

	// Custom Fields
	if utils.IsSet(data.CustomFields) {
		var customFields []utils.CustomFieldModel
		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if diags.HasError() {
			return
		}
		ipsecRequest.CustomFields = utils.CustomFieldsToMap(customFields)
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

	// Tags
	if len(ipsec.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(ipsec.Tags)
		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields
	switch {
	case len(ipsec.CustomFields) > 0 && !data.CustomFields.IsNull():
		var stateCustomFields []utils.CustomFieldModel
		data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		customFields := utils.MapToCustomFieldModels(ipsec.CustomFields, stateCustomFields)
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	case len(ipsec.CustomFields) > 0:
		customFields := utils.MapToCustomFieldModels(ipsec.CustomFields, []utils.CustomFieldModel{})
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	default:
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
