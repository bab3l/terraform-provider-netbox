// Package datasources contains Terraform data source implementations for the Netbox provider.
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
var _ datasource.DataSource = &FHRPGroupDataSource{}

func NewFHRPGroupDataSource() datasource.DataSource {
	return &FHRPGroupDataSource{}
}

// FHRPGroupDataSource defines the data source implementation.
type FHRPGroupDataSource struct {
	client *netbox.APIClient
}

// FHRPGroupDataSourceModel describes the data source data model.
type FHRPGroupDataSourceModel struct {
	ID          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Protocol    types.String `tfsdk:"protocol"`
	GroupID     types.Int32  `tfsdk:"group_id"`
	AuthType    types.String `tfsdk:"auth_type"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
}

// Metadata returns the data source type name.
func (d *FHRPGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fhrp_group"
}

// Schema defines the schema for the data source.
func (d *FHRPGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an FHRP (First Hop Redundancy Protocol) group in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the FHRP group. Use for lookup when specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the FHRP group.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The redundancy protocol (vrrp2, vrrp3, carp, clusterxl, hsrp, glbp, other). Used with group_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"group_id": schema.Int32Attribute{
				MarkdownDescription: "The FHRP group identifier. Used with protocol for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type (plaintext, md5).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the FHRP group.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments about the FHRP group.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *FHRPGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves the FHRP group data from NetBox.
func (d *FHRPGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FHRPGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fhrpGroup *netbox.FHRPGroup

	// Lookup by ID if provided
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id := data.ID.ValueInt32()

		tflog.Debug(ctx, "Looking up FHRP group by ID", map[string]interface{}{
			"id": id,
		})

		result, httpResp, err := d.client.IpamAPI.IpamFhrpGroupsRetrieve(ctx, id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading FHRP Group",
				utils.FormatAPIError("reading FHRP group by ID", err, httpResp),
			)
			return
		}
		fhrpGroup = result
	} else if !data.Protocol.IsNull() && !data.Protocol.IsUnknown() && !data.GroupID.IsNull() && !data.GroupID.IsUnknown() {
		// Lookup by protocol and group_id
		protocol := data.Protocol.ValueString()
		groupID := data.GroupID.ValueInt32()

		tflog.Debug(ctx, "Looking up FHRP group by protocol and group_id", map[string]interface{}{
			"protocol": protocol,
			"group_id": groupID,
		})

		listReq := d.client.IpamAPI.IpamFhrpGroupsList(ctx)
		listReq = listReq.Protocol([]string{protocol})
		listReq = listReq.GroupId([]int32{groupID})

		list, httpResp, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading FHRP Group",
				utils.FormatAPIError("listing FHRP groups", err, httpResp),
			)
			return
		}

		if list.Count == 0 {
			resp.Diagnostics.AddError(
				"FHRP Group Not Found",
				fmt.Sprintf("No FHRP group found with protocol %q and group_id %d", protocol, groupID),
			)
			return
		}

		if list.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple FHRP Groups Found",
				fmt.Sprintf("Found %d FHRP groups with protocol %q and group_id %d, expected exactly 1", list.Count, protocol, groupID),
			)
			return
		}

		fhrpGroup = &list.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Either 'id' or both 'protocol' and 'group_id' must be specified for FHRP group lookup",
		)
		return
	}

	// Map response to state
	data.ID = types.Int32Value(fhrpGroup.GetId())
	data.Protocol = types.StringValue(string(fhrpGroup.Protocol))
	data.GroupID = types.Int32Value(fhrpGroup.GetGroupId())

	// Name
	if name := fhrpGroup.GetName(); name != "" {
		data.Name = types.StringValue(name)
	} else {
		data.Name = types.StringNull()
	}

	// Auth Type
	if fhrpGroup.AuthType != nil {
		authTypeValue := string(*fhrpGroup.AuthType)
		if authTypeValue != "" {
			data.AuthType = types.StringValue(authTypeValue)
		} else {
			data.AuthType = types.StringNull()
		}
	} else {
		data.AuthType = types.StringNull()
	}

	// Description
	if description := fhrpGroup.GetDescription(); description != "" {
		data.Description = types.StringValue(description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if comments := fhrpGroup.GetComments(); comments != "" {
		data.Comments = types.StringValue(comments)
	} else {
		data.Comments = types.StringNull()
	}

	tflog.Debug(ctx, "Read FHRP group", map[string]interface{}{
		"id":       data.ID.ValueInt32(),
		"protocol": data.Protocol.ValueString(),
		"group_id": data.GroupID.ValueInt32(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
