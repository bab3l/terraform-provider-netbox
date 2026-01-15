# Required Acceptance Tests for Resources

## Standard Test Suite

Every resource should have the following acceptance tests:

### Core Tests (Required)
- ✅ `TestAcc{Resource}Resource_basic` - Minimal configuration (required fields only)
- ✅ `TestAcc{Resource}Resource_full` - Complete configuration (all optional fields)
- ✅ `TestAcc{Resource}Resource_update` - Modify existing resource
- ✅ `TestAcc{Resource}Resource_import` - Terraform import functionality

### Reliability Tests (Required)
- ✅ `TestAcc{Resource}Resource_IDPreservation` - ID stability across updates
- ✅ `TestAcc{Resource}Resource_externalDeletion` - Handle external deletion
- ✅ `TestAcc{Resource}Resource_removeOptionalFields` - Clear optional fields

### Additional Tests (Recommended)
- `TestAcc{Resource}Resource_removeOptionalFields_extended` - Comprehensive optional field testing
- `Consistency_{Resource}_LiteralNames` - Literal values vs resource references

## CheckDestroy Function

Each resource must have a corresponding CheckDestroy function in `internal/testutil/check_destroy*.go`:

```go
func Check{Resource}Destroy(s *terraform.State) error
```

## Test Helper Usage

Use helpers from `internal/testutil/` to standardize test implementations:
- `RunImportTest()` - Import testing
- `RunUpdateTest()` - Update testing
- `RunExternalDeletionTest()` - External deletion testing
- `TestRemoveOptionalFields()` - Optional field removal

## Cleanup Registration

All tests must register cleanup:

```go
cleanup := testutil.NewCleanupResource(t)
cleanup.Register{Resource}Cleanup(resourceName)
```

---

*Refer to ACCEPTANCE_TEST_COVERAGE_ANALYSIS.md for detailed gap analysis and implementation guidelines.*
