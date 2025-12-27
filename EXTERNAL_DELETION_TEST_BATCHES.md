# External Deletion Test Implementation - Batch Progress

## Current Status
- **Overall Progress**: 8/99 resources (8.1%) - Batch 1 Complete
- **Last Update**: December 28, 2025 - Batch 1 Committed
- **Strategy**: Batch-by-batch implementation with commits after each batch

---

## Batch 1: Core Infrastructure Resources (8 resources)
**Status**: ✅ COMPLETE

Priority: CRITICAL - Most frequently used resources

| Resource | File | Status | Notes |
|----------|------|--------|-------|
| device | device_resource_test.go | ✅ DONE | Core infrastructure |
| interface | interface_resource_test.go | ✅ DONE | Critical for networking |
| ip_address | ip_address_resource_test.go | ✅ DONE | Core IPAM resource |
| site | site_resource_test.go | ✅ DONE | Top-level organizational |
| vlan | vlan_resource_test.go | ✅ DONE | Network configuration |
| circuit | circuit_resource_test.go | ✅ DONE | Core interconnect |
| cable | cable_resource_test.go | ✅ DONE | Physical layer |
| cluster | cluster_resource_test.go | ✅ DONE | Virtualization core |

---

## Batch 2: Port Resources (15 resources)
**Status**: ⏳ NOT STARTED

Similar implementation patterns, can be templated.

| Resource | File | Status |
|----------|------|--------|
| console_port | console_port_resource_test.go | ⏳ TODO |
| console_port_template | console_port_template_resource_test.go | ⏳ TODO |
| console_server_port | console_server_port_resource_test.go | ⏳ TODO |
| console_server_port_template | console_server_port_template_resource_test.go | ⏳ TODO |
| front_port | front_port_resource_test.go | ⏳ TODO |
| front_port_template | front_port_template_resource_test.go | ⏳ TODO |
| rear_port | rear_port_resource_test.go | ⏳ TODO |
| rear_port_template | rear_port_template_resource_test.go | ⏳ TODO |
| power_port | power_port_resource_test.go | ⏳ TODO |
| power_port_template | power_port_template_resource_test.go | ⏳ TODO |
| power_outlet | power_outlet_resource_test.go | ⏳ TODO |
| power_outlet_template | power_outlet_template_resource_test.go | ⏳ TODO |
| inventory_item | inventory_item_resource_test.go | ⏳ TODO |
| inventory_item_template | inventory_item_template_resource_test.go | ⏳ TODO |
| inventory_item_role | inventory_item_role_resource_test.go | ⏳ TODO |

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

---

## Batch Completion Tracking

- [ ] Batch 1 (8) - CURRENT
- [ ] Batch 2 (15)
- [ ] Batch 3 (10)
- [ ] Batch 4 (10)
- [ ] Batch 5 (8)
- [ ] Batch 6 (8)
- [ ] Batch 7 (6)
- [ ] Batch 8 (8)
- [ ] Batch 9 (8)
- [ ] Batch 10 (17)

**Target**: 99/99 resources with external deletion tests
