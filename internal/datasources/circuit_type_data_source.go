// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &CircuitTypeDataSource{}

// NewCircuitTypeDataSource returns a new Circuit Type data source.

func NewCircuitTypeDataSource() datasource.DataSource {

	return &CircuitTypeDataSource{}

}

// CircuitTypeDataSource defines the data source implementation for circuit types.

type CircuitTypeDataSource struct {
	client *netbox.APIClient
}

// CircuitTypeDataSourceModel describes the data source data model.

type CircuitTypeDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	Color types.String `tfsdk:"color"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *CircuitTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_circuit_type"

}

// Schema defines the schema for the data source.

func (d *CircuitTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Retrieves information about a circuit type in Netbox. Circuit types categorize the various types of circuits used by your organization (e.g., Internet Transit, MPLS, Point-to-Point, etc.).",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.DSIDAttribute("circuit type"),

			"name": nbschema.DSNameAttribute("circuit type"),

			"slug": nbschema.DSSlugAttribute("circuit type"),

			"description": nbschema.DSComputedStringAttribute("Description of the circuit type."),

			"color": nbschema.DSComputedStringAttribute("Color of the circuit type (6-character hex code)."),

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

// Configure sets up the data source with the provider client.

func (d *CircuitTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read reads the data source.

func (d *CircuitTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data CircuitTypeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var circuitType *netbox.CircuitType

	var err error

	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name

	switch {

	case !data.ID.IsNull():

		// Search by ID

		circuitTypeID := data.ID.ValueString()

		tflog.Debug(ctx, "Reading circuit type by ID", map[string]interface{}{

			"id": circuitTypeID,
		})

		var circuitTypeIDInt int32

		if _, parseErr := fmt.Sscanf(circuitTypeID, "%d", &circuitTypeIDInt); parseErr != nil {

			resp.Diagnostics.AddError(

				"Invalid Circuit Type ID",

				fmt.Sprintf("Circuit Type ID must be a number, got: %s", circuitTypeID),
			)

			return

		}

		circuitType, httpResp, err = d.client.CircuitsAPI.CircuitsCircuitTypesRetrieve(ctx, circuitTypeIDInt).Execute()

		defer utils.CloseResponseBody(httpResp)

	case !data.Slug.IsNull():

		// Search by slug

		circuitTypeSlug := data.Slug.ValueString()

		tflog.Debug(ctx, "Reading circuit type by slug", map[string]interface{}{

			"slug": circuitTypeSlug,
		})

		var circuitTypes *netbox.PaginatedCircuitTypeList

		circuitTypes, httpResp, err = d.client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Slug([]string{circuitTypeSlug}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading circuit type",

				utils.FormatAPIError("read circuit type by slug", err, httpResp),
			)

			return

		}

		if len(circuitTypes.GetResults()) == 0 {

			resp.Diagnostics.AddError(

				"Circuit Type Not Found",

				fmt.Sprintf("No circuit type found with slug: %s", circuitTypeSlug),
			)

			return

		}

		circuitType = &circuitTypes.GetResults()[0]

	case !data.Name.IsNull():

		// Search by name

		circuitTypeName := data.Name.ValueString()

		tflog.Debug(ctx, "Reading circuit type by name", map[string]interface{}{

			"name": circuitTypeName,
		})

		var circuitTypes *netbox.PaginatedCircuitTypeList

		circuitTypes, httpResp, err = d.client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Name([]string{circuitTypeName}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading circuit type",

				utils.FormatAPIError("read circuit type by name", err, httpResp),
			)

			return

		}

		if len(circuitTypes.GetResults()) == 0 {

			resp.Diagnostics.AddError(

				"Circuit Type Not Found",

				fmt.Sprintf("No circuit type found with name: %s", circuitTypeName),
			)

			return

		}

		circuitType = &circuitTypes.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"At least one of 'id', 'name', or 'slug' must be specified to look up a circuit type.",
		)

		return

	}

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading circuit type",

			utils.FormatAPIError("read circuit type", err, httpResp),
		)

		return

	}

	// Map the circuit type to state

	data.ID = types.StringValue(fmt.Sprintf("%d", circuitType.GetId()))

	data.Name = types.StringValue(circuitType.GetName())

	data.Slug = types.StringValue(circuitType.GetSlug())

	// Handle description

	if circuitType.HasDescription() && circuitType.GetDescription() != "" {

		data.Description = types.StringValue(circuitType.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Handle color

	if circuitType.HasColor() && circuitType.GetColor() != "" {

		data.Color = types.StringValue(circuitType.GetColor())

	} else {

		data.Color = types.StringNull()

	}

	// Handle tags

	if circuitType.HasTags() {

		tags := utils.NestedTagsToTagModels(circuitType.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	if circuitType.HasCustomFields() {

		customFields := utils.MapToCustomFieldModels(circuitType.GetCustomFields(), nil)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		resp.Diagnostics.Append(cfDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

	// Map display name

	if circuitType.GetDisplay() != "" {
		data.DisplayName = types.StringValue(circuitType.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	tflog.Debug(ctx, "Read circuit type", map[string]interface{}{

		"id": circuitType.GetId(),

		"name": circuitType.GetName(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
