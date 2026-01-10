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
var _ datasource.DataSource = &JournalEntryDataSource{}

func NewJournalEntryDataSource() datasource.DataSource {
	return &JournalEntryDataSource{}
}

// JournalEntryDataSource defines the data source implementation.
type JournalEntryDataSource struct {
	client *netbox.APIClient
}

// JournalEntryDataSourceModel describes the data source data model.
type JournalEntryDataSourceModel struct {
	ID                 types.Int32  `tfsdk:"id"`
	DisplayName        types.String `tfsdk:"display_name"`
	AssignedObjectType types.String `tfsdk:"assigned_object_type"`
	AssignedObjectID   types.Int64  `tfsdk:"assigned_object_id"`
	Kind               types.String `tfsdk:"kind"`
	Comments           types.String `tfsdk:"comments"`
	CustomFields       types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *JournalEntryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_journal_entry"
}

// Schema defines the schema for the data source.
func (d *JournalEntryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a journal entry in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the journal entry. Required for lookup.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the journal entry.",
				Computed:            true,
			},
			"assigned_object_type": schema.StringAttribute{
				MarkdownDescription: "The content type of the assigned object (e.g., `dcim.device`, `dcim.site`, `ipam.ipaddress`).",
				Computed:            true,
			},
			"assigned_object_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the assigned object.",
				Computed:            true,
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "The kind/severity of the journal entry (info, success, warning, danger).",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "The content of the journal entry.",
				Computed:            true,
			}, "custom_fields": nbschema.DSCustomFieldsAttribute()},
	}
}

// Configure adds the provider configured client to the data source.
func (d *JournalEntryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves the journal entry data from NetBox.
func (d *JournalEntryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JournalEntryDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueInt32()
	tflog.Debug(ctx, "Looking up journal entry by ID", map[string]interface{}{
		"id": id,
	})
	journalEntry, httpResp, err := d.client.ExtrasAPI.ExtrasJournalEntriesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Journal Entry Not Found",
			fmt.Sprintf("No journal entry found with ID: %d", id),
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Journal Entry",
			utils.FormatAPIError("reading journal entry by ID", err, httpResp),
		)
		return
	}

	// Map response to state
	data.ID = types.Int32Value(journalEntry.GetId())

	// Display Name
	if journalEntry.GetDisplay() != "" {
		data.DisplayName = types.StringValue(journalEntry.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
	data.AssignedObjectType = types.StringValue(journalEntry.GetAssignedObjectType())
	data.AssignedObjectID = types.Int64Value(journalEntry.GetAssignedObjectId())
	data.Comments = types.StringValue(journalEntry.GetComments())

	// Kind
	if journalEntry.Kind != nil && journalEntry.Kind.Value != nil {
		data.Kind = types.StringValue(string(*journalEntry.Kind.Value))
	} else {
		data.Kind = types.StringValue("info")
	}

	// Map custom fields
	if journalEntry.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(journalEntry.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
		resp.Diagnostics.Append(cfDiags...)
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Debug(ctx, "Read journal entry", map[string]interface{}{
		"id":                   data.ID.ValueInt32(),
		"assigned_object_type": data.AssignedObjectType.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
