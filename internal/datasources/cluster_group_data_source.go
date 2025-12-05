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

var _ datasource.DataSource = &ClusterGroupDataSource{}

func NewClusterGroupDataSource() datasource.DataSource {
	return &ClusterGroupDataSource{}
}

type ClusterGroupDataSource struct {
	client *netbox.APIClient
}

type ClusterGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
}

func (d *ClusterGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_group"
}

func (d *ClusterGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a cluster group in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id":          nbschema.DSIDAttribute("cluster group"),
			"name":        nbschema.DSNameAttribute("cluster group"),
			"slug":        nbschema.DSSlugAttribute("cluster group"),
			"description": nbschema.DSComputedStringAttribute("Description of the cluster group."),
		},
	}
}

func (d *ClusterGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClusterGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clusterGroup *netbox.ClusterGroup
	var httpResp *http.Response
	var err error

	// Lookup by ID first
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id := utils.ParseInt32FromString(data.ID.ValueString())
		if id == 0 {
			resp.Diagnostics.AddError("Invalid ID", "ID must be a number")
			return
		}
		tflog.Debug(ctx, "Looking up cluster group by ID", map[string]interface{}{"id": id})
		clusterGroup, httpResp, err = d.client.VirtualizationAPI.VirtualizationClusterGroupsRetrieve(ctx, id).Execute()
	} else if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
		// Lookup by slug
		slug := data.Slug.ValueString()
		tflog.Debug(ctx, "Looking up cluster group by slug", map[string]interface{}{"slug": slug})
		list, listResp, listErr := d.client.VirtualizationAPI.VirtualizationClusterGroupsList(ctx).Slug([]string{slug}).Execute()
		httpResp = listResp
		err = listErr
		if err == nil && list != nil && len(list.Results) > 0 {
			clusterGroup = &list.Results[0]
		} else if err == nil {
			resp.Diagnostics.AddError("Cluster group not found", fmt.Sprintf("No cluster group found with slug: %s", slug))
			return
		}
	} else if !data.Name.IsNull() && !data.Name.IsUnknown() {
		// Lookup by name
		name := data.Name.ValueString()
		tflog.Debug(ctx, "Looking up cluster group by name", map[string]interface{}{"name": name})
		list, listResp, listErr := d.client.VirtualizationAPI.VirtualizationClusterGroupsList(ctx).Name([]string{name}).Execute()
		httpResp = listResp
		err = listErr
		if err == nil && list != nil {
			if len(list.Results) == 0 {
				resp.Diagnostics.AddError("Cluster group not found", fmt.Sprintf("No cluster group found with name: %s", name))
				return
			}
			if len(list.Results) > 1 {
				resp.Diagnostics.AddError("Multiple cluster groups found", fmt.Sprintf("Found %d cluster groups with name: %s. Use slug or ID for unique lookup.", len(list.Results), name))
				return
			}
			clusterGroup = &list.Results[0]
		}
	} else {
		resp.Diagnostics.AddError("Missing identifier", "Either 'id', 'slug', or 'name' must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster group", utils.FormatAPIError("read cluster group", err, httpResp))
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", clusterGroup.GetId()))
	data.Name = types.StringValue(clusterGroup.GetName())
	data.Slug = types.StringValue(clusterGroup.GetSlug())

	if clusterGroup.HasDescription() && clusterGroup.GetDescription() != "" {
		data.Description = types.StringValue(clusterGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
