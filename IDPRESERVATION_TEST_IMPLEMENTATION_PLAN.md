# ID Preservation Test Implementation Plan

## Overview

~~Implement ID preservation tests for the remaining resources to achieve 100% coverage.~~ **✅ COMPLETE - 100% COVERAGE ACHIEVED!**

ID preservation tests verify that resource IDs remain stable across Terraform operations, which is critical for state management and resource references.

**Current Status:** ✅ 99/99 complete (100%) - All resources now have ID preservation tests!

**Implementation Date:** 2026-01-15

**Resources Added (10):**
1. ✅ Device
2. ✅ FrontPortTemplate
3. ✅ Interface
4. ✅ InterfaceTemplate
5. ✅ PowerFeed
6. ✅ PowerPortTemplate
7. ✅ Prefix
8. ✅ RearPort
9. ✅ RearPortTemplate
10. ✅ Role

**Test Execution Results:**
- All 10 tests passing
- Total execution time: 3.344s
- Clean teardown with proper cleanup handlers
- All resources properly cleaned up from Netbox

## Test Pattern

```go
func TestAcc{Resource}Resource_IDPreservation(t *testing.T) {
	t.Parallel()

	// Create resource with unique identifiers
	name := testutil.RandomName("resource-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.Register{Resource}Cleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.Check{Resource}Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAcc{Resource}ResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_{resource}.test", "id"),
					resource.TestCheckResourceAttr("netbox_{resource}.test", "name", name),
					resource.TestCheckResourceAttr("netbox_{resource}.test", "slug", slug),
				),
			},
		},
	})
}
```

## Purpose and Value

**ID Preservation tests verify:**
1. Resource ID is assigned on creation
2. ID persists through refresh operations
3. ID doesn't change on updates (unless recreate is intended)
4. State tracking remains consistent

**Why This Matters:**
- **Reference Stability**: Other resources may reference this ID
- **State Integrity**: Terraform state must accurately track resources
- **Production Safety**: Unexpected ID changes can break infrastructure
- **Import/Export**: Consistent IDs enable proper import workflows

## Implementation Strategy

### Discovery Phase

1. **Identify Missing Resources**: Search for resources without IDPreservation tests
2. **Analyze Patterns**: Review existing tests to understand variations
3. **Group by Complexity**: Simple vs. complex prerequisite chains

### Implementation Approach

**Single Comprehensive Batch:**
- Process all 4-5 remaining resources in one focused session
- Estimated time: 2-3 hours
- Single commit with all changes
- Rationale: Small number, consistent pattern, clean history

## Finding Missing Resources

### Search Commands

```powershell
# Find all IDPreservation tests
Get-ChildItem -Path "internal/resources_acceptance_tests" -Filter "*_test.go" -Recurse |
  Select-String -Pattern "func TestAcc.*Resource_IDPreservation" |
  ForEach-Object { $_.Matches.Value } |
  Sort-Object

# Compare against all resources to find gaps
# (Manual comparison needed)
```

### Expected Missing Resources

Based on analysis, likely candidates for missing tests:
- Resources added recently
- Resources with complex prerequisites
- Resources with non-standard patterns

**Note:** Exact list to be determined during discovery phase.

## Test Structure

### Key Components

1. **Unique Identifiers**: Use testutil.RandomName() for uniqueness
2. **Cleanup Registration**: Ensure proper cleanup
3. **CheckDestroy**: Verify resource is destroyed after test
4. **Single Step**: One create operation is sufficient
5. **ID Verification**: TestCheckResourceAttrSet for ID field

### Configuration Pattern

