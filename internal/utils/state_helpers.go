// Package utils provides utility functions for working with Netbox provider data structures.
package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// =====================================================
// STATE MAPPING HELPERS
// =====================================================
// These helpers reduce boilerplate in mapXToState functions by providing
// consistent patterns for handling optional/nullable API responses.

// StringFromAPI maps an API string value to a Terraform types.String.
// Use this for simple string fields that are always present when HasField returns true.
//
// Example:
//
//	data.Name = StringFromAPI(device.HasName(), device.GetName, data.Name)
func StringFromAPI(hasValue bool, getValue func() string, current types.String) types.String {
	if hasValue {
		val := getValue()
		if val != "" {
			return types.StringValue(val)
		}
	}
	// If API has no value or empty string, keep null if already null, otherwise set null
	if !current.IsNull() {
		return types.StringNull()
	}
	return current
}

// StringFromAPIPreserveEmpty maps an API string value, but keeps empty strings as values.
// Use this when an empty string is semantically different from null.
func StringFromAPIPreserveEmpty(hasValue bool, getValue func() string, current types.String) types.String {
	if hasValue {
		return types.StringValue(getValue())
	}
	if !current.IsNull() {
		return types.StringNull()
	}
	return current
}

// NullableStringFromAPI maps a nullable API string pointer to a Terraform types.String.
// Use this for fields that use nullable wrappers like NullableString in the API.
//
// Example:
//
//	data.AssetTag = NullableStringFromAPI(
//	    device.HasAssetTag() && device.AssetTag.Get() != nil,
//	    func() string { return *device.AssetTag.Get() },
//	    data.AssetTag,
//	)
func NullableStringFromAPI(hasValue bool, getValue func() string, current types.String) types.String {
	if hasValue {
		val := getValue()
		if val != "" {
			return types.StringValue(val)
		}
	}
	if !current.IsNull() {
		return types.StringNull()
	}
	return current
}

// Int64FromAPI maps an API int value to a Terraform types.Int64.
// Use this for optional integer fields.
func Int64FromAPI(hasValue bool, getValue func() int64, current types.Int64) types.Int64 {
	if hasValue {
		return types.Int64Value(getValue())
	}
	if !current.IsNull() {
		return types.Int64Null()
	}
	return current
}

// Int64FromInt32API maps an API int32 value to a Terraform types.Int64.
// Use this for optional integer fields that come from the API as int32.
func Int64FromInt32API(hasValue bool, getValue func() int32, current types.Int64) types.Int64 {
	if hasValue {
		return types.Int64Value(int64(getValue()))
	}
	if !current.IsNull() {
		return types.Int64Null()
	}
	return current
}

// NullableInt64FromAPI maps a nullable API int pointer to a Terraform types.Int64.
// Use this for fields that use nullable wrappers in the API.
func NullableInt64FromAPI(hasValue bool, getValue func() *int32, current types.Int64) types.Int64 {
	if hasValue {
		ptr := getValue()
		if ptr != nil {
			return types.Int64Value(int64(*ptr))
		}
	}
	if !current.IsNull() {
		return types.Int64Null()
	}
	return current
}

// Float64FromAPI maps an API float value to a Terraform types.Float64.
// Use this for optional float fields.
func Float64FromAPI(hasValue bool, getValue func() float64, current types.Float64) types.Float64 {
	if hasValue {
		return types.Float64Value(getValue())
	}
	if !current.IsNull() {
		return types.Float64Null()
	}
	return current
}

// NullableFloat64FromAPI maps a nullable API float pointer to a Terraform types.Float64.
// Use this for fields that use nullable wrappers in the API.
func NullableFloat64FromAPI(hasValue bool, getValue func() *float64, current types.Float64) types.Float64 {
	if hasValue {
		ptr := getValue()
		if ptr != nil {
			return types.Float64Value(*ptr)
		}
	}
	if !current.IsNull() {
		return types.Float64Null()
	}
	return current
}

// BoolFromAPI maps an API bool value to a Terraform types.Bool.
// Use this for optional boolean fields.
func BoolFromAPI(hasValue bool, getValue func() bool, current types.Bool) types.Bool {
	if hasValue {
		return types.BoolValue(getValue())
	}
	if !current.IsNull() {
		return types.BoolNull()
	}
	return current
}

// =====================================================
// REFERENCE FIELD HELPERS
// =====================================================
// These helpers handle fields that reference other Netbox objects.

// ReferenceIDFromAPI maps a referenced object's ID to a Terraform types.String.
// Use this for foreign key fields where we want to preserve the user's original
// input (ID or slug) when possible.
//
// Example:
//
//	data.Tenant = ReferenceIDFromAPI(
//	    device.HasTenant() && device.Tenant.Get() != nil,
//	    func() int32 { return device.Tenant.Get().GetId() },
//	    data.Tenant,
//	)
func ReferenceIDFromAPI(hasValue bool, getID func() int32, current types.String) types.String {
	if hasValue {
		id := getID()
		if id != 0 {
			// Only update if current is null/unknown (preserve user's original input)
			if current.IsNull() || current.IsUnknown() {
				return types.StringValue(fmt.Sprintf("%d", id))
			}
			return current
		}
	}
	// No value from API
	if !current.IsNull() && !current.IsUnknown() {
		// User had a value but API returned null - this shouldn't normally happen
		// Keep the current value and let Terraform detect the drift
		return current
	}
	return types.StringNull()
}

