package utils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type CustomFieldValueFilter struct {
	Name  string
	Value string
}

// ParseCustomFieldValueFilters parses filters of the form "<field_name>=<value>".
func ParseCustomFieldValueFilters(values []string) ([]CustomFieldValueFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	filters := make([]CustomFieldValueFilter, 0, len(values))
	for _, raw := range values {
		name, value, ok := strings.Cut(raw, "=")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if !ok || name == "" {
			diags.AddError(
				"Invalid filter values",
				fmt.Sprintf("custom_field_value entries must be formatted as <field_name>=<value>; got %q", raw),
			)
			continue
		}

		filters = append(filters, CustomFieldValueFilter{Name: name, Value: value})
	}

	return filters, diags
}

// MatchesCustomFieldFilters applies custom field existence/value filters to a custom field map.
//
// Semantics:
// - If existsNames is non-empty, at least one named custom field must exist (key present and value non-nil).
// - If valueFilters is non-empty, at least one <field>=<value> must match.
// - Both groups are ANDed together.
func MatchesCustomFieldFilters(customFields map[string]interface{}, existsNames []string, valueFilters []CustomFieldValueFilter) bool {
	if len(existsNames) > 0 {
		matched := false
		for _, name := range existsNames {
			if v, ok := customFields[name]; ok && v != nil {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(valueFilters) > 0 {
		matched := false
		for _, f := range valueFilters {
			if v, ok := customFields[f.Name]; ok && v != nil && fmt.Sprint(v) == f.Value {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
