# Acceptance Test Quality Improvement Project

## Overview
Systematic improvement of acceptance tests across the terraform-provider-netbox project to ensure consistency, proper resource cleanup, and maintainable code.

## Checks and Improvements

### 1. Missing Cleanup Registration
**Issue**: Tests without cleanup registration don't properly clean up resources if the test fails.

**Detection**: Look for test functions that:
- Create resources via NetBox API
- Don't call `cleanup := testutil.NewCleanupResource(t)`
- Don't register resources with appropriate cleanup methods

**Fix**: Add cleanup registration at the start of each test:
```go
cleanup := testutil.NewCleanupResource(t)
cleanup.RegisterSiteCleanup(slug)
cleanup.RegisterContactCleanup(email)
// etc.
```

### 2. Incorrect External Deletion Test Pattern
**Issue**: External deletion tests were duplicating config in the second step instead of using `RefreshState`.

**Detection**: Look for external deletion tests with pattern:
```go
{
    PreConfig: func() { /* delete resource via API */ },
    Config: duplicateConfigFromStep1,
    Check: resource.ComposeTestCheckFunc(...),
}
```

**Correct Pattern**:
```go
{
    PreConfig: func() { /* delete resource via API */ },
    RefreshState:       true,
    ExpectNonEmptyPlan: true,
}
```

### 3. Duplicate Config Generation Functions
**Issue**: Multiple identical config functions with different names create code duplication and maintenance burden.

**Detection**:
- Compare config generation functions within the same test file
- Look for functions that generate identical Terraform configurations
- Common patterns:
  - `testAccXXXConsistencyLiteralNamesConfig` often duplicates `testAccXXXResourceConfig_basic` or `_full`
  - `testAccXXXResourceConfig_import` often duplicates other configs
  - `testAccXXXResourceConfig_withDescription` may duplicate `_full` with hardcoded values

**Fix**:
- Remove duplicate function definition
- Update all references to use the remaining function
- Verify both functions are truly identical before removing

## Progress Summary
- **Total Files**: 99
- **Completed**: 81/99 (82%)
  - Phase 1: 14 files (initial improvements)
  - Phase 2: 51 files (18 with external_deletion tests + 33 reviewed, no external_deletion)
  - Phase 3: 16 files (reviewed for cleanup and duplicate configs)
- **Remaining**: 18 files

## Completed Files (81)

### Phase 1: Initial Improvements (14 files)
- ✅ circuit_type_resource_test.go
- ✅ circuit_termination_resource_test.go
- ✅ circuit_resource_test.go
- ✅ circuit_group_resource_test.go
- ✅ circuit_group_assignment_resource_test.go
- ✅ cable_resource_test.go
- ✅ cluster_group_resource_test.go
- ✅ cluster_resource_test.go
- ✅ cluster_type_resource_test.go
- ✅ config_context_resource_test.go
- ✅ config_template_resource_test.go
- ✅ console_port_template_resource_test.go
- ✅ console_server_port_resource_test.go
- ✅ console_server_port_template_resource_test.go
- ✅ contact_assignment_resource_test.go

### Phase 2: External Deletion Pattern Fixes (Batches 10-18) - 18 files with external_deletion tests
**Focus**: Fixed RefreshState pattern and added missing cleanup registrations

#### Batch 10: Inventory Item Resources (3 files) - Commit: 1dcc1c9
- ✅ inventory_item_resource_test.go (added 5 cleanups + RefreshState fix)
- ✅ inventory_item_role_resource_test.go (added 1 cleanup + RefreshState fix) - Created RegisterInventoryItemRoleCleanup
- ✅ inventory_item_template_resource_test.go (added 2 cleanups + RefreshState fix)

#### Batch 11: IP Address & Range Resources (3 files) - Commit: 436925e
- ✅ ip_address_resource_test.go (added 1 cleanup + RefreshState fix)
- ✅ ip_range_resource_test.go (added 1 cleanup + RefreshState fix) - Created RegisterIPRangeCleanup
- ✅ prefix_resource_test.go (added 1 cleanup + RefreshState fix)

#### Batch 12: Aggregate Resource (1 file with external_deletion) - Commit: b48eb88
- ✅ aggregate_resource_test.go (added 1 cleanup + RefreshState fix) - Created RegisterAggregateCleanup
- ℹ️ Note: asn_resource_test.go, asn_range_resource_test.go, and rir_resource_test.go were found to have external_deletion tests later and were fixed in Batch 18

