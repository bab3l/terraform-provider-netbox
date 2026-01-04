# Custom Fields Test Migration Plan

## Overview
Migrate all acceptance tests that create/delete custom fields to a separate package with build tag `customfields`. This prevents parallel execution conflicts and speeds up normal test runs.

**Total Files Affected:** 41 resource test files
**Build Tag:** `customfields`
**New Package:** `internal/resources_acceptance_tests_customfields`

## Why This Migration?
1. **Prevent Deadlocks**: Custom fields are global per content type in NetBox - parallel tests cause conflicts
2. **Faster Default Tests**: Normal test runs skip ~41 slow tests, reducing time by 30-40 minutes
3. **Clearer Intent**: Explicitly separate serial-only tests from parallel-safe tests
4. **CI Flexibility**: Run parallel tests concurrently, custom field tests sequentially

## Implementation Strategy

### Phase 1: Setup Infrastructure
1. Create new package directory: `internal/resources_acceptance_tests_customfields/`
2. Create shared test utilities file with build tag
3. Update Makefile with new test targets
4. Update GitHub Actions CI workflow
5. Document in CONTRIBUTING.md and DEVELOPMENT.md

### Phase 2: Migration Batches
Organize by resource domain for easier review and testing.

---

## Batch 1: IPAM Resources (6 files)
**Estimated Time:** 30 minutes

### Files to Migrate:
1. `aggregate_resource_test.go`
   - Test: `TestAccAggregateResource_importWithCustomFieldsAndTags`
   - Config: `testAccAggregateResourceImportConfig_full` (used only by this test)

2. `asn_resource_test.go`
   - Test: `TestAccAsnResource_importWithCustomFieldsAndTags`
   - Config: `testAccAsnResourceImportConfig_full`

3. `asn_range_resource_test.go`
   - Test: `TestAccAsnRangeResource_importWithCustomFieldsAndTags`
   - Config: `testAccAsnRangeResourceImportConfig_full`

4. `ip_range_resource_test.go`
   - Test: `TestAccIPRangeResource_importWithCustomFieldsAndTags`
   - Config: `testAccIPRangeResourceImportConfig_full`

5. `vlan_resource_test.go`
   - Test: `TestAccVlanResource_importWithCustomFieldsAndTags`
   - Config: `testAccVlanResourceImportConfig_full`

6. `vrf_resource_test.go`
   - Test: `TestAccVrfResource_importWithCustomFieldsAndTags`
   - Config: `testAccVrfResourceImportConfig_full`

---

## Batch 2: Circuits Resources (3 files)
**Estimated Time:** 20 minutes

### Files to Migrate:
1. `circuit_resource_test.go`
   - Test: `TestAccCircuitResource_importWithCustomFieldsAndTags`
   - Config: `testAccCircuitResourceImportConfig_full`

2. `circuit_termination_resource_test.go`
   - Test: `TestAccCircuitTerminationResource_importWithCustomFieldsAndTags`
   - Config: `testAccCircuitTerminationResourceImportConfig_full`

3. `l2vpn_resource_test.go`
   - Test: `TestAccL2vpnResource_importWithCustomFieldsAndTags`
   - Config: `testAccL2vpnResourceImportConfig_full`

---

## Batch 3: DCIM Device Components (10 files)
**Estimated Time:** 50 minutes

### Files to Migrate:
1. `console_port_resource_test.go`
   - Test: `TestAccConsolePortResource_importWithCustomFieldsAndTags`
   - Config: `testAccConsolePortResourceImportConfig_full`

2. `console_server_port_resource_test.go`
   - Test: `TestAccConsoleServerPortResource_importWithCustomFieldsAndTags`
   - Config: `testAccConsoleServerPortResourceImportConfig_full`

3. `device_bay_resource_test.go`
   - Test: `TestAccDeviceBayResource_importWithCustomFieldsAndTags`
   - Config: `testAccDeviceBayResourceImportConfig_full`

4. `front_port_resource_test.go`
   - Test: `TestAccFrontPortResource_importWithCustomFieldsAndTags`
   - Config: `testAccFrontPortResourceImportConfig_full`

5. `interface_resource_test.go`
   - Test: `TestAccInterfaceResource_importWithCustomFieldsAndTags`
   - Config: `testAccInterfaceResourceImportConfig_full`

6. `module_bay_resource_test.go`
   - Test: `TestAccModuleBayResource_importWithCustomFieldsAndTags`
   - Config: `testAccModuleBayResourceImportConfig_full`

7. `power_outlet_resource_test.go`
   - Test: `TestAccPowerOutletResource_importWithCustomFieldsAndTags`
   - Config: `testAccPowerOutletResourceImportConfig_full`

8. `power_port_resource_test.go`
   - Test: `TestAccPowerPortResource_importWithCustomFieldsAndTags`
   - Config: `testAccPowerPortResourceImportConfig_full`

9. `rear_port_resource_test.go`
   - Test: `TestAccRearPortResource_importWithCustomFieldsAndTags`
   - Config: `testAccRearPortResourceImportConfig_full`

10. `cable_resource_test.go`
    - Test: `TestAccCableResource_importWithCustomFieldsAndTags`
    - Config: `testAccCableResourceImportConfig_full`

