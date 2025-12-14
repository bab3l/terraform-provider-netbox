// Package datasources contains Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &RackTypeDataSource{}

	_ datasource.DataSourceWithConfigure = &RackTypeDataSource{}
)

// NewRackTypeDataSource returns a new data source implementing the RackType data source.

func NewRackTypeDataSource() datasource.DataSource {

	return &RackTypeDataSource{}

}

// RackTypeDataSource defines the data source implementation.

type RackTypeDataSource struct {
	client *netbox.APIClient
}

// RackTypeDataSourceModel describes the data source data model.

type RackTypeDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Manufacturer types.String `tfsdk:"manufacturer"`

	Model types.String `tfsdk:"model"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	FormFactor types.String `tfsdk:"form_factor"`

	Width types.Int64 `tfsdk:"width"`

	UHeight types.Int64 `tfsdk:"u_height"`

	StartingUnit types.Int64 `tfsdk:"starting_unit"`

	DescUnits types.Bool `tfsdk:"desc_units"`

	OuterWidth types.Int64 `tfsdk:"outer_width"`

	OuterDepth types.Int64 `tfsdk:"outer_depth"`

	OuterUnit types.String `tfsdk:"outer_unit"`

	Weight types.Float64 `tfsdk:"weight"`

	MaxWeight types.Int64 `tfsdk:"max_weight"`

	WeightUnit types.String `tfsdk:"weight_unit"`

	MountingDepth types.Int64 `tfsdk:"mounting_depth"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *RackTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_rack_type"

}

// Schema defines the schema for the data source.

func (d *RackTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Retrieves information about a rack type in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the rack type. Use this to look up by ID.",

				Optional: true,

				Computed: true,
			},

			"manufacturer": schema.StringAttribute{

				MarkdownDescription: "The manufacturer of this rack type.",

				Optional: true,

				Computed: true,
			},

			"model": schema.StringAttribute{

				MarkdownDescription: "The model name of the rack type. Use this with manufacturer to look up by model.",

				Optional: true,

				Computed: true,
			},

			"slug": schema.StringAttribute{

				MarkdownDescription: "URL-friendly identifier for the rack type.",

				Optional: true,

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the rack type.",

				Computed: true,
			},

			"form_factor": schema.StringAttribute{

				MarkdownDescription: "Form factor of the rack type.",

				Computed: true,
			},

			"width": schema.Int64Attribute{

				MarkdownDescription: "Rail-to-rail width in inches.",

				Computed: true,
			},

			"u_height": schema.Int64Attribute{

				MarkdownDescription: "Height in rack units (U).",

				Computed: true,
			},

			"starting_unit": schema.Int64Attribute{

				MarkdownDescription: "Starting unit number for the rack.",

				Computed: true,
			},

			"desc_units": schema.BoolAttribute{

				MarkdownDescription: "Whether units are numbered top-to-bottom (descending).",

				Computed: true,
			},

			"outer_width": schema.Int64Attribute{

				MarkdownDescription: "Outer dimension of rack (width).",

				Computed: true,
			},

			"outer_depth": schema.Int64Attribute{

				MarkdownDescription: "Outer dimension of rack (depth).",

				Computed: true,
			},

			"outer_unit": schema.StringAttribute{

				MarkdownDescription: "Unit for outer dimensions (mm or in).",

				Computed: true,
			},

			"weight": schema.Float64Attribute{

				MarkdownDescription: "Weight of the rack.",

				Computed: true,
			},

			"max_weight": schema.Int64Attribute{

				MarkdownDescription: "Maximum load capacity for the rack.",

				Computed: true,
			},

			"weight_unit": schema.StringAttribute{

				MarkdownDescription: "Unit for weight.",

				Computed: true,
			},

			"mounting_depth": schema.Int64Attribute{

				MarkdownDescription: "Maximum depth of a mounted device, in millimeters.",

				Computed: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about this rack type.",

				Computed: true,
			},

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *RackTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {

		return

	}

	client, ok := req.ProviderData.(*netbox.APIClient)

	if !ok {

		resp.Diagnostics.AddError(

			"Unexpected Data Source Configure Type",

			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return

	}

	d.client = client

}

// Read refreshes the data source data.

