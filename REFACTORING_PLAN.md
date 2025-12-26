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

**Status**: 🟡 In Progress

### 3.1 Common Resource Attributes Helpers

**Target**: Compose schemas from reusable attribute sets to reduce repetition

**Strategy Refined (Phase 3b)**: Created three-tier composition system to handle different resource patterns:

1. **CommonDescriptiveAttributes** (description + comments) - for resources with both fields
2. **DescriptionOnlyAttributes** (description only) - for resources without comments field
3. **CommonMetadataAttributes** (tags + custom_fields) - universal for all resources

**Current Pattern** - Most resources repeat these common attributes:
```go
"description": nbschema.DescriptionAttribute("resource"),
"tags": nbschema.TagsAttribute(),
"custom_fields": nbschema.CustomFieldsAttribute(),
// Some resources also have:
"comments": nbschema.CommentsAttribute("resource"),
```

**Refined Pattern** - Use appropriate helper based on resource needs:
```go
// Option 1: Resource with description + comments + tags + custom_fields
maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("resource"))
maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())

// Option 2: Resource with description + tags + custom_fields (no comments)
maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("resource"))
maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())

// Option 3: Resource with only tags + custom_fields (no description)
maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
```

**Benefits**:
- Reduces 4 lines to 2 lines (50% savings) for full set
- Reduces 3 lines to 2 lines (33% savings) for description-only
- Reduces 2 lines to 1 line (50% savings) for metadata-only
- Makes it easier to add new common attributes to all resources
- Ensures consistency across all resource schemas

**Progress**:
- [x] Create `CommonDescriptiveAttributes()` helper (description + comments)
- [x] Create `DescriptionOnlyAttributes()` helper (description only)
- [x] Create `CommonMetadataAttributes()` helper (tags + custom_fields)
- [x] Refactor Phase 3a pilot resources (cluster, tenant, site) - using CommonDescriptiveAttributes
- [x] Refactor Phase 3b resources (circuit_type, cluster_group, cluster_type, rir, region)
- [x] Validate with tests - all 56 tests passing (32 Phase 3a + 24 Phase 3b)
- [ ] Create batches for systematic rollout
- [ ] Execute remaining batches

**Initial Results**:

*Phase 3a* (3 resources - description + comments):
- **cluster_resource.go**: 4 lines → 2 lines (**2 lines saved**)
- **tenant_resource.go**: 4 lines → 2 lines (**2 lines saved**)
- **site_resource.go**: 4 lines → 2 lines (**2 lines saved**)

*Phase 3b* (5 resources - description only):
- **circuit_type_resource.go**: 2 lines → 1 line (**1 line saved**)
- **cluster_group_resource.go**: 3 lines → 2 lines (**1 line saved**)
- **cluster_type_resource.go**: 3 lines → 2 lines (**1 line saved**)
- **rir_resource.go**: 3 lines → 2 lines (**1 line saved**)
- **region_resource.go**: 3 lines → 2 lines (**1 line saved**)

**Subtotal Phase 3**: **11 lines saved** across 8 resources

### 3.2 Phase 3 Rollout Batches

**Batch Strategy**: Group resources by schema pattern for efficient refactoring

#### Batch 1: Description + Comments + Tags + Custom Fields (CommonDescriptiveAttributes + CommonMetadataAttributes)
Resources with full descriptive metadata (like cluster, tenant, site):
- [ ] contact_resource.go
- [ ] device_resource.go
- [ ] device_type_resource.go
- [ ] fhrp_group_resource.go
- [ ] ike_proposal_resource.go
- [ ] ip_address_resource.go
- [ ] ip_range_resource.go
- [ ] ipsec_policy_resource.go
- [ ] ipsec_profile_resource.go
- [ ] ipsec_proposal_resource.go
- [ ] l2vpn_resource.go
- [ ] module_resource.go
- [ ] module_type_resource.go
- [ ] power_panel_resource.go
- [ ] prefix_resource.go
- [ ] provider_account_resource.go
- [ ] provider_network_resource.go
- [ ] provider_resource.go
- [ ] rack_reservation_resource.go
- [ ] rack_type_resource.go

**Estimated savings**: ~2 lines × 20 resources = **40 lines**

