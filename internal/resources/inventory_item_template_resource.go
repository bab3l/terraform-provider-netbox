// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource = &InventoryItemTemplateResource{}

	_ resource.ResourceWithConfigure = &InventoryItemTemplateResource{}

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
	ID types.String `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	Parent types.String `tfsdk:"parent"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Role types.String `tfsdk:"role"`

	Manufacturer types.String `tfsdk:"manufacturer"`

	PartID types.String `tfsdk:"part_id"`

	Description types.String `tfsdk:"description"`

	ComponentType types.String `tfsdk:"component_type"`

	ComponentID types.String `tfsdk:"component_id"`
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

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device_type": schema.StringAttribute{

				MarkdownDescription: "The device type this inventory item template belongs to (ID or model name).",

				Required: true,
			},

			"parent": schema.StringAttribute{

				MarkdownDescription: "Parent inventory item template (ID).",

				Optional: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the inventory item template. {module} is accepted as a substitution for the module bay position when attached to a module type.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the inventory item template.",

				Optional: true,
			},

			"role": schema.StringAttribute{

				MarkdownDescription: "The inventory item role (ID or slug).",

				Optional: true,
			},

			"manufacturer": schema.StringAttribute{

				MarkdownDescription: "The manufacturer of the inventory item (ID or slug).",

				Optional: true,
			},

			"part_id": schema.StringAttribute{

				MarkdownDescription: "Manufacturer-assigned part identifier.",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the inventory item template.",

				Optional: true,
			},

			"component_type": schema.StringAttribute{

				MarkdownDescription: "The type of component this inventory item represents (e.g., `dcim.interface`).",

				Optional: true,
			},

			"component_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the component this inventory item represents.",

				Optional: true,
			},
		},
	}

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

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

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

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

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

		"name": data.Name.ValueString(),

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

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

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

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

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

	}

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

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

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

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

	tflog.Debug(ctx, "Updating inventory item template", map[string]interface{}{

		"id": id,

		"name": data.Name.ValueString(),

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

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

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

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapToState maps the API response to the Terraform state.

func (r *InventoryItemTemplateResource) mapToState(ctx context.Context, result *netbox.InventoryItemTemplate, data *InventoryItemTemplateResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	// Map device type (required field)

	deviceType := result.GetDeviceType()

	data.DeviceType = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))

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

		data.Role = types.StringValue(fmt.Sprintf("%d", role.GetId()))

	} else {

		data.Role = types.StringNull()

	}

	// Map manufacturer

	if result.HasManufacturer() && result.GetManufacturer().Id != 0 {

		manufacturer := result.GetManufacturer()

		data.Manufacturer = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))

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
