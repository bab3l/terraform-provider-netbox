package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CustomFieldsMergeModifier is a plan modifier that merges planned custom fields
// with custom fields from prior state. This enables partial custom field management
// where users can manage some fields in Terraform while preserving others managed
// externally (NetBox UI, automation, etc.).
//
// Behavior:
//   - If config is null: Set plan to null (don't manage custom fields)
//   - If config is empty set []: Set plan to empty (remove all custom fields)
//   - If config has values: Merge config with prior state, set plan to merged result
//
// This modifier ensures that unmanaged custom fields appear in the plan so they
// will be preserved in the state after apply.
func CustomFieldsMergeModifier() planmodifier.Set {
	return customFieldsMergeModifier{}
}

type customFieldsMergeModifier struct{}

func (m customFieldsMergeModifier) Description(_ context.Context) string {
	return "Merges planned custom fields with custom fields from prior state to enable partial management."
}

func (m customFieldsMergeModifier) MarkdownDescription(_ context.Context) string {
	return "Merges planned custom fields with custom fields from prior state to enable partial management."
}

func (m customFieldsMergeModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// If config is null/unknown, keep plan as-is (don't manage)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		tflog.Debug(ctx, "CustomFieldsMergeModifier: Config is null/unknown, preserving plan")
		return
	}

	// If this is a create operation (no prior state), use config as-is
	if req.StateValue.IsNull() {
		tflog.Debug(ctx, "CustomFieldsMergeModifier: No prior state (create), using config")
		return
	}

	// Parse config and state custom fields
	var configFields, stateFields []customFieldElement
	diags := req.ConfigValue.ElementsAs(ctx, &configFields, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.StateValue.ElementsAs(ctx, &stateFields, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If config is empty set, use it as-is (explicit removal of all fields)
	if len(configFields) == 0 {
		tflog.Debug(ctx, "CustomFieldsMergeModifier: Config is empty set, removing all custom fields")
		return
	}

	// Build map of config field names
	configFieldNames := make(map[string]bool)
	for _, cf := range configFields {
		name := cf.Name.ValueString()
		configFieldNames[name] = true
	}

	// Add state fields that aren't in config (preserve unmanaged fields)
	mergedFields := make([]customFieldElement, 0, len(configFields)+len(stateFields))
	mergedFields = append(mergedFields, configFields...)

	for _, sf := range stateFields {
		name := sf.Name.ValueString()
		if !configFieldNames[name] {
			// This field is in state but not in config - preserve it
			mergedFields = append(mergedFields, sf)
			tflog.Debug(ctx, "CustomFieldsMergeModifier: Preserving unmanaged field from state", map[string]interface{}{
				"field": name,
			})
		}
	}

	// Convert merged fields to a Set
	elementType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	}

	elements := make([]attr.Value, len(mergedFields))
	for i, mf := range mergedFields {
		elements[i] = types.ObjectValueMust(elementType.AttrTypes, map[string]attr.Value{
			"name":  mf.Name,
			"type":  mf.Type,
			"value": mf.Value,
		})
	}

	mergedSet, diags := types.SetValue(elementType, elements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "CustomFieldsMergeModifier: Merged custom fields", map[string]interface{}{
		"config_count": len(configFields),
		"state_count":  len(stateFields),
		"merged_count": len(mergedFields),
	})

	resp.PlanValue = mergedSet
}

// customFieldElement represents a custom field element in the set.
// This matches the structure of the custom_fields set elements.
type customFieldElement struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}
