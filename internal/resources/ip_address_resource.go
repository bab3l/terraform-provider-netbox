// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IPAddressResource{}
	_ resource.ResourceWithConfigure   = &IPAddressResource{}
	_ resource.ResourceWithImportState = &IPAddressResource{}
)

// NewIPAddressResource returns a new IP Address resource.
func NewIPAddressResource() resource.Resource {
	return &IPAddressResource{}
}

// IPAddressResource defines the resource implementation.
type IPAddressResource struct {
	client *netbox.APIClient
}

// IPAddressResourceModel describes the resource data model.
type IPAddressResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Address            types.String `tfsdk:"address"`
	VRF                types.String `tfsdk:"vrf"`
	Tenant             types.String `tfsdk:"tenant"`
	Status             types.String `tfsdk:"status"`
	Role               types.String `tfsdk:"role"`
	AssignedObjectType types.String `tfsdk:"assigned_object_type"`
	AssignedObjectID   types.Int64  `tfsdk:"assigned_object_id"`
	DNSName            types.String `tfsdk:"dns_name"`
	Description        types.String `tfsdk:"description"`
	Comments           types.String `tfsdk:"comments"`
	Tags               types.Set    `tfsdk:"tags"`
	CustomFields       types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *IPAddressResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_address"
}

// Schema defines the schema for the resource.
func (r *IPAddressResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IP address in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the IP address.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The IP address with prefix length (e.g., 192.168.1.1/24).",
				Required:            true,
			},
			"vrf":    nbschema.ReferenceAttributeWithDiffSuppress("VRF", "ID or name of the VRF this IP address is assigned to."),
			"tenant": nbschema.ReferenceAttributeWithDiffSuppress("tenant", "ID or slug of the tenant this IP address is assigned to."),
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the IP address. Valid values are: `active`, `reserved`, `deprecated`, `dhcp`, `slaac`. Defaults to `active`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the IP address. Valid values are: `loopback`, `secondary`, `anycast`, `vip`, `vrrp`, `hsrp`, `glbp`, `carp`.",
				Optional:            true,
			},
			"assigned_object_type": schema.StringAttribute{
				MarkdownDescription: "The content type of the assigned object (e.g., `dcim.interface`, `virtualization.vminterface`).",
				Optional:            true,
			},
			"assigned_object_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the assigned object (interface or VM interface).",
				Optional:            true,
			},
			"dns_name": schema.StringAttribute{
				MarkdownDescription: "Hostname or FQDN (not case-sensitive).",
				Optional:            true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("IP address"))

	// Add tags and custom_fields attributes
	resp.Schema.Attributes["tags"] = nbschema.TagsAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

