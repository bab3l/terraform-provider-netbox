// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"regexp"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &CircuitResource{}
	_ resource.ResourceWithConfigure   = &CircuitResource{}
	_ resource.ResourceWithImportState = &CircuitResource{}
	_ resource.ResourceWithIdentity    = &CircuitResource{}
)

// NewCircuitResource returns a new circuit resource.
func NewCircuitResource() resource.Resource {
	return &CircuitResource{}
}

// CircuitResource defines the circuit resource implementation.
type CircuitResource struct {
	client *netbox.APIClient
}

// CircuitResourceModel describes the circuit resource data model.
type CircuitResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Cid             types.String `tfsdk:"cid"`
	CircuitProvider types.String `tfsdk:"circuit_provider"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	Tenant          types.String `tfsdk:"tenant"`
	InstallDate     types.String `tfsdk:"install_date"`
	TerminationDate types.String `tfsdk:"termination_date"`
	CommitRate      types.Int64  `tfsdk:"commit_rate"`
	Description     types.String `tfsdk:"description"`
	Comments        types.String `tfsdk:"comments"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *CircuitResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit"
}

// Schema defines the schema for the resource.
func (r *CircuitResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a circuit in Netbox. Circuits represent physical or logical network connections provided by external carriers or service providers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the circuit.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cid": schema.StringAttribute{
				MarkdownDescription: "The unique circuit ID assigned by the provider. This is typically a service order number or circuit identifier from the carrier.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"circuit_provider": schema.StringAttribute{
				MarkdownDescription: "The circuit provider (carrier or ISP) supplying this circuit. Can be specified by name, slug, or ID.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of circuit (e.g., Internet Transit, MPLS, Point-to-Point). Can be specified by name, slug, or ID.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The operational status of the circuit. Valid values are: `planned`, `provisioning`, `active`, `offline`, `deprovisioning`, `decommissioned`. Defaults to `active`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("planned", "provisioning", "active", "offline", "deprovisioning", "decommissioned"),
				},
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The tenant that owns this circuit. Can be specified by name, slug, or ID.",
				Optional:            true,
			},
			"install_date": schema.StringAttribute{
				MarkdownDescription: "The date when the circuit was installed, in YYYY-MM-DD format.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
						"must be in YYYY-MM-DD format",
					),
				},
			},
			"termination_date": schema.StringAttribute{
				MarkdownDescription: "The date when the circuit will be or was terminated, in YYYY-MM-DD format.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
						"must be in YYYY-MM-DD format",
					),
				},
			},
			"commit_rate": schema.Int64Attribute{
				MarkdownDescription: "The committed information rate (CIR) in Kbps for this circuit.",
				Optional:            true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("circuit"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *CircuitResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure sets up the resource with the provider client.
func (r *CircuitResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new circuit resource.
func (r *CircuitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CircuitResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request (pass nil state since this is Create)
	createReq, diags := r.buildCircuitRequest(ctx, &data, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating circuit", map[string]interface{}{
		"cid":      data.Cid.ValueString(),
		"provider": data.CircuitProvider.ValueString(),
		"type":     data.Type.ValueString(),
	})

	// Create the circuit
	circuit, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitsCreate(ctx).WritableCircuitRequest(*createReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating circuit",
			utils.FormatAPIError("create circuit", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Created circuit", map[string]interface{}{
		"id":  circuit.GetId(),
		"cid": circuit.GetCid(),
	})

	// Map the response to state
	r.mapCircuitToState(ctx, circuit, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the circuit resource.
func (r *CircuitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CircuitResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse circuit ID: %s", err))
		return
	}
	tflog.Debug(ctx, "Reading circuit", map[string]interface{}{
		"id": id,
	})
	circuit, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Circuit not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading circuit",
			utils.FormatAPIError("read circuit", err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state if null or empty
	originalCustomFields := data.CustomFields

	// Map the response to state
	r.mapCircuitToState(ctx, circuit, &data, &resp.Diagnostics)

	// Restore null/empty custom_fields to prevent unwanted updates
	if originalCustomFields.IsNull() || (!originalCustomFields.IsUnknown() && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}

	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the circuit resource.
func (r *CircuitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CircuitResourceModel
	var state CircuitResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := plan
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse circuit ID: %s", err))
		return
	}

	// Build the update request (pass state for merge-aware custom fields)
	updateReq, diags := r.buildCircuitRequest(ctx, &data, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating circuit", map[string]interface{}{
		"id":  id,
		"cid": data.Cid.ValueString(),
	})

	// Update the circuit
	circuit, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitsUpdate(ctx, id).WritableCircuitRequest(*updateReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating circuit",
			utils.FormatAPIError("update circuit", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Updated circuit", map[string]interface{}{
		"id":  circuit.GetId(),
		"cid": circuit.GetCid(),
	})

	// Map the response to state
	r.mapCircuitToState(ctx, circuit, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the circuit resource.
func (r *CircuitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CircuitResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse circuit ID: %s", err))
		return
	}
	tflog.Debug(ctx, "Deleting circuit", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.CircuitsAPI.CircuitsCircuitsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Circuit already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting circuit",
			utils.FormatAPIError("delete circuit", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted circuit", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports a circuit resource.
func (r *CircuitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		id, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid circuit ID",
				fmt.Sprintf("Circuit ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		circuit, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing circuit",
				utils.FormatAPIError("read circuit", err, httpResp),
			)
			return
		}

		var data CircuitResourceModel
		data.Tags = types.SetNull(types.StringType)
		if circuit.HasTags() {
			var tagSlugs []string
			for _, tag := range circuit.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		}
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

		r.mapCircuitToState(ctx, circuit, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, circuit.GetCustomFields(), &resp.Diagnostics)
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildCircuitRequest builds a WritableCircuitRequest from the resource model.
// state parameter: pass nil during Create, pass state during Update for merge-aware custom_fields.
func (r *CircuitResource) buildCircuitRequest(ctx context.Context, data *CircuitResourceModel, state *CircuitResourceModel) (*netbox.WritableCircuitRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Lookup provider (required)
	provider, providerDiags := netboxlookup.LookupProvider(ctx, r.client, data.CircuitProvider.ValueString())
	diags.Append(providerDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Lookup circuit type (required)
	circuitType, typeDiags := netboxlookup.LookupCircuitType(ctx, r.client, data.Type.ValueString())
	diags.Append(typeDiags...)
	if diags.HasError() {
		return nil, diags
	}
	circuitReq := &netbox.WritableCircuitRequest{
		Cid:      data.Cid.ValueString(),
		Provider: *provider,
		Type:     *circuitType,
	}

	// Status
	if utils.IsSet(data.Status) {
		status := netbox.CircuitStatusValue(data.Status.ValueString())
		circuitReq.Status = &status
	}

	// Tenant
	if utils.IsSet(data.Tenant) {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return nil, diags
		}
		circuitReq.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
	} else if data.Tenant.IsNull() {
		// Explicitly set to null to clear the field
		circuitReq.SetTenantNil()
	}

	// Install date
	if utils.IsSet(data.InstallDate) {
		circuitReq.InstallDate = *netbox.NewNullableString(netbox.PtrString(data.InstallDate.ValueString()))
	} else if data.InstallDate.IsNull() {
		// Explicitly set to null to clear the field
		circuitReq.InstallDate = *netbox.NewNullableString(nil)
	}

	// Termination date
	if utils.IsSet(data.TerminationDate) {
		circuitReq.TerminationDate = *netbox.NewNullableString(netbox.PtrString(data.TerminationDate.ValueString()))
	} else if data.TerminationDate.IsNull() {
		// Explicitly set to null to clear the field
		circuitReq.TerminationDate = *netbox.NewNullableString(nil)
	}

	// Commit rate
	if utils.IsSet(data.CommitRate) {
		commitRate, err := utils.SafeInt32FromValue(data.CommitRate)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("CommitRate value overflow: %s", err))
			return nil, diags
		}
		circuitReq.CommitRate = *netbox.NewNullableInt32(netbox.PtrInt32(commitRate))
	} else if data.CommitRate.IsNull() {
		// Explicitly set to null to clear the field
		circuitReq.CommitRate = *netbox.NewNullableInt32(nil)
	}

	// Apply common fields (description, comments, tags)
	utils.ApplyDescription(circuitReq, data.Description)
	utils.ApplyComments(circuitReq, data.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, circuitReq, data.Tags, &diags)
	if diags.HasError() {
		return nil, diags
	}

	// Handle custom fields with merge-aware logic
	if state != nil {
		// Update: merge plan custom fields with existing state custom fields
		utils.ApplyCustomFieldsWithMerge(ctx, circuitReq, data.CustomFields, state.CustomFields, &diags)
	} else {
		// Create: apply plan custom fields directly
		utils.ApplyCustomFields(ctx, circuitReq, data.CustomFields, &diags)
	}

	if diags.HasError() {
		return nil, diags
	}
	return circuitReq, diags
}

// mapCircuitToState maps a Circuit to the Terraform state model.
func (r *CircuitResource) mapCircuitToState(ctx context.Context, circuit *netbox.Circuit, data *CircuitResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", circuit.GetId()))
	data.Cid = types.StringValue(circuit.GetCid())

	// Provider - preserve user input if it matches, otherwise normalize to slug/name
	providerObj := circuit.GetProvider()
	if data.CircuitProvider.IsUnknown() || data.CircuitProvider.IsNull() {
		// During initial creation, set to ID so plan matches apply
		data.CircuitProvider = types.StringValue(fmt.Sprintf("%d", providerObj.GetId()))
	} else {
		userProvider := data.CircuitProvider.ValueString()
		if userProvider == providerObj.GetName() || userProvider == providerObj.GetSlug() || userProvider == providerObj.GetDisplay() || userProvider == fmt.Sprintf("%d", providerObj.GetId()) {
			// Keep user's original value
		} else {
			// Reference changed, update to slug/name
			if providerObj.GetSlug() != "" {
				data.CircuitProvider = types.StringValue(providerObj.GetSlug())
			} else {
				data.CircuitProvider = types.StringValue(providerObj.GetName())
			}
		}
	}

	// Type - preserve user input if it matches, otherwise normalize to slug/name
	typeObj := circuit.GetType()
	if data.Type.IsUnknown() || data.Type.IsNull() {
		// During initial creation, set to ID so plan matches apply
		data.Type = types.StringValue(fmt.Sprintf("%d", typeObj.GetId()))
	} else {
		userType := data.Type.ValueString()
		if userType == typeObj.GetName() || userType == typeObj.GetSlug() || userType == typeObj.GetDisplay() || userType == fmt.Sprintf("%d", typeObj.GetId()) {
			// Keep user's original value
		} else {
			// Reference changed, update to slug/name
			if typeObj.GetSlug() != "" {
				data.Type = types.StringValue(typeObj.GetSlug())
			} else {
				data.Type = types.StringValue(typeObj.GetName())
			}
		}
	}

	// Status
	if circuit.HasStatus() {
		data.Status = types.StringValue(string(circuit.Status.GetValue()))
	} else {
		data.Status = types.StringValue("active")
	}

	// Tenant - preserve user input if it matches, otherwise normalize to slug/name
	if circuit.Tenant.IsSet() && circuit.Tenant.Get() != nil {
		tenantObj := circuit.Tenant.Get()
		if data.Tenant.IsUnknown() || data.Tenant.IsNull() {
			// During initial creation, set to ID so plan matches apply
			data.Tenant = types.StringValue(fmt.Sprintf("%d", tenantObj.GetId()))
		} else {
			userTenant := data.Tenant.ValueString()
			if userTenant == tenantObj.GetName() || userTenant == tenantObj.GetSlug() || userTenant == tenantObj.GetDisplay() || userTenant == fmt.Sprintf("%d", tenantObj.GetId()) {
				// Keep user's original value
			} else {
				// Reference changed, update to slug/name
				if tenantObj.GetSlug() != "" {
					data.Tenant = types.StringValue(tenantObj.GetSlug())
				} else {
					data.Tenant = types.StringValue(tenantObj.GetName())
				}
			}
		}
	} else {
		data.Tenant = types.StringNull()
	}

	// Install date
	if circuit.InstallDate.IsSet() && circuit.InstallDate.Get() != nil {
		data.InstallDate = types.StringValue(*circuit.InstallDate.Get())
	} else {
		data.InstallDate = types.StringNull()
	}

	// Termination date
	if circuit.TerminationDate.IsSet() && circuit.TerminationDate.Get() != nil {
		data.TerminationDate = types.StringValue(*circuit.TerminationDate.Get())
	} else {
		data.TerminationDate = types.StringNull()
	}

	// Commit rate
	if circuit.CommitRate.IsSet() && circuit.CommitRate.Get() != nil {
		data.CommitRate = types.Int64Value(int64(*circuit.CommitRate.Get()))
	} else {
		data.CommitRate = types.Int64Null()
	}

	// Description
	if circuit.HasDescription() && circuit.GetDescription() != "" {
		data.Description = types.StringValue(circuit.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if circuit.HasComments() && circuit.GetComments() != "" {
		data.Comments = types.StringValue(circuit.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Tags (slug list)
	var tagSlugs []string
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	case circuit.HasTags():
		for _, tag := range circuit.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if diags.HasError() {
		return
	}

	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, circuit.GetCustomFields(), diags)
}
