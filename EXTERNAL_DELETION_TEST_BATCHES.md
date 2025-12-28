3
## Current Status
- **Overall Progress**: 15/99 resources (15.2%) - Batch 2B Complete
- **Last Update**: December 28, 2025 - Batch 2B Completed
- **Strategy**: Sub-batch implementation with commits after each sub-batch

---

## Batch 1: Core Infrastructure Resources (8 resources)
**Status**: ✅ COMPLETE

Priority: CRITICAL - Most frequently used resources

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| device | device_resource_test.go | ✅ DONE | CRUD + Import | Core infrastructure |
| interface | interface_resource_test.go | ✅ DONE | CRUD + Import | Critical for networking |
| ip_address | ip_address_resource_test.go | ✅ DONE | CRUD + Import | Core IPAM resource |
| site | site_resource_test.go | ✅ DONE | CRUD + Import | Top-level organizational |
| vlan | vlan_resource_test.go | ✅ DONE | CRUD + Import | Network configuration |
| circuit | circuit_resource_test.go | ✅ DONE | CRUD + Import | Core interconnect |
| cable | cable_resource_test.go | ✅ DONE | CRUD + Import | Physical layer |
| cluster | cluster_resource_test.go | ✅ DONE | CRUD + Import | Virtualization core |

**Recent Improvements**:
- Added 404 handling to Delete methods for site and ip_address resources
- Added full, update, and import tests to cable resource for complete CRUD coverage
- All external deletion tests now pass (including cleanup phase)

---

## Batch 2A: Inventory Resources - Quick Wins (3 resources)
**Status**: ✅ COMPLETE

Priority: HIGH - Already have complete CRUD coverage, only need external deletion tests

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| inventory_item | inventory_item_resource_test.go | ✅ DONE | CRUD + Import + Ext Del ✅ | All tests complete |
| inventory_item_template | inventory_item_template_resource_test.go | ✅ DONE | CRUD + Import + Ext Del ✅ | All tests complete |
| inventory_item_role | inventory_item_role_resource_test.go | ✅ DONE | CRUD + Import + Ext Del ✅ | All tests complete |

**Implementation Notes**:
- Added context import to all test files
- All external deletion tests follow the standard pattern from Batch 1
- Uses DcimAPI methods: DcimInventoryItemsList/Destroy, DcimInventoryItemTemplatesList/Destroy, DcimInventoryItemRolesList/Destroy

---

## Batch 2B: Console Port Resources (4 resources)
**Status**: ✅ COMPLETE

Priority: HIGH - Most commonly used port types

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| console_port | console_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |
| console_port_template | console_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |
| console_server_port | console_server_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |
| console_server_port_template | console_server_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |

**Implementation Notes**:
- Added context import to all test files
- Added update tests for description/label field changes
- Added external deletion tests following Batch 1 pattern
- Uses DcimAPI methods: DcimConsolePorts*/DcimConsoleServerPorts* List/Destroy and Templates variants

---

## Batch 2C: Power Infrastructure Resources (4 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - Critical for power management

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| power_port | power_port_resource_test.go | ⏳ TODO | Update + External Deletion |
| power_port_template | power_port_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| power_outlet | power_outlet_resource_test.go | ⏳ TODO | Update + External Deletion |
| power_outlet_template | power_outlet_template_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 2D: Patch Panel Port Resources (4 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - Structured cabling

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| front_port | front_port_resource_test.go | ⏳ TODO | Update + External Deletion |
| front_port_template | front_port_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| rear_port | rear_port_resource_test.go | ⏳ TODO | Update + External Deletion |
| rear_port_template | rear_port_template_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 3-10: Other Resources (76 resources)
**Status**: ⏳ NOT STARTED

Will be organized by category and added after Batches 1-2 are complete.

---

## Test Pattern

All external deletion tests follow this structure:
```go
func TestAcc{ResourceName}Resource_externalDeletion(t *testing.T) {
	t.Parallel()
	tex] Batch 1 (8) - ✅ COMPLETE
- [ ] Batch 2A (3) - 🔄 IN PROGRESS
- [ ] Batch 2B (4) - Console Ports
- [ ] Batch 2C (4) - Power Infrastructure
- [ ] Batch 2D (4) - Patch Panel Ports
- [ ] Batch 3 (10)
- [ ] Batch 4 (10)
- [ ] Batch 5 (8)
- [ ] Batch 6 (8)
- [ ] Batch 7 (6)
## Test Pattern

All external deletion tests follow this structure:
```go
func TestAcc{ResourceName}Resource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)
	// Setup: Create resource with unique name
	// Step 1: Create and verify resource exists
	// Step 2: PreConfig hook deletes resource via API
	//         Reapply config (Terraform detects 404, recreates)
	//         Verify ID is set (resource was recreated)
}
```

---

## Implementation Notes

1. Each test verifies that when a resource is deleted externally via the NetBox API, Terraform properly:
   - Detects the 404 error during refresh
   - Marks the resource as missing from state
   - Recreates the resource on next apply

2. Imports required:
   - `"context"` - for API calls
   - `testutil.GetSharedClient()` - to get NetBox API client
   - Resource-specific API methods (e.g., `IpamAPI.IpamDevicesList()`)

3. Safe commit strategy:
   - Commit after each batch
   - Run `go build .` to verify no compilation errors
   - Verify test syntax is correct

---x] Batch 2B (4) - ✅ COMPLETE

## Batch Completion Tracking

- [x] Batch 1 (8) - ✅ COMPLETE
- [x] Batch 2A (3) - ✅ COMPLETE
- [ ] Batch 2B (4) - Console Ports
- [ ] Batch 2C (4) - Power Infrastructure
- [ ] Batch 2D (4) - Patch Panel Ports
- [ ] Batch 3 (10)
- [ ] Batch 4 (10)
- [ ] Batch 5 (8)
- [ ] Batch 6 (8)
- [ ] Batch 7 5/99 (15.2%)
**Next Milestone**: 19/99 (19.2%) after Batch 2C
- [ ] Batch 9 (8)
- [ ] Batch 10 (17)

**Target**: 99/99 resources with external deletion tests
**Current**: 11/99 (11.1%)
**Next Milestone**: 15/99 (15.2%) after Batch 2B
