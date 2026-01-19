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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &PowerOutletResource{}
	_ resource.ResourceWithConfigure   = &PowerOutletResource{}
	_ resource.ResourceWithImportState = &PowerOutletResource{}
	_ resource.ResourceWithIdentity    = &PowerOutletResource{}
)

// NewPowerOutletResource returns a new resource implementing the power outlet resource.
func NewPowerOutletResource() resource.Resource {
	return &PowerOutletResource{}
}

// PowerOutletResource defines the resource implementation.
type PowerOutletResource struct {
	client *netbox.APIClient
}

// PowerOutletResourceModel describes the resource data model.
type PowerOutletResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	PowerPort     types.Int32  `tfsdk:"power_port"`
	FeedLeg       types.String `tfsdk:"feed_leg"`
	Description   types.String `tfsdk:"description"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *PowerOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_outlet"
}

// Schema defines the schema for the resource.
func (r *PowerOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a power outlet in NetBox. Power outlets represent power distribution connections on PDUs and other power distribution devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the power outlet.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device",
				"The device this power outlet belongs to (ID or name).",
			),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the power outlet.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the power outlet.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Power outlet type (e.g., `iec-60320-c13`, `nema-5-15r`, etc.).",
				Optional:            true,
			},
			"power_port": schema.Int32Attribute{
				MarkdownDescription: "The power port ID that feeds this outlet.",
				Optional:            true,
			},
			"feed_leg": schema.StringAttribute{
				MarkdownDescription: "Phase leg for three-phase power. Valid values: `A`, `B`, `C`.",
				Optional:            true,
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
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("power outlet"))

	// Add tags and custom_fields attributes
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *PowerOutletResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *PowerOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PowerOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PowerOutletResourceModel
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
	apiReq := netbox.NewWritablePowerOutletRequest(*device, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		outletType := netbox.PatchedWritablePowerOutletRequestType(data.Type.ValueString())
		apiReq.SetType(outletType)
	}

	if !data.PowerPort.IsNull() && !data.PowerPort.IsUnknown() {
		// Look up the power port to get its device and name
		powerPortID := fmt.Sprintf("%d", data.PowerPort.ValueInt32())
		powerPortReq, ppDiags := lookup.LookupPowerPort(ctx, r.client, powerPortID)
		resp.Diagnostics.Append(ppDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetPowerPort(*powerPortReq)
	}

	if !data.FeedLeg.IsNull() && !data.FeedLeg.IsUnknown() {
		feedLeg := netbox.PatchedWritablePowerOutletRequestFeedLeg(data.FeedLeg.ValueString())
		apiReq.SetFeedLeg(feedLeg)
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
	tflog.Debug(ctx, "Creating power outlet", map[string]interface{}{
		"device": data.Device.ValueString(),
		"name":   data.Name.ValueString(),
	})
	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsCreate(ctx).WritablePowerOutletRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating power outlet",
			utils.FormatAPIError(fmt.Sprintf("create power outlet %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Created power outlet", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *PowerOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PowerOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outletID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Power Outlet ID",
			fmt.Sprintf("Power Outlet ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading power outlet", map[string]interface{}{
		"id": outletID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsRetrieve(ctx, outletID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading power outlet",
			utils.FormatAPIError(fmt.Sprintf("read power outlet ID %d", outletID), err, httpResp),
		)
		return
	}

	// Preserve original custom_fields state
	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore original custom_fields if it was null/empty and API returned none
	if !utils.IsSet(originalCustomFields) && !utils.IsSet(data.CustomFields) {
		data.CustomFields = originalCustomFields
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *PowerOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data PowerOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outletID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Power Outlet ID",
			fmt.Sprintf("Power Outlet ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritablePowerOutletRequest(*device, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		outletType := netbox.PatchedWritablePowerOutletRequestType(data.Type.ValueString())
		apiReq.SetType(outletType)
	} else if data.Type.IsNull() {
		// Clear the type by sending empty string
		apiReq.SetType("")
	}

	if !data.PowerPort.IsNull() && !data.PowerPort.IsUnknown() {
		// Look up the power port to get its device and name
		powerPortID := fmt.Sprintf("%d", data.PowerPort.ValueInt32())
		powerPortReq, ppDiags := lookup.LookupPowerPort(ctx, r.client, powerPortID)
		resp.Diagnostics.Append(ppDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetPowerPort(*powerPortReq)
	} else if data.PowerPort.IsNull() {
		// Clear the power port
		apiReq.SetPowerPortNil()
	}

	if !data.FeedLeg.IsNull() && !data.FeedLeg.IsUnknown() {
		feedLeg := netbox.PatchedWritablePowerOutletRequestFeedLeg(data.FeedLeg.ValueString())
		apiReq.SetFeedLeg(feedLeg)
	} else if data.FeedLeg.IsNull() {
		// Clear the feed leg by sending empty string
		apiReq.SetFeedLeg("")
	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields with merge-aware behavior
	utils.ApplyDescription(apiReq, data.Description)

	// Apply tags using plan tags (slug list format)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)

	// Apply custom fields with merge logic (preserves unmanaged fields)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, data.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating power outlet", map[string]interface{}{
		"id": outletID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsUpdate(ctx, outletID).WritablePowerOutletRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating power outlet",
			utils.FormatAPIError(fmt.Sprintf("update power outlet ID %d", outletID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.
func (r *PowerOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PowerOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outletID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Power Outlet ID",
			fmt.Sprintf("Power Outlet ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting power outlet", map[string]interface{}{
		"id": outletID,
	})
	httpResp, err := r.client.DcimAPI.DcimPowerOutletsDestroy(ctx, outletID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting power outlet",
			utils.FormatAPIError(fmt.Sprintf("delete power outlet ID %d", outletID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *PowerOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		outletID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Power Outlet ID must be a number, got: %s", parsed.ID))
			return
		}
		response, httpResp, err := r.client.DcimAPI.DcimPowerOutletsRetrieve(ctx, outletID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing power outlet", utils.FormatAPIError(fmt.Sprintf("import power outlet ID %d", outletID), err, httpResp))
			return
		}

		var data PowerOutletResourceModel
		if response.HasTags() {
			var tagSlugs []string
			for _, tag := range response.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		} else {
			data.Tags = types.SetNull(types.StringType)
		}
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *PowerOutletResource) mapResponseToModel(ctx context.Context, powerOutlet *netbox.PowerOutlet, data *PowerOutletResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", powerOutlet.GetId()))
	data.Name = types.StringValue(powerOutlet.GetName())

	// Map device - preserve user's input format
	if device := powerOutlet.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map label
	if label, ok := powerOutlet.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map type
	if powerOutlet.Type.IsSet() && powerOutlet.Type.Get() != nil {
		data.Type = types.StringValue(string(powerOutlet.Type.Get().GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map power_port
	if powerOutlet.PowerPort.IsSet() && powerOutlet.PowerPort.Get() != nil {
		data.PowerPort = types.Int32Value(powerOutlet.PowerPort.Get().Id)
	} else {
		data.PowerPort = types.Int32Null()
	}

	// Map feed_leg
	if powerOutlet.FeedLeg.IsSet() && powerOutlet.FeedLeg.Get() != nil {
		data.FeedLeg = types.StringValue(string(powerOutlet.FeedLeg.Get().GetValue()))
	} else {
		data.FeedLeg = types.StringNull()
	}

	// Map description
	if desc, ok := powerOutlet.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if mc, ok := powerOutlet.GetMarkConnectedOk(); ok && mc != nil {
		data.MarkConnected = types.BoolValue(*mc)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Handle tags - filter to owned slugs
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	case powerOutlet.HasTags():
		// Extract slugs from API tags
		var tagSlugs []string
		for _, tag := range powerOutlet.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Handle custom fields - use filtered-to-owned for partial management
	if powerOutlet.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, powerOutlet.GetCustomFields(), diags)
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
