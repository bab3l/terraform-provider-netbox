// Package utils provides utility functions for working with Netbox provider data structures.

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/bab3l/go-netbox"
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

// UpdateReferenceAttribute updates a reference attribute in the state.
// It preserves the user's input (Name or Slug) if it matches the API response.

// If the user provided an ID, or if the reference changed, it updates to the new value (preferring Name/Slug if available, or ID).

func UpdateReferenceAttribute(current types.String, apiName string, apiSlug string, apiID int32) types.String {
	apiIDStr := fmt.Sprintf("%d", apiID)

	// If current state is null, keep it null (user didn't configure this attribute)
	// This prevents drift when the API returns values for optional attributes not in config

	if current.IsNull() {
		return current
	}

	// If current state is unknown, return the ID (during initial resource creation)

	if current.IsUnknown() {
		return types.StringValue(apiIDStr)
	}

	// Check if current value matches any of the API identifiers

	val := current.ValueString()

	// Exact matches - preserve the current format

	if val == apiName {
		return current
	}

	if val == apiSlug {
		return current
	}

	if val == apiIDStr {
		return current
	}

	// Case-insensitive name match - preserve the user's casing

	if apiName != "" && strings.EqualFold(val, apiName) {
		return current
	}

	// Case-insensitive slug match - preserve the user's casing

	if apiSlug != "" && strings.EqualFold(val, apiSlug) {
		return current
	}

	// If current value is not numeric but API name/slug exist, it might be an old name
	// that still resolves to the same resource ID. Keep the current value to avoid

	// unnecessary plan diffs. The lookup functions will validate it on next apply.

	if _, err := strconv.ParseInt(val, 10, 32); err != nil {
		// Current is not a number, keep it (it's likely a name/slug)

		return current
	}

	// Current is a numeric ID but doesn't match the API ID - this is actual drift
	// Update to the correct ID

	return types.StringValue(apiIDStr)
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
// REFERENCE FIELD HELPERS

// =====================================================
// PreserveReferenceFormat preserves the user's configured format (ID, name, or slug) for reference fields.
// This is a simpler alternative to UpdateReferenceAttribute for required reference fields
// that always have a value from the API.
//
// Parameters:
//   - stateValue: The current state value (user's configured format)
//   - apiID: The ID from the API response
//   - apiName: The name/display from the API response
//   - apiSlug: The slug from the API response (can be empty if not applicable)
//
// Returns the value that preserves the user's format to prevent unnecessary plan diffs.
//
// Example:
//
//	data.Type = utils.PreserveReferenceFormat(data.Type, cluster.Type.GetId(), cluster.Type.GetName(), cluster.Type.GetSlug())
func PreserveReferenceFormat(stateValue types.String, apiID int32, apiName, apiSlug string) types.String {
	apiIDStr := fmt.Sprintf("%d", apiID)

	// If state is null or unknown, return the name (default for new/imported resources)
	if stateValue.IsNull() || stateValue.IsUnknown() {
		return types.StringValue(apiName)
	}

	// Check if the configured value matches any API identifier
	configValue := stateValue.ValueString()

	// Exact matches - preserve the user's format
	if configValue == apiIDStr {
		return types.StringValue(apiIDStr)
	}
	if configValue == apiSlug && apiSlug != "" {
		return types.StringValue(apiSlug)
	}
	if configValue == apiName {
		return types.StringValue(apiName)
	}

	// Case-insensitive matches for name/slug - preserve user's casing
	if apiName != "" && strings.EqualFold(configValue, apiName) {
		return stateValue
	}
	if apiSlug != "" && strings.EqualFold(configValue, apiSlug) {
		return stateValue
	}

	// Default to name for any other case (reference resolved but format changed)
	return types.StringValue(apiName)
}

// PreserveOptionalReferenceFormat is like PreserveReferenceFormat but handles nullable references.
// It returns a null string if the API value indicates no reference is set.
//
// Example:
//
//	data.Group = utils.PreserveOptionalReferenceFormat(data.Group, group.IsSet() && group.Get() != nil, ...)
func PreserveOptionalReferenceFormat(stateValue types.String, hasValue bool, apiID int32, apiName, apiSlug string) types.String {
	if !hasValue {
		return types.StringNull()
	}
	return PreserveReferenceFormat(stateValue, apiID, apiName, apiSlug)
}

// ReferenceWithID holds both the reference field value and the computed ID field value.
// This is used by PreserveOptionalReferenceWithID for the dual-field pattern.
type ReferenceWithID struct {
	Reference types.String
	ID        types.String
}

// PreserveOptionalReferenceWithID handles the dual-field pattern where a resource has both
// a Reference field (e.g., Tenant) and a computed ReferenceID field (e.g., TenantID).
// This is common in older resources or resources that expose the ID separately.
//
// The Reference field preserves user input format (ID/name/slug), while the ID field
// always contains the numeric ID for computed display.
//
// Parameters:
//   - stateValue: The current state value for the reference field
//   - hasValue: Whether the API returned a valid reference
//   - apiID: The ID from the API response
//   - apiName: The name/display from the API response
//   - apiSlug: The slug from the API response (can be empty)
//
// Returns a ReferenceWithID struct containing both field values.
//
// Example:
//
//	result := utils.PreserveOptionalReferenceWithID(data.Tenant, site.HasTenant(), tenant.GetId(), tenant.GetName(), tenant.GetSlug())
//	data.Tenant = result.Reference
//	data.TenantID = result.ID
func PreserveOptionalReferenceWithID(stateValue types.String, hasValue bool, apiID int32, apiName, apiSlug string) ReferenceWithID {
	if !hasValue {
		return ReferenceWithID{
			Reference: types.StringNull(),
			ID:        types.StringNull(),
		}
	}

	return ReferenceWithID{
		Reference: PreserveReferenceFormat(stateValue, apiID, apiName, apiSlug),
		ID:        types.StringValue(fmt.Sprintf("%d", apiID)),
	}
}

// =====================================================
// TAGS AND CUSTOM FIELDS HELPERS

// =====================================================
// TagsFromAPI converts API tags to a Terraform Set value.
// Returns a null set if the API has no tags.
//
// Deprecated: Use PopulateTagsFromNestedTags instead which provides a cleaner API.
func TagsFromAPI(ctx context.Context, hasTags bool, getTags func() []netbox.NestedTag, diags *diag.Diagnostics) types.Set {
	if hasTags {
		tags := NestedTagsToTagModels(getTags())

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

// PopulateTagsFromNestedTags converts Netbox NestedTag slice to a Terraform Set value.
// This is the preferred helper for resources using []netbox.NestedTag.
//
// Example:
//
//	data.Tags = utils.PopulateTagsFromNestedTags(ctx, cluster.HasTags(), cluster.GetTags(), &diags)
func PopulateTagsFromNestedTags(ctx context.Context, hasTags bool, tags []netbox.NestedTag, diags *diag.Diagnostics) types.Set {
	if !hasTags || len(tags) == 0 {
		return types.SetNull(GetTagsAttributeType().ElemType)
	}

	tagModels := NestedTagsToTagModels(tags)

	tagsValue, tagDiags := types.SetValueFrom(ctx, GetTagsAttributeType().ElemType, tagModels)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return types.SetNull(GetTagsAttributeType().ElemType)
	}

	return tagsValue
}

// PopulateCustomFieldsFromMap converts Netbox custom fields map to a Terraform Set value.
// It preserves type information from the existing state custom fields.
//
// Example:
//
//	data.CustomFields = utils.PopulateCustomFieldsFromMap(ctx, cluster.HasCustomFields(), cluster.GetCustomFields(), data.CustomFields, &diags)
func PopulateCustomFieldsFromMap(ctx context.Context, hasCustomFields bool, customFieldsMap map[string]interface{}, stateCustomFields types.Set, diags *diag.Diagnostics) types.Set {
	// If the API has no custom fields or state doesn't have custom fields configured, return null
	if !hasCustomFields || stateCustomFields.IsNull() {
		return types.SetNull(GetCustomFieldsAttributeType().ElemType)
	}

	// Extract existing state custom fields to preserve type information
	var existingFields []CustomFieldModel
	cfDiags := stateCustomFields.ElementsAs(ctx, &existingFields, false)
	diags.Append(cfDiags...)
	if diags.HasError() {
		return stateCustomFields
	}

	// Convert API custom fields using state type information
	customFields := MapToCustomFieldModels(customFieldsMap, existingFields)

	customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, GetCustomFieldsAttributeType().ElemType, customFields)
	diags.Append(cfValueDiags...)
	if diags.HasError() {
		return stateCustomFields
	}

	return customFieldsValue
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
// Use this for optional integer fields in API requests where overflow is not a concern.

// For cases where overflow checking is needed, use SafeInt32FromValue instead.

func Int32Ptr(value types.Int64) *int32 {
	if IsSet(value) {
		v := int32(value.ValueInt64()) // #nosec G115 -- Netbox IDs are within int32 range

		return &v
	}

	return nil
}

// Int32Value returns the int32 value, or 0 if not set.
// Use this for integer fields where overflow is not a concern.

// For cases where overflow checking is needed, use SafeInt32FromValue instead.

func Int32Value(value types.Int64) int32 {
	if IsSet(value) {
		return int32(value.ValueInt64()) // #nosec G115 -- Netbox IDs are within int32 range
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

	_, _ = fmt.Sscanf(value.ValueString(), "%d", &result)

	return result
}

// ParseInt32FromString parses a string to int32.
// Returns 0 if the string is empty or cannot be parsed.

func ParseInt32FromString(s string) int32 {
	if s == "" {
		return 0
	}

	var result int32

	_, _ = fmt.Sscanf(s, "%d", &result)

	return result
}

// =====================================================
// SAFE INTEGER CONVERSION HELPERS

// =====================================================
// These helpers safely convert between int64 (Terraform's standard integer type)

// and int32 (Netbox API's integer type) with overflow checking.
//

// Background: Terraform Plugin Framework uses types.Int64 as its standard integer
// type, but the Netbox API (and go-netbox client) uses int32 for IDs and most

// numeric fields. While Netbox IDs will never realistically exceed int32 range
// (~2.1 billion), we perform explicit overflow checks to satisfy security linters

// (gosec G115) and ensure robust error handling.
// SafeInt32 safely converts an int64 to int32, returning an error if the value

// would overflow. Use this when converting Terraform int64 values to Netbox API
// int32 parameters.

//
// Example:

//
//	id, err := utils.SafeInt32(data.ID.ValueInt64())

//	if err != nil {
//	    resp.Diagnostics.AddError("Invalid ID", err.Error())
//	    return

//	}
//	result, _, err := client.API.Retrieve(ctx, id).Execute()

func SafeInt32(v int64) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, fmt.Errorf("value %d overflows int32 range [%d, %d]", v, math.MinInt32, math.MaxInt32)
	}

	return int32(v), nil
}

// MustSafeInt32 safely converts an int64 to int32, panicking if the value would
// overflow. Use this only in tests or when you're certain the value is within range.

func MustSafeInt32(v int64) int32 {
	result, err := SafeInt32(v)

	if err != nil {
		panic(err)
	}

	return result
}

// SafeInt32FromValue safely extracts an int32 from a types.Int64 Terraform value.
// Returns 0 and nil error if the value is null or unknown.

// Returns an error if the value would overflow int32.
//

// Example:
//

//	weight, err := utils.SafeInt32FromValue(data.Weight)
//	if err != nil {
//	    resp.Diagnostics.AddError("Invalid weight", err.Error())
//	    return

//	}
//	if weight != 0 {
//	    req.SetWeight(weight)
//	}

func SafeInt32FromValue(v types.Int64) (int32, error) {
	if v.IsNull() || v.IsUnknown() {
		return 0, nil
	}

	return SafeInt32(v.ValueInt64())
}

// ParseID parses a string ID to int32, returning an error if parsing fails.
// This is the preferred method for parsing resource IDs in Read, Update, and Delete

// operations where invalid IDs should result in an error.
//

// Example:
//

//	id, err := utils.ParseID(data.ID.ValueString())
//	if err != nil {
//	    resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse resource ID: %s", err))
//	    return

//	}

func ParseID(idString string) (int32, error) {
	if idString == "" {
		return 0, fmt.Errorf("ID cannot be empty")
	}

	// Try parsing as int64 first to handle potential overflow gracefully

	parsed, err := strconv.ParseInt(idString, 10, 32)

	if err != nil {
		return 0, fmt.Errorf("invalid ID %q: %w", idString, err)
	}

	return int32(parsed), nil
}

// MustParseID parses a string ID to int32, panicking if parsing fails.
// Use this only in tests or when you're certain the ID is valid.

func MustParseID(idString string) int32 {
	id, err := ParseID(idString)

	if err != nil {
		panic(err)
	}

	return id
}

// ParseID64 parses a string ID to int64.
// Returns an error if the string cannot be parsed as a valid 64-bit integer.

func ParseID64(idString string) (int64, error) {
	if idString == "" {
		return 0, fmt.Errorf("ID cannot be empty")
	}

	parsed, err := strconv.ParseInt(idString, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("invalid ID %q: %w", idString, err)
	}

	return parsed, nil
}

// ToJSONString converts an interface{} to a JSON string.
// Returns an empty string if the value is nil or if serialization fails.

func ToJSONString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}

	bytes, err := json.Marshal(v)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// =====================================================
