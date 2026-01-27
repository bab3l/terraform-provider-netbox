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
	_ resource.Resource                = &InventoryItemTemplateResource{}
	_ resource.ResourceWithConfigure   = &InventoryItemTemplateResource{}
	_ resource.ResourceWithImportState = &InventoryItemTemplateResource{}
)

// NewInventoryItemTemplateResource returns a new resource implementing the inventory item template resource.
func NewInventoryItemTemplateResource() resource.Resource {
	return &InventoryItemTemplateResource{}
}

// InventoryItemTemplateResource defines the resource implementation.
type InventoryItemTemplateResource struct {
	client *netbox.APIClient
}

// InventoryItemTemplateResourceModel describes the resource data model.
type InventoryItemTemplateResourceModel struct {
	ID            types.String `tfsdk:"id"`
	DeviceType    types.String `tfsdk:"device_type"`
	Parent        types.String `tfsdk:"parent"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Role          types.String `tfsdk:"role"`
	Manufacturer  types.String `tfsdk:"manufacturer"`
	PartID        types.String `tfsdk:"part_id"`
	Description   types.String `tfsdk:"description"`
	ComponentType types.String `tfsdk:"component_type"`
	ComponentID   types.String `tfsdk:"component_id"`
}

// Metadata returns the resource type name.
func (r *InventoryItemTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_item_template"
}

// Schema defines the schema for the resource.
func (r *InventoryItemTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an inventory item template in NetBox. Inventory item templates define inventory items for device types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the inventory item template.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_type": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device_type",
				"The device type this inventory item template belongs to (ID or model name).",
			),
			"parent": schema.StringAttribute{
				MarkdownDescription: "Parent inventory item template (ID).",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the inventory item template. {module} is accepted as a substitution for the module bay position when attached to a module type.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the inventory item template.",
				Optional:            true,
			},
			"role": nbschema.ReferenceAttributeWithDiffSuppress(
				"inventory item role",
				"The inventory item role (ID or slug).",
			),
			"manufacturer": nbschema.ReferenceAttributeWithDiffSuppress(
				"manufacturer",
				"The manufacturer of the inventory item (ID or slug).",
			),
			"part_id": schema.StringAttribute{
				MarkdownDescription: "Manufacturer-assigned part identifier.",
				Optional:            true,
			},
			"component_type": schema.StringAttribute{
				MarkdownDescription: "The type of component this inventory item represents (e.g., `dcim.interface`).",
				Optional:            true,
			},
			"component_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the component this inventory item represents.",
				Optional:            true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("inventory item template"))
}

// Configure adds the provider configured client to the resource.
func (r *InventoryItemTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *InventoryItemTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InventoryItemTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup device type
	deviceType, diags := lookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewInventoryItemTemplateRequest(*deviceType, data.Name.ValueString())

	// Set optional fields
	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		var parentID int32
		_, err := fmt.Sscanf(data.Parent.ValueString(), "%d", &parentID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Could not parse parent ID '%s': %s", data.Parent.ValueString(), err.Error()),
			)
			return
		}
		apiReq.SetParent(parentID)
	}

	utils.ApplyLabel(apiReq, data.Label)

	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		role, roleDiags := lookup.LookupInventoryItemRole(ctx, r.client, data.Role.ValueString())
		resp.Diagnostics.Append(roleDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetRole(*role)
	}

	if !data.Manufacturer.IsNull() && !data.Manufacturer.IsUnknown() {
		manufacturer, mfrDiags := lookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(mfrDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetManufacturer(*manufacturer)
	}

	if !data.PartID.IsNull() && !data.PartID.IsUnknown() {
		apiReq.SetPartId(data.PartID.ValueString())
	}

	// Apply description
	utils.ApplyDescription(apiReq, data.Description)

	if !data.ComponentType.IsNull() && !data.ComponentType.IsUnknown() {
		apiReq.SetComponentType(data.ComponentType.ValueString())
	}

	if !data.ComponentID.IsNull() && !data.ComponentID.IsUnknown() {
		var componentID int64
		_, err := fmt.Sscanf(data.ComponentID.ValueString(), "%d", &componentID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Component ID",
				fmt.Sprintf("Could not parse component ID '%s': %s", data.ComponentID.ValueString(), err.Error()),
			)
			return
		}
		apiReq.SetComponentId(componentID)
	}
	tflog.Debug(ctx, "Creating inventory item template", map[string]interface{}{
		"name":        data.Name.ValueString(),
		"device_type": data.DeviceType.ValueString(),
	})

	// Create the resource
	result, httpResp, err := r.client.DcimAPI.DcimInventoryItemTemplatesCreate(ctx).InventoryItemTemplateRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating inventory item template",
			utils.FormatAPIError("create inventory item template", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapToState(result, &data)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the resource.
func (r *InventoryItemTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InventoryItemTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing inventory item template ID",
			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Read from API
	result, httpResp, err := r.client.DcimAPI.DcimInventoryItemTemplatesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading inventory item template",
			utils.FormatAPIError(fmt.Sprintf("read inventory item template ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapToState(result, &data)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *InventoryItemTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InventoryItemTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get prior state to check what fields were previously set
	var state InventoryItemTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing inventory item template ID",
			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Lookup device type
	deviceType, diags := lookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewInventoryItemTemplateRequest(*deviceType, data.Name.ValueString())

	// Set optional fields
	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		var parentID int32
		_, err := fmt.Sscanf(data.Parent.ValueString(), "%d", &parentID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Could not parse parent ID '%s': %s", data.Parent.ValueString(), err.Error()),
			)
			return
		}
		apiReq.SetParent(parentID)
	} else if data.Parent.IsNull() && !state.Parent.IsNull() {
		// Only send null if we're removing a previously-set value
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["parent"] = nil
	}
	utils.ApplyLabel(apiReq, data.Label)

	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		role, roleDiags := lookup.LookupInventoryItemRole(ctx, r.client, data.Role.ValueString())
		resp.Diagnostics.Append(roleDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetRole(*role)
	} else if data.Role.IsNull() && !state.Role.IsNull() {
		// Only send null if we're removing a previously-set value
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["role"] = nil
	}

	if !data.Manufacturer.IsNull() && !data.Manufacturer.IsUnknown() {
		manufacturer, mfrDiags := lookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(mfrDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetManufacturer(*manufacturer)
	} else if data.Manufacturer.IsNull() && !state.Manufacturer.IsNull() {
		// Only send null if we're removing a previously-set value
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["manufacturer"] = nil
	}

	if !data.PartID.IsNull() && !data.PartID.IsUnknown() {
		apiReq.SetPartId(data.PartID.ValueString())
	} else if data.PartID.IsNull() && !state.PartID.IsNull() {
		// Only send null if we're removing a previously-set value
		// NOTE: NetBox API may reject this with "This field may not be null"
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["part_id"] = nil
	}

	// Apply description
	utils.ApplyDescription(apiReq, data.Description)
	if !data.ComponentType.IsNull() && !data.ComponentType.IsUnknown() {
		apiReq.SetComponentType(data.ComponentType.ValueString())
	} else if data.ComponentType.IsNull() && !state.ComponentType.IsNull() {
		// Only send null if we're removing a previously-set value
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["component_type"] = nil
	}

	if !data.ComponentID.IsNull() && !data.ComponentID.IsUnknown() {
		var componentID int64
		_, err := fmt.Sscanf(data.ComponentID.ValueString(), "%d", &componentID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Component ID",
				fmt.Sprintf("Could not parse component ID '%s': %s", data.ComponentID.ValueString(), err.Error()),
			)
			return
		}
		apiReq.SetComponentId(componentID)
	} else if data.ComponentID.IsNull() && !state.ComponentID.IsNull() {
		// Only send null if we're removing a previously-set value
		if apiReq.AdditionalProperties == nil {
			apiReq.AdditionalProperties = make(map[string]interface{})
		}
		apiReq.AdditionalProperties["component_id"] = nil
	}
	tflog.Debug(ctx, "Updating inventory item template", map[string]interface{}{
		"id":          id,
		"name":        data.Name.ValueString(),
		"device_type": data.DeviceType.ValueString(),
	})

	// Update the resource
	result, httpResp, err := r.client.DcimAPI.DcimInventoryItemTemplatesUpdate(ctx, id).InventoryItemTemplateRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating inventory item template",
			utils.FormatAPIError(fmt.Sprintf("update inventory item template ID %d", id), err, httpResp),
		)
		return
	}

	// Preserve display_name since it's computed but might change when other attributes update
	// Map response to state
	r.mapToState(result, &data)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.
func (r *InventoryItemTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InventoryItemTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing inventory item template ID",
			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting inventory item template", map[string]interface{}{"id": id})

	// Delete the resource
	httpResp, err := r.client.DcimAPI.DcimInventoryItemTemplatesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting inventory item template",
			utils.FormatAPIError(fmt.Sprintf("delete inventory item template ID %d", id), err, httpResp),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *InventoryItemTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapToState maps the API response to the Terraform state.
func (r *InventoryItemTemplateResource) mapToState(result *netbox.InventoryItemTemplate, data *InventoryItemTemplateResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())

	// Map device type (required field)
	deviceType := result.GetDeviceType()
	data.DeviceType = utils.UpdateReferenceAttribute(data.DeviceType, deviceType.GetSlug(), deviceType.GetModel(), deviceType.GetId())

	// Map parent (NullableInt32 - just the ID)
	// Check if the value is set and non-nil (parent: null in JSON means no parent)
	if parentID, ok := result.GetParentOk(); ok && parentID != nil {
		data.Parent = types.StringValue(fmt.Sprintf("%d", *parentID))
	} else {
		data.Parent = types.StringNull()
	}

	// Map label
	if result.HasLabel() && result.GetLabel() != "" {
		data.Label = types.StringValue(result.GetLabel())
	} else {
		data.Label = types.StringNull()
	}

	// Map role
	if result.HasRole() && result.GetRole().Id != 0 {
		role := result.GetRole()
		data.Role = utils.UpdateReferenceAttribute(data.Role, role.GetName(), role.GetSlug(), role.GetId())
	} else {
		data.Role = types.StringNull()
	}

	// Map manufacturer
	if result.HasManufacturer() && result.GetManufacturer().Id != 0 {
		manufacturer := result.GetManufacturer()
		data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, manufacturer.GetName(), manufacturer.GetSlug(), manufacturer.GetId())
	} else {
		data.Manufacturer = types.StringNull()
	}

	// Map part ID
	if result.HasPartId() && result.GetPartId() != "" {
		data.PartID = types.StringValue(result.GetPartId())
	} else {
		data.PartID = types.StringNull()
	}

	// Map description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map component type
	if result.HasComponentType() {
		componentTypePtr, ok := result.GetComponentTypeOk()

		if ok && componentTypePtr != nil && *componentTypePtr != "" {
			data.ComponentType = types.StringValue(*componentTypePtr)
		} else {
			data.ComponentType = types.StringNull()
		}
	} else {
		data.ComponentType = types.StringNull()
	}

	// Map component ID
	if result.HasComponentId() {
		componentIDPtr, ok := result.GetComponentIdOk()

		if ok && componentIDPtr != nil {
			data.ComponentID = types.StringValue(fmt.Sprintf("%d", *componentIDPtr))
		} else {
			data.ComponentID = types.StringNull()
		}
	} else {
		data.ComponentID = types.StringNull()
	}
}
