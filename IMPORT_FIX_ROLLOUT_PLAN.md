# Import State Preservation Fix - Rollout Plan

## Overview

This plan outlines the systematic fix for import state preservation issues across all NetBox provider resources. The root cause is helper functions that check if current state is null and skip populating fields from the API, which breaks import since state starts empty.

## Phase 1: Update Utility Functions

### File: `internal/utils/state_helpers.go`

Replace the existing helper functions with the fixed versions from `state_helpers_fixed.go`:

#### Functions to Update:

1. **UpdateReferenceAttribute**
   - Remove: `if current.IsNull() { return current }`
   - Keep: Logic to preserve user's input format when it matches
   - Add: Default to name/slug/ID when current is null or doesn't match

2. **StringFromAPI**
   - Remove: `if !current.IsNull() { return types.StringNull() } return current`
   - Replace with: `return types.StringNull()`

3. **StringFromAPIPreserveEmpty**
   - Remove: `if !current.IsNull() { return types.StringNull() } return current`
   - Replace with: `return types.StringNull()`

4. **NullableStringFromAPI**
   - Remove: `if !current.IsNull() { return types.StringNull() } return current`
   - Replace with: `return types.StringNull()`

5. **Int64FromAPI**
   - Remove: `if !current.IsNull() { return types.Int64Null() } return current`
   - Replace with: `return types.Int64Null()`

6. **Int64FromInt32API**
   - Remove: `if !current.IsNull() { return types.Int64Null() } return current`
   - Replace with: `return types.Int64Null()`

7. **NullableInt64FromAPI**
   - Remove: `if !current.IsNull() { return types.Int64Null() } return current`
   - Replace with: `return types.Int64Null()`

8. **Float64FromAPI**
   - Remove: `if !current.IsNull() { return types.Float64Null() } return current`
   - Replace with: `return types.Float64Null()`

9. **NullableFloat64FromAPI**
   - Remove: `if !current.IsNull() { return types.Float64Null() } return current`
   - Replace with: `return types.Float64Null()`

### Implementation Steps:

```bash
# Step 1: Backup current state_helpers.go
cp internal/utils/state_helpers.go internal/utils/state_helpers.go.backup

# Step 2: Apply fixes to state_helpers.go (merge fixed versions)

# Step 3: Delete temporary fixed file
rm internal/utils/state_helpers_fixed.go

# Step 4: Run all unit tests
go test ./internal/utils/...
```

## Phase 2: Update Resource Mapping Functions

### Two-Part Update Strategy

#### Part A: Fix Mapping Logic (IsNull checks)
Resources with explicit problematic `else if data.Field.IsNull()` patterns

#### Part B: Enhance Import Tests (ALL resources with import tests)
Add PlanOnly step after import to verify no changes occur

### Part A: Resources Found with IsNull Issues (via grep search)

Based on actual code analysis, these resources have the problematic pattern:

#### Batch 1: Already Fixed ✅
- **aggregate** - Fixed with UpdateReferenceAttributeFixed

#### Batch 2: CustomFields Pattern (12 resources)
These have `else if data.CustomFields.IsNull()` blocks:
- circuit_termination (also has other fields - see Batch 3)
- device_role
- device_type
- fhrp_group
- journal_entry
- provider
- rack_role
- site_group
- tunnel_group
- virtual_machine
- vlan
- vlan_group
- vm_interface

#### Batch 3: Circuit Termination Special Fields
- **circuit_termination** - Has 7 different fields with the pattern:
  - Site
  - ProviderNetwork
  - PortSpeed
  - UpstreamSpeed
  - XconnectID
  - PPInfo
  - Description
  - CustomFields (also in Batch 2)

### Part B: ALL Resources with Import Tests (40 resources)

These resources have `TestAccXxxResource_importWithCustomFieldsAndTags` tests that need PlanOnly step added:

#### Import Test Enhancement Batches (by functional area):

