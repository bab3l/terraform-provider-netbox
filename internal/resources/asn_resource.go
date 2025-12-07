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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ASNResource{}
	_ resource.ResourceWithConfigure   = &ASNResource{}
	_ resource.ResourceWithImportState = &ASNResource{}
)

// NewASNResource returns a new ASN resource.
func NewASNResource() resource.Resource {
	return &ASNResource{}
}

// ASNResource defines the resource implementation.
type ASNResource struct {
	client *netbox.APIClient
}

// ASNResourceModel describes the resource data model.
type ASNResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ASN          types.Int64  `tfsdk:"asn"`
	RIR          types.String `tfsdk:"rir"`
	Tenant       types.String `tfsdk:"tenant"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ASNResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asn"
}

// Schema defines the schema for the resource.
func (r *ASNResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Autonomous System Number (ASN) in NetBox. ASNs are used for BGP routing and network identification.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the ASN resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"asn": schema.Int64Attribute{
				MarkdownDescription: "The 16- or 32-bit autonomous system number.",
				Required:            true,
			},
			"rir": schema.StringAttribute{
				MarkdownDescription: "The Regional Internet Registry (RIR) that manages this ASN. Can be specified by name, slug, or ID.",
				Optional:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The tenant this ASN is assigned to. Can be specified by name, slug, or ID.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of this ASN.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about this ASN.",
				Optional:            true,
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ASNResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ASNResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ASNResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating ASN", map[string]interface{}{
		"asn": data.ASN.ValueInt64(),
	})

	// Build the ASN request
	asnRequest, diags := r.buildASNRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	asn, httpResp, err := r.client.IpamAPI.IpamAsnsCreate(ctx).ASNRequest(*asnRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ASN",
			utils.FormatAPIError(fmt.Sprintf("create ASN %d", data.ASN.ValueInt64()), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Created ASN", map[string]interface{}{
		"id":  asn.GetId(),
		"asn": asn.GetAsn(),
	})

	// Map response to state
	r.mapResponseToModel(ctx, asn, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ASNResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ASNResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	asnID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ASN ID",
			fmt.Sprintf("ASN ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Reading ASN", map[string]interface{}{
		"id": asnID,
	})

	// Call the API
	asn, httpResp, err := r.client.IpamAPI.IpamAsnsRetrieve(ctx, asnID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "ASN not found, removing from state", map[string]interface{}{
				"id": asnID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading ASN",
			utils.FormatAPIError(fmt.Sprintf("read ASN ID %d", asnID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToModel(ctx, asn, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ASNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ASNResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	asnID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ASN ID",
			fmt.Sprintf("ASN ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Updating ASN", map[string]interface{}{
		"id":  asnID,
		"asn": data.ASN.ValueInt64(),
	})

	// Build the ASN request
	asnRequest, diags := r.buildASNRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	asn, httpResp, err := r.client.IpamAPI.IpamAsnsUpdate(ctx, asnID).ASNRequest(*asnRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ASN",
			utils.FormatAPIError(fmt.Sprintf("update ASN ID %d", asnID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Updated ASN", map[string]interface{}{
		"id":  asn.GetId(),
		"asn": asn.GetAsn(),
	})

	// Map response to state
	r.mapResponseToModel(ctx, asn, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *ASNResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ASNResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	asnID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ASN ID",
			fmt.Sprintf("ASN ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Deleting ASN", map[string]interface{}{
		"id":  asnID,
		"asn": data.ASN.ValueInt64(),
	})

	// Call the API
	httpResp, err := r.client.IpamAPI.IpamAsnsDestroy(ctx, asnID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting ASN",
			utils.FormatAPIError(fmt.Sprintf("delete ASN ID %d", asnID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted ASN", map[string]interface{}{
		"id": asnID,
	})
}

// ImportState imports the resource state.
func (r *ASNResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildASNRequest builds an ASNRequest from the Terraform model.
func (r *ASNResource) buildASNRequest(ctx context.Context, data *ASNResourceModel) (*netbox.ASNRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create the request with required fields
	asnRequest := netbox.NewASNRequest(data.ASN.ValueInt64())

	// Handle RIR (optional)
	if !data.RIR.IsNull() && !data.RIR.IsUnknown() {
		rir, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
		diags.Append(rirDiags...)
		if diags.HasError() {
			return nil, diags
		}
		asnRequest.Rir = *netbox.NewNullableBriefRIRRequest(rir)
	}

	// Handle Tenant (optional)
	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return nil, diags
		}
		asnRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	}

	// Handle description (optional)
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		asnRequest.Description = &desc
	}

	// Handle comments (optional)
	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		asnRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return nil, diags
		}
		asnRequest.Tags = tags
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFieldModels []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return nil, diags
		}
		asnRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)
	}

	return asnRequest, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ASNResource) mapResponseToModel(ctx context.Context, asn *netbox.ASN, data *ASNResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", asn.GetId()))
	data.ASN = types.Int64Value(asn.GetAsn())

	// Map RIR - preserve the user's input format (ID or name)
	if asn.Rir.IsSet() && asn.Rir.Get() != nil {
		// Always return the ID to match what Terraform sent
		data.RIR = types.StringValue(fmt.Sprintf("%d", asn.Rir.Get().GetId()))
	} else {
		data.RIR = types.StringNull()
	}

	// Map Tenant - preserve the user's input format (ID or name)
	if asn.Tenant.IsSet() && asn.Tenant.Get() != nil {
		// Always return the ID to match what Terraform sent
		data.Tenant = types.StringValue(fmt.Sprintf("%d", asn.Tenant.Get().GetId()))
	} else {
		data.Tenant = types.StringNull()
	}

	// Map description
	if desc, ok := asn.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := asn.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if asn.HasTags() {
		tags := utils.NestedTagsToTagModels(asn.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if asn.HasCustomFields() {
		apiCustomFields := asn.GetCustomFields()
		var stateCustomFieldModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			data.CustomFields.ElementsAs(ctx, &stateCustomFieldModels, false)
		}
		customFields := utils.MapToCustomFieldModels(apiCustomFields, stateCustomFieldModels)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
