package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCustomFieldsSetter is a mock implementation of CustomFieldsSetter for testing.
type mockCustomFieldsSetter struct {
	customFields map[string]interface{}
}

func (m *mockCustomFieldsSetter) SetCustomFields(v map[string]interface{}) {
	m.customFields = v
}

// Helper to create a custom field Set from models.
func createCustomFieldSet(t *testing.T, models []CustomFieldModel) types.Set {
	if len(models) == 0 {
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
		})
	}

	elements := make([]attr.Value, len(models))
	for i, model := range models {
		obj, diags := types.ObjectValue(
			map[string]attr.Type{
				"name":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
			map[string]attr.Value{
				"name":  model.Name,
				"type":  model.Type,
				"value": model.Value,
			},
		)
		require.False(t, diags.HasError(), "Failed to create object: %v", diags)
		elements[i] = obj
	}

	set, diags := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
		},
		elements,
	)
	require.False(t, diags.HasError(), "Failed to create set: %v", diags)
	return set
}

func TestMergeCustomFieldSets_BothNull(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}

	planSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})
	stateSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})

	result := MergeCustomFieldSets(ctx, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	assert.Empty(t, result, "Result should be empty when both plan and state are null")
}

func TestMergeCustomFieldSets_PlanNullStateHasFields(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}

	planSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_a"),
		},
		{
			Name:  types.StringValue("field_b"),
			Type:  types.StringValue("integer"),
			Value: types.StringValue("42"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	result := MergeCustomFieldSets(ctx, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, "value_a", result["field_a"], "State field_a should be preserved")
	assert.Equal(t, 42, result["field_b"], "State field_b should be preserved as integer")
	assert.Len(t, result, 2, "Should have 2 fields from state")
}

func TestMergeCustomFieldSets_PlanHasFieldsStateNull(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}

	planModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_a"),
		},
	}
	planSet := createCustomFieldSet(t, planModels)

	stateSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})

	result := MergeCustomFieldSets(ctx, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, "value_a", result["field_a"], "Plan field_a should be set")
	assert.Len(t, result, 1, "Should have 1 field from plan")
}

func TestMergeCustomFieldSets_PlanOverridesState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}

	planModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("new_value"),
		},
	}
	planSet := createCustomFieldSet(t, planModels)

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("old_value"),
		},
		{
			Name:  types.StringValue("field_b"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_b"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	result := MergeCustomFieldSets(ctx, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, "new_value", result["field_a"], "Plan should override state for field_a")
	assert.Equal(t, "value_b", result["field_b"], "State field_b should be preserved")
	assert.Len(t, result, 2, "Should have 2 fields (1 overridden, 1 preserved)")
}

func TestMergeCustomFieldSets_EmptyValueRemovesField(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}

	planModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue(""), // Empty value should remove field
		},
	}
	planSet := createCustomFieldSet(t, planModels)

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_a"),
		},
		{
			Name:  types.StringValue("field_b"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_b"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	result := MergeCustomFieldSets(ctx, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	assert.NotContains(t, result, "field_a", "field_a should be removed (empty value)")
	assert.Equal(t, "value_b", result["field_b"], "field_b should be preserved")
	assert.Len(t, result, 1, "Should have 1 field (1 removed, 1 preserved)")
}

func TestMergeCustomFieldSets_MultipleTypes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}

	planModels := []CustomFieldModel{
		{
			Name:  types.StringValue("text_field"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("text_value"),
		},
		{
			Name:  types.StringValue("int_field"),
			Type:  types.StringValue("integer"),
			Value: types.StringValue("42"),
		},
		{
			Name:  types.StringValue("bool_field"),
			Type:  types.StringValue("boolean"),
			Value: types.StringValue("true"),
		},
	}
	planSet := createCustomFieldSet(t, planModels)

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("old_field"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("old_value"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	result := MergeCustomFieldSets(ctx, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, "text_value", result["text_field"])
	assert.Equal(t, 42, result["int_field"], "Integer should be converted")
	assert.Equal(t, true, result["bool_field"], "Boolean should be converted")
	assert.Equal(t, "old_value", result["old_field"], "State field should be preserved")
	assert.Len(t, result, 4, "Should have 4 fields (3 from plan, 1 from state)")
}

func TestApplyCustomFieldsWithMerge_PlanNullPreservesState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}
	mock := &mockCustomFieldsSetter{}

	planSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_a"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	ApplyCustomFieldsWithMerge(ctx, mock, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	require.NotNil(t, mock.customFields)
	assert.Equal(t, "value_a", mock.customFields["field_a"], "State should be preserved when plan is null")
}

func TestApplyCustomFieldsWithMerge_BothNullSendsEmptyMap(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}
	mock := &mockCustomFieldsSetter{}

	planSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})
	stateSet := types.SetNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		},
	})

	ApplyCustomFieldsWithMerge(ctx, mock, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	require.NotNil(t, mock.customFields)
	assert.Empty(t, mock.customFields, "Empty map should be sent when both are null")
}

func TestApplyCustomFieldsWithMerge_PlanSetMergesWithState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}
	mock := &mockCustomFieldsSetter{}

	planModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("new_value"),
		},
	}
	planSet := createCustomFieldSet(t, planModels)

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("old_value"),
		},
		{
			Name:  types.StringValue("field_b"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_b"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	ApplyCustomFieldsWithMerge(ctx, mock, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	require.NotNil(t, mock.customFields)
	assert.Equal(t, "new_value", mock.customFields["field_a"], "Plan should override state")
	assert.Equal(t, "value_b", mock.customFields["field_b"], "State should be preserved")
}

func TestApplyCustomFieldsWithMerge_EmptyStringRemovesField(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	diags := diag.Diagnostics{}
	mock := &mockCustomFieldsSetter{}

	planModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue(""), // Empty removes field
		},
	}
	planSet := createCustomFieldSet(t, planModels)

	stateModels := []CustomFieldModel{
		{
			Name:  types.StringValue("field_a"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_a"),
		},
		{
			Name:  types.StringValue("field_b"),
			Type:  types.StringValue("text"),
			Value: types.StringValue("value_b"),
		},
	}
	stateSet := createCustomFieldSet(t, stateModels)

	ApplyCustomFieldsWithMerge(ctx, mock, planSet, stateSet, &diags)

	assert.False(t, diags.HasError())
	require.NotNil(t, mock.customFields)
	assert.NotContains(t, mock.customFields, "field_a", "field_a should be removed")
	assert.Equal(t, "value_b", mock.customFields["field_b"], "field_b should be preserved")
}