// HTTP RESPONSE HELPERS

// =====================================================
// CloseResponseBody safely closes an HTTP response body if it's not nil.

// This should be called via defer immediately after any API call that returns
// an *http.Response to prevent resource leaks.

//
// Example:

//
//	result, httpResp, err := client.API.SomeEndpoint(ctx).Execute()

//	defer utils.CloseResponseBody(httpResp)
//	if err != nil { ... }

func CloseResponseBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

// =====================================================
// REFERENCE RESOLUTION HELPERS (Create/Update operations)

// =====================================================
// These helpers standardize the pattern of looking up related resources by ID,

// name, or slug during Create/Update operations. They reduce boilerplate and
// ensure consistent error handling across all resources.

// LookupFunc is a function that resolves a reference by ID, name, or slug.
// All lookup functions in the netboxlookup package follow this signature.

//
// Example: netboxlookup.LookupClusterType, netboxlookup.LookupTenant, etc.

type LookupFunc[T any] func(ctx context.Context, client *netbox.APIClient, value string) (*T, diag.Diagnostics)

// ResolveRequiredReference resolves a required reference field during Create/Update operations.
// It calls the lookup function and appends any errors to the diagnostics.

// Returns nil if the lookup fails (diagnostics will contain the error).
//

// Example usage in buildRequest:
//

