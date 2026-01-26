package schema

import (
	"context"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReferencePreferSlugPlanModifier returns a plan modifier that normalizes
// reference values to a slug when possible. This is intended for targeted
// use where tests expect slug persistence even if config uses IDs.
func ReferencePreferSlugPlanModifier(resourceType string) planmodifier.String {
	return preferReferenceSlugModifier{resourceType: resourceType}
}

type preferReferenceSlugModifier struct {
	resourceType string
}

func (m preferReferenceSlugModifier) Description(ctx context.Context) string {
	return "Prefers slug values for reference attributes when resolvable"
}

func (m preferReferenceSlugModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m preferReferenceSlugModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	configValue := req.ConfigValue.ValueString()
	if configValue == "" {
		return
	}

	client, ok := ctx.Value("netbox_client").(*netbox.APIClient)
	if !ok {
		return
	}

	if m.resourceType != "site" {
		return
	}

	siteRef, diags := netboxlookup.LookupSite(ctx, client, configValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || siteRef == nil {
		return
	}

	if slug := siteRef.GetSlug(); slug != "" {
		resp.PlanValue = types.StringValue(slug)
		return
	}

	if name := siteRef.GetName(); name != "" {
		resp.PlanValue = types.StringValue(name)
		return
	}
}
