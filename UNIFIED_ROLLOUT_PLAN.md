# Unified Import & State Consistency Fix - Complete Rollout Plan

## Executive Summary

**Problems Identified:**
1. Import doesn't preserve optional fields because helper functions check `IsNull()` and skip population when state is empty
2. Custom fields and tags handling is duplicated across ~68 resources with inconsistent patterns
3. Tests don't validate state consistency after Create/Update/Import operations
4. Custom field import tests can't run in parallel due to global NetBox custom fields

**Solution:**
1. Consolidated helper functions for custom fields (`PopulateCustomFieldsFromAPI`) and tags (`PopulateTagsFromAPI`)
2. Migrate all resources to use these helpers consistently
3. Add `PlanOnly` validation steps to **Create**, **Update**, AND **Import** tests for ALL resources
4. Remove `t.Parallel()` from tests that create/delete global custom fields

**Scope:**
- **101 total resources** in the provider
- **68 resources** with custom fields and/or tags support
- **~250+ tests** needing PlanOnly step enhancement (create + update + import)

---

## Phase 1: Helper Function Consolidation ✅ COMPLETE

### New Helper Functions (in `internal/utils/state_helpers.go`)

#### 1. `PopulateCustomFieldsFromAPI`
Comprehensive custom fields handling that replaces ~60-80 lines of inline code per resource:
- **Normal Create/Read/Update**: Uses state type information to preserve field types
- **Import**: Infers types from API values when state is null/unknown
- **Empty preservation**: Maintains explicit empty sets (`custom_fields = []`) vs null

```go
data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, obj.HasCustomFields(), obj.GetCustomFields(), data.CustomFields, diags)
```

#### 2. `PopulateTagsFromAPI`
Comprehensive tags handling:
- **Normal Create/Read/Update**: Converts API tags to TagModels
- **Empty preservation**: Maintains explicit empty sets (`tags = []`) vs null

```go
data.Tags = utils.PopulateTagsFromAPI(ctx, obj.HasTags(), obj.GetTags(), data.Tags, diags)
```

### Deprecated Functions (kept for backwards compatibility)
- `PopulateCustomFieldsFromMap` → Use `PopulateCustomFieldsFromAPI`
- `PopulateTagsFromNestedTags` → Use `PopulateTagsFromAPI`

---

## Phase 2: Resource Migration Pattern

### For EACH Resource with Custom Fields/Tags (68 resources)

#### A. Update `mapResponseToModel` Function
Replace inline custom fields/tags handling with consolidated helpers:

**Before (~80+ lines):**
```go
// Handle tags
wasTagsExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
if obj.HasTags() && len(obj.GetTags()) > 0 {
    // ... 20+ lines
}
// Handle custom fields
wasExplicitlyEmpty := !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() && len(data.CustomFields.Elements()) == 0
if obj.HasCustomFields() && len(obj.GetCustomFields()) > 0 {
    // ... 60+ lines with import handling
}
```

**After (4-6 lines):**
```go
// Handle tags using consolidated helper
data.Tags = utils.PopulateTagsFromAPI(ctx, obj.HasTags(), obj.GetTags(), data.Tags, diags)
if diags.HasError() {
    return
}

// Handle custom fields using consolidated helper
data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, obj.HasCustomFields(), obj.GetCustomFields(), data.CustomFields, diags)
```

#### B. Update Acceptance Tests

For EVERY test function in `*_resource_test.go`:

1. **Basic Test (`TestAccXxxResource_basic`)**: Add PlanOnly step after create
2. **Update Test (`TestAccXxxResource_update`)**: Add PlanOnly step after each update
3. **Import Test (`TestAccXxxResource_import`)**: Add PlanOnly step after import
4. **Import with CF Test (`TestAccXxxResource_importWithCustomFieldsAndTags`)**:
   - Add PlanOnly step after import
   - Remove `t.Parallel()` (custom fields are global in NetBox)

**PlanOnly Step Pattern:**
```go
{
    // Verify no changes after create/update/import
    Config:   testAccXxxResourceConfig(...),
    PlanOnly: true,
},
```

