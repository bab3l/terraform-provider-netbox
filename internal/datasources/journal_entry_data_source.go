// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
	AssignedObjectType types.String `tfsdk:"assigned_object_type"`
	AssignedObjectID   types.Int64  `tfsdk:"assigned_object_id"`
	Kind               types.String `tfsdk:"kind"`
	Comments           types.String `tfsdk:"comments"`
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
			},
		},
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
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Journal Entry",
			utils.FormatAPIError("reading journal entry by ID", err, httpResp),
		)
		return
	}

	// Map response to state
	data.ID = types.Int32Value(journalEntry.GetId())
	data.AssignedObjectType = types.StringValue(journalEntry.GetAssignedObjectType())
	data.AssignedObjectID = types.Int64Value(journalEntry.GetAssignedObjectId())
	data.Comments = types.StringValue(journalEntry.GetComments())

	// Kind
	if journalEntry.Kind != nil && journalEntry.Kind.Value != nil {
		data.Kind = types.StringValue(string(*journalEntry.Kind.Value))
	} else {
		data.Kind = types.StringValue("info")
	}

	tflog.Debug(ctx, "Read journal entry", map[string]interface{}{
		"id":                   data.ID.ValueInt32(),
		"assigned_object_type": data.AssignedObjectType.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
