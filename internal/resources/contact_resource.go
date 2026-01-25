package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ContactResource{}
var _ resource.ResourceWithImportState = &ContactResource{}

func NewContactResource() resource.Resource {
	return &ContactResource{}
}

// ContactResource defines the contact resource implementation.
type ContactResource struct {
	client *netbox.APIClient
}

// ContactResourceModel describes the contact resource data model.
type ContactResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Group       types.String `tfsdk:"group"`
	Title       types.String `tfsdk:"title"`
	Phone       types.String `tfsdk:"phone"`
	Email       types.String `tfsdk:"email"`
	Address     types.String `tfsdk:"address"`
	Link        types.String `tfsdk:"link"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	Tags        types.Set    `tfsdk:"tags"`
}

func (r *ContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (r *ContactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a contact in Netbox. Contacts represent people or organizations that can be assigned to various resources.",
		Attributes: map[string]schema.Attribute{
			"id":    nbschema.IDAttribute("contact"),
			"name":  nbschema.NameAttribute("contact", 100),
			"group": nbschema.ReferenceAttributeWithDiffSuppress("contact group", "ID or slug of the contact group this contact belongs to."),

			"tags": nbschema.TagsSlugAttribute(),
			"title": schema.StringAttribute{
				MarkdownDescription: "Job title or role of the contact.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"phone": schema.StringAttribute{
				MarkdownDescription: "Phone number of the contact.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of the contact.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(254),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Physical address of the contact.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"link": schema.StringAttribute{
				MarkdownDescription: "URL link associated with the contact (e.g., a personal website or profile page).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("contact"))

	// Note: This resource does not have custom_fields
}

func (r *ContactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ContactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	contactRequest := netbox.NewContactRequest(data.Name.ValueString())

	// Set optional group reference
	if !data.Group.IsNull() && !data.Group.IsUnknown() {
		group, diags := netboxlookup.LookupContactGroup(ctx, r.client, data.Group.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contactRequest.SetGroup(*group)
	}

	// Set optional fields with explicit null clearing
	if utils.IsSet(data.Title) {
		title := data.Title.ValueString()
		contactRequest.Title = &title
	} else if data.Title.IsNull() {
		empty := ""
		contactRequest.Title = &empty
	}

	if utils.IsSet(data.Phone) {
		phone := data.Phone.ValueString()
		contactRequest.Phone = &phone
	} else if data.Phone.IsNull() {
		empty := ""
		contactRequest.Phone = &empty
	}

	if utils.IsSet(data.Email) {
		email := data.Email.ValueString()
		contactRequest.Email = &email
	} else if data.Email.IsNull() {
		empty := ""
		contactRequest.Email = &empty
	}

	if utils.IsSet(data.Address) {
		address := data.Address.ValueString()
		contactRequest.Address = &address
	} else if data.Address.IsNull() {
		empty := ""
		contactRequest.Address = &empty
	}

	if utils.IsSet(data.Link) {
		link := data.Link.ValueString()
		contactRequest.Link = &link
	} else if data.Link.IsNull() {
		empty := ""
		contactRequest.Link = &empty
	}

	// Apply description and comments
	utils.ApplyDescription(contactRequest, data.Description)
	utils.ApplyComments(contactRequest, data.Comments)

	// Store plan values for filter-to-owned pattern
	planTags := data.Tags

	// Handle tags
	utils.ApplyTagsFromSlugs(ctx, r.client, contactRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	contact, httpResp, err := r.client.TenancyAPI.TenancyContactsCreate(ctx).ContactRequest(*contactRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating contact", utils.FormatAPIError("create contact", err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating contact", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	if contact == nil {
		resp.Diagnostics.AddError("Contact API returned nil", "No contact object returned from Netbox API.")
		return
	}

	// Map response to state
	r.mapContactToState(contact, &data)

	// Apply filter-to-owned pattern for tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, contact.HasTags(), contact.GetTags(), planTags)
	tflog.Debug(ctx, "Created contact", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Contact ID", fmt.Sprintf("Contact ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	contact, httpResp, err := r.client.TenancyAPI.TenancyContactsRetrieve(ctx, contactID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading contact", utils.FormatAPIError("read contact", err, httpResp))
		return
	}
	// Store state tags before mapping
	stateTags := data.Tags
	r.mapContactToState(contact, &data)
	// Apply filter-to-owned pattern for tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, contact.HasTags(), contact.GetTags(), stateTags)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Contact ID", fmt.Sprintf("Contact ID must be a number, got: %s", plan.ID.ValueString()))
		return
	}
	contactRequest := netbox.NewContactRequest(plan.Name.ValueString())

	// Set optional group reference
	if !plan.Group.IsNull() && !plan.Group.IsUnknown() {
		group, diags := netboxlookup.LookupContactGroup(ctx, r.client, plan.Group.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contactRequest.SetGroup(*group)
	} else if plan.Group.IsNull() {
		contactRequest.SetGroupNil()
	}

	// Set optional fields with explicit null clearing
	if utils.IsSet(plan.Title) {
		title := plan.Title.ValueString()
		contactRequest.Title = &title
	} else if plan.Title.IsNull() {
		empty := ""
		contactRequest.Title = &empty
	}

	if utils.IsSet(plan.Phone) {
		phone := plan.Phone.ValueString()
		contactRequest.Phone = &phone
	} else if plan.Phone.IsNull() {
		empty := ""
		contactRequest.Phone = &empty
	}

	if utils.IsSet(plan.Email) {
		email := plan.Email.ValueString()
		contactRequest.Email = &email
	} else if plan.Email.IsNull() {
		empty := ""
		contactRequest.Email = &empty
	}

	if utils.IsSet(plan.Address) {
		address := plan.Address.ValueString()
		contactRequest.Address = &address
	} else if plan.Address.IsNull() {
		empty := ""
		contactRequest.Address = &empty
	}

	if utils.IsSet(plan.Link) {
		link := plan.Link.ValueString()
		contactRequest.Link = &link
	} else if plan.Link.IsNull() {
		empty := ""
		contactRequest.Link = &empty
	}

	// Apply description and comments
	utils.ApplyDescription(contactRequest, plan.Description)
	utils.ApplyComments(contactRequest, plan.Comments)

	// Handle tags (tags use replace-all semantics)
	utils.ApplyTagsFromSlugs(ctx, r.client, contactRequest, plan.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	contact, httpResp, err := r.client.TenancyAPI.TenancyContactsUpdate(ctx, contactID).ContactRequest(*contactRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating contact", utils.FormatAPIError("update contact", err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating contact", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	r.mapContactToState(contact, &plan)

	// After update, populate tags based on what the user specified
	// If tags were null in plan (not specified), keep them null to match config
	// Otherwise, populate from API response (replace-all semantics)
	plan.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, contact.HasTags(), contact.GetTags(), plan.Tags)

	tflog.Debug(ctx, "Updated contact", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Contact ID", fmt.Sprintf("Contact ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	httpResp, err := r.client.TenancyAPI.TenancyContactsDestroy(ctx, contactID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Already deleted, nothing to do
			return
		}
		resp.Diagnostics.AddError("Error deleting contact", utils.FormatAPIError("delete contact", err, httpResp))
		return
	}
	tflog.Debug(ctx, "Deleted contact", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
}

func (r *ContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapContactToState maps a Netbox Contact to the Terraform state model.
func (r *ContactResource) mapContactToState(contact *netbox.Contact, data *ContactResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", contact.GetId()))
	data.Name = types.StringValue(contact.GetName())

	// Handle optional group - preserve user's input format
	if contact.HasGroup() && contact.GetGroup().Id != 0 {
		group := contact.GetGroup()
		data.Group = utils.UpdateReferenceAttribute(data.Group, group.GetName(), group.GetSlug(), group.Id)
	} else {
		data.Group = types.StringNull()
	}

	// Map optional string fields
	if title := contact.GetTitle(); title != "" {
		data.Title = types.StringValue(title)
	} else {
		data.Title = types.StringNull()
	}

	if phone := contact.GetPhone(); phone != "" {
		data.Phone = types.StringValue(phone)
	} else {
		data.Phone = types.StringNull()
	}

	if email := contact.GetEmail(); email != "" {
		data.Email = types.StringValue(email)
	} else {
		data.Email = types.StringNull()
	}

	if address := contact.GetAddress(); address != "" {
		data.Address = types.StringValue(address)
	} else {
		data.Address = types.StringNull()
	}

	if link := contact.GetLink(); link != "" {
		data.Link = types.StringValue(link)
	} else {
		data.Link = types.StringNull()
	}

	if desc := contact.GetDescription(); desc != "" {
		data.Description = types.StringValue(desc)
	} else {
		data.Description = types.StringNull()
	}

	if comments := contact.GetComments(); comments != "" {
		data.Comments = types.StringValue(comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags are handled in Create/Read/Update with filter-to-owned pattern.
}
