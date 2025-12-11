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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	client *netbox.APIClient
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Username  types.String `tfsdk:"username"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Email     types.String `tfsdk:"email"`
	IsStaff   types.Bool   `tfsdk:"is_staff"`
	IsActive  types.Bool   `tfsdk:"is_active"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a user in Netbox. You can identify the user using `id` or `username`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique ID of the user. Use this or `username` to look up the user.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username of the user. Use this or `id` to look up the user.",
				Optional:            true,
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name of the user.",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name of the user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of the user.",
				Computed:            true,
			},
			"is_staff": schema.BoolAttribute{
				MarkdownDescription: "Whether the user can log into the admin site.",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the user account is active.",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var user *netbox.User
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID or username
	if !data.ID.IsNull() {
		// Search by ID
		userID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading user by ID", map[string]interface{}{
			"id": userID,
		})

		// Parse the user ID to int32 for the API call
		var userIDInt int32
		if _, parseErr := fmt.Sscanf(userID, "%d", &userIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid User ID",
				fmt.Sprintf("User ID must be a number, got: %s", userID),
			)
			return
		}

		// Retrieve the user via API
		user, httpResp, err = d.client.UsersAPI.UsersUsersRetrieve(ctx, userIDInt).Execute()
	} else if !data.Username.IsNull() {
		// Search by username
		username := data.Username.ValueString()
		tflog.Debug(ctx, "Reading user by username", map[string]interface{}{
			"username": username,
		})

		// List users with username filter
		var users *netbox.PaginatedUserList
		users, httpResp, err = d.client.UsersAPI.UsersUsersList(ctx).Username([]string{username}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading user",
				utils.FormatAPIError("read user by username", err, httpResp),
			)
			return
		}
		if len(users.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"User Not Found",
				fmt.Sprintf("No user found with username: %s", username),
			)
			return
		}
		if len(users.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Users Found",
				fmt.Sprintf("Multiple users found with username: %s. This should not happen as usernames should be unique.", username),
			)
			return
		}
		user = &users.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Identifier",
			"Either 'id' or 'username' must be specified to look up a user.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			utils.FormatAPIError("read user", err, httpResp),
		)
		return
	}

	// Map response to model
	data.ID = types.StringValue(fmt.Sprintf("%d", user.GetId()))
	data.Username = types.StringValue(user.GetUsername())

	if user.HasFirstName() && user.GetFirstName() != "" {
		data.FirstName = types.StringValue(user.GetFirstName())
	} else {
		data.FirstName = types.StringNull()
	}

	if user.HasLastName() && user.GetLastName() != "" {
		data.LastName = types.StringValue(user.GetLastName())
	} else {
		data.LastName = types.StringNull()
	}

	if user.HasEmail() && user.GetEmail() != "" {
		data.Email = types.StringValue(user.GetEmail())
	} else {
		data.Email = types.StringNull()
	}

	if user.HasIsStaff() {
		data.IsStaff = types.BoolValue(user.GetIsStaff())
	} else {
		data.IsStaff = types.BoolNull()
	}

	if user.HasIsActive() {
		data.IsActive = types.BoolValue(user.GetIsActive())
	} else {
		data.IsActive = types.BoolNull()
	}

	tflog.Trace(ctx, "read user data source", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"username": data.Username.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