**Batch A: IPAM Resources (5)** - Already has PlanOnly ✅
- aggregate ✅
- asn
- asn_range
- ip_range
- tenant

**Batch B: Tenancy/Organizational (5)**
- contact_assignment
- contact_group
- contact_role
- inventory_item_role
- tenant_group

**Batch C: Virtualization (6)**
- cluster
- cluster_group
- cluster_type
- virtual_chassis
- virtual_device_context
- virtual_disk

**Batch D: VM/Interfaces (3)**
- vm_interface
- vrf
- vlan

**Batch E: Device Components Part 1 (6)**
- console_port
- console_server_port
- device_bay
- front_port
- interface
- inventory_item

**Batch F: Device Components Part 2 (5)**
- module
- module_bay
- power_feed
- power_outlet
- power_port

**Batch G: Device Infrastructure (5)**
- device_role
- device_type
- location
- rack
- rear_port

**Batch H: Circuits & Miscellaneous (5)**
- cable
- circuit
- circuit_termination
- l2vpn
- site

### Combined Batch Organization (Simplified - handles both mapping AND tests)

#### Batch 1: Aggregate ✅ COMPLETE (Reference implementation)
**Status**: Fixed and tested
- Mapping: Fixed tenant, date_added, description, comments
- Import Test: Added PlanOnly step
- Tests: All passing

#### Batch 2: IPAM Resources ✅ COMPLETE (4 resources)
**Status**: Fixed and tested
- Mapping: Fixed aggregate_resource.go (tenant, date_added, description, comments)
- Import Tests: Added PlanOnly steps to all import tests
- Resources: asn, asn_range, ip_range, tenant
- Tests: All passing (40-50s each)
  - TestAccASNResource_import: PASS (50.65s)
  - TestAccTenantResource_import: PASS (46.59s)
  - TestAccTenantResource_importWithCustomFieldsAndTags: PASS (47.55s)
  - TestAccASNRangeResource_import: PASS (49.89s)
  - TestAccASNRangeResource_importWithCustomFieldsAndTags: PASS (49.30s)
  - TestAccIPRangeResource_import: PASS (49.29s)
  - TestAccIPRangeResource_importWithCustomFieldsAndTags: PASS (40.86s)

#### Batch 3: CustomFields-Only Resources ✅ COMPLETE (12 resources)
**Status**: Fixed and tested
**Mapping Fix**: Removed `else if data.CustomFields.IsNull()` blocks from 12 resources
**Import Test Enhancement**: Added PlanOnly steps to 13 import tests

Resources:
1. device_role (comprehensive test)
2. device_type (comprehensive test)
3. fhrp_group
4. journal_entry
5. provider
6. rack_role
7. site_group
8. tunnel_group
9. virtual_machine (no import test)
10. vlan (basic + comprehensive tests)
11. vlan_group
12. vm_interface (basic + comprehensive tests)

**Test Results**: All passing (53-65s each)
- TestAccDeviceRoleResource_importWithCustomFieldsAndTags: PASS (53.93s)
- TestAccDeviceTypeResource_importWithCustomFieldsAndTags: PASS (55.28s)
- TestAccFHRPGroupResource_import: PASS (64.90s)
- TestAccJournalEntryResource_import: PASS (65.26s)
- TestAccProviderResource_import: PASS (64.83s)
- TestAccRackRoleResource_import: PASS (62.00s)
- TestAccSiteGroupResource_import: PASS (64.24s)
- TestAccTunnelGroupResource_import: PASS (62.70s)
- TestAccVLANGroupResource_import: PASS (64.97s)
- VLAN (both tests): PASS
- VMInterface (both tests): PASS

#### Batch 4: Circuit Termination ✅ COMPLETE (1 resource)
**Status**: Fixed and tested
**Mapping Fix**: Removed IsNull checks from 7 fields (Site, ProviderNetwork, PortSpeed, UpstreamSpeed, XconnectID, PPInfo, Description)
**Import Test Enhancement**: Added PlanOnly validation to both import tests

