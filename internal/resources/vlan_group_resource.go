// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
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

var _ resource.Resource = &VLANGroupResource{}

var _ resource.ResourceWithImportState = &VLANGroupResource{}

func NewVLANGroupResource() resource.Resource {

	return &VLANGroupResource{}

}

// VLANGroupResource defines the resource implementation.

type VLANGroupResource struct {
	client *netbox.APIClient
}

// VLANGroupResourceModel describes the resource data model.

type VLANGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	DisplayName types.String `tfsdk:"display_name"`

	ScopeType types.String `tfsdk:"scope_type"`

	ScopeID types.String `tfsdk:"scope_id"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *VLANGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_vlan_group"

}

func (r *VLANGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a VLAN Group in Netbox. VLAN Groups are used to organize VLANs and ensure VLAN ID uniqueness within a specific scope such as a site, location, or region.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("VLAN Group"),

			"name": nbschema.NameAttribute("VLAN Group", 100),

			"slug": nbschema.SlugAttribute("VLAN Group"),

			"display_name": nbschema.DisplayNameAttribute("VLAN Group"),

			"description": nbschema.DescriptionAttribute("VLAN Group"),

			"scope_type": schema.StringAttribute{

				MarkdownDescription: "The type of object to scope this VLAN Group to. Valid values: `dcim.site`, `dcim.sitegroup`, `dcim.region`, `dcim.location`, `dcim.rack`, `virtualization.clustergroup`, `virtualization.cluster`.",

				Optional: true,
			},

			"scope_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the object to scope this VLAN Group to. Must be used together with `scope_type`.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *VLANGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *VLANGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data VLANGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating VLAN Group", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Prepare the VLAN Group request

	vlanGroupRequest := netbox.VLANGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Set optional fields

	r.setOptionalFields(ctx, &vlanGroupRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Create the VLAN Group via API

	vlanGroup, httpResp, err := r.client.IpamAPI.IpamVlanGroupsCreate(ctx).VLANGroupRequest(vlanGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		handler := utils.CreateErrorHandler{

			ResourceType: "netbox_vlan_group",

			ResourceName: "this.vlan_group",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {

				list, _, lookupErr := r.client.IpamAPI.IpamVlanGroupsList(lookupCtx).Slug([]string{slug}).Execute()

				if lookupErr != nil {

					return "", lookupErr

				}

				if list != nil && len(list.Results) > 0 {

					return fmt.Sprintf("%d", list.Results[0].GetId()), nil

				}

				return "", nil

			},
		}

		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)

		return

	}

	tflog.Debug(ctx, "Created VLAN Group", map[string]interface{}{

		"id": vlanGroup.GetId(),

		"name": vlanGroup.GetName(),
	})

	// Map response back to state

	r.mapVLANGroupToState(ctx, vlanGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *VLANGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data VLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	vlanGroupID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vlanGroupID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid VLAN Group ID", fmt.Sprintf("VLAN Group ID must be a number, got: %s", vlanGroupID))

		return

	}

	tflog.Debug(ctx, "Reading VLAN Group", map[string]interface{}{

		"id": id,
	})

	// Read the VLAN Group via API

	vlanGroup, httpResp, err := r.client.IpamAPI.IpamVlanGroupsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "VLAN Group not found, removing from state", map[string]interface{}{

				"id": id,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading VLAN Group",

			utils.FormatAPIError(fmt.Sprintf("read VLAN Group ID %d", id), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Read VLAN Group", map[string]interface{}{

		"id": vlanGroup.GetId(),

		"name": vlanGroup.GetName(),
	})

	// Map response back to state

	r.mapVLANGroupToState(ctx, vlanGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *VLANGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data VLANGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	vlanGroupID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vlanGroupID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid VLAN Group ID", fmt.Sprintf("VLAN Group ID must be a number, got: %s", vlanGroupID))

		return

	}

	tflog.Debug(ctx, "Updating VLAN Group", map[string]interface{}{

		"id": id,

		"name": data.Name.ValueString(),
	})

	// Prepare the VLAN Group request

	vlanGroupRequest := netbox.VLANGroupRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Set optional fields

	r.setOptionalFields(ctx, &vlanGroupRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Update the VLAN Group via API

	vlanGroup, httpResp, err := r.client.IpamAPI.IpamVlanGroupsUpdate(ctx, id).VLANGroupRequest(vlanGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating VLAN Group",

			utils.FormatAPIError(fmt.Sprintf("update VLAN Group ID %d", id), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Updated VLAN Group", map[string]interface{}{

		"id": vlanGroup.GetId(),

		"name": vlanGroup.GetName(),
	})

	// Preserve display_name from current state since it's computed and managed by Terraform
	displayNameBeforeMapping := data.DisplayName

	// Map response back to state

	r.mapVLANGroupToState(ctx, vlanGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Restore the previous display_name to avoid inconsistent result error
	// (display_name is computed, so Terraform handles its updates separately)
	data.DisplayName = displayNameBeforeMapping

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *VLANGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data VLANGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	vlanGroupID := data.ID.ValueString()

	var id int32

	id, err := utils.ParseID(vlanGroupID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid VLAN Group ID", fmt.Sprintf("VLAN Group ID must be a number, got: %s", vlanGroupID))

		return

	}

	tflog.Debug(ctx, "Deleting VLAN Group", map[string]interface{}{

		"id": id,

		"name": data.Name.ValueString(),
	})

	// Delete the VLAN Group via API

	httpResp, err := r.client.IpamAPI.IpamVlanGroupsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "VLAN Group already deleted", map[string]interface{}{

				"id": id,
			})

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting VLAN Group",

			utils.FormatAPIError(fmt.Sprintf("delete VLAN Group ID %d", id), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted VLAN Group", map[string]interface{}{

		"id": id,
	})

}

func (r *VLANGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// setOptionalFields sets optional fields on the VLAN Group request from the resource model.

func (r *VLANGroupResource) setOptionalFields(ctx context.Context, vlanGroupRequest *netbox.VLANGroupRequest, data *VLANGroupResourceModel, diags *diag.Diagnostics) {

	// Description

	vlanGroupRequest.Description = utils.StringPtr(data.Description)

	// Scope type

	if utils.IsSet(data.ScopeType) {

		scopeType := data.ScopeType.ValueString()

		vlanGroupRequest.ScopeType = *netbox.NewNullableString(&scopeType)

	}

	// Scope ID

	if utils.IsSet(data.ScopeID) {

		scopeID, err := utils.ParseID(data.ScopeID.ValueString())

		if err != nil {

			diags.AddError("Invalid Scope ID", fmt.Sprintf("Scope ID must be a number, got: %s", data.ScopeID.ValueString()))

			return

		}

		vlanGroupRequest.ScopeId = *netbox.NewNullableInt32(&scopeID)

	}

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		vlanGroupRequest.Tags = tags

	}

	// Handle custom fields

	if utils.IsSet(data.CustomFields) {

		var customFields []utils.CustomFieldModel

		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if diags.HasError() {

			return

		}

		vlanGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

}

// mapVLANGroupToState maps a VLANGroup API response to the resource model.

func (r *VLANGroupResource) mapVLANGroupToState(ctx context.Context, vlanGroup *netbox.VLANGroup, data *VLANGroupResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", vlanGroup.GetId()))

	data.Name = types.StringValue(vlanGroup.GetName())

	data.Slug = types.StringValue(vlanGroup.GetSlug())

	data.DisplayName = types.StringValue(vlanGroup.GetDisplay())

	// Scope type

	if scopeType, ok := vlanGroup.GetScopeTypeOk(); ok && scopeType != nil && *scopeType != "" {

		data.ScopeType = types.StringValue(*scopeType)

	} else {

		data.ScopeType = types.StringNull()

	}

	// Scope ID

	if vlanGroup.HasScopeId() && vlanGroup.ScopeId.Get() != nil {

		data.ScopeID = types.StringValue(fmt.Sprintf("%d", *vlanGroup.ScopeId.Get()))

	} else {

		data.ScopeID = types.StringNull()

	}

	// Description

	if desc, ok := vlanGroup.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Tags

	if vlanGroup.HasTags() {

		tags := utils.NestedTagsToTagModels(vlanGroup.GetTags())

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

	if vlanGroup.HasCustomFields() && !data.CustomFields.IsNull() {

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		customFields := utils.MapToCustomFieldModels(vlanGroup.GetCustomFields(), stateCustomFields)

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
