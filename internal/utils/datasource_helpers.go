package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

// CustomFieldsSetFromAPI maps API custom fields to a Terraform Set value for datasources.
// It returns a null set when custom fields are absent or mapping fails.
func CustomFieldsSetFromAPI(ctx context.Context, hasCustomFields bool, customFields map[string]interface{}, diags *diag.Diagnostics) types.Set {
	attrType := GetCustomFieldsAttributeType().ElemType
	if !hasCustomFields || len(customFields) == 0 {
		return types.SetNull(attrType)
	}

	models := MapAllCustomFieldsToModels(customFields)
	setValue, setDiags := types.SetValueFrom(ctx, attrType, models)
	if diags != nil {
		diags.Append(setDiags...)
	}
	if setDiags.HasError() {
		return types.SetNull(attrType)
	}
	return setValue
}
