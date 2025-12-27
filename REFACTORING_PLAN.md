# Terraform Provider Netbox - Refactoring Plan

## Overview

This document tracks the progress of refactoring resources and datasources to use common helper functions, reducing code duplication and improving maintainability.

**Branch**: `refactor/extract-common-helpers`
**Start Date**: December 26, 2025
**Status**: ✅ Phase 3 Batch 4 Complete - 86 resources refactored total

---

## Progress Summary

### Completed Phases

| Phase | Description | Resources | Lines Saved | Commit(s) |
|-------|-------------|-----------|-------------|-----------|
| Phase 1 | State Mapping Helpers | 13 | 1,070 | 43f8e21 |
| Phase 2 | Reference Resolution | 3 | 80 | f9995a4 |
| Phase 3a | Schema Composition Pilot | 3 | 6 | 37f8733 |
| Phase 3b | DescriptionOnly Helper | 5 | 5 | 8eb12a3 |
| **Phase 3 Batch 1** | **Full Schema Refactor** | **20** | **81** | **6 commits** |
| **Phase 3 Batch 2** | **DescriptionOnly Refactor** | **27** | **62** | **3 commits** |
| **Phase 3 Batch 1+** | **Batch 1 Additions** | **6** | **77** | **afca597** |
| **Phase 3 Batch 3** | **Metadata-Only Refactor** | **3** | **10** | **be04836** |
| **Phase 3 Batch 4** | **Misc Description/Comments** | **6** | **13** | **e8bc4c9** |
| **TOTAL** | **All Phases** | **86** | **1,404** | **18 commits** |

### Phase 3 Batch 1 Details
- Part 1 (5 resources): 9 lines saved - commit f4612ae
- Part 2 (5 resources): 11 lines saved - commit 526e6cb
- Part 3 (4 resources): 31 lines saved - commit 91d2214
- Part 4 (4 resources): 26 lines saved - commit e282a4d
- prefix_resource: 9 lines saved - commit 8e68d01
- contact_resource: -5 lines (infrastructure) - commit acccce3

### Phase 3 Batch 2 Details
- Part 1 (4 resources): 4 lines saved - commit f4b0b32
- Part 2 (10 resources): 39 lines saved - commit f1d3048
- Part 3 (13 resources): 19 lines saved - commit e5faf8f

---

## Goals

1. **Reduce Code Duplication**: Extract repetitive patterns into reusable helper functions
2. **Improve Maintainability**: Centralize common logic for easier updates
3. **Maintain Stability**: All existing tests must continue to pass
4. **Enable Future Development**: Make adding new resources easier

---

## Phase 1: State Mapping Helpers

**Status**: ✅ Complete - 13 resources refactored, 1,070 lines saved

**Helpers Created**:
- `PreserveReferenceFormat()` - Preserves user's input format (ID/name/slug) for references
- `PreserveOptionalReferenceFormat()` - Same as above but for nullable references
- `PreserveOptionalReferenceWithID()` - Handles dual-field pattern (Reference + ReferenceID)
- `PopulateTagsFromNestedTags()` - Converts Netbox nested tags to Terraform tag models
- `PopulateCustomFieldsFromMap()` - Converts Netbox custom fields map to Terraform models

**Refactored Resources**: cluster, tenant, site, circuit_type, cluster_group, cluster_type, rir, region, device_role, vrf, asn, route_target, wireless_lan

**Commit**: 43f8e21

---

## Phase 2: Reference Resolution Helpers

**Status**: ✅ Complete - 3 resources refactored, 80 lines saved

**Helpers Created**:
- `LookupFunc[T]` - Generic type for netboxlookup functions
- `ResolveRequiredReference[T]()` - Standardized required reference lookup with error handling
- `ResolveOptionalReference[T]()` - Standardized optional reference lookup with error handling

**Refactored Resources**: cluster, tenant, site

**Commit**: f9995a4

---

## Phase 3: Schema Composition Helpers

**Status**: 🟢 Batch 1 Complete (20/71 resources) - Moving to Batch 2

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

**Status**: ✅ Batches 1-4 Complete (62/~99 resources)

### 3.1 Common Resource Attributes Helpers

**Target**: Compose schemas from reusable attribute sets to reduce repetition

**Strategy**: Three-tier composition system to handle different resource patterns:

1. **CommonDescriptiveAttributes** (description + comments) - for resources with both fields
2. **DescriptionOnlyAttributes** (description only) - for resources without comments field
3. **CommonMetadataAttributes** (tags + custom_fields) - universal for all resources

**Pattern** - Use appropriate helper based on resource needs:
```go
// Option 1: Resource with description + comments + tags + custom_fields
maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("resource"))
maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())

// Option 2: Resource with description + tags + custom_fields (no comments)
maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("resource"))
maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())

// Option 3: Resource with only tags + custom_fields (no description)
maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())

// Special Case: Resource without custom_fields (rare)
maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("resource"))
// Then add tags directly: "tags": nbschema.TagsAttribute()
```