// RequiredReferenceIDFromAPI maps a required referenced object's ID to a Terraform types.String.
// Use this for required foreign key fields.
func RequiredReferenceIDFromAPI(getID func() int32, current types.String) types.String {
	// For required fields, preserve user's input if they provided one
	if current.IsNull() || current.IsUnknown() {
		return types.StringValue(fmt.Sprintf("%d", getID()))
	}
	return current
}

// =====================================================
// ENUM/STATUS FIELD HELPERS
// =====================================================

// EnumFromAPI maps an API enum value to a Terraform types.String.
// Use this for status/enum fields that have Value() methods.
//
// Example:
//
//	data.Status = EnumFromAPI(device.HasStatus() && device.Status != nil, device.Status.GetValue)
func EnumFromAPI[T ~string](hasValue bool, getValue func() T) types.String {
	if hasValue {
		return types.StringValue(string(getValue()))
	}
	return types.StringNull()
}

// EnumFromAPIWithDefault maps an API enum value to a Terraform types.String,
// preserving the current value if the API returns nothing.
func EnumFromAPIWithDefault[T ~string](hasValue bool, getValue func() T, current types.String) types.String {
	if hasValue {
		return types.StringValue(string(getValue()))
	}
	return current
}

// =====================================================
// TAGS AND CUSTOM FIELDS HELPERS
// =====================================================

// TagsFromAPI converts API tags to a Terraform Set value.
// Returns a null set if the API has no tags.
func TagsFromAPI(ctx context.Context, hasTags bool, getTags func() []interface {
	GetName() string
	GetSlug() string
}, diags *diag.Diagnostics) types.Set {
	if hasTags {
		tags := make([]TagModel, 0, len(getTags()))
		for _, tag := range getTags() {
			tags = append(tags, TagModel{
				Name: types.StringValue(tag.GetName()),
				Slug: types.StringValue(tag.GetSlug()),
			})
		}
		tagsValue, tagDiags := types.SetValueFrom(ctx, GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return types.SetNull(GetTagsAttributeType().ElemType)
		}
		return tagsValue
	}
	return types.SetNull(GetTagsAttributeType().ElemType)
}

// CustomFieldsFromAPI converts API custom fields to a Terraform Set value.
// Uses the stateCustomFields to preserve type information.
func CustomFieldsFromAPI(ctx context.Context, hasCustomFields bool, getCustomFields func() map[string]interface{}, stateCustomFields types.Set, diags *diag.Diagnostics) types.Set {
	if hasCustomFields && !stateCustomFields.IsNull() {
		var existingFields []CustomFieldModel
		cfDiags := stateCustomFields.ElementsAs(ctx, &existingFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return stateCustomFields
		}

		customFields := MapToCustomFieldModels(getCustomFields(), existingFields)
		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfValueDiags...)
		if diags.HasError() {
			return stateCustomFields
		}
		return customFieldsValue
	}
	if stateCustomFields.IsNull() {
		return types.SetNull(GetCustomFieldsAttributeType().ElemType)
	}
	return stateCustomFields
}

// =====================================================
// REQUEST BUILDING HELPERS
// =====================================================
// These helpers reduce boilerplate in Create/Update methods.

// IsSet returns true if the value is not null and not unknown.
// Use this for conditional field setting in Create/Update methods.
//
// Example:
//
//	if IsSet(data.Description) {
//	    request.Description = data.Description.ValueStringPointer()
//	}
func IsSet(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// StringPtr returns a pointer to the string value if set, nil otherwise.
// Use this for optional string fields in API requests.
func StringPtr(value types.String) *string {
	if IsSet(value) {
		v := value.ValueString()
		return &v
	}
	return nil
}

// Int32Ptr returns a pointer to the int32 value if set, nil otherwise.
// Use this for optional integer fields in API requests.
func Int32Ptr(value types.Int64) *int32 {
	if IsSet(value) {
		v := int32(value.ValueInt64())
		return &v
	}
	return nil
}

// Int32Value returns the int32 value, or 0 if not set.
func Int32Value(value types.Int64) int32 {
	if IsSet(value) {
		return int32(value.ValueInt64())
	}
	return 0
}

// Float64Ptr returns a pointer to the float64 value if set, nil otherwise.
func Float64Ptr(value types.Float64) *float64 {
	if IsSet(value) {
		v := value.ValueFloat64()
		return &v
	}
	return nil
}

// BoolPtr returns a pointer to the bool value if set, nil otherwise.
func BoolPtr(value types.Bool) *bool {
	if IsSet(value) {
		v := value.ValueBool()
		return &v
	}
	return nil
}

// ParseInt32 parses a types.String to int32.
// Returns 0 if the value is null/unknown or cannot be parsed.
func ParseInt32(value types.String) int32 {
	if !IsSet(value) {
		return 0
	}
	var result int32
	fmt.Sscanf(value.ValueString(), "%d", &result)
	return result
}

// ParseInt32FromString parses a string to int32.
// Returns 0 if the string is empty or cannot be parsed.
func ParseInt32FromString(s string) int32 {
	if s == "" {
		return 0
	}
	var result int32
	fmt.Sscanf(s, "%d", &result)
	return result
}
