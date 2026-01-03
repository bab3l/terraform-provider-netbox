# Import State Preservation Fix

## Problem Statement

When importing resources into Terraform with a minimal configuration (only required fields), optional fields that exist in NetBox are not preserved in state. This causes Terraform to want to remove/modify those fields on the next apply, even though the user didn't intend to change them.

## Root Cause

The helper functions in `internal/utils/state_helpers.go` have logic that says:

```go
// If current state is null, keep it null (user didn't configure this attribute)
if current.IsNull() {
    return current
}
```

During import, the state starts completely empty (all fields null). This causes the helpers to skip populating fields from the API response.

## Solution Pattern

We need **two different behaviors**:
1. **During Import/Read**: Always populate from API response (unconditionally)
2. **During Update/Refresh**: Preserve user intent (keep null if user didn't configure)

### Option 1: Context-Aware Helpers (Recommended)

Add new helper functions that always populate from API, to be used specifically during import:

```go
// StringFromAPIAlways maps an API string value to a Terraform types.String.
// Unlike StringFromAPI, this ALWAYS sets the value from the API response,
// even if current is null. Use this during import to preserve all fields.
func StringFromAPIAlways(hasValue bool, getValue func() string) types.String {
	if hasValue {
		val := getValue()
		if val != "" {
			return types.StringValue(val)
		}
	}
	return types.StringNull()
}

// UpdateReferenceAttributeAlways updates a reference attribute, always setting from API.
// Use during import to preserve references even if not in config.
func UpdateReferenceAttributeAlways(apiName string, apiSlug string, apiID int32) types.String {
	apiIDStr := fmt.Sprintf("%d", apiID)

	// Prefer name, then slug, then ID
	if apiName != "" {
		return types.StringValue(apiName)
	}
	if apiSlug != "" {
		return types.StringValue(apiSlug)
	}
	return types.StringValue(apiIDStr)
}

// Int64FromAPIAlways maps an API int value, always setting from API.
func Int64FromAPIAlways(hasValue bool, getValue func() int64) types.Int64 {
	if hasValue {
		return types.Int64Value(getValue())
	}
	return types.Int64Null()
}

// BoolFromAPIAlways maps an API bool value, always setting from API.
func BoolFromAPIAlways(hasValue bool, getValue func() bool) types.Bool {
	if hasValue {
		return types.BoolValue(getValue())
	}
	return types.BoolNull()
}
```

### Option 2: Simplified Approach - Remove Null Checking

Simply fix the existing helpers to always populate from API:

```go
func StringFromAPI(hasValue bool, getValue func() string) types.String {
	if hasValue {
		val := getValue()
		if val != "" {
			return types.StringValue(val)
		}
	}
	return types.StringNull()
}

func UpdateReferenceAttribute(apiName string, apiSlug string, apiID int32, current types.String) types.String {
	apiIDStr := fmt.Sprintf("%d", apiID)

	// If API has no value (ID is 0), return null
	if apiID == 0 {
		return types.StringNull()
	}

	// If current state is unknown, return the ID (during initial resource creation)
	if current.IsUnknown() {
		return types.StringValue(apiIDStr)
	}

	// Check if current value matches any of the API identifiers
	if !current.IsNull() {
		val := current.ValueString()

		// Exact matches - preserve the current format
		if val == apiName || val == apiSlug || val == apiIDStr {
			return current
		}

		// Case-insensitive matches
		if (apiName != "" && strings.EqualFold(val, apiName)) ||
		   (apiSlug != "" && strings.EqualFold(val, apiSlug)) {
			return current
		}
	}

	// Default: prefer name, then slug, then ID
	if apiName != "" {
		return types.StringValue(apiName)
	}
	if apiSlug != "" {
		return types.StringValue(apiSlug)
	}
	return types.StringValue(apiIDStr)
}
```

### Option 3: Direct Mapping (No Helpers)

For resources, directly map without checking current state:

```go
// BEFORE (problematic):
if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
    data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
} else if data.Tenant.IsNull() {
    // Keep null if it was null
} else {
    data.Tenant = types.StringNull()
}

// AFTER (fixed):
if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
    data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
} else {
    data.Tenant = types.StringNull()
}
```

## Recommended Approach

**Option 2** is recommended because:
1. Less code duplication
2. Import and regular reads should both populate all fields
3. User intent is preserved through Terraform's config matching, not our helpers
4. Simpler mental model: helpers always reflect API state accurately

## Implementation Plan

1. **Update helper functions** in `internal/utils/state_helpers.go`:
   - Remove the `if current.IsNull() { return current }` logic
   - Simplify to always return API value or null

2. **Update all resource mapResponseToModel functions**:
   - Remove `else if data.Field.IsNull() { /* keep null */ }` blocks
   - Simplify to: if API has value, set it; else set null

3. **Add import preservation tests** for all resources:
   - Create with full optional fields
   - Import to minimal config
   - Verify no plan changes

4. **Test with actual import scenarios**:
   - Import existing resources from NetBox
   - Verify all fields are preserved
   - Verify no unwanted diffs

## Example: Aggregate Resource Fix

```go
func (r *AggregateResource) mapResponseToModel(ctx context.Context, aggregate *netbox.Aggregate, data *AggregateResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", aggregate.GetId()))
	data.Prefix = types.StringValue(aggregate.GetPrefix())

	// Map RIR (required)
	if rir := aggregate.GetRir(); rir.Id != 0 {
		data.RIR = utils.UpdateReferenceAttribute(data.RIR, rir.Name, rir.Slug, rir.Id)
	}

	// Map tenant (optional) - FIXED: always populate from API
	if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Map date_added (optional) - FIXED: always populate from API
	if dateAdded := aggregate.GetDateAdded(); dateAdded != "" {
		data.DateAdded = types.StringValue(dateAdded)
	} else {
		data.DateAdded = types.StringNull()
	}

	// Map description (optional) - FIXED: always populate from API
	if description, ok := aggregate.GetDescriptionOk(); ok && description != nil && *description != "" {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments (optional) - FIXED: always populate from API
	if comments, ok := aggregate.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags and custom fields (already handled correctly by existing helpers)
	// ...
}
```

## Testing Strategy

For each resource, add this test pattern:

```go
func TestAccXxxResource_importPreservesOptionalFields(t *testing.T) {
	// Step 1: Create resource with ALL optional fields
	// Step 2: Import into config with ONLY required fields
	// Step 3: PlanOnly step - should show NO changes
}
```

## Files Requiring Changes

1. `internal/utils/state_helpers.go` - Fix helper functions
2. All files in `internal/resources/*_resource.go` - Update mapResponseToModel functions
3. All files in `internal/resources_acceptance_tests/*_resource_test.go` - Add import preservation tests

## Migration Strategy

1. Fix helpers first (smallest, most impactful change)
2. Run existing tests to identify broken resources
3. Fix resources one batch at a time
4. Add import preservation tests
5. Validate with real NetBox instance imports
