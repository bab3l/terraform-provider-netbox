package utils

import "github.com/hashicorp/terraform-plugin-framework/diag"

// ExpectSingleResult validates that a list result has exactly one item.
// It appends diagnostics for not found or multiple results and returns the item and true when valid.
func ExpectSingleResult[T any](results []T, notFoundTitle, notFoundMsg, multipleTitle, multipleMsg string, diags *diag.Diagnostics) (*T, bool) {
	if len(results) == 0 {
		diags.AddError(notFoundTitle, notFoundMsg)
		return nil, false
	}
	if len(results) > 1 {
		diags.AddError(multipleTitle, multipleMsg)
		return nil, false
	}
	return &results[0], true
}