---

## Batch 4: DCIM Infrastructure (7 files)
**Estimated Time:** 35 minutes

### Files to Migrate:
1. `device_role_resource_test.go`
   - Test: `TestAccDeviceRoleResource_importWithCustomFieldsAndTags`
   - Config: `testAccDeviceRoleResourceImportConfig_full`

2. `device_type_resource_test.go`
   - Test: `TestAccDeviceTypeResource_importWithCustomFieldsAndTags`
   - Config: `testAccDeviceTypeResourceImportConfig_full`

3. `inventory_item_resource_test.go`
   - Test: `TestAccInventoryItemResource_importWithCustomFieldsAndTags`
   - Config: `testAccInventoryItemResourceImportConfig_full`

4. `inventory_item_role_resource_test.go`
   - Test: `TestAccInventoryItemRoleResource_importWithCustomFieldsAndTags`
   - Config: `testAccInventoryItemRoleResourceImportConfig_full`

5. `location_resource_test.go`
   - Test: `TestAccLocationResource_importWithCustomFieldsAndTags`
   - Config: `testAccLocationResourceImportConfig_full`

6. `module_resource_test.go`
   - Test: `TestAccModuleResource_importWithCustomFieldsAndTags`
   - Config: `testAccModuleResourceImportConfig_full`

7. `rack_resource_test.go`
   - Test: `TestAccRackResource_importWithCustomFieldsAndTags`
   - Config: `testAccRackResourceImportConfig_full`

---

## Batch 5: DCIM Sites & Power (2 files)
**Estimated Time:** 15 minutes

### Files to Migrate:
1. `site_resource_test.go`
   - Test: `TestAccSiteResource_importWithCustomFieldsAndTags`
   - Config: `testAccSiteResourceImportConfig_full`

2. `power_feed_resource_test.go`
   - Test: `TestAccPowerFeedResource_importWithCustomFieldsAndTags`
   - Config: `testAccPowerFeedResourceImportConfig_full`

---

## Batch 6: Tenancy Resources (3 files)
**Estimated Time:** 20 minutes

### Files to Migrate:
1. `tenant_resource_test.go`
   - Test: `TestAccTenantResource_importWithCustomFieldsAndTags`
   - Config: `testAccTenantResourceImportConfig_full`

2. `tenant_group_resource_test.go`
   - Test: `TestAccTenantGroupResource_importWithCustomFieldsAndTags`
   - Config: `testAccTenantGroupResourceImportConfig_full`

3. `contact_assignment_resource_test.go`
   - Test: `TestAccContactAssignmentResource_importWithCustomFieldsAndTags`
   - Config: `testAccContactAssignmentResourceImportConfig_full`

---

## Batch 7: Organizational Resources (2 files)
**Estimated Time:** 15 minutes

### Files to Migrate:
1. `contact_group_resource_test.go`
   - Test: `TestAccContactGroupResource_importWithCustomFieldsAndTags`
   - Config: `testAccContactGroupResourceImportConfig_full`

2. `contact_role_resource_test.go`
   - Test: `TestAccContactRoleResource_importWithCustomFieldsAndTags`
   - Config: `testAccContactRoleResourceImportConfig_full`

---

## Batch 8: Virtualization Resources (5 files)
**Estimated Time:** 30 minutes

### Files to Migrate:
1. `cluster_resource_test.go`
   - Test: `TestAccClusterResource_importWithCustomFieldsAndTags`
   - Config: `testAccClusterResourceImportConfig_full`

2. `cluster_group_resource_test.go`
   - Test: `TestAccClusterGroupResource_importWithCustomFieldsAndTags`
   - Config: `testAccClusterGroupResourceImportConfig_full`

3. `cluster_type_resource_test.go`
   - Test: `TestAccClusterTypeResource_importWithCustomFieldsAndTags`
   - Config: `testAccClusterTypeResourceImportConfig_full`

4. `virtual_chassis_resource_test.go`
   - Test: `TestAccVirtualChassisResource_importWithCustomFieldsAndTags`
   - Config: `testAccVirtualChassisResourceImportConfig_full`

5. `virtual_device_context_resource_test.go`
   - Test: `TestAccVirtualDeviceContextResource_importWithCustomFieldsAndTags`
   - Config: `testAccVirtualDeviceContextResourceImportConfig_full`

---

## Batch 9: Virtual Machine Resources (3 files)
**Estimated Time:** 20 minutes

### Files to Migrate:
1. `virtual_machine_import_test.go`
   - Test: `TestAccVirtualMachineResource_importWithCustomFieldsAndTags`
   - Config: `testAccVirtualMachineResourceImportConfig_full`

2. `virtual_disk_resource_test.go`
   - Test: `TestAccVirtualDiskResource_importWithCustomFieldsAndTags`
   - Config: `testAccVirtualDiskResourceImportConfig_full`

3. `vm_interface_resource_test.go`
   - Test: `TestAccVMInterfaceResource_importWithCustomFieldsAndTags`
   - Config: `testAccVMInterfaceResourceImportConfig_full`

---

