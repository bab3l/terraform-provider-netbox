# Optional Fields Test Infrastructure

This document describes the test infrastructure for verifying proper null handling of optional fields in Terraform resources.

## Overview

When optional fields are removed from Terraform configuration, the provider must explicitly clear them in the API. Without proper handling, fields retain old values, causing "Provider produced inconsistent result after apply" errors.

The test helpers in `internal/testutil/optional_fields_test_helpers.go` provide reusable patterns for testing this behavior across all resources.

## Quick Start

### Basic Pattern

For simple resources with a few optional fields:

```go
func TestAccResourceName_removeOptionalFields(t *testing.T) {
    t.Parallel()

    config := testutil.OptionalFieldTestConfig{
        ResourceType: "netbox_resource",
        ResourceName: "test",
        CreateConfigWithOptionalFields: func() string {
            return `
resource "netbox_resource" "test" {
  name        = "test"
  description = "Test description"
  comments    = "Test comments"
}
`
        },
        CreateConfigWithoutOptionalFields: func() string {
            return `
resource "netbox_resource" "test" {
  name = "test"
}
`
        },
        OptionalFields: map[string]string{
            "description": "Test description",
            "comments":    "Test comments",
        },
        RequiredFields: map[string]string{
            "name": "test",
        },
    }

    testutil.TestRemoveOptionalFields(t, config)
}
```

### Multi-Step Pattern

For complex scenarios testing different combinations:

```go
func TestAccAggregate_removeOptionalFieldsMultiStep(t *testing.T) {
    t.Parallel()

    testutil.TestOptionalFieldsMultiStep(t, testutil.OptionalFieldsMultiStepConfig{
        ResourceType: "netbox_aggregate",
        ResourceName: "test",
        Steps: []testutil.OptionalFieldsTestStep{
            {
                Description: "Create with all optional fields",
                Config:      testAccAggregateConfig_allFields(),
                ExpectFields: map[string]string{
                    "description": "Test description",
                    "comments":    "Test comments",
                },
            },
            {
                Description: "Remove description only",
                Config:      testAccAggregateConfig_noDescription(),
                ExpectFields: map[string]string{
                    "comments": "Test comments",
                },
                ExpectNoFields: []string{"description"},
            },
            {
                Description:    "Remove all optional fields",
                Config:         testAccAggregateConfig_minimal(),
                ExpectNoFields: []string{"description", "comments"},
            },
            {
                Description: "Re-add optional fields",
                Config:      testAccAggregateConfig_allFields(),
                ExpectFields: map[string]string{
                    "description": "Test description",
                    "comments":    "Test comments",
                },
            },
        },
    })
}
```

## Test Coverage Checklist

When fixing a resource's null handling, ensure you:

### 1. Code Changes
- [ ] Add `else if data.FieldName.IsNull()` handling in `buildXRequest()`
- [ ] For strings: `request.SetFieldName("")`
- [ ] For nullable types: `request.SetFieldNameNil()`
- [ ] For all optional fields in the resource

### 2. Test Coverage
- [ ] Add test case using `TestRemoveOptionalFields()` helper
- [ ] Or enhance existing test with null removal step
- [ ] Verify test fails before fix (optional but recommended)
- [ ] Verify test passes after fix

### 3. Test Configuration
- [ ] Test removes ALL optional fields (not just some)
- [ ] Test re-adds fields to verify they can be set again
- [ ] Test includes proper cleanup registration

## Field Type Patterns

### String Fields (description, comments, label, etc.)

**Code:**
```go
if !data.Description.IsNull() && !data.Description.IsUnknown() {
    request.SetDescription(data.Description.ValueString())
} else if data.Description.IsNull() {
    request.SetDescription("")  // Clear with empty string
}
```

**Test:**
```go
OptionalFields: map[string]string{
    "description": "Expected description value",
    "comments":    "Expected comments value",
}
```

### Numeric Fields (port, weight, etc.)

**Code:**
```go
if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
    request.SetWeight(int32(data.Weight.ValueInt64()))
} else if data.Weight.IsNull() {
    request.SetWeightNil()  // Use Nil setter if available
}
```

### Reference Fields (tenant, site, etc.)

Most reference fields already handle this correctly:
```go
if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
    tenant, diags := lookup.LookupTenant(...)
    request.SetTenant(*tenant)
} else if data.Tenant.IsNull() {
    request.SetTenantNil()  // Already implemented
}
```

