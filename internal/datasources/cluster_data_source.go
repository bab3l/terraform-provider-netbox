// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ClusterDataSource{}

// NewClusterDataSource returns a new Cluster data source.
func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

// ClusterDataSource defines the data source implementation.
type ClusterDataSource struct {
	client *netbox.APIClient
}

// ClusterDataSourceModel describes the data source data model.
type ClusterDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Group        types.String `tfsdk:"group"`
	Status       types.String `tfsdk:"status"`
	Tenant       types.String `tfsdk:"tenant"`
	Site         types.String `tfsdk:"site"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the data source.
func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a virtualization cluster in Netbox. Clusters represent a pool of physical resources that can be used to run virtual machines. You can identify the cluster using `id` or `name`.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("cluster"),
			"name":          nbschema.DSNameAttribute("cluster"),
			"type":          nbschema.DSComputedStringAttribute("The cluster type (e.g., 'VMware vSphere', 'Proxmox')."),
			"group":         nbschema.DSComputedStringAttribute("The cluster group this cluster belongs to."),
			"status":        nbschema.DSComputedStringAttribute("The status of the cluster (planned, staging, active, decommissioning, offline)."),
			"tenant":        nbschema.DSComputedStringAttribute("The tenant this cluster is assigned to."),
			"site":          nbschema.DSComputedStringAttribute("The site where this cluster is located."),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the cluster."),
			"comments":      nbschema.DSComputedStringAttribute("Additional comments or notes about the cluster."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure sets up the data source with the provider client.
func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves data from the API.
func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cluster *netbox.Cluster
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID or name
	if !data.ID.IsNull() {
		// Search by ID
		clusterID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading cluster by ID", map[string]interface{}{
			"id": clusterID,
		})

		var clusterIDInt int32
		if _, parseErr := fmt.Sscanf(clusterID, "%d", &clusterIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Cluster ID",
				fmt.Sprintf("Cluster ID must be a number, got: %s", clusterID),
			)
			return
		}

		cluster, httpResp, err = d.client.VirtualizationAPI.VirtualizationClustersRetrieve(ctx, clusterIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
	} else if !data.Name.IsNull() {
		// Search by name
		clusterName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading cluster by name", map[string]interface{}{
			"name": clusterName,
		})

		var clusters *netbox.PaginatedClusterList
		clusters, httpResp, err = d.client.VirtualizationAPI.VirtualizationClustersList(ctx).Name([]string{clusterName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading cluster",
				utils.FormatAPIError("read cluster by name", err, httpResp),
			)
			return
		}
		if len(clusters.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Cluster Not Found",
				fmt.Sprintf("No cluster found with name: %s", clusterName),
			)
			return
		}
		if len(clusters.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Clusters Found",
				fmt.Sprintf("Multiple clusters found with name: %s. Cluster names may not be unique in Netbox.", clusterName),
			)
			return
		}
		cluster = &clusters.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Cluster Identifier",
			"Either 'id' or 'name' must be specified to identify the cluster.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			utils.FormatAPIError("read cluster", err, httpResp),
		)
		return
	}

	if httpResp != nil && httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError(
			"Cluster Not Found",
			fmt.Sprintf("No cluster found with ID: %s", data.ID.ValueString()),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", cluster.GetId()))
	data.Name = types.StringValue(cluster.GetName())

	// Type (always present - required field)
	data.Type = types.StringValue(cluster.Type.GetName())

	// Group
	if cluster.Group.IsSet() && cluster.Group.Get() != nil {
		data.Group = types.StringValue(cluster.Group.Get().GetName())
	} else {
		data.Group = types.StringNull()
	}

	// Status
	if cluster.HasStatus() {
		data.Status = types.StringValue(string(cluster.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Tenant
	if cluster.Tenant.IsSet() && cluster.Tenant.Get() != nil {
		data.Tenant = types.StringValue(cluster.Tenant.Get().GetName())
	} else {
		data.Tenant = types.StringNull()
	}

	// Site
	if cluster.Site.IsSet() && cluster.Site.Get() != nil {
		data.Site = types.StringValue(cluster.Site.Get().GetName())
	} else {
		data.Site = types.StringNull()
	}

	// Description
	if cluster.HasDescription() && cluster.GetDescription() != "" {
		data.Description = types.StringValue(cluster.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if cluster.HasComments() && cluster.GetComments() != "" {
		data.Comments = types.StringValue(cluster.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if cluster.HasTags() && len(cluster.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(cluster.GetTags())
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
	if cluster.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(cluster.GetCustomFields(), nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(cfDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Debug(ctx, "Read cluster", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