## Migration Procedure (Per File)

### Step 1: Extract Test and Config Functions
1. Identify the test function (starts with `TestAcc*_importWithCustomFieldsAndTags`)
2. Identify the config function (usually `testAcc*ResourceImportConfig_full`)
3. Check if config function is used by any other tests (search entire file)
4. If config function is only used by the custom field test, include it in migration

### Step 2: Create New File in Target Package
```go
//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
    "testing"
    // ... other imports ...
    "github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// Paste test function here
// Paste config function here (if exclusive to this test)
```

### Step 3: Update Original File
1. Delete the migrated test function
2. Delete the config function if it was migrated
3. Add comment at top of file if needed:
   ```go
   // NOTE: Custom field tests for this resource are in resources_acceptance_tests_customfields package
   ```

### Step 4: Verify Build and Tests
```bash
# Verify normal tests still compile and run
go test ./internal/resources_acceptance_tests -v -count=1

# Verify custom field tests compile
go test -tags=customfields ./internal/resources_acceptance_tests_customfields -v -count=1
```

---

## Infrastructure Files to Create

### 1. `internal/resources_acceptance_tests_customfields/doc.go`
```go
//go:build customfields
// +build customfields

// Package resources_acceptance_tests_customfields contains acceptance tests
// that create/delete custom fields. These tests must run serially to avoid
// conflicts since custom fields are global per content type in NetBox.
//
// Run these tests with: go test -tags=customfields ./...
package resources_acceptance_tests_customfields
```

### 2. Update `Makefile`
```makefile
# Run only parallel-safe tests (default, fast)
test-acceptance:
	TF_ACC=1 go test ./internal/resources_acceptance_tests/... -v -timeout 60m

# Run only custom field tests (serial execution)
test-acceptance-customfields:
	TF_ACC=1 go test -tags=customfields ./internal/resources_acceptance_tests_customfields/... -v -timeout 90m -p 1

# Run all acceptance tests (parallel + serial)
test-acceptance-all: test-acceptance test-acceptance-customfields

# Fast test (unit tests only, no acceptance)
test-fast:
	go test ./internal/resources_unit_tests/... -v
```

### 3. Update `.github/workflows/test.yml`
```yaml
- name: Run Acceptance Tests (Parallel)
  run: make test-acceptance
  env:
    TF_ACC: "1"
    NETBOX_SERVER_URL: ${{ secrets.NETBOX_SERVER_URL }}
    NETBOX_API_TOKEN: ${{ secrets.NETBOX_API_TOKEN }}

- name: Run Custom Field Tests (Serial)
  run: make test-acceptance-customfields
  env:
    TF_ACC: "1"
    NETBOX_SERVER_URL: ${{ secrets.NETBOX_SERVER_URL }}
    NETBOX_API_TOKEN: ${{ secrets.NETBOX_API_TOKEN }}
```

### 4. Update `CONTRIBUTING.md`
Add section:
```markdown
### Running Tests

#### Fast Development Cycle (Parallel Tests Only)
```bash
make test-acceptance
```
Runs ~150 parallel-safe tests in 30-40 minutes.

#### Custom Field Tests (Serial Execution)
```bash
make test-acceptance-customfields
```
Runs 41 custom field tests serially in 60-90 minutes.

#### Full Test Suite
```bash
make test-acceptance-all
```
Runs all acceptance tests (2-3 hours total).
```

### 5. Update `DEVELOPMENT.md`
Add explanation of the split test architecture.

---

## Testing Strategy

### After Each Batch:
1. **Compile check**: `go build ./...`
2. **Original package**: `go test ./internal/resources_acceptance_tests -run TestAccAggregate -v`
3. **New package**: `go test -tags=customfields ./internal/resources_acceptance_tests_customfields -run TestAccAggregate -v`
4. **Git commit**: `git commit -m "Migrate batch X custom field tests to separate package"`

### Final Verification:
1. Run full parallel test suite: `make test-acceptance`
2. Run full custom field suite: `make test-acceptance-customfields`
3. Verify both pass without conflicts
4. Measure time savings (expect 30-40 minute reduction in default runs)

---

## Rollback Plan
If issues arise:
1. Each batch is a separate commit - can revert individual batches
2. Build tags prevent compilation issues - old tests still work
3. Original test files remain - just missing specific tests
4. Can move tests back by reversing the process

---

## Success Metrics
- ✅ All 41 custom field tests migrated
- ✅ Normal test runs reduced by 30-40 minutes
- ✅ No parallel execution conflicts
- ✅ CI pipeline runs both test suites
- ✅ Documentation updated
- ✅ Clear separation of concerns

---

## Estimated Total Time
- **Phase 1 Setup**: 1 hour
- **Phase 2 Migration**: 4-5 hours (9 batches)
- **Testing & Documentation**: 1 hour
- **Total**: 6-7 hours

---

## Notes
- Tests can still run individually: `go test -tags=customfields -run TestAccAggregateResource_importWithCustomFieldsAndTags`
- The `-p 1` flag in Makefile forces serial execution for custom field tests
- Build tag `customfields` prevents accidental parallel execution
- Each resource's custom field test is completely independent - no shared state
