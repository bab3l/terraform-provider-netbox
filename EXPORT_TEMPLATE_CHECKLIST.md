# Export Template Resource Standardization Checklist

**Resource:** `netbox_export_template`
**Test File:** `internal/resources_acceptance_tests/export_template_resource_test.go`
**Completed:** 2026-01-16

## Changes Made

### ✅ Removed IDPreservation Test
- **Lines removed:** 28 lines (lines 57-84)
- **Test name:** `TestAccExportTemplateResource_IDPreservation`
- **Reason:** Redundant with basic test's import verification

### ❌ No Tag Support
- Export Template resource does not support tags (no tags field in schema)
- No tag lifecycle or order tests added

## Test Results

### Test Execution
```
7 tests PASSED (~3.5s total)
- TestAccExportTemplateResource_basic
- TestAccExportTemplateResource_full
- TestAccExportTemplateResource_update
- TestAccExportTemplateResource_externalDeletion
- TestAccExportTemplateResource_removeOptionalFields
- TestAccExportTemplateResource_removeOptionalFields_extended
- TestAccExportTemplateResource_validationErrors (with 3 subtests)
```

### Test Coverage
- ✅ Basic CRUD operations
- ✅ Update validation
- ✅ Import state verification
- ✅ External deletion handling
- ✅ Optional field removal
- ✅ Validation error handling
- ❌ No tag tests (not supported)

## Resource Details

**Primary Fields:**
- `object_types` (Set, Required) - Object types this template applies to
- `name` (String, Required) - Template name (max 100 chars)
- `template_code` (String, Required) - Jinja2 template code
- `mime_type` (String, Optional) - MIME type for output
- `file_extension` (String, Optional) - File extension
- `as_attachment` (Bool, Optional, Default: true) - Download as attachment
- `description` (String, Optional) - Description

**Dependencies:** None

**Special Considerations:**
- Template code uses Jinja2 syntax
- Objects are passed as `queryset` context variable
- Default MIME type: `text/plain; charset=utf-8`

## Commit Information

**Files Modified:**
- `internal/resources_acceptance_tests/export_template_resource_test.go` (-28 lines)

**Commit Message:**
```
Standardize Export Template resource tests - remove IDPreservation test

- Remove redundant TestAccExportTemplateResource_IDPreservation test (28 lines)
- Import verification already covered in basic test
- All 7 tests passing (~3.5s)
- No tag support (no tags field in schema)

Resource 31/86 complete (36.0%)
```