---

## Phase 3: Complete Resource List (101 Resources)

### Resources WITH Custom Fields/Tags Support (68 resources)
These need full migration (helper functions + test updates):

| # | Resource | Has CF | Has Tags | Current Pattern | Migration Status |
|---|----------|--------|----------|-----------------|------------------|
| 1 | aggregate | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 2 | asn | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 3 | asn_range | ✅ | ✅ | Inline | Pending |
| 4 | cable | ✅ | ✅ | Inline | Pending |
| 5 | circuit | ✅ | ✅ | Inline | Pending |
| 6 | circuit_group | ✅ | ✅ | Inline | Pending |
| 7 | circuit_group_assignment | ✅ | ✅ | Inline | Pending |
| 8 | circuit_termination | ✅ | ✅ | Inline | Pending |
| 9 | console_port | ✅ | ✅ | Inline | Pending |
| 10 | console_server_port | ✅ | ✅ | Inline | Pending |
| 11 | contact | ✅ | ✅ | Inline | Pending |
| 12 | contact_assignment | ✅ | ✅ | Inline | Pending |
| 13 | contact_group | ✅ | ✅ | Inline | Pending |
| 14 | contact_role | ✅ | ✅ | Inline | Pending |
| 15 | device | ✅ | ✅ | **NEW HELPERS** | ✅ DONE |
| 16 | device_bay | ✅ | ✅ | Inline | Pending |
| 17 | device_role | ✅ | ✅ | Inline | Pending |
| 18 | device_type | ✅ | ✅ | Inline | Pending |
| 19 | event_rule | ✅ | ✅ | Inline | Pending |
| 20 | fhrp_group | ✅ | ✅ | Inline | Pending |
| 21 | front_port | ✅ | ✅ | Inline | Pending |
| 22 | ike_policy | ✅ | ✅ | Inline | Pending |
| 23 | ike_proposal | ✅ | ✅ | Inline | Pending |
| 24 | interface | ✅ | ✅ | Inline | Pending |
| 25 | inventory_item | ✅ | ✅ | Inline | Pending |
| 26 | inventory_item_role | ✅ | ✅ | Inline | Pending |
| 27 | ip_address | ✅ | ✅ | Inline | Pending |
| 28 | ip_range | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 29 | ipsec_policy | ✅ | ✅ | Inline | Pending |
| 30 | ipsec_profile | ✅ | ✅ | Inline | Pending |
| 31 | ipsec_proposal | ✅ | ✅ | Inline | Pending |
| 32 | journal_entry | ✅ | ✅ | Inline | Pending |
| 33 | l2vpn | ✅ | ✅ | Inline | Pending |
| 34 | l2vpn_termination | ✅ | ✅ | Inline | Pending |
| 35 | location | ✅ | ✅ | Inline | Pending |
| 36 | module | ✅ | ✅ | Inline | Pending |
| 37 | module_bay | ✅ | ✅ | Inline | Pending |
| 38 | module_type | ✅ | ✅ | Inline | Pending |
| 39 | power_feed | ✅ | ✅ | Inline | Pending |
| 40 | power_outlet | ✅ | ✅ | Inline | Pending |
| 41 | power_panel | ✅ | ✅ | Inline | Pending |
| 42 | power_port | ✅ | ✅ | Inline | Pending |
| 43 | prefix | ✅ | ✅ | Inline | Pending |
| 44 | provider | ✅ | ✅ | Inline | Pending |
| 45 | provider_account | ✅ | ✅ | Inline | Pending |
| 46 | provider_network | ✅ | ✅ | Inline | Pending |
| 47 | rack | ✅ | ✅ | Inline | Pending |
| 48 | rack_reservation | ✅ | ✅ | Inline | Pending |
| 49 | rack_role | ✅ | ✅ | Inline | Pending |
| 50 | rack_type | ✅ | ✅ | Inline | Pending |
| 51 | rear_port | ✅ | ✅ | Inline | Pending |
| 52 | route_target | ✅ | ✅ | Inline | Pending |
| 53 | service | ✅ | ✅ | Inline | Pending |
| 54 | service_template | ✅ | ✅ | Inline | Pending |
| 55 | site | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 56 | site_group | ✅ | ✅ | Inline | Pending |
| 57 | tenant | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 58 | tenant_group | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 59 | tunnel | ✅ | ✅ | Inline | Pending |
| 60 | tunnel_group | ✅ | ✅ | Inline | Pending |
| 61 | tunnel_termination | ✅ | ✅ | Inline | Pending |
| 62 | virtual_chassis | ✅ | ✅ | Inline | Pending |
| 63 | virtual_device_context | ✅ | ✅ | Inline | Pending |
| 64 | virtual_disk | ✅ | ✅ | Inline | Pending |
| 65 | virtual_machine | ✅ | ✅ | Inline | Pending |
| 66 | vlan | ✅ | ✅ | Inline | Pending |
| 67 | vlan_group | ✅ | ✅ | Inline | Pending |
| 68 | vm_interface | ✅ | ✅ | Inline | Pending |
| 69 | vrf | ✅ | ✅ | PopulateCustomFieldsFromMap | Pending |
| 70 | webhook | ✅ | ✅ | Inline | Pending |
| 71 | wireless_lan | ✅ | ✅ | Inline | Pending |
| 72 | wireless_lan_group | ✅ | ✅ | Inline | Pending |
| 73 | wireless_link | ✅ | ✅ | Inline | Pending |

