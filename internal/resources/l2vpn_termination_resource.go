// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &L2VPNTerminationResource{}
	_ resource.ResourceWithImportState = &L2VPNTerminationResource{}
	_ resource.ResourceWithIdentity    = &L2VPNTerminationResource{}
)

func NewL2VPNTerminationResource() resource.Resource {
	return &L2VPNTerminationResource{}
}

// L2VPNTerminationResource defines the resource implementation.
type L2VPNTerminationResource struct {
	client *netbox.APIClient
}

// L2VPNTerminationResourceModel describes the resource data model.
type L2VPNTerminationResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	L2VPN              types.String `tfsdk:"l2vpn"`
	AssignedObjectType types.String `tfsdk:"assigned_object_type"`
	AssignedObjectID   types.Int64  `tfsdk:"assigned_object_id"`
	Tags               types.Set    `tfsdk:"tags"`
	CustomFields       types.Set    `tfsdk:"custom_fields"`
}

func (r *L2VPNTerminationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn_termination"
}

func (r *L2VPNTerminationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Layer 2 VPN termination in Netbox. L2VPN terminations associate L2VPNs with interfaces or VLANs.",
		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("L2VPN termination"),
			"l2vpn": schema.StringAttribute{
				MarkdownDescription: "ID of the L2VPN this termination belongs to.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer ID",
					),
				},
			},
			"assigned_object_type": schema.StringAttribute{
				MarkdownDescription: "Content type of the assigned object. Valid values: `dcim.interface`, `ipam.vlan`, `virtualization.vminterface`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"dcim.interface",
						"ipam.vlan",
						"virtualization.vminterface",
					),
				},
			},
			"assigned_object_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the assigned object (interface or VLAN).",
				Required:            true,
			},
		},
	}

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *L2VPNTerminationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *L2VPNTerminationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *L2VPNTerminationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data L2VPNTerminationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Look up L2VPN to get name and slug for BriefL2VPNRequest
	var l2vpnID int32
	if _, err := fmt.Sscanf(data.L2VPN.ValueString(), "%d", &l2vpnID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid L2VPN ID format",
			fmt.Sprintf("Could not parse L2VPN ID '%s': %s", data.L2VPN.ValueString(), err.Error()),
		)
		return
	}
	l2vpn, httpResp, err := r.client.VpnAPI.VpnL2vpnsRetrieve(ctx, l2vpnID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error looking up L2VPN",
			utils.FormatAPIError("lookup L2VPN", err, httpResp),
		)
		return
	}

	// Build the API request
	l2vpnRef := netbox.NewBriefL2VPNRequest(l2vpn.GetName(), l2vpn.GetSlug())
	terminationRequest := netbox.NewL2VPNTerminationRequest(
		*l2vpnRef,
		data.AssignedObjectType.ValueString(),
		data.AssignedObjectID.ValueInt64(),
	)

	// Apply metadata fields (tags, custom_fields)
	utils.ApplyTagsFromSlugs(ctx, r.client, terminationRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, terminationRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating L2VPN termination", map[string]interface{}{
		"l2vpn_id":             l2vpnID,
		"assigned_object_type": data.AssignedObjectType.ValueString(),
		"assigned_object_id":   data.AssignedObjectID.ValueInt64(),
	})

	// Call the API
	termination, httpResp, err := r.client.VpnAPI.VpnL2vpnTerminationsCreate(ctx).L2VPNTerminationRequest(*terminationRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating L2VPN termination",
			utils.FormatAPIError("create L2VPN termination", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToState(ctx, termination, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Created L2VPN termination resource", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNTerminationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data L2VPNTerminationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var idInt int32
	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse L2VPN termination ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Read from API
	termination, httpResp, err := r.client.VpnAPI.VpnL2vpnTerminationsRetrieve(ctx, idInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading L2VPN termination",
			utils.FormatAPIError("read L2VPN termination", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToState(ctx, termination, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNTerminationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data L2VPNTerminationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var idInt int32
	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse L2VPN termination ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Look up L2VPN to get name and slug for BriefL2VPNRequest
	var l2vpnID int32
	if _, err := fmt.Sscanf(data.L2VPN.ValueString(), "%d", &l2vpnID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid L2VPN ID format",
			fmt.Sprintf("Could not parse L2VPN ID '%s': %s", data.L2VPN.ValueString(), err.Error()),
		)
		return
	}
	l2vpn, httpResp, err := r.client.VpnAPI.VpnL2vpnsRetrieve(ctx, l2vpnID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error looking up L2VPN",
			utils.FormatAPIError("lookup L2VPN", err, httpResp),
		)
		return
	}

	// Build the API request
	l2vpnRef := netbox.NewBriefL2VPNRequest(l2vpn.GetName(), l2vpn.GetSlug())
	terminationRequest := netbox.NewL2VPNTerminationRequest(
		*l2vpnRef,
		data.AssignedObjectType.ValueString(),
		data.AssignedObjectID.ValueInt64(),
	)

	// Handle tags and custom fields - merge-aware
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		utils.ApplyTagsFromSlugs(ctx, r.client, terminationRequest, data.Tags, &resp.Diagnostics)
	} else if !state.Tags.IsNull() && !state.Tags.IsUnknown() {
		utils.ApplyTagsFromSlugs(ctx, r.client, terminationRequest, state.Tags, &resp.Diagnostics)
	}

	utils.ApplyCustomFieldsWithMerge(ctx, terminationRequest, data.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating L2VPN termination", map[string]interface{}{
		"id":       idInt,
		"l2vpn_id": l2vpnID,
	})

	// Call the API
	termination, httpResp, err := r.client.VpnAPI.VpnL2vpnTerminationsUpdate(ctx, idInt).L2VPNTerminationRequest(*terminationRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating L2VPN termination",
			utils.FormatAPIError("update L2VPN termination", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToState(ctx, termination, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNTerminationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data L2VPNTerminationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var idInt int32
	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse L2VPN termination ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting L2VPN termination", map[string]interface{}{
		"id": idInt,
	})

	// Call the API
	httpResp, err := r.client.VpnAPI.VpnL2vpnTerminationsDestroy(ctx, idInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting L2VPN termination",
			utils.FormatAPIError("delete L2VPN termination", err, httpResp),
		)
		return
	}
}

func (r *L2VPNTerminationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		idInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID format", fmt.Sprintf("Could not parse L2VPN termination ID '%s': %s", parsed.ID, err.Error()))
			return
		}
		termination, httpResp, err := r.client.VpnAPI.VpnL2vpnTerminationsRetrieve(ctx, idInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing L2VPN termination", utils.FormatAPIError("read L2VPN termination", err, httpResp))
			return
		}

		var data L2VPNTerminationResourceModel
		l2vpn := termination.GetL2vpn()
		if l2vpn.GetId() != 0 {
			data.L2VPN = types.StringValue(fmt.Sprintf("%d", l2vpn.GetId()))
		}
		data.AssignedObjectType = types.StringValue(termination.GetAssignedObjectType())
		data.AssignedObjectID = types.Int64Value(termination.GetAssignedObjectId())
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, termination.HasTags(), termination.GetTags(), data.Tags)
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

		r.mapResponseToState(ctx, termination, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, termination.GetCustomFields(), &resp.Diagnostics)
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

// mapResponseToState maps an L2VPNTermination API response to the Terraform state model.
func (r *L2VPNTerminationResource) mapResponseToState(ctx context.Context, termination *netbox.L2VPNTermination, data *L2VPNTerminationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", termination.GetId()))

	// L2VPN
	l2vpn := termination.GetL2vpn()
	data.L2VPN = types.StringValue(fmt.Sprintf("%d", l2vpn.GetId()))

	// Assigned object
	data.AssignedObjectType = types.StringValue(termination.GetAssignedObjectType())
	data.AssignedObjectID = types.Int64Value(termination.GetAssignedObjectId())

	// Handle tags using filter-to-owned approach
	planTags := data.Tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, termination.HasTags(), termination.GetTags(), planTags)

	// Handle custom fields using consolidated helper
	if termination.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, termination.GetCustomFields(), diags)
	}
}
