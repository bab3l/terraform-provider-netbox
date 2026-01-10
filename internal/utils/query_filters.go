package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// QueryFilterModel is a reusable model for AWS-style filter blocks:
//
//	filter {
//	  name   = "..."
//	  values = ["..."]
//	}
//
// It is intentionally generic so multiple plural/query datasources can reuse it.
type QueryFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// ExpandQueryFilters validates and normalizes filter blocks into a map.
//
// - Filter names are lowercased and trimmed.
// - Values are trimmed; empty values are discarded.
// - Multiple blocks of the same name are merged.
func ExpandQueryFilters(ctx context.Context, filters []QueryFilterModel) (map[string][]string, diag.Diagnostics) {
	result := make(map[string][]string)
	var diags diag.Diagnostics

	for i, f := range filters {
		if f.Name.IsNull() || f.Name.IsUnknown() {
			diags.AddError("Invalid filter", fmt.Sprintf("filter[%d].name must be set", i))
			continue
		}

		name := strings.ToLower(strings.TrimSpace(f.Name.ValueString()))
		if name == "" {
			diags.AddError("Invalid filter", fmt.Sprintf("filter[%d].name must not be empty", i))
			continue
		}

		if f.Values.IsNull() || f.Values.IsUnknown() {
			diags.AddError("Invalid filter", fmt.Sprintf("filter[%d].values must be set", i))
			continue
		}

		var raw []string
		valueDiags := f.Values.ElementsAs(ctx, &raw, false)
		diags.Append(valueDiags...)
		if valueDiags.HasError() {
			continue
		}

		values := make([]string, 0, len(raw))
		seen := make(map[string]struct{}, len(raw))
		for _, v := range raw {
			vv := strings.TrimSpace(v)
			if vv == "" {
				continue
			}
			if _, exists := seen[vv]; exists {
				continue
			}
			seen[vv] = struct{}{}
			values = append(values, vv)
		}

		if len(values) == 0 {
			diags.AddError("Invalid filter", fmt.Sprintf("filter[%d].values must include at least one non-empty value", i))
			continue
		}

		result[name] = append(result[name], values...)
	}

	return result, diags
}