### Resources WITHOUT Custom Fields/Tags (33 resources)
These only need PlanOnly test updates (no helper migration):

| # | Resource | Notes |
|---|----------|-------|
| 1 | circuit_type | Tags only via different pattern |
| 2 | cluster | Uses different pattern |
| 3 | cluster_group | Uses PopulateCustomFieldsFromMap |
| 4 | cluster_type | Uses PopulateCustomFieldsFromMap |
| 5 | config_context | No CF/Tags |
| 6 | config_template | No CF/Tags |
| 7 | console_port_template | Template - no CF/Tags |
| 8 | console_server_port_template | Template - no CF/Tags |
| 9 | custom_field | Meta resource |
| 10 | custom_field_choice_set | Meta resource |
| 11 | custom_link | No CF/Tags |
| 12 | device_bay_template | Template - no CF/Tags |
| 13 | export_template | No CF/Tags |
| 14 | fhrp_group_assignment | Assignment only |
| 15 | front_port_template | Template - no CF/Tags |
| 16 | interface_template | Template - no CF/Tags |
| 17 | inventory_item_template | Template - no CF/Tags |
| 18 | manufacturer | No CF/Tags |
| 19 | module_bay_template | Template - no CF/Tags |
| 20 | notification_group | No CF/Tags |
| 21 | platform | No CF/Tags |
| 22 | power_outlet_template | Template - no CF/Tags |
| 23 | power_port_template | Template - no CF/Tags |
| 24 | rear_port_template | Template - no CF/Tags |
| 25 | region | Uses PopulateCustomFieldsFromMap |
| 26 | rir | Uses PopulateCustomFieldsFromMap |
| 27 | role | Uses PopulateCustomFieldsFromMap |
| 28 | tag | Meta resource |

---

## Phase 4: Batch Organization

Resources are organized into batches by functional area for efficient context switching.

### Batch 1: Reference Implementation ✅ COMPLETE
**Resource:** device
**Status:** Complete! Ready for replication

- [x] Migrate to `PopulateCustomFieldsFromAPI`
- [x] Migrate to `PopulateTagsFromAPI`
- [x] Add PlanOnly to basic test (already had it)
- [x] Add PlanOnly to update test
- [x] Add PlanOnly to import tests
- [x] Remove t.Parallel() from CF import test (already removed)

### Batch 2: IPAM Core (6 resources)
**Resources:** aggregate, asn, asn_range, ip_address, ip_range, prefix

For each resource:
- [ ] Migrate to new helper functions
- [ ] Add PlanOnly to all test steps
- [ ] Remove t.Parallel() from CF tests