//	clusterType := utils.ResolveRequiredReference(ctx, r.client, data.Type, netboxlookup.LookupClusterType, diags)
//	if diags.HasError() {
//	    return nil
//	}

//	request.Type = *clusterType

func ResolveRequiredReference[T any](

	ctx context.Context,

	client *netbox.APIClient,

	field types.String,

	lookupFunc LookupFunc[T],

	diags *diag.Diagnostics,

) *T {
	result, lookupDiags := lookupFunc(ctx, client, field.ValueString())

	diags.Append(lookupDiags...)

	return result
}

// ResolveOptionalReference resolves an optional reference field during Create/Update operations.
// Returns nil if the field is not set or if the lookup fails.

// Any lookup errors are appended to the diagnostics.
//

// Example usage in buildRequest:
//

//	if group := utils.ResolveOptionalReference(ctx, r.client, data.Group, netboxlookup.LookupClusterGroup, diags); group != nil {
//	    request.Group = *netbox.NewNullableBriefClusterGroupRequest(group)
//	}

func ResolveOptionalReference[T any](

	ctx context.Context,

	client *netbox.APIClient,

	field types.String,

	lookupFunc LookupFunc[T],

	diags *diag.Diagnostics,

) *T {
	if !IsSet(field) {
		return nil
	}

	result, lookupDiags := lookupFunc(ctx, client, field.ValueString())

	diags.Append(lookupDiags...)

	return result
}