Resource: circuit_termination

**Test Results**: All passing (65-69s each)
- TestAccCircuitTerminationResource_import: PASS (69.38s)
- TestAccCircuitTerminationResource_importWithCustomFieldsAndTags: PASS (65.67s)

Note: Comprehensive test was refactored to pass random values as parameters
rather than generating them inside the config function, ensuring consistent
values across create and PlanOnly validation steps.

#### Batch 5: Import Test Enhancements Only - Tenancy/Org (5 resources - 1 hour)
**No Mapping Changes** - Just add PlanOnly steps

Resources (Batch B):
1. contact_assignment
2. contact_group
3. contact_role
4. inventory_item_role
5. tenant_group

**Test Command**:
```bash
TF_ACC=1 go test -v -run 'TestAcc(ContactAssignment|ContactGroup|ContactRole|InventoryItemRole|TenantGroup)Resource_import' ./internal/resources_acceptance_tests/ -timeout 20m
```

#### Batch 5: Import Test Enhancements Only - Virtualization (6 resources - 1 hour)
**No Mapping Changes** - Just add PlanOnly steps

Resources (Batch C):
1. cluster
2. cluster_group
3. cluster_type
4. virtual_chassis
5. virtual_device_context
6. virtual_disk

**Test Command**:
```bash
TF_ACC=1 go test -v -run 'TestAcc(Cluster|ClusterGroup|ClusterType|VirtualChassis|VirtualDeviceContext|VirtualDisk)Resource_import' ./internal/resources_acceptance_tests/ -timeout 25m
```

#### Batch 6: Import Test Enhancements Only - Device Components (11 resources - 1.5 hours)
**No Mapping Changes** - Just add PlanOnly steps

Resources (Batches E + F):
1. console_port
2. console_server_port
3. device_bay
4. front_port
5. interface
6. inventory_item
7. module
8. module_bay
9. power_feed
10. power_outlet
11. power_port

**Test Command**:
```bash
TF_ACC=1 go test -v -run 'TestAcc(ConsolePort|ConsoleServerPort|DeviceBay|FrontPort|Interface|InventoryItem|Module|ModuleBay|PowerFeed|PowerOutlet|PowerPort)Resource_import' ./internal/resources_acceptance_tests/ -timeout 35m
```

#### Batch 7: Import Test Enhancements Only - Infrastructure (7 resources - 1 hour)
**No Mapping Changes** - Just add PlanOnly steps

Resources (Batch G + H + remaining):
1. location
2. rack
3. rear_port
4. cable
5. circuit
6. l2vpn
7. site

**Test Command**:
```bash
TF_ACC=1 go test -v -run 'TestAcc(Location|Rack|RearPort|Cable|Circuit|L2vpn|Site)Resource_import' ./internal/resources_acceptance_tests/ -timeout 25m
```

#### Batch 8: Import Test Enhancements - Remaining IPAM (4 resources - 45 minutes)
**No Mapping Changes** - Just add PlanOnly steps

Resources (Batch A - excluding aggregate):
1. asn
2. asn_range
3. ip_range
4. tenant

**Test Command**:
```bash
TF_ACC=1 go test -v -run 'TestAcc(ASN|ASNRange|IPRange|Tenant)Resource_import' ./internal/resources_acceptance_tests/ -timeout 20m
```

#### Batch 9: Import Test Enhancements - Remaining VRF (1 resource - 15 minutes)
**No Mapping Changes** - Just add PlanOnly step

Resource (Batch D):
1. vrf

**Test Command**:
```bash
TF_ACC=1 go test -v -run 'TestAccVRFResource_import' ./internal/resources_acceptance_tests/ -timeout 10m
```

## Phase 3: Testing Strategy

### For Each Batch:

1. **Update Resource Files**
   - Apply mapping fixes to all resources in batch
   - Remove `else if data.Field.IsNull()` blocks

