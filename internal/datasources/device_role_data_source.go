// Package datasources contains Terraform data source implementations for the Netbox provider.
//

// This package integrates with the go-netbox OpenAPI client to provide
// read-only access to Netbox resources via Terraform data sources.

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

var _ datasource.DataSource = &DeviceRoleDataSource{}

func NewDeviceRoleDataSource() datasource.DataSource {
	return &DeviceRoleDataSource{}
}

// DeviceRoleDataSource defines the data source implementation.

type DeviceRoleDataSource struct {
	client *netbox.APIClient
}

// DeviceRoleDataSourceModel describes the data source data model.

type DeviceRoleDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Color types.String `tfsdk:"color"`

	VMRole types.Bool `tfsdk:"vm_role"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`

	DisplayName types.String `tfsdk:"display_name"`
}

func (d *DeviceRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_role"
}

func (d *DeviceRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a device role in Netbox. Device roles categorize devices by their function (e.g., Router, Switch, Server, Firewall). You can identify the device role using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("device role"),
			"name":          nbschema.DSNameAttribute("device role"),
			"slug":          nbschema.DSSlugAttribute("device role"),
			"color":         nbschema.DSComputedStringAttribute("Color for the device role in 6-character hexadecimal format (e.g., 'aa1409')."),
			"vm_role":       nbschema.DSComputedBoolAttribute("Whether virtual machines may be assigned to this role."),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the device role."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the device role."),
		},
	}
}

func (d *DeviceRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.

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

func (d *DeviceRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeviceRoleDataSourceModel

	// Read Terraform configuration data into the model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var deviceRole *netbox.DeviceRole

	var err error

	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name

	switch {
	case !data.ID.IsNull():

		// Search by ID

		deviceRoleID := data.ID.ValueString()

		tflog.Debug(ctx, "Reading device role by ID", map[string]interface{}{
			"id": deviceRoleID,
		})

		// Parse the device role ID to int32 for the API call

		var deviceRoleIDInt int32

		if _, parseErr := fmt.Sscanf(deviceRoleID, "%d", &deviceRoleIDInt); parseErr != nil {
			resp.Diagnostics.AddError(

				"Invalid Device Role ID",

				fmt.Sprintf("Device Role ID must be a number, got: %s", deviceRoleID),
			)

			return
		}

		// Retrieve the device role via API

		deviceRole, httpResp, err = d.client.DcimAPI.DcimDeviceRolesRetrieve(ctx, deviceRoleIDInt).Execute()

		defer utils.CloseResponseBody(httpResp)

	case !data.Slug.IsNull():

		// Search by slug

		deviceRoleSlug := data.Slug.ValueString()

		tflog.Debug(ctx, "Reading device role by slug", map[string]interface{}{
			"slug": deviceRoleSlug,
		})

		// List device roles with slug filter

		var deviceRoles *netbox.PaginatedDeviceRoleList

		deviceRoles, httpResp, err = d.client.DcimAPI.DcimDeviceRolesList(ctx).Slug([]string{deviceRoleSlug}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError(

				"Error reading device role",

				utils.FormatAPIError("read device role by slug", err, httpResp),
			)

			return
		}

		if len(deviceRoles.GetResults()) == 0 {
			resp.Diagnostics.AddError(

				"Device Role Not Found",

				fmt.Sprintf("No device role found with slug: %s", deviceRoleSlug),
			)

			return
		}

		if len(deviceRoles.GetResults()) > 1 {
			resp.Diagnostics.AddError(

				"Multiple Device Roles Found",

				fmt.Sprintf("Multiple device roles found with slug: %s. This should not happen as slugs should be unique.", deviceRoleSlug),
			)

			return
		}

		deviceRole = &deviceRoles.GetResults()[0]

	case !data.Name.IsNull():

		// Search by name

		deviceRoleName := data.Name.ValueString()

		tflog.Debug(ctx, "Reading device role by name", map[string]interface{}{
			"name": deviceRoleName,
		})

		// List device roles with name filter

		var deviceRoles *netbox.PaginatedDeviceRoleList

		deviceRoles, httpResp, err = d.client.DcimAPI.DcimDeviceRolesList(ctx).Name([]string{deviceRoleName}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError(

				"Error reading device role",

				utils.FormatAPIError("read device role by name", err, httpResp),
			)

			return
		}

		if len(deviceRoles.GetResults()) == 0 {
			resp.Diagnostics.AddError(

				"Device Role Not Found",

				fmt.Sprintf("No device role found with name: %s", deviceRoleName),
			)

			return
		}

		if len(deviceRoles.GetResults()) > 1 {
			resp.Diagnostics.AddError(

				"Multiple Device Roles Found",

				fmt.Sprintf("Multiple device roles found with name: %s. Device role names may not be unique in Netbox.", deviceRoleName),
			)

			return
		}

		deviceRole = &deviceRoles.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Device Role Identifier",

			"Either 'id', 'slug', or 'name' must be specified to identify the device role.",
		)

		return
	}

	if err != nil {
		resp.Diagnostics.AddError(

			"Error reading device role",

			utils.FormatAPIError("read device role", err, httpResp),
		)

		return
	}

	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError(

			"Device Role Not Found",

			"The specified device role was not found in Netbox.",
		)

		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(

			"Error reading device role",

			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)

		return
	}

	// Update the model with the response from the API

	data.ID = types.StringValue(fmt.Sprintf("%d", deviceRole.GetId()))

	data.Name = types.StringValue(deviceRole.GetName())

	data.Slug = types.StringValue(deviceRole.GetSlug())

	// Handle color

	if deviceRole.HasColor() && deviceRole.GetColor() != "" {
		data.Color = types.StringValue(deviceRole.GetColor())
	} else {
		data.Color = types.StringNull()
	}

	// Handle vm_role

	if deviceRole.HasVmRole() {
		data.VMRole = types.BoolValue(deviceRole.GetVmRole())
	} else {
		data.VMRole = types.BoolValue(true) // Default to true per Netbox API
	}

	// Handle description

	if deviceRole.HasDescription() && deviceRole.GetDescription() != "" {
		data.Description = types.StringValue(deviceRole.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags

	if deviceRole.HasTags() {
		tags := utils.NestedTagsToTagModels(deviceRole.GetTags())

		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields

	if deviceRole.HasCustomFields() {
		// For data sources, we extract all available custom fields

		customFields := utils.MapToCustomFieldModels(deviceRole.GetCustomFields(), nil)

		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Map display_name
	if deviceRole.GetDisplay() != "" {
		data.DisplayName = types.StringValue(deviceRole.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
