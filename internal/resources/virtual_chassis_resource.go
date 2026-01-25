// Package resources contains Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &VirtualChassisResource{}

	_ resource.ResourceWithConfigure = &VirtualChassisResource{}

	_ resource.ResourceWithImportState = &VirtualChassisResource{}

	_ resource.ResourceWithIdentity = &VirtualChassisResource{}
)

// NewVirtualChassisResource returns a new resource implementing the VirtualChassis resource.

func NewVirtualChassisResource() resource.Resource {
	return &VirtualChassisResource{}
}

// VirtualChassisResource defines the resource implementation.

type VirtualChassisResource struct {
	client *netbox.APIClient
}

// VirtualChassisResourceModel describes the resource data model.

type VirtualChassisResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Domain types.String `tfsdk:"domain"`

	Master types.String `tfsdk:"master"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	MemberCount types.Int64 `tfsdk:"member_count"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *VirtualChassisResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_chassis"
}

// Schema defines the schema for the resource.

func (r *VirtualChassisResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual chassis in NetBox. A virtual chassis represents a set of devices that are physically stacked or clustered together and managed as a single logical device.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("virtual chassis"),

			"name": nbschema.NameAttribute("virtual chassis", 64),

			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain for this virtual chassis.",

				Optional: true,
			},

			"master": schema.StringAttribute{
				MarkdownDescription: "ID of the master device for this virtual chassis.",

				Optional: true,
			},

			"member_count": schema.Int64Attribute{
				MarkdownDescription: "Number of member devices in this virtual chassis.",

				Computed: true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("virtual chassis"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *VirtualChassisResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.

func (r *VirtualChassisResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new virtual chassis resource.

func (r *VirtualChassisResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtualChassisResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request

	vcRequest, diags := r.buildRequest(ctx, &data, nil)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating virtual chassis", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	vc, httpResp, err := r.client.DcimAPI.DcimVirtualChassisCreate(ctx).WritableVirtualChassisRequest(*vcRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating virtual chassis",

			utils.FormatAPIError("create virtual chassis", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToModel(ctx, vc, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created virtual chassis", map[string]interface{}{
		"id": vc.GetId(),

		"name": vc.GetName(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the virtual chassis resource.

func (r *VirtualChassisResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtualChassisResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vcID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Virtual Chassis ID",

			fmt.Sprintf("Could not parse virtual chassis ID: %s", err),
		)

		return
	}

	tflog.Debug(ctx, "Reading virtual chassis", map[string]interface{}{
		"id": vcID,
	})

	vc, httpResp, err := r.client.DcimAPI.DcimVirtualChassisRetrieve(ctx, vcID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading virtual chassis",

			utils.FormatAPIError(fmt.Sprintf("read virtual chassis ID %d", vcID), err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToModel(ctx, vc, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the virtual chassis resource.

func (r *VirtualChassisResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan VirtualChassisResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vcID, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Virtual Chassis ID",

			fmt.Sprintf("Could not parse virtual chassis ID: %s", err),
		)

		return
	}

	// Build the request

	vcRequest, diags := r.buildRequest(ctx, &plan, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating virtual chassis", map[string]interface{}{
		"id": vcID,
	})

	vc, httpResp, err := r.client.DcimAPI.DcimVirtualChassisUpdate(ctx, vcID).WritableVirtualChassisRequest(*vcRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating virtual chassis",

			utils.FormatAPIError(fmt.Sprintf("update virtual chassis ID %d", vcID), err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToModel(ctx, vc, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the virtual chassis resource.

func (r *VirtualChassisResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtualChassisResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vcID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Virtual Chassis ID",

			fmt.Sprintf("Could not parse virtual chassis ID: %s", err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting virtual chassis", map[string]interface{}{
		"id": vcID,
	})

	httpResp, err := r.client.DcimAPI.DcimVirtualChassisDestroy(ctx, vcID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting virtual chassis",

			utils.FormatAPIError(fmt.Sprintf("delete virtual chassis ID %d", vcID), err, httpResp),
		)

		return
	}
}

// ImportState imports an existing virtual chassis resource.

func (r *VirtualChassisResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		vcID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Virtual Chassis ID", fmt.Sprintf("Could not parse virtual chassis ID: %s", parsed.ID))
			return
		}

		vc, httpResp, err := r.client.DcimAPI.DcimVirtualChassisRetrieve(ctx, vcID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing virtual chassis", utils.FormatAPIError(fmt.Sprintf("read virtual chassis ID %d", vcID), err, httpResp))
			return
		}

		var data VirtualChassisResourceModel
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapResponseToModel(ctx, vc, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, vc.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// buildRequest builds the API request from the Terraform model.

func (r *VirtualChassisResource) buildRequest(ctx context.Context, plan *VirtualChassisResourceModel, state *VirtualChassisResourceModel) (*netbox.WritableVirtualChassisRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	vcRequest := netbox.NewWritableVirtualChassisRequest(plan.Name.ValueString())

	// Set optional fields

	if !plan.Domain.IsNull() && !plan.Domain.IsUnknown() {
		vcRequest.SetDomain(plan.Domain.ValueString())
	} else if !plan.Domain.IsUnknown() {
		// Explicitly set empty string to clear the field
		vcRequest.SetDomain("")
	}

	if !plan.Master.IsNull() && !plan.Master.IsUnknown() {
		masterID, err := utils.ParseID(plan.Master.ValueString())

		if err != nil {
			diags.AddError(

				"Invalid Master Device ID",

				fmt.Sprintf("Could not parse master device ID: %s", err),
			)

			return nil, diags
		}

		vcRequest.SetMaster(masterID)
	}

	// Set common fields with merge-aware helpers
	utils.ApplyDescription(vcRequest, plan.Description)
	utils.ApplyComments(vcRequest, plan.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, vcRequest, plan.Tags, &diags)
	// For custom fields, use state if available (Update), otherwise nil (Create)
	var stateCustomFields types.Set
	if state != nil {
		stateCustomFields = state.CustomFields
	}
	utils.ApplyCustomFieldsWithMerge(ctx, vcRequest, plan.CustomFields, stateCustomFields, &diags)
	if diags.HasError() {
		return nil, diags
	}

	return vcRequest, diags
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *VirtualChassisResource) mapResponseToModel(ctx context.Context, vc *netbox.VirtualChassis, data *VirtualChassisResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", vc.GetId()))

	data.Name = types.StringValue(vc.GetName())

	// Map domain

	if domain, ok := vc.GetDomainOk(); ok && domain != nil && *domain != "" {
		data.Domain = types.StringValue(*domain)
	} else {
		data.Domain = types.StringNull()
	}

	// Map master

	if vc.Master.IsSet() && vc.Master.Get() != nil {
		master := vc.Master.Get()

		userMaster := data.Master.ValueString()

		if userMaster == master.GetName() || userMaster == master.GetDisplay() || userMaster == fmt.Sprintf("%d", master.GetId()) {
			// Keep user's original value
		} else {
			data.Master = types.StringValue(master.GetName())
		}
	} else {
		data.Master = types.StringNull()
	}

	// Map description

	if desc, ok := vc.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments

	if comments, ok := vc.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Map member_count

	data.MemberCount = types.Int64Value(int64(vc.GetMemberCount()))

	// Handle tags (slug list) with empty-set preservation
	data.Tags = utils.PopulateTagsSlugFromAPI(ctx, vc.HasTags(), vc.GetTags(), data.Tags)

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, vc.GetCustomFields(), diags)
}