### Boolean Fields

Optional booleans are rare but follow the same pattern:
```go
if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
    request.SetEnabled(data.Enabled.ValueBool())
} else if data.Enabled.IsNull() {
    request.SetEnabledNil()
}
```

## Running Tests

### Single Resource
```bash
$env:TF_ACC='1'
$env:NETBOX_SERVER_URL='http://localhost:8000'
$env:NETBOX_API_TOKEN='0123456789abcdef0123456789abcdef01234567'
go test -v ./internal/resources_acceptance_tests -run TestAccResourceName_removeOptionalFields
```

### All Optional Field Tests
```bash
go test -v ./internal/resources_acceptance_tests -run "removeOptionalFields"
```

## Common Pitfalls

### 1. Not Clearing Empty Strings
❌ **Wrong:**
```go
if !data.Description.IsNull() {
    request.SetDescription(data.Description.ValueString())
}
// Missing: else if data.Description.IsNull() { ... }
```

✅ **Correct:**
```go
if !data.Description.IsNull() && !data.Description.IsUnknown() {
    request.SetDescription(data.Description.ValueString())
} else if data.Description.IsNull() {
    request.SetDescription("")
}
```

### 2. Testing Only Some Fields
❌ **Wrong:** Test removes `description` but not `comments`

✅ **Correct:** Test removes ALL optional fields to catch all issues

### 3. Not Testing Re-Add
Always verify fields can be re-added after removal:
```go
Steps: []resource.TestStep{
    {Config: withFields()},      // Set
    {Config: withoutFields()},   // Remove
    {Config: withFields()},      // Re-add ✓
}
```

## Example: Complete Fix

### Before (Broken)
```go
// aggregate_resource.go
func (r *AggregateResource) buildCreateRequest(...) {
    // ...
    if !data.Description.IsNull() {
        request.SetDescription(data.Description.ValueString())
    }
    // Missing null handling!
}
```

### After (Fixed)
```go
// aggregate_resource.go
func (r *AggregateResource) buildCreateRequest(...) {
    // ...
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        request.SetDescription(data.Description.ValueString())
    } else if data.Description.IsNull() {
        request.SetDescription("")
    }
}
```

### Test
```go
// aggregate_resource_test.go
func TestAccAggregate_removeOptionalFields(t *testing.T) {
    t.Parallel()

    testutil.TestRemoveOptionalFields(t, testutil.OptionalFieldTestConfig{
        ResourceType: "netbox_aggregate",
        ResourceName: "test",
        CreateConfigWithOptionalFields: func() string {
            return testAccAggregateConfig_withFields()
        },
        CreateConfigWithoutOptionalFields: func() string {
            return testAccAggregateConfig_minimal()
        },
        OptionalFields: map[string]string{
            "description": "Test description",
            "comments":    "Test comments",
        },
    })
}
```

## Integration with Existing Tests

You can also add optional field removal steps to existing tests:

```go
func TestAccResource_full(t *testing.T) {
    // ... existing test ...
    Steps: []resource.TestStep{
        {Config: fullConfig()},
        {Config: updatedConfig()},
        // Add optional field removal test
        {
            Config: minimalConfig(),
            Check: resource.ComposeTestCheckFunc(
                resource.TestCheckNoResourceAttr("netbox_resource.test", "description"),
                resource.TestCheckNoResourceAttr("netbox_resource.test", "comments"),
            ),
        },
    }
}
```

## Best Practices

1. **Always test null removal** for optional fields when fixing a resource
2. **Use the helper functions** rather than writing custom test logic
3. **Group related fields** - test description+comments together
4. **Document field types** in test comments (string vs reference vs numeric)
5. **Run tests before committing** to ensure fix works
6. **Check test output** for actual vs expected values in failures

## Troubleshooting

### Test fails with "inconsistent result"
✓ This means the fix isn't complete - check all optional fields have null handling

### Test passes but fields aren't cleared
✓ Check the `mapResponseToModel()` function - it might be mapping empty strings back to values

### Test is flaky
✓ Ensure proper cleanup with `testutil.NewCleanupResource(t)`

## Reference

- Analysis results: `BUGFIX_ANALYSIS_null_handling.csv`
- Planning document: `BUGFIX_PLAN_optional_field_null_handling.md`
- Example implementation: `internal/resources/asn_resource.go` (fixed)
- Example test: `internal/resources_acceptance_tests/asn_resource_test.go`
