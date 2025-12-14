// Package validators provides custom validators for the Netbox provider.
package validators

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// CustomFieldValueValidator validates that a custom field value is appropriate for its type.
type CustomFieldValueValidator struct {
	customFieldType string
}

// Description returns a description of the validator.
func (v CustomFieldValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be valid for custom field type '%s'", v.customFieldType)
}

// MarkdownDescription returns a markdown description of the validator.
func (v CustomFieldValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v CustomFieldValueValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// Get the custom field type from the configuration
	// This is a simplified approach - in a real implementation, you'd extract the type from the same object
	cfType := v.customFieldType

	switch cfType {
	case "integer":
		if _, err := strconv.Atoi(value); err != nil {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Integer Value",
				fmt.Sprintf("The value '%s' is not a valid integer: %s", value, err),
			)
		}
	case "boolean":
		if value != "true" && value != "false" {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Boolean Value",
				fmt.Sprintf("The value '%s' must be either 'true' or 'false'", value),
			)
		}
	case "json":
		var jsonValue interface{}
		if err := json.Unmarshal([]byte(value), &jsonValue); err != nil {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid JSON Value",
				fmt.Sprintf("The value '%s' is not valid JSON: %s", value, err),
			)
		}
	case "multiselect", "multiple":
		// Validate comma-separated values (basic check)
		if strings.TrimSpace(value) == "" {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Multiselect Value",
				"Multiselect values cannot be empty",
			)
		}
	case "url":
		// Basic URL validation - starts with http:// or https://
		if value != "" && !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid URL Value",
				fmt.Sprintf("The value '%s' must be a valid URL starting with http:// or https://", value),
			)
		}
	}
}

// ValidCustomFieldValue returns a validator which ensures that the custom field value
// is appropriate for the specified type.
func ValidCustomFieldValue(customFieldType string) validator.String {
	return CustomFieldValueValidator{
		customFieldType: customFieldType,
	}
}

// SlugValidator validates that a string is a valid slug.
type SlugValidator struct{}

// Description returns a description of the validator.
func (v SlugValidator) Description(_ context.Context) string {
	return "value must be a valid slug (lowercase letters, numbers, hyphens, and underscores only)"
}

// MarkdownDescription returns a markdown description of the validator.
func (v SlugValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v SlugValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// Check for valid slug characters
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Slug Format",
				fmt.Sprintf("Slug '%s' contains invalid character '%c'. Only lowercase letters, numbers, hyphens, and underscores are allowed.", value, char),
			)
			return
		}
	}

	// Check that it doesn't start or end with hyphens/underscores
	if strings.HasPrefix(value, "-") || strings.HasPrefix(value, "_") ||
		strings.HasSuffix(value, "-") || strings.HasSuffix(value, "_") {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Slug Format",
			fmt.Sprintf("Slug '%s' cannot start or end with hyphens or underscores", value),
		)
	}
}

// ValidSlug returns a validator which ensures that the value is a valid slug.
func ValidSlug() validator.String {
	return SlugValidator{}
}

// IntegerRegex returns a regex for validating integer values.
func IntegerRegex() *regexp.Regexp {
	return regexp.MustCompile(`^[0-9]+$`)
}

// HexColorRegex returns a regex for validating 6-character hexadecimal color codes.
func HexColorRegex() *regexp.Regexp {
	return regexp.MustCompile(`^[0-9a-fA-F]{6}$`)
}

// CustomFieldNameValidator validates custom field names.
type CustomFieldNameValidator struct{}

// Description returns a description of the validator.
func (v CustomFieldNameValidator) Description(_ context.Context) string {
	return "custom field name must start with a letter and contain only letters, numbers, and underscores"
}

// MarkdownDescription returns a markdown description of the validator.
func (v CustomFieldNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v CustomFieldNameValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// Check if name starts with letter and contains only valid characters
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, value)
	if !matched {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Custom Field Name",
			fmt.Sprintf("Custom field name '%s' must start with a letter and contain only letters, numbers, and underscores", value),
		)
	}
}

// ValidCustomFieldName returns a validator for custom field names.
func ValidCustomFieldName() validator.String {
	return CustomFieldNameValidator{}
}

// CustomFieldTypeValidator validates custom field types.
type CustomFieldTypeValidator struct{}

// Description returns a description of the validator.
func (v CustomFieldTypeValidator) Description(_ context.Context) string {
	return "custom field type must be one of: text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject"
}

// MarkdownDescription returns a markdown description of the validator.
func (v CustomFieldTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v CustomFieldTypeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	validTypes := []string{
		"text", "longtext", "integer", "boolean", "date", "url", "json",
		"select", "multiselect", "object", "multiobject",
		"multiple", "selection", // legacy types
	}

	for _, validType := range validTypes {
		if value == validType {
			return
		}
	}

	response.Diagnostics.AddAttributeError(
		request.Path,
		"Invalid Custom Field Type",
		fmt.Sprintf("Custom field type '%s' is not valid. Must be one of: %s", value, strings.Join(validTypes, ", ")),
	)
}

// ValidCustomFieldType returns a validator for custom field types.
func ValidCustomFieldType() validator.String {
	return CustomFieldTypeValidator{}
}

// SimpleValidCustomFieldValue returns a basic validator for custom field values (without type checking).
func SimpleValidCustomFieldValue() validator.String {
	return ValidCustomFieldValue("text") // Default to text validation
}