**Benefits**:
- Reduces 4 lines to 2 lines (50% savings) for full set
- Reduces 3 lines to 2 lines (33% savings) for description-only
- Reduces 2 lines to 1 line (50% savings) for metadata-only
- Makes it easier to add new common attributes to all resources
- Ensures consistency across all resource schemas

**Completed**:
- [x] Create all three composition helpers in `internal/schema/attributes.go`
- [x] Phase 3a: Pilot resources (cluster, tenant, site) - commit 37f8733
- [x] Phase 3b: DescriptionOnly pattern (circuit_type, cluster_group, cluster_type, rir, region) - commit 8eb12a3
- [x] **Batch 1: Full descriptive resources (20 resources)** - commits f4612ae, 526e6cb, 91d2214, e282a4d, 8e68d01, acccce3
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
Resources with full descriptive metadata:
- [x] device_resource.go ✅
- [x] device_type_resource.go ✅
- [x] fhrp_group_resource.go ✅
- [x] ike_proposal_resource.go ✅
- [x] ip_address_resource.go ✅ (special case - has tags but no custom_fields)
- [x] ip_range_resource.go ✅
- [x] ipsec_policy_resource.go ✅
- [x] ipsec_profile_resource.go ✅
- [x] ipsec_proposal_resource.go ✅
- [x] l2vpn_resource.go ✅
- [x] module_resource.go ✅
- [x] module_type_resource.go ✅
- [x] power_panel_resource.go ✅
- [x] prefix_resource.go ✅ (special case - has tags but no custom_fields)
- [x] provider_account_resource.go ✅
- [x] provider_network_resource.go ✅
- [x] provider_resource.go ✅
- [x] rack_reservation_resource.go ✅
- [x] rack_type_resource.go ✅
- [x] contact_resource.go ✅ (special case - has tags but no custom_fields)

**Completed: 20/20 ✅**
**Lines saved: 81 lines** (Part 1: 9, Part 2: 11, Part 3: 31, Part 4: 26, prefix: 9, contact: -5)
**Commits**: f4612ae, 526e6cb, 91d2214, e282a4d, 8e68d01, acccce3

**Note**: Two resources (prefix, contact) lack custom_fields and use TagsAttribute() directly instead of CommonMetadataAttributes().

#### Batch 2: Description + Tags + Custom Fields (DescriptionOnlyAttributes + CommonMetadataAttributes)
Resources with description but no comments field:
- [x] asn_range_resource.go ✅
- [x] circuit_termination_resource.go ✅
- [x] circuit_type_resource.go ✅ (already in Phase 3b)
- [x] cluster_group_resource.go ✅ (already in Phase 3b)
- [x] cluster_type_resource.go ✅ (already in Phase 3b)
- [x] console_server_port_template_resource.go ✅
- [x] contact_group_resource.go ✅
- [x] contact_role_resource.go ✅
- [x] device_bay_template_resource.go ✅
- [x] device_role_resource.go ✅
- [x] front_port_template_resource.go ✅
- [x] interface_template_resource.go ✅
- [x] inventory_item_role_resource.go ✅
- [x] inventory_item_template_resource.go ✅
- [x] location_resource.go ✅
- [x] manufacturer_resource.go ✅
- [x] module_bay_template_resource.go ✅
- [x] platform_resource.go ✅
- [x] power_outlet_template_resource.go ✅
- [x] power_port_template_resource.go ✅
- [x] rack_role_resource.go ✅
- [x] rear_port_template_resource.go ✅
- [x] region_resource.go ✅ (already in Phase 3b)
- [x] rir_resource.go ✅ (already in Phase 3b)
- [x] site_group_resource.go ✅
- [x] tenant_group_resource.go ✅
- [x] tunnel_group_resource.go ✅
- [x] virtual_disk_resource.go ✅
- [x] vlan_group_resource.go ✅

**Completed: 27/27 ✅** (includes 5 from Phase 3b)
**Lines saved: 62 lines** (Part 1: 4, Part 2: 39, Part 3: 19)
**Commits**: f4b0b32, f1d3048, e5faf8f

**Note**: During audit, discovered misclassifications:
- **Moved to Batch 1** (have comments field): aggregate, cable, power_feed, rack, site, virtual_machine, vlan, vrf
- **Moved to Batch 3** (no description field): contact_assignment, fhrp_group_assignment, l2vpn_termination, tunnel_termination
- **Excluded** (no tags/custom_fields support): console_port_template, and other port templates without metadata

#### Batch 1 Additions: Description + Comments + Tags + Custom Fields (CommonDescriptiveAttributes + CommonMetadataAttributes)
Misclassified resources moved from Batch 2 that actually have comments field:
- [x] aggregate_resource.go ✅
- [x] cable_resource.go ✅
- [x] power_feed_resource.go ✅
- [x] rack_resource.go ✅
- [x] site_resource.go ✅ (already in Phase 3a)
- [x] virtual_machine_resource.go ✅
- [x] vlan_resource.go ✅
- [x] vrf_resource.go ✅ (already in Phase 1)

