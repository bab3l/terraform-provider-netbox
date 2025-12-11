// Package datasources provides Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &NotificationGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &NotificationGroupDataSource{}
)

// NewNotificationGroupDataSource returns a new data source implementing the notification group data source.
func NewNotificationGroupDataSource() datasource.DataSource {
	return &NotificationGroupDataSource{}
}

// NotificationGroupDataSource defines the data source implementation.
type NotificationGroupDataSource struct {
	client *netbox.APIClient
}

// NotificationGroupDataSourceModel describes the data source model.
type NotificationGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	GroupIDs    types.Set    `tfsdk:"group_ids"`
	Groups      types.Set    `tfsdk:"groups"`
	UserIDs     types.Set    `tfsdk:"user_ids"`
	Users       types.Set    `tfsdk:"users"`
}

// GroupModel represents a group in the data source.
type GroupModel struct {
	ID          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// UserModel represents a user in the data source.
type UserModel struct {
	ID       types.Int32  `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
}

// Metadata returns the data source type name.
func (d *NotificationGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_group"
}

// Schema defines the schema for the data source.
func (d *NotificationGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a notification group in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the notification group to look up.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the notification group.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the notification group.",
				Computed:            true,
			},
			"group_ids": schema.SetAttribute{
				MarkdownDescription: "Set of user group IDs included in this notification group.",
				Computed:            true,
				ElementType:         types.Int32Type,
			},
			"groups": schema.SetNestedAttribute{
				MarkdownDescription: "The user groups included in this notification group.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							MarkdownDescription: "The unique numeric ID of the group.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the group.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A description of the group.",
							Computed:            true,
						},
					},
				},
			},
			"user_ids": schema.SetAttribute{
				MarkdownDescription: "Set of user IDs included in this notification group.",
				Computed:            true,
				ElementType:         types.Int32Type,
			},
			"users": schema.SetNestedAttribute{
				MarkdownDescription: "The users included in this notification group.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							MarkdownDescription: "The unique numeric ID of the user.",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "The username of the user.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *NotificationGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the data source.
func (d *NotificationGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NotificationGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tflog.Debug(ctx, "Reading notification group", map[string]interface{}{"id": id})

	result, httpResp, err := d.client.ExtrasAPI.ExtrasNotificationGroupsRetrieve(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Notification Group",
			utils.FormatAPIError(fmt.Sprintf("read notification group ID %d", id), err, httpResp))
		return
	}

	d.mapToDataSourceModel(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapToDataSourceModel maps the API response to the data source model.
func (d *NotificationGroupDataSource) mapToDataSourceModel(ctx context.Context, result *netbox.NotificationGroup, data *NotificationGroupDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())

	// Map description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map groups
	if result.HasGroups() {
		groups := result.GetGroups()
		groupIDs := make([]int32, len(groups))
		groupModels := make([]GroupModel, len(groups))
		for i, g := range groups {
			groupIDs[i] = g.GetId()
			groupModels[i] = GroupModel{
				ID:   types.Int32Value(g.GetId()),
				Name: types.StringValue(g.GetName()),
			}
			if g.HasDescription() && g.GetDescription() != "" {
				groupModels[i].Description = types.StringValue(g.GetDescription())
			} else {
				groupModels[i].Description = types.StringNull()
			}
		}
		groupIDsValue, groupIDDiags := types.SetValueFrom(ctx, types.Int32Type, groupIDs)
		diags.Append(groupIDDiags...)
		data.GroupIDs = groupIDsValue

		groupObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":          types.Int32Type,
				"name":        types.StringType,
				"description": types.StringType,
			},
		}
		groupsValue, groupsDiags := types.SetValueFrom(ctx, groupObjectType, groupModels)
		diags.Append(groupsDiags...)
		data.Groups = groupsValue
	} else {
		data.GroupIDs = types.SetNull(types.Int32Type)
		groupObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":          types.Int32Type,
				"name":        types.StringType,
				"description": types.StringType,
			},
		}
		data.Groups = types.SetNull(groupObjectType)
	}

	// Map users
	if result.HasUsers() {
		users := result.GetUsers()
		userIDs := make([]int32, len(users))
		userModels := make([]UserModel, len(users))
		for i, u := range users {
			userIDs[i] = u.GetId()
			userModels[i] = UserModel{
				ID:       types.Int32Value(u.GetId()),
				Username: types.StringValue(u.GetUsername()),
			}
		}
		userIDsValue, userIDDiags := types.SetValueFrom(ctx, types.Int32Type, userIDs)
		diags.Append(userIDDiags...)
		data.UserIDs = userIDsValue

		userObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":       types.Int32Type,
				"username": types.StringType,
			},
		}
		usersValue, usersDiags := types.SetValueFrom(ctx, userObjectType, userModels)
		diags.Append(usersDiags...)
		data.Users = usersValue
	} else {
		data.UserIDs = types.SetNull(types.Int32Type)
		userObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":       types.Int32Type,
				"username": types.StringType,
			},
		}
		data.Users = types.SetNull(userObjectType)
	}
}
