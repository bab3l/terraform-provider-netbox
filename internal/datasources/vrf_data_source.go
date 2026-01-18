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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &VRFDataSource{}

func NewVRFDataSource() datasource.DataSource {
	return &VRFDataSource{}
}

// VRFDataSource defines the data source implementation.
type VRFDataSource struct {
	client *netbox.APIClient
}

// VRFDataSourceModel describes the data source data model.
type VRFDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	RD            types.String `tfsdk:"rd"`
	Tenant        types.String `tfsdk:"tenant"`
	EnforceUnique types.Bool   `tfsdk:"enforce_unique"`
	Description   types.String `tfsdk:"description"`
	DisplayName   types.String `tfsdk:"display_name"`
	Comments      types.String `tfsdk:"comments"`
	Tags          types.List   `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

func (d *VRFDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrf"
}

func (d *VRFDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a VRF (Virtual Routing and Forwarding) instance in Netbox. You can identify the VRF using `id` or `name`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the VRF. Use to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the VRF. Use to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"rd": schema.StringAttribute{
				MarkdownDescription: "Route distinguisher (RD) as defined in RFC 4364.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant this VRF belongs to.",
				Computed:            true,
			},
			"enforce_unique": schema.BoolAttribute{
				MarkdownDescription: "Prevent duplicate prefixes/IP addresses within this VRF.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Brief description of the VRF.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name for the VRF.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the VRF.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "Tags assigned to this VRF.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *VRFDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VRFDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VRFDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var vrf *netbox.VRF

	// Look up by ID or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "":
		var id int32
		if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id); err != nil {
			resp.Diagnostics.AddError("Invalid VRF ID", fmt.Sprintf("VRF ID must be a number, got: %s", data.ID.ValueString()))
			return
		}
		tflog.Debug(ctx, "Looking up VRF by ID", map[string]interface{}{
			"id": id,
		})
		vrfResp, httpResp, err := d.client.IpamAPI.IpamVrfsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VRF",
				utils.FormatAPIError(fmt.Sprintf("read VRF ID %d", id), err, httpResp),
			)
			return
		}
		vrf = vrfResp

	case !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != "":
		// Look up by name
		name := data.Name.ValueString()
		tflog.Debug(ctx, "Looking up VRF by name", map[string]interface{}{
			"name": name,
		})
		list, httpResp, err := d.client.IpamAPI.IpamVrfsList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VRF",
				utils.FormatAPIError(fmt.Sprintf("list VRFs with name %s", name), err, httpResp),
			)
			return
		}
		if list == nil || len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"VRF not found",
				fmt.Sprintf("No VRF found with name: %s", name),
			)
			return
		}
		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple VRFs found",
				fmt.Sprintf("Found %d VRFs with name '%s'. Use 'id' for a unique lookup.", len(list.Results), name),
			)
			return
		}
		vrf = &list.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Either 'id' or 'name' must be specified",
		)
		return
	}
	tflog.Debug(ctx, "Found VRF", map[string]interface{}{
		"id":   vrf.GetId(),
		"name": vrf.GetName(),
	})

	// Map response to state
	d.mapVRFToState(ctx, vrf, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapVRFToState maps a VRF API response to the data source model.
func (d *VRFDataSource) mapVRFToState(ctx context.Context, vrf *netbox.VRF, data *VRFDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vrf.GetId()))
	data.Name = types.StringValue(vrf.GetName())

	// Route distinguisher
	if rd, ok := vrf.GetRdOk(); ok && rd != nil && *rd != "" {
		data.RD = types.StringValue(*rd)
	} else {
		data.RD = types.StringNull()
	}

	// Tenant
	if vrf.HasTenant() && vrf.Tenant.Get() != nil {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", vrf.Tenant.Get().GetId()))
	} else {
		data.Tenant = types.StringNull()
	}

	// Enforce unique
	data.EnforceUnique = types.BoolValue(vrf.GetEnforceUnique())

	// Description
	if desc, ok := vrf.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if comments, ok := vrf.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Display name
	if displayName := vrf.GetDisplay(); displayName != "" {
		data.DisplayName = types.StringValue(displayName)
	} else {
		data.DisplayName = types.StringNull()
	}

	// Tags (slug list)
	if vrf.HasTags() && len(vrf.GetTags()) > 0 {
		tagSlugs := make([]string, 0, len(vrf.GetTags()))
		for _, tag := range vrf.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		tagsValue, tagDiags := types.ListValueFrom(ctx, types.StringType, tagSlugs)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Custom fields - datasources return ALL fields
	if vrf.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(vrf.GetCustomFields())
		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfValueDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