**Completed: 6/6 new resources ✅** (2 already done)
**Lines saved: 77 lines**
**Commit**: afca597

#### Batch 3: Tags + Custom Fields Only (CommonMetadataAttributes)
Resources with only metadata, no description or comments:
- [x] contact_assignment_resource.go ✅
- [x] fhrp_group_assignment_resource.go ✅
- [x] tunnel_termination_resource.go ✅

**Completed: 3/3 ✅**
**Lines saved: 10 lines**
**Commit**: be04836

**Note**: Original Batch 3 list was inaccurate - most listed resources actually have description fields:
- **Moved to other batches**: console_port (has description), console_server_port (has description), device_bay (has description), module_bay (has description), module_bay_template (has description), power_port (has description)
- **Has comments instead**: journal_entry (has comments field - should be Batch 1)
- **No metadata support**: l2vpn_termination (no tags/custom_fields), console_port_template, console_server_port_template

#### Batch 4: Miscellaneous Description/Comments Resources
Resources from Batch 3 audit that have description or comments fields:
- [x] console_port_resource.go ✅ (description only)
- [x] console_server_port_resource.go ✅ (description only)
- [x] device_bay_resource.go ✅ (description only)
- [x] module_bay_resource.go ✅ (description only)
- [x] power_port_resource.go ✅ (description only)
- [x] journal_entry_resource.go ✅ (comments only - no CommentsOnlyAttributes helper, used inline)

**Completed: 6/6 ✅**
**Lines saved: 13 lines**
**Commit**: e8bc4c9

#### Batch 5: Remaining Category 1 Resources (Description + Comments + Tags + Custom Fields)
Resources with full descriptive metadata - organized into sub-batches:

**Batch 5a - Circuits & Services (3 resources):** ✅
- [x] circuit_resource.go ✅
- [x] service_resource.go ✅
- [x] service_template_resource.go ✅

**Completed: 3/3 ✅**
**Lines saved: 26 lines**
**Commit**: cc4aaea

**Note**: circuit_group_resource.go moved to Batch 6 (lacks comments field - Category 2 pattern)

**Batch 5b - Interfaces & Inventory (1 resource):** ✅
- [x] route_target_resource.go ✅

**Completed: 1/1 ✅**
**Lines saved: 2 lines**
**Commit**: ad298d5

**Note**: interface, inventory_item, role lack comments field - moved to Batch 6 (Category 2 pattern)

**Batch 5c - VPN & Wireless (4 resources):** ✅
- [x] ike_policy_resource.go ✅
- [x] ike_proposal_resource.go ✅
- [x] tunnel_resource.go ✅
- [x] wireless_lan_resource.go ✅

**Completed: 4/4 ✅**
**Lines saved: 7 lines**
**Commit**: 0829646

**Batch 5d - Virtual Resources (2 resources):**
- [ ] virtual_chassis_resource.go
- [ ] virtual_device_context_resource.go

**Total Batch 5: 10 resources (8 complete, 2 remaining)**
**Pattern**: Use `CommonDescriptiveAttributes()` + `CommonMetadataAttributes()`
**Estimated savings for remaining**: ~4-6 lines

#### Batch 6: Category 2 Resources (Description Only, No Comments)
Resources with description + tags + custom_fields but no comments field:

**Batch 6a (9 resources):**
- [ ] circuit_group_resource.go (from Batch 5a)
- [ ] interface_resource.go (from Batch 5b)
- [ ] inventory_item_resource.go (from Batch 5b)
- [ ] role_resource.go (from Batch 5b)
- [ ] front_port_resource.go
- [ ] notification_group_resource.go
- [ ] power_outlet_resource.go
- [ ] rear_port_resource.go
- [ ] vm_interface_resource.go

**Batch 6b (3 resources):**
- [ ] wireless_lan_group_resource.go
- [ ] wireless_link_resource.go
- [ ] circuit_group_assignment_resource.go

**Total Batch 6: 12 resources**
**Pattern**: Use `DescriptionOnlyAttributes()` + `CommonMetadataAttributes()`
**Estimated savings**: ~12-24 lines

#### Batch 7: Special Cases (Non-Standard Metadata)
Resources with unique schemas or missing standard metadata support:
- [ ] config_context_resource.go (description + tags, no custom_fields)
- [ ] config_template_resource.go (description only, no metadata)
- [ ] custom_field_resource.go (description + comments, no tags/custom_fields)
- [ ] custom_field_choice_set_resource.go (description only, no metadata)
- [ ] custom_link_resource.go (no standard attributes)
- [ ] event_rule_resource.go (needs verification)
- [ ] export_template_resource.go (description only, no metadata)
- [ ] tag_resource.go (description only, can't tag a tag)
- [ ] webhook_resource.go (description + tags, no custom_fields)

**Total Batch 7: 9 resources**
**Pattern**: Individual handling required
**Estimated savings**: ~5-15 lines

**Total Remaining Phase 3 Savings**: ~41-73 lines across 31 resources

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
