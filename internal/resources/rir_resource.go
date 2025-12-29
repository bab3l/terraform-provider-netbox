// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &RIRResource{}
	_ resource.ResourceWithConfigure   = &RIRResource{}
	_ resource.ResourceWithImportState = &RIRResource{}
)

// NewRIRResource returns a new RIR resource.
func NewRIRResource() resource.Resource {
	return &RIRResource{}
}

// RIRResource defines the resource implementation.
type RIRResource struct {
	client *netbox.APIClient
}

// RIRResourceModel describes the resource data model.
type RIRResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	IsPrivate    types.Bool   `tfsdk:"is_private"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *RIRResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rir"
}

// Schema defines the schema for the resource.
func (r *RIRResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Regional Internet Registry (RIR) in Netbox. RIRs are organizations that manage the allocation and registration of Internet number resources (IP addresses, ASNs).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the RIR.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":         nbschema.NameAttribute("RIR", 100),
			"slug":         nbschema.SlugAttribute("RIR"),
			"is_private": schema.BoolAttribute{
				MarkdownDescription: "Whether IP space managed by this RIR is considered private. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("RIR"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.
func (r *RIRResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *RIRResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RIRResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the RIR request
	rirRequest := netbox.NewRIRRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, rirRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating RIR", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Create the RIR
	rir, httpResp, err := r.client.IpamAPI.IpamRirsCreate(ctx).RIRRequest(*rirRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating RIR",
			utils.FormatAPIError("create RIR", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *RIRResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RIRResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Reading RIR", map[string]interface{}{
		"id": id,
	})

	// Get the RIR from Netbox
	rir, httpResp, err := r.client.IpamAPI.IpamRirsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading RIR",
			utils.FormatAPIError(fmt.Sprintf("read RIR ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Read RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *RIRResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RIRResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Create the RIR request
	rirRequest := netbox.NewRIRRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields
	r.setOptionalFields(ctx, rirRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating RIR", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Update the RIR
	rir, httpResp, err := r.client.IpamAPI.IpamRirsUpdate(ctx, id).RIRRequest(*rirRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating RIR",
			utils.FormatAPIError(fmt.Sprintf("update RIR ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapRIRToState(ctx, rir, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *RIRResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RIRResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting RIR", map[string]interface{}{
		"id": id,
	})

	// Delete the RIR
	httpResp, err := r.client.IpamAPI.IpamRirsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting RIR",
			utils.FormatAPIError(fmt.Sprintf("delete RIR ID %d", id), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted RIR", map[string]interface{}{
		"id": id,
	})
}

func (r *RIRResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the RIR request from the resource model.
func (r *RIRResource) setOptionalFields(ctx context.Context, rirRequest *netbox.RIRRequest, data *RIRResourceModel, diags *diag.Diagnostics) {
	// Is Private
	if utils.IsSet(data.IsPrivate) {
		isPrivate := data.IsPrivate.ValueBool()
		rirRequest.IsPrivate = &isPrivate
	}

	// Apply description and metadata fields
	utils.ApplyDescription(rirRequest, data.Description)
	utils.ApplyMetadataFields(ctx, rirRequest, data.Tags, data.CustomFields, diags)
}

// mapRIRToState maps a Netbox RIR to the Terraform state model.
func (r *RIRResource) mapRIRToState(ctx context.Context, rir *netbox.RIR, data *RIRResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rir.Id))
	data.Name = types.StringValue(rir.Name)
	data.Slug = types.StringValue(rir.Slug)

	// Is Private
	if rir.IsPrivate != nil {
		data.IsPrivate = types.BoolValue(*rir.IsPrivate)
	} else {
		data.IsPrivate = types.BoolValue(false)
	}

	// Description
	if rir.Description != nil && *rir.Description != "" {
		data.Description = types.StringValue(*rir.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags
	data.Tags = utils.PopulateTagsFromNestedTags(ctx, len(rir.Tags) > 0, rir.Tags, diags)

	// Custom Fields
	data.CustomFields = utils.PopulateCustomFieldsFromMap(ctx, len(rir.CustomFields) > 0, rir.CustomFields, data.CustomFields, diags)
}
