// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &IKEProposalResource{}

	_ resource.ResourceWithConfigure = &IKEProposalResource{}

	_ resource.ResourceWithImportState = &IKEProposalResource{}
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
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Description types.String `tfsdk:"description"`

	AuthenticationMethod types.String `tfsdk:"authentication_method"`

	EncryptionAlgorithm types.String `tfsdk:"encryption_algorithm"`

	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`

	Group types.Int64 `tfsdk:"group"`

	SALifetime types.Int64 `tfsdk:"sa_lifetime"`

	Comments types.String `tfsdk:"comments"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
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

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the IKE proposal. Required.",

				Required: true,
			},

			"description": nbschema.DescriptionAttribute("IKE proposal"),

			"authentication_method": schema.StringAttribute{

				MarkdownDescription: "The authentication method for the IKE proposal. Required. Valid values: `preshared-keys`, `certificates`, `rsa-signatures`, `dsa-signatures`.",

				Required: true,

				Validators: []validator.String{

					stringvalidator.OneOf("preshared-keys", "certificates", "rsa-signatures", "dsa-signatures"),
				},
			},

			"encryption_algorithm": schema.StringAttribute{

				MarkdownDescription: "The encryption algorithm for the IKE proposal. Required. Valid values: `aes-128-cbc`, `aes-128-gcm`, `aes-192-cbc`, `aes-192-gcm`, `aes-256-cbc`, `aes-256-gcm`, `3des-cbc`, `des-cbc`.",

				Required: true,

				Validators: []validator.String{

					stringvalidator.OneOf("aes-128-cbc", "aes-128-gcm", "aes-192-cbc", "aes-192-gcm", "aes-256-cbc", "aes-256-gcm", "3des-cbc", "des-cbc"),
				},
			},

			"authentication_algorithm": schema.StringAttribute{

				MarkdownDescription: "The authentication algorithm (hash) for the IKE proposal. Optional. Valid values: `hmac-sha1`, `hmac-sha256`, `hmac-sha384`, `hmac-sha512`, `hmac-md5`.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.OneOf("hmac-sha1", "hmac-sha256", "hmac-sha384", "hmac-sha512", "hmac-md5"),
				},
			},

			"group": schema.Int64Attribute{

				MarkdownDescription: "The Diffie-Hellman group for the IKE proposal. Required. Valid values: 1, 2, 5, 14-34.",

				Required: true,

				Validators: []validator.Int64{

					int64validator.OneOf(1, 2, 5, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34),
				},
			},

			"sa_lifetime": schema.Int64Attribute{

				MarkdownDescription: "Security association lifetime in seconds. Optional.",

				Optional: true,

				PlanModifiers: []planmodifier.Int64{

					int64planmodifier.UseStateForUnknown(),
				},
			},

			"comments": nbschema.CommentsAttribute("IKE proposal"),

			"display_name": nbschema.DisplayNameAttribute("IKE proposal"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

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

	r.setOptionalFields(ctx, ikeRequest, &data, &resp.Diagnostics)

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

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

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

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	// Map response to model

	r.mapIKEProposalToState(ctx, ike, &data, &resp.Diagnostics)

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state on success.

func (r *IKEProposalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data IKEProposalResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	r.setOptionalFields(ctx, ikeRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Updating IKEProposal", map[string]interface{}{

		"id": id,

		"name": data.Name.ValueString(),
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

	// Map response to model

	r.mapIKEProposalToState(ctx, ike, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated IKEProposal", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

		"id": id,

		"name": data.Name.ValueString(),
	})

	// Delete the IKEProposal

	httpResp, err := r.client.VpnAPI.VpnIkeProposalsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// setOptionalFields sets optional fields on the WritableIKEProposalRequest.

func (r *IKEProposalResource) setOptionalFields(ctx context.Context, ikeRequest *netbox.WritableIKEProposalRequest, data *IKEProposalResourceModel, diags *diag.Diagnostics) {

	// Description

	ikeRequest.Description = utils.StringPtr(data.Description)

	// Authentication Algorithm

	if utils.IsSet(data.AuthenticationAlgorithm) {

		authAlg := netbox.PatchedWritableIKEProposalRequestAuthenticationAlgorithm(data.AuthenticationAlgorithm.ValueString())

		ikeRequest.AuthenticationAlgorithm = &authAlg

	}

	// SA Lifetime

	if utils.IsSet(data.SALifetime) {

		lifetime, err := utils.SafeInt32FromValue(data.SALifetime)

		if err != nil {

			diags.AddError("Invalid value", fmt.Sprintf("SALifetime value overflow: %s", err))

			return

		}

		ikeRequest.SaLifetime = *netbox.NewNullableInt32(&lifetime)

	}

	// Comments

	ikeRequest.Comments = utils.StringPtr(data.Comments)

	// Tags

	if utils.IsSet(data.Tags) {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		ikeRequest.Tags = tags

	}

	// Custom Fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if diags.HasError() {

			return

		}

		ikeRequest.CustomFields = utils.CustomFieldsToMap(customFields)

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

	if ike.AuthenticationAlgorithm != nil && ike.AuthenticationAlgorithm.Value != nil {

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

	// Display Name

	if ike.Display != "" {

		data.DisplayName = types.StringValue(ike.Display)

	} else {

		data.DisplayName = types.StringNull()

	}

	// Tags

	if len(ike.Tags) > 0 {

		tags := utils.NestedTagsToTagModels(ike.Tags)

		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Custom Fields

	switch {

	case len(ike.CustomFields) > 0 && !data.CustomFields.IsNull():

		var stateCustomFields []utils.CustomFieldModel

		data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		customFields := utils.MapToCustomFieldModels(ike.CustomFields, stateCustomFields)

		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		data.CustomFields = customFieldsValue

	case len(ike.CustomFields) > 0:

		customFields := utils.MapToCustomFieldModels(ike.CustomFields, []utils.CustomFieldModel{})

		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		data.CustomFields = customFieldsValue

	default:

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
