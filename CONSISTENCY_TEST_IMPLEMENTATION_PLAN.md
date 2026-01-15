# Consistency/LiteralNames Test Implementation Plan

## Overview

~~Implement Consistency/LiteralNames tests for the remaining 5 resources to achieve 100% coverage.~~ **✅ COMPLETE - 100% COVERAGE ACHIEVED!**

These tests verify that resources work correctly whether configured with literal values (slugs/names) or references to other resources, and that no unexpected plan differences occur.

**Current Status:** ✅ 99/99 complete (100%) - All resources now have Consistency/LiteralNames tests!

**Implementation Date:** 2026-01-15

**Resources Added (5):**
1. ✅ event_rule
2. ✅ notification_group
3. ✅ service_template
4. ✅ virtual_chassis
5. ✅ wireless_lan_group

**Test Execution Results:**
- All 5 tests passing
- Total execution time: 2.029s
- Clean teardown with proper cleanup handlers
- All resources properly cleaned up from Netbox

## Test Pattern

```go
func TestAccConsistency_{Resource}_LiteralNames(t *testing.T) {
	t.Parallel()

	// Create resource with unique identifiers
	name := testutil.RandomName("resource-lit")
	// ... other dependencies

	cleanup := testutil.NewCleanupResource(t)
	// Register cleanup handlers

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcc{Resource}ConsistencyLiteralNamesConfig(...),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_{resource}.test", "id"),
					resource.TestCheckResourceAttr("netbox_{resource}.test", "name", name),
					// Verify literal values are stored correctly
				),
			},
			{
				// Critical: Verify no plan drift when refreshing state
				Config:   testAcc{Resource}ConsistencyLiteralNamesConfig(...),
				PlanOnly: true,
			},
		},
	})
}
```

## Purpose and Value

**Consistency tests verify:**
1. Resources can be created using literal string values (slugs/names)
2. State correctly stores these literal values
3. No unexpected plan differences occur on refresh
4. Provider handles both reference IDs and literal names correctly

**Why This Matters:**
- **User Flexibility**: Users can choose their preferred configuration style
- **State Consistency**: Prevents drift between config and state
- **Reference Handling**: Validates provider's normalization logic
- **Production Reliability**: Ensures predictable behavior across Terraform operations

## Missing Resources (5)

### 1. event_rule
**Dependencies:**
- content_types (array of strings)
- webhook (optional reference)
- action_type, action_object_type (optional)

**Complexity:** High - complex object with many optional fields
**Estimated Time:** 30 min

### 2. notification_group
**Dependencies:**
- None (simple resource)

**Complexity:** Low - only has name
**Estimated Time:** 15 min

### 3. service_template
**Dependencies:**
- None (simple resource with protocol/ports)

**Complexity:** Low
**Estimated Time:** 15 min

### 4. virtual_chassis
**Dependencies:**
- master (device reference)
- domain (optional string)

**Complexity:** Medium - requires device setup
**Estimated Time:** 20 min

### 5. wireless_lan_group
**Dependencies:**
- parent (optional self-reference)

**Complexity:** Low
**Estimated Time:** 15 min

## Implementation Strategy

### Single Batch Implementation
- Implement all 5 tests in one session (~1.5-2 hours)
- Validate each test individually as implemented
- Single commit with all changes

### Test Validation Criteria
✅ Test passes consistently
✅ PlanOnly step shows no differences
✅ Cleanup handlers properly registered
✅ Config uses literal values (not resource references)
✅ Checks verify correct attribute values

## Common Patterns

### Simple Resources (no dependencies)
```go
func testAcc{Resource}ConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_{resource}" "test" {
  name = %q
}
`, name)
}
```

### Resources with Dependencies
```go
func testAcc{Resource}ConsistencyLiteralNamesConfig(name, depName, depSlug string) string {
	return fmt.Sprintf(`
resource "netbox_dependency" "dep" {
  name = %q
  slug = %q
}

resource "netbox_{resource}" "test" {
  name       = %q
  dependency = netbox_dependency.dep.slug  # Use slug, not ID
}
`, depName, depSlug, name)
}
```

### Self-Referencing Resources
```go
func testAcc{Resource}ConsistencyLiteralNamesConfig(parentName, parentSlug, childName, childSlug string) string {
	return fmt.Sprintf(`
resource "netbox_{resource}" "parent" {
  name = %q
  slug = %q
}

resource "netbox_{resource}" "test" {
  name   = %q
  slug   = %q
  parent = netbox_{resource}.parent.slug  # Use slug for parent reference
}
`, parentName, parentSlug, childName, childSlug)
}
```

## Progress Tracking

- [x] event_rule
- [x] notification_group
- [x] service_template
- [x] virtual_chassis
- [x] wireless_lan_group
- [x] Validate all tests pass
- [x] Update ACCEPTANCE_TEST_COVERAGE_ANALYSIS.md
- [x] Commit changes

## Success Criteria

- ✅ All 5 resources have TestAccConsistency_{Resource}_LiteralNames
- ✅ All tests passing consistently
- ✅ PlanOnly steps show no drift
- ✅ Documentation updated to 100% coverage

## Implementation Results

**All success criteria met!**

- Added Consistency tests for all 5 missing resources
- All tests pass with clean execution (2.029s total)
- Each test includes:
  - Initial creation with literal values
  - PlanOnly step verifying no drift
  - Proper cleanup handlers
- EventRule required fix: event_types value corrected from "created" to "object_created"
- Documentation updated to reflect 100% coverage