// Configure adds the provider configured client to the resource.
func (r *IPAddressResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IPAddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IPAddressResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the IP address request
	ipRequest := netbox.NewWritableIPAddressRequest(data.Address.ValueString())

	// Set optional fields (pass nil for state since this is Create)
	r.setOptionalFields(ctx, ipRequest, &data, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating IP address", map[string]interface{}{
		"address": data.Address.ValueString(),
	})

	// Create the IP address
	ipAddress, httpResp, err := r.client.IpamAPI.IpamIpAddressesCreate(ctx).WritableIPAddressRequest(*ipRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IP address",
			utils.FormatAPIError("create IP address", err, httpResp),
		)
		return
	}

	// Map response to model (use filter-to-owned for custom fields)
	r.mapIPAddressToState(ctx, ipAddress, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created IP address", map[string]interface{}{
		"id":      data.ID.ValueString(),
		"address": data.Address.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *IPAddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IPAddressResourceModel

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
	tflog.Debug(ctx, "Reading IP address", map[string]interface{}{
		"id": id,
	})

	// Get the IP address from Netbox
	ipAddress, httpResp, err := r.client.IpamAPI.IpamIpAddressesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IP address",
			utils.FormatAPIError(fmt.Sprintf("read IP address ID %d", id), err, httpResp),
		)
		return
	}

	// Save original custom_fields state before mapping
	originalCustomFields := data.CustomFields

	// Map response to model (use filter-to-owned for custom fields)
	r.mapIPAddressToState(ctx, ipAddress, &data, &resp.Diagnostics)

	// Preserve original custom_fields state if it was null or empty
	// This prevents unmanaged/cleared fields from reappearing in state
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}

	tflog.Debug(ctx, "Read IP address", map[string]interface{}{
		"id":      data.ID.ValueString(),
		"address": data.Address.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *IPAddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan IPAddressResourceModel

	// Read Terraform state and plan data into the models
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Create the IP address request
	ipRequest := netbox.NewWritableIPAddressRequest(plan.Address.ValueString())

	// Set optional fields (including merge-aware custom fields)
	r.setOptionalFields(ctx, ipRequest, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating IP address", map[string]interface{}{
		"id":      id,
		"address": plan.Address.ValueString(),
	})

	// Update the IP address
	ipAddress, httpResp, err := r.client.IpamAPI.IpamIpAddressesUpdate(ctx, id).WritableIPAddressRequest(*ipRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IP address",
			utils.FormatAPIError(fmt.Sprintf("update IP address ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model (use filter-to-owned for custom fields)
	r.mapIPAddressToState(ctx, ipAddress, &plan, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated IP address", map[string]interface{}{
		"id":      plan.ID.ValueString(),
		"address": plan.Address.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IPAddressResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IPAddressResourceModel

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
	tflog.Debug(ctx, "Deleting IP address", map[string]interface{}{
		"id": id,
	})

	// Delete the IP address
	httpResp, err := r.client.IpamAPI.IpamIpAddressesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		// Ignore 404 errors (resource already deleted)
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "IP address already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting IP address",
			utils.FormatAPIError(fmt.Sprintf("delete IP address ID %d", id), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted IP address", map[string]interface{}{
		"id": id,
	})
}

func (r *IPAddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the IP address request from the resource model.
func (r *IPAddressResource) setOptionalFields(ctx context.Context, ipRequest *netbox.WritableIPAddressRequest, plan *IPAddressResourceModel, state *IPAddressResourceModel, diags *diag.Diagnostics) {
	// VRF
	if utils.IsSet(plan.VRF) {
		vrf, vrfDiags := netboxlookup.LookupVRF(ctx, r.client, plan.VRF.ValueString())
		diags.Append(vrfDiags...)
		if diags.HasError() {
			return
		}
		ipRequest.Vrf = *netbox.NewNullableBriefVRFRequest(vrf)
	} else if plan.VRF.IsNull() {
		ipRequest.SetVrfNil()
	}

	// Tenant
	if utils.IsSet(plan.Tenant) {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, plan.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return
		}
		ipRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	} else if plan.Tenant.IsNull() {
		ipRequest.SetTenantNil()
	}

	// Status
	if utils.IsSet(plan.Status) {
		status := netbox.PatchedWritableIPAddressRequestStatus(plan.Status.ValueString())
		ipRequest.Status = &status
	}

	// Role
	if utils.IsSet(plan.Role) {
		role := netbox.PatchedWritableIPAddressRequestRole(plan.Role.ValueString())
		ipRequest.Role = &role
	}

	// Assigned Object Type
	if utils.IsSet(plan.AssignedObjectType) {
		objType := plan.AssignedObjectType.ValueString()
		ipRequest.AssignedObjectType = *netbox.NewNullableString(&objType)
	}

	// Assigned Object ID
	if utils.IsSet(plan.AssignedObjectID) {
		objID := plan.AssignedObjectID.ValueInt64()
		ipRequest.AssignedObjectId = *netbox.NewNullableInt64(&objID)
	}

	// DNS Name
	if utils.IsSet(plan.DNSName) {
		dnsName := plan.DNSName.ValueString()
		ipRequest.DnsName = &dnsName
	}

	// Description
	ipRequest.Description = utils.StringPtr(plan.Description)

	// Comments
	ipRequest.Comments = utils.StringPtr(plan.Comments)

	// Handle tags
	utils.ApplyTags(ctx, ipRequest, plan.Tags, diags)
	if diags.HasError() {
		return
	}
	// Handle custom fields - use merge-aware helper if state exists (Update), otherwise regular (Create)
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, ipRequest, plan.CustomFields, state.CustomFields, diags)
	} else {
		utils.ApplyCustomFields(ctx, ipRequest, plan.CustomFields, diags)
	}
}

// mapIPAddressToState maps a Netbox IPAddress to the Terraform state model.
func (r *IPAddressResource) mapIPAddressToState(ctx context.Context, ipAddress *netbox.IPAddress, data *IPAddressResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", ipAddress.Id))
	data.Address = types.StringValue(ipAddress.Address)

	// VRF
	if ipAddress.Vrf.IsSet() && ipAddress.Vrf.Get() != nil {
		vrfObj := ipAddress.Vrf.Get()
		data.VRF = utils.UpdateReferenceAttribute(data.VRF, vrfObj.Name, "", vrfObj.Id)
	} else {
		data.VRF = types.StringNull()
	}

	// Tenant
	if ipAddress.Tenant.IsSet() && ipAddress.Tenant.Get() != nil {
		tenantObj := ipAddress.Tenant.Get()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenantObj.Name, tenantObj.Slug, tenantObj.Id)
	} else {
		data.Tenant = types.StringNull()
	}

	// Status
	if ipAddress.Status != nil {
		data.Status = types.StringValue(string(ipAddress.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Role
	if ipAddress.Role != nil {
		data.Role = types.StringValue(string(ipAddress.Role.GetValue()))
	} else {
		data.Role = types.StringNull()
	}

	// Assigned Object Type
	if ipAddress.AssignedObjectType.IsSet() && ipAddress.AssignedObjectType.Get() != nil {
		data.AssignedObjectType = types.StringValue(*ipAddress.AssignedObjectType.Get())
	} else {
		data.AssignedObjectType = types.StringNull()
	}

	// Assigned Object ID
	if ipAddress.AssignedObjectId.IsSet() && ipAddress.AssignedObjectId.Get() != nil {
		data.AssignedObjectID = types.Int64Value(*ipAddress.AssignedObjectId.Get())
	} else {
		data.AssignedObjectID = types.Int64Null()
	}

	// DNS Name
	if ipAddress.DnsName != nil && *ipAddress.DnsName != "" {
		data.DNSName = types.StringValue(*ipAddress.DnsName)
	} else {
		data.DNSName = types.StringNull()
	}

	// Description
	if ipAddress.Description != nil && *ipAddress.Description != "" {
		data.Description = types.StringValue(*ipAddress.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if ipAddress.Comments != nil && *ipAddress.Comments != "" {
		data.Comments = types.StringValue(*ipAddress.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags
	data.Tags = utils.PopulateTagsFromAPI(ctx, len(ipAddress.Tags) > 0, ipAddress.Tags, data.Tags, diags)

	// Custom Fields - use filter-to-owned pattern
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, ipAddress.CustomFields, diags)
}
