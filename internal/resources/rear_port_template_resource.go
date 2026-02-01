// Package resources provides Terraform resource implementations for NetBox objects.

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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RearPortTemplateResource{}
	_ resource.ResourceWithImportState = &RearPortTemplateResource{}
)

// NewRearPortTemplateResource returns a new resource implementing the rear port template resource.
func NewRearPortTemplateResource() resource.Resource {
	return &RearPortTemplateResource{}
}

// RearPortTemplateResource defines the resource implementation.
type RearPortTemplateResource struct {
	client *netbox.APIClient
}

// RearPortTemplateResourceModel describes the resource data model.
type RearPortTemplateResourceModel struct {
	ID          types.Int32  `tfsdk:"id"`
	DeviceType  types.String `tfsdk:"device_type"`
	ModuleType  types.String `tfsdk:"module_type"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Type        types.String `tfsdk:"type"`
	Color       types.String `tfsdk:"color"`
	Positions   types.Int32  `tfsdk:"positions"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.
func (r *RearPortTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rear_port_template"
}

// Schema defines the schema for the resource.
func (r *RearPortTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rear port template in NetBox. Rear port templates define the default rear ports that will be created on new devices of a specific device type or modules of a module type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the rear port template.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"device_type": nbschema.ReferenceAttributeWithDiffSuppress(
				"device_type",
				"The device type ID or slug. Either device_type or module_type must be specified.",
			),
			"module_type": nbschema.ReferenceAttributeWithDiffSuppress(
				"module_type",
				"The module type ID or model name. Either device_type or module_type must be specified.",
			),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the rear port template. Use {module} as a substitution for the module bay position when attached to a module type.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the rear port template.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of rear port (e.g., `8p8c`, `8p6c`, `110-punch`, `bnc`, `f`, `n`, `mrj21`, `fc`, `lc`, `lc-pc`, `lc-upc`, `lc-apc`, `lsh`, `lsh-pc`, `lsh-upc`, `lsh-apc`, `mpo`, `mtrj`, `sc`, `sc-pc`, `sc-upc`, `sc-apc`, `st`, `cs`, `sn`, `sma-905`, `sma-906`, `splice`, `other`).",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the rear port in hex format (e.g., `aa1409`).",
				Optional:            true,
			},
			"positions": schema.Int32Attribute{
				MarkdownDescription: "Number of front ports that may be mapped to this rear port (1-1024). Default is 1.",
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(1),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("rear port template"))
}

