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
var _ datasource.DataSource = &CircuitDataSource{}

// NewCircuitDataSource returns a new Circuit data source.
func NewCircuitDataSource() datasource.DataSource {
	return &CircuitDataSource{}
}

// CircuitDataSource defines the data source implementation for circuits.
type CircuitDataSource struct {
	client *netbox.APIClient
}

// CircuitDataSourceModel describes the data source data model.
type CircuitDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Cid             types.String `tfsdk:"cid"`
	CircuitProvider types.String `tfsdk:"circuit_provider"`
	ProviderAccount types.String `tfsdk:"provider_account"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	Tenant          types.String `tfsdk:"tenant"`
	InstallDate     types.String `tfsdk:"install_date"`
	TerminationDate types.String `tfsdk:"termination_date"`
	CommitRate      types.Int64  `tfsdk:"commit_rate"`
	Description     types.String `tfsdk:"description"`
	Comments        types.String `tfsdk:"comments"`
	DisplayName     types.String `tfsdk:"display_name"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *CircuitDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit"
}

// Schema defines the schema for the data source.
func (d *CircuitDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a circuit in Netbox. Circuits represent physical or logical network connections provided by external carriers or service providers. You can identify the circuit using `id` or `cid`.",
		Attributes: map[string]schema.Attribute{
			"id":               nbschema.DSIDAttribute("circuit"),
			"cid":              nbschema.DSNameAttribute("circuit"), // Using DSNameAttribute since cid is similar to name
			"circuit_provider": nbschema.DSComputedStringAttribute("The circuit provider (carrier or ISP) name."),
			"provider_account": nbschema.DSComputedStringAttribute("The provider account for this circuit (account identifier)."),
			"type":             nbschema.DSComputedStringAttribute("The type of circuit."),
			"status":           nbschema.DSComputedStringAttribute("The operational status of the circuit."),
			"tenant":           nbschema.DSComputedStringAttribute("The tenant that owns this circuit."),
			"install_date":     nbschema.DSComputedStringAttribute("The date when the circuit was installed."),
			"termination_date": nbschema.DSComputedStringAttribute("The date when the circuit will be or was terminated."),
			"commit_rate":      nbschema.DSComputedInt64Attribute("The committed information rate (CIR) in Kbps for this circuit."),
			"description":      nbschema.DSComputedStringAttribute("Description of the circuit."),
			"comments":         nbschema.DSComputedStringAttribute("Additional comments or notes about the circuit."),
			"tags":             nbschema.DSTagsAttribute(),
			"display_name":     nbschema.DSComputedStringAttribute("The display name of the circuit."),
			"custom_fields":    nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure sets up the data source with the provider client.

func (d *CircuitDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *CircuitDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CircuitDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var circuit *netbox.Circuit
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID or cid
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "":
		// Search by ID
		circuitID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading circuit by ID", map[string]interface{}{
			"id": circuitID,
		})
		var circuitIDInt int32
		if _, parseErr := fmt.Sscanf(circuitID, "%d", &circuitIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Circuit ID",
				fmt.Sprintf("Circuit ID must be a number, got: %s", circuitID),
			)
			return
		}
		circuit, httpResp, err = d.client.CircuitsAPI.CircuitsCircuitsRetrieve(ctx, circuitIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)

	case !data.Cid.IsNull() && !data.Cid.IsUnknown() && data.Cid.ValueString() != "":
		// Search by cid
		circuitCid := data.Cid.ValueString()
		tflog.Debug(ctx, "Reading circuit by cid", map[string]interface{}{
			"cid": circuitCid,
		})
		var circuits *netbox.PaginatedCircuitList
		circuits, httpResp, err = d.client.CircuitsAPI.CircuitsCircuitsList(ctx).Cid([]string{circuitCid}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading circuit",
				utils.FormatAPIError("read circuit by cid", err, httpResp),
			)
			return
		}
		if len(circuits.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Circuit Not Found",
				fmt.Sprintf("No circuit found with cid: %s", circuitCid),
			)
			return
		}
		if len(circuits.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Circuits Found",
				fmt.Sprintf("Multiple circuits found with cid: %s. Please use ID for a specific circuit.", circuitCid),
			)
			return
		}
		circuit = &circuits.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"At least one of 'id' or 'cid' must be specified to look up a circuit.",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading circuit",
			utils.FormatAPIError("read circuit", err, httpResp),
		)
		return
	}

	// Map the circuit to state
	data.ID = types.StringValue(fmt.Sprintf("%d", circuit.GetId()))
	data.Cid = types.StringValue(circuit.GetCid())
	data.CircuitProvider = types.StringValue(circuit.GetProvider().Name)
	data.Type = types.StringValue(circuit.GetType().Name)

	// Provider account
	if circuit.ProviderAccount.IsSet() && circuit.ProviderAccount.Get() != nil {
		data.ProviderAccount = types.StringValue(circuit.ProviderAccount.Get().GetAccount())
	} else {
		data.ProviderAccount = types.StringNull()
	}

	// Handle status
	if circuit.HasStatus() {
		data.Status = types.StringValue(string(circuit.Status.GetValue()))
	} else {
		data.Status = types.StringValue("active")
	}

	// Handle tenant
	if circuit.Tenant.IsSet() && circuit.Tenant.Get() != nil {
		data.Tenant = types.StringValue(circuit.Tenant.Get().GetName())
	} else {
		data.Tenant = types.StringNull()
	}

	// Handle install date
	if circuit.InstallDate.IsSet() && circuit.InstallDate.Get() != nil {
		data.InstallDate = types.StringValue(*circuit.InstallDate.Get())
	} else {
		data.InstallDate = types.StringNull()
	}

	// Handle termination date
	if circuit.TerminationDate.IsSet() && circuit.TerminationDate.Get() != nil {
		data.TerminationDate = types.StringValue(*circuit.TerminationDate.Get())
	} else {
		data.TerminationDate = types.StringNull()
	}

	// Handle commit rate
	if circuit.CommitRate.IsSet() && circuit.CommitRate.Get() != nil {
		data.CommitRate = types.Int64Value(int64(*circuit.CommitRate.Get()))
	} else {
		data.CommitRate = types.Int64Null()
	}

	// Handle description
	if circuit.HasDescription() && circuit.GetDescription() != "" {
		data.Description = types.StringValue(circuit.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments
	if circuit.HasComments() && circuit.GetComments() != "" {
		data.Comments = types.StringValue(circuit.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if circuit.HasTags() {
		tags := utils.NestedTagsToTagModels(circuit.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields - datasources return ALL fields
	if circuit.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(circuit.GetCustomFields())
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
	if circuit.GetDisplay() != "" {
		data.DisplayName = types.StringValue(circuit.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
	tflog.Debug(ctx, "Read circuit", map[string]interface{}{
		"id":  circuit.GetId(),
		"cid": circuit.GetCid(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
