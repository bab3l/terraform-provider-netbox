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
- [x] Add `PreserveOptionalReferenceWithID()` for dual-field pattern (Reference + ReferenceID)
- [x] Add unit tests for the helpers (21 test cases)
- [x] Refactor pilot resource (cluster_resource.go)
- [x] Run acceptance tests to validate
- [x] Apply to 13+ resources

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

**Status**: 🟡 In Progress

### 2.1 Generic Reference Resolver for Create/Update

**Target**: Standardize reference lookup in request building, eliminate repetitive lookup patterns

**Current Pattern** (cluster_resource.go, lines 269-276):
```go
clusterType, typeDiags := netboxlookup.LookupClusterType(ctx, r.client, data.Type.ValueString())
diags.Append(typeDiags...)
if diags.HasError() {
    return nil
}
clusterRequest := &netbox.WritableClusterRequest{
    Name: data.Name.ValueString(),
    Type: *clusterType,
}
```

**After** (with helper):
```go
clusterType := utils.ResolveRequiredReference(ctx, r.client, data.Type, netboxlookup.LookupClusterType, diags)
if diags.HasError() {
    return nil
}
clusterRequest := &netbox.WritableClusterRequest{
    Name: data.Name.ValueString(),
    Type: *clusterType,
}
```

**Optional Reference Pattern** (cluster_resource.go, lines 288-299):
```go
if utils.IsSet(data.Group) {
    group, groupDiags := netboxlookup.LookupClusterGroup(ctx, r.client, data.Group.ValueString())
    diags.Append(groupDiags...)
    if diags.HasError() {
        return nil
    }
    clusterRequest.Group = *netbox.NewNullableBriefClusterGroupRequest(group)
}
```

**After** (with helper):
```go
if group := utils.ResolveOptionalReference(ctx, r.client, data.Group, netboxlookup.LookupClusterGroup, diags); group != nil {
    clusterRequest.Group = *netbox.NewNullableBriefClusterGroupRequest(group)
}
```

**Benefits**:
- Reduces ~8 lines to ~2 lines per required reference (75% savings)
- Reduces ~12 lines to ~3 lines per optional reference (75% savings)
- Standardizes error handling across all resources
- Makes lookup function explicit and testable

**Progress**:
- [x] Add `ResolveRequiredReference()` to `state_helpers.go`
- [x] Add `ResolveOptionalReference()` to `state_helpers.go`
- [x] Add unit tests for both helpers
- [x] Refactor pilot resource (cluster_resource.go) - 26 lines saved
- [x] Run acceptance tests to validate - 12 tests passing
- [x] Apply to tenant_resource.go - 18 lines saved
- [x] Apply to site_resource.go - 36 lines saved (both Create & Update)
- [ ] Apply to remaining resources systematically
- [ ] Document final savings and update metrics

**Initial Results** (3 resources refactored):
- **cluster_resource.go**: 4 reference lookups (1 required, 3 optional) → **26 lines saved**
- **tenant_resource.go**: 2 reference lookups (both Create & Update) → **18 lines saved**
- **site_resource.go**: 6 reference lookups (both Create & Update) → **36 lines saved**
- **Total**: **80 lines saved** across 3 resources
- **Average**: ~27 lines per resource × 99 resources = **~2,670 lines potential savings**

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
| 2025-12-26 | 6d28ac5 | Add PreserveReferenceFormat and PopulateTagsFromNestedTags helpers |
| 2025-12-26 | 04280bf | Refactor cluster_resource.go |
| 2025-12-26 | c95ab48 | Refactor cluster_group and cluster_type resources |
| 2025-12-26 | 4a34e7f | Refactor tenant_resource.go |
| 2025-12-26 | a6ca219 | Refactor vrf_resource.go and ip_range_resource.go |
| 2025-12-26 | 053944c | Add PreserveOptionalReferenceWithID helper and refactor dual-field resources |
| 2025-12-26 | 43f8e21 | Complete dual-field migration for tenant and vrf resources |
| 2025-12-26 | TBD | Phase 2: Add ResolveRequiredReference and ResolveOptionalReference helpers |

---

## Metrics

### Code Reduction (Estimated vs Actual)

| Pattern | Before (lines) | After (lines) | Savings | Status |
|---------|----------------|---------------|---------|--------|
| PreserveReferenceFormat (required) | 20 | 1 | 95% | ✅ Implemented (Phase 1) |
| PreserveOptionalReferenceFormat | 15 | 1 | 93% | ✅ Implemented (Phase 1) |
| PreserveOptionalReferenceWithID (dual-field) | 28 | 5 | 82% | ✅ Implemented (Phase 1) |
| PopulateTagsFromNestedTags | 14 | 1 | 93% | ✅ Implemented (Phase 1) |
| PopulateCustomFieldsFromMap | 20 | 1 | 95% | ✅ Implemented (Phase 1) |
| ResolveRequiredReference | 10 | 2 | 80% | ✅ Implemented (Phase 2) |
| ResolveOptionalReference | 12 | 3 | 75% | ✅ Implemented (Phase 2) |

### Resources Refactored (Phase 1)

| Resource | Lines Removed (Phase 1) | Lines Removed (Phase 2) | Total Lines Removed | Tests Passing |
|----------|--------------------------|-------------------------|---------------------|---------------|
| cluster_resource.go | 148 | 26 | 174 | ✅ All |
| cluster_group_resource.go | 59 | 0 | 59 | ✅ All |
| cluster_type_resource.go | 73 | 0 | 73 | ✅ All |
| tenant_resource.go | 60 | 18 | 78 | ✅ All |
| vrf_resource.go | 78 | 0 | 78 | ✅ All |
| ip_range_resource.go | 100 | 0 | 100 | ✅ All |
| rir_resource.go | 60 | 0 | 60 | ✅ All |
| tenant_group_resource.go | 95 | 0 | 95 | ✅ All |
| region_resource.go | 98 | 0 | 98 | ✅ All |
| site_resource.go | 118 | 36 | 154 | ✅ All |
| role_resource.go | 65 | 0 | 65 | ✅ All |
| circuit_type_resource.go | 66 | 0 | 66 | ✅ All |
| asn_resource.go | 50 | 0 | 50 | ✅ All |
| **Total** | **1,070 lines** | **80 lines** | **1,150 lines** | ✅ |

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
