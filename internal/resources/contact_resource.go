package resources

import (
	"context"
	"fmt"

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
			"id":          nbschema.IDAttribute("contact"),
			"name":        nbschema.NameAttribute("contact", 100),
			"group":       nbschema.ReferenceAttribute("contact group", "ID or slug of the contact group this contact belongs to."),
			"description": nbschema.DescriptionAttribute("contact"),
			"comments":    nbschema.CommentsAttribute("contact"),
			"tags":        nbschema.TagsAttribute(),
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

	// Set optional fields
	contactRequest.Title = utils.StringPtr(data.Title)
	contactRequest.Phone = utils.StringPtr(data.Phone)
	contactRequest.Email = utils.StringPtr(data.Email)
	contactRequest.Address = utils.StringPtr(data.Address)
	contactRequest.Link = utils.StringPtr(data.Link)
	contactRequest.Description = utils.StringPtr(data.Description)
	contactRequest.Comments = utils.StringPtr(data.Comments)

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contactRequest.Tags = tags
	}

	contact, httpResp, err := r.client.TenancyAPI.TenancyContactsCreate(ctx).ContactRequest(*contactRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating contact", utils.FormatAPIError("create contact", err, httpResp))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError("Error creating contact", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	if contact == nil {
		resp.Diagnostics.AddError("Contact API returned nil", "No contact object returned from Netbox API.")
		return
	}

	// Map response to state
	r.mapContactToState(ctx, contact, &data)

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
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading contact", utils.FormatAPIError("read contact", err, httpResp))
		return
	}

	r.mapContactToState(ctx, contact, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ContactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Contact ID", fmt.Sprintf("Contact ID must be a number, got: %s", data.ID.ValueString()))
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

	// Set optional fields
	contactRequest.Title = utils.StringPtr(data.Title)
	contactRequest.Phone = utils.StringPtr(data.Phone)
	contactRequest.Email = utils.StringPtr(data.Email)
	contactRequest.Address = utils.StringPtr(data.Address)
	contactRequest.Link = utils.StringPtr(data.Link)
	contactRequest.Description = utils.StringPtr(data.Description)
	contactRequest.Comments = utils.StringPtr(data.Comments)

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contactRequest.Tags = tags
	}

	contact, httpResp, err := r.client.TenancyAPI.TenancyContactsUpdate(ctx, contactID).ContactRequest(*contactRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating contact", utils.FormatAPIError("update contact", err, httpResp))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating contact", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	r.mapContactToState(ctx, contact, &data)

	tflog.Debug(ctx, "Updated contact", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapContactToState maps a Netbox Contact to the Terraform state model.
func (r *ContactResource) mapContactToState(ctx context.Context, contact *netbox.Contact, data *ContactResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", contact.GetId()))
	data.Name = types.StringValue(contact.GetName())

	// Handle optional group
	if contact.HasGroup() && contact.GetGroup().Id != 0 {
		data.Group = types.StringValue(fmt.Sprintf("%d", contact.GetGroup().Id))
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

	// Handle tags
	if contact.HasTags() {
		tags := utils.NestedTagsToTagModels(contact.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		if !tagDiags.HasError() {
			data.Tags = tagsValue
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}
}
