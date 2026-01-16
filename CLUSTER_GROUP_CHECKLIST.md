# Cluster Group Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: netbox_cluster_group
- **Test File**: internal/resources_acceptance_tests/cluster_group_resource_test.go
- **Completion Date**: 2025-01-15

## Test Summary
- **Total Tests**: 8
- **Test Duration**: ~5.6 seconds
- **All Tests Passing**: ✅ Yes

## Test Breakdown

### Core CRUD Tests (4)
- ✅ TestAccClusterGroupResource_basic - Basic cluster group creation
- ✅ TestAccClusterGroupResource_full - Full attribute cluster group creation
- ✅ TestAccClusterGroupResource_update - Update cluster group attributes
- ✅ TestAccConsistency_ClusterGroup_LiteralNames - Import consistency check

### Reliability Tests (2)
- ✅ TestAccClusterGroupResource_externalDeletion - Handle external deletion
- ✅ TestAccClusterGroupResource_removeDescription - Remove description field

### Tag Tests (2) - Using Helper Functions
- ✅ TestAccClusterGroupResource_tagLifecycle - Complete tag lifecycle (add/update/remove)
- ✅ TestAccClusterGroupResource_tagOrderInvariance - Tag order independence

### Validation Tests (1)
- ✅ TestAccClusterGroupResource_validationErrors - Input validation (2 sub-tests)

## Changes Made

### Removed
- ❌ TestAccClusterGroupResource_IDPreservation (duplicate of import test)

### Added
- ✅ TestAccClusterGroupResource_tagLifecycle using RunTagLifecycleTest helper
- ✅ TestAccClusterGroupResource_tagOrderInvariance using RunTagOrderTest helper
- ✅ Config helper functions for nested tag format

## Technical Details

### Tag Format
Cluster Group uses **nested tag format** (like Cluster, Circuit Termination, and Circuit Type):
```hcl
tags = [
  { name = netbox_tag.tag1.name, slug = netbox_tag.tag1.slug },
  { name = netbox_tag.tag2.name, slug = netbox_tag.tag2.slug }
]
```

### Dependencies
- `netbox_tag` - For tag tests

### Test Duration Breakdown
- Basic tests: ~1.7-1.9 seconds each
- Update test: ~2.2 seconds
- Reliability tests: ~2-2.7 seconds
- Tag tests: ~4-5 seconds
- Validation test: ~0.6 seconds

## Compliance
- ✅ Follows REQUIRED_TESTS.md standard pattern (8-10 tests)
- ✅ Uses standardized tag test helpers
- ✅ Proper cleanup registration for all resources
- ✅ Parallel execution where appropriate
- ✅ Comprehensive validation coverage

## Notes
- Simple resource with only name, slug, description, and tags fields
- Fast test execution due to minimal dependencies
- Tag tests follow nested format pattern consistently
- Has a consistency check test in addition to standard import test
