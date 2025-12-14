// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &RackResource{}

var _ resource.ResourceWithImportState = &RackResource{}

func NewRackResource() resource.Resource {

	return &RackResource{}

}

// RackResource defines the resource implementation.

type RackResource struct {
	client *netbox.APIClient
}

// RackResourceModel describes the resource data model.

type RackResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Site types.String `tfsdk:"site"`

	Location types.String `tfsdk:"location"`

	Tenant types.String `tfsdk:"tenant"`

	Status types.String `tfsdk:"status"`

	Role types.String `tfsdk:"role"`

	Serial types.String `tfsdk:"serial"`

	AssetTag types.String `tfsdk:"asset_tag"`

	RackType types.String `tfsdk:"rack_type"`

	FormFactor types.String `tfsdk:"form_factor"`

	Width types.String `tfsdk:"width"`

	UHeight types.String `tfsdk:"u_height"`

	StartingUnit types.String `tfsdk:"starting_unit"`

	Weight types.String `tfsdk:"weight"`

	MaxWeight types.String `tfsdk:"max_weight"`

	WeightUnit types.String `tfsdk:"weight_unit"`

	DescUnits types.Bool `tfsdk:"desc_units"`

	OuterWidth types.String `tfsdk:"outer_width"`

	OuterDepth types.String `tfsdk:"outer_depth"`

	OuterUnit types.String `tfsdk:"outer_unit"`

	MountingDepth types.String `tfsdk:"mounting_depth"`

	Airflow types.String `tfsdk:"airflow"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *RackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_rack"

}

func (r *RackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a rack in Netbox. Racks represent physical equipment enclosures used to organize network infrastructure within a site or location.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("rack"),

			"name": nbschema.NameAttribute("rack", 100),

			"site": nbschema.RequiredReferenceAttribute("site", "ID or slug of the site where this rack is located. Required."),

			"location": nbschema.IDOnlyReferenceAttribute("location", "ID of the location within the site (e.g., building, floor, room)."),

			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this rack."),

			"status": nbschema.StatusAttribute([]string{"reserved", "available", "planned", "active", "deprecated"}, "Operational status of the rack. Defaults to `active`."),

			"role": nbschema.ReferenceAttribute("rack role", "ID or slug of the functional role of the rack."),

			"serial": nbschema.SerialAttribute(),

			"asset_tag": nbschema.AssetTagAttribute(),

			"rack_type": nbschema.IDOnlyReferenceAttribute("rack type", "ID of the rack type (model/form factor definition)."),

			"form_factor": schema.StringAttribute{

				MarkdownDescription: "Physical form factor of the rack. Valid values: `2-post-frame`, `4-post-frame`, `4-post-cabinet`, `wall-frame`, `wall-frame-vertical`, `wall-cabinet`, `wall-cabinet-vertical`.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.OneOf("2-post-frame", "4-post-frame", "4-post-cabinet", "wall-frame", "wall-frame-vertical", "wall-cabinet", "wall-cabinet-vertical"),
				},
			},

			"width": schema.StringAttribute{

				MarkdownDescription: "Rail-to-rail width of the rack in inches. Valid values: `10`, `19`, `21`, `23`. Defaults to 19.",

				Optional: true,

				Computed: true,

				Validators: []validator.String{

					stringvalidator.OneOf("10", "19", "21", "23"),
				},
			},

			"u_height": schema.StringAttribute{

				MarkdownDescription: "Height of the rack in rack units. Defaults to 42.",

				Optional: true,

				Computed: true,

				Validators: []validator.String{

					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},

			"starting_unit": schema.StringAttribute{

				MarkdownDescription: "Starting unit number for the rack (bottom). Defaults to 1.",

				Optional: true,

				Computed: true,

				Validators: []validator.String{

					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},

			"weight": schema.StringAttribute{

				MarkdownDescription: "Weight of the rack itself (numeric value).",

				Optional: true,
			},

			"max_weight": schema.StringAttribute{

				MarkdownDescription: "Maximum weight capacity of the rack.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},

			"weight_unit": schema.StringAttribute{

				MarkdownDescription: "Unit of measurement for weight. Valid values: `kg`, `g`, `lb`, `oz`.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.OneOf("kg", "g", "lb", "oz"),
				},
			},

			"desc_units": schema.BoolAttribute{

				MarkdownDescription: "If true, rack units are numbered in descending order (top to bottom).",

				Optional: true,

				Computed: true,
			},

			"outer_width": schema.StringAttribute{

				MarkdownDescription: "Outer width of the rack.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},

			"outer_depth": schema.StringAttribute{

				MarkdownDescription: "Outer depth of the rack.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},

			"outer_unit": schema.StringAttribute{

				MarkdownDescription: "Unit of measurement for outer dimensions. Valid values: `mm`, `in`.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.OneOf("mm", "in"),
				},
			},

			"mounting_depth": schema.StringAttribute{

				MarkdownDescription: "Maximum depth of equipment that can be installed (in mm).",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer",
					),
				},
			},

			"airflow": schema.StringAttribute{

				MarkdownDescription: "Direction of airflow through the rack. Valid values: `front-to-rear`, `rear-to-front`, `passive`, `mixed`.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.OneOf("front-to-rear", "rear-to-front", "passive", "mixed"),
				},
			},

			"description": nbschema.DescriptionAttribute("rack"),

			"comments": nbschema.CommentsAttribute("rack"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *RackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

// buildRackRequest creates a WritableRackRequest from the model.

func (r *RackResource) buildRackRequest(ctx context.Context, data *RackResourceModel, resp *resource.CreateResponse) *netbox.WritableRackRequest {

	// Lookup required site

	siteRef, diags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return nil

	}

	// Create the request

	rackRequest := netbox.WritableRackRequest{

		Name: data.Name.ValueString(),

		Site: *siteRef,
	}

	// Handle location relationship

	if !data.Location.IsNull() {

		locationRef, diags := netboxlookup.LookupLocation(ctx, r.client, data.Location.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Location = *netbox.NewNullableBriefLocationRequest(locationRef)

	}

	// Handle tenant relationship

	if !data.Tenant.IsNull() {

		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)

	}

	// Handle role relationship

	if !data.Role.IsNull() {

		roleRef, diags := netboxlookup.LookupRackRole(ctx, r.client, data.Role.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Role = *netbox.NewNullableBriefRackRoleRequest(roleRef)

	}

	// Handle rack_type relationship

	if !data.RackType.IsNull() {

		rackTypeRef, diags := netboxlookup.LookupRackType(ctx, r.client, data.RackType.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.RackType = *netbox.NewNullableBriefRackTypeRequest(rackTypeRef)

	}

	// Set status (default to "active" if not specified - required by API)

	if !data.Status.IsNull() && !data.Status.IsUnknown() {

		statusValue := netbox.PatchedWritableRackRequestStatus(data.Status.ValueString())

		rackRequest.Status = &statusValue

	} else {

		statusValue := netbox.PatchedWritableRackRequestStatus("active")

		rackRequest.Status = &statusValue

	}

	// Set form_factor

	if !data.FormFactor.IsNull() {

		formFactorValue := netbox.PatchedWritableRackRequestFormFactor(data.FormFactor.ValueString())

		rackRequest.FormFactor = &formFactorValue

	}

	// Set width (integer value: 10, 19, 21, 23)

	if !data.Width.IsNull() {

		var widthInt int32

		if _, err := fmt.Sscanf(data.Width.ValueString(), "%d", &widthInt); err == nil {

			widthValue, err := netbox.NewPatchedWritableRackRequestWidthFromValue(widthInt)

			if err == nil {

				rackRequest.Width = widthValue

			}

		}

	}

	// Set u_height

	if !data.UHeight.IsNull() {

		var uHeight int32

		if _, err := fmt.Sscanf(data.UHeight.ValueString(), "%d", &uHeight); err == nil {

			rackRequest.UHeight = &uHeight

		}

	}

	// Set starting_unit

	if !data.StartingUnit.IsNull() {

		var startingUnit int32

		if _, err := fmt.Sscanf(data.StartingUnit.ValueString(), "%d", &startingUnit); err == nil {

			rackRequest.StartingUnit = &startingUnit

		}

	}

	// Set weight

	if !data.Weight.IsNull() {

		var weight float64

		if _, err := fmt.Sscanf(data.Weight.ValueString(), "%f", &weight); err == nil {

			rackRequest.Weight = *netbox.NewNullableFloat64(&weight)

		}

	}

	// Set max_weight

	if !data.MaxWeight.IsNull() {

		var maxWeight int32

		if _, err := fmt.Sscanf(data.MaxWeight.ValueString(), "%d", &maxWeight); err == nil {

			rackRequest.MaxWeight = *netbox.NewNullableInt32(&maxWeight)

		}

	}

	// Set weight_unit

	if !data.WeightUnit.IsNull() {

		weightUnitValue := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())

		rackRequest.WeightUnit = &weightUnitValue

	}

	// Set desc_units

	if !data.DescUnits.IsNull() {

		descUnits := data.DescUnits.ValueBool()

		rackRequest.DescUnits = &descUnits

	}

	// Set outer_width

	if !data.OuterWidth.IsNull() {

		var outerWidth int32

		if _, err := fmt.Sscanf(data.OuterWidth.ValueString(), "%d", &outerWidth); err == nil {

			rackRequest.OuterWidth = *netbox.NewNullableInt32(&outerWidth)

		}

	}

	// Set outer_depth

	if !data.OuterDepth.IsNull() {

		var outerDepth int32

		if _, err := fmt.Sscanf(data.OuterDepth.ValueString(), "%d", &outerDepth); err == nil {

			rackRequest.OuterDepth = *netbox.NewNullableInt32(&outerDepth)

		}

	}

	// Set outer_unit

	if !data.OuterUnit.IsNull() {

		outerUnitValue := netbox.PatchedWritableRackRequestOuterUnit(data.OuterUnit.ValueString())

		rackRequest.OuterUnit = &outerUnitValue

	}

	// Set mounting_depth

	if !data.MountingDepth.IsNull() {

		var mountingDepth int32

		if _, err := fmt.Sscanf(data.MountingDepth.ValueString(), "%d", &mountingDepth); err == nil {

			rackRequest.MountingDepth = *netbox.NewNullableInt32(&mountingDepth)

		}

	}

	// Set airflow

	if !data.Airflow.IsNull() {

		airflowValue := netbox.PatchedWritableRackRequestAirflow(data.Airflow.ValueString())

		rackRequest.Airflow = &airflowValue

	}

	// Set serial

	if !data.Serial.IsNull() {

		serial := data.Serial.ValueString()

		rackRequest.Serial = &serial

	}

	// Set asset_tag

	if !data.AssetTag.IsNull() {

		assetTag := data.AssetTag.ValueString()

		rackRequest.AssetTag = *netbox.NewNullableString(&assetTag)

	}

	// Set description

	if !data.Description.IsNull() {

		description := data.Description.ValueString()

		rackRequest.Description = &description

	}

	// Set comments

	if !data.Comments.IsNull() {

		comments := data.Comments.ValueString()

		rackRequest.Comments = &comments

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	return &rackRequest

}

func (r *RackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data RackResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating rack", map[string]interface{}{

		"name": data.Name.ValueString(),

		"site": data.Site.ValueString(),
	})

	rackRequest := r.buildRackRequest(ctx, &data, resp)

	if rackRequest == nil {

		return

	}

	// Create the rack via API

	rack, httpResp, err := r.client.DcimAPI.DcimRacksCreate(ctx).WritableRackRequest(*rackRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		handler := utils.CreateErrorHandler{

			ResourceType: "netbox_rack",

			ResourceName: "this.rack",

			SlugValue: data.Name.ValueString(),

			LookupFunc: func(lookupCtx context.Context, name string) (string, error) {

				list, _, lookupErr := r.client.DcimAPI.DcimRacksList(lookupCtx).Name([]string{name}).Execute()

				if lookupErr != nil {

					return "", lookupErr

				}

				if list != nil && len(list.Results) > 0 {

					return fmt.Sprintf("%d", list.Results[0].GetId()), nil

				}

				return "", nil

			},
		}

		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)

		return

	}

	if httpResp.StatusCode != 201 {

		resp.Diagnostics.AddError(

			"Error creating rack",

			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)

		return

	}

	// Map response to state

	mapRackToState(ctx, rack, &data)

	tflog.Debug(ctx, "Created rack", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data RackResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Reading rack", map[string]interface{}{

		"id": data.ID.ValueString(),
	})

	// Parse ID

	rackID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Error parsing rack ID",

			fmt.Sprintf("Could not parse rack ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	rack, httpResp, err := r.client.DcimAPI.DcimRacksRetrieve(ctx, rackID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Rack not found, removing from state", map[string]interface{}{

				"id": data.ID.ValueString(),
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading rack",

			utils.FormatAPIError("read rack", err, httpResp),
		)

		return

	}

	// Preserve state custom fields for proper mapping

	stateCustomFields := data.CustomFields

	// Map response to state

	mapRackToState(ctx, rack, &data)

	// Restore custom fields from state (needed for proper type handling)

	if !stateCustomFields.IsNull() && !stateCustomFields.IsUnknown() {

		var stateCF []utils.CustomFieldModel

		if diags := stateCustomFields.ElementsAs(ctx, &stateCF, false); !diags.HasError() && len(stateCF) > 0 {

			if cfMap := rack.GetCustomFields(); cfMap != nil {

				mappedCF := utils.MapToCustomFieldModels(cfMap, stateCF)

				if len(mappedCF) > 0 {

					data.CustomFields, _ = types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, mappedCF)

				}

			}

		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// buildRackRequestForUpdate creates a WritableRackRequest for update operations.

func (r *RackResource) buildRackRequestForUpdate(ctx context.Context, data *RackResourceModel, resp *resource.UpdateResponse) *netbox.WritableRackRequest {

	// Lookup required site

	siteRef, diags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return nil

	}

	// Create the request

	rackRequest := netbox.WritableRackRequest{

		Name: data.Name.ValueString(),

		Site: *siteRef,
	}

	// Handle location relationship

	if !data.Location.IsNull() {

		locationRef, diags := netboxlookup.LookupLocation(ctx, r.client, data.Location.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Location = *netbox.NewNullableBriefLocationRequest(locationRef)

	}

	// Handle tenant relationship

	if !data.Tenant.IsNull() {

		tenantRef, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)

	}

	// Handle role relationship

	if !data.Role.IsNull() {

		roleRef, diags := netboxlookup.LookupRackRole(ctx, r.client, data.Role.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Role = *netbox.NewNullableBriefRackRoleRequest(roleRef)

	}

	// Handle rack_type relationship

	if !data.RackType.IsNull() {

		rackTypeRef, diags := netboxlookup.LookupRackType(ctx, r.client, data.RackType.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.RackType = *netbox.NewNullableBriefRackTypeRequest(rackTypeRef)

	}

	// Set status (default to "active" if not specified - required by API)

	if !data.Status.IsNull() && !data.Status.IsUnknown() {

		statusValue := netbox.PatchedWritableRackRequestStatus(data.Status.ValueString())

		rackRequest.Status = &statusValue

	} else {

		statusValue := netbox.PatchedWritableRackRequestStatus("active")

		rackRequest.Status = &statusValue

	}

	// Set form_factor

	if !data.FormFactor.IsNull() {

		formFactorValue := netbox.PatchedWritableRackRequestFormFactor(data.FormFactor.ValueString())

		rackRequest.FormFactor = &formFactorValue

	}

	// Set width (integer value: 10, 19, 21, 23)

	if !data.Width.IsNull() {

		var widthInt int32

		if _, err := fmt.Sscanf(data.Width.ValueString(), "%d", &widthInt); err == nil {

			widthValue, err := netbox.NewPatchedWritableRackRequestWidthFromValue(widthInt)

			if err == nil {

				rackRequest.Width = widthValue

			}

		}

	}

	// Set u_height

	if !data.UHeight.IsNull() {

		var uHeight int32

		if _, err := fmt.Sscanf(data.UHeight.ValueString(), "%d", &uHeight); err == nil {

			rackRequest.UHeight = &uHeight

		}

	}

	// Set starting_unit

	if !data.StartingUnit.IsNull() {

		var startingUnit int32

		if _, err := fmt.Sscanf(data.StartingUnit.ValueString(), "%d", &startingUnit); err == nil {

			rackRequest.StartingUnit = &startingUnit

		}

	}

	// Set weight

	if !data.Weight.IsNull() {

		var weight float64

		if _, err := fmt.Sscanf(data.Weight.ValueString(), "%f", &weight); err == nil {

			rackRequest.Weight = *netbox.NewNullableFloat64(&weight)

		}

	}

	// Set max_weight

	if !data.MaxWeight.IsNull() {

		var maxWeight int32

		if _, err := fmt.Sscanf(data.MaxWeight.ValueString(), "%d", &maxWeight); err == nil {

			rackRequest.MaxWeight = *netbox.NewNullableInt32(&maxWeight)

		}

	}

	// Set weight_unit

	if !data.WeightUnit.IsNull() {

		weightUnitValue := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())

		rackRequest.WeightUnit = &weightUnitValue

	}

	// Set desc_units

	if !data.DescUnits.IsNull() {

		descUnits := data.DescUnits.ValueBool()

		rackRequest.DescUnits = &descUnits

	}

	// Set outer_width

	if !data.OuterWidth.IsNull() {

		var outerWidth int32

		if _, err := fmt.Sscanf(data.OuterWidth.ValueString(), "%d", &outerWidth); err == nil {

			rackRequest.OuterWidth = *netbox.NewNullableInt32(&outerWidth)

		}

	}

	// Set outer_depth

	if !data.OuterDepth.IsNull() {

		var outerDepth int32

		if _, err := fmt.Sscanf(data.OuterDepth.ValueString(), "%d", &outerDepth); err == nil {

			rackRequest.OuterDepth = *netbox.NewNullableInt32(&outerDepth)

		}

	}

	// Set outer_unit

	if !data.OuterUnit.IsNull() {

		outerUnitValue := netbox.PatchedWritableRackRequestOuterUnit(data.OuterUnit.ValueString())

		rackRequest.OuterUnit = &outerUnitValue

	}

	// Set mounting_depth

	if !data.MountingDepth.IsNull() {

		var mountingDepth int32

		if _, err := fmt.Sscanf(data.MountingDepth.ValueString(), "%d", &mountingDepth); err == nil {

			rackRequest.MountingDepth = *netbox.NewNullableInt32(&mountingDepth)

		}

	}

	// Set airflow

	if !data.Airflow.IsNull() {

		airflowValue := netbox.PatchedWritableRackRequestAirflow(data.Airflow.ValueString())

		rackRequest.Airflow = &airflowValue

	}

	// Set serial

	if !data.Serial.IsNull() {

		serial := data.Serial.ValueString()

		rackRequest.Serial = &serial

	}

	// Set asset_tag

	if !data.AssetTag.IsNull() {

		assetTag := data.AssetTag.ValueString()

		rackRequest.AssetTag = *netbox.NewNullableString(&assetTag)

	}

	// Set description

	if !data.Description.IsNull() {

		description := data.Description.ValueString()

		rackRequest.Description = &description

	}

	// Set comments

	if !data.Comments.IsNull() {

		comments := data.Comments.ValueString()

		rackRequest.Comments = &comments

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return nil

		}

		rackRequest.CustomFields = utils.CustomFieldsToMap(customFields)

	}

	return &rackRequest

}

func (r *RackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data RackResourceModel

	var state RackResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Updating rack", map[string]interface{}{

		"id": state.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Parse ID

	rackID, err := utils.ParseID(state.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Error parsing rack ID",

			fmt.Sprintf("Could not parse rack ID '%s': %s", state.ID.ValueString(), err.Error()),
		)

		return

	}

	rackRequest := r.buildRackRequestForUpdate(ctx, &data, resp)

	if rackRequest == nil {

		return

	}

	// Update the rack via API

	rack, httpResp, err := r.client.DcimAPI.DcimRacksUpdate(ctx, rackID).WritableRackRequest(*rackRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating rack",

			utils.FormatAPIError("update rack", err, httpResp),
		)

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError(

			"Error updating rack",

			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)

		return

	}

	// Map response to state

	mapRackToState(ctx, rack, &data)

	tflog.Debug(ctx, "Updated rack", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data RackResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Deleting rack", map[string]interface{}{

		"id": data.ID.ValueString(),
	})

	// Parse ID

	rackID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Error parsing rack ID",

			fmt.Sprintf("Could not parse rack ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	httpResp, err := r.client.DcimAPI.DcimRacksDestroy(ctx, rackID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Rack already deleted", map[string]interface{}{

				"id": data.ID.ValueString(),
			})

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting rack",

			utils.FormatAPIError("delete rack", err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted rack", map[string]interface{}{

		"id": data.ID.ValueString(),
	})

}

func (r *RackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapRackToState maps a Netbox Rack to the Terraform state model.

func mapRackToState(ctx context.Context, rack *netbox.Rack, data *RackResourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", rack.GetId()))

	data.Name = types.StringValue(rack.GetName())

	// Map site

	if site := rack.GetSite(); site.Id != 0 {

		data.Site = types.StringValue(fmt.Sprintf("%d", site.GetId()))

	}

	// Map location

	if location, ok := rack.GetLocationOk(); ok && location != nil && location.Id != 0 {

		data.Location = types.StringValue(fmt.Sprintf("%d", location.GetId()))

	} else {

		data.Location = types.StringNull()

	}

	// Map tenant

	if tenant, ok := rack.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {

		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))

	} else {

		data.Tenant = types.StringNull()

	}

	// Map status

	if status, ok := rack.GetStatusOk(); ok && status != nil {

		if value, ok := status.GetValueOk(); ok && value != nil {

			data.Status = types.StringValue(string(*value))

		} else {

			data.Status = types.StringNull()

		}

	} else {

		data.Status = types.StringNull()

	}

	// Map role

	if role, ok := rack.GetRoleOk(); ok && role != nil && role.Id != 0 {

		data.Role = types.StringValue(fmt.Sprintf("%d", role.GetId()))

	} else {

		data.Role = types.StringNull()

	}

	// Map serial

	if serial := rack.GetSerial(); serial != "" {

		data.Serial = types.StringValue(serial)

	} else {

		data.Serial = types.StringNull()

	}

	// Map asset_tag

	if assetTag, ok := rack.GetAssetTagOk(); ok && assetTag != nil && *assetTag != "" {

		data.AssetTag = types.StringValue(*assetTag)

	} else {

		data.AssetTag = types.StringNull()

	}

	// Map rack_type

	if rackType, ok := rack.GetRackTypeOk(); ok && rackType != nil && rackType.Id != 0 {

		data.RackType = types.StringValue(fmt.Sprintf("%d", rackType.GetId()))

	} else {

		data.RackType = types.StringNull()

	}

	// Map form_factor

	if formFactor, ok := rack.GetFormFactorOk(); ok && formFactor != nil {

		if value, ok := formFactor.GetValueOk(); ok && value != nil {

			data.FormFactor = types.StringValue(string(*value))

		} else {

			data.FormFactor = types.StringNull()

		}

	} else {

		data.FormFactor = types.StringNull()

	}

	// Map width

	if width, ok := rack.GetWidthOk(); ok && width != nil {

		if value, ok := width.GetValueOk(); ok && value != nil {

			data.Width = types.StringValue(fmt.Sprintf("%d", *value))

		} else {

			data.Width = types.StringNull()

		}

	} else {

		data.Width = types.StringNull()

	}

	// Map u_height

	if uHeight, ok := rack.GetUHeightOk(); ok && uHeight != nil {

		data.UHeight = types.StringValue(fmt.Sprintf("%d", *uHeight))

	} else {

		data.UHeight = types.StringNull()

	}

	// Map starting_unit

	if startingUnit, ok := rack.GetStartingUnitOk(); ok && startingUnit != nil {

		data.StartingUnit = types.StringValue(fmt.Sprintf("%d", *startingUnit))

	} else {

		data.StartingUnit = types.StringNull()

	}

	// Map weight

	if weight, ok := rack.GetWeightOk(); ok && weight != nil {

		data.Weight = types.StringValue(fmt.Sprintf("%g", *weight))

	} else {

		data.Weight = types.StringNull()

	}

	// Map max_weight

	if maxWeight, ok := rack.GetMaxWeightOk(); ok && maxWeight != nil {

		data.MaxWeight = types.StringValue(fmt.Sprintf("%d", *maxWeight))

	} else {

		data.MaxWeight = types.StringNull()

	}

	// Map weight_unit

	if weightUnit, ok := rack.GetWeightUnitOk(); ok && weightUnit != nil {

		if value, ok := weightUnit.GetValueOk(); ok && value != nil {

			data.WeightUnit = types.StringValue(string(*value))

		} else {

			data.WeightUnit = types.StringNull()

		}

	} else {

		data.WeightUnit = types.StringNull()

	}

	// Map desc_units

	if descUnits, ok := rack.GetDescUnitsOk(); ok && descUnits != nil {

		data.DescUnits = types.BoolValue(*descUnits)

	} else {

		data.DescUnits = types.BoolNull()

	}

	// Map outer_width

	if outerWidth, ok := rack.GetOuterWidthOk(); ok && outerWidth != nil {

		data.OuterWidth = types.StringValue(fmt.Sprintf("%d", *outerWidth))

	} else {

		data.OuterWidth = types.StringNull()

	}

	// Map outer_depth

	if outerDepth, ok := rack.GetOuterDepthOk(); ok && outerDepth != nil {

		data.OuterDepth = types.StringValue(fmt.Sprintf("%d", *outerDepth))

	} else {

		data.OuterDepth = types.StringNull()

	}

	// Map outer_unit

	if outerUnit, ok := rack.GetOuterUnitOk(); ok && outerUnit != nil {

		if value, ok := outerUnit.GetValueOk(); ok && value != nil {

			data.OuterUnit = types.StringValue(string(*value))

		} else {

			data.OuterUnit = types.StringNull()

		}

	} else {

		data.OuterUnit = types.StringNull()

	}

	// Map mounting_depth

	if mountingDepth, ok := rack.GetMountingDepthOk(); ok && mountingDepth != nil {

		data.MountingDepth = types.StringValue(fmt.Sprintf("%d", *mountingDepth))

	} else {

		data.MountingDepth = types.StringNull()

	}

	// Map airflow

	if airflow, ok := rack.GetAirflowOk(); ok && airflow != nil {

		if value, ok := airflow.GetValueOk(); ok && value != nil {

			data.Airflow = types.StringValue(string(*value))

		} else {

			data.Airflow = types.StringNull()

		}

	} else {

		data.Airflow = types.StringNull()

	}

	// Map description

	if description := rack.GetDescription(); description != "" {

		data.Description = types.StringValue(description)

	} else {

		data.Description = types.StringNull()

	}

	// Map comments

	if comments := rack.GetComments(); comments != "" {

		data.Comments = types.StringValue(comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Map tags

	if tags := rack.GetTags(); len(tags) > 0 {

		tagModels := utils.NestedTagsToTagModels(tags)

		if tagSet, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tagModels); !diags.HasError() {

			data.Tags = tagSet

		}

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Note: Custom fields are handled separately in Read to preserve type information

}
