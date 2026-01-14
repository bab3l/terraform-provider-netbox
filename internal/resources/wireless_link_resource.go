// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
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
	_ resource.Resource = &WirelessLinkResource{}

	_ resource.ResourceWithConfigure = &WirelessLinkResource{}

	_ resource.ResourceWithImportState = &WirelessLinkResource{}
)

// NewWirelessLinkResource returns a new resource implementing the wireless link resource.

func NewWirelessLinkResource() resource.Resource {
	return &WirelessLinkResource{}
}

// WirelessLinkResource defines the resource implementation.

type WirelessLinkResource struct {
	client *netbox.APIClient
}

// WirelessLinkResourceModel describes the resource data model.

type WirelessLinkResourceModel struct {
	ID types.String `tfsdk:"id"`

	InterfaceA types.String `tfsdk:"interface_a"`

	InterfaceB types.String `tfsdk:"interface_b"`

	SSID types.String `tfsdk:"ssid"`

	Status types.String `tfsdk:"status"`

	Tenant types.String `tfsdk:"tenant"`

	AuthType types.String `tfsdk:"auth_type"`

	AuthCipher types.String `tfsdk:"auth_cipher"`

	AuthPSK types.String `tfsdk:"auth_psk"`

	Distance types.Float64 `tfsdk:"distance"`

	DistanceUnit types.String `tfsdk:"distance_unit"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *WirelessLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_link"
}

// Schema defines the schema for the resource.

func (r *WirelessLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a wireless link between two interfaces in NetBox. Wireless links represent point-to-point wireless connections between devices.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the wireless link.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"interface_a": schema.StringAttribute{
				MarkdownDescription: "ID of the first interface (A-side) of the wireless link.",

				Required: true,
			},

			"interface_b": schema.StringAttribute{
				MarkdownDescription: "ID of the second interface (B-side) of the wireless link.",

				Required: true,
			},

			"ssid": schema.StringAttribute{
				MarkdownDescription: "The SSID (network name) for the wireless link.",

				Optional: true,

				Validators: []validator.String{
					stringvalidator.LengthAtMost(32),
				},
			},

			"status": schema.StringAttribute{
				MarkdownDescription: "Connection status. Valid values: `connected`, `planned`, `decommissioning`.",

				Optional: true,

				Computed: true,

				Validators: []validator.String{
					stringvalidator.OneOf("connected", "planned", "decommissioning"),
				},
			},

			"tenant": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant that owns this wireless link.",

				Optional: true,
			},

			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type. Valid values: `open`, `wep`, `wpa-personal`, `wpa-enterprise`.",

				Optional: true,

				Validators: []validator.String{
					stringvalidator.OneOf("open", "wep", "wpa-personal", "wpa-enterprise", ""),
				},
			},

			"auth_cipher": schema.StringAttribute{
				MarkdownDescription: "Authentication cipher. Valid values: `auto`, `tkip`, `aes`.",

				Optional: true,

				Validators: []validator.String{
					stringvalidator.OneOf("auto", "tkip", "aes", ""),
				},
			},

			"auth_psk": schema.StringAttribute{
				MarkdownDescription: "Pre-shared key for authentication.",

				Optional: true,

				Sensitive: true,
			},

			"distance": schema.Float64Attribute{
				MarkdownDescription: "Distance of the wireless link.",

				Optional: true,

				Validators: []validator.Float64{
					float64validator.AtLeast(0),
				},
			},

			"distance_unit": schema.StringAttribute{
				MarkdownDescription: "Unit for distance. Valid values: `km`, `m`, `mi`, `ft`.",

				Optional: true,

				Validators: []validator.String{
					stringvalidator.OneOf("km", "m", "mi", "ft", ""),
				},
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("wireless link"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.

func (r *WirelessLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// lookupInterfaceBrief looks up an interface by ID and returns a BriefInterfaceRequest.

func (r *WirelessLinkResource) lookupInterfaceBrief(ctx context.Context, interfaceID string) (*netbox.BriefInterfaceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	id, err := utils.ParseID(interfaceID)

	if err != nil {
		diags.AddError("Invalid Interface ID", fmt.Sprintf("Interface ID must be a number, got: %s", interfaceID))

		return nil, diags
	}

	iface, httpResp, err := r.client.DcimAPI.DcimInterfacesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		diags.AddError("Error Looking Up Interface",

			utils.FormatAPIError(fmt.Sprintf("lookup interface ID %d", id), err, httpResp))

		return nil, diags
	}

	device := iface.GetDevice()

	// BriefDeviceRequest uses NullableString for Name

	briefDeviceRequest := netbox.BriefDeviceRequest{}

	briefDeviceRequest.SetName(device.GetName())

	briefInterfaceRequest := netbox.NewBriefInterfaceRequest(briefDeviceRequest, iface.GetName())

	return briefInterfaceRequest, diags
}

// Create creates the resource.

func (r *WirelessLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WirelessLinkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating wireless link", map[string]interface{}{
		"interface_a": data.InterfaceA.ValueString(),

		"interface_b": data.InterfaceB.ValueString(),
	})

	// Look up interface A

	interfaceA, diags := r.lookupInterfaceBrief(ctx, data.InterfaceA.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Look up interface B

	interfaceB, diags := r.lookupInterfaceBrief(ctx, data.InterfaceB.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request

	request := netbox.NewWritableWirelessLinkRequest(*interfaceA, *interfaceB)

	// Set optional fields

	if !data.SSID.IsNull() && !data.SSID.IsUnknown() {
		ssid := data.SSID.ValueString()

		request.Ssid = &ssid
	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.CableStatusValue(data.Status.ValueString())

		request.Status = &status
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenantID, err := utils.ParseID(data.Tenant.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Tenant ID must be a number, got: %s", data.Tenant.ValueString()))

			return
		}

		tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError("Error Looking Up Tenant",

				utils.FormatAPIError(fmt.Sprintf("lookup tenant ID %d", tenantID), err, httpResp))

			return
		}

		tenantRequest := netbox.BriefTenantRequest{
			Name: tenant.GetName(),

			Slug: tenant.GetSlug(),
		}

		request.Tenant = *netbox.NewNullableBriefTenantRequest(&tenantRequest)
	} else if data.Tenant.IsNull() {
		request.SetTenantNil()
	}

	if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
		authType := netbox.AuthenticationType1(data.AuthType.ValueString())

		request.AuthType = &authType
	}

	if !data.AuthCipher.IsNull() && !data.AuthCipher.IsUnknown() {
		authCipher := netbox.AuthenticationCipher(data.AuthCipher.ValueString())

		request.AuthCipher = &authCipher
	}

	if !data.AuthPSK.IsNull() && !data.AuthPSK.IsUnknown() {
		psk := data.AuthPSK.ValueString()

		request.AuthPsk = &psk
	}

	if !data.Distance.IsNull() && !data.Distance.IsUnknown() {
		distance := data.Distance.ValueFloat64()

		request.Distance = *netbox.NewNullableFloat64(&distance)
	}

	if !data.DistanceUnit.IsNull() && !data.DistanceUnit.IsUnknown() {
		distanceUnit := netbox.PatchedWritableWirelessLinkRequestDistanceUnit(data.DistanceUnit.ValueString())

		request.DistanceUnit = &distanceUnit
	}

	// Store plan values for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Apply common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, request, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the wireless link

	result, httpResp, err := r.client.WirelessAPI.WirelessWirelessLinksCreate(ctx).
		WritableWirelessLinkRequest(*request).Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error Creating Wireless Link",

			utils.FormatAPIError("create wireless link", err, httpResp))

		return
	}

	// Map the response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only
	data.Tags = utils.PopulateTagsFromAPI(ctx, result.HasTags(), result.GetTags(), planTags, &resp.Diagnostics)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, result.GetCustomFields(), &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created wireless link", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.

func (r *WirelessLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WirelessLinkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))

		return
	}

	result, httpResp, err := r.client.WirelessAPI.WirelessWirelessLinksRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Wireless link not found, removing from state", map[string]interface{}{"id": id})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError("Error Reading Wireless Link",

			utils.FormatAPIError(fmt.Sprintf("read wireless link ID %d", id), err, httpResp))

		return
	}

	// Store state values for filter-to-owned pattern (preserve null vs empty set)
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only
	// Only populate tags if they were in the state (owned by config)
	if utils.IsSet(stateTags) {
		data.Tags = utils.PopulateTagsFromAPI(ctx, result.HasTags(), result.GetTags(), stateTags, &resp.Diagnostics)
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, result.GetCustomFields(), &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.

func (r *WirelessLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan WirelessLinkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", plan.ID.ValueString()))

		return
	}

	tflog.Debug(ctx, "Updating wireless link", map[string]interface{}{
		"id": id,

		"interface_a": plan.InterfaceA.ValueString(),

		"interface_b": plan.InterfaceB.ValueString(),
	})

	// Look up interface A

	interfaceA, diags := r.lookupInterfaceBrief(ctx, plan.InterfaceA.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Look up interface B

	interfaceB, diags := r.lookupInterfaceBrief(ctx, plan.InterfaceB.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request

	request := netbox.NewWritableWirelessLinkRequest(*interfaceA, *interfaceB)

	// Set optional fields (same as Create)

	if !plan.SSID.IsNull() && !plan.SSID.IsUnknown() {
		ssid := plan.SSID.ValueString()

		request.Ssid = &ssid
	}

	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		status := netbox.CableStatusValue(plan.Status.ValueString())

		request.Status = &status
	}

	if !plan.Tenant.IsNull() && !plan.Tenant.IsUnknown() {
		tenantID, err := utils.ParseID(plan.Tenant.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Tenant ID must be a number, got: %s", plan.Tenant.ValueString()))

			return
		}

		tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError("Error Looking Up Tenant",

				utils.FormatAPIError(fmt.Sprintf("lookup tenant ID %d", tenantID), err, httpResp))

			return
		}

		tenantRequest := netbox.BriefTenantRequest{
			Name: tenant.GetName(),

			Slug: tenant.GetSlug(),
		}

		request.Tenant = *netbox.NewNullableBriefTenantRequest(&tenantRequest)
	} else if plan.Tenant.IsNull() {
		request.SetTenantNil()
	}

	if !plan.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
		authType := netbox.AuthenticationType1(plan.AuthType.ValueString())

		request.AuthType = &authType
	} else if plan.AuthType.IsNull() && !state.AuthType.IsNull() {
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["auth_type"] = ""
	}

	if !plan.AuthCipher.IsNull() && !plan.AuthCipher.IsUnknown() {
		authCipher := netbox.AuthenticationCipher(plan.AuthCipher.ValueString())

		request.AuthCipher = &authCipher
	} else if plan.AuthCipher.IsNull() && !state.AuthCipher.IsNull() {
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["auth_cipher"] = ""
	}

	if !plan.AuthPSK.IsNull() && !plan.AuthPSK.IsUnknown() {
		psk := plan.AuthPSK.ValueString()

		request.AuthPsk = &psk
	} else if plan.AuthPSK.IsNull() && !state.AuthPSK.IsNull() {
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["auth_psk"] = ""
	}

	if !plan.Distance.IsNull() && !plan.Distance.IsUnknown() {
		distance := plan.Distance.ValueFloat64()

		request.Distance = *netbox.NewNullableFloat64(&distance)
	} else if plan.Distance.IsNull() && !state.Distance.IsNull() {
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["distance"] = nil
	}

	if !plan.DistanceUnit.IsNull() && !plan.DistanceUnit.IsUnknown() {
		distanceUnit := netbox.PatchedWritableWirelessLinkRequestDistanceUnit(plan.DistanceUnit.ValueString())

		request.DistanceUnit = &distanceUnit
	} else if plan.DistanceUnit.IsNull() && !state.DistanceUnit.IsNull() {
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["distance_unit"] = ""
	}

	// Store plan values before mapping for filter-to-owned pattern
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Apply common fields with merge-aware helpers
	utils.ApplyCommonFieldsWithMerge(ctx, request, plan.Description, plan.Comments, plan.Tags, state.Tags, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the wireless link

	result, httpResp, err := r.client.WirelessAPI.WirelessWirelessLinksUpdate(ctx, id).
		WritableWirelessLinkRequest(*request).Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error Updating Wireless Link",

			utils.FormatAPIError(fmt.Sprintf("update wireless link ID %d", id), err, httpResp))

		return
	}

	// Map the response to plan model

	r.mapToState(ctx, result, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Populate tags and custom fields filtered to owned fields only
	// Only populate tags if they were in the plan (owned by config)
	if utils.IsSet(planTags) {
		plan.Tags = utils.PopulateTagsFromAPI(ctx, result.HasTags(), result.GetTags(), planTags, &resp.Diagnostics)
	} else {
		plan.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, result.GetCustomFields(), &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.

func (r *WirelessLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WirelessLinkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))

		return
	}

	tflog.Debug(ctx, "Deleting wireless link", map[string]interface{}{"id": id})

	httpResp, err := r.client.WirelessAPI.WirelessWirelessLinksDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		// If the resource was already deleted (404), consider it a success
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Wireless link already deleted", map[string]interface{}{"id": id})
			return
		}

		resp.Diagnostics.AddError("Error Deleting Wireless Link",

			utils.FormatAPIError(fmt.Sprintf("delete wireless link ID %d", id), err, httpResp))

		return
	}
}

// ImportState imports the resource state.

func (r *WirelessLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapToState maps the API response to the Terraform state.

func (r *WirelessLinkResource) mapToState(ctx context.Context, result *netbox.WirelessLink, data *WirelessLinkResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map interface IDs - on first read (unknown/null), set to ID; otherwise preserve current value to avoid drift

	interfaceA := result.GetInterfaceA()

	if data.InterfaceA.IsUnknown() || data.InterfaceA.IsNull() {
		data.InterfaceA = types.StringValue(fmt.Sprintf("%d", interfaceA.GetId()))
	}

	interfaceB := result.GetInterfaceB()

	if data.InterfaceB.IsUnknown() || data.InterfaceB.IsNull() {
		data.InterfaceB = types.StringValue(fmt.Sprintf("%d", interfaceB.GetId()))
	}

	// Map optional fields

	if result.HasSsid() && result.GetSsid() != "" {
		data.SSID = types.StringValue(result.GetSsid())
	} else {
		data.SSID = types.StringNull()
	}

	if result.HasStatus() {
		status := result.GetStatus()

		data.Status = types.StringValue(string(status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()

		userTenant := data.Tenant.ValueString()

		if userTenant == tenant.GetName() || userTenant == tenant.GetSlug() || userTenant == tenant.GetDisplay() || userTenant == fmt.Sprintf("%d", tenant.GetId()) {
			// Keep user's original value
		} else {
			data.Tenant = types.StringValue(tenant.GetName())
		}
	} else {
		data.Tenant = types.StringNull()
	}

	if result.HasAuthType() {
		authType := result.GetAuthType()

		data.AuthType = types.StringValue(string(authType.GetValue()))
	} else {
		data.AuthType = types.StringNull()
	}

	if result.HasAuthCipher() {
		authCipher := result.GetAuthCipher()

		data.AuthCipher = types.StringValue(string(authCipher.GetValue()))
	} else {
		data.AuthCipher = types.StringNull()
	}

	if result.HasAuthPsk() && result.GetAuthPsk() != "" {
		data.AuthPSK = types.StringValue(result.GetAuthPsk())
	} else {
		data.AuthPSK = types.StringNull()
	}

	if result.HasDistance() {
		distance, ok := result.GetDistanceOk()

		if ok && distance != nil {
			data.Distance = types.Float64Value(*distance)
		} else {
			data.Distance = types.Float64Null()
		}
	} else {
		data.Distance = types.Float64Null()
	}

	if result.HasDistanceUnit() {
		distanceUnit := result.GetDistanceUnit()

		if distanceUnit.Value != nil && *distanceUnit.Value != "" {
			data.DistanceUnit = types.StringValue(string(*distanceUnit.Value))
		} else {
			data.DistanceUnit = types.StringNull()
		}
	} else {
		data.DistanceUnit = types.StringNull()
	}

	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	if result.HasComments() && result.GetComments() != "" {
		data.Comments = types.StringValue(result.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Note: Tags and custom_fields are populated using filter-to-owned pattern in Create(), Read(), and Update()
}
