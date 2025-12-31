// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

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

var _ resource.Resource = &VLANResource{}

var _ resource.ResourceWithImportState = &VLANResource{}

func NewVLANResource() resource.Resource {
	return &VLANResource{}
}

// VLANResource defines the resource implementation.

type VLANResource struct {
	client *netbox.APIClient
}

// VLANResourceModel describes the resource data model.

type VLANResourceModel struct {
	ID types.String `tfsdk:"id"`

	VID types.Int64 `tfsdk:"vid"`

	Name types.String `tfsdk:"name"`

	Site types.String `tfsdk:"site"`

	Group types.String `tfsdk:"group"`

	Tenant types.String `tfsdk:"tenant"`

	Status types.String `tfsdk:"status"`

	Role types.String `tfsdk:"role"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *VLANResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vlan"
}

func (r *VLANResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a VLAN in Netbox. VLANs (Virtual Local Area Networks) represent layer 2 broadcast domains that can be assigned to interfaces and used to organize network resources.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("VLAN"),

			"vid": schema.Int64Attribute{
				MarkdownDescription: "VLAN ID (1-4094). This is the numeric identifier used on network devices. Required.",

				Required: true,
			},

			"name": nbschema.NameAttribute("VLAN", 64),

			"site": nbschema.ReferenceAttribute("site", "ID or slug of the site this VLAN belongs to."),

			"group": schema.StringAttribute{
				MarkdownDescription: "ID or slug of the VLAN group this VLAN belongs to.",

				Optional: true,
			},

			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant this VLAN belongs to."),

			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the VLAN. Valid values: `active`, `reserved`, `deprecated`. Defaults to `active`.",

				Optional: true,

				Computed: true,
			},

			"role": schema.StringAttribute{
				MarkdownDescription: "ID or slug of the role assigned to this VLAN.",

				Optional: true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("VLAN"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *VLANResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VLANResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating VLAN", map[string]interface{}{
		"vid": data.VID.ValueInt64(),

		"name": data.Name.ValueString(),
	})

	// Convert VID to int32 with overflow check

	vid, err := utils.SafeInt32FromValue(data.VID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VLAN ID", fmt.Sprintf("VID overflow: %s", err))

		return
	}

	// Prepare the VLAN request

	vlanRequest := netbox.WritableVLANRequest{
		Vid: vid,

		Name: data.Name.ValueString(),
	}

	// Set optional fields

	r.setOptionalFields(ctx, &vlanRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the VLAN via API

	vlan, httpResp, err := r.client.IpamAPI.IpamVlansCreate(ctx).WritableVLANRequest(vlanRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_vlan",

			ResourceName: "this.vlan",

			SlugValue: "",

			LookupFunc: nil,
		}

		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)

		return
	}

	tflog.Debug(ctx, "Created VLAN", map[string]interface{}{
		"id": vlan.GetId(),

		"vid": vlan.GetVid(),

		"name": vlan.GetName(),
	})

	// Map response back to state

	r.mapVLANToState(ctx, vlan, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VLANResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	vlanID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vlanID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VLAN ID", fmt.Sprintf("VLAN ID must be a number, got: %s", vlanID))

		return
	}

	tflog.Debug(ctx, "Reading VLAN", map[string]interface{}{
		"id": id,
	})

	// Read the VLAN via API

	vlan, httpResp, err := r.client.IpamAPI.IpamVlansRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "VLAN not found, removing from state", map[string]interface{}{
				"id": id,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading VLAN",

			utils.FormatAPIError(fmt.Sprintf("read VLAN ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Read VLAN", map[string]interface{}{
		"id": vlan.GetId(),

		"vid": vlan.GetVid(),

		"name": vlan.GetName(),
	})

	// Map response back to state

	r.mapVLANToState(ctx, vlan, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VLANResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	vlanID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vlanID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VLAN ID", fmt.Sprintf("VLAN ID must be a number, got: %s", vlanID))

		return
	}

	tflog.Debug(ctx, "Updating VLAN", map[string]interface{}{
		"id": id,

		"vid": data.VID.ValueInt64(),

		"name": data.Name.ValueString(),
	})

	// Convert VID to int32 with overflow check

	vid, err := utils.SafeInt32FromValue(data.VID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VLAN ID", fmt.Sprintf("VID overflow: %s", err))

		return
	}

	// Prepare the VLAN request

	vlanRequest := netbox.WritableVLANRequest{
		Vid: vid,

		Name: data.Name.ValueString(),
	}

	// Set optional fields

	r.setOptionalFields(ctx, &vlanRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the VLAN via API

	vlan, httpResp, err := r.client.IpamAPI.IpamVlansUpdate(ctx, id).WritableVLANRequest(vlanRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating VLAN",

			utils.FormatAPIError(fmt.Sprintf("update VLAN ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Updated VLAN", map[string]interface{}{
		"id": vlan.GetId(),

		"vid": vlan.GetVid(),

		"name": vlan.GetName(),
	})

	// Map response back to state

	r.mapVLANToState(ctx, vlan, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VLANResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	vlanID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vlanID)

	if err != nil {
		resp.Diagnostics.AddError("Invalid VLAN ID", fmt.Sprintf("VLAN ID must be a number, got: %s", vlanID))

		return
	}

	tflog.Debug(ctx, "Deleting VLAN", map[string]interface{}{
		"id": id,

		"vid": data.VID.ValueInt64(),

		"name": data.Name.ValueString(),
	})

	// Delete the VLAN via API

	httpResp, err := r.client.IpamAPI.IpamVlansDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "VLAN already deleted", map[string]interface{}{
				"id": id,
			})

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting VLAN",

			utils.FormatAPIError(fmt.Sprintf("delete VLAN ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted VLAN", map[string]interface{}{
		"id": id,
	})
}

func (r *VLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the VLAN request from the resource model.

func (r *VLANResource) setOptionalFields(ctx context.Context, vlanRequest *netbox.WritableVLANRequest, data *VLANResourceModel, diags *diag.Diagnostics) {
	// Site

	if utils.IsSet(data.Site) {
		site, siteDiags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())

		diags.Append(siteDiags...)

		if diags.HasError() {
			return
		}

		vlanRequest.Site = *netbox.NewNullableBriefSiteRequest(site)
	}

	// Group

	if utils.IsSet(data.Group) {
		group, groupDiags := netboxlookup.LookupVLANGroup(ctx, r.client, data.Group.ValueString())

		diags.Append(groupDiags...)

		if diags.HasError() {
			return
		}

		vlanRequest.Group = *netbox.NewNullableBriefVLANGroupRequest(group)
	}

	// Tenant

	if utils.IsSet(data.Tenant) {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		diags.Append(tenantDiags...)

		if diags.HasError() {
			return
		}

		vlanRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	}

	// Status

	if utils.IsSet(data.Status) {
		status := netbox.PatchedWritableVLANRequestStatus(data.Status.ValueString())

		vlanRequest.Status = &status
	}

	// Role

	if utils.IsSet(data.Role) {
		role, roleDiags := netboxlookup.LookupRole(ctx, r.client, data.Role.ValueString())

		diags.Append(roleDiags...)

		if diags.HasError() {
			return
		}

		vlanRequest.Role = *netbox.NewNullableBriefRoleRequest(role)
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, vlanRequest, data.Description, data.Comments, data.Tags, data.CustomFields, diags)
	if diags.HasError() {
		return
	}
}

// mapVLANToState maps a VLAN API response to the resource model.

func (r *VLANResource) mapVLANToState(ctx context.Context, vlan *netbox.VLAN, data *VLANResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vlan.GetId()))

	data.VID = types.Int64Value(int64(vlan.GetVid()))

	data.Name = types.StringValue(vlan.GetName())

	// Site

	if vlan.HasSite() && vlan.Site.Get() != nil {
		data.Site = utils.UpdateReferenceAttribute(data.Site, vlan.Site.Get().Name, vlan.Site.Get().Slug, vlan.Site.Get().Id)
	} else {
		data.Site = types.StringNull()
	}

	// Group - preserve user's input format

	if vlan.HasGroup() && vlan.Group.Get() != nil {
		group := vlan.Group.Get()

		data.Group = utils.UpdateReferenceAttribute(data.Group, group.GetName(), group.GetSlug(), group.GetId())
	} else {
		data.Group = types.StringNull()
	}

	// Tenant

	if vlan.HasTenant() && vlan.Tenant.Get() != nil {
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, vlan.Tenant.Get().Name, vlan.Tenant.Get().Slug, vlan.Tenant.Get().Id)
	} else {
		data.Tenant = types.StringNull()
	}

	// Status - only set if user specified it in config, or during import (when current is unknown)
	// This prevents Terraform from seeing drift when the API returns a status but config doesn't specify one

	// Status - always set since it's computed (defaults to "active")
	if vlan.HasStatus() {
		data.Status = types.StringValue(string(vlan.Status.GetValue()))
	} else {
		data.Status = types.StringValue("active")
	}

	// Role - preserve user's input format

	if vlan.HasRole() && vlan.Role.Get() != nil {
		role := vlan.Role.Get()

		data.Role = utils.UpdateReferenceAttribute(data.Role, role.GetName(), role.GetSlug(), role.GetId())
	} else {
		data.Role = types.StringNull()
	}

	// Description

	if desc, ok := vlan.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Comments

	if comments, ok := vlan.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags

	if vlan.HasTags() {
		tags := utils.NestedTagsToTagModels(vlan.GetTags())

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

	if vlan.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(vlan.GetCustomFields(), stateCustomFields)

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