```bash
# Test command
$env:TF_ACC="1"; go test -v -run 'TestAcc(Aggregate|ASN|ASNRange|IPAddress|IPRange|Prefix)Resource' ./internal/resources_acceptance_tests/ -timeout 60m
```

### Batch 3: IPAM Supporting (4 resources)
**Resources:** rir, role, route_target, vrf

### Batch 4: Tenancy (4 resources)
**Resources:** tenant, tenant_group, contact, contact_assignment

### Batch 5: Contacts (2 resources)
**Resources:** contact_group, contact_role

### Batch 6: Sites & Locations (4 resources) ✅
**Resources:** site, site_group, location, region

### Batch 7: Racks (5 resources) ✅
**Resources:** rack, rack_role, rack_type, power_panel, rack_reservation
**Test Results:** 32/32 tests passed
**Commits:** 287bade, c523a1a
**Note:** Fixed rack_reservation import issue - now uses ID format during import for consistency

### Batch 8: Device Types & Roles (3 resources) ✅
**Resources:** device_type, device_role, manufacturer
**Test Results:** 18/18 tests passed
**Commit:** a999d1b
**Note:** Added full Tags/CustomFields support to manufacturer (previously missing)

### Batch 9: Device Components - Ports (6 resources) ✅
**Resources:** console_port, console_server_port, front_port, rear_port, power_port, power_outlet
**Test Results:** 36/36 tests passed
**Commit:** 6cc36f1

### Batch 10: Device Components - Other (4 resources) ✅
**Resources:** device_bay, interface, inventory_item, inventory_item_role
**Line Changes:** -22, -24, -22, -22 (net -90 lines)
**Test Results:** 26/26 tests passed
**Commit:** 8b2b42e

### Batch 11: Modules (3 resources) ✅
**Resources:** module, module_bay, module_type
**Line Changes:** -33, -33, -33 (net -81 lines)
**Test Results:** 17/17 tests passed
**Commit:** f9a82c9

### Batch 12: Power (1 resource) ✅
**Resources:** power_feed
**Line Changes:** -26 lines
**Test Results:** 7/7 tests passed
**Commit:** 49dfff0
**Note:** power_panel was already migrated in Batch 7

### Batch 13: Circuits (5 resources) ✅
**Resources:** circuit, circuit_group, circuit_group_assignment, circuit_termination, circuit_type
**Line Changes:** -25, -17, -4, -23, +2 (net -67 lines)
**Test Results:** 35/35 tests passed
**Commit:** aecf4c6

### Batch 14: Providers (3 resources) ✅
**Resources:** provider, provider_account, provider_network
**Line Changes:** -42, -22, -42 (net -106 lines)
**Test Results:** 18/18 tests passed
**Commit:** b46851b

### Batch 15: Virtualization Core (4 resources) ✅
**Resources:** cluster, cluster_group, cluster_type, virtual_machine
**Line Changes:** -25, -18, -18, -23 (net -84 lines)
**Test Results:** 24/24 tests passed
**Commit:** c50b3ac
**Note:** Fixed virtual_machine import handling for site inheritance and cluster reference consistency. Uses standard ImportStateVerifyIgnore pattern for custom_fields.

### Batch 16: Virtualization Components (4 resources) ✅
**Resources:** virtual_chassis, virtual_device_context, virtual_disk, vm_interface
**Line Changes:** -35, -27, -31, -35 (net -128 lines)
**Test Results:** 28/28 tests passed
**Commit:** 08f2322
**Note:** All resources successfully migrated to unified helpers. Special handling for virtual_disk which uses pointer checks instead of Has* methods.

### Batch 17: VLANs (2 resources) ✅
**Resources:** vlan, vlan_group
**Line Changes:** -45, -35 (net -80 lines)
**Test Results:** 15/15 tests passed
**Commit:** a354573
**Note:** Fixed vlan import test to pass random names as parameters instead of regenerating them inside config function.

### Batch 18: VPN - L2 (2 resources) ✅
**Resources:** l2vpn, l2vpn_termination
**Line Changes:** -17, -17 (net -34 lines)
**Test Results:** 10/10 tests passed
**Commit:** b3f1ba4