2. **Update Acceptance Tests**
   - Add PlanOnly step after import in existing import tests
   - Ensure ImportStateVerifyIgnore includes reference fields that may change format

3. **Run Tests**
   ```bash
   # Run acceptance tests for the batch
   TF_ACC=1 go test -v -run 'TestAcc(Resource1|Resource2|...)' ./internal/resources_acceptance_tests/ -timeout 20m
   ```

4. **Verify Results**
   - All tests must pass
   - Import tests should show no changes in PlanOnly step
   - Reference fields format may change (ID→name/slug) - that's OK if ignored

### Test Template for Each Resource:

Existing import tests should work correctly after the fix. The key is ensuring they:
1. Import with the SAME config used to create (not minimal config)
2. Add PlanOnly step to verify no changes
3. Ignore reference fields that may change format (rir, tenant, etc.)

## Phase 4: Cleanup

After all batches are complete and tested:

1. **Delete temporary files:**
   ```bash
   rm internal/utils/state_helpers_fixed.go
   rm internal/utils/state_helpers.go.backup
   rm internal/examples/import_fix_pattern.go
   ```

2. **Update documentation:**
   - Update CONTRIBUTING.md with correct mapping patterns
   - Document that Optional fields should NOT use IsNull() checks in mapping
   - Keep IMPORT_FIX_IMPLEMENTATION.md as reference

3. **Final verification:**
   ```bash
   # Run ALL acceptance tests
   TF_ACC=1 go test -v ./internal/resources_acceptance_tests/ -timeout 120m

   # Run Terraform integration tests
   .\scripts\run-terraform-tests.ps1
   ```

## Implementation Checklist

### Phase 1: Utilities (30 minutes)
- [ ] Backup state_helpers.go
- [ ] Update UpdateReferenceAttribute (remove IsNull check)
- [ ] Update StringFromAPI (remove IsNull check)
- [ ] Update NullableStringFromAPI (remove IsNull check)
- [ ] Update Int64FromAPI and variants (remove IsNull checks)
- [ ] Update Float64FromAPI and variants (remove IsNull checks)
- [ ] Run utility unit tests: `go test ./internal/utils/...`
- [ ] Delete state_helpers_fixed.go

### Phase 2: Resource Updates (Organized into 9 Batches)

#### Part A: Fix Mapping Logic (IsNull checks) - 14 Resources

##### Batch 1: Aggregate ✅ COMPLETE (Reference implementation)
- [x] Update aggregate_resource.go mapping (removed 4 IsNull checks)
- [x] Add PlanOnly step to import test
- [x] Run tests - ALL PASSING
- [x] Commit: "fix: Aggregate resource import state preservation"

##### Batch 2: IPAM CustomFields Resources (1 hour)
**Resources:** asn, asn_range, ip_range, tenant
- [ ] asn_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] asn_range_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] ip_range_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] tenant_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAccAsnResource_importWithCustomFieldsAndTags|TestAccAsnRangeResource_importWithCustomFieldsAndTags|TestAccIpRangeResource_importWithCustomFieldsAndTags|TestAccTenantResource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "fix: IPAM resources import state preservation"

##### Batch 3: Device Role/Type/Components CustomFields (1.5 hours)
**Resources:** device_role, device_type, fhrp_group, rack_role, vm_interface
- [ ] device_role_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] device_type_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] fhrp_group_resource.go - Remove CustomFields IsNull check
- [ ] rack_role_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] vm_interface_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAccDeviceRoleResource_importWithCustomFieldsAndTags|TestAccDeviceTypeResource_importWithCustomFieldsAndTags|TestAccFhrpGroupResource|TestAccRackRoleResource_importWithCustomFieldsAndTags|TestAccVmInterfaceResource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "fix: Device/rack role and type import state preservation"

