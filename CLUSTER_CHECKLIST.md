# Cluster Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: netbox_cluster
- **Test File**: internal/resources_acceptance_tests/cluster_resource_test.go
- **Completion Date**: 2025-01-15

## Test Summary
- **Total Tests**: 10 (plus 1 extended variant)
- **Test Duration**: ~10.9 seconds
- **All Tests Passing**: ✅ Yes

## Test Breakdown

### Core CRUD Tests (4)
- ✅ TestAccClusterResource_basic - Basic cluster creation
- ✅ TestAccClusterResource_full - Full attribute cluster creation
- ✅ TestAccClusterResource_update - Update cluster attributes
- ✅ TestAccClusterResource_import - Import existing cluster

### Reliability Tests (3)
- ✅ TestAccClusterResource_externalDeletion - Handle external deletion
- ✅ TestAccClusterResource_removeOptionalFields - Remove optional attributes
- ✅ TestAccClusterResource_removeDescriptionAndComments - Remove description and comments

### Tag Tests (2) - Using Helper Functions
- ✅ TestAccClusterResource_tagLifecycle - Complete tag lifecycle (add/update/remove)
- ✅ TestAccClusterResource_tagOrderInvariance - Tag order independence

### Validation Tests (1)
- ✅ TestAccClusterResource_validationErrors - Input validation (6 sub-tests)

### Extended Variants (1)
- ✅ TestAccClusterResource_removeOptionalFields_extended - Status field removal

## Changes Made

### Removed
- ❌ TestAccClusterResource_IDPreservation (duplicate of import test)

### Added
- ✅ TestAccClusterResource_tagLifecycle using RunTagLifecycleTest helper
- ✅ TestAccClusterResource_tagOrderInvariance using RunTagOrderTest helper
- ✅ Config helper functions for nested tag format

## Technical Details

### Tag Format
Cluster uses **nested tag format** (like Circuit Termination and Circuit Type):
```hcl
tags = [
  { name = netbox_tag.tag1.name, slug = netbox_tag.tag1.slug },
  { name = netbox_tag.tag2.name, slug = netbox_tag.tag2.slug }
]
```

### Dependencies
- `netbox_cluster_type` - Required for cluster creation
- `netbox_site` - Optional site assignment
- `netbox_tenant` - Optional tenant assignment
- `netbox_cluster_group` - Optional group assignment
- `netbox_tag` - For tag tests

### Test Duration Breakdown
- Basic tests: ~2-3 seconds each
- Update test: ~3.5 seconds
- Reliability tests: ~4-6 seconds
- Tag tests: ~4-8 seconds
- Validation test: ~3 seconds

## Compliance
- ✅ Follows REQUIRED_TESTS.md standard pattern (8-10 tests)
- ✅ Uses standardized tag test helpers
- ✅ Proper cleanup registration for all resources
- ✅ Parallel execution where appropriate
- ✅ Comprehensive validation coverage

## Notes
- Cluster has 11 total test functions (10 standard + 1 extended variant)
- Extended variant tests status field removal
- All tests use proper cleanup mechanisms
- Tag tests follow nested format pattern consistently
