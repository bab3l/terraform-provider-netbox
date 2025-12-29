// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &L2VPNTerminationResource{}

var _ resource.ResourceWithImportState = &L2VPNTerminationResource{}

func NewL2VPNTerminationResource() resource.Resource {
	return &L2VPNTerminationResource{}
}

// L2VPNTerminationResource defines the resource implementation.

type L2VPNTerminationResource struct {
	client *netbox.APIClient
}

// L2VPNTerminationResourceModel describes the resource data model.

type L2VPNTerminationResourceModel struct {
	ID types.String `tfsdk:"id"`

	L2VPN types.String `tfsdk:"l2vpn"`

	AssignedObjectType types.String `tfsdk:"assigned_object_type"`

	AssignedObjectID types.Int64 `tfsdk:"assigned_object_id"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
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

				Required: true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer ID",
					),
				},
			},

			"assigned_object_type": schema.StringAttribute{
				MarkdownDescription: "Content type of the assigned object. Valid values: `dcim.interface`, `ipam.vlan`, `virtualization.vminterface`.",

				Required: true,

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

				Required: true,
			},
		},
	}

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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

	utils.ApplyMetadataFields(ctx, terminationRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating L2VPN termination", map[string]interface{}{
		"l2vpn_id": l2vpnID,

		"assigned_object_type": data.AssignedObjectType.ValueString(),

		"assigned_object_id": data.AssignedObjectID.ValueInt64(),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNTerminationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data L2VPNTerminationResourceModel

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

	// Handle tags and custom fields

	utils.ApplyMetadataFields(ctx, terminationRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating L2VPN termination", map[string]interface{}{
		"id": idInt,

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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapResponseToState maps an L2VPNTermination API response to the Terraform state model.

func (r *L2VPNTerminationResource) mapResponseToState(ctx context.Context, termination *netbox.L2VPNTermination, data *L2VPNTerminationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", termination.GetId()))

	// DisplayName
	if termination.Display != "" {
	} else {
	}

	// L2VPN

	l2vpn := termination.GetL2vpn()

	data.L2VPN = types.StringValue(fmt.Sprintf("%d", l2vpn.GetId()))

	// Assigned object

	data.AssignedObjectType = types.StringValue(termination.GetAssignedObjectType())

	data.AssignedObjectID = types.Int64Value(termination.GetAssignedObjectId())

	// Tags

	if termination.HasTags() && len(termination.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(termination.GetTags())

		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(d...)

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields

	if termination.HasCustomFields() && len(termination.GetCustomFields()) > 0 {
		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {
			d := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(d...)
		}

		customFields := utils.MapToCustomFieldModels(termination.GetCustomFields(), existingModels)

		customFieldsValue, d := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(d...)

		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
