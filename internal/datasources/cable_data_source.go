// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &CableDataSource{}

func NewCableDataSource() datasource.DataSource {
	return &CableDataSource{}
}

// CableDataSource defines the data source implementation.

type CableDataSource struct {
	client *netbox.APIClient
}

// CableDataSourceModel describes the data source data model.

type CableDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	ATerminations types.List `tfsdk:"a_terminations"`

	BTerminations types.List `tfsdk:"b_terminations"`

	Type types.String `tfsdk:"type"`

	Status types.String `tfsdk:"status"`

	Tenant types.String `tfsdk:"tenant"`

	TenantID types.String `tfsdk:"tenant_id"`

	Label types.String `tfsdk:"label"`

	Color types.String `tfsdk:"color"`

	Length types.Float64 `tfsdk:"length"`

	LengthUnit types.String `tfsdk:"length_unit"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`
}

// TerminationDataSourceModel represents a cable termination point.

type TerminationDataSourceModel struct {
	ObjectType types.String `tfsdk:"object_type"`

	ObjectID types.Int64 `tfsdk:"object_id"`
}

func (d *CableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cable"
}

func (d *CableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	terminationNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"object_type": schema.StringAttribute{
				MarkdownDescription: "Content type of the termination object.",

				Computed: true,
			},

			"object_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the termination object.",

				Computed: true,
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a cable connection in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.DSIDAttribute("cable"),

			"a_terminations": schema.ListNestedAttribute{
				MarkdownDescription: "A-side termination points for this cable.",

				Computed: true,

				NestedObject: terminationNestedObject,
			},

			"b_terminations": schema.ListNestedAttribute{
				MarkdownDescription: "B-side termination points for this cable.",

				Computed: true,

				NestedObject: terminationNestedObject,
			},

			"type": nbschema.DSComputedStringAttribute("Type of cable."),

			"status": nbschema.DSComputedStringAttribute("Connection status."),

			"tenant": nbschema.DSComputedStringAttribute("Name of the tenant that owns this cable."),

			"tenant_id": nbschema.DSComputedStringAttribute("ID of the tenant that owns this cable."),

			"label": nbschema.DSComputedStringAttribute("Physical label attached to the cable."),

			"color": nbschema.DSComputedStringAttribute("Color of the cable (hex code)."),

			"length": schema.Float64Attribute{
				MarkdownDescription: "Length of the cable.",

				Computed: true,
			},

			"length_unit": nbschema.DSComputedStringAttribute("Unit for cable length."),

			"description": nbschema.DSComputedStringAttribute("Description of the cable."),

			"comments": nbschema.DSComputedStringAttribute("Comments about the cable."),

			"display_name": nbschema.DSComputedStringAttribute("The display name of the cable."),

			"tags": nbschema.DSTagsAttribute(),
		},
	}
}

func (d *CableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CableDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var result *netbox.Cable

	// Lookup by ID (required for cables since they don't have name/slug)

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id, err := utils.ParseID(data.ID.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cable ID: %s", err))

			return
		}

		tflog.Debug(ctx, "Reading cable by ID", map[string]interface{}{"id": id})

		cable, httpResp, err := d.client.DcimAPI.DcimCablesRetrieve(ctx, id).Execute()

		defer utils.CloseResponseBody(httpResp)

		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Cable Not Found",
				fmt.Sprintf("No cable found with ID: %d", id),
			)
			return
		}

		if err != nil {
			resp.Diagnostics.AddError(

				"Error reading cable",

				utils.FormatAPIError("read cable", err, httpResp),
			)

			return
		}

		result = cable
	} else {
		resp.Diagnostics.AddError(

			"Missing Identifier",

			"The 'id' attribute must be specified to look up a cable.",
		)

		return
	}

	// Map response to state

	resp.Diagnostics.Append(d.mapResponseToState(ctx, result, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps API response to Terraform state.

func (d *CableDataSource) mapResponseToState(ctx context.Context, result *netbox.Cable, data *CableDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map A terminations

	if result.HasATerminations() {
		aTerms, termDiags := d.mapTerminationsToState(ctx, result.GetATerminations())

		diags.Append(termDiags...)

		data.ATerminations = aTerms
	} else {
		data.ATerminations = types.ListNull(getTerminationDataSourceObjectType())
	}

	// Map B terminations

	if result.HasBTerminations() {
		bTerms, termDiags := d.mapTerminationsToState(ctx, result.GetBTerminations())

		diags.Append(termDiags...)

		data.BTerminations = bTerms
	} else {
		data.BTerminations = types.ListNull(getTerminationDataSourceObjectType())
	}

	// Type

	if result.HasType() && result.GetType() != "" {
		data.Type = types.StringValue(string(result.GetType()))
	} else {
		data.Type = types.StringNull()
	}

	// Status

	if result.HasStatus() {
		status := result.GetStatus()

		data.Status = types.StringValue(string(status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Tenant

	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()

		data.Tenant = types.StringValue(tenant.GetName())

		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()

		data.TenantID = types.StringNull()
	}

	// Label

	if result.HasLabel() && result.GetLabel() != "" {
		data.Label = types.StringValue(result.GetLabel())
	} else {
		data.Label = types.StringNull()
	}

	// Color

	if result.HasColor() && result.GetColor() != "" {
		data.Color = types.StringValue(result.GetColor())
	} else {
		data.Color = types.StringNull()
	}

	// Length

	if result.HasLength() && result.GetLength() != 0 {
		data.Length = types.Float64Value(result.GetLength())
	} else {
		data.Length = types.Float64Null()
	}

	// Length unit

	if result.HasLengthUnit() {
		lengthUnit := result.GetLengthUnit()

		if lengthUnit.GetValue() != "" {
			data.LengthUnit = types.StringValue(string(lengthUnit.GetValue()))
		} else {
			data.LengthUnit = types.StringNull()
		}
	} else {
		data.LengthUnit = types.StringNull()
	}

	// Description

	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments

	if result.HasComments() && result.GetComments() != "" {
		data.Comments = types.StringValue(result.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Tags

	if result.HasTags() && len(result.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(result.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Map display name

	if result.GetDisplay() != "" {
		data.DisplayName = types.StringValue(result.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	return diags
}

// mapTerminationsToState converts API terminations to Terraform state.

func (d *CableDataSource) mapTerminationsToState(ctx context.Context, terminations []netbox.GenericObject) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(terminations) == 0 {
		return types.ListNull(getTerminationDataSourceObjectType()), diags
	}

	models := make([]TerminationDataSourceModel, len(terminations))

	for i, t := range terminations {
		models[i] = TerminationDataSourceModel{
			ObjectType: types.StringValue(t.GetObjectType()),

			ObjectID: types.Int64Value(int64(t.GetObjectId())),
		}
	}

	result, listDiags := types.ListValueFrom(ctx, getTerminationDataSourceObjectType(), models)

	diags.Append(listDiags...)

	return result, diags
}

// getTerminationDataSourceObjectType returns the Terraform object type for terminations.

func getTerminationDataSourceObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"object_type": types.StringType,

			"object_id": types.Int64Type,
		},
	}
}