##### Batch 4: Network/VM CustomFields (1 hour)
**Resources:** journal_entry, provider, site_group, tunnel_group, virtual_machine, vlan, vlan_group
- [ ] journal_entry_resource.go - Remove CustomFields IsNull check
- [ ] provider_resource.go - Remove CustomFields IsNull check
- [ ] site_group_resource.go - Remove CustomFields IsNull check
- [ ] tunnel_group_resource.go - Remove CustomFields IsNull check
- [ ] virtual_machine_resource.go - Remove CustomFields IsNull check
- [ ] vlan_resource.go - Remove CustomFields IsNull check, add PlanOnly to import test
- [ ] vlan_group_resource.go - Remove CustomFields IsNull check
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAcc(JournalEntry|Provider|SiteGroup|TunnelGroup|VirtualMachine|Vlan|VlanGroup)Resource' -v`
- [ ] Commit: "fix: Network/VM resources import state preservation"

##### Batch 5: Circuit Termination (30 minutes)
**Resources:** circuit_termination (7 fields + CustomFields)
- [ ] circuit_termination_resource.go - Remove 8 IsNull checks (Site, ProviderNetwork, PortSpeed, UpstreamSpeed, XconnectID, PPInfo, Description, CustomFields), add PlanOnly to import test
- [ ] Run tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAccCircuitTerminationResource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "fix: Circuit termination import state preservation"

#### Part B: Enhance Import Tests ONLY (PlanOnly steps) - 33 Additional Resources

##### Batch 6: Tenancy & Contacts (1 hour)
**Resources:** contact_assignment, contact_group, contact_role, inventory_item_role, tenant_group
- [ ] Add PlanOnly step to contact_assignment import test
- [ ] Add PlanOnly step to contact_group import test
- [ ] Add PlanOnly step to contact_role import test
- [ ] Add PlanOnly step to inventory_item_role import test
- [ ] Add PlanOnly step to tenant_group import test
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAcc(ContactAssignment|ContactGroup|ContactRole|InventoryItemRole|TenantGroup)Resource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "test: Add PlanOnly validation to tenancy import tests"

##### Batch 7: Virtualization (1.5 hours)
**Resources:** cluster, cluster_group, cluster_type, virtual_chassis, virtual_device_context, virtual_disk, vrf
- [ ] Add PlanOnly step to cluster import test
- [ ] Add PlanOnly step to cluster_group import test
- [ ] Add PlanOnly step to cluster_type import test
- [ ] Add PlanOnly step to virtual_chassis import test
- [ ] Add PlanOnly step to virtual_device_context import test
- [ ] Add PlanOnly step to virtual_disk import test
- [ ] Add PlanOnly step to vrf import test
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAcc(Cluster|ClusterGroup|ClusterType|VirtualChassis|VirtualDeviceContext|VirtualDisk|Vrf)Resource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "test: Add PlanOnly validation to virtualization import tests"

##### Batch 8: Device Components Part 1 (1.5 hours)
**Resources:** console_port, console_server_port, device_bay, front_port, interface, inventory_item, module, module_bay
- [ ] Add PlanOnly step to console_port import test
- [ ] Add PlanOnly step to console_server_port import test
- [ ] Add PlanOnly step to device_bay import test
- [ ] Add PlanOnly step to front_port import test
- [ ] Add PlanOnly step to interface import test
- [ ] Add PlanOnly step to inventory_item import test
- [ ] Add PlanOnly step to module import test
- [ ] Add PlanOnly step to module_bay import test
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAcc(ConsolePort|ConsoleServerPort|DeviceBay|FrontPort|Interface|InventoryItem|Module|ModuleBay)Resource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "test: Add PlanOnly validation to device component import tests (part 1)"

