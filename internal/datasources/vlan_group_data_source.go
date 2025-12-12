// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &VLANGroupDataSource{}

func NewVLANGroupDataSource() datasource.DataSource {
	return &VLANGroupDataSource{}
}

// VLANGroupDataSource defines the data source implementation.
type VLANGroupDataSource struct {
	client *netbox.APIClient
}

// VLANGroupDataSourceModel describes the data source data model.
type VLANGroupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	ScopeType    types.String `tfsdk:"scope_type"`
	ScopeID      types.String `tfsdk:"scope_id"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *VLANGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vlan_group"
}

func (d *VLANGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a VLAN Group in Netbox. You can identify the VLAN Group using `id`, `name`, or `slug`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the VLAN Group. Use to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the VLAN Group. Use to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly unique identifier for the VLAN Group. Use to look up by slug.",
				Optional:            true,
				Computed:            true,
			},
			"scope_type": schema.StringAttribute{
				MarkdownDescription: "The type of object this VLAN Group is scoped to.",
				Computed:            true,
			},
			"scope_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the object this VLAN Group is scoped to.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Brief description of the VLAN Group.",
				Computed:            true,
			},
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *VLANGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VLANGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VLANGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var vlanGroup *netbox.VLANGroup

	// Look up by ID if provided
	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		var id int32
		if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id); err != nil {
			resp.Diagnostics.AddError("Invalid VLAN Group ID", fmt.Sprintf("VLAN Group ID must be a number, got: %s", data.ID.ValueString()))
			return
		}

		tflog.Debug(ctx, "Looking up VLAN Group by ID", map[string]interface{}{
			"id": id,
		})

		vlanGroupResp, httpResp, err := d.client.IpamAPI.IpamVlanGroupsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VLAN Group",
				utils.FormatAPIError(fmt.Sprintf("read VLAN Group ID %d", id), err, httpResp),
			)
			return
		}
		vlanGroup = vlanGroupResp
	} else if !data.Slug.IsNull() && data.Slug.ValueString() != "" {
		// Look up by slug
		slug := data.Slug.ValueString()

		tflog.Debug(ctx, "Looking up VLAN Group by slug", map[string]interface{}{
			"slug": slug,
		})

		list, httpResp, err := d.client.IpamAPI.IpamVlanGroupsList(ctx).Slug([]string{slug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VLAN Group",
				utils.FormatAPIError(fmt.Sprintf("list VLAN Groups with slug %s", slug), err, httpResp),
			)
			return
		}

		if list == nil || len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"VLAN Group not found",
				fmt.Sprintf("No VLAN Group found with slug: %s", slug),
			)
			return
		}

		vlanGroup = &list.Results[0]
	} else if !data.Name.IsNull() && data.Name.ValueString() != "" {
		// Look up by name
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Looking up VLAN Group by name", map[string]interface{}{
			"name": name,
		})

		list, httpResp, err := d.client.IpamAPI.IpamVlanGroupsList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VLAN Group",
				utils.FormatAPIError(fmt.Sprintf("list VLAN Groups with name %s", name), err, httpResp),
			)
			return
		}

		if list == nil || len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"VLAN Group not found",
				fmt.Sprintf("No VLAN Group found with name: %s", name),
			)
			return
		}

		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple VLAN Groups found",
				fmt.Sprintf("Found %d VLAN Groups with name '%s'. Use 'id' or 'slug' for a unique lookup.", len(list.Results), name),
			)
			return
		}

		vlanGroup = &list.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Either 'id', 'name', or 'slug' must be specified",
		)
		return
	}

	tflog.Debug(ctx, "Found VLAN Group", map[string]interface{}{
		"id":   vlanGroup.GetId(),
		"name": vlanGroup.GetName(),
	})

	// Map response to state
	d.mapVLANGroupToState(ctx, vlanGroup, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapVLANGroupToState maps a VLANGroup API response to the data source model.
func (d *VLANGroupDataSource) mapVLANGroupToState(ctx context.Context, vlanGroup *netbox.VLANGroup, data *VLANGroupDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vlanGroup.GetId()))
	data.Name = types.StringValue(vlanGroup.GetName())
	data.Slug = types.StringValue(vlanGroup.GetSlug())

	// Scope type
	if scopeType, ok := vlanGroup.GetScopeTypeOk(); ok && scopeType != nil && *scopeType != "" {
		data.ScopeType = types.StringValue(*scopeType)
	} else {
		data.ScopeType = types.StringNull()
	}

	// Scope ID
	if vlanGroup.HasScopeId() && vlanGroup.ScopeId.Get() != nil {
		data.ScopeID = types.StringValue(fmt.Sprintf("%d", *vlanGroup.ScopeId.Get()))
	} else {
		data.ScopeID = types.StringNull()
	}

	// Description
	if desc, ok := vlanGroup.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Tags
	if vlanGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(vlanGroup.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields
	if vlanGroup.HasCustomFields() {
		cf := vlanGroup.GetCustomFields()
		customFields := utils.MapToCustomFieldModels(cf, nil)
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
