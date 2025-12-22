// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ServiceResource{}

	_ resource.ResourceWithConfigure = &ServiceResource{}

	_ resource.ResourceWithImportState = &ServiceResource{}
)

// NewServiceResource returns a new resource implementing the service resource.

func NewServiceResource() resource.Resource {

	return &ServiceResource{}

}

// ServiceResource defines the resource implementation.

type ServiceResource struct {
	client *netbox.APIClient
}

// ServiceResourceModel describes the resource data model.

type ServiceResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	VirtualMachine types.String `tfsdk:"virtual_machine"`

	Name types.String `tfsdk:"name"`

	DisplayName types.String `tfsdk:"display_name"`

	Protocol types.String `tfsdk:"protocol"`

	Ports types.List `tfsdk:"ports"`

	IPAddresses types.List `tfsdk:"ipaddresses"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_service"

}

// Schema defines the schema for the resource.

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a network service in NetBox. Services represent TCP/UDP services running on devices or virtual machines.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the service.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The device this service runs on (ID or name). Mutually exclusive with virtual_machine.",

				Optional: true,
			},

			"virtual_machine": schema.StringAttribute{

				MarkdownDescription: "The virtual machine this service runs on (ID or name). Mutually exclusive with device.",

				Optional: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the service (e.g., 'ssh', 'http', 'https').",

				Required: true,
			},

			"display_name": nbschema.DisplayNameAttribute("service"),

			"protocol": schema.StringAttribute{

				MarkdownDescription: "The protocol used by the service. Valid values: `tcp`, `udp`, `sctp`.",

				Required: true,
			},

			"ports": schema.ListAttribute{

				MarkdownDescription: "List of port numbers the service listens on.",

				Required: true,

				ElementType: types.Int64Type,
			},

			"ipaddresses": schema.ListAttribute{

				MarkdownDescription: "List of IP address IDs associated with this service.",

				Optional: true,

				ElementType: types.Int64Type,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the service.",

				Optional: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about the service.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *ServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Extract ports

	var ports []int64

	resp.Diagnostics.Append(data.Ports.ElementsAs(ctx, &ports, false)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portsInt32 := make([]int32, len(ports))

	for i, p := range ports {

		p32, err := utils.SafeInt32(p)

		if err != nil {

			resp.Diagnostics.AddError("Invalid port", fmt.Sprintf("Port value overflow: %s", err))

			return

		}

		portsInt32[i] = p32

	}

	// Build request

	protocol := netbox.PatchedWritableServiceRequestProtocol(data.Protocol.ValueString())

	apiReq := netbox.NewWritableServiceRequest(data.Name.ValueString(), protocol, portsInt32)

	// Set device or virtual_machine

	if !data.Device.IsNull() && !data.Device.IsUnknown() {

		device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetDevice(*device)

	}

	if !data.VirtualMachine.IsNull() && !data.VirtualMachine.IsUnknown() {

		vm, diags := lookup.LookupVirtualMachine(ctx, r.client, data.VirtualMachine.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetVirtualMachine(*vm)

	}

	// Set IP addresses

	if !data.IPAddresses.IsNull() && !data.IPAddresses.IsUnknown() {

		var ipIDs []int64

		resp.Diagnostics.Append(data.IPAddresses.ElementsAs(ctx, &ipIDs, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		ipIDsInt32 := make([]int32, len(ipIDs))

		for i, id := range ipIDs {

			id32, err := utils.SafeInt32(id)

			if err != nil {

				resp.Diagnostics.AddError("Invalid IP address ID", fmt.Sprintf("IP address ID overflow: %s", err))

				return

			}

			ipIDsInt32[i] = id32

		}

		apiReq.SetIpaddresses(ipIDsInt32)

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

	}

	tflog.Debug(ctx, "Creating service", map[string]interface{}{

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.IpamAPI.IpamServicesCreate(ctx).WritableServiceRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating service",

			utils.FormatAPIError(fmt.Sprintf("create service %s", data.Name.ValueString()), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created service", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	svcID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Service ID",

			fmt.Sprintf("Service ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading service", map[string]interface{}{

		"id": svcID,
	})

	response, httpResp, err := r.client.IpamAPI.IpamServicesRetrieve(ctx, svcID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading service",

			utils.FormatAPIError(fmt.Sprintf("read service ID %d", svcID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource.

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	svcID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Service ID",

			fmt.Sprintf("Service ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	// Extract ports

	var ports []int64

	resp.Diagnostics.Append(data.Ports.ElementsAs(ctx, &ports, false)...)

	if resp.Diagnostics.HasError() {

		return

	}

	portsInt32 := make([]int32, len(ports))

	for i, p := range ports {

		p32, err := utils.SafeInt32(p)

		if err != nil {

			resp.Diagnostics.AddError("Invalid port number", fmt.Sprintf("Port number overflow: %s", err))

			return

		}

		portsInt32[i] = p32

	}

	// Build request

	protocol := netbox.PatchedWritableServiceRequestProtocol(data.Protocol.ValueString())

	apiReq := netbox.NewWritableServiceRequest(data.Name.ValueString(), protocol, portsInt32)

	// Set device or virtual_machine

	if !data.Device.IsNull() && !data.Device.IsUnknown() {

		device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetDevice(*device)

	}

	if !data.VirtualMachine.IsNull() && !data.VirtualMachine.IsUnknown() {

		vm, diags := lookup.LookupVirtualMachine(ctx, r.client, data.VirtualMachine.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetVirtualMachine(*vm)

	}

	// Set IP addresses

	if !data.IPAddresses.IsNull() && !data.IPAddresses.IsUnknown() {

		var ipIDs []int64

		resp.Diagnostics.Append(data.IPAddresses.ElementsAs(ctx, &ipIDs, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		ipIDsInt32 := make([]int32, len(ipIDs))

		for i, id := range ipIDs {

			id32, err := utils.SafeInt32(id)

			if err != nil {

				resp.Diagnostics.AddError("Invalid IP address ID", fmt.Sprintf("IP address ID overflow: %s", err))

				return

			}

			ipIDsInt32[i] = id32

		}

		apiReq.SetIpaddresses(ipIDsInt32)

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

	}

	tflog.Debug(ctx, "Updating service", map[string]interface{}{

		"id": svcID,
	})

	response, httpResp, err := r.client.IpamAPI.IpamServicesUpdate(ctx, svcID).WritableServiceRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating service",

			utils.FormatAPIError(fmt.Sprintf("update service ID %d", svcID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource.

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	svcID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Service ID",

			fmt.Sprintf("Service ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting service", map[string]interface{}{

		"id": svcID,
	})

	httpResp, err := r.client.IpamAPI.IpamServicesDestroy(ctx, svcID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting service",

			utils.FormatAPIError(fmt.Sprintf("delete service ID %d", svcID), err, httpResp),
		)

		return

	}

}

// ImportState imports an existing resource.

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	svcID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Service ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.IpamAPI.IpamServicesRetrieve(ctx, svcID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing service",

			utils.FormatAPIError(fmt.Sprintf("import service ID %d", svcID), err, httpResp),
		)

		return

	}

	var data ServiceResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *ServiceResource) mapResponseToModel(ctx context.Context, svc *netbox.Service, data *ServiceResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", svc.GetId()))

	data.Name = types.StringValue(svc.GetName())

	// Map device

	// Map device - preserve user's input format

	if svc.Device.IsSet() && svc.Device.Get() != nil {

		device := svc.Device.Get()

		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())

	} else {

		data.Device = types.StringNull()

	}

	// Map virtual_machine - preserve user's input format

	if svc.VirtualMachine.IsSet() && svc.VirtualMachine.Get() != nil {

		vm := svc.VirtualMachine.Get()

		data.VirtualMachine = utils.UpdateReferenceAttribute(data.VirtualMachine, vm.GetName(), "", vm.GetId())

	} else {

		data.VirtualMachine = types.StringNull()

	}

	// Map protocol

	if protocol, ok := svc.GetProtocolOk(); ok && protocol != nil {

		data.Protocol = types.StringValue(string(protocol.GetValue()))

	} else {

		data.Protocol = types.StringNull()

	}

	// Map ports

	ports := svc.GetPorts()

	portsInt64 := make([]int64, len(ports))

	for i, p := range ports {

		portsInt64[i] = int64(p)

	}

	portsValue, portsDiags := types.ListValueFrom(ctx, types.Int64Type, portsInt64)

	diags.Append(portsDiags...)

	if diags.HasError() {

		return

	}

	data.Ports = portsValue

	// Map IP addresses

	if svc.HasIpaddresses() && len(svc.GetIpaddresses()) > 0 {

		ipAddrs := svc.GetIpaddresses()

		ipIDs := make([]int64, len(ipAddrs))

		for i, ip := range ipAddrs {

			ipIDs[i] = int64(ip.GetId())

		}

		ipValue, ipDiags := types.ListValueFrom(ctx, types.Int64Type, ipIDs)

		diags.Append(ipDiags...)

		if diags.HasError() {

			return

		}

		data.IPAddresses = ipValue

	} else {

		data.IPAddresses = types.ListNull(types.Int64Type)

	}

	// Map description

	if desc, ok := svc.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map comments

	if comments, ok := svc.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if svc.HasTags() && len(svc.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(svc.GetTags())

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

	if svc.HasCustomFields() {

		apiCustomFields := svc.GetCustomFields()

		var stateCustomFieldModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

			data.CustomFields.ElementsAs(ctx, &stateCustomFieldModels, false)

		}

		customFields := utils.MapToCustomFieldModels(apiCustomFields, stateCustomFieldModels)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
