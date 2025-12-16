// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &PrefixResource{}

	_ resource.ResourceWithConfigure = &PrefixResource{}

	_ resource.ResourceWithImportState = &PrefixResource{}
)

// NewPrefixResource returns a new Prefix resource.

func NewPrefixResource() resource.Resource {

	return &PrefixResource{}

}

// PrefixResource defines the resource implementation.

type PrefixResource struct {
	client *netbox.APIClient
}

// PrefixResourceModel describes the resource data model.

type PrefixResourceModel struct {
	ID types.String `tfsdk:"id"`

	Prefix types.String `tfsdk:"prefix"`

	Site types.String `tfsdk:"site"`

	VRF types.String `tfsdk:"vrf"`

	Tenant types.String `tfsdk:"tenant"`

	VLAN types.String `tfsdk:"vlan"`

	Status types.String `tfsdk:"status"`

	Role types.String `tfsdk:"role"`

	IsPool types.Bool `tfsdk:"is_pool"`

	MarkUtilized types.Bool `tfsdk:"mark_utilized"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`
}

// Metadata returns the resource type name.

func (r *PrefixResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_prefix"

}

// Schema defines the schema for the resource.

func (r *PrefixResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a prefix in Netbox. A prefix represents an IP address space (CIDR).",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the prefix.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"prefix": schema.StringAttribute{

				MarkdownDescription: "The IP prefix in CIDR notation (e.g., 192.168.1.0/24).",

				Required: true,
			},

			"site": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the site this prefix is assigned to.",

				Optional: true,
			},

			"vrf": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the VRF this prefix is assigned to.",

				Optional: true,
			},

			"tenant": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the tenant this prefix is assigned to.",

				Optional: true,
			},

			"vlan": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the VLAN this prefix is assigned to.",

				Optional: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "The status of the prefix. Valid values are: `container`, `active`, `reserved`, `deprecated`. Defaults to `active`.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString("active"),
			},

			"role": schema.StringAttribute{

				MarkdownDescription: "The name or ID of the role for this prefix.",

				Optional: true,
			},

			"is_pool": schema.BoolAttribute{

				MarkdownDescription: "If true, all IP addresses within this prefix are considered usable. Defaults to false.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(false),
			},

			"mark_utilized": schema.BoolAttribute{

				MarkdownDescription: "If true, treat the prefix as fully utilized. Defaults to false.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(false),
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description for the prefix.",

				Optional: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Comments for the prefix.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *PrefixResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *PrefixResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data PrefixResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Create the prefix request

	prefixRequest := netbox.NewWritablePrefixRequest(data.Prefix.ValueString())

	// Set optional fields

	r.setOptionalFields(ctx, prefixRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating prefix", map[string]interface{}{

		"prefix": data.Prefix.ValueString(),
	})

	// Create the prefix

	prefix, httpResp, err := r.client.IpamAPI.IpamPrefixesCreate(ctx).WritablePrefixRequest(*prefixRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating prefix",

			utils.FormatAPIError("create prefix", err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapPrefixToState(ctx, prefix, &data)

	tflog.Debug(ctx, "Created prefix", map[string]interface{}{

		"id": data.ID.ValueString(),

		"prefix": data.Prefix.ValueString(),
	})

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *PrefixResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data PrefixResourceModel

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

	tflog.Debug(ctx, "Reading prefix", map[string]interface{}{

		"id": id,
	})

	// Get the prefix from Netbox

	prefix, httpResp, err := r.client.IpamAPI.IpamPrefixesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading prefix",

			utils.FormatAPIError(fmt.Sprintf("read prefix ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapPrefixToState(ctx, prefix, &data)

	tflog.Debug(ctx, "Read prefix", map[string]interface{}{

		"id": data.ID.ValueString(),

		"prefix": data.Prefix.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state on success.

func (r *PrefixResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data PrefixResourceModel

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

	// Create the prefix request

	prefixRequest := netbox.NewWritablePrefixRequest(data.Prefix.ValueString())

	// Set optional fields

	r.setOptionalFields(ctx, prefixRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Updating prefix", map[string]interface{}{

		"id": id,

		"prefix": data.Prefix.ValueString(),
	})

	// Update the prefix

	prefix, httpResp, err := r.client.IpamAPI.IpamPrefixesUpdate(ctx, id).WritablePrefixRequest(*prefixRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating prefix",

			utils.FormatAPIError(fmt.Sprintf("update prefix ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapPrefixToState(ctx, prefix, &data)

	tflog.Debug(ctx, "Updated prefix", map[string]interface{}{

		"id": data.ID.ValueString(),

		"prefix": data.Prefix.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state on success.

func (r *PrefixResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data PrefixResourceModel

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

	tflog.Debug(ctx, "Deleting prefix", map[string]interface{}{

		"id": id,
	})

	// Delete the prefix

	httpResp, err := r.client.IpamAPI.IpamPrefixesDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error deleting prefix",

			utils.FormatAPIError(fmt.Sprintf("delete prefix ID %d", id), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted prefix", map[string]interface{}{

		"id": id,
	})

}

func (r *PrefixResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// setOptionalFields sets optional fields on the prefix request from the resource model.

func (r *PrefixResource) setOptionalFields(ctx context.Context, prefixRequest *netbox.WritablePrefixRequest, data *PrefixResourceModel, diags *diag.Diagnostics) {

	// Site

	if utils.IsSet(data.Site) {

		site, siteDiags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())

		diags.Append(siteDiags...)

		if diags.HasError() {

			return

		}

		prefixRequest.Site = *netbox.NewNullableBriefSiteRequest(site)

	}

	// VRF

	if utils.IsSet(data.VRF) {

		vrf, vrfDiags := netboxlookup.LookupVRF(ctx, r.client, data.VRF.ValueString())

		diags.Append(vrfDiags...)

		if diags.HasError() {

			return

		}

		prefixRequest.Vrf = *netbox.NewNullableBriefVRFRequest(vrf)

	}

	// Tenant

	if utils.IsSet(data.Tenant) {

		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		diags.Append(tenantDiags...)

		if diags.HasError() {

			return

		}

		prefixRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)

	}

	// VLAN

	if utils.IsSet(data.VLAN) {

		vlan, vlanDiags := netboxlookup.LookupVLAN(ctx, r.client, data.VLAN.ValueString())

		diags.Append(vlanDiags...)

		if diags.HasError() {

			return

		}

		prefixRequest.Vlan = *netbox.NewNullableBriefVLANRequest(vlan)

	}

	// Status

	if utils.IsSet(data.Status) {

		status := netbox.PatchedWritablePrefixRequestStatus(data.Status.ValueString())

		prefixRequest.Status = &status

	}

	// Role

	if utils.IsSet(data.Role) {

		role, roleDiags := netboxlookup.LookupRole(ctx, r.client, data.Role.ValueString())

		diags.Append(roleDiags...)

		if diags.HasError() {

			return

		}

		prefixRequest.Role = *netbox.NewNullableBriefRoleRequest(role)

	}

	// IsPool

	if utils.IsSet(data.IsPool) {

		isPool := data.IsPool.ValueBool()

		prefixRequest.IsPool = &isPool

	}

	// MarkUtilized

	if utils.IsSet(data.MarkUtilized) {

		markUtilized := data.MarkUtilized.ValueBool()

		prefixRequest.MarkUtilized = &markUtilized

	}

	// Description

	prefixRequest.Description = utils.StringPtr(data.Description)

	// Comments

	prefixRequest.Comments = utils.StringPtr(data.Comments)

	// Handle tags

	if utils.IsSet(data.Tags) {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		prefixRequest.Tags = tags

	}

}

// mapPrefixToState maps a Netbox Prefix to the Terraform state model.

func (r *PrefixResource) mapPrefixToState(ctx context.Context, prefix *netbox.Prefix, data *PrefixResourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", prefix.Id))

	data.Prefix = types.StringValue(prefix.Prefix)

	// Site - preserve user input if it matches

	if prefix.Site.IsSet() && prefix.Site.Get() != nil {

		siteObj := prefix.Site.Get()

		userSite := data.Site.ValueString()

		if userSite == siteObj.Name || userSite == siteObj.Slug || userSite == siteObj.Display || userSite == fmt.Sprintf("%d", siteObj.Id) {

			// Keep user's original value

		} else {

			data.Site = types.StringValue(siteObj.Name)

		}

	} else {

		data.Site = types.StringNull()

	}

	// VRF - preserve user input if it matches

	if prefix.Vrf.IsSet() && prefix.Vrf.Get() != nil {

		vrfObj := prefix.Vrf.Get()

		userVrf := data.VRF.ValueString()

		if userVrf == vrfObj.Name || userVrf == vrfObj.Display || userVrf == fmt.Sprintf("%d", vrfObj.Id) {

			// Keep user's original value

		} else {

			data.VRF = types.StringValue(vrfObj.Name)

		}

	} else {

		data.VRF = types.StringNull()

	}

	// Tenant - preserve user input if it matches

	if prefix.Tenant.IsSet() && prefix.Tenant.Get() != nil {

		tenantObj := prefix.Tenant.Get()

		userTenant := data.Tenant.ValueString()

		if userTenant == tenantObj.Name || userTenant == tenantObj.Slug || userTenant == tenantObj.Display || userTenant == fmt.Sprintf("%d", tenantObj.Id) {

			// Keep user's original value

		} else {

			data.Tenant = types.StringValue(tenantObj.Name)

		}

	} else {

		data.Tenant = types.StringNull()

	}

	// VLAN - preserve user input if it matches

	if prefix.Vlan.IsSet() && prefix.Vlan.Get() != nil {

		vlanObj := prefix.Vlan.Get()

		userVlan := data.VLAN.ValueString()

		if userVlan == vlanObj.Display || userVlan == vlanObj.Name || userVlan == fmt.Sprintf("%d", vlanObj.Id) || userVlan == fmt.Sprintf("%d", vlanObj.Vid) {

			// Keep user's original value

		} else {

			data.VLAN = types.StringValue(vlanObj.Display)

		}

	} else {

		data.VLAN = types.StringNull()

	}

	// Status

	if prefix.Status != nil {

		data.Status = types.StringValue(string(prefix.Status.GetValue()))

	} else {

		data.Status = types.StringNull()

	}

	// Role - preserve user input if it matches

	if prefix.Role.IsSet() && prefix.Role.Get() != nil {

		roleObj := prefix.Role.Get()

		userRole := data.Role.ValueString()

		if userRole == roleObj.Name || userRole == roleObj.Slug || userRole == fmt.Sprintf("%d", roleObj.Id) {

			// Keep user's original value

		} else {

			data.Role = types.StringValue(roleObj.Name)

		}

	} else {

		data.Role = types.StringNull()

	}

	// IsPool

	if prefix.IsPool != nil {

		data.IsPool = types.BoolValue(*prefix.IsPool)

	} else {

		data.IsPool = types.BoolNull()

	}

	// MarkUtilized

	if prefix.MarkUtilized != nil {

		data.MarkUtilized = types.BoolValue(*prefix.MarkUtilized)

	} else {

		data.MarkUtilized = types.BoolNull()

	}

	// Description

	if prefix.Description != nil && *prefix.Description != "" {

		data.Description = types.StringValue(*prefix.Description)

	} else {

		data.Description = types.StringNull()

	}

	// Comments

	if prefix.Comments != nil && *prefix.Comments != "" {

		data.Comments = types.StringValue(*prefix.Comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Tags

	if len(prefix.Tags) > 0 {

		tags := utils.NestedTagsToTagModels(prefix.Tags)

		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

}
