// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FHRPGroupResource{}
var _ resource.ResourceWithImportState = &FHRPGroupResource{}

func NewFHRPGroupResource() resource.Resource {
	return &FHRPGroupResource{}
}

// FHRPGroupResource defines the resource implementation.
type FHRPGroupResource struct {
	client *netbox.APIClient
}

// FHRPGroupResourceModel describes the resource data model.
type FHRPGroupResourceModel struct {
	ID           types.Int32  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Protocol     types.String `tfsdk:"protocol"`
	GroupID      types.Int32  `tfsdk:"group_id"`
	AuthType     types.String `tfsdk:"auth_type"`
	AuthKey      types.String `tfsdk:"auth_key"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *FHRPGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fhrp_group"
}

func (r *FHRPGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an FHRP (First Hop Redundancy Protocol) group in Netbox. FHRP groups represent virtual IP configurations for protocols like VRRP, HSRP, CARP, GLBP, and others that provide gateway redundancy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the FHRP group.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the FHRP group.",
				Optional:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The redundancy protocol. Valid values: `vrrp2`, `vrrp3`, `carp`, `clusterxl`, `hsrp`, `glbp`, `other`.",
				Required:            true,
			},
			"group_id": schema.Int32Attribute{
				MarkdownDescription: "The FHRP group identifier (e.g., VRRP group ID, HSRP group number).",
				Required:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type. Valid values: `plaintext`, `md5`, or empty string.",
				Optional:            true,
			},
			"auth_key": schema.StringAttribute{
				MarkdownDescription: "Authentication key/password for the FHRP group.",
				Optional:            true,
				Sensitive:           true,
			},
			"description":   nbschema.DescriptionAttribute("FHRP group"),
			"comments":      nbschema.CommentsAttribute("FHRP group"),
			"tags":          nbschema.TagsSlugAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *FHRPGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FHRPGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FHRPGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating FHRP Group", map[string]interface{}{
		"protocol": data.Protocol.ValueString(),
		"group_id": data.GroupID.ValueInt32(),
	})

	// Prepare the FHRP Group request
	protocol, err := netbox.NewBriefFHRPGroupProtocolFromValue(data.Protocol.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Protocol", fmt.Sprintf("Invalid FHRP protocol value: %s", data.Protocol.ValueString()))
		return
	}
	fhrpGroupRequest := netbox.FHRPGroupRequest{
		Protocol: *protocol,
		GroupId:  data.GroupID.ValueInt32(),
	}

	// Set optional fields
	r.setOptionalFields(ctx, &fhrpGroupRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the FHRP Group via API
	fhrpGroup, httpResp, err := r.client.IpamAPI.IpamFhrpGroupsCreate(ctx).FHRPGroupRequest(fhrpGroupRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating FHRP Group",
			utils.FormatAPIError("creating FHRP group", err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Created FHRP Group", map[string]interface{}{
		"id":       fhrpGroup.GetId(),
		"protocol": string(fhrpGroup.Protocol),
		"group_id": fhrpGroup.GetGroupId(),
	})

	// Map response back to state
	r.mapFHRPGroupToState(ctx, fhrpGroup, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FHRPGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FHRPGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields from state for potential restoration
	originalCustomFields := data.CustomFields

	// Get the FHRP Group via API
	id := data.ID.ValueInt32()
	fhrpGroup, httpResp, err := r.client.IpamAPI.IpamFhrpGroupsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "FHRP Group not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading FHRP Group",
			utils.FormatAPIError("reading FHRP group", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapFHRPGroupToState(ctx, fhrpGroup, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was null or empty before (not managed or explicitly cleared),
	// restore that state after mapping.
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		tflog.Debug(ctx, "Custom fields unmanaged/cleared, preserving original state during Read", map[string]interface{}{
			"was_null":  originalCustomFields.IsNull(),
			"was_empty": !originalCustomFields.IsNull() && len(originalCustomFields.Elements()) == 0,
		})
		data.CustomFields = originalCustomFields
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FHRPGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read BOTH state and plan for merge-aware custom fields
	var state, plan FHRPGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := plan.ID.ValueInt32()
	tflog.Debug(ctx, "Updating FHRP Group", map[string]interface{}{
		"id": id,
	})

	// Prepare the FHRP Group request
	protocol, err := netbox.NewBriefFHRPGroupProtocolFromValue(plan.Protocol.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Protocol", fmt.Sprintf("Invalid FHRP protocol value: %s", plan.Protocol.ValueString()))
		return
	}
	fhrpGroupRequest := netbox.FHRPGroupRequest{
		Protocol: *protocol,
		GroupId:  plan.GroupID.ValueInt32(),
	}

	// Set optional fields with state for merge
	r.setOptionalFields(ctx, &fhrpGroupRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Store plan tags/customfields for filter-to-owned population
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Update the FHRP Group via API
	fhrpGroup, httpResp, err := r.client.IpamAPI.IpamFhrpGroupsUpdate(ctx, id).FHRPGroupRequest(fhrpGroupRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating FHRP Group",
			utils.FormatAPIError("updating FHRP group", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Updated FHRP Group", map[string]interface{}{
		"id": fhrpGroup.GetId(),
	})

	// Map response to state
	plan.Tags = planTags
	plan.CustomFields = planCustomFields
	r.mapFHRPGroupToState(ctx, fhrpGroup, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FHRPGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FHRPGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueInt32()
	tflog.Debug(ctx, "Deleting FHRP Group", map[string]interface{}{
		"id": id,
	})

	// Delete the FHRP Group via API
	httpResp, err := r.client.IpamAPI.IpamFhrpGroupsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting FHRP Group",
			utils.FormatAPIError("deleting FHRP group", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted FHRP Group", map[string]interface{}{
		"id": id,
	})
}

func (r *FHRPGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing FHRP Group",
			fmt.Sprintf("Could not parse ID %q: %s", req.ID, err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// setOptionalFields sets optional fields on the FHRP Group request.
func (r *FHRPGroupResource) setOptionalFields(ctx context.Context, fhrpGroupRequest *netbox.FHRPGroupRequest, plan *FHRPGroupResourceModel, state *FHRPGroupResourceModel, diags *diag.Diagnostics) {
	// Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		fhrpGroupRequest.Name = &name
	} else if plan.Name.IsNull() {
		// Explicitly clear the name
		emptyString := ""
		fhrpGroupRequest.Name = &emptyString
	}

	// Auth Type
	if !plan.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
		authType, err := netbox.NewAuthenticationTypeFromValue(plan.AuthType.ValueString())
		if err != nil {
			diags.AddError("Invalid Auth Type", fmt.Sprintf("Invalid authentication type value: %s", plan.AuthType.ValueString()))
			return
		}
		fhrpGroupRequest.AuthType = authType
	} else if plan.AuthType.IsNull() {
		// Explicitly clear the auth type by setting to empty value
		emptyAuthType := netbox.AUTHENTICATIONTYPE_EMPTY
		fhrpGroupRequest.AuthType = &emptyAuthType
	}

	// Auth Key
	if !plan.AuthKey.IsNull() && !plan.AuthKey.IsUnknown() {
		authKey := plan.AuthKey.ValueString()
		fhrpGroupRequest.AuthKey = &authKey
	}

	// Set common fields with merge-aware custom fields
	utils.ApplyDescription(fhrpGroupRequest, plan.Description)
	utils.ApplyComments(fhrpGroupRequest, plan.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, fhrpGroupRequest, plan.Tags, diags)
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, fhrpGroupRequest, plan.CustomFields, state.CustomFields, diags)
	} else {
		utils.ApplyCustomFields(ctx, fhrpGroupRequest, plan.CustomFields, diags)
	}
	if diags.HasError() {
		return
	}
}

// mapFHRPGroupToState maps an FHRP Group API response to the Terraform state model.
func (r *FHRPGroupResource) mapFHRPGroupToState(ctx context.Context, fhrpGroup *netbox.FHRPGroup, data *FHRPGroupResourceModel, diags *diag.Diagnostics) {
	data.ID = types.Int32Value(fhrpGroup.GetId())
	data.Protocol = types.StringValue(string(fhrpGroup.Protocol))
	data.GroupID = types.Int32Value(fhrpGroup.GetGroupId())

	// Name
	if name := fhrpGroup.GetName(); name != "" {
		data.Name = types.StringValue(name)
	} else {
		data.Name = types.StringNull()
	}

	// Auth Type
	if fhrpGroup.AuthType != nil {
		authTypeValue := string(*fhrpGroup.AuthType)
		if authTypeValue != "" {
			data.AuthType = types.StringValue(authTypeValue)
		} else {
			data.AuthType = types.StringNull()
		}
	} else {
		data.AuthType = types.StringNull()
	}

	// Auth Key - API doesn't return auth_key for security, keep existing value
	// Only set to null if it was null before (import case)
	if data.AuthKey.IsNull() {
		// Keep null
	}

	// Description
	if description := fhrpGroup.GetDescription(); description != "" {
		data.Description = types.StringValue(description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if comments := fhrpGroup.GetComments(); comments != "" {
		data.Comments = types.StringValue(comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags using slug list handling
	wasExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
	switch {
	case fhrpGroup.HasTags() && len(fhrpGroup.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(fhrpGroup.GetTags()))
		for _, tag := range fhrpGroup.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields using filter-to-owned helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, fhrpGroup.GetCustomFields(), diags)
}
