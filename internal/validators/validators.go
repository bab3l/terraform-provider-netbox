// Package validators provides custom validators for the Netbox provider.

package validators

import (
	"context"
	"encoding/json"
	"fmt"
	"net/netip"
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

// LatitudeValidator validates that a float is between -90 and 90 (inclusive).
type LatitudeValidator struct{}

func (v LatitudeValidator) Description(_ context.Context) string {
	return "value must be a valid latitude between -90 and 90"
}

func (v LatitudeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v LatitudeValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat64()
	if value < -90 || value > 90 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Latitude",
			fmt.Sprintf("Latitude must be between -90 and 90 (inclusive). Got %v", value),
		)
	}
}

// ValidLatitude returns a validator for latitude values.
func ValidLatitude() validator.Float64 {
	return LatitudeValidator{}
}

// LongitudeValidator validates that a float is between -180 and 180 (inclusive).
type LongitudeValidator struct{}

func (v LongitudeValidator) Description(_ context.Context) string {
	return "value must be a valid longitude between -180 and 180"
}

func (v LongitudeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v LongitudeValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat64()
	if value < -180 || value > 180 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Longitude",
			fmt.Sprintf("Longitude must be between -180 and 180 (inclusive). Got %v", value),
		)
	}
}

// ValidLongitude returns a validator for longitude values.
func ValidLongitude() validator.Float64 {
	return LongitudeValidator{}
}

// VLANIDInt64Validator validates that a VLAN ID is between 1 and 4094 (inclusive).
type VLANIDInt64Validator struct{}

func (v VLANIDInt64Validator) Description(_ context.Context) string {
	return "value must be a valid VLAN ID between 1 and 4094"
}

func (v VLANIDInt64Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v VLANIDInt64Validator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueInt64()
	if value < 1 || value > 4094 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid VLAN ID",
			fmt.Sprintf("VLAN ID must be between 1 and 4094 (inclusive). Got %d", value),
		)
	}
}

// ValidVLANIDInt64 returns a validator for VLAN IDs stored in int64 attributes.
func ValidVLANIDInt64() validator.Int64 {
	return VLANIDInt64Validator{}
}

// VLANIDInt32Validator validates that a VLAN ID is between 1 and 4094 (inclusive).
type VLANIDInt32Validator struct{}

func (v VLANIDInt32Validator) Description(_ context.Context) string {
	return "value must be a valid VLAN ID between 1 and 4094"
}

func (v VLANIDInt32Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v VLANIDInt32Validator) ValidateInt32(ctx context.Context, request validator.Int32Request, response *validator.Int32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueInt32()
	if value < 1 || value > 4094 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid VLAN ID",
			fmt.Sprintf("VLAN ID must be between 1 and 4094 (inclusive). Got %d", value),
		)
	}
}

// ValidVLANIDInt32 returns a validator for VLAN IDs stored in int32 attributes.
func ValidVLANIDInt32() validator.Int32 {
	return VLANIDInt32Validator{}
}

// ASNInt64Validator validates that an ASN is between 1 and 4294967295 (inclusive).
type ASNInt64Validator struct{}

func (v ASNInt64Validator) Description(_ context.Context) string {
	return "value must be a valid ASN between 1 and 4294967295"
}

func (v ASNInt64Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ASNInt64Validator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueInt64()
	if value < 1 || value > 4294967295 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid ASN",
			fmt.Sprintf("ASN must be between 1 and 4294967295 (inclusive). Got %d", value),
		)
	}
}

// ValidASNInt64 returns a validator for ASN values stored in int64 attributes.
func ValidASNInt64() validator.Int64 {
	return ASNInt64Validator{}
}

// ASNStringValidator validates that a string ASN is between 1 and 4294967295 (inclusive).
type ASNStringValidator struct{}

func (v ASNStringValidator) Description(_ context.Context) string {
	return "value must be a valid ASN between 1 and 4294967295"
}

func (v ASNStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ASNStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	asn, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid ASN",
			fmt.Sprintf("ASN must be a valid integer between 1 and 4294967295. Got '%s'", value),
		)
		return
	}

	if asn < 1 || asn > 4294967295 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid ASN",
			fmt.Sprintf("ASN must be between 1 and 4294967295 (inclusive). Got %d", asn),
		)
	}
}

// ValidASNString returns a validator for ASN values stored as strings.
func ValidASNString() validator.String {
	return ASNStringValidator{}
}

// MACAddressValidator validates MAC addresses in colon-separated hex format.
type MACAddressValidator struct{}

func (v MACAddressValidator) Description(_ context.Context) string {
	return "value must be a valid MAC address in format AA:BB:CC:DD:EE:FF"
}

func (v MACAddressValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v MACAddressValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	matched, _ := regexp.MatchString(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`, value)
	if !matched {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid MAC Address",
			fmt.Sprintf("MAC address must be in format AA:BB:CC:DD:EE:FF. Got '%s'", value),
		)
	}
}

// ValidMACAddress returns a validator for MAC address values.
func ValidMACAddress() validator.String {
	return MACAddressValidator{}
}

// IPAddressValidator validates that a string is a valid IPv4 or IPv6 address.
type IPAddressValidator struct{}

// Description returns a description of the validator.
func (v IPAddressValidator) Description(_ context.Context) string {
	return "value must be a valid IPv4 or IPv6 address"
}

// MarkdownDescription returns a markdown description of the validator.
func (v IPAddressValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v IPAddressValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	if _, err := netip.ParseAddr(value); err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid IP Address",
			fmt.Sprintf("The value '%s' is not a valid IPv4 or IPv6 address: %s", value, err),
		)
	}
}

// ValidIPAddress returns a validator which ensures the value is a valid IPv4 or IPv6 address.
func ValidIPAddress() validator.String {
	return IPAddressValidator{}
}

// IPAddressWithPrefixValidator validates that a string is a valid IP address with prefix length (CIDR).
type IPAddressWithPrefixValidator struct{}

// Description returns a description of the validator.
func (v IPAddressWithPrefixValidator) Description(_ context.Context) string {
	return "value must be a valid IPv4 or IPv6 address with prefix length (CIDR)"
}

// MarkdownDescription returns a markdown description of the validator.
func (v IPAddressWithPrefixValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v IPAddressWithPrefixValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	if _, err := netip.ParsePrefix(value); err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid IP Address with Prefix",
			fmt.Sprintf("The value '%s' is not a valid IPv4 or IPv6 address with prefix length: %s", value, err),
		)
	}
}

// ValidIPAddressWithPrefix returns a validator which ensures the value is a valid IP address with prefix length.
func ValidIPAddressWithPrefix() validator.String {
	return IPAddressWithPrefixValidator{}
}

// IPPrefixValidator validates that a string is a valid network prefix (CIDR) with a masked network address.
type IPPrefixValidator struct{}

// Description returns a description of the validator.
func (v IPPrefixValidator) Description(_ context.Context) string {
	return "value must be a valid IPv4 or IPv6 network prefix (CIDR)"
}

// MarkdownDescription returns a markdown description of the validator.
func (v IPPrefixValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v IPPrefixValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Prefix",
			fmt.Sprintf("The value '%s' is not a valid IPv4 or IPv6 network prefix: %s", value, err),
		)
		return
	}

	masked := prefix.Masked()
	if prefix.Addr() != masked.Addr() {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Prefix",
			fmt.Sprintf("The value '%s' must be a network address (host bits must be zero)", value),
		)
	}
}

// ValidIPPrefix returns a validator which ensures the value is a valid network prefix (CIDR).
func ValidIPPrefix() validator.String {
	return IPPrefixValidator{}
}
