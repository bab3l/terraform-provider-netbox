// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &VRFResource{}

var _ resource.ResourceWithImportState = &VRFResource{}

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

			"tenant": nbschema.IDOnlyReferenceAttribute("tenant", "ID of the tenant this VRF belongs to."),

			"enforce_unique": schema.BoolAttribute{

				MarkdownDescription: "Prevent duplicate prefixes/IP addresses within this VRF. Defaults to `true`.",

				Optional: true,

				Computed: true,
			},

			"description": nbschema.DescriptionAttribute("VRF"),

			"comments": nbschema.CommentsAttribute("VRF"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

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

	// Set optional fields

	r.setOptionalFields(ctx, &vrfRequest, &data, &resp.Diagnostics)

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

	// Map response back to state

	r.mapVRFToState(ctx, vrf, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *VRFResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data VRFResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	tflog.Debug(ctx, "Updating VRF", map[string]interface{}{

		"id": id,

		"name": data.Name.ValueString(),
	})

	// Prepare the VRF request

	vrfRequest := netbox.VRFRequest{

		Name: data.Name.ValueString(),
	}

	// Set optional fields

	r.setOptionalFields(ctx, &vrfRequest, &data, &resp.Diagnostics)

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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// setOptionalFields sets optional fields on the VRF request from the resource model.

func (r *VRFResource) setOptionalFields(ctx context.Context, vrfRequest *netbox.VRFRequest, data *VRFResourceModel, diags *diag.Diagnostics) {

	// Route distinguisher

	if utils.IsSet(data.RD) {

		rdValue := data.RD.ValueString()

		vrfRequest.Rd = *netbox.NewNullableString(&rdValue)

	}

	// Tenant

	if utils.IsSet(data.Tenant) {

		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		diags.Append(tenantDiags...)

		if diags.HasError() {

			return

		}

		vrfRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)

	}

	// Enforce unique

	if utils.IsSet(data.EnforceUnique) {

		vrfRequest.EnforceUnique = utils.BoolPtr(data.EnforceUnique)

	}

	// Description

	vrfRequest.Description = utils.StringPtr(data.Description)

	// Comments

	vrfRequest.Comments = utils.StringPtr(data.Comments)

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		vrfRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if diags.HasError() {

			return

		}

		vrfRequest.CustomFields = utils.CustomFieldsToMap(customFields)

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

	if vrf.HasTenant() && vrf.Tenant.Get() != nil {

		if data.Tenant.IsNull() || data.Tenant.IsUnknown() {

			data.Tenant = types.StringValue(fmt.Sprintf("%d", vrf.Tenant.Get().GetId()))

		}

		// Otherwise keep the original value the user provided

	} else {

		data.Tenant = types.StringNull()

	}

	// Enforce unique - default is true

	data.EnforceUnique = types.BoolValue(vrf.GetEnforceUnique())

	// Description

	if desc, ok := vrf.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Comments

	if comments, ok := vrf.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Tags

	if vrf.HasTags() {

		tags := utils.NestedTagsToTagModels(vrf.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Custom fields

	if vrf.HasCustomFields() && !data.CustomFields.IsNull() {

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		customFields := utils.MapToCustomFieldModels(vrf.GetCustomFields(), stateCustomFields)

		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfValueDiags...)

		if diags.HasError() {

			return

		}

		data.CustomFields = customFieldsValue

	} else if data.CustomFields.IsNull() {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
