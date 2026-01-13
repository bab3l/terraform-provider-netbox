// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"strconv"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &DeviceBayTemplateResource{}
	_ resource.ResourceWithConfigure   = &DeviceBayTemplateResource{}
	_ resource.ResourceWithImportState = &DeviceBayTemplateResource{}
)

// NewDeviceBayTemplateResource returns a new DeviceBayTemplate resource.
func NewDeviceBayTemplateResource() resource.Resource {
	return &DeviceBayTemplateResource{}
}

// DeviceBayTemplateResource defines the resource implementation.
type DeviceBayTemplateResource struct {
	client *netbox.APIClient
}

// DeviceBayTemplateResourceModel describes the resource data model.
type DeviceBayTemplateResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DeviceType  types.String `tfsdk:"device_type"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.
func (r *DeviceBayTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_bay_template"
}

// Schema defines the schema for the resource.
func (r *DeviceBayTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Device Bay Template in Netbox. Device bay templates define device bays that will be created on devices of the associated device type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the device bay template.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_type": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device_type",
				"The ID or slug of the device type this template belongs to.",
			),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the device bay template. Use {module} as a substitution for the module bay position when attached to a module type.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label for the device bay.",
				Optional:            true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("device bay template"))
}

// Configure adds the provider configured client to the resource.
func (r *DeviceBayTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DeviceBayTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceBayTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Look up device type
	deviceTypeRef, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request
	createReq := *netbox.NewDeviceBayTemplateRequest(*deviceTypeRef, data.Name.ValueString())

	// Set optional fields
	if utils.IsSet(data.Label) {
		label := data.Label.ValueString()
		createReq.Label = &label
	}

	// Apply description
	utils.ApplyDescription(&createReq, data.Description)
	tflog.Debug(ctx, "Creating DeviceBayTemplate", map[string]interface{}{
		"device_type": data.DeviceType.ValueString(),
		"name":        data.Name.ValueString(),
	})

	// Create the device bay template
	template, httpResp, err := r.client.DcimAPI.DcimDeviceBayTemplatesCreate(ctx).
		DeviceBayTemplateRequest(createReq).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating DeviceBayTemplate",
			utils.FormatAPIError("create device bay template", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapTemplateToModel(template, &data)
	tflog.Debug(ctx, "Created DeviceBayTemplate", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *DeviceBayTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceBayTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing DeviceBayTemplate ID",
			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading DeviceBayTemplate", map[string]interface{}{
		"id": id,
	})

	// Get the device bay template
	id32, err := utils.SafeInt32(int64(id))
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
		return
	}

	template, httpResp, err := r.client.DcimAPI.DcimDeviceBayTemplatesRetrieve(ctx, id32).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "DeviceBayTemplate not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading DeviceBayTemplate",
			utils.FormatAPIError(fmt.Sprintf("read device bay template %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapTemplateToModel(template, &data)
	tflog.Debug(ctx, "Read DeviceBayTemplate", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DeviceBayTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceBayTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing DeviceBayTemplate ID",
			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)
		return
	}

	// Look up device type
	deviceTypeRef, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the update request
	updateReq := *netbox.NewDeviceBayTemplateRequest(*deviceTypeRef, data.Name.ValueString())

	// Set optional fields - use empty string to clear
	if utils.IsSet(data.Label) {
		label := data.Label.ValueString()
		updateReq.Label = &label
	} else {
		emptyLabel := ""
		updateReq.Label = &emptyLabel
	}

	// Apply description
	utils.ApplyDescription(&updateReq, data.Description)
	tflog.Debug(ctx, "Updating DeviceBayTemplate", map[string]interface{}{
		"id": id,
	})

	// Update the device bay template
	id32, convErr := utils.SafeInt32(int64(id))
	if convErr != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", convErr))
		return
	}

	template, httpResp, err := r.client.DcimAPI.DcimDeviceBayTemplatesUpdate(ctx, id32).
		DeviceBayTemplateRequest(updateReq).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating DeviceBayTemplate",
			utils.FormatAPIError(fmt.Sprintf("update device bay template %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapTemplateToModel(template, &data)
	tflog.Debug(ctx, "Updated DeviceBayTemplate", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *DeviceBayTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceBayTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing DeviceBayTemplate ID",
			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting DeviceBayTemplate", map[string]interface{}{
		"id": id,
	})

	// Delete the device bay template
	id32, convErr := utils.SafeInt32(int64(id))
	if convErr != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", convErr))
		return
	}
	httpResp, err := r.client.DcimAPI.DcimDeviceBayTemplatesDestroy(ctx, id32).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "DeviceBayTemplate already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting DeviceBayTemplate",
			utils.FormatAPIError(fmt.Sprintf("delete device bay template %d", id), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted DeviceBayTemplate", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state from Terraform.
func (r *DeviceBayTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapTemplateToModel maps a Netbox DeviceBayTemplate to the Terraform resource model.
func (r *DeviceBayTemplateResource) mapTemplateToModel(template *netbox.DeviceBayTemplate, data *DeviceBayTemplateResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", template.Id))
	data.DeviceType = utils.UpdateReferenceAttribute(data.DeviceType, template.DeviceType.GetModel(), template.DeviceType.GetSlug(), template.DeviceType.GetId())
	data.Name = types.StringValue(template.Name)

	// Label
	if template.Label != nil && *template.Label != "" {
		data.Label = types.StringValue(*template.Label)
	} else {
		data.Label = types.StringNull()
	}

	// Description
	if template.Description != nil && *template.Description != "" {
		data.Description = types.StringValue(*template.Description)
	} else {
		data.Description = types.StringNull()
	}
}
