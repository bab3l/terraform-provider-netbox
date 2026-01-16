# Custom Link Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: Custom Link
- **Test File**: `internal/resources_acceptance_tests/custom_link_resource_test.go`
- **Completion Date**: 2025-01-16
- **Status**: ✅ COMPLETE

## Test Structure Analysis

### Original Test Count
- **Before Standardization**: 7 tests (including IDPreservation)
- **After Standardization**: 6 tests

### Test Categories

#### Core CRUD Tests (3)
- ✅ `TestAccCustomLinkResource_basic` - Basic resource creation with required fields
- ✅ `TestAccCustomLinkResource_full` - Complete resource with all optional fields
- ✅ `TestAccCustomLinkResource_update` - Update name and other attributes

#### Reliability Tests (2)
- ✅ `TestAccCustomLinkResource_externalDeletion` - Handles external resource deletion
- ✅ `TestAccCustomLinkResource_removeOptionalFields` - Handles removal of optional fields

#### Validation Tests (1)
- ✅ `TestAccCustomLinkResource_validationErrors` - API validation errors (4 subtests: missing name, link_text, link_url, object_types)

#### Special Tests (1)
- ✅ `TestAccConsistency_CustomLink_LiteralNames` - Plan consistency with literal names

## Changes Made

### 1. Removed IDPreservation Test ✅
- **Removed**: `TestAccCustomLinkResource_IDPreservation`
- **Reason**: Duplicate of basic test functionality
- **Lines Removed**: ~30 lines

### 2. No Tag Tests Added
- **Reason**: Custom Link resource does not support tags
- **Verification**: No tag-related code found in test file

## Test Execution Results

### All Tests Passing ✅
```
TestAccCustomLinkResource_basic                        PASS (1.95s)
TestAccCustomLinkResource_full                         PASS (1.42s)
TestAccCustomLinkResource_update                       PASS (2.40s)
TestAccCustomLinkResource_externalDeletion             PASS (1.49s)
TestAccCustomLinkResource_removeOptionalFields         PASS (1.88s)
TestAccCustomLinkResource_validationErrors             PASS (1.09s)

Total: 6 tests, ~3.7 seconds
```

## Technical Details

### Dependencies
- None (standalone resource)

### Key Attributes
- `name`: Custom link name (required)
- `object_types`: List of content types this link applies to (required)
- `link_text`: Display text for the link (required)
- `link_url`: URL template for the link (required)
- `weight`: Display order weight (optional)
- `group_name`: Group name for organizing links (optional)
- `button_class`: CSS class for button styling (optional)
- `new_window`: Whether to open in new window (optional)
- `enabled`: Whether link is enabled (optional)

### Resource Characteristics
- Creates custom links in NetBox UI for specific object types
- Does not support tagging
- Uses Jinja2 templates for dynamic URL generation
- Supports grouping and ordering via weight
- Can be enabled/disabled without deletion

## Compliance Checklist

- ✅ IDPreservation test removed
- ✅ No tag tests needed (resource doesn't support tags)
- ✅ All core CRUD operations tested
- ✅ Reliability tests present (externalDeletion, removeOptionalFields)
- ✅ Validation tests present with comprehensive subtests
- ✅ All tests passing
- ✅ Cleanup registered for all resources
- ✅ Special consistency test for literal names

## Notes
- Custom Link is an extensibility feature allowing admins to add custom links to NetBox object views
- Does not support tagging
- Uses Jinja2 templating for dynamic URL generation
- Object types are specified as content type identifiers (e.g., "dcim.device")
- Special consistency test validates literal name handling
