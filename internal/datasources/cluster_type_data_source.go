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
var _ datasource.DataSource = &ClusterTypeDataSource{}

// NewClusterTypeDataSource returns a new Cluster Type data source.
func NewClusterTypeDataSource() datasource.DataSource {
	return &ClusterTypeDataSource{}
}

// ClusterTypeDataSource defines the data source implementation.
type ClusterTypeDataSource struct {
	client *netbox.APIClient
}

// ClusterTypeDataSourceModel describes the data source data model.
type ClusterTypeDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	DisplayName  types.String `tfsdk:"display_name"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ClusterTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_type"
}

// Schema defines the schema for the data source.
func (d *ClusterTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a cluster type in Netbox. Cluster types define the technology or platform used for virtualization clusters (e.g., 'VMware vSphere', 'Proxmox', 'Kubernetes'). You can identify the cluster type using `id`, `slug`, or `name`.",
		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("cluster type"),
			"name":          nbschema.DSNameAttribute("cluster type"),
			"slug":          nbschema.DSSlugAttribute("cluster type"),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the cluster type."),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the cluster type."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure sets up the data source with the provider client.
func (d *ClusterTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ClusterTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var clusterType *netbox.ClusterType
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Search by ID
		clusterTypeID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading cluster type by ID", map[string]interface{}{
			"id": clusterTypeID,
		})
		var clusterTypeIDInt int32
		if _, parseErr := fmt.Sscanf(clusterTypeID, "%d", &clusterTypeIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Cluster Type ID",
				fmt.Sprintf("Cluster Type ID must be a number, got: %s", clusterTypeID),
			)
			return
		}
		clusterType, httpResp, err = d.client.VirtualizationAPI.VirtualizationClusterTypesRetrieve(ctx, clusterTypeIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)

	case !data.Slug.IsNull() && !data.Slug.IsUnknown():
		// Search by slug
		clusterTypeSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading cluster type by slug", map[string]interface{}{
			"slug": clusterTypeSlug,
		})
		var clusterTypes *netbox.PaginatedClusterTypeList
		clusterTypes, httpResp, err = d.client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Slug([]string{clusterTypeSlug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading cluster type",
				utils.FormatAPIError("read cluster type by slug", err, httpResp),
			)
			return
		}
		if len(clusterTypes.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Cluster Type Not Found",
				fmt.Sprintf("No cluster type found with slug: %s", clusterTypeSlug),
			)
			return
		}
		if len(clusterTypes.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Cluster Types Found",
				fmt.Sprintf("Multiple cluster types found with slug: %s. This should not happen as slugs should be unique.", clusterTypeSlug),
			)
			return
		}
		clusterType = &clusterTypes.GetResults()[0]

	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Search by name
		clusterTypeName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading cluster type by name", map[string]interface{}{
			"name": clusterTypeName,
		})
		var clusterTypes *netbox.PaginatedClusterTypeList
		clusterTypes, httpResp, err = d.client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Name([]string{clusterTypeName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading cluster type",
				utils.FormatAPIError("read cluster type by name", err, httpResp),
			)
			return
		}
		if len(clusterTypes.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Cluster Type Not Found",
				fmt.Sprintf("No cluster type found with name: %s", clusterTypeName),
			)
			return
		}
		if len(clusterTypes.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Cluster Types Found",
				fmt.Sprintf("Multiple cluster types found with name: %s. Cluster type names may not be unique in Netbox.", clusterTypeName),
			)
			return
		}
		clusterType = &clusterTypes.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Cluster Type Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the cluster type.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster type",
			utils.FormatAPIError("read cluster type", err, httpResp),
		)
		return
	}

	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Cluster Type Not Found",
			fmt.Sprintf("No cluster type found with ID: %s", data.ID.ValueString()),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", clusterType.GetId()))
	data.Name = types.StringValue(clusterType.GetName())
	data.Slug = types.StringValue(clusterType.GetSlug())

	// Handle description
	if clusterType.HasDescription() && clusterType.GetDescription() != "" {
		data.Description = types.StringValue(clusterType.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if clusterType.HasTags() && len(clusterType.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(clusterType.GetTags())
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
	data.CustomFields = utils.CustomFieldsSetFromAPI(ctx, clusterType.HasCustomFields(), clusterType.GetCustomFields(), &resp.Diagnostics)

	// Map display name
	if clusterType.GetDisplay() != "" {
		data.DisplayName = types.StringValue(clusterType.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
	tflog.Debug(ctx, "Read cluster type", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
