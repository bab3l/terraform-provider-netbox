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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IPRangeResource{}
	_ resource.ResourceWithConfigure   = &IPRangeResource{}
	_ resource.ResourceWithImportState = &IPRangeResource{}
)

// NewIPRangeResource returns a new IP Range resource.
func NewIPRangeResource() resource.Resource {
	return &IPRangeResource{}
}

// IPRangeResource defines the resource implementation.
type IPRangeResource struct {
	client *netbox.APIClient
}

// IPRangeResourceModel describes the resource data model.
type IPRangeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
	Size         types.Int64  `tfsdk:"size"`
	VRF          types.String `tfsdk:"vrf"`
	Tenant       types.String `tfsdk:"tenant"`
	Status       types.String `tfsdk:"status"`
	Role         types.String `tfsdk:"role"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	MarkUtilized types.Bool   `tfsdk:"mark_utilized"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IPRangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_range"
}

// Schema defines the schema for the resource.
func (r *IPRangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IP address range in Netbox. IP ranges are used to define a contiguous range of IP addresses within a prefix.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IP range.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_address": schema.StringAttribute{
				MarkdownDescription: "The starting IP address of the range (e.g., 192.168.1.10/24).",
				Required:            true,
			},
			"end_address": schema.StringAttribute{
				MarkdownDescription: "The ending IP address of the range (e.g., 192.168.1.20/24).",
				Required:            true,
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "The number of IP addresses in the range (computed).",
				Computed:            true,
			},
			"vrf": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the VRF this IP range is assigned to.",
				Optional:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the tenant this IP range is assigned to.",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the IP range. Valid values are: `active`, `reserved`, `deprecated`. Defaults to `active`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the IPAM role for this IP range.",
				Optional:            true,
			},
			"description": nbschema.DescriptionAttribute("IP range"),
			"comments":    nbschema.CommentsAttribute("IP range"),
			"mark_utilized": schema.BoolAttribute{
				MarkdownDescription: "Treat this range as fully utilized regardless of actual usage. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *IPRangeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *IPRangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IPRangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IP range request
	ipRangeRequest := netbox.NewWritableIPRangeRequest(data.StartAddress.ValueString(), data.EndAddress.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, ipRangeRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating IP range", map[string]interface{}{
		"start_address": data.StartAddress.ValueString(),
		"end_address":   data.EndAddress.ValueString(),
	})

	// Create the IP range
	ipRange, httpResp, err := r.client.IpamAPI.IpamIpRangesCreate(ctx).WritableIPRangeRequest(*ipRangeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IP range",
			utils.FormatAPIError("create IP range", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIPRangeToState(ctx, ipRange, &data)

	tflog.Debug(ctx, "Created IP range", map[string]interface{}{
		"id":            data.ID.ValueString(),
		"start_address": data.StartAddress.ValueString(),
		"end_address":   data.EndAddress.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *IPRangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IPRangeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Reading IP range", map[string]interface{}{
		"id": id,
	})

	// Get the IP range from Netbox
	ipRange, httpResp, err := r.client.IpamAPI.IpamIpRangesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IP range",
			utils.FormatAPIError(fmt.Sprintf("read IP range ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIPRangeToState(ctx, ipRange, &data)

	tflog.Debug(ctx, "Read IP range", map[string]interface{}{
		"id":            data.ID.ValueString(),
		"start_address": data.StartAddress.ValueString(),
		"end_address":   data.EndAddress.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *IPRangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IPRangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Create the IP range request
	ipRangeRequest := netbox.NewWritableIPRangeRequest(data.StartAddress.ValueString(), data.EndAddress.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, ipRangeRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating IP range", map[string]interface{}{
		"id":            id,
		"start_address": data.StartAddress.ValueString(),
		"end_address":   data.EndAddress.ValueString(),
	})

	// Update the IP range
	ipRange, httpResp, err := r.client.IpamAPI.IpamIpRangesUpdate(ctx, id).WritableIPRangeRequest(*ipRangeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IP range",
			utils.FormatAPIError(fmt.Sprintf("update IP range ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapIPRangeToState(ctx, ipRange, &data)

	tflog.Debug(ctx, "Updated IP range", map[string]interface{}{
		"id":            data.ID.ValueString(),
		"start_address": data.StartAddress.ValueString(),
		"end_address":   data.EndAddress.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IPRangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IPRangeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Deleting IP range", map[string]interface{}{
		"id": id,
	})

	// Delete the IP range
	httpResp, err := r.client.IpamAPI.IpamIpRangesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting IP range",
			utils.FormatAPIError(fmt.Sprintf("delete IP range ID %d", id), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted IP range", map[string]interface{}{
		"id": id,
	})
}

func (r *IPRangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the IP range request from the resource model.
func (r *IPRangeResource) setOptionalFields(ctx context.Context, ipRangeRequest *netbox.WritableIPRangeRequest, data *IPRangeResourceModel, diags *diag.Diagnostics) {
	// VRF
	if utils.IsSet(data.VRF) {
		vrf, vrfDiags := netboxlookup.LookupVRF(ctx, r.client, data.VRF.ValueString())
		diags.Append(vrfDiags...)
		if diags.HasError() {
			return
		}
		ipRangeRequest.Vrf = *netbox.NewNullableBriefVRFRequest(vrf)
	}

	// Tenant
	if utils.IsSet(data.Tenant) {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return
		}
		ipRangeRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	}

	// Status
	if utils.IsSet(data.Status) {
		status := netbox.PatchedWritableIPRangeRequestStatus(data.Status.ValueString())
		ipRangeRequest.Status = &status
	}

	// Role (IPAM Role)
	if utils.IsSet(data.Role) {
		role, roleDiags := netboxlookup.LookupRole(ctx, r.client, data.Role.ValueString())
		diags.Append(roleDiags...)
		if diags.HasError() {
			return
		}
		ipRangeRequest.Role = *netbox.NewNullableBriefRoleRequest(role)
	}

	// Description
	ipRangeRequest.Description = utils.StringPtr(data.Description)

	// Comments
	ipRangeRequest.Comments = utils.StringPtr(data.Comments)

	// Mark Utilized
	if utils.IsSet(data.MarkUtilized) {
		markUtilized := data.MarkUtilized.ValueBool()
		ipRangeRequest.MarkUtilized = &markUtilized
	}

	// Handle tags
	if utils.IsSet(data.Tags) {
		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		ipRangeRequest.Tags = tags
	}

	// Handle custom fields
	if utils.IsSet(data.CustomFields) {
		var customFields []utils.CustomFieldModel
		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if diags.HasError() {
			return
		}
		ipRangeRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}
}

// mapIPRangeToState maps a Netbox IPRange to the Terraform state model.
func (r *IPRangeResource) mapIPRangeToState(ctx context.Context, ipRange *netbox.IPRange, data *IPRangeResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", ipRange.Id))
	data.StartAddress = types.StringValue(ipRange.StartAddress)
	data.EndAddress = types.StringValue(ipRange.EndAddress)
	data.Size = types.Int64Value(int64(ipRange.Size))

	// VRF - preserve user input if it matches
	if ipRange.Vrf.IsSet() && ipRange.Vrf.Get() != nil {
		vrfObj := ipRange.Vrf.Get()
		userVrf := data.VRF.ValueString()
		if userVrf == vrfObj.Name || userVrf == fmt.Sprintf("%d", vrfObj.Id) {
			// Keep user's original value
		} else {
			data.VRF = types.StringValue(vrfObj.Name)
		}
	} else {
		data.VRF = types.StringNull()
	}

	// Tenant - preserve user input if it matches
	if ipRange.Tenant.IsSet() && ipRange.Tenant.Get() != nil {
		tenantObj := ipRange.Tenant.Get()
		userTenant := data.Tenant.ValueString()
		if userTenant == tenantObj.Name || userTenant == tenantObj.Slug || userTenant == fmt.Sprintf("%d", tenantObj.Id) {
			// Keep user's original value
		} else {
			data.Tenant = types.StringValue(tenantObj.Name)
		}
	} else {
		data.Tenant = types.StringNull()
	}

	// Status
	if ipRange.Status != nil {
		data.Status = types.StringValue(string(ipRange.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Role - preserve user input if it matches
	if ipRange.Role.IsSet() && ipRange.Role.Get() != nil {
		roleObj := ipRange.Role.Get()
		userRole := data.Role.ValueString()
		if userRole == roleObj.Name || userRole == roleObj.Slug || userRole == fmt.Sprintf("%d", roleObj.Id) {
			// Keep user's original value
		} else {
			data.Role = types.StringValue(roleObj.Name)
		}
	} else {
		data.Role = types.StringNull()
	}

	// Description
	if ipRange.Description != nil && *ipRange.Description != "" {
		data.Description = types.StringValue(*ipRange.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if ipRange.Comments != nil && *ipRange.Comments != "" {
		data.Comments = types.StringValue(*ipRange.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Mark Utilized
	if ipRange.MarkUtilized != nil {
		data.MarkUtilized = types.BoolValue(*ipRange.MarkUtilized)
	} else {
		data.MarkUtilized = types.BoolValue(false)
	}

	// Tags
	if len(ipRange.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(ipRange.Tags)
		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields
	switch {
	case len(ipRange.CustomFields) > 0 && !data.CustomFields.IsNull():
		var stateCustomFields []utils.CustomFieldModel
		data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		customFields := utils.MapToCustomFieldModels(ipRange.CustomFields, stateCustomFields)
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	case len(ipRange.CustomFields) > 0:
		customFields := utils.MapToCustomFieldModels(ipRange.CustomFields, []utils.CustomFieldModel{})
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	default:
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
