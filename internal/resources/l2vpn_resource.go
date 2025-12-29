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

var _ resource.Resource = &L2VPNResource{}

var _ resource.ResourceWithImportState = &L2VPNResource{}

func NewL2VPNResource() resource.Resource {
	return &L2VPNResource{}
}

// L2VPNResource defines the resource implementation.

type L2VPNResource struct {
	client *netbox.APIClient
}

// L2VPNResourceModel describes the resource data model.

type L2VPNResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Type types.String `tfsdk:"type"`

	Identifier types.Int64 `tfsdk:"identifier"`

	ImportTargets types.Set `tfsdk:"import_targets"`

	ExportTargets types.Set `tfsdk:"export_targets"`

	Tenant types.String `tfsdk:"tenant"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *L2VPNResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn"
}

func (r *L2VPNResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Layer 2 VPN in Netbox. L2VPNs represent layer 2 virtual private network services such as VPLS, VXLAN, EVPN, etc.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("L2VPN"),

			"name": nbschema.NameAttribute("L2VPN", 100),

			"slug": nbschema.SlugAttribute("L2VPN"),

			"type": schema.StringAttribute{
				MarkdownDescription: "L2VPN type. Valid values: `vpws`, `vpls`, `vxlan`, `vxlan-evpn`, `mpls-evpn`, `pbb-evpn`, `evpn-vpws`, `epl`, `evpl`, `ep-lan`, `evp-lan`, `ep-tree`, `evp-tree`.",

				Required: true,

				Validators: []validator.String{
					stringvalidator.OneOf(

						"vpws",

						"vpls",

						"vxlan",

						"vxlan-evpn",

						"mpls-evpn",

						"pbb-evpn",

						"evpn-vpws",

						"epl",

						"evpl",

						"ep-lan",

						"evp-lan",

						"ep-tree",

						"evp-tree",
					),
				},
			},

			"identifier": schema.Int64Attribute{
				MarkdownDescription: "Numeric identifier unique to the parent L2VPN.",

				Optional: true,
			},

			"import_targets": schema.SetAttribute{
				MarkdownDescription: "Set of route target IDs to import.",

				Optional: true,

				ElementType: types.StringType,
			},

			"export_targets": schema.SetAttribute{
				MarkdownDescription: "Set of route target IDs to export.",

				Optional: true,

				ElementType: types.StringType,
			},

			"tenant": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant this L2VPN belongs to.",

				Optional: true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer ID",
					),
				},
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("L2VPN"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *L2VPNResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *L2VPNResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data L2VPNResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request

	l2vpnRequest := netbox.NewWritableL2VPNRequest(

		data.Name.ValueString(),

		data.Slug.ValueString(),

		netbox.BriefL2VPNTypeValue(data.Type.ValueString()),
	)

	// Set optional fields

	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		l2vpnRequest.SetIdentifier(data.Identifier.ValueInt64())
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, l2vpnRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle tenant reference

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		l2vpnRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	}

	// Handle import targets

	if !data.ImportTargets.IsNull() && !data.ImportTargets.IsUnknown() {
		var targetIDs []string

		diags := data.ImportTargets.ElementsAs(ctx, &targetIDs, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		var importTargets []int32

		for _, idStr := range targetIDs {
			var id int32

			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				importTargets = append(importTargets, id)
			}
		}

		l2vpnRequest.ImportTargets = importTargets
	}

	// Handle export targets

	if !data.ExportTargets.IsNull() && !data.ExportTargets.IsUnknown() {
		var targetIDs []string

		diags := data.ExportTargets.ElementsAs(ctx, &targetIDs, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		var exportTargets []int32

		for _, idStr := range targetIDs {
			var id int32

			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				exportTargets = append(exportTargets, id)
			}
		}

		l2vpnRequest.ExportTargets = exportTargets
	}

	tflog.Debug(ctx, "Creating L2VPN", map[string]interface{}{
		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),

		"type": data.Type.ValueString(),
	})

	// Call the API

	l2vpn, httpResp, err := r.client.VpnAPI.VpnL2vpnsCreate(ctx).WritableL2VPNRequest(*l2vpnRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating L2VPN",

			utils.FormatAPIError("create L2VPN", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, l2vpn, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created L2VPN resource", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data L2VPNResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse L2VPN ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	// Read from API

	l2vpn, httpResp, err := r.client.VpnAPI.VpnL2vpnsRetrieve(ctx, idInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading L2VPN",

			utils.FormatAPIError("read L2VPN", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, l2vpn, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data L2VPNResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse L2VPN ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	// Build the API request

	l2vpnRequest := netbox.NewWritableL2VPNRequest(

		data.Name.ValueString(),

		data.Slug.ValueString(),

		netbox.BriefL2VPNTypeValue(data.Type.ValueString()),
	)

	// Set optional fields

	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		l2vpnRequest.SetIdentifier(data.Identifier.ValueInt64())
	} else {
		l2vpnRequest.SetIdentifierNil()
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, l2vpnRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle tenant reference

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		l2vpnRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	} else {
		l2vpnRequest.Tenant = *netbox.NewNullableBriefTenantRequest(nil)
	}

	// Handle import targets

	if !data.ImportTargets.IsNull() && !data.ImportTargets.IsUnknown() {
		var targetIDs []string

		diags := data.ImportTargets.ElementsAs(ctx, &targetIDs, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		var importTargets []int32

		for _, idStr := range targetIDs {
			var id int32

			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				importTargets = append(importTargets, id)
			}
		}

		l2vpnRequest.ImportTargets = importTargets
	} else {
		l2vpnRequest.ImportTargets = []int32{}
	}

	// Handle export targets

	if !data.ExportTargets.IsNull() && !data.ExportTargets.IsUnknown() {
		var targetIDs []string

		diags := data.ExportTargets.ElementsAs(ctx, &targetIDs, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		var exportTargets []int32

		for _, idStr := range targetIDs {
			var id int32

			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				exportTargets = append(exportTargets, id)
			}
		}

		l2vpnRequest.ExportTargets = exportTargets
	} else {
		l2vpnRequest.ExportTargets = []int32{}
	}

	tflog.Debug(ctx, "Updating L2VPN", map[string]interface{}{
		"id": idInt,

		"name": data.Name.ValueString(),
	})

	// Call the API

	l2vpn, httpResp, err := r.client.VpnAPI.VpnL2vpnsUpdate(ctx, idInt).WritableL2VPNRequest(*l2vpnRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating L2VPN",

			utils.FormatAPIError("update L2VPN", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, l2vpn, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2VPNResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data L2VPNResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse L2VPN ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting L2VPN", map[string]interface{}{
		"id": idInt,

		"name": data.Name.ValueString(),
	})

	// Call the API

	httpResp, err := r.client.VpnAPI.VpnL2vpnsDestroy(ctx, idInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting L2VPN",

			utils.FormatAPIError("delete L2VPN", err, httpResp),
		)

		return
	}
}

func (r *L2VPNResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapResponseToState maps an L2VPN API response to the Terraform state model.

func (r *L2VPNResource) mapResponseToState(ctx context.Context, l2vpn *netbox.L2VPN, data *L2VPNResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", l2vpn.GetId()))

	data.Name = types.StringValue(l2vpn.GetName())

	data.Slug = types.StringValue(l2vpn.GetSlug())

	// Type

	if l2vpn.HasType() {
		typeObj := l2vpn.GetType()

		data.Type = types.StringValue(string(typeObj.GetValue()))
	}

	// Identifier

	if l2vpn.HasIdentifier() && l2vpn.GetIdentifier() != 0 {
		data.Identifier = types.Int64Value(l2vpn.GetIdentifier())
	} else {
		data.Identifier = types.Int64Null()
	}

	// Description

	if l2vpn.HasDescription() && l2vpn.GetDescription() != "" {
		data.Description = types.StringValue(l2vpn.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments

	if l2vpn.HasComments() && l2vpn.GetComments() != "" {
		data.Comments = types.StringValue(l2vpn.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Tenant - preserve user's input format

	if l2vpn.HasTenant() && l2vpn.GetTenant().Id != 0 {
		tenant := l2vpn.GetTenant()

		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Import targets

	if l2vpn.HasImportTargets() && len(l2vpn.GetImportTargets()) > 0 {
		var targetIDs []string

		for _, target := range l2vpn.GetImportTargets() {
			targetIDs = append(targetIDs, fmt.Sprintf("%d", target.GetId()))
		}

		targetSet, d := types.SetValueFrom(ctx, types.StringType, targetIDs)

		diags.Append(d...)

		data.ImportTargets = targetSet
	} else {
		data.ImportTargets = types.SetNull(types.StringType)
	}

	// Export targets

	if l2vpn.HasExportTargets() && len(l2vpn.GetExportTargets()) > 0 {
		var targetIDs []string

		for _, target := range l2vpn.GetExportTargets() {
			targetIDs = append(targetIDs, fmt.Sprintf("%d", target.GetId()))
		}

		targetSet, d := types.SetValueFrom(ctx, types.StringType, targetIDs)

		diags.Append(d...)

		data.ExportTargets = targetSet
	} else {
		data.ExportTargets = types.SetNull(types.StringType)
	}

	// Tags

	if l2vpn.HasTags() && len(l2vpn.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(l2vpn.GetTags())

		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(d...)

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields

	if l2vpn.HasCustomFields() && len(l2vpn.GetCustomFields()) > 0 {
		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {
			d := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(d...)
		}

		customFields := utils.MapToCustomFieldModels(l2vpn.GetCustomFields(), existingModels)

		customFieldsValue, d := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(d...)

		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
