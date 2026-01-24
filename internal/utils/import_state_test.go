package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestImportStatePassthroughIDWithValidation_EmptyID(t *testing.T) {
	ctx := context.Background()

	req := resource.ImportStateRequest{ID: " "}
	resp := &resource.ImportStateResponse{}

	ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)

	require.True(t, resp.Diagnostics.HasError())
}

func TestImportStatePassthroughIDWithValidation_InvalidNumericID(t *testing.T) {
	ctx := context.Background()

	req := resource.ImportStateRequest{ID: "abc"}
	resp := &resource.ImportStateResponse{}

	ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)

	require.True(t, resp.Diagnostics.HasError())
}

func TestImportStatePassthroughIDWithValidation_CustomFieldsIdentity(t *testing.T) {
	ctx := context.Background()

	identityValue, diags := types.ObjectValue(
		map[string]attr.Type{
			"id":            types.StringType,
			"custom_fields": types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"id": types.StringValue("123"),
			"custom_fields": func() attr.Value {
				list, listDiags := types.ListValueFrom(ctx, types.StringType, []string{"cf_text:text"})
				require.False(t, listDiags.HasError())
				return list
			}(),
		},
	)
	require.False(t, diags.HasError())
	terraformValue, valueErr := identityValue.ToTerraformValue(ctx)
	require.NoError(t, valueErr)

	identitySchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Optional: true},
			"custom_fields": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
	stateSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		},
	}
	stateRaw := tftypes.NewValue(
		tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}},
		map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, nil)},
	)

	identity := &tfsdk.ResourceIdentity{Raw: terraformValue, Schema: identitySchema}

	req := resource.ImportStateRequest{ID: " ", Identity: identity}
	resp := &resource.ImportStateResponse{
		Identity: &tfsdk.ResourceIdentity{Schema: identitySchema},
		State:    tfsdk.State{Schema: stateSchema, Raw: stateRaw},
	}

	ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)

	require.False(t, resp.Diagnostics.HasError())

	var stateID types.String
	stateDiags := resp.State.GetAttribute(ctx, path.Root("id"), &stateID)
	require.False(t, stateDiags.HasError())
	require.Equal(t, "123", stateID.ValueString())

	var model ImportIdentityCustomFieldsModel
	identityDiags := resp.Identity.Get(ctx, &model)
	require.False(t, identityDiags.HasError())
	require.Equal(t, "123", model.ID.ValueString())

	var entries []string
	listDiags := model.CustomFields.ElementsAs(ctx, &entries, false)
	require.False(t, listDiags.HasError())
	require.Equal(t, []string{"cf_text:text"}, entries)
}
