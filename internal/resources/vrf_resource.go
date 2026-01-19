// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &VRFResource{}

var _ resource.ResourceWithImportState = &VRFResource{}
var _ resource.ResourceWithIdentity = &VRFResource{}

func NewVRFResource() resource.Resource {
	return &VRFResource{}
}

// VRFResource defines the resource implementation.

type VRFResource struct {
	client *netbox.APIClient
}

// VRFResourceModel describes the resource data model.

type VRFResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	RD types.String `tfsdk:"rd"`

	Tenant types.String `tfsdk:"tenant"`

	EnforceUnique types.Bool `tfsdk:"enforce_unique"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *VRFResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrf"
}

func (r *VRFResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a VRF (Virtual Routing and Forwarding) instance in Netbox. VRFs are used to implement virtual routing tables, enabling multiple routing instances to coexist within the same network infrastructure. They are commonly used to provide network isolation for different customers, departments, or network functions.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("VRF"),

			"name": nbschema.NameAttribute("VRF", 100),

			"rd": schema.StringAttribute{
				MarkdownDescription: "Route distinguisher (RD) as defined in RFC 4364. Format: `ASN:nn` or `IP:nn`. Example: `65000:1` or `192.168.1.1:1`.",

				Optional: true,
			},

			"tenant": nbschema.ReferenceAttribute("tenant", "Name, Slug, or ID of the tenant this VRF belongs to."),

			"enforce_unique": schema.BoolAttribute{
				MarkdownDescription: "Prevent duplicate prefixes/IP addresses within this VRF. Defaults to `true`.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(true),
			},

			"description": nbschema.DescriptionAttribute("VRF"),

			"comments": nbschema.CommentsAttribute("VRF"),

			"tags": nbschema.TagsSlugAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *VRFResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *VRFResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VRFResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VRFResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating VRF", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Prepare the VRF request

	vrfRequest := netbox.VRFRequest{
		Name: data.Name.ValueString(),
	}

	// Set optional fields (pass nil for state since this is Create)

	r.setOptionalFields(ctx, &vrfRequest, &data, nil, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the VRF via API

	vrf, httpResp, err := r.client.IpamAPI.IpamVrfsCreate(ctx).VRFRequest(vrfRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_vrf",

			ResourceName: "this.vrf",

			SlugValue: "", // VRF doesn't have slug

			LookupFunc: nil,
		}

		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)

		return
	}

	tflog.Debug(ctx, "Created VRF", map[string]interface{}{
		"id": vrf.GetId(),

		"name": vrf.GetName(),
	})

	// Map response back to state

	r.mapVRFToState(ctx, vrf, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VRFResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VRFResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	vrfID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vrfID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VRF ID", fmt.Sprintf("VRF ID must be a number, got: %s", vrfID))

		return
	}

	tflog.Debug(ctx, "Reading VRF", map[string]interface{}{
		"id": id,
	})

	// Read the VRF via API

	vrf, httpResp, err := r.client.IpamAPI.IpamVrfsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "VRF not found, removing from state", map[string]interface{}{
				"id": id,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading VRF",

			utils.FormatAPIError(fmt.Sprintf("read VRF ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Read VRF", map[string]interface{}{
		"id": vrf.GetId(),

		"name": vrf.GetName(),
	})

	// Save original custom_fields state before mapping
	originalCustomFields := data.CustomFields

	// Map response back to state

	r.mapVRFToState(ctx, vrf, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields state if it was null or empty
	// This prevents unmanaged/cleared fields from reappearing in state
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VRFResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VRFResourceModel
	var state VRFResourceModel

	// Read both plan and state for merge-aware custom fields handling
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan as the data source
	data := plan

	// Parse the ID

	vrfID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vrfID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VRF ID", fmt.Sprintf("VRF ID must be a number, got: %s", vrfID))

		return
	}

	tflog.Debug(ctx, "Updating VRF", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Prepare the VRF request

	vrfRequest := netbox.VRFRequest{
		Name: data.Name.ValueString(),
	}

	// Set optional fields (pass state for merge-aware custom fields handling)

	r.setOptionalFields(ctx, &vrfRequest, &data, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the VRF via API

	vrf, httpResp, err := r.client.IpamAPI.IpamVrfsUpdate(ctx, id).VRFRequest(vrfRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating VRF",

			utils.FormatAPIError(fmt.Sprintf("update VRF ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Updated VRF", map[string]interface{}{
		"id": vrf.GetId(),

		"name": vrf.GetName(),
	})

	// Map response back to state

	r.mapVRFToState(ctx, vrf, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VRFResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VRFResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	vrfID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vrfID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VRF ID", fmt.Sprintf("VRF ID must be a number, got: %s", vrfID))

		return
	}

	tflog.Debug(ctx, "Deleting VRF", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Delete the VRF via API

	httpResp, err := r.client.IpamAPI.IpamVrfsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "VRF already deleted", map[string]interface{}{
				"id": id,
			})

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting VRF",

			utils.FormatAPIError(fmt.Sprintf("delete VRF ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted VRF", map[string]interface{}{
		"id": id,
	})
}

func (r *VRFResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			resp.Diagnostics.AddError("Invalid VRF ID", fmt.Sprintf("VRF ID must be a number, got: %s", parsed.ID))
			return
		}
		vrf, httpResp, err := r.client.IpamAPI.IpamVrfsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing VRF", utils.FormatAPIError(fmt.Sprintf("read VRF ID %d", id), err, httpResp))
			return
		}

		var data VRFResourceModel
		if vrf.HasTenant() && vrf.Tenant.Get() != nil && vrf.Tenant.Get().GetId() != 0 {
			tenant := vrf.Tenant.Get()
			data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		}
		if vrf.HasTags() && len(vrf.GetTags()) > 0 {
			tagSlugs := make([]string, 0, len(vrf.GetTags()))
			for _, tag := range vrf.GetTags() {
				tagSlugs = append(tagSlugs, tag.Slug)
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

		r.mapVRFToState(ctx, vrf, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, vrf.GetCustomFields(), &resp.Diagnostics)
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

// setOptionalFields sets optional fields on the VRF request from the resource model.

func (r *VRFResource) setOptionalFields(ctx context.Context, vrfRequest *netbox.VRFRequest, data *VRFResourceModel, state *VRFResourceModel, diags *diag.Diagnostics) {
	// Route distinguisher

	if utils.IsSet(data.RD) {
		rdValue := data.RD.ValueString()

		vrfRequest.Rd = *netbox.NewNullableString(&rdValue)
	} else if state != nil && data.RD.IsNull() {
		// NetBox PATCH semantics: omitting a field does not clear it.
		// When removed from config during Update, explicitly clear it.
		vrfRequest.SetRdNil()
	}

	// Tenant

	if utils.IsSet(data.Tenant) {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		diags.Append(tenantDiags...)

		if diags.HasError() {
			return
		}

		vrfRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	} else if data.Tenant.IsNull() {
		// Explicitly set to null to clear the field
		vrfRequest.SetTenantNil()
	}

	// Enforce unique

	if utils.IsSet(data.EnforceUnique) {
		vrfRequest.EnforceUnique = utils.BoolPtr(data.EnforceUnique)
	}

	// Description
	utils.ApplyDescription(vrfRequest, data.Description)

	// Comments
	utils.ApplyComments(vrfRequest, data.Comments)

	// Tags
	utils.ApplyTagsFromSlugs(ctx, r.client, vrfRequest, data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Custom fields with merge-aware handling
	if state != nil {
		// Update operation - merge custom fields to preserve unmanaged fields
		utils.ApplyCustomFieldsWithMerge(ctx, vrfRequest, data.CustomFields, state.CustomFields, diags)
	} else {
		// Create operation - apply custom fields directly
		utils.ApplyCustomFields(ctx, vrfRequest, data.CustomFields, diags)
	}
	if diags.HasError() {
		return
	}
}

// mapVRFToState maps a VRF API response to the resource model.

func (r *VRFResource) mapVRFToState(ctx context.Context, vrf *netbox.VRF, data *VRFResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vrf.GetId()))

	data.Name = types.StringValue(vrf.GetName())

	// Route distinguisher
	if rd, ok := vrf.GetRdOk(); ok && rd != nil && *rd != "" {
		data.RD = types.StringValue(*rd)
	} else {
		data.RD = types.StringNull()
	}

	// Tenant
	tenantRef := utils.PreserveOptionalReferenceWithID(
		data.Tenant,
		vrf.HasTenant() && vrf.Tenant.Get() != nil,
		vrf.Tenant.Get().GetId(),
		vrf.Tenant.Get().GetName(),
		vrf.Tenant.Get().GetSlug(),
	)
	data.Tenant = tenantRef.Reference

	// Enforce unique - default is true
	data.EnforceUnique = types.BoolValue(vrf.GetEnforceUnique())

	// Description
	data.Description = utils.NullableStringFromAPI(
		vrf.HasDescription() && vrf.GetDescription() != "",
		vrf.GetDescription,
		data.Description,
	)

	// Comments
	data.Comments = utils.NullableStringFromAPI(
		vrf.HasComments() && vrf.GetComments() != "",
		vrf.GetComments,
		data.Comments,
	)

	// Tags (slug list) with empty-set preservation
	wasExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
	switch {
	case vrf.HasTags() && len(vrf.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(vrf.GetTags()))
		for _, tag := range vrf.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Custom fields
	// Custom Fields - filter to owned fields only
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, vrf.GetCustomFields(), diags)
}
