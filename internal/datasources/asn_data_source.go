// Package datasources contains Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ASNDataSource{}
	_ datasource.DataSourceWithConfigure = &ASNDataSource{}
)

// NewASNDataSource returns a new data source implementing the ASN data source.
func NewASNDataSource() datasource.DataSource {
	return &ASNDataSource{}
}

// ASNDataSource defines the data source implementation.
type ASNDataSource struct {
	client *netbox.APIClient
}

// ASNDataSourceModel describes the data source data model.
type ASNDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	ASN           types.Int64  `tfsdk:"asn"`
	RIR           types.String `tfsdk:"rir"`
	RIRID         types.String `tfsdk:"rir_id"`
	Tenant        types.String `tfsdk:"tenant"`
	TenantID      types.String `tfsdk:"tenant_id"`
	Description   types.String `tfsdk:"description"`
	Comments      types.String `tfsdk:"comments"`
	DisplayName   types.String `tfsdk:"display_name"`
	SiteCount     types.Int64  `tfsdk:"site_count"`
	ProviderCount types.Int64  `tfsdk:"provider_count"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ASNDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asn"
}

// Schema defines the schema for the data source.
func (d *ASNDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an Autonomous System Number (ASN) in NetBox. You can identify the ASN using `id` or `asn`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the ASN resource. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"asn": schema.Int64Attribute{
				MarkdownDescription: "The 16- or 32-bit autonomous system number. Use this to look up by ASN.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					validators.ValidASNInt64(),
				},
			},
			"rir": schema.StringAttribute{
				MarkdownDescription: "The Regional Internet Registry (RIR) that manages this ASN.",
				Computed:            true,
			},
			"rir_id": nbschema.DSComputedStringAttribute("ID of the Regional Internet Registry (RIR) that manages this ASN."),
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The tenant this ASN is assigned to.",
				Computed:            true,
			},
			"tenant_id": nbschema.DSComputedStringAttribute("ID of the tenant this ASN is assigned to."),
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of this ASN.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about this ASN.",
				Computed:            true,
			},
			"site_count": schema.Int64Attribute{
				MarkdownDescription: "Number of sites using this ASN.",
				Computed:            true,
			},
			"provider_count": schema.Int64Attribute{
				MarkdownDescription: "Number of providers using this ASN.",
				Computed:            true,
			},
			"tags":          nbschema.DSTagsAttribute(),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the ASN."),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ASNDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ASNDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ASNDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var asn *netbox.ASN

	// Look up by ID if provided
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		asnID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ASN ID",
				fmt.Sprintf("ASN ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}
		tflog.Debug(ctx, "Reading ASN by ID", map[string]interface{}{
			"id": asnID,
		})
		a, httpResp, err := d.client.IpamAPI.IpamAsnsRetrieve(ctx, asnID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading ASN",
				utils.FormatAPIError(fmt.Sprintf("read ASN ID %d", asnID), err, httpResp),
			)
			return
		}
		asn = a

	case !data.ASN.IsNull() && !data.ASN.IsUnknown():
		// Look up by ASN number
		tflog.Debug(ctx, "Reading ASN by number", map[string]interface{}{
			"asn": data.ASN.ValueInt64(),
		})
		asn32, err := utils.SafeInt32FromValue(data.ASN)
		if err != nil {
			resp.Diagnostics.AddError("Invalid ASN", fmt.Sprintf("ASN value overflow: %s", err))
			return
		}
		listResp, httpResp, err := d.client.IpamAPI.IpamAsnsList(ctx).Asn([]int32{asn32}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading ASN",
				utils.FormatAPIError(fmt.Sprintf("read ASN %d", data.ASN.ValueInt64()), err, httpResp),
			)
			return
		}
		result, ok := utils.ExpectSingleResult(
			listResp.GetResults(),
			"ASN not found",
			fmt.Sprintf("No ASN found with number: %d", data.ASN.ValueInt64()),
			"Multiple ASNs found",
			fmt.Sprintf("Found %d ASNs with number: %d", listResp.GetCount(), data.ASN.ValueInt64()),
			&resp.Diagnostics,
		)
		if !ok {
			return
		}
		asn = result

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'asn' must be specified to look up an ASN.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, asn, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ASNDataSource) mapResponseToModel(ctx context.Context, asn *netbox.ASN, data *ASNDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", asn.GetId()))
	data.ASN = types.Int64Value(asn.GetAsn())

	// Map RIR
	if asn.Rir.IsSet() && asn.Rir.Get() != nil {
		data.RIR = types.StringValue(asn.Rir.Get().GetName())
		data.RIRID = types.StringValue(fmt.Sprintf("%d", asn.Rir.Get().GetId()))
	} else {
		data.RIR = types.StringNull()
		data.RIRID = types.StringNull()
	}

	// Map Tenant
	if asn.Tenant.IsSet() && asn.Tenant.Get() != nil {
		data.Tenant = types.StringValue(asn.Tenant.Get().GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", asn.Tenant.Get().GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Map description
	if desc, ok := asn.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := asn.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Map counts
	data.SiteCount = types.Int64Value(asn.GetSiteCount())
	data.ProviderCount = types.Int64Value(asn.GetProviderCount())

	// Handle tags
	if asn.HasTags() && len(asn.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(asn.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields - datasources return ALL fields
	data.CustomFields = utils.CustomFieldsSetFromAPI(ctx, asn.HasCustomFields(), asn.GetCustomFields(), diags)

	// Map display name
	if asn.GetDisplay() != "" {
		data.DisplayName = types.StringValue(asn.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
