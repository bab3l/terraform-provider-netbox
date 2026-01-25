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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &FrontPortResource{}
	_ resource.ResourceWithConfigure   = &FrontPortResource{}
	_ resource.ResourceWithImportState = &FrontPortResource{}
	_ resource.ResourceWithIdentity    = &FrontPortResource{}
)

// NewFrontPortResource returns a new resource implementing the front port resource.
func NewFrontPortResource() resource.Resource {
	return &FrontPortResource{}
}

// FrontPortResource defines the resource implementation.
type FrontPortResource struct {
	client *netbox.APIClient
}

// FrontPortResourceModel describes the resource data model.
type FrontPortResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Device           types.String `tfsdk:"device"`
	Name             types.String `tfsdk:"name"`
	Label            types.String `tfsdk:"label"`
	Type             types.String `tfsdk:"type"`
	Color            types.String `tfsdk:"color"`
	RearPort         types.String `tfsdk:"rear_port"`
	RearPortPosition types.Int32  `tfsdk:"rear_port_position"`
	Description      types.String `tfsdk:"description"`
	MarkConnected    types.Bool   `tfsdk:"mark_connected"`
	Tags             types.Set    `tfsdk:"tags"`
	CustomFields     types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *FrontPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_front_port"
}

// Schema defines the schema for the resource.
func (r *FrontPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a front port in NetBox. Front ports represent physical ports on the front of a device, typically used for patch panels and fiber distribution. They are mapped to rear ports.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the front port.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device",
				"The device this front port belongs to (ID or name).",
			),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the front port.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the front port.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of front port (e.g., `8p8c`, `8p6c`, `110-punch`, `bnc`, `f`, `n`, `mrj21`, `fc`, `lc`, `lc-pc`, `lc-upc`, `lc-apc`, `lsh`, `mpo`, `mtrj`, `sc`, `sc-pc`, `sc-upc`, `sc-apc`, `st`, `cs`, `sn`, `splice`, `other`).",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the front port in hex format (e.g., `aa1409`).",
				Optional:            true,
			},
			"rear_port": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"rear_port",
				"The rear port that this front port maps to (ID).",
			),
			"rear_port_position": schema.Int32Attribute{
				MarkdownDescription: "Position on the rear port (1-1024). Default is 1.",
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(1),
			},
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("front port"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *FrontPortResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *FrontPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *FrontPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FrontPortResourceModel
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

	// Parse rear port ID
	rearPortID, err := utils.ParseID(data.RearPort.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Rear Port ID",
			fmt.Sprintf("Could not parse rear port ID %q: %s", data.RearPort.ValueString(), err),
		)
		return
	}

	// Build request
	apiReq := netbox.NewWritableFrontPortRequest(*device, data.Name.ValueString(), netbox.FrontPortTypeValue(data.Type.ValueString()), rearPortID)

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		apiReq.SetColor(data.Color.ValueString())
	}

	if !data.RearPortPosition.IsNull() && !data.RearPortPosition.IsUnknown() {
		apiReq.SetRearPortPosition(data.RearPortPosition.ValueInt32())
	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating front port", map[string]interface{}{
		"device":    data.Device.ValueString(),
		"name":      data.Name.ValueString(),
		"rear_port": rearPortID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsCreate(ctx).WritableFrontPortRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating front port",
			utils.FormatAPIError("create front port", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "Created front port", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves the resource.
func (r *FrontPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FrontPortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading front port", map[string]interface{}{
		"id": portID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsRetrieve(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Front port not found, removing from state", map[string]interface{}{
				"id": portID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading front port",
			utils.FormatAPIError(fmt.Sprintf("read front port ID %d", portID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *FrontPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read both state and plan for merge-aware custom fields handling
	var state, plan FrontPortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse ID %q: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, plan.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse rear port ID
	rearPortID, err := utils.ParseID(plan.RearPort.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Rear Port ID",
			fmt.Sprintf("Could not parse rear port ID %q: %s", plan.RearPort.ValueString(), err),
		)
		return
	}

	// Build request
	apiReq := netbox.NewWritableFrontPortRequest(*device, plan.Name.ValueString(), netbox.FrontPortTypeValue(plan.Type.ValueString()), rearPortID)

	// Set optional fields
	utils.ApplyLabel(apiReq, plan.Label)

	// For nullable string fields, explicitly clear if null in plan
	if plan.Color.IsNull() {
		apiReq.SetColor("")
	} else if !plan.Color.IsUnknown() {
		apiReq.SetColor(plan.Color.ValueString())
	}

	if !plan.RearPortPosition.IsNull() && !plan.RearPortPosition.IsUnknown() {
		apiReq.SetRearPortPosition(plan.RearPortPosition.ValueInt32())
	}

	if !plan.MarkConnected.IsNull() && !plan.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(plan.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields with merge-aware helpers
	utils.ApplyDescription(apiReq, plan.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating front port", map[string]interface{}{
		"id": portID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsUpdate(ctx, portID).WritableFrontPortRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating front port",
			utils.FormatAPIError(fmt.Sprintf("update front port ID %d", portID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the resource.
func (r *FrontPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FrontPortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting front port", map[string]interface{}{
		"id": portID,
	})
	httpResp, err := r.client.DcimAPI.DcimFrontPortsDestroy(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting front port",
			utils.FormatAPIError(fmt.Sprintf("delete front port ID %d", portID), err, httpResp),
		)
		return
	}
}

// ImportState imports the resource.
func (r *FrontPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		portID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Import ID",
				fmt.Sprintf("Front Port ID must be a number, got: %s", parsed.ID),
			)
			return
		}
		response, httpResp, err := r.client.DcimAPI.DcimFrontPortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing front port",
				utils.FormatAPIError(fmt.Sprintf("import front port ID %d", portID), err, httpResp),
			)
			return
		}

		var data FrontPortResourceModel
		if device := response.GetDevice(); device.Id != 0 {
			data.Device = types.StringValue(device.GetName())
		}
		if rearPort := response.GetRearPort(); rearPort.Id != 0 {
			data.RearPort = types.StringValue(fmt.Sprintf("%d", rearPort.GetId()))
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, response.HasTags(), response.GetTags(), data.Tags)
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

		r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, response.GetCustomFields(), &resp.Diagnostics)
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

	portID, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Front Port ID must be a number, got: %s", req.ID),
		)
		return
	}
	response, httpResp, err := r.client.DcimAPI.DcimFrontPortsRetrieve(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing front port",
			utils.FormatAPIError(fmt.Sprintf("import front port ID %d", portID), err, httpResp),
		)
		return
	}

	var data FrontPortResourceModel
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *FrontPortResource) mapResponseToModel(ctx context.Context, port *netbox.FrontPort, data *FrontPortResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", port.GetId()))
	data.Name = types.StringValue(port.GetName())

	// Map device - preserve user's input format
	if device := port.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map type
	data.Type = types.StringValue(string(port.Type.GetValue()))

	// Map label
	if label, ok := port.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map color
	if color, ok := port.GetColorOk(); ok && color != nil && *color != "" {
		data.Color = types.StringValue(*color)
	} else {
		data.Color = types.StringNull()
	}

	// Map rear port - preserve user's input format
	if rearPort := port.GetRearPort(); rearPort.Id != 0 {
		data.RearPort = utils.UpdateReferenceAttribute(data.RearPort, rearPort.GetName(), "", rearPort.GetId())
	}

	// Map rear port position
	if rearPortPos, ok := port.GetRearPortPositionOk(); ok && rearPortPos != nil {
		data.RearPortPosition = types.Int32Value(*rearPortPos)
	} else {
		data.RearPortPosition = types.Int32Value(1)
	}

	// Map description
	if desc, ok := port.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if markConnected, ok := port.GetMarkConnectedOk(); ok && markConnected != nil {
		data.MarkConnected = types.BoolValue(*markConnected)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Handle tags with filter-to-owned pattern
	planTags := data.Tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, port.HasTags(), port.GetTags(), planTags)

	// Handle custom fields with filter-to-owned pattern
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, port.GetCustomFields(), diags)
}