func (d *RackTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data RackTypeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var rackType *netbox.RackType

	// Look up by ID if provided

	switch {

	case !data.ID.IsNull() && !data.ID.IsUnknown():

		rtID, err := utils.ParseID(data.ID.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid Rack Type ID",

				fmt.Sprintf("Rack type ID must be a number, got: %s", data.ID.ValueString()),
			)

			return

		}

		tflog.Debug(ctx, "Reading rack type by ID", map[string]interface{}{

			"id": rtID,
		})

		rt, httpResp, err := d.client.DcimAPI.DcimRackTypesRetrieve(ctx, rtID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading rack type",

				utils.FormatAPIError(fmt.Sprintf("read rack type ID %d", rtID), err, httpResp),
			)

			return

		}

		rackType = rt

	case !data.Model.IsNull() && !data.Model.IsUnknown():

		// Look up by model name

		tflog.Debug(ctx, "Reading rack type by model", map[string]interface{}{

			"model": data.Model.ValueString(),
		})

		listReq := d.client.DcimAPI.DcimRackTypesList(ctx).Model([]string{data.Model.ValueString()})

		// Optionally filter by manufacturer

		if !data.Manufacturer.IsNull() && !data.Manufacturer.IsUnknown() {

			listReq = listReq.Manufacturer([]string{data.Manufacturer.ValueString()})

		}

		listResp, httpResp, err := listReq.Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading rack type",

				utils.FormatAPIError(fmt.Sprintf("read rack type by model %s", data.Model.ValueString()), err, httpResp),
			)

			return

		}

		if listResp.GetCount() == 0 {

			resp.Diagnostics.AddError(

				"Rack type not found",

				fmt.Sprintf("No rack type found with model: %s", data.Model.ValueString()),
			)

			return

		}

		if listResp.GetCount() > 1 {

			resp.Diagnostics.AddError(

				"Multiple rack types found",

				fmt.Sprintf("Found %d rack types with model: %s. Consider filtering by manufacturer as well.", listResp.GetCount(), data.Model.ValueString()),
			)

			return

		}

		rackType = &listResp.GetResults()[0]

	case !data.Slug.IsNull() && !data.Slug.IsUnknown():

		// Look up by slug

		tflog.Debug(ctx, "Reading rack type by slug", map[string]interface{}{

			"slug": data.Slug.ValueString(),
		})

		listResp, httpResp, err := d.client.DcimAPI.DcimRackTypesList(ctx).Slug([]string{data.Slug.ValueString()}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading rack type",

				utils.FormatAPIError(fmt.Sprintf("read rack type by slug %s", data.Slug.ValueString()), err, httpResp),
			)

			return

		}

		if listResp.GetCount() == 0 {

			resp.Diagnostics.AddError(

				"Rack type not found",

				fmt.Sprintf("No rack type found with slug: %s", data.Slug.ValueString()),
			)

			return

		}

		rackType = &listResp.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either 'id', 'model', or 'slug' must be specified to look up a rack type.",
		)

		return

	}

	// Map response to model

	d.mapResponseToModel(ctx, rackType, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (d *RackTypeDataSource) mapResponseToModel(ctx context.Context, rackType *netbox.RackType, data *RackTypeDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", rackType.GetId()))

	data.Model = types.StringValue(rackType.GetModel())

	data.Slug = types.StringValue(rackType.GetSlug())

	// Map manufacturer

	data.Manufacturer = types.StringValue(rackType.Manufacturer.GetName())

	// Map description

	if desc, ok := rackType.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map form_factor

	if ff, ok := rackType.GetFormFactorOk(); ok && ff != nil {

		data.FormFactor = types.StringValue(string(ff.GetValue()))

	} else {

		data.FormFactor = types.StringNull()

	}

	// Map width

	if width, ok := rackType.GetWidthOk(); ok && width != nil {

		data.Width = types.Int64Value(int64(width.GetValue()))

	} else {

		data.Width = types.Int64Null()

	}

	// Map u_height

	if uHeight, ok := rackType.GetUHeightOk(); ok && uHeight != nil {

		data.UHeight = types.Int64Value(int64(*uHeight))

	} else {

		data.UHeight = types.Int64Null()

	}

	// Map starting_unit

	if startingUnit, ok := rackType.GetStartingUnitOk(); ok && startingUnit != nil {

		data.StartingUnit = types.Int64Value(int64(*startingUnit))

	} else {

		data.StartingUnit = types.Int64Null()

	}

	// Map desc_units

	if descUnits, ok := rackType.GetDescUnitsOk(); ok && descUnits != nil {

		data.DescUnits = types.BoolValue(*descUnits)

	} else {

		data.DescUnits = types.BoolNull()

	}

	// Map outer_width

	if outerWidth, ok := rackType.GetOuterWidthOk(); ok && outerWidth != nil {

		data.OuterWidth = types.Int64Value(int64(*outerWidth))

	} else {

		data.OuterWidth = types.Int64Null()

	}

	// Map outer_depth

	if outerDepth, ok := rackType.GetOuterDepthOk(); ok && outerDepth != nil {

		data.OuterDepth = types.Int64Value(int64(*outerDepth))

	} else {

		data.OuterDepth = types.Int64Null()

	}

	// Map outer_unit

	if outerUnit, ok := rackType.GetOuterUnitOk(); ok && outerUnit != nil {

		data.OuterUnit = types.StringValue(string(outerUnit.GetValue()))

	} else {

		data.OuterUnit = types.StringNull()

	}

	// Map weight

	if weight, ok := rackType.GetWeightOk(); ok && weight != nil {

		data.Weight = types.Float64Value(*weight)

	} else {

		data.Weight = types.Float64Null()

	}

	// Map max_weight

	if maxWeight, ok := rackType.GetMaxWeightOk(); ok && maxWeight != nil {

		data.MaxWeight = types.Int64Value(int64(*maxWeight))

	} else {

		data.MaxWeight = types.Int64Null()

	}

	// Map weight_unit

	if weightUnit, ok := rackType.GetWeightUnitOk(); ok && weightUnit != nil {

		data.WeightUnit = types.StringValue(string(weightUnit.GetValue()))

	} else {

		data.WeightUnit = types.StringNull()

	}

	// Map mounting_depth

	if mountingDepth, ok := rackType.GetMountingDepthOk(); ok && mountingDepth != nil {

		data.MountingDepth = types.Int64Value(int64(*mountingDepth))

	} else {

		data.MountingDepth = types.Int64Null()

	}

	// Map comments

	if comments, ok := rackType.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if rackType.HasTags() && len(rackType.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(rackType.GetTags())

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

	if rackType.HasCustomFields() {

		apiCustomFields := rackType.GetCustomFields()

		customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)

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
