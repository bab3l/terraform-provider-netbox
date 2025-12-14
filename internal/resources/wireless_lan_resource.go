// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource = &WirelessLANResource{}

	_ resource.ResourceWithConfigure = &WirelessLANResource{}

	_ resource.ResourceWithImportState = &WirelessLANResource{}
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
	ID types.String `tfsdk:"id"`

	SSID types.String `tfsdk:"ssid"`

	Description types.String `tfsdk:"description"`

	Group types.String `tfsdk:"group"`

	Status types.String `tfsdk:"status"`

	VLAN types.String `tfsdk:"vlan"`

	Tenant types.String `tfsdk:"tenant"`

	AuthType types.String `tfsdk:"auth_type"`

	AuthCipher types.String `tfsdk:"auth_cipher"`

	AuthPSK types.String `tfsdk:"auth_psk"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
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

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"ssid": schema.StringAttribute{

				MarkdownDescription: "The SSID (network name) of the wireless LAN.",

				Required: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the wireless LAN.",

				Optional: true,
			},

			"group": schema.StringAttribute{

				MarkdownDescription: "The wireless LAN group this network belongs to (ID or slug).",

				Optional: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "Status of the wireless LAN. Valid values: `active`, `reserved`, `disabled`, `deprecated`. Default: `active`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("active"),
			},

			"vlan": schema.StringAttribute{

				MarkdownDescription: "The VLAN associated with this wireless LAN (ID or name).",

				Optional: true,
			},

			"tenant": schema.StringAttribute{

				MarkdownDescription: "The tenant this wireless LAN belongs to (ID or slug).",

				Optional: true,
			},

			"auth_type": schema.StringAttribute{

				MarkdownDescription: "Authentication type. Valid values: `open`, `wep`, `wpa-personal`, `wpa-enterprise`.",

				Optional: true,
			},

			"auth_cipher": schema.StringAttribute{

				MarkdownDescription: "Authentication cipher. Valid values: `auto`, `tkip`, `aes`.",

				Optional: true,
			},

			"auth_psk": schema.StringAttribute{

				MarkdownDescription: "Pre-shared key for authentication.",

				Optional: true,

				Sensitive: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about the wireless LAN.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

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

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

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

	tflog.Trace(ctx, "Created wireless LAN", map[string]interface{}{

		"id": data.ID.ValueString(),

		"ssid": data.SSID.ValueString(),
	})

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

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Restore sensitive field

	data.AuthPSK = existingPSK

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource.

func (r *WirelessLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data WirelessLANResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

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

	// Preserve sensitive field

	existingPSK := data.AuthPSK

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Restore sensitive field

	data.AuthPSK = existingPSK

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	wlanID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Wireless LAN ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.WirelessAPI.WirelessWirelessLansRetrieve(ctx, wlanID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing wireless LAN",

			utils.FormatAPIError(fmt.Sprintf("import wireless LAN ID %d", wlanID), err, httpResp),
		)

		return

	}

	var data WirelessLANResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

		data.Group = types.StringValue(fmt.Sprintf("%d", wlan.Group.Get().GetId()))

	} else {

		data.Group = types.StringNull()

	}

	// Map status

	if status, ok := wlan.GetStatusOk(); ok && status != nil {

		data.Status = types.StringValue(string(status.GetValue()))

	} else {

		data.Status = types.StringNull()

	}

	// Map VLAN

	if wlan.Vlan.IsSet() && wlan.Vlan.Get() != nil {

		data.VLAN = types.StringValue(fmt.Sprintf("%d", wlan.Vlan.Get().GetId()))

	} else {

		data.VLAN = types.StringNull()

	}

	// Map tenant

	if wlan.Tenant.IsSet() && wlan.Tenant.Get() != nil {

		data.Tenant = types.StringValue(fmt.Sprintf("%d", wlan.Tenant.Get().GetId()))

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

	// Handle tags

	if wlan.HasTags() && len(wlan.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(wlan.GetTags())

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

	if wlan.HasCustomFields() {

		apiCustomFields := wlan.GetCustomFields()

		var stateCustomFieldModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

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
