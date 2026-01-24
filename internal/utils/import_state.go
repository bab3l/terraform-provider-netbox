package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportStatePassthroughIDWithValidation sets the import identifier on state with validation.
//
// Behavior:
// - Rejects empty IDs with a clear error message.
// - Optionally enforces numeric IDs (int32).
// - Reuses custom_fields identity parsing for import blocks.
func ImportStatePassthroughIDWithValidation(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
	idPath path.Path,
	requireNumeric bool,
) {
	id := strings.TrimSpace(req.ID)

	if req.Identity != nil {
		parsed, ok := ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics)
		if ok {
			if resp.Diagnostics.HasError() {
				return
			}
			if parsed.ID != "" {
				id = strings.TrimSpace(parsed.ID)
			}
			if parsed.HasCustomFields && resp.Identity != nil {
				listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
				resp.Diagnostics.Append(listDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				resp.Diagnostics.Append(resp.Identity.Set(ctx, &ImportIdentityCustomFieldsModel{
					ID:           types.StringValue(id),
					CustomFields: listValue,
				})...)
				if resp.Diagnostics.HasError() {
					return
				}
			}
		}
	}

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be provided; example: terraform import netbox_* <id>",
		)
		return
	}

	if requireNumeric {
		if _, err := ParseID(id); err != nil {
			resp.Diagnostics.AddError(
				"Invalid import ID",
				fmt.Sprintf("Import ID must be a numeric ID, got: %q", id),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, idPath, id)...)
}
