// Copyright (c) 2024 Kevin Pelzel
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &VLANDataSource{}
	_ datasource.DataSourceWithConfigure = &VLANDataSource{}
)

// NewVLANDataSource returns a new VLAN data source.
func NewVLANDataSource() datasource.DataSource {
	return &VLANDataSource{}
}

// VLANDataSource defines the data source implementation.
type VLANDataSource struct {
	client *netbox.APIClient
}

// VLANDataSourceModel describes the data source data model.
type VLANDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	VID         types.Int32  `tfsdk:"vid"`
	Name        types.String `tfsdk:"name"`
	Site        types.String `tfsdk:"site"`
	SiteID      types.Int64  `tfsdk:"site_id"`
	Group       types.String `tfsdk:"group"`
	GroupID     types.Int64  `tfsdk:"group_id"`
	Tenant      types.String `tfsdk:"tenant"`
	TenantID    types.Int64  `tfsdk:"tenant_id"`
	Status      types.String `tfsdk:"status"`
	Role        types.String `tfsdk:"role"`
	RoleID      types.Int64  `tfsdk:"role_id"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	Tags        types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *VLANDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vlan"
}

// Schema defines the schema for the data source.
func (d *VLANDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a VLAN in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the VLAN. Either `id` or `name` and `vid` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"vid": schema.Int32Attribute{
				MarkdownDescription: "The VLAN ID (numeric identifier).",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the VLAN.",
				Optional:            true,
				Computed:            true,
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "The name of the site this VLAN is assigned to.",
				Computed:            true,
			},
			"site_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the site this VLAN is assigned to.",
				Computed:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The name of the VLAN group this VLAN belongs to.",
				Computed:            true,
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VLAN group this VLAN belongs to.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant this VLAN is assigned to.",
				Computed:            true,
			},
			"tenant_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the tenant this VLAN is assigned to.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the VLAN (active, reserved, deprecated).",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The name of the role assigned to this VLAN.",
				Computed:            true,
			},
			"role_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the role assigned to this VLAN.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the VLAN.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments for the VLAN.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this VLAN.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *VLANDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *VLANDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VLANDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var vlan *netbox.VLAN

	// Check if we're looking up by ID, VID, or name
	switch {
	case utils.IsSet(data.ID):
		var idInt int
		_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Reading VLAN by ID", map[string]interface{}{
			"id": idInt,
		})

		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}

		result, httpResp, err := d.client.IpamAPI.IpamVlansRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading VLAN",
				utils.FormatAPIError(fmt.Sprintf("retrieve VLAN ID %d", idInt), err, httpResp),
			)
			return
		}
		vlan = result
	case utils.IsSet(data.VID):
		// Looking up by VID (optionally with name and group)
		tflog.Debug(ctx, "Reading VLAN by VID", map[string]interface{}{
			"vid": data.VID.ValueInt32(),
		})

		listReq := d.client.IpamAPI.IpamVlansList(ctx)
		listReq = listReq.Vid([]int32{data.VID.ValueInt32()})

		// Optionally filter by name
		if utils.IsSet(data.Name) {
			listReq = listReq.Name([]string{data.Name.ValueString()})
		}

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing VLANs",
				utils.FormatAPIError(fmt.Sprintf("list VLANs with VID %d", data.VID.ValueInt32()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"VLAN not found",
				fmt.Sprintf("No VLAN found with VID %d", data.VID.ValueInt32()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple VLANs found",
				fmt.Sprintf("Found %d VLANs with VID %d. Please specify additional filters to narrow down the results.", results.Count, data.VID.ValueInt32()),
			)
			return
		}

		vlan = &results.Results[0]
	case utils.IsSet(data.Name):
		// Looking up by name
		tflog.Debug(ctx, "Reading VLAN by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listReq := d.client.IpamAPI.IpamVlansList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing VLANs",
				utils.FormatAPIError(fmt.Sprintf("list VLANs with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"VLAN not found",
				fmt.Sprintf("No VLAN found with name %q", data.Name.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple VLANs found",
				fmt.Sprintf("Found %d VLANs with name %q. Please specify additional filters to narrow down the results.", results.Count, data.Name.ValueString()),
			)
			return
		}

		vlan = &results.Results[0]
	default:
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Either 'id', 'vid', or 'name' must be specified to look up a VLAN.",
		)
		return
	}

	// Map the VLAN to our model
	d.mapVLANToState(ctx, vlan, &data)

	tflog.Debug(ctx, "Read VLAN", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapVLANToState maps a Netbox VLAN to the Terraform state model.
func (d *VLANDataSource) mapVLANToState(ctx context.Context, vlan *netbox.VLAN, data *VLANDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vlan.Id))
	data.VID = types.Int32Value(vlan.Vid)
	data.Name = types.StringValue(vlan.Name)

	// Site
	if vlan.Site.IsSet() && vlan.Site.Get() != nil {
		data.Site = types.StringValue(vlan.Site.Get().Name)
		data.SiteID = types.Int64Value(int64(vlan.Site.Get().Id))
	} else {
		data.Site = types.StringNull()
		data.SiteID = types.Int64Null()
	}

	// Group
	if vlan.Group.IsSet() && vlan.Group.Get() != nil {
		data.Group = types.StringValue(vlan.Group.Get().Name)
		data.GroupID = types.Int64Value(int64(vlan.Group.Get().Id))
	} else {
		data.Group = types.StringNull()
		data.GroupID = types.Int64Null()
	}

	// Tenant
	if vlan.Tenant.IsSet() && vlan.Tenant.Get() != nil {
		data.Tenant = types.StringValue(vlan.Tenant.Get().Name)
		data.TenantID = types.Int64Value(int64(vlan.Tenant.Get().Id))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.Int64Null()
	}

	// Status
	if vlan.Status != nil {
		data.Status = types.StringValue(string(vlan.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Role
	if vlan.Role.IsSet() && vlan.Role.Get() != nil {
		data.Role = types.StringValue(vlan.Role.Get().Name)
		data.RoleID = types.Int64Value(int64(vlan.Role.Get().Id))
	} else {
		data.Role = types.StringNull()
		data.RoleID = types.Int64Null()
	}

	// Description
	if vlan.Description != nil && *vlan.Description != "" {
		data.Description = types.StringValue(*vlan.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if vlan.Comments != nil && *vlan.Comments != "" {
		data.Comments = types.StringValue(*vlan.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags
	if len(vlan.Tags) > 0 {
		tagNames := make([]string, 0, len(vlan.Tags))
		for _, tag := range vlan.Tags {
			tagNames = append(tagNames, tag.Name)
		}
		tagList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}
}
