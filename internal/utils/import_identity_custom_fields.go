package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ImportIdentityCustomFieldsModel struct {
	ID           types.String `tfsdk:"id"`
	CustomFields types.List   `tfsdk:"custom_fields"`
}

type ParsedImportIdentityCustomFields struct {
	ID               string
	CustomFields     []CustomFieldModel
	CustomFieldItems []string
	HasCustomFields  bool
}

var allowedCustomFieldIdentityTypes = map[string]struct{}{
	"text":        {},
	"longtext":    {},
	"integer":     {},
	"boolean":     {},
	"date":        {},
	"url":         {},
	"json":        {},
	"select":      {},
	"multiselect": {},
	"object":      {},
	"multiobject": {},
	"multiple":    {},
	"selection":   {},
}

// ParseImportIdentityCustomFields extracts an import identity with optional custom field hints.
// Returns ok=false when identity is nil.
func ParseImportIdentityCustomFields(ctx context.Context, identity *tfsdk.ResourceIdentity, diags *diag.Diagnostics) (ParsedImportIdentityCustomFields, bool) {
	if identity == nil || identity.Raw.IsNull() {
		return ParsedImportIdentityCustomFields{}, false
	}
	var parsed ParsedImportIdentityCustomFields
	var model ImportIdentityCustomFieldsModel
	diags.Append(identity.Get(ctx, &model)...)
	if diags.HasError() {
		return ParsedImportIdentityCustomFields{}, true
	}

	parsed.ID = model.ID.ValueString()
	if model.CustomFields.IsNull() || model.CustomFields.IsUnknown() {
		return parsed, true
	}
	parsed.HasCustomFields = true
	var entries []string
	diags.Append(model.CustomFields.ElementsAs(ctx, &entries, false)...)
	if diags.HasError() {
		return parsed, true
	}

	parsed.CustomFieldItems = entries
	if len(entries) == 0 {
		return parsed, true
	}

	parsed.CustomFields = ParseCustomFieldIdentityEntries(entries, diags)
	return parsed, true
}

// ParseCustomFieldIdentityEntries parses identity entries in the form "name[:type]".
func ParseCustomFieldIdentityEntries(entries []string, diags *diag.Diagnostics) []CustomFieldModel {
	if len(entries) == 0 {
		return nil
	}

	result := make([]CustomFieldModel, 0, len(entries))
	for _, entry := range entries {
		raw := strings.TrimSpace(entry)
		if raw == "" {
			diags.AddError("Invalid custom_fields identity entry", "custom_fields entries must not be empty")
			continue
		}
		name, rawType, _ := strings.Cut(raw, ":")
		name = strings.TrimSpace(name)
		cfType := strings.TrimSpace(rawType)
		if name == "" {
			diags.AddError("Invalid custom_fields identity entry", fmt.Sprintf("custom_fields entry %q must include a field name", entry))
			continue
		}
		if cfType == "" {
			cfType = "text"
		}
		if _, ok := allowedCustomFieldIdentityTypes[cfType]; !ok {
			diags.AddError("Invalid custom_fields identity entry", fmt.Sprintf("custom_fields entry %q has unsupported type %q", entry, cfType))
			continue
		}
		result = append(result, CustomFieldModel{
			Name:  types.StringValue(name),
			Type:  types.StringValue(cfType),
			Value: types.StringValue(""),
		})
	}
	return result
}

// CustomFieldIdentityEntriesFromSet builds identity entry strings from a custom_fields set.
func CustomFieldIdentityEntriesFromSet(ctx context.Context, customFields types.Set, diags *diag.Diagnostics) []string {
	if customFields.IsNull() || customFields.IsUnknown() {
		return nil
	}

	var models []CustomFieldModel
	diags.Append(customFields.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	entries := make([]string, 0, len(models))
	for _, cf := range models {
		name := cf.Name.ValueString()
		if name == "" {
			continue
		}
		cfType := cf.Type.ValueString()
		if cfType == "" {
			cfType = "text"
		}
		entries = append(entries, fmt.Sprintf("%s:%s", name, cfType))
	}
	return entries
}

// SetIdentityCustomFields sets identity values while preserving existing custom_fields when not set.
func SetIdentityCustomFields(ctx context.Context, identity *tfsdk.ResourceIdentity, id types.String, customFields types.Set, diags *diag.Diagnostics) {
	if identity == nil {
		return
	}
	if id.IsNull() || id.IsUnknown() {
		return
	}

	listValue := types.ListNull(types.StringType)
	if identity.Raw.IsNull() {
		if IsSet(customFields) {
			entries := CustomFieldIdentityEntriesFromSet(ctx, customFields, diags)
			listValueNew, listDiags := types.ListValueFrom(ctx, types.StringType, entries)
			diags.Append(listDiags...)
			if diags.HasError() {
				return
			}
			listValue = listValueNew
		}
	} else {
		existing := ImportIdentityCustomFieldsModel{
			ID:           types.StringNull(),
			CustomFields: types.ListNull(types.StringType),
		}
		if getDiags := identity.Get(ctx, &existing); getDiags.HasError() {
			diags.Append(getDiags...)
			return
		}
		listValue = existing.CustomFields
	}

	setDiags := identity.Set(ctx, &ImportIdentityCustomFieldsModel{
		ID:           id,
		CustomFields: listValue,
	})
	diags.Append(setDiags...)
}