##### Batch 9: Device Components Part 2 + Infrastructure (2 hours)
**Resources:** power_feed, power_outlet, power_port, location, rear_port, cable, circuit, l2vpn, site
- [ ] Add PlanOnly step to power_feed import test
- [ ] Add PlanOnly step to power_outlet import test
- [ ] Add PlanOnly step to power_port import test
- [ ] Add PlanOnly step to location import test
- [ ] Add PlanOnly step to rear_port import test
- [ ] Add PlanOnly step to cable import test
- [ ] Add PlanOnly step to circuit import test
- [ ] Add PlanOnly step to l2vpn import test
- [ ] Add PlanOnly step to site import test
- [ ] Run batch tests: `TF_ACC=1 go test ./internal/resources_acceptance_tests -run 'TestAcc(PowerFeed|PowerOutlet|PowerPort|Location|RearPort|Cable|Circuit|L2vpn|Site)Resource_importWithCustomFieldsAndTags' -v`
- [ ] Commit: "test: Add PlanOnly validation to device component and infrastructure import tests (part 2)"

### Phase 3: Verification (1 hour)
- [ ] Run full acceptance test suite: `TF_ACC=1 go test ./internal/resources_acceptance_tests/... -v`
- [ ] Run Terraform integration tests for key resources
- [ ] Spot-check manual imports of key resources:
  - [ ] aggregate (has tenant reference)
  - [ ] circuit_termination (multiple optional fields)
  - [ ] virtual_machine (uses helpers extensively)
- [ ] Verify no regressions in existing CRUD tests

### Phase 4: Cleanup (30 minutes)
- [ ] Delete internal/utils/state_helpers_fixed.go
- [ ] Delete internal/utils/state_helpers.go.backup
- [ ] Delete internal/examples/import_fix_pattern.go
- [ ] Update CONTRIBUTING.md with correct patterns
- [ ] Final commit: "chore: Clean up import fix temporary files"

## Success Criteria

1. ✅ All helper functions always populate from API response
2. ✅ All resource mapResponseToModel functions simplified (no IsNull checks)
3. ✅ All acceptance tests pass
4. ✅ Import tests verify no changes after import
5. ✅ No temporary/duplicate helper functions remain
6. ✅ Documentation updated with correct patterns

## Risk Mitigation

1. **Batch approach**: Small incremental changes, easy to rollback individual batches
2. **Test coverage**: Every batch tested before moving to next
3. **Backup**: Keep state_helpers.go.backup until all batches complete
4. **Git commits**: One commit per batch for easy reversion
5. **Reference fields**: Use ImportStateVerifyIgnore for fields that can change format

## Timeline Estimate

- **Phase 1 (Utilities)**: 30 minutes
- **Phase 2 (Resources)**:
  - Batch 1: ✅ COMPLETE (aggregate reference)
  - Batch 2: 1 hour (4 IPAM resources - mapping + tests)
  - Batch 3: 1.5 hours (5 device/rack resources - mapping + tests)
  - Batch 4: 1 hour (7 network/VM resources - mapping only)
  - Batch 5: 30 minutes (1 circuit_termination - 8 fields + test)
  - Batch 6: 1 hour (5 tenancy resources - tests only)
  - Batch 7: 1.5 hours (7 virtualization resources - tests only)
  - Batch 8: 1.5 hours (8 device components - tests only)
  - Batch 9: 2 hours (9 components + infrastructure - tests only)
- **Phase 3 (Verification)**: 1 hour
- **Phase 4 (Cleanup)**: 30 minutes

**Total remaining estimate: ~10.5 hours** (can be done in 2-3 sessions)

**Already complete: Batch 1 (aggregate) - 1 hour**

**Summary by Work Type:**
- Mapping fixes (IsNull removal): 14 resources across Batches 2-5 (~4 hours)
- Test enhancements (PlanOnly steps): 40 resources across all batches (~6.5 hours)
- Verification & cleanup: ~1.5 hours

## Notes

- Reference fields (tenant, site, etc.) may change format during import (ID→name/slug)
  - This is EXPECTED and CORRECT behavior
  - Add these to ImportStateVerifyIgnore
- Optional fields remain Optional (NOT Computed)
- Import should be done with config matching the resource's actual state
- The fix ensures state accurately reflects NetBox, Terraform handles config matching
