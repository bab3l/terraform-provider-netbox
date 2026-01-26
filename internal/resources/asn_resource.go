// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
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
	_ resource.Resource                = &ASNResource{}
	_ resource.ResourceWithConfigure   = &ASNResource{}
	_ resource.ResourceWithImportState = &ASNResource{}
	_ resource.ResourceWithIdentity    = &ASNResource{}
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
				Validators: []validator.Int64{
					validators.ValidASNInt64(),
				},
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
			"tags":          nbschema.TagsSlugAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *ASNResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
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

	// Build the ASN request (pass nil state since this is a new resource)
	asnRequest, diags := r.buildASNRequest(ctx, &data, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	asn, httpResp, err := r.client.IpamAPI.IpamAsnsCreate(ctx).ASNRequest(*asnRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
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
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
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
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ASNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ASNResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state for merge-aware custom fields
	var state ASNResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Build the ASN request with state for merge-aware custom fields
	asnRequest, diags := r.buildASNRequest(ctx, &data, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	asn, httpResp, err := r.client.IpamAPI.IpamAsnsUpdate(ctx, asnID).ASNRequest(*asnRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
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
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
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
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		asnID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ASN ID",
				fmt.Sprintf("ASN ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		asn, httpResp, err := r.client.IpamAPI.IpamAsnsRetrieve(ctx, asnID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing ASN",
				utils.FormatAPIError(fmt.Sprintf("read ASN ID %d", asnID), err, httpResp),
			)
			return
		}

		var data ASNResourceModel
		if asn.Rir.IsSet() && asn.Rir.Get() != nil {
			rir := asn.Rir.Get()
			if rir.GetId() != 0 {
				data.RIR = types.StringValue(fmt.Sprintf("%d", rir.GetId()))
			} else {
				data.RIR = types.StringNull()
			}
		} else {
			data.RIR = types.StringNull()
		}
		if asn.Tenant.IsSet() && asn.Tenant.Get() != nil {
			tenant := asn.Tenant.Get()
			if tenant.GetId() != 0 {
				data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
			} else {
				data.Tenant = types.StringNull()
			}
		} else {
			data.Tenant = types.StringNull()
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, asn.HasTags(), asn.GetTags(), data.Tags)
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

		r.mapResponseToModel(ctx, asn, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, asn.GetCustomFields(), &resp.Diagnostics)
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

// buildASNRequest builds an ASNRequest from the Terraform model.
// state is optional and only provided during updates for merge-aware custom fields.
func (r *ASNResource) buildASNRequest(ctx context.Context, data *ASNResourceModel, state *ASNResourceModel) (*netbox.ASNRequest, diag.Diagnostics) {
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
	} else if data.RIR.IsNull() {
		// Explicitly set to null to clear the field
		asnRequest.SetRirNil()
	}

	// Handle Tenant (optional)
	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return nil, diags
		}
		asnRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	} else if data.Tenant.IsNull() {
		// Explicitly set to null to clear the field
		asnRequest.SetTenantNil()
	}

	// Apply description and comments
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		asnRequest.SetDescription(desc)
	} else if data.Description.IsNull() {
		// Explicitly set to empty string to clear the field
		asnRequest.SetDescription("")
	}
	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		asnRequest.SetComments(comments)
	} else if data.Comments.IsNull() {
		// Explicitly set to empty string to clear the field
		asnRequest.SetComments("")
	}

	// Apply tags (slug list)
	utils.ApplyTagsFromSlugs(ctx, r.client, asnRequest, data.Tags, &diags)
	if diags.HasError() {
		return nil, diags
	}

	// Apply custom fields with merge awareness
	if state != nil {
		// Update: use merge-aware helper
		utils.ApplyCustomFieldsWithMerge(ctx, asnRequest, data.CustomFields, state.CustomFields, &diags)
	} else {
		// Create: apply custom fields directly
		utils.ApplyCustomFields(ctx, asnRequest, data.CustomFields, &diags)
	}
	if diags.HasError() {
		return nil, diags
	}

	return asnRequest, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ASNResource) mapResponseToModel(ctx context.Context, asn *netbox.ASN, data *ASNResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", asn.GetId()))
	data.ASN = types.Int64Value(asn.GetAsn())

	// Map RIR (store ID to avoid import drift)
	if asn.Rir.IsSet() && asn.Rir.Get() != nil {
		rir := asn.Rir.Get()
		if rir.GetId() != 0 {
			data.RIR = types.StringValue(fmt.Sprintf("%d", rir.GetId()))
		} else {
			data.RIR = types.StringNull()
		}
	} else {
		data.RIR = types.StringNull()
	}

	// Map Tenant (store ID to avoid import drift)
	if asn.Tenant.IsSet() && asn.Tenant.Get() != nil {
		tenant := asn.Tenant.Get()
		if tenant.GetId() != 0 {
			data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		} else {
			data.Tenant = types.StringNull()
		}
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

	// Tags (slug list)
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, asn.HasTags(), asn.GetTags(), data.Tags)
	if diags.HasError() {
		return
	}

	// Custom Fields - filter to owned fields only
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, asn.GetCustomFields(), diags)
}
