// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &PowerOutletTemplateResource{}

var _ resource.ResourceWithImportState = &PowerOutletTemplateResource{}

// NewPowerOutletTemplateResource returns a new resource implementing the power outlet template resource.

func NewPowerOutletTemplateResource() resource.Resource {
	return &PowerOutletTemplateResource{}
}

// PowerOutletTemplateResource defines the resource implementation.

type PowerOutletTemplateResource struct {
	client *netbox.APIClient
}

// PowerOutletTemplateResourceModel describes the resource data model.

type PowerOutletTemplateResourceModel struct {
	ID types.Int32 `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	ModuleType types.String `tfsdk:"module_type"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	PowerPort types.Int32 `tfsdk:"power_port"`

	FeedLeg types.String `tfsdk:"feed_leg"`

	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.

func (r *PowerOutletTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_outlet_template"
}

// Schema defines the schema for the resource.

func (r *PowerOutletTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a power outlet template in NetBox. Power outlet templates define the default power outlets that will be created on new devices of a specific device type or modules of a module type.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the power outlet template.",

				Computed: true,

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
				MarkdownDescription: "The name of the power outlet template. Use {module} as a substitution for the module bay position when attached to a module type.",

				Required: true,
			},

			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the power outlet template.",

				Optional: true,
			},

			"type": schema.StringAttribute{
				MarkdownDescription: "The type of power outlet (e.g., iec-60320-c5, iec-60320-c7, iec-60320-c13, iec-60320-c15, iec-60320-c19, nema-1-15r, nema-5-15r, nema-5-20r, etc.).",

				Optional: true,
			},

			"power_port": schema.Int32Attribute{
				MarkdownDescription: "The power port template that feeds this power outlet.",
				Optional:            true,
			},

			"feed_leg": schema.StringAttribute{
				MarkdownDescription: "Feed leg for three-phase power (A, B, or C).",
				Optional:            true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("power outlet template"))
}

// Configure adds the provider configured client to the resource.