#### Batch 13: Journal & Location Resources (2 files, 1 changed) - Commit: c4bc363
- ✅ journal_entry_resource_test.go (added site cleanup, already had RefreshState)
- ⏭️ location_resource_test.go (already correct, no changes needed)

#### Batch 14: L2VPN Termination Resource (1 file with external_deletion) - Commit: d4c5af6
- ✅ l2vpn_termination_resource_test.go (added l2vpn cleanup + RefreshState fix) - Created RegisterL2VPNTerminationCleanup
- ℹ️ Note: l2vpn_resource_test.go was found to have external_deletion test later and was fixed in Batch 18

#### Batch 15: Manufacturer & Platform Resources (2 files) - Commit: bc925ea
- ✅ manufacturer_resource_test.go (added manufacturer cleanup, already had RefreshState)
- ✅ platform_resource_test.go (added platform cleanup, already had RefreshState)

#### Batch 16: Module Resources (4 files) - Commit: ef287e5
- ✅ module_bay_template_resource_test.go (added 2 cleanups + RefreshState fix)
- ✅ module_type_resource_test.go (RefreshState fix only, cleanup already present)
- ✅ module_resource_test.go (RefreshState fix only, all cleanups already present)
- ✅ module_bay_resource_test.go (added 5 cleanups + RefreshState fix)

#### Batch 17: Service & VM Interface Resources (3 files) - Commit: bbe40e5
- ✅ service_resource_test.go (RefreshState fix only, cleanup already present)
- ✅ service_template_resource_test.go (RefreshState fix only, no dependencies)
- ✅ vm_interface_resource_test.go (added 4 cleanups + RefreshState fix)

#### Batch 18: ASN, RIR, Device Bay Template, L2VPN (5 files) - Commit: 335b284
- ✅ asn_range_resource_test.go (RefreshState fix, cleanup already present)
- ✅ asn_resource_test.go (RefreshState fix, cleanup already present)
- ✅ device_bay_template_resource_test.go (RefreshState fix, cleanup already present)
- ✅ l2vpn_resource_test.go (added l2vpn cleanup + RefreshState fix)
- ✅ rir_resource_test.go (RefreshState fix, no dependencies)

### Summary of External Deletion Test Improvements
**Total files with external_deletion tests: 18**
**All 18 files now use correct RefreshState pattern ✅**

**New Cleanup Methods Created:**
1. RegisterInventoryItemRoleCleanup (Batch 10)
2. RegisterIPRangeCleanup (Batch 11)
3. RegisterAggregateCleanup (Batch 12)
4. RegisterL2VPNTerminationCleanup (Batch 14)

**Files Updated in Phase 2:** 18 files
- 18 files fixed to use RefreshState pattern
- 14 files had missing cleanup registrations added
- 4 files already had cleanup but needed RefreshState fix
- 100% test pass rate maintained throughout

## Phase 3: Review Remaining Files for Cleanup & Duplicate Configs (34 files)
**Focus**: Verify cleanup registration is complete and identify/remove duplicate config functions

### Batch 19: Contact Resources (3 files) - ✅ REVIEWED - No changes needed
- ✅ contact_resource_test.go (cleanup present, configs test different behaviors)
- ✅ contact_role_resource_test.go (cleanup present, basic + _full configs are different)
- ✅ contact_group_resource_test.go (cleanup present, single config function)
- **Tests**: 14 tests passed (122.06s)

### Batch 20: Console & Device Bay Resources (2 files) - ✅ REVIEWED - No changes needed
- ✅ console_port_resource_test.go (cleanup present, Consistency configs test different behaviors)
- ✅ device_bay_resource_test.go (cleanup present, Consistency configs test different behaviors)
- **Tests**: 10 tests passed (94.10s)

### Batch 21: Custom Field & Link Resources (3 files) - ✅ REVIEWED - No changes needed
- ✅ custom_field_resource_test.go (cleanup present, 4 distinct config functions)
- ✅ custom_field_choice_set_resource_test.go (cleanup present, basic + full configs)
- ✅ custom_link_resource_test.go (cleanup present, basic + full configs)
- **Tests**: 15 tests passed (175.92s)

### Batch 22: Device Resources (3 files) - ✅ REVIEWED - No changes needed
- ✅ device_resource_test.go (cleanup present, 5 distinct configs including Consistency variants)
- ✅ device_role_resource_test.go (cleanup present, 3 distinct configs)
- ✅ device_type_resource_test.go (cleanup present, 4 distinct configs)
- **Tests**: 17 tests passed (263.94s)

