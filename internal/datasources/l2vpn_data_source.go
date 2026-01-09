// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &L2VPNDataSource{}

func NewL2VPNDataSource() datasource.DataSource {
	return &L2VPNDataSource{}
}

// L2VPNDataSource defines the data source implementation.
type L2VPNDataSource struct {
	client *netbox.APIClient
}

// L2VPNDataSourceModel describes the data source data model.
type L2VPNDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	Slug          types.String `tfsdk:"slug"`
	Type          types.String `tfsdk:"type"`
	Identifier    types.Int64  `tfsdk:"identifier"`
	ImportTargets types.Set    `tfsdk:"import_targets"`
	ExportTargets types.Set    `tfsdk:"export_targets"`
	Tenant        types.String `tfsdk:"tenant"`
	TenantID      types.String `tfsdk:"tenant_id"`
	Description   types.String `tfsdk:"description"`
	Comments      types.String `tfsdk:"comments"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

func (d *L2VPNDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn"
}

func (d *L2VPNDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Layer 2 VPN in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id":           nbschema.DSIDAttribute("L2VPN"),
			"display_name": nbschema.DSComputedStringAttribute("The display name of the L2VPN."),
			"name":         nbschema.DSNameAttribute("L2VPN"),
			"slug":         nbschema.DSSlugAttribute("L2VPN"),
			"type":         nbschema.DSComputedStringAttribute("L2VPN type."),
			"identifier": schema.Int64Attribute{
				MarkdownDescription: "Numeric identifier unique to the parent L2VPN.",
				Computed:            true,
			},
			"import_targets": schema.SetAttribute{
				MarkdownDescription: "Set of route target IDs to import.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"export_targets": schema.SetAttribute{
				MarkdownDescription: "Set of route target IDs to export.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"tenant":        nbschema.DSComputedStringAttribute("Name of the tenant."),
			"tenant_id":     nbschema.DSComputedStringAttribute("ID of the tenant."),
			"description":   nbschema.DSComputedStringAttribute("Description of the L2VPN."),
			"comments":      nbschema.DSComputedStringAttribute("Comments for the L2VPN."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *L2VPNDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *L2VPNDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data L2VPNDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var l2vpn *netbox.L2VPN

	// Lookup by ID
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		var idInt int32
		if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid ID format",
				fmt.Sprintf("Could not parse L2VPN ID '%s': %s", data.ID.ValueString(), parseErr.Error()),
			)
			return
		}
		tflog.Debug(ctx, "Looking up L2VPN by ID", map[string]interface{}{
			"id": idInt,
		})
		result, httpResp, err := d.client.VpnAPI.VpnL2vpnsRetrieve(ctx, idInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading L2VPN",
				utils.FormatAPIError("read L2VPN", err, httpResp),
			)
			return
		}
		l2vpn = result

	case !data.Slug.IsNull() && !data.Slug.IsUnknown():
		// Lookup by slug
		tflog.Debug(ctx, "Looking up L2VPN by slug", map[string]interface{}{
			"slug": data.Slug.ValueString(),
		})
		list, httpResp, err := d.client.VpnAPI.VpnL2vpnsList(ctx).
			Slug([]string{data.Slug.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading L2VPN",
				utils.FormatAPIError("find L2VPN by slug", err, httpResp),
			)
			return
		}
		if list == nil || len(list.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"L2VPN not found",
				fmt.Sprintf("No L2VPN found with slug: %s", data.Slug.ValueString()),
			)
			return
		}
		result := list.GetResults()[0]
		l2vpn = &result

	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by name
		tflog.Debug(ctx, "Looking up L2VPN by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		list, httpResp, err := d.client.VpnAPI.VpnL2vpnsList(ctx).
			Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading L2VPN",
				utils.FormatAPIError("find L2VPN by name", err, httpResp),
			)
			return
		}
		if list == nil || len(list.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"L2VPN not found",
				fmt.Sprintf("No L2VPN found with name: %s", data.Name.ValueString()),
			)
			return
		}
		if len(list.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple L2VPNs found",
				fmt.Sprintf("Found %d L2VPNs with name '%s'. Please use a more specific identifier like 'id' or 'slug'.",
					len(list.GetResults()), data.Name.ValueString()),
			)
			return
		}
		result := list.GetResults()[0]
		l2vpn = &result

	default:
		resp.Diagnostics.AddError(
			"Missing required identifier",
			"Either 'id', 'slug', or 'name' must be specified to lookup an L2VPN.",
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, l2vpn, &data, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps an L2VPN API response to the Terraform state model.
func (d *L2VPNDataSource) mapResponseToState(ctx context.Context, l2vpn *netbox.L2VPN, data *L2VPNDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", l2vpn.GetId()))

	// Display Name
	if l2vpn.GetDisplay() != "" {
		data.DisplayName = types.StringValue(l2vpn.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
	data.Name = types.StringValue(l2vpn.GetName())
	data.Slug = types.StringValue(l2vpn.GetSlug())

	// Type
	if l2vpn.HasType() {
		typeObj := l2vpn.GetType()
		data.Type = types.StringValue(string(typeObj.GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Identifier
	if l2vpn.HasIdentifier() && l2vpn.GetIdentifier() != 0 {
		data.Identifier = types.Int64Value(l2vpn.GetIdentifier())
	} else {
		data.Identifier = types.Int64Null()
	}

	// Description
	if l2vpn.HasDescription() && l2vpn.GetDescription() != "" {
		data.Description = types.StringValue(l2vpn.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if l2vpn.HasComments() && l2vpn.GetComments() != "" {
		data.Comments = types.StringValue(l2vpn.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Tenant
	if l2vpn.HasTenant() && l2vpn.GetTenant().Id != 0 {
		tenant := l2vpn.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Import targets
	if l2vpn.HasImportTargets() && len(l2vpn.GetImportTargets()) > 0 {
		var targetIDs []string
		for _, target := range l2vpn.GetImportTargets() {
			targetIDs = append(targetIDs, fmt.Sprintf("%d", target.GetId()))
		}
		targetSet, diags := types.SetValueFrom(ctx, types.StringType, targetIDs)
		resp.Diagnostics.Append(diags...)
		data.ImportTargets = targetSet
	} else {
		data.ImportTargets = types.SetNull(types.StringType)
	}

	// Export targets
	if l2vpn.HasExportTargets() && len(l2vpn.GetExportTargets()) > 0 {
		var targetIDs []string
		for _, target := range l2vpn.GetExportTargets() {
			targetIDs = append(targetIDs, fmt.Sprintf("%d", target.GetId()))
		}
		targetSet, diags := types.SetValueFrom(ctx, types.StringType, targetIDs)
		resp.Diagnostics.Append(diags...)
		data.ExportTargets = targetSet
	} else {
		data.ExportTargets = types.SetNull(types.StringType)
	}

	// Tags
	if l2vpn.HasTags() && len(l2vpn.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(l2vpn.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields
	if l2vpn.HasCustomFields() && len(l2vpn.GetCustomFields()) > 0 {
		customFields := utils.MapAllCustomFieldsToModels(l2vpn.GetCustomFields())
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
