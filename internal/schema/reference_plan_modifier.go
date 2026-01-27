package schema

import (
	"context"
	"strconv"
	"strings"

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

// ReferenceResolveToIDPlanModifier returns a plan modifier that normalizes
// configured reference values to IDs when resolvable.
func ReferenceResolveToIDPlanModifier(resourceType string) planmodifier.String {
	return resolveReferenceToIDModifier{resourceType: normalizeReferenceType(resourceType)}
}

type resolveReferenceToIDModifier struct {
	resourceType string
}

func (m resolveReferenceToIDModifier) Description(ctx context.Context) string {
	return "Resolves reference values to IDs when possible"
}

func (m resolveReferenceToIDModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m resolveReferenceToIDModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
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

	resourceType := getResourceTypeFromAttribute(req.Path.String())
	if m.resourceType != "" {
		resourceType = m.resourceType
	}
	if resourceType == "" {
		return
	}

	id, diags := netboxlookup.LookupReferenceID(ctx, client, resourceType, configValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || id == 0 {
		return
	}

	resp.PlanValue = types.StringValue(strconv.FormatInt(int64(id), 10))
}

func normalizeReferenceType(resourceType string) string {
	resourceType = strings.TrimSpace(strings.ToLower(resourceType))
	resourceType = strings.ReplaceAll(resourceType, " ", "_")

	return resourceType
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