### Batch 19: VPN - Tunnels (3 resources) ✅
**Resources:** tunnel, tunnel_group, tunnel_termination
**Line Changes:** -42, -36, -41 (net -119 lines)
**Test Results:** 19/19 tests passed
**Commit:** ed755ba

### Batch 20: VPN - IPSec (3 resources) ✅
**Resources:** ipsec_policy, ipsec_profile, ipsec_proposal
**Line Changes:** -19, -19, -19 (net -57 lines)
**Test Results:** 15/15 tests passed
**Commit:** 86f30b4

### Batch 21: VPN - IKE (2 resources) ✅
**Resources:** ike_policy, ike_proposal
**Line Changes:** -18, -18 (net -36 lines)
**Test Results:** 12/12 tests passed
**Commit:** 1ab0fdb

### Batch 22: Wireless (3 resources) ✅
**Resources:** wireless_lan, wireless_lan_group, wireless_link
**Line Changes:** -35, -35, -27 (net -97 lines)
**Test Results:** 14/14 tests passed
**Commit:** 7d5d7ce

### Batch 23: Services (2 resources) ✅
**Resources:** service, service_template
**Line Changes:** -36, -32 (net -68 lines)
**Test Results:** 10/10 tests passed
**Commit:** Pending

### Batch 24: Extras (4 resources)
**Resources:** custom_field, custom_field_choice_set, custom_link, tag

### Batch 25: Events & Automation (3 resources)
**Resources:** event_rule, webhook, journal_entry

### Batch 26: Config & Templates (2 resources)
**Resources:** config_context, config_template

### Batch 27: Export & Notification (2 resources)
**Resources:** export_template, notification_group

### Batch 28: Miscellaneous (3 resources)
**Resources:** platform, fhrp_group, fhrp_group_assignment

### Batch 29: Templates (10 resources)
**Resources:** All `*_template` resources (console_port_template, etc.)

---

## Phase 5: Per-Resource Checklist Template

For each resource, complete this checklist:

### Resource: `xxx_resource.go`

#### Code Changes
- [ ] Update imports (add `netbox` if using tags)
- [ ] Replace inline tags handling with `PopulateTagsFromAPI`
- [ ] Replace inline custom fields handling with `PopulateCustomFieldsFromAPI`
- [ ] Remove any `IsNull()` checks that skip API population
- [ ] Build successfully: `go build .`

#### Test Changes (`xxx_resource_test.go`)
- [ ] **Basic test**: Add PlanOnly step after create
- [ ] **Update test**: Add PlanOnly step after each update
- [ ] **Import test**: Add PlanOnly step after import
- [ ] **CF Import test**: Add PlanOnly step + remove `t.Parallel()`
- [ ] All tests pass

#### Verification
- [ ] Run resource tests: `$env:TF_ACC="1"; go test -v -run 'TestAccXxxResource' ./internal/resources_acceptance_tests/ -timeout 30m`
- [ ] No regressions in CRUD operations
- [ ] Import produces clean plan

---

## Phase 6: Testing Strategy

### Test Enhancement Pattern

**For Basic/Create Tests:**
```go
func TestAccXxxResource_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        // ...
        Steps: []resource.TestStep{
            {
                Config: testAccXxxResourceConfig_basic(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    // existing checks
                ),
            },
            {
                // NEW: Verify no drift after create
                Config:   testAccXxxResourceConfig_basic(),
                PlanOnly: true,
            },
        },
    })
}
```

**For Update Tests:**
```go
func TestAccXxxResource_update(t *testing.T) {
    resource.Test(t, resource.TestCase{
        // ...
        Steps: []resource.TestStep{
            {
                Config: testAccXxxResourceConfig_initial(),
                Check: /* ... */,
            },
            {
                // NEW: Verify no drift after create
                Config:   testAccXxxResourceConfig_initial(),
                PlanOnly: true,
            },
            {
                Config: testAccXxxResourceConfig_updated(),
                Check: /* ... */,
            },
            {
                // NEW: Verify no drift after update
                Config:   testAccXxxResourceConfig_updated(),
                PlanOnly: true,
            },
        },
    })
}
```

