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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &VirtualMachineResource{}

	_ resource.ResourceWithConfigure = &VirtualMachineResource{}

	_ resource.ResourceWithImportState = &VirtualMachineResource{}
)

// NewVirtualMachineResource returns a new Virtual Machine resource.

func NewVirtualMachineResource() resource.Resource {

	return &VirtualMachineResource{}

}

// VirtualMachineResource defines the resource implementation.

type VirtualMachineResource struct {
	client *netbox.APIClient
}

// VirtualMachineResourceModel describes the resource data model.

type VirtualMachineResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Status types.String `tfsdk:"status"`

	Site types.String `tfsdk:"site"`

	SiteID types.String `tfsdk:"site_id"`

	Cluster types.String `tfsdk:"cluster"`

	ClusterID types.String `tfsdk:"cluster_id"`

	Role types.String `tfsdk:"role"`

	RoleID types.String `tfsdk:"role_id"`

	Tenant types.String `tfsdk:"tenant"`

	TenantID types.String `tfsdk:"tenant_id"`

	Platform types.String `tfsdk:"platform"`

	PlatformID types.String `tfsdk:"platform_id"`

	Vcpus types.Float64 `tfsdk:"vcpus"`

	Memory types.Int64 `tfsdk:"memory"`

	Disk types.Int64 `tfsdk:"disk"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *VirtualMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_virtual_machine"

}

// Schema defines the schema for the resource.

func (r *VirtualMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a virtual machine in Netbox. Virtual machines represent virtualized compute instances that run on clusters or hypervisors.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the virtual machine.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": nbschema.NameAttribute("virtual machine", 64),

			"status": schema.StringAttribute{

				MarkdownDescription: "The status of the virtual machine. Valid values are: `offline`, `active`, `planned`, `staged`, `failed`, `decommissioning`. Defaults to `active`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("active"),
			},

			"site": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the site where this virtual machine is located.",

				Optional: true,
			},

			"site_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the site where this virtual machine is located.",

				Computed: true,
			},

			"cluster": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the cluster this virtual machine belongs to.",

				Optional: true,
			},

			"cluster_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the cluster this virtual machine belongs to.",

				Computed: true,
			},

			"role": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the device role for this virtual machine.",

				Optional: true,
			},

			"role_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the device role for this virtual machine.",

				Computed: true,
			},

			"tenant": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the tenant this virtual machine is assigned to.",

				Optional: true,
			},

			"tenant_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the tenant this virtual machine is assigned to.",

				Computed: true,
			},

			"platform": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the platform (operating system) running on this virtual machine.",

				Optional: true,
			},

			"platform_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the platform (operating system) running on this virtual machine.",

				Computed: true,
			},

			"vcpus": schema.Float64Attribute{

				MarkdownDescription: "The number of virtual CPUs allocated to this virtual machine.",

				Optional: true,
			},

			"memory": schema.Int64Attribute{

				MarkdownDescription: "The amount of memory (in MB) allocated to this virtual machine.",

				Optional: true,

				PlanModifiers: []planmodifier.Int64{

					int64planmodifier.UseStateForUnknown(),
				},
			},

			"disk": schema.Int64Attribute{

				MarkdownDescription: "The total disk space (in GB) allocated to this virtual machine.",

				Optional: true,

				PlanModifiers: []planmodifier.Int64{

					int64planmodifier.UseStateForUnknown(),
				},
			},

			"description": nbschema.DescriptionAttribute("virtual machine"),

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about the virtual machine.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure sets up the resource with the provider client.

func (r *VirtualMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

// mapVirtualMachineToState maps a VirtualMachine from the API to the Terraform state model.

func (r *VirtualMachineResource) mapVirtualMachineToState(ctx context.Context, vm *netbox.VirtualMachineWithConfigContext, data *VirtualMachineResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", vm.GetId()))

	data.Name = types.StringValue(vm.GetName())

	// Status

	if vm.HasStatus() {

		data.Status = types.StringValue(string(vm.Status.GetValue()))

	} else {

		data.Status = types.StringValue("active")

	}

	// Site

	if vm.Site.IsSet() && vm.Site.Get() != nil {

		siteObj := vm.Site.Get()

		data.SiteID = types.StringValue(fmt.Sprintf("%d", siteObj.GetId()))

		data.Site = utils.UpdateReferenceAttribute(data.Site, siteObj.GetName(), siteObj.GetSlug(), siteObj.GetId())

	} else {

		data.Site = types.StringNull()

		data.SiteID = types.StringNull()

	}

	// Cluster

	if vm.Cluster.IsSet() && vm.Cluster.Get() != nil {

		clusterObj := vm.Cluster.Get()

		data.ClusterID = types.StringValue(fmt.Sprintf("%d", clusterObj.GetId()))

		data.Cluster = utils.UpdateReferenceAttribute(data.Cluster, clusterObj.GetName(), "", clusterObj.GetId())

	} else {

		data.Cluster = types.StringNull()

		data.ClusterID = types.StringNull()

	}

	// Role

	if vm.Role.IsSet() && vm.Role.Get() != nil {

		roleObj := vm.Role.Get()

		data.RoleID = types.StringValue(fmt.Sprintf("%d", roleObj.GetId()))

		data.Role = utils.UpdateReferenceAttribute(data.Role, roleObj.GetName(), roleObj.GetSlug(), roleObj.GetId())

	} else {

		data.Role = types.StringNull()

		data.RoleID = types.StringNull()

	}

	// Tenant

	if vm.Tenant.IsSet() && vm.Tenant.Get() != nil {

		tenantObj := vm.Tenant.Get()

		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenantObj.GetId()))

		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenantObj.GetName(), tenantObj.GetSlug(), tenantObj.GetId())

	} else {

		data.Tenant = types.StringNull()

		data.TenantID = types.StringNull()

	}

	// Platform

	if vm.Platform.IsSet() && vm.Platform.Get() != nil {

		platformObj := vm.Platform.Get()

		data.PlatformID = types.StringValue(fmt.Sprintf("%d", platformObj.GetId()))

		data.Platform = utils.UpdateReferenceAttribute(data.Platform, platformObj.GetName(), platformObj.GetSlug(), platformObj.GetId())

	} else {

		data.Platform = types.StringNull()

		data.PlatformID = types.StringNull()

	}

	// Vcpus

	if vm.Vcpus.IsSet() && vm.Vcpus.Get() != nil {

		data.Vcpus = types.Float64Value(*vm.Vcpus.Get())

	} else {

		data.Vcpus = types.Float64Null()

	}

	// Memory

	if vm.Memory.IsSet() && vm.Memory.Get() != nil {

		data.Memory = types.Int64Value(int64(*vm.Memory.Get()))

	} else {

		data.Memory = types.Int64Null()

	}

	// Disk

	if vm.Disk.IsSet() && vm.Disk.Get() != nil {

		data.Disk = types.Int64Value(int64(*vm.Disk.Get()))

	} else {

		data.Disk = types.Int64Null()

	}

	// Description

	if vm.HasDescription() && vm.GetDescription() != "" {

		data.Description = types.StringValue(vm.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Comments

	if vm.HasComments() && vm.GetComments() != "" {

		data.Comments = types.StringValue(vm.GetComments())

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if vm.HasTags() {

		tags := utils.NestedTagsToTagModels(vm.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	if vm.HasCustomFields() && !data.CustomFields.IsNull() {

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		customFields := utils.MapToCustomFieldModels(vm.GetCustomFields(), stateCustomFields)

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

// buildVirtualMachineRequest builds a WritableVirtualMachineWithConfigContextRequest from the resource model.

func (r *VirtualMachineResource) buildVirtualMachineRequest(ctx context.Context, data *VirtualMachineResourceModel, diags *diag.Diagnostics) *netbox.WritableVirtualMachineWithConfigContextRequest {

	vmRequest := &netbox.WritableVirtualMachineWithConfigContextRequest{

		Name: data.Name.ValueString(),
	}

	// Status

	if utils.IsSet(data.Status) {

		status := netbox.ModuleStatusValue(data.Status.ValueString())

		vmRequest.Status = &status

	}

	// Site

	if utils.IsSet(data.Site) {

		site, siteDiags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())

		diags.Append(siteDiags...)

		if diags.HasError() {

			return nil

		}

		vmRequest.Site = *netbox.NewNullableBriefSiteRequest(site)

	}

	// Cluster

	if utils.IsSet(data.Cluster) {

		cluster, clusterDiags := netboxlookup.LookupCluster(ctx, r.client, data.Cluster.ValueString())

		diags.Append(clusterDiags...)

		if diags.HasError() {

			return nil

		}

		vmRequest.Cluster = *netbox.NewNullableBriefClusterRequest(cluster)

	}

	// Role

	if utils.IsSet(data.Role) {

		role, roleDiags := netboxlookup.LookupDeviceRole(ctx, r.client, data.Role.ValueString())

		diags.Append(roleDiags...)

		if diags.HasError() {

			return nil

		}

		vmRequest.Role = *netbox.NewNullableBriefDeviceRoleRequest(role)

	}

	// Tenant

	if utils.IsSet(data.Tenant) {

		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		diags.Append(tenantDiags...)

		if diags.HasError() {

			return nil

		}

		vmRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)

	}

	// Platform

	if utils.IsSet(data.Platform) {

		platform, platformDiags := netboxlookup.LookupPlatform(ctx, r.client, data.Platform.ValueString())

		diags.Append(platformDiags...)

		if diags.HasError() {

			return nil

		}

		vmRequest.Platform = *netbox.NewNullableBriefPlatformRequest(platform)

	}

	// Vcpus

	if utils.IsSet(data.Vcpus) {

		vcpus := data.Vcpus.ValueFloat64()

		vmRequest.Vcpus = *netbox.NewNullableFloat64(&vcpus)

	}

	// Memory

	if utils.IsSet(data.Memory) {

		memory, err := utils.SafeInt32FromValue(data.Memory)

		if err != nil {

			diags.AddError("Invalid memory value", fmt.Sprintf("Memory value overflow: %s", err))

			return nil

		}

		vmRequest.Memory = *netbox.NewNullableInt32(&memory)

	}

	// Disk

	if utils.IsSet(data.Disk) {

		disk, err := utils.SafeInt32FromValue(data.Disk)

		if err != nil {

			diags.AddError("Invalid disk value", fmt.Sprintf("Disk value overflow: %s", err))

			return nil

		}

		vmRequest.Disk = *netbox.NewNullableInt32(&disk)

	}

	// Description

	if utils.IsSet(data.Description) {

		description := data.Description.ValueString()

		vmRequest.Description = &description

	}

	// Comments

	if utils.IsSet(data.Comments) {

		comments := data.Comments.ValueString()

		vmRequest.Comments = &comments

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		var tags []utils.TagModel

		diags.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if diags.HasError() {

			return nil

		}

		vmRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFields []utils.CustomFieldModel

		diags.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if diags.HasError() {

			return nil

		}

		vmRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	return vmRequest

}

// Create creates a new virtual machine resource.

func (r *VirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data VirtualMachineResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating virtual machine", map[string]interface{}{

		"name": data.Name.ValueString(),
	})

	// Build the VM request

	vmRequest := r.buildVirtualMachineRequest(ctx, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Call the API

	vm, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesCreate(ctx).WritableVirtualMachineWithConfigContextRequest(*vmRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating virtual machine",

			utils.FormatAPIError("create virtual machine", err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Created virtual machine", map[string]interface{}{

		"id": vm.GetId(),

		"name": vm.GetName(),
	})

	// Map response to state

	r.mapVirtualMachineToState(ctx, vm, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *VirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data VirtualMachineResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	vmID := data.ID.ValueString()

	var vmIDInt int32

	vmIDInt, err := utils.ParseID(vmID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Virtual Machine ID",

			fmt.Sprintf("Virtual Machine ID must be a number, got: %s", vmID),
		)

		return

	}

	tflog.Debug(ctx, "Reading virtual machine", map[string]interface{}{

		"id": vmID,
	})

	// Call the API

	vm, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx, vmIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Virtual machine not found, removing from state", map[string]interface{}{

				"id": vmID,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading virtual machine",

			utils.FormatAPIError(fmt.Sprintf("read virtual machine ID %s", vmID), err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapVirtualMachineToState(ctx, vm, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state.

func (r *VirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data VirtualMachineResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	vmID := data.ID.ValueString()

	var vmIDInt int32

	vmIDInt, err := utils.ParseID(vmID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Virtual Machine ID",

			fmt.Sprintf("Virtual Machine ID must be a number, got: %s", vmID),
		)

		return

	}

	tflog.Debug(ctx, "Updating virtual machine", map[string]interface{}{

		"id": vmID,

		"name": data.Name.ValueString(),
	})

	// Build the VM request

	vmRequest := r.buildVirtualMachineRequest(ctx, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Call the API

	vm, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesUpdate(ctx, vmIDInt).WritableVirtualMachineWithConfigContextRequest(*vmRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating virtual machine",

			utils.FormatAPIError(fmt.Sprintf("update virtual machine ID %s", vmID), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Updated virtual machine", map[string]interface{}{

		"id": vm.GetId(),

		"name": vm.GetName(),
	})

	// Map response to state

	r.mapVirtualMachineToState(ctx, vm, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state.

func (r *VirtualMachineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data VirtualMachineResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	vmID := data.ID.ValueString()

	var vmIDInt int32

	vmIDInt, err := utils.ParseID(vmID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Virtual Machine ID",

			fmt.Sprintf("Virtual Machine ID must be a number, got: %s", vmID),
		)

		return

	}

	tflog.Debug(ctx, "Deleting virtual machine", map[string]interface{}{

		"id": vmID,
	})

	// Call the API

	httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualMachinesDestroy(ctx, vmIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Already deleted, consider success

			tflog.Debug(ctx, "Virtual machine already deleted", map[string]interface{}{

				"id": vmID,
			})

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting virtual machine",

			utils.FormatAPIError(fmt.Sprintf("delete virtual machine ID %s", vmID), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted virtual machine", map[string]interface{}{

		"id": vmID,
	})

}

// ImportState imports an existing resource into Terraform.

func (r *VirtualMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}
