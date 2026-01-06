package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ContactRoleResource{}
var _ resource.ResourceWithImportState = &ContactRoleResource{}

func NewContactRoleResource() resource.Resource {
	return &ContactRoleResource{}
}

type ContactRoleResource struct {
	client *netbox.APIClient
}

// GetClient returns the API client for testing purposes.
func (r *ContactRoleResource) GetClient() *netbox.APIClient {
	return r.client
}

type ContactRoleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *ContactRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact_role"
}

func (r *ContactRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a contact role in Netbox. Contact roles define the function or responsibility of a contact within an organization (e.g., Technical, Administrative, Billing).",
		Attributes: map[string]schema.Attribute{
			"id":          nbschema.IDAttribute("contact role"),
			"name":        nbschema.NameAttribute("contact role", 100),
			"slug":        nbschema.SlugAttribute("contact role"),
			"description": nbschema.DescriptionAttribute("contact role"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("contact role"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *ContactRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ContactRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ContactRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating contact role", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the request
	contactRoleRequest := netbox.ContactRoleRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&contactRoleRequest, data.Description)

	// Handle tags and custom_fields
	utils.ApplyMetadataFields(ctx, &contactRoleRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	// Create via API
	contactRole, httpResp, err := r.client.TenancyAPI.TenancyContactRolesCreate(ctx).ContactRoleRequest(contactRoleRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_contact_role",
			ResourceName: "this.contact_role",
			SlugValue:    data.Slug.ValueString(),
			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				list, _, lookupErr := r.client.TenancyAPI.TenancyContactRolesList(lookupCtx).Slug([]string{slug}).Execute()
				if lookupErr != nil {
					return "", lookupErr
				}
				if list != nil && len(list.Results) > 0 {
					return fmt.Sprintf("%d", list.Results[0].GetId()), nil
				}
				return "", nil
			},
		}
		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating contact role", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	r.mapContactRoleToState(ctx, contactRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a contact role resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ContactRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	contactRoleID := data.ID.ValueString()
	contactRoleIDInt := utils.ParseInt32FromString(contactRoleID)
	if contactRoleIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Contact Role ID", fmt.Sprintf("Contact Role ID must be a number, got: %s", contactRoleID))
		return
	}
	contactRole, httpResp, err := r.client.TenancyAPI.TenancyContactRolesRetrieve(ctx, contactRoleIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading contact role", utils.FormatAPIError(fmt.Sprintf("read contact role ID %s", contactRoleID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error reading contact role", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	r.mapContactRoleToState(ctx, contactRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ContactRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	contactRoleID := data.ID.ValueString()
	contactRoleIDInt := utils.ParseInt32FromString(contactRoleID)
	if contactRoleIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Contact Role ID", fmt.Sprintf("Contact Role ID must be a number, got: %s", contactRoleID))
		return
	}
	tflog.Debug(ctx, "Updating contact role", map[string]interface{}{
		"id":   contactRoleID,
		"name": data.Name.ValueString(),
	})

	// Build the request
	contactRoleRequest := netbox.ContactRoleRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Apply description
	utils.ApplyDescription(&contactRoleRequest, data.Description)

	// Handle tags and custom_fields
	utils.ApplyMetadataFields(ctx, &contactRoleRequest, data.Tags, data.CustomFields, &resp.Diagnostics)

	// Update via API
	contactRole, httpResp, err := r.client.TenancyAPI.TenancyContactRolesUpdate(ctx, contactRoleIDInt).ContactRoleRequest(contactRoleRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating contact role", utils.FormatAPIError(fmt.Sprintf("update contact role ID %s", contactRoleID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating contact role", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	r.mapContactRoleToState(ctx, contactRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ContactRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	contactRoleID := data.ID.ValueString()
	contactRoleIDInt := utils.ParseInt32FromString(contactRoleID)
	if contactRoleIDInt == 0 {
		resp.Diagnostics.AddError("Invalid Contact Role ID", fmt.Sprintf("Contact Role ID must be a number, got: %s", contactRoleID))
		return
	}
	tflog.Debug(ctx, "Deleting contact role", map[string]interface{}{"id": contactRoleID})
	httpResp, err := r.client.TenancyAPI.TenancyContactRolesDestroy(ctx, contactRoleIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting contact role", utils.FormatAPIError(fmt.Sprintf("delete contact role ID %s", contactRoleID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError("Error deleting contact role", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
	tflog.Trace(ctx, "deleted a contact role resource")
}

func (r *ContactRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapContactRoleToState maps API response to Terraform state.
func (r *ContactRoleResource) mapContactRoleToState(ctx context.Context, contactRole *netbox.ContactRole, data *ContactRoleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", contactRole.GetId()))
	data.Name = types.StringValue(contactRole.GetName())
	data.Slug = types.StringValue(contactRole.GetSlug())
	data.Description = utils.StringFromAPI(contactRole.HasDescription(), contactRole.GetDescription, data.Description)

	// Handle tags
	data.Tags = utils.PopulateTagsFromAPI(ctx, contactRole.HasTags(), contactRole.GetTags(), data.Tags, diags)

	// Handle custom fields
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, contactRole.HasCustomFields(), contactRole.GetCustomFields(), data.CustomFields, diags)
}