// Configure adds the provider configured client to the resource.
func (r *RearPortTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *RearPortTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RearPortTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	apiReq := netbox.NewWritableRearPortTemplateRequest(data.Name.ValueString(), netbox.FrontPortTypeValue(data.Type.ValueString()))

	// Set device type or module type
	if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
		deviceType, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetDeviceType(*deviceType)
	}

	if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {
		moduleType, diags := netboxlookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetModuleType(*moduleType)
	}

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)
	if utils.IsSet(data.Color) {
		apiReq.SetColor(data.Color.ValueString())
	} else if data.Color.IsNull() {
		// Explicitly clear color when removed from config
		apiReq.SetColor("")
	}
	if !data.Positions.IsNull() && !data.Positions.IsUnknown() {
		apiReq.SetPositions(data.Positions.ValueInt32())
	}

	// Apply description
	utils.ApplyDescription(apiReq, data.Description)
	tflog.Debug(ctx, "Creating rear port template", map[string]interface{}{
		"name": data.Name.ValueString(),
	})
	response, httpResp, err := r.client.DcimAPI.DcimRearPortTemplatesCreate(ctx).WritableRearPortTemplateRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating rear port template",
			utils.FormatAPIError("create rear port template", err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "create rear port template", httpResp, http.StatusCreated) {
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)
	tflog.Trace(ctx, "Created rear port template", map[string]interface{}{
		"id": data.ID.ValueInt32(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *RearPortTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RearPortTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID := data.ID.ValueInt32()
	tflog.Debug(ctx, "Reading rear port template", map[string]interface{}{
		"id": templateID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimRearPortTemplatesRetrieve(ctx, templateID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() {
			tflog.Debug(ctx, "Rear port template not found, removing from state", map[string]interface{}{
				"id": templateID,
			})
			resp.State.RemoveResource(ctx)
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error reading rear port template",
			utils.FormatAPIError(fmt.Sprintf("read rear port template ID %d", templateID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read rear port template", httpResp, http.StatusOK) {
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *RearPortTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RearPortTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID := data.ID.ValueInt32()

	// Build the API request
	apiReq := netbox.NewWritableRearPortTemplateRequest(data.Name.ValueString(), netbox.FrontPortTypeValue(data.Type.ValueString()))

	// Set device type or module type
	if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
		deviceType, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetDeviceType(*deviceType)
	}

	if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {
		moduleType, diags := netboxlookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetModuleType(*moduleType)
	}

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)
	if utils.IsSet(data.Color) {
		apiReq.SetColor(data.Color.ValueString())
	} else if data.Color.IsNull() {
		// Explicitly clear color when removed from config
		apiReq.SetColor("")
	}
	if !data.Positions.IsNull() && !data.Positions.IsUnknown() {
		apiReq.SetPositions(data.Positions.ValueInt32())
	}

	// Apply description
	utils.ApplyDescription(apiReq, data.Description)
	tflog.Debug(ctx, "Updating rear port template", map[string]interface{}{
		"id": templateID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimRearPortTemplatesUpdate(ctx, templateID).WritableRearPortTemplateRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating rear port template",
			utils.FormatAPIError(fmt.Sprintf("update rear port template ID %d", templateID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "update rear port template", httpResp, http.StatusOK) {
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *RearPortTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RearPortTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID := data.ID.ValueInt32()
	tflog.Debug(ctx, "Deleting rear port template", map[string]interface{}{
		"id": templateID,
	})
	httpResp, err := r.client.DcimAPI.DcimRearPortTemplatesDestroy(ctx, templateID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, nil) {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting rear port template",
			utils.FormatAPIError(fmt.Sprintf("delete rear port template ID %d", templateID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "delete rear port template", httpResp, http.StatusNoContent) {
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *RearPortTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID as an integer
	id, err := utils.ParseInt32ID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Could not parse import ID %q as integer: %s", req.ID, err),
		)
		return
	}

	// Set the ID in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *RearPortTemplateResource) mapResponseToModel(template *netbox.RearPortTemplate, data *RearPortTemplateResourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())

	// Map device type - preserve user's input format
	if template.DeviceType.IsSet() && template.DeviceType.Get() != nil {
		dt := template.DeviceType.Get()
		data.DeviceType = utils.UpdateReferenceAttribute(data.DeviceType, dt.GetSlug(), "", dt.Id)
	} else {
		data.DeviceType = types.StringNull()
	}

	// Map module type - preserve user's input format
	if template.ModuleType.IsSet() && template.ModuleType.Get() != nil {
		mt := template.ModuleType.Get()
		data.ModuleType = utils.UpdateReferenceAttribute(data.ModuleType, mt.GetModel(), "", mt.Id)
	} else {
		data.ModuleType = types.StringNull()
	}

	// Map type
	data.Type = types.StringValue(string(template.Type.GetValue()))

	// Map label
	data.Label = utils.StringFromAPI(template.HasLabel(), template.GetLabel, data.Label)

	// Map color
	data.Color = utils.StringFromAPI(template.HasColor(), template.GetColor, data.Color)

	// Map positions
	if positions, ok := template.GetPositionsOk(); ok && positions != nil {
		data.Positions = types.Int32Value(*positions)
	} else {
		data.Positions = types.Int32Value(1)
	}

	// Map description
	data.Description = utils.StringFromAPI(template.HasDescription(), template.GetDescription, data.Description)
}
