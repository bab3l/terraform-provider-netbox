# Terraform Provider Netbox - Refactoring Plan

## Overview

This document tracks the progress of refactoring resources and datasources to use common helper functions, reducing code duplication and improving maintainability.

**Branch**: `refactor/extract-common-helpers`
**Start Date**: December 26, 2025
**Status**: 🟡 In Progress

---

## Goals

1. **Reduce Code Duplication**: Extract repetitive patterns into reusable helper functions
2. **Improve Maintainability**: Centralize common logic for easier updates
3. **Maintain Stability**: All existing tests must continue to pass
4. **Enable Future Development**: Make adding new resources easier

---

## Phase 1: Low-Risk, High-Impact Helpers

**Status**: 🟡 In Progress

### 1.1 PreserveReferenceFormat() Helper

**Target**: Eliminate ~300 duplicate code blocks for reference field handling

**Before** (cluster_resource.go, lines 211-230):
```go
clusterTypeID := fmt.Sprintf("%d", cluster.Type.GetId())
clusterTypeName := cluster.Type.GetName()
clusterTypeSlug := cluster.Type.GetSlug()

if !data.Type.IsNull() && !data.Type.IsUnknown() {
    configuredValue := data.Type.ValueString()
    switch configuredValue {
    case clusterTypeID:
        data.Type = types.StringValue(clusterTypeID)
    case clusterTypeSlug:
        data.Type = types.StringValue(clusterTypeSlug)
    default:
        data.Type = types.StringValue(clusterTypeName)
    }
} else {
    data.Type = types.StringValue(clusterTypeName)
}
```

**After**:
```go
data.Type = utils.PreserveReferenceFormat(data.Type, cluster.Type.GetId(), cluster.Type.GetName(), cluster.Type.GetSlug())
```

**Progress**:
- [x] Add `PreserveReferenceFormat()` to `state_helpers.go`
- [x] Add `PreserveOptionalReferenceFormat()` for nullable references
- [x] Add unit tests for the helpers (14 test cases)
- [x] Refactor pilot resource (cluster_resource.go)
- [ ] Run acceptance tests to validate
- [ ] Apply to remaining resources (53 resources not yet using helpers)

### 1.2 PopulateTagsFromNestedTags() Helper

**Target**: Simplify tags handling from 20 lines to 1 line per resource

**Before**:
```go
if cluster.HasTags() {
    tags := utils.NestedTagsToTagModels(cluster.GetTags())
    tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
    diags.Append(tagDiags...)
    if diags.HasError() {
        return
    }
    data.Tags = tagsValue
} else {
    data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
}
```

**After**:
```go
data.Tags = utils.PopulateTagsFromNestedTags(ctx, cluster.HasTags(), cluster.GetTags(), diags)
```

**Progress**:
- [x] Add `PopulateTagsFromNestedTags()` to `state_helpers.go`
- [x] Refactor pilot resource (cluster_resource.go)
- [ ] Add unit tests for the helper
- [ ] Run acceptance tests to validate
- [ ] Apply to remaining resources

### 1.3 PopulateCustomFieldsFromMap() Helper

**Target**: Simplify custom fields handling

**Before**:
```go
if cluster.HasCustomFields() && !data.CustomFields.IsNull() {
    var stateCustomFields []utils.CustomFieldModel
    cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
    diags.Append(cfDiags...)
    if diags.HasError() {
        return
    }
    customFields := utils.MapToCustomFieldModels(cluster.GetCustomFields(), stateCustomFields)
    customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
    diags.Append(cfValueDiags...)
    if diags.HasError() {
        return
    }
    data.CustomFields = customFieldsValue
} else if data.CustomFields.IsNull() {
    data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
}
```

**After**:
```go
data.CustomFields = utils.PopulateCustomFieldsFromMap(ctx, cluster.HasCustomFields(), cluster.GetCustomFields(), data.CustomFields, diags)
```

**Progress**:
- [x] Add `PopulateCustomFieldsFromMap()` to `state_helpers.go`
- [x] Refactor pilot resource (cluster_resource.go)
- [ ] Add unit tests for the helper
- [ ] Run acceptance tests to validate
- [ ] Apply to remaining resources
- [ ] Refactor pilot resource (cluster_resource.go)
- [ ] Run acceptance tests to validate
- [ ] Apply to remaining resources

---

## Phase 2: Reference Resolution Helpers

**Status**: ⬜ Not Started

### 2.1 Generic Reference Resolver for Create/Update

**Target**: Standardize reference lookup in request building

---

## Phase 3: Schema Composition

**Status**: ⬜ Not Started

### 3.1 Common Resource Attributes Helper

**Target**: Compose schemas from reusable attribute sets

---

## Testing Strategy

### Unit Tests
- All new helpers must have comprehensive unit tests
- Test edge cases: null, unknown, empty string, various input formats
- Location: `internal/utils/state_helpers_test.go`

### Acceptance Tests
- Run existing acceptance tests after each refactoring batch
- No new acceptance tests required (helpers don't change behavior)
- Command: `go test ./internal/resources_acceptance_tests/... -run <ResourceName> -v`

---

## Commit History

| Date | Commit | Description |
|------|--------|-------------|
| 2025-12-26 | 5e13472 | Initial plan document |
| 2025-12-26 | - | Add helpers and refactor cluster_resource.go |

---

## Metrics

### Code Reduction (Estimated vs Actual)

| Pattern | Before (lines) | After (lines) | Savings | Status |
|---------|----------------|---------------|---------|--------|
| PreserveReferenceFormat (required) | 20 | 1 | 95% | ✅ Implemented |
| PreserveOptionalReferenceFormat | 15 | 1 | 93% | ✅ Implemented |
| PopulateTagsFromNestedTags | 14 | 1 | 93% | ✅ Implemented |
| PopulateCustomFieldsFromMap | 20 | 1 | 95% | ✅ Implemented |

### cluster_resource.go mapClusterToState()

| Metric | Before | After |
|--------|--------|-------|
| Lines of code | ~140 | ~55 |
| Reduction | - | **61%** |

---

## Rollback Plan

If issues are discovered:
1. Identify the problematic commit using git bisect
2. Revert the commit: `git revert <commit-hash>`
3. Investigate and fix before re-applying

---

## Notes

- Keep changes atomic - one helper at a time
- Test thoroughly before moving to next phase
- Document any API differences discovered during refactoring