**For Import Tests:**
```go
func TestAccXxxResource_importWithCustomFieldsAndTags(t *testing.T) {
    // NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
    // that would affect other tests of the same resource type running in parallel.

    resource.Test(t, resource.TestCase{
        // ...
        Steps: []resource.TestStep{
            {
                Config: testAccXxxResourceConfig_withCF(params...),
                Check: /* ... */,
            },
            {
                ResourceName:      "netbox_xxx.test",
                ImportState:       true,
                ImportStateVerify: true,
                ImportStateVerifyIgnore: []string{/* reference fields */},
            },
            {
                // NEW: Verify import produces clean plan
                Config:   testAccXxxResourceConfig_withCF(params...),
                PlanOnly: true,
            },
        },
    })
}
```

### Custom Fields Test Parallelism Rule

**CRITICAL:** Tests that create/delete custom fields for a resource type CANNOT run in parallel.

NetBox custom fields are **GLOBAL** - when a test creates a custom field for `dcim.device`, it appears on ALL device objects. If Test A deletes its custom fields while Test B's device still references them, Test B fails with "Unknown field name" errors.

**Solution:** Remove `t.Parallel()` from all `*_importWithCustomFieldsAndTags` tests and add explanatory comment.

---

## Phase 7: Verification & Cleanup

### Final Verification Steps

1. **Build check:**
   ```bash
   go build .
   ```

2. **Unit tests:**
   ```bash
   go test ./internal/utils/... -v
   ```

3. **Full acceptance test suite:**
   ```bash
   $env:TF_ACC="1"; go test ./internal/resources_acceptance_tests/... -v -timeout 180m
   ```

4. **Spot-check manual imports:**
   - aggregate (reference fields)
   - device (complex resource)
   - circuit_termination (many optional fields)
   - virtual_machine (VM context)

### Cleanup Tasks

- [ ] Delete `IMPORT_FIX_ROLLOUT_PLAN.md` (superseded)
- [ ] Delete `IMPORT_FIX_ROLLOUT_PLAN_V2.md` (superseded)
- [ ] Delete any `.backup` files
- [ ] Update `CONTRIBUTING.md` with correct patterns
- [ ] Update `DEVELOPMENT.md` if needed

---

## Timeline Estimate

| Phase | Work | Duration |
|-------|------|----------|
| Phase 1 | Helper consolidation | ✅ COMPLETE |
| Phase 2-4 | Resource migration (101 resources) | ~25-30 hours |
| Phase 5 | Test enhancements (~250 tests) | Included above |
| Phase 6 | Verification | 2 hours |
| Phase 7 | Cleanup | 30 minutes |
| **TOTAL** | | **~28-33 hours** |

**Per-resource estimate:** 15-20 minutes average
- Simple resources (no CF/Tags): 5-10 minutes
- Complex resources (with CF/Tags): 20-30 minutes

---

## Progress Tracking

### Completed
- [x] Phase 1: Helper function consolidation
- [x] device_resource.go helper migration

### In Progress
- [ ] Batch 1: device (test updates)

### Pending
- [ ] Batches 2-29 (100 resources)
- [ ] Phase 6: Verification
- [ ] Phase 7: Cleanup

---

## Key Differences from Previous Plans

1. **Unified approach**: Single plan replacing IMPORT_FIX_ROLLOUT_PLAN.md and V2
2. **Complete scope**: All 101 resources, not just 40-68
3. **Test coverage**: PlanOnly in Create + Update + Import (not just Import)
4. **Parallel test fix**: Documented t.Parallel() removal for CF tests
5. **Helper consolidation**: Integrated as Phase 1 (already complete)
6. **Realistic timeline**: 28-33 hours based on full scope

---

## Notes

- **Reference fields**: May change format during import (ID→name/slug) - use `ImportStateVerifyIgnore`
- **Optional fields**: Remain Optional (NOT Computed)
- **Import config**: Must match resource's actual state exactly
- **Template resources**: Only need PlanOnly test updates (no CF/Tags)
- **Meta resources**: (custom_field, tag, etc.) have different patterns