func (r *PowerOutletTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PowerOutletTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PowerOutletTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request

	apiReq := netbox.NewWritablePowerOutletTemplateRequest(data.Name.ValueString())

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

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		apiReq.SetType(netbox.PatchedWritablePowerOutletTemplateRequestType(data.Type.ValueString()))
	}

	if !data.PowerPort.IsNull() && !data.PowerPort.IsUnknown() {
		// Lookup the power port template by ID - we need the name for BriefPowerPortTemplateRequest

		powerPortID := data.PowerPort.ValueInt32()

		powerPort, httpResp, err := r.client.DcimAPI.DcimPowerPortTemplatesRetrieve(ctx, powerPortID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError(

				"Error looking up power port template",

				utils.FormatAPIError(fmt.Sprintf("read power port template ID %d", powerPortID), err, httpResp),
			)

			return
		}

		apiReq.SetPowerPort(netbox.BriefPowerPortTemplateRequest{
			Name: powerPort.GetName(),
		})
	}

	if !data.FeedLeg.IsNull() && !data.FeedLeg.IsUnknown() {
		apiReq.SetFeedLeg(netbox.PatchedWritablePowerOutletRequestFeedLeg(data.FeedLeg.ValueString()))
	}

	// Apply description
	utils.ApplyDescription(apiReq, data.Description)

	tflog.Debug(ctx, "Creating power outlet template", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletTemplatesCreate(ctx).WritablePowerOutletTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating power outlet template",

			utils.FormatAPIError("create power outlet template", err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	tflog.Trace(ctx, "Created power outlet template", map[string]interface{}{
		"id": data.ID.ValueInt32(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.

func (r *PowerOutletTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PowerOutletTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Reading power outlet template", map[string]interface{}{
		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletTemplatesRetrieve(ctx, templateID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Power outlet template not found, removing from state", map[string]interface{}{
				"id": templateID,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading power outlet template",

			utils.FormatAPIError(fmt.Sprintf("read power outlet template ID %d", templateID), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.

func (r *PowerOutletTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PowerOutletTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	// Build the API request

	apiReq := netbox.NewWritablePowerOutletTemplateRequest(data.Name.ValueString())

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

	if utils.IsSet(data.Type) {
		apiReq.SetType(netbox.PatchedWritablePowerOutletTemplateRequestType(data.Type.ValueString()))
	} else if data.Type.IsNull() {
		// Explicitly clear type when removed from config
		apiReq.SetType("")
	}

	if !data.PowerPort.IsNull() && !data.PowerPort.IsUnknown() {
		// Lookup the power port template by ID - we need the name for BriefPowerPortTemplateRequest

		powerPortID := data.PowerPort.ValueInt32()

		powerPort, httpResp, err := r.client.DcimAPI.DcimPowerPortTemplatesRetrieve(ctx, powerPortID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError(

				"Error looking up power port template",

				utils.FormatAPIError(fmt.Sprintf("read power port template ID %d", powerPortID), err, httpResp),
			)

			return
		}

		apiReq.SetPowerPort(netbox.BriefPowerPortTemplateRequest{
			Name: powerPort.GetName(),
		})
	} else if data.PowerPort.IsNull() {
		// Explicitly clear power_port when removed from config
		apiReq.SetPowerPortNil()
	}

	if utils.IsSet(data.FeedLeg) {
		apiReq.SetFeedLeg(netbox.PatchedWritablePowerOutletRequestFeedLeg(data.FeedLeg.ValueString()))
	} else if data.FeedLeg.IsNull() {
		// Explicitly clear feed_leg when removed from config
		apiReq.SetFeedLeg("")
	}

	// Apply description
	utils.ApplyDescription(apiReq, data.Description)

	tflog.Debug(ctx, "Updating power outlet template", map[string]interface{}{
		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletTemplatesUpdate(ctx, templateID).WritablePowerOutletTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating power outlet template",

			utils.FormatAPIError(fmt.Sprintf("update power outlet template ID %d", templateID), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.

func (r *PowerOutletTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PowerOutletTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Deleting power outlet template", map[string]interface{}{
		"id": templateID,
	})

	httpResp, err := r.client.DcimAPI.DcimPowerOutletTemplatesDestroy(ctx, templateID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting power outlet template",

			utils.FormatAPIError(fmt.Sprintf("delete power outlet template ID %d", templateID), err, httpResp),
		)

		return
	}
}

// ImportState imports the resource state from Terraform.

func (r *PowerOutletTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *PowerOutletTemplateResource) mapResponseToModel(template *netbox.PowerOutletTemplate, data *PowerOutletTemplateResourceModel) {
	data.ID = types.Int32Value(template.GetId())

	data.Name = types.StringValue(template.GetName())

	// Map device type - preserve user's input format

	if template.DeviceType.IsSet() && template.DeviceType.Get() != nil {
		dt := template.DeviceType.Get()

		data.DeviceType = utils.UpdateReferenceAttribute(data.DeviceType, dt.GetModel(), dt.GetSlug(), dt.Id)
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

	// Map label

	// Map label - always set since it's computed

	data.Label = utils.StringFromAPI(template.HasLabel(), template.GetLabel, data.Label)

	// Map type

	if template.Type.IsSet() && template.Type.Get() != nil {
		data.Type = types.StringValue(string(template.Type.Get().GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map power_port

	if template.PowerPort.IsSet() && template.PowerPort.Get() != nil {
		data.PowerPort = types.Int32Value(template.PowerPort.Get().Id)
	} else {
		data.PowerPort = types.Int32Null()
	}

	// Map feed_leg

	if template.FeedLeg.IsSet() && template.FeedLeg.Get() != nil {
		data.FeedLeg = types.StringValue(string(template.FeedLeg.Get().GetValue()))
	} else {
		data.FeedLeg = types.StringNull()
	}

	// Map description

	data.Description = utils.StringFromAPI(template.HasDescription(), template.GetDescription, data.Description)
}
