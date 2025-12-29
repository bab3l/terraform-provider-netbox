// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

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

var _ resource.Resource = &TunnelGroupResource{}

var _ resource.ResourceWithImportState = &TunnelGroupResource{}

// NewTunnelGroupResource creates a new TunnelGroupResource.

func NewTunnelGroupResource() resource.Resource {
	return &TunnelGroupResource{}
}

// TunnelGroupResource defines the resource implementation.

type TunnelGroupResource struct {
	client *netbox.APIClient
}

// TunnelGroupResourceModel describes the resource data model.

type TunnelGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *TunnelGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunnel_group"
}

// Schema defines the schema for the resource.

func (r *TunnelGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tunnel group in Netbox. Tunnel groups are used to organize VPN tunnels by function, location, or other criteria.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("tunnel group"),

			"name": nbschema.NameAttribute("tunnel group", 100),

			"slug": nbschema.SlugAttribute("tunnel group"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("tunnel group"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure configures the resource with the provider client.

func (r *TunnelGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new tunnel group resource.

func (r *TunnelGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TunnelGroupResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating tunnel group", map[string]interface{}{
		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Create the tunnel group request

	tunnelGroupRequest := netbox.TunnelGroupRequest{
		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Set optional fields

	utils.ApplyDescription(&tunnelGroupRequest, data.Description)

	utils.ApplyMetadataFields(ctx, &tunnelGroupRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create via API

	tunnelGroup, httpResp, err := r.client.VpnAPI.VpnTunnelGroupsCreate(ctx).TunnelGroupRequest(tunnelGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_tunnel_group",

			ResourceName: "this.tunnel_group",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				list, _, lookupErr := r.client.VpnAPI.VpnTunnelGroupsList(lookupCtx).Slug([]string{slug}).Execute()

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

	// Map response to state

	r.mapTunnelGroupToState(ctx, tunnelGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created tunnel group", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the tunnel group resource.

func (r *TunnelGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TunnelGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Error parsing ID",

			fmt.Sprintf("Could not parse tunnel group ID %s: %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Reading tunnel group", map[string]interface{}{
		"id": id,
	})

	tunnelGroup, httpResp, err := r.client.VpnAPI.VpnTunnelGroupsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading tunnel group",

			utils.FormatAPIError(fmt.Sprintf("read tunnel group ID %d", id), err, httpResp),
		)

		return
	}

	r.mapTunnelGroupToState(ctx, tunnelGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the tunnel group resource.

func (r *TunnelGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TunnelGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Error parsing ID",

			fmt.Sprintf("Could not parse tunnel group ID %s: %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Updating tunnel group", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Create the tunnel group request

	tunnelGroupRequest := netbox.TunnelGroupRequest{
		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Set optional fields

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()

		tunnelGroupRequest.Description = &description
	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {
			return
		}

		tunnelGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Apply metadata fields (tags, custom_fields)

	utils.ApplyMetadataFields(ctx, &tunnelGroupRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update via API

	tunnelGroup, httpResp, err := r.client.VpnAPI.VpnTunnelGroupsUpdate(ctx, id).TunnelGroupRequest(tunnelGroupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating tunnel group",

			utils.FormatAPIError(fmt.Sprintf("update tunnel group ID %d", id), err, httpResp),
		)

		return
	}

	r.mapTunnelGroupToState(ctx, tunnelGroup, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated tunnel group", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the tunnel group resource.

func (r *TunnelGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TunnelGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Error parsing ID",

			fmt.Sprintf("Could not parse tunnel group ID %s: %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting tunnel group", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.VpnAPI.VpnTunnelGroupsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error deleting tunnel group",

			utils.FormatAPIError(fmt.Sprintf("delete tunnel group ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted tunnel group", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing tunnel group resource.

func (r *TunnelGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapTunnelGroupToState maps a TunnelGroup API response to the Terraform state model.

func (r *TunnelGroupResource) mapTunnelGroupToState(ctx context.Context, tunnelGroup *netbox.TunnelGroup, data *TunnelGroupResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tunnelGroup.Id))

	data.Name = types.StringValue(tunnelGroup.Name)

	data.Slug = types.StringValue(tunnelGroup.Slug)

	// Description

	if tunnelGroup.Description != nil && *tunnelGroup.Description != "" {
		data.Description = types.StringValue(*tunnelGroup.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Map display_name
	if tunnelGroup.Display != "" {
	} else {
	}

	// Tags

	if len(tunnelGroup.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(tunnelGroup.Tags)

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

	if tunnelGroup.CustomFields != nil && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(tunnelGroup.CustomFields, stateCustomFields)

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
