// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &WirelessLANResource{}
	_ resource.ResourceWithConfigure   = &WirelessLANResource{}
	_ resource.ResourceWithImportState = &WirelessLANResource{}
	_ resource.ResourceWithIdentity    = &WirelessLANResource{}
)

// NewWirelessLANResource returns a new resource implementing the wireless LAN resource.
func NewWirelessLANResource() resource.Resource {
	return &WirelessLANResource{}
}

// WirelessLANResource defines the resource implementation.
type WirelessLANResource struct {
	client *netbox.APIClient
}

// WirelessLANResourceModel describes the resource data model.
type WirelessLANResourceModel struct {
	ID           types.String `tfsdk:"id"`
	SSID         types.String `tfsdk:"ssid"`
	Description  types.String `tfsdk:"description"`
	Group        types.String `tfsdk:"group"`
	Status       types.String `tfsdk:"status"`
	VLAN         types.String `tfsdk:"vlan"`
	Tenant       types.String `tfsdk:"tenant"`
	AuthType     types.String `tfsdk:"auth_type"`
	AuthCipher   types.String `tfsdk:"auth_cipher"`
	AuthPSK      types.String `tfsdk:"auth_psk"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *WirelessLANResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_lan"
}

// Schema defines the schema for the resource.
func (r *WirelessLANResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a wireless LAN (WiFi network) in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the wireless LAN.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssid": schema.StringAttribute{
				MarkdownDescription: "The SSID (network name) of the wireless LAN.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the wireless LAN.",
				Optional:            true,
			},
			"group": nbschema.ReferenceAttributeWithDiffSuppress(
				"wireless LAN group",
				"The wireless LAN group this network belongs to (ID or slug).",
			),
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the wireless LAN. Valid values: `active`, `reserved`, `disabled`, `deprecated`. Default: `active`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
			},
			"vlan": nbschema.ReferenceAttributeWithDiffSuppress(
				"vlan",
				"The VLAN associated with this wireless LAN (ID or name).",
			),
			"tenant": nbschema.ReferenceAttributeWithDiffSuppress(
				"tenant",
				"The tenant this wireless LAN belongs to (ID or slug).",
			),
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type. Valid values: `open`, `wep`, `wpa-personal`, `wpa-enterprise`.",
				Optional:            true,
			},
			"auth_cipher": schema.StringAttribute{
				MarkdownDescription: "Authentication cipher. Valid values: `auto`, `tkip`, `aes`.",
				Optional:            true,
			},
			"auth_psk": schema.StringAttribute{
				MarkdownDescription: "Pre-shared key for authentication.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("wireless LAN"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *WirelessLANResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *WirelessLANResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.
func (r *WirelessLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WirelessLANResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritableWirelessLANRequest(data.SSID.ValueString())

	// Set optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	if !data.Group.IsNull() && !data.Group.IsUnknown() {
		group, diags := lookup.LookupWirelessLANGroup(ctx, r.client, data.Group.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetGroup(*group)
	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.PatchedWritableWirelessLANRequestStatus(data.Status.ValueString())
		apiReq.SetStatus(status)
	}

	if !data.VLAN.IsNull() && !data.VLAN.IsUnknown() {
		vlan, diags := lookup.LookupVLAN(ctx, r.client, data.VLAN.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetVlan(*vlan)
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, diags := lookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTenant(*tenant)
	}

	if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
		authType := netbox.AuthenticationType1(data.AuthType.ValueString())
		apiReq.SetAuthType(authType)
	}

	if !data.AuthCipher.IsNull() && !data.AuthCipher.IsUnknown() {
		authCipher := netbox.AuthenticationCipher(data.AuthCipher.ValueString())
		apiReq.SetAuthCipher(authCipher)
	}

	if !data.AuthPSK.IsNull() && !data.AuthPSK.IsUnknown() {
		apiReq.SetAuthPsk(data.AuthPSK.ValueString())
	}

	// Handle description and comments
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyComments(apiReq, data.Comments)

	// Store plan values for filter-to-owned population later
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Set tags and custom fields
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating wireless LAN", map[string]interface{}{
		"ssid": data.SSID.ValueString(),
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLansCreate(ctx).WritableWirelessLANRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating wireless LAN",
			utils.FormatAPIError(fmt.Sprintf("create wireless LAN %s", data.SSID.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, response.HasTags(), response.GetTags(), planTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Created wireless LAN", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"ssid": data.SSID.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *WirelessLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WirelessLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wlanID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Wireless LAN ID",
			fmt.Sprintf("Wireless LAN ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading wireless LAN", map[string]interface{}{
		"id": wlanID,
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLansRetrieve(ctx, wlanID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading wireless LAN",
			utils.FormatAPIError(fmt.Sprintf("read wireless LAN ID %d", wlanID), err, httpResp),
		)
		return
	}

	// Preserve sensitive field from state since API doesn't return it
	existingPSK := data.AuthPSK

	// Store state values for filter-to-owned (preserve null vs empty set distinction)
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore sensitive field
	data.AuthPSK = existingPSK

	// Populate tags and custom fields filtered to owned fields only (preserves null/empty state)
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, response.HasTags(), response.GetTags(), stateTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *WirelessLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan WirelessLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wlanID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Wireless LAN ID",
			fmt.Sprintf("Wireless LAN ID must be a number, got: %s", plan.ID.ValueString()),
		)
		return
	}

	// Build request
	apiReq := netbox.NewWritableWirelessLANRequest(plan.SSID.ValueString())

	// Set optional fields
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiReq.SetDescription(plan.Description.ValueString())
	}

	if !plan.Group.IsNull() && !plan.Group.IsUnknown() {
		group, diags := lookup.LookupWirelessLANGroup(ctx, r.client, plan.Group.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetGroup(*group)
	}

	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		status := netbox.PatchedWritableWirelessLANRequestStatus(plan.Status.ValueString())
		apiReq.SetStatus(status)
	}

	if !plan.VLAN.IsNull() && !plan.VLAN.IsUnknown() {
		vlan, diags := lookup.LookupVLAN(ctx, r.client, plan.VLAN.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetVlan(*vlan)
	}

	if !plan.Tenant.IsNull() && !plan.Tenant.IsUnknown() {
		tenant, diags := lookup.LookupTenant(ctx, r.client, plan.Tenant.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTenant(*tenant)
	}

	if !plan.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
		authType := netbox.AuthenticationType1(plan.AuthType.ValueString())
		apiReq.SetAuthType(authType)
	} else if plan.AuthType.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["auth_type"] = nil
	}

	if !plan.AuthCipher.IsNull() && !plan.AuthCipher.IsUnknown() {
		authCipher := netbox.AuthenticationCipher(plan.AuthCipher.ValueString())
		apiReq.SetAuthCipher(authCipher)
	} else if plan.AuthCipher.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["auth_cipher"] = nil
	}

	if !plan.AuthPSK.IsNull() && !plan.AuthPSK.IsUnknown() {
		apiReq.SetAuthPsk(plan.AuthPSK.ValueString())
	}

	// Handle description and comments
	utils.ApplyDescription(apiReq, plan.Description)
	utils.ApplyComments(apiReq, plan.Comments)

	// Handle tags and custom fields with merge-aware helpers
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating wireless LAN", map[string]interface{}{
		"id": wlanID,
	})

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLansUpdate(ctx, wlanID).WritableWirelessLANRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating wireless LAN",
			utils.FormatAPIError(fmt.Sprintf("update wireless LAN ID %d", wlanID), err, httpResp),
		)
		return
	}

	// Preserve sensitive field from plan AND store plan custom fields/tags before mapping
	planTags := plan.Tags
	planCustomFields := plan.CustomFields
	existingPSK := plan.AuthPSK

	// Map response to plan model
	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore sensitive field
	plan.AuthPSK = existingPSK

	// Populate tags and custom fields filtered to owned fields only
	plan.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, response.HasTags(), response.GetTags(), planTags)
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *WirelessLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WirelessLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wlanID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Wireless LAN ID",
			fmt.Sprintf("Wireless LAN ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting wireless LAN", map[string]interface{}{
		"id": wlanID,
	})

	httpResp, err := r.client.WirelessAPI.WirelessWirelessLansDestroy(ctx, wlanID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting wireless LAN",
			utils.FormatAPIError(fmt.Sprintf("delete wireless LAN ID %d", wlanID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *WirelessLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		wlanID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Wireless LAN ID must be a number, got: %s", parsed.ID))
			return
		}

		response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLansRetrieve(ctx, wlanID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing wireless LAN", utils.FormatAPIError(fmt.Sprintf("import wireless LAN ID %d", wlanID), err, httpResp))
			return
		}

		var data WirelessLANResourceModel
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

		r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, response.GetCustomFields(), &resp.Diagnostics)
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

// mapResponseToModel maps the API response to the Terraform model.
func (r *WirelessLANResource) mapResponseToModel(ctx context.Context, wlan *netbox.WirelessLAN, data *WirelessLANResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", wlan.GetId()))
	data.SSID = types.StringValue(wlan.GetSsid())

	// Map description
	if desc, ok := wlan.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map group
	if wlan.Group.IsSet() && wlan.Group.Get() != nil {
		group := wlan.Group.Get()
		data.Group = utils.UpdateReferenceAttribute(data.Group, group.GetName(), group.GetSlug(), group.GetId())
	} else {
		data.Group = types.StringNull()
	}

	// Map status
	if status, ok := wlan.GetStatusOk(); ok && status != nil {
		data.Status = types.StringValue(string(status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Map VLAN - preserve user's input format
	if wlan.Vlan.IsSet() && wlan.Vlan.Get() != nil {
		vlan := wlan.Vlan.Get()
		data.VLAN = utils.UpdateReferenceAttribute(data.VLAN, vlan.GetName(), "", vlan.GetId())
	} else {
		data.VLAN = types.StringNull()
	}

	// Map tenant - preserve user's input format
	if wlan.Tenant.IsSet() && wlan.Tenant.Get() != nil {
		tenant := wlan.Tenant.Get()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Map auth_type
	if authType, ok := wlan.GetAuthTypeOk(); ok && authType != nil {
		data.AuthType = types.StringValue(string(authType.GetValue()))
	} else {
		data.AuthType = types.StringNull()
	}

	// Map auth_cipher
	if authCipher, ok := wlan.GetAuthCipherOk(); ok && authCipher != nil {
		data.AuthCipher = types.StringValue(string(authCipher.GetValue()))
	} else {
		data.AuthCipher = types.StringNull()
	}

	// Note: auth_psk is not returned by API, handled separately in Read/Update
	// Map comments
	if comments, ok := wlan.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags (slug list) with empty-set preservation
	data.Tags = utils.PopulateTagsSlugFromAPI(ctx, wlan.HasTags(), wlan.GetTags(), data.Tags)

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, wlan.GetCustomFields(), diags)
}