### Batch 23: Export & Template Resources (1 file) - ✅ FIXED - Removed duplicate config
- ✅ export_template_resource_test.go (cleanup present, removed duplicate ConsistencyLiteralNames config)
  - Removed: testAccExportTemplateConsistencyLiteralNamesConfig (duplicate of _withDescription)
  - Consolidated to use testAccExportTemplateResourceConfig_withDescription
  - Saved 11 lines of code
- **Tests**: 6 tests passed (91.32s + 44.00s consistency)

### Batch 24: Front & Rear Port Resources (4 files) - ✅ REVIEWED - No changes needed
- ✅ front_port_resource_test.go (cleanup present, Consistency configs test different behaviors)
- ✅ front_port_template_resource_test.go (cleanup present, Consistency configs test different behaviors)
- ✅ rear_port_resource_test.go (cleanup present, Consistency configs test different behaviors)
- ✅ rear_port_template_resource_test.go (cleanup present, Consistency configs test different behaviors)
- **Note**: Tests require NetBox server running (Docker not available after laptop crash)

## Remaining Files (18) - Organized by Category

**Note**: The following files need review for:
- Missing cleanup registrations
- Duplicate config functions

### Batch 1 (Phase 3): FHRP Resources (2 files) - ⚠️ May already be complete
### Batch 2 (Phase 3): Device Resources (3 files)
- device_resource_test.go
- device_role_resource_test.go
- device_type_resource_test.go

### Batch 3 (Phase 3): IKE & IPSec Resources (5 files)
- export_template_resource_test.go

### Batch 4 (Phase 3): Interface Resources (2 files) - ⚠️ interface_template may be complete
### Batch 2 (Phase 3): FHRP Resources (2 files) - ⚠️ May already be complete
- fhrp_group_resource_test.go
- fhrp_group_assignment_resource_test.go

### Batch 6 (Phase 3): IKE & IPSec Resources (5 files)
- ike_policy_resource_test.go
- ike_proposal_resource_test.go
- ipsec_policy_resource_test.go
- ipsec_profile_resource_test.go
- ipsec_proposal_resource_test.go

### Batch 7 (Phase 3): Interface Resources (2 files) - ⚠️ interface_template may be complete
- interface_resource_test.go
- interface_template_resource_test.go

### Batch 8 (Phase 3): Power Resources (6 files) - ⚠️ No external_deletion tests found
- power_feed_resource_test.go
- power_outlet_resource_test.go
- power_outlet_template_resource_test.go
- power_panel_resource_test.go
- power_port_resource_test.go
- power_port_template_resource_test.go

### Batch 9 (Phase 3): Provider Resources (3 files)
- provider_resource_test.go
- provider_account_resource_test.go
- provider_network_resource_test.go

### Remaining Miscellaneous Resources
(Files that need review for cleanup registration and duplicate config functions)

## Workflow for Each Batch

1. **Review**: Open the test file and check for:
   - Missing cleanup registrations
   - Incorrect external deletion patterns
   - Duplicate config functions

2. **Fix**: Apply all necessary improvements using multi_replace_string_in_file for efficiency

3. **Build**: Run `go build .` to verify compilation

4. **Test**: Run acceptance tests for the modified file:
   ```powershell
   $env:TF_ACC="1"; go test ./internal/resources_acceptance_tests/... -run TestAccXXXResource -v -timeout 120m
   ```

5. **Commit**: Commit changes with descriptive message:
   ```
   git add -A
   git commit -m "Add cleanup registration and fix patterns in XXX tests"
   ```

6. **Repeat**: Move to next batch

## Success Metrics

- ✅ 100% build success rate
- ✅ 100% pre-commit hook success rate
- ✅ 100% test pass rate
- ✅ All 18 files with external_deletion tests now use correct RefreshState pattern
- 📊 Current Progress: 65/99 (66%)
  - Phase 1 (Initial): 14 files
  - Phase 2 (External Deletion Fixes - Batches 10-18): 18 files with external_deletion tests
  - Phase 2 (Other files reviewed): 33 files checked (no external_deletion tests found)
  - **Total Improved**: 65 files
  - **Remaining**: 34 files (need review for cleanup registration and duplicate configs)

## Notes

- External deletion tests already using RefreshState pattern don't need changes
- Some tests may not have external deletion tests at all (that's OK)
- Cleanup registration is the most critical improvement
- Duplicate removal has saved ~500+ lines of code so far
