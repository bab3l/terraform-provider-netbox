// Package datasources contains Terraform data source implementations for NetBox objects.

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
var (
	_ datasource.DataSource              = &VirtualChassisDataSource{}
	_ datasource.DataSourceWithConfigure = &VirtualChassisDataSource{}
)

// NewVirtualChassisDataSource returns a new data source implementing the VirtualChassis data source.
func NewVirtualChassisDataSource() datasource.DataSource {
	return &VirtualChassisDataSource{}
}

// VirtualChassisDataSource defines the data source implementation.
type VirtualChassisDataSource struct {
	client *netbox.APIClient
}

// VirtualChassisDataSourceModel describes the data source data model.
type VirtualChassisDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Domain       types.String `tfsdk:"domain"`
	Master       types.String `tfsdk:"master"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	MemberCount  types.Int64  `tfsdk:"member_count"`
	DisplayName  types.String `tfsdk:"display_name"`
	Tags         types.List   `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *VirtualChassisDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_chassis"
}

// Schema defines the schema for the data source.
func (d *VirtualChassisDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a virtual chassis in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the virtual chassis. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the virtual chassis. Use this to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain for this virtual chassis.",
				Computed:            true,
			},
			"master": schema.StringAttribute{
				MarkdownDescription: "ID of the master device for this virtual chassis.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the virtual chassis.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about this virtual chassis.",
				Computed:            true,
			},
			"member_count": schema.Int64Attribute{
				MarkdownDescription: "Number of member devices in this virtual chassis.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the virtual chassis.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "Tags assigned to this virtual chassis.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *VirtualChassisDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *VirtualChassisDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtualChassisDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var vc *netbox.VirtualChassis

	// Look up by ID or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		vcID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Virtual Chassis ID",
				fmt.Sprintf("Virtual chassis ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}
		tflog.Debug(ctx, "Reading virtual chassis by ID", map[string]interface{}{
			"id": vcID,
		})
		v, httpResp, err := d.client.DcimAPI.DcimVirtualChassisRetrieve(ctx, vcID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading virtual chassis",
				utils.FormatAPIError(fmt.Sprintf("read virtual chassis ID %d", vcID), err, httpResp),
			)
			return
		}
		vc = v

	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Look up by name
		tflog.Debug(ctx, "Reading virtual chassis by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		listResp, httpResp, err := d.client.DcimAPI.DcimVirtualChassisList(ctx).Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading virtual chassis",
				utils.FormatAPIError(fmt.Sprintf("read virtual chassis by name %s", data.Name.ValueString()), err, httpResp),
			)
			return
		}
		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Virtual chassis not found",
				fmt.Sprintf("No virtual chassis found with name: %s", data.Name.ValueString()),
			)
			return
		}
		if listResp.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple virtual chassis found",
				fmt.Sprintf("Found %d virtual chassis with name: %s.", listResp.GetCount(), data.Name.ValueString()),
			)
			return
		}
		vc = &listResp.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to look up a virtual chassis.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, vc, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *VirtualChassisDataSource) mapResponseToModel(ctx context.Context, vc *netbox.VirtualChassis, data *VirtualChassisDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vc.GetId()))
	data.Name = types.StringValue(vc.GetName())

	// Map domain
	if domain, ok := vc.GetDomainOk(); ok && domain != nil && *domain != "" {
		data.Domain = types.StringValue(*domain)
	} else {
		data.Domain = types.StringNull()
	}

	// Map master
	if vc.Master.IsSet() && vc.Master.Get() != nil {
		data.Master = types.StringValue(vc.Master.Get().GetName())
	} else {
		data.Master = types.StringNull()
	}

	// Map description
	if desc, ok := vc.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := vc.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle display_name
	if vc.GetDisplay() != "" {
		data.DisplayName = types.StringValue(vc.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Map member_count
	data.MemberCount = types.Int64Value(int64(vc.GetMemberCount()))

	// Handle tags (slug list)
	if vc.HasTags() && len(vc.GetTags()) > 0 {
		tagSlugs := make([]string, 0, len(vc.GetTags()))
		for _, tag := range vc.GetTags() {
			tagSlugs = append(tagSlugs, tag.Slug)
		}
		tagList, tagDiags := types.ListValueFrom(ctx, types.StringType, tagSlugs)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Handle custom fields
	if vc.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(vc.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