#### Batch 2: Description + Tags + Custom Fields (DescriptionOnlyAttributes + CommonMetadataAttributes)
Resources with description but no comments field:
- [ ] aggregate_resource.go
- [ ] asn_range_resource.go
- [ ] cable_resource.go
- [ ] contact_assignment_resource.go
- [ ] contact_group_resource.go
- [ ] contact_role_resource.go
- [ ] device_bay_template_resource.go
- [ ] front_port_resource.go
- [ ] front_port_template_resource.go
- [ ] interface_resource.go
- [ ] interface_template_resource.go
- [ ] inventory_item_resource.go
- [ ] inventory_item_role_resource.go
- [ ] inventory_item_template_resource.go
- [ ] l2vpn_termination_resource.go
- [ ] location_resource.go
- [ ] manufacturer_resource.go
- [ ] platform_resource.go
- [ ] power_feed_resource.go
- [ ] power_outlet_resource.go
- [ ] power_outlet_template_resource.go
- [ ] power_port_resource.go
- [ ] power_port_template_resource.go
- [ ] rack_resource.go
- [ ] rack_role_resource.go
- [ ] rear_port_resource.go
- [ ] rear_port_template_resource.go
- [ ] service_resource.go
- [ ] service_template_resource.go
- [ ] site_group_resource.go
- [ ] tunnel_resource.go
- [ ] tunnel_group_resource.go
- [ ] tunnel_termination_resource.go
- [ ] virtual_chassis_resource.go
- [ ] virtual_disk_resource.go
- [ ] vlan_resource.go
- [ ] vlan_group_resource.go
- [ ] vm_interface_resource.go
- [ ] vpn_tunnel_resource.go
- [ ] wireless_lan_resource.go
- [ ] wireless_lan_group_resource.go
- [ ] wireless_link_resource.go

**Estimated savings**: ~1 line × 42 resources = **42 lines**

#### Batch 3: Tags + Custom Fields Only (CommonMetadataAttributes)
Resources with only metadata, no description or comments:
- [ ] console_port_resource.go
- [ ] console_port_template_resource.go
- [ ] console_server_port_resource.go
- [ ] console_server_port_template_resource.go
- [ ] device_bay_resource.go
- [ ] journal_entry_resource.go
- [ ] module_bay_resource.go
- [ ] module_bay_template_resource.go
- [ ] power_port_resource.go (if no description)

**Estimated savings**: ~1 line × 9 resources = **9 lines**

#### Batch 4: Special Cases
Resources with inline descriptions or other variations:
- [ ] Review and handle individually

**Total Phase 3 Potential Savings**: 11 (completed) + 91 (estimated) = **~102 lines**

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
| 2025-12-26 | f9995a4 | Phase 2: Add ResolveRequiredReference and ResolveOptionalReference helpers |
| 2025-12-26 | 37f8733 | Phase 3a: Add CommonDescriptiveAttributes and CommonMetadataAttributes helpers |
| 2025-12-26 | 8eb12a3 | Phase 3b: Add DescriptionOnlyAttributes helper and refactor 5 more resources |

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
| CommonDescriptiveAttributes | 2 | 1 | 50% | ✅ Implemented (Phase 3) |
| CommonMetadataAttributes | 2 | 1 | 50% | ✅ Implemented (Phase 3) |

### Resources Refactored (Phase 1)

| Resource | Lines Removed (Phase 1) | Lines Removed (Phase 2) | Lines Removed (Phase 3) | Total Lines Removed | Tests Passing |
|----------|--------------------------|-------------------------|-------------------------|---------------------|---------------|
| cluster_resource.go | 148 | 26 | 2 | 176 | ✅ All |
| cluster_group_resource.go | 59 | 0 | 1 | 60 | ✅ All |
| cluster_type_resource.go | 73 | 0 | 1 | 74 | ✅ All |
| tenant_resource.go | 60 | 18 | 2 | 80 | ✅ All |
| vrf_resource.go | 78 | 0 | 0 | 78 | ✅ All |
| ip_range_resource.go | 100 | 0 | 0 | 100 | ✅ All |
| rir_resource.go | 60 | 0 | 1 | 61 | ✅ All |
| tenant_group_resource.go | 95 | 0 | 0 | 95 | ✅ All |
| region_resource.go | 98 | 0 | 1 | 99 | ✅ All |
| site_resource.go | 118 | 36 | 2 | 156 | ✅ All |
| role_resource.go | 65 | 0 | 0 | 65 | ✅ All |
| circuit_type_resource.go | 66 | 0 | 1 | 67 | ✅ All |
| asn_resource.go | 50 | 0 | 0 | 50 | ✅ All |
| **Total** | **1,070 lines** | **80 lines** | **11 lines** | **1,161 lines** | ✅ |

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
