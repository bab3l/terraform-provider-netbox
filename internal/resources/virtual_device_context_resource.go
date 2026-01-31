// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &VirtualDeviceContextResource{}
	_ resource.ResourceWithConfigure   = &VirtualDeviceContextResource{}
	_ resource.ResourceWithImportState = &VirtualDeviceContextResource{}
	_ resource.ResourceWithIdentity    = &VirtualDeviceContextResource{}
)

// NewVirtualDeviceContextResource returns a new resource implementing the virtual device context resource.
func NewVirtualDeviceContextResource() resource.Resource {
	return &VirtualDeviceContextResource{}
}

// VirtualDeviceContextResource defines the resource implementation.
type VirtualDeviceContextResource struct {
	client *netbox.APIClient
}

// VirtualDeviceContextResourceModel describes the resource data model.
type VirtualDeviceContextResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Device       types.String `tfsdk:"device"`
	Identifier   types.Int64  `tfsdk:"identifier"`
	Tenant       types.String `tfsdk:"tenant"`
	PrimaryIP4   types.String `tfsdk:"primary_ip4"`
	PrimaryIP6   types.String `tfsdk:"primary_ip6"`
	Status       types.String `tfsdk:"status"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *VirtualDeviceContextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_device_context"
}

// Schema defines the schema for the resource.
func (r *VirtualDeviceContextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual device context (VDC) in NetBox. Virtual device contexts allow a single physical device to be logically partitioned into multiple virtual devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the virtual device context.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the virtual device context.",
				Required:            true,
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device",
				"The device this VDC belongs to (ID or name).",
			),
			"identifier": schema.Int64Attribute{
				MarkdownDescription: "Numeric identifier unique to the parent device.",
				Optional:            true,
			},
			"tenant": nbschema.ReferenceAttributeWithDiffSuppress(
				"tenant",
				"The tenant associated with this VDC (ID or slug).",
			),
			"primary_ip4": schema.StringAttribute{
				MarkdownDescription: "Primary IPv4 address assigned to this VDC (ID).",
				Optional:            true,
			},
			"primary_ip6": schema.StringAttribute{
				MarkdownDescription: "Primary IPv6 address assigned to this VDC (ID).",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the VDC. Valid values: `active`, `planned`, `offline`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("active", "planned", "offline"),
				},
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("virtual device context"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *VirtualDeviceContextResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *VirtualDeviceContextResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.
func (r *VirtualDeviceContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualDeviceContextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	status := netbox.PatchedWritableVirtualDeviceContextRequestStatus(data.Status.ValueString())
	apiReq := netbox.NewWritableVirtualDeviceContextRequest(data.Name.ValueString(), *device, status)

	// Set optional fields
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		identifier, err := utils.SafeInt32FromValue(data.Identifier)
		if err != nil {
			resp.Diagnostics.AddError("Invalid identifier", fmt.Sprintf("Identifier overflow: %s", err))
			return
		}
		apiReq.SetIdentifier(identifier)
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := lookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		resp.Diagnostics.Append(tenantDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTenant(*tenant)
	}

	if !data.PrimaryIP4.IsNull() && !data.PrimaryIP4.IsUnknown() {
		ipAddr, ipDiags := lookup.LookupIPAddress(ctx, r.client, data.PrimaryIP4.ValueString())
		resp.Diagnostics.Append(ipDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetPrimaryIp4(*ipAddr)
	}

	if !data.PrimaryIP6.IsNull() && !data.PrimaryIP6.IsUnknown() {
		ipAddr, ipDiags := lookup.LookupIPAddress(ctx, r.client, data.PrimaryIP6.ValueString())
		resp.Diagnostics.Append(ipDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetPrimaryIp6(*ipAddr)
	}

	// Apply common fields (description, comments, tags, custom_fields)
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyComments(apiReq, data.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating virtual device context", map[string]interface{}{
		"name":   data.Name.ValueString(),
		"device": data.Device.ValueString(),
		"status": data.Status.ValueString(),
	})

	// Create the resource
	result, httpResp, err := r.client.DcimAPI.DcimVirtualDeviceContextsCreate(ctx).WritableVirtualDeviceContextRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating virtual device context",
			utils.FormatAPIError("create virtual device context", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the resource.
func (r *VirtualDeviceContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualDeviceContextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields to detect null/empty cases
	originalCustomFields := data.CustomFields

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing virtual device context ID",
			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Read from API
	result, httpResp, err := r.client.DcimAPI.DcimVirtualDeviceContextsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading virtual device context",
			utils.FormatAPIError(fmt.Sprintf("read virtual device context ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was null or empty before (not managed or explicitly cleared),
	// restore that state after mapping.
	// This prevents Terraform from trying to manage fields that aren't in the configuration.
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		tflog.Debug(ctx, "Custom fields unmanaged/cleared, preserving original state during Read", map[string]interface{}{
			"was_null":  originalCustomFields.IsNull(),
			"was_empty": !originalCustomFields.IsNull() && len(originalCustomFields.Elements()) == 0,
		})
		data.CustomFields = originalCustomFields
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *VirtualDeviceContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan VirtualDeviceContextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(plan.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing virtual device context ID",
			fmt.Sprintf("Could not parse ID '%s': %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, plan.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	status := netbox.PatchedWritableVirtualDeviceContextRequestStatus(plan.Status.ValueString())
	apiReq := netbox.NewWritableVirtualDeviceContextRequest(plan.Name.ValueString(), *device, status)

	// Set optional fields
	if !plan.Identifier.IsNull() && !plan.Identifier.IsUnknown() {
		identifier, err := utils.SafeInt32FromValue(plan.Identifier)
		if err != nil {
			resp.Diagnostics.AddError("Invalid identifier", fmt.Sprintf("Identifier overflow: %s", err))
			return
		}
		apiReq.SetIdentifier(identifier)
	} else if plan.Identifier.IsNull() {
		// Explicitly clear identifier
		apiReq.SetIdentifierNil()
	}

	if !plan.Tenant.IsNull() && !plan.Tenant.IsUnknown() {
		tenant, tenantDiags := lookup.LookupTenant(ctx, r.client, plan.Tenant.ValueString())
		resp.Diagnostics.Append(tenantDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTenant(*tenant)
	} else if plan.Tenant.IsNull() {
		// Explicitly clear tenant
		apiReq.SetTenantNil()
	}

	if !plan.PrimaryIP4.IsNull() && !plan.PrimaryIP4.IsUnknown() {
		ipAddr, ipDiags := lookup.LookupIPAddress(ctx, r.client, plan.PrimaryIP4.ValueString())
		resp.Diagnostics.Append(ipDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetPrimaryIp4(*ipAddr)
	} else if plan.PrimaryIP4.IsNull() {
		// Explicitly clear primary_ip4
		apiReq.SetPrimaryIp4Nil()
	}

	if !plan.PrimaryIP6.IsNull() && !plan.PrimaryIP6.IsUnknown() {
		ipAddr, ipDiags := lookup.LookupIPAddress(ctx, r.client, plan.PrimaryIP6.ValueString())
		resp.Diagnostics.Append(ipDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetPrimaryIp6(*ipAddr)
	} else if plan.PrimaryIP6.IsNull() {
		// Explicitly clear primary_ip6
		apiReq.SetPrimaryIp6Nil()
	}

	// Apply common fields individually with merge-aware helpers
	utils.ApplyDescription(apiReq, plan.Description)
	utils.ApplyComments(apiReq, plan.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating virtual device context", map[string]interface{}{
		"id":     id,
		"name":   plan.Name.ValueString(),
		"device": plan.Device.ValueString(),
		"status": plan.Status.ValueString(),
	})

	// Update the resource
	result, httpResp, err := r.client.DcimAPI.DcimVirtualDeviceContextsUpdate(ctx, id).WritableVirtualDeviceContextRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating virtual device context",
			utils.FormatAPIError(fmt.Sprintf("update virtual device context ID %d", id), err, httpResp),
		)
		return
	}

	// Store plan's custom_fields to filter the response
	planCustomFields := plan.CustomFields

	// Map response to state
	r.mapToState(ctx, result, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Filter custom_fields to only those owned by this resource
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, result.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *VirtualDeviceContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualDeviceContextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing virtual device context ID",
			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting virtual device context", map[string]interface{}{"id": id})

	// Delete the resource
	httpResp, err := r.client.DcimAPI.DcimVirtualDeviceContextsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting virtual device context",
			utils.FormatAPIError(fmt.Sprintf("delete virtual device context ID %d", id), err, httpResp),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *VirtualDeviceContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		var id int32
		if _, err := fmt.Sscanf(parsed.ID, "%d", &id); err != nil {
			resp.Diagnostics.AddError("Error parsing virtual device context ID", fmt.Sprintf("Could not parse ID '%s': %s", parsed.ID, err.Error()))
			return
		}

		result, httpResp, err := r.client.DcimAPI.DcimVirtualDeviceContextsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing virtual device context", utils.FormatAPIError(fmt.Sprintf("read virtual device context ID %d", id), err, httpResp))
			return
		}

		var data VirtualDeviceContextResourceModel
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

		r.mapToState(ctx, result, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, result.GetCustomFields(), &resp.Diagnostics)
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

// mapToState maps the API response to the Terraform state.
func (r *VirtualDeviceContextResource) mapToState(ctx context.Context, result *netbox.VirtualDeviceContext, data *VirtualDeviceContextResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())

	// Map device (required field) - preserve user's input format
	device := result.GetDevice()
	data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())

	// Map identifier
	if result.HasIdentifier() {
		identifierPtr, ok := result.GetIdentifierOk()
		if ok && identifierPtr != nil {
			data.Identifier = types.Int64Value(int64(*identifierPtr))
		} else {
			data.Identifier = types.Int64Null()
		}
	} else {
		data.Identifier = types.Int64Null()
	}

	// Map tenant - preserve user's input format
	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Map primary IPs
	if result.HasPrimaryIp4() && result.GetPrimaryIp4().Id != 0 {
		ip := result.GetPrimaryIp4()
		data.PrimaryIP4 = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.PrimaryIP4 = types.StringNull()
	}

	if result.HasPrimaryIp6() && result.GetPrimaryIp6().Id != 0 {
		ip := result.GetPrimaryIp6()
		data.PrimaryIP6 = types.StringValue(fmt.Sprintf("%d", ip.GetId()))
	} else {
		data.PrimaryIP6 = types.StringNull()
	}

	// Map status (required field)
	status := result.GetStatus()
	data.Status = types.StringValue(string(status.GetValue()))

	// Map description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if result.HasComments() && result.GetComments() != "" {
		data.Comments = types.StringValue(result.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags (slug list) with empty-set preservation
	data.Tags = utils.PopulateTagsSlugFromAPI(ctx, result.HasTags(), result.GetTags(), data.Tags)

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, result.GetCustomFields(), diags)
}
