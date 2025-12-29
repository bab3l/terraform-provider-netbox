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
)

var _ datasource.DataSource = &ContactDataSource{}

func NewContactDataSource() datasource.DataSource {
	return &ContactDataSource{}
}

// ContactDataSource defines the contact data source implementation.

type ContactDataSource struct {
	client *netbox.APIClient
}

// ContactDataSourceModel describes the contact data source data model.

type ContactDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Group types.String `tfsdk:"group"`

	Title types.String `tfsdk:"title"`

	Phone types.String `tfsdk:"phone"`

	Email types.String `tfsdk:"email"`

	Address types.String `tfsdk:"address"`

	Link types.String `tfsdk:"link"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	DisplayName types.String `tfsdk:"display_name"`
}

func (d *ContactDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (d *ContactDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a contact in Netbox. You can identify the contact using `id`, `name`, or `email`.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.DSIDAttribute("contact"),

			"name": nbschema.DSNameAttribute("contact"),

			"group": nbschema.DSComputedStringAttribute("ID of the contact group this contact belongs to."),

			"title": nbschema.DSComputedStringAttribute("Job title or role of the contact."),

			"phone": nbschema.DSComputedStringAttribute("Phone number of the contact."),

			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of the contact. Use to look up by email.",

				Optional: true,

				Computed: true,
			},

			"address": nbschema.DSComputedStringAttribute("Physical address of the contact."),

			"link": nbschema.DSComputedStringAttribute("URL link associated with the contact."),

			"description": nbschema.DSComputedStringAttribute("Description of the contact."),

			"comments": nbschema.DSComputedStringAttribute("Comments about the contact."),

			"tags": nbschema.DSTagsAttribute(),

			"display_name": nbschema.DSComputedStringAttribute("The display name of the contact."),
		},
	}
}

func (d *ContactDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContactDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var contact *netbox.Contact

	var err error

	var httpResp *http.Response

	// Lookup by id, name, or email

	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():

		contactID, parseErr := utils.ParseID(data.ID.ValueString())

		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid Contact ID", "Contact ID must be a number.")

			return
		}

		contact, httpResp, err = d.client.TenancyAPI.TenancyContactsRetrieve(ctx, contactID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {
			resp.Diagnostics.AddError("Error reading contact", utils.FormatAPIError("read contact", err, httpResp))

			return
		}

	case !data.Name.IsNull() && !data.Name.IsUnknown():

		name := data.Name.ValueString()

		contacts, httpResp, listErr := d.client.TenancyAPI.TenancyContactsList(ctx).Name([]string{name}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if listErr != nil {
			resp.Diagnostics.AddError("Error reading contact", utils.FormatAPIError("read contact by name", listErr, httpResp))

			return
		}

		if contacts == nil || len(contacts.GetResults()) == 0 {
			resp.Diagnostics.AddError("Contact Not Found", fmt.Sprintf("No contact found with name: %s", name))

			return
		}

		contact = &contacts.GetResults()[0]

	case !data.Email.IsNull() && !data.Email.IsUnknown():

		email := data.Email.ValueString()

		contacts, httpResp, listErr := d.client.TenancyAPI.TenancyContactsList(ctx).Email([]string{email}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if listErr != nil {
			resp.Diagnostics.AddError("Error reading contact", utils.FormatAPIError("read contact by email", listErr, httpResp))

			return
		}

		if contacts == nil || len(contacts.GetResults()) == 0 {
			resp.Diagnostics.AddError("Contact Not Found", fmt.Sprintf("No contact found with email: %s", email))

			return
		}

		contact = &contacts.GetResults()[0]

	default:

		resp.Diagnostics.AddError("Missing Contact Identifier", "Either 'id', 'name', or 'email' must be specified.")

		return
	}

	if contact == nil {
		resp.Diagnostics.AddError("Contact Not Found", "No contact found with the specified identifier.")

		return
	}

	// Map API response to state

	d.mapContactToState(ctx, contact, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapContactToState maps API response to Terraform state.

func (d *ContactDataSource) mapContactToState(ctx context.Context, contact *netbox.Contact, data *ContactDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", contact.GetId()))

	data.Name = types.StringValue(contact.GetName())

	// Handle optional group

	if group, ok := contact.GetGroupOk(); ok && group != nil && group.Id != 0 {
		data.Group = types.StringValue(group.GetName())
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

	// Map display name
	if contact.GetDisplay() != "" {
		data.DisplayName = types.StringValue(contact.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