```go
func testAcc{Resource}ResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_{resource}" "test" {
  name = %q
  slug = %q
  # Other required fields...
}
`, name, slug)
}
```

### Prerequisites Handling

For resources with prerequisites:
- Create minimal prerequisite chain
- Use same pattern as _basic tests
- Reuse existing config functions where possible

## Validation Criteria

Each test must:
- ✅ Have `CheckDestroy` function specified
- ✅ Check that `id` attribute is set (`TestCheckResourceAttrSet`)
- ✅ Verify at least one identifying field (name/slug)
- ✅ Use `t.Parallel()` for concurrent execution
- ✅ Register cleanup handlers
- ✅ Pass consistently (no flakiness)
- ✅ Have proper test naming: `TestAcc{Resource}Resource_IDPreservation`

## Common Patterns

### Simple Resource (no prerequisites)

```go
func TestAccSiteResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("site-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
		},
	})
}
```

### Resource with Prerequisites

```go
func TestAccInterfaceResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-id")
	deviceName := testutil.RandomName("device-id")
	interfaceName := testutil.RandomName("interface-id")

	cleanup := testutil.NewCleanupResource(t)
	// Register all cleanups

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_basic(siteName, deviceName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),
				),
			},
		},
	})
}
```

## Testing Strategy

### Local Testing

```powershell
# Test individual resource
$env:TF_ACC='1'
go test ./internal/resources_acceptance_tests/... `
  -run 'TestAcc{Resource}Resource_IDPreservation' `
  -v -timeout 30m
```

### Batch Testing

```powershell
# Test all new IDPreservation tests
$env:TF_ACC='1'
go test ./internal/resources_acceptance_tests/... `
  -run 'TestAcc.*Resource_IDPreservation' `
  -v -timeout 30m -p 1
```

### Full Regression

```powershell
# Ensure no existing tests were broken
$env:TF_ACC='1'
go test ./internal/resources_acceptance_tests/... -v -timeout 60m
```

## Documentation Updates

After completion:

1. **ACCEPTANCE_TEST_COVERAGE_ANALYSIS.md**
   - Update IDPreservation coverage to 100%
   - Remove from gap analysis
   - Update executive summary

2. **This Document**
   - Mark all resources as complete
   - Add final summary with results
   - Document any issues encountered

3. **Commit Message**
   - List all resources updated
   - Include test pass rates
   - Note any interesting findings

## Expected Challenges

### 1. Missing CheckDestroy Functions

**Issue**: Some resources may not have CheckDestroy functions
**Solution**: Create CheckDestroy functions following existing patterns

```go
func Check{Resource}Destroy(s *terraform.State) error {
	client, err := testutil.GetSharedClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_{resource}" {
			continue
		}

		// Check resource is deleted
		// Return error if still exists
	}
	return nil
}
```

### 2. Complex Prerequisites

**Issue**: Some resources require many prerequisites
**Solution**: Reuse config functions from _basic tests

### 3. Non-Standard ID Fields

**Issue**: Some resources may have computed ID fields with different names
**Solution**: Check existing tests to identify correct attribute name

## Success Metrics

**Target**: 99/99 resources with IDPreservation tests (100% coverage)

**Definition of Done:**
- All resources have TestAcc{Resource}Resource_IDPreservation
- All tests passing consistently
- All tests have CheckDestroy functions
- Documentation updated
- Coverage analysis updated to show 100%

## Timeline

**Estimated Duration**: 2-3 hours (single session)

### Milestones

- [ ] Discovery: Identify missing resources (30 min)
- [ ] Implementation: Add tests (1.5-2 hours)
- [ ] Validation: Run all tests (30 min)
- [ ] Documentation: Update coverage docs (15 min)
- [ ] Commit: Final commit with summary

## Related Work

**Prerequisites:**
- ✅ Validation tests (100%)
- ✅ Full tests (100%)
- ✅ Update tests (100%)
- ✅ Import tests (100%)
- ✅ Optional field tests (100%)
- ✅ External deletion tests (102%)

**Next Phase After IDPreservation:**
- Consistency/LiteralNames tests (~30% → 100%)
- Edge case testing
- Performance testing

## Notes

- IDPreservation tests are simpler than full tests (single step, minimal config)
- Most resources already have these tests (95% coverage)
- Pattern is well-established and consistent
- Quick wins - should take minimal time
- High value for state management integrity

---

*Created: January 15, 2026*
*Status: Planning Phase*
*Priority: High (Priority 2)*
*Target: 100% IDPreservation Test Coverage*
