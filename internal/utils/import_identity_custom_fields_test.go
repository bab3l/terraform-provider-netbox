package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestParseCustomFieldIdentityEntries(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	t.Run("valid entries", func(t *testing.T) {
		t.Parallel()
		var diags diag.Diagnostics
		models := ParseCustomFieldIdentityEntries([]string{"cf_text:text", "cf_bool:boolean"}, &diags)
		require.False(t, diags.HasError())
		require.Len(t, models, 2)
		require.Equal(t, "cf_text", models[0].Name.ValueString())
		require.Equal(t, "text", models[0].Type.ValueString())
		require.Equal(t, "cf_bool", models[1].Name.ValueString())
		require.Equal(t, "boolean", models[1].Type.ValueString())
	})

	t.Run("default type", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics
		models := ParseCustomFieldIdentityEntries([]string{"cf_default"}, &diags)
		require.False(t, diags.HasError())
		require.Len(t, models, 1)
		require.Equal(t, "text", models[0].Type.ValueString())
	})

	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics
		models := ParseCustomFieldIdentityEntries([]string{"cf_bad:unknown"}, &diags)
		require.True(t, diags.HasError())
		require.Empty(t, models)
	})

	t.Run("entries from set", func(t *testing.T) {
		t.Parallel()

		customFields := []CustomFieldModel{
			{Name: types.StringValue("cf_text"), Type: types.StringValue("text"), Value: types.StringValue("")},
			{Name: types.StringValue("cf_int"), Type: types.StringValue("integer"), Value: types.StringValue("")},
		}
		set, diags := types.SetValueFrom(ctx, GetCustomFieldsAttributeType().ElemType, customFields)
		require.False(t, diags.HasError())
		var entryDiags diag.Diagnostics
		entries := CustomFieldIdentityEntriesFromSet(ctx, set, &entryDiags)
		require.False(t, entryDiags.HasError())
		require.ElementsMatch(t, []string{"cf_text:text", "cf_int:integer"}, entries)
	})
}
