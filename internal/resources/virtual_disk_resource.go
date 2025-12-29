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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &VirtualDiskResource{}

	_ resource.ResourceWithConfigure = &VirtualDiskResource{}

	_ resource.ResourceWithImportState = &VirtualDiskResource{}
)

// NewVirtualDiskResource returns a new VirtualDisk resource.

func NewVirtualDiskResource() resource.Resource {
	return &VirtualDiskResource{}
}

// VirtualDiskResource defines the resource implementation.

type VirtualDiskResource struct {
	client *netbox.APIClient
}

// VirtualDiskResourceModel describes the resource data model.

type VirtualDiskResourceModel struct {
	ID types.String `tfsdk:"id"`

	VirtualMachine types.String `tfsdk:"virtual_machine"`

	Name types.String `tfsdk:"name"`

	Size types.String `tfsdk:"size"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *VirtualDiskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_disk"
}

// Schema defines the schema for the resource.

func (r *VirtualDiskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual disk attached to a virtual machine in Netbox. Virtual disks represent storage volumes associated with VMs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the virtual disk.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"virtual_machine": schema.StringAttribute{
				MarkdownDescription: "ID or name of the virtual machine this disk belongs to. Required.",

				Required: true,
			},

			"name": nbschema.NameAttribute("virtual disk", 64),

			"size": schema.StringAttribute{
				MarkdownDescription: "Size of the virtual disk in GB. Required.",

				Required: true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("virtual disk"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.

func (r *VirtualDiskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VirtualDiskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualDiskResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup virtual machine

	vmRef, vmDiags := netboxlookup.LookupVirtualMachine(ctx, r.client, data.VirtualMachine.ValueString())

	resp.Diagnostics.Append(vmDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse size to int32

	var size int32

	if _, err := fmt.Sscanf(data.Size.ValueString(), "%d", &size); err != nil {
		resp.Diagnostics.AddError(

			"Invalid Size",

			fmt.Sprintf("Unable to parse size %q: %s", data.Size.ValueString(), err.Error()),
		)

		return
	}

	// Create the VirtualDisk request

	vdRequest := netbox.NewVirtualDiskRequest(*vmRef, data.Name.ValueString(), size)

	// Set optional fields

	r.setOptionalFields(ctx, vdRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating VirtualDisk", map[string]interface{}{
		"name": data.Name.ValueString(),

		"virtual_machine": data.VirtualMachine.ValueString(),

		"size": size,
	})

	// Create the VirtualDisk

	vd, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualDisksCreate(ctx).VirtualDiskRequest(*vdRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating VirtualDisk",

			utils.FormatAPIError("create VirtualDisk", err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapVirtualDiskToState(ctx, vd, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Created VirtualDisk", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.

func (r *VirtualDiskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualDiskResourceModel

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

	tflog.Debug(ctx, "Reading VirtualDisk", map[string]interface{}{
		"id": id,
	})

	// Get the VirtualDisk from Netbox

	vd, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualDisksRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading VirtualDisk",

			utils.FormatAPIError(fmt.Sprintf("read VirtualDisk ID %d", id), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapVirtualDiskToState(ctx, vd, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Read VirtualDisk", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.

func (r *VirtualDiskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtualDiskResourceModel

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

	// Lookup virtual machine

	vmRef, vmDiags := netboxlookup.LookupVirtualMachine(ctx, r.client, data.VirtualMachine.ValueString())

	resp.Diagnostics.Append(vmDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse size to int32

	var size int32

	if _, err := fmt.Sscanf(data.Size.ValueString(), "%d", &size); err != nil {
		resp.Diagnostics.AddError(

			"Invalid Size",

			fmt.Sprintf("Unable to parse size %q: %s", data.Size.ValueString(), err.Error()),
		)

		return
	}

	// Create the VirtualDisk request

	vdRequest := netbox.NewVirtualDiskRequest(*vmRef, data.Name.ValueString(), size)

	// Set optional fields

	r.setOptionalFields(ctx, vdRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating VirtualDisk", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Update the VirtualDisk

	vd, httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualDisksUpdate(ctx, id).VirtualDiskRequest(*vdRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating VirtualDisk",

			utils.FormatAPIError(fmt.Sprintf("update VirtualDisk ID %d", id), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapVirtualDiskToState(ctx, vd, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated VirtualDisk", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.

func (r *VirtualDiskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualDiskResourceModel

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

	tflog.Debug(ctx, "Deleting VirtualDisk", map[string]interface{}{
		"id": id,
	})

	// Delete the VirtualDisk

	httpResp, err := r.client.VirtualizationAPI.VirtualizationVirtualDisksDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return // Already deleted
		}

		resp.Diagnostics.AddError(

			"Error deleting VirtualDisk",

			utils.FormatAPIError(fmt.Sprintf("delete VirtualDisk ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted VirtualDisk", map[string]interface{}{
		"id": id,
	})
}

func (r *VirtualDiskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the VirtualDisk request from the resource model.

func (r *VirtualDiskResource) setOptionalFields(ctx context.Context, vdRequest *netbox.VirtualDiskRequest, data *VirtualDiskResourceModel, diags *diag.Diagnostics) {
	// Apply description and metadata fields

	utils.ApplyDescription(vdRequest, data.Description)

	utils.ApplyMetadataFields(ctx, vdRequest, data.Tags, data.CustomFields, diags)
}

// mapVirtualDiskToState maps a Netbox VirtualDisk to the Terraform state model.

func (r *VirtualDiskResource) mapVirtualDiskToState(ctx context.Context, vd *netbox.VirtualDisk, data *VirtualDiskResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vd.Id))

	data.Name = types.StringValue(vd.Name)

	// DisplayName
	if vd.Display != "" {
	} else {
	}

	// VirtualMachine - preserve user's input format

	data.VirtualMachine = utils.UpdateReferenceAttribute(data.VirtualMachine, vd.VirtualMachine.GetName(), "", vd.VirtualMachine.GetId())

	data.Size = types.StringValue(fmt.Sprintf("%d", vd.Size))

	// Description

	if vd.Description != nil && *vd.Description != "" {
		data.Description = types.StringValue(*vd.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags

	if len(vd.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(vd.Tags)

		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields

	switch {
	case len(vd.CustomFields) > 0 && !data.CustomFields.IsNull():

		var stateCustomFields []utils.CustomFieldModel

		data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		customFields := utils.MapToCustomFieldModels(vd.CustomFields, stateCustomFields)

		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		data.CustomFields = customFieldsValue

	case len(vd.CustomFields) > 0:

		customFields := utils.MapToCustomFieldModels(vd.CustomFields, []utils.CustomFieldModel{})

		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		data.CustomFields = customFieldsValue

	default:

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
