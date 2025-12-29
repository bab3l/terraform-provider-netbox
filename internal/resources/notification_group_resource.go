// Package resources provides Terraform resource implementations for NetBox objects.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationGroupResource{}
	_ resource.ResourceWithConfigure   = &NotificationGroupResource{}
	_ resource.ResourceWithImportState = &NotificationGroupResource{}
)

// NewNotificationGroupResource returns a new resource implementing the notification group resource.
func NewNotificationGroupResource() resource.Resource {
	return &NotificationGroupResource{}
}

// NotificationGroupResource defines the resource implementation.
type NotificationGroupResource struct {
	client *netbox.APIClient
}

// NotificationGroupResourceModel describes the resource data model.
type NotificationGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	GroupIDs    types.Set    `tfsdk:"group_ids"`
	UserIDs     types.Set    `tfsdk:"user_ids"`
	DisplayName types.String `tfsdk:"display_name"`
}

// Metadata returns the resource type name.
func (r *NotificationGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_group"
}

// Schema defines the schema for the resource.
func (r *NotificationGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a notification group in NetBox. Notification groups are used to define sets of users and groups that receive notifications from event rules.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the notification group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":        nbschema.NameAttribute("notification group", 100),
			"description": nbschema.DescriptionAttribute("notification group"),
			"group_ids": schema.SetAttribute{
				MarkdownDescription: "Set of user group IDs to include in this notification group.",
				Optional:            true,
				ElementType:         types.Int32Type,
			},
			"user_ids": schema.SetAttribute{
				MarkdownDescription: "Set of user IDs to include in this notification group.",
				Optional:            true,
				ElementType:         types.Int32Type,
			},
			"display_name": nbschema.DisplayNameAttribute("notification group"),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NotificationGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource.
func (r *NotificationGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotificationGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating notification group", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Build the request
	request := netbox.NewNotificationGroupRequest(data.Name.ValueString())

	// Set optional fields
	utils.ApplyDescription(request, data.Description)

	if !data.GroupIDs.IsNull() && !data.GroupIDs.IsUnknown() {
		var groupIDs []int32
		resp.Diagnostics.Append(data.GroupIDs.ElementsAs(ctx, &groupIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Groups = groupIDs
	}

	if !data.UserIDs.IsNull() && !data.UserIDs.IsUnknown() {
		var userIDs []int32
		resp.Diagnostics.Append(data.UserIDs.ElementsAs(ctx, &userIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Users = userIDs
	}

	// Create the notification group
	result, httpResp, err := r.client.ExtrasAPI.ExtrasNotificationGroupsCreate(ctx).
		NotificationGroupRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Notification Group",
			utils.FormatAPIError("create notification group", err, httpResp))
		return
	}

	// Map the response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created notification group", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *NotificationGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NotificationGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	result, httpResp, err := r.client.ExtrasAPI.ExtrasNotificationGroupsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Notification group not found, removing from state", map[string]interface{}{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Notification Group",
			utils.FormatAPIError(fmt.Sprintf("read notification group ID %d", id), err, httpResp))
		return
	}

	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *NotificationGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NotificationGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tflog.Debug(ctx, "Updating notification group", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Build the request
	request := netbox.NewNotificationGroupRequest(data.Name.ValueString())

	// Set optional fields (same as Create)
	utils.ApplyDescription(request, data.Description)

	if !data.GroupIDs.IsNull() && !data.GroupIDs.IsUnknown() {
		var groupIDs []int32
		resp.Diagnostics.Append(data.GroupIDs.ElementsAs(ctx, &groupIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Groups = groupIDs
	}

	if !data.UserIDs.IsNull() && !data.UserIDs.IsUnknown() {
		var userIDs []int32
		resp.Diagnostics.Append(data.UserIDs.ElementsAs(ctx, &userIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Users = userIDs
	}

	// Update the notification group
	result, httpResp, err := r.client.ExtrasAPI.ExtrasNotificationGroupsUpdate(ctx, id).
		NotificationGroupRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Notification Group",
			utils.FormatAPIError(fmt.Sprintf("update notification group ID %d", id), err, httpResp))
		return
	}

	// Map the response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.
func (r *NotificationGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NotificationGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tflog.Debug(ctx, "Deleting notification group", map[string]interface{}{"id": id})

	httpResp, err := r.client.ExtrasAPI.ExtrasNotificationGroupsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Notification Group",
			utils.FormatAPIError(fmt.Sprintf("delete notification group ID %d", id), err, httpResp))
		return
	}
}

// ImportState imports the resource state.
func (r *NotificationGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapToState maps the API response to the Terraform state.
func (r *NotificationGroupResource) mapToState(ctx context.Context, result *netbox.NotificationGroup, data *NotificationGroupResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())

	// Map description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map group IDs - extract IDs from Group objects
	// Only set if there are actual groups AND the original data had group_ids set
	if result.HasGroups() && len(result.GetGroups()) > 0 {
		groups := result.GetGroups()
		groupIDs := make([]int32, len(groups))
		for i, g := range groups {
			groupIDs[i] = g.GetId()
		}
		groupIDsValue, groupDiags := types.SetValueFrom(ctx, types.Int32Type, groupIDs)
		diags.Append(groupDiags...)
		data.GroupIDs = groupIDsValue
	} else {
		// No groups - set to null
		data.GroupIDs = types.SetNull(types.Int32Type)
	}

	// Map user IDs - extract IDs from User objects
	// Only set if there are actual users AND the original data had user_ids set
	if result.HasUsers() && len(result.GetUsers()) > 0 {
		users := result.GetUsers()
		userIDs := make([]int32, len(users))
		for i, u := range users {
			userIDs[i] = u.GetId()
		}
		userIDsValue, userDiags := types.SetValueFrom(ctx, types.Int32Type, userIDs)
		diags.Append(userDiags...)
		data.UserIDs = userIDsValue
	} else {
		// No users - set to null
		data.UserIDs = types.SetNull(types.Int32Type)
	}
	// Map display_name
	if result.Display != "" {
		data.DisplayName = types.StringValue(result.Display)
	} else {
		data.DisplayName = types.StringNull()
	}
}
