# Fix display_name and reference attribute display issues

## Summary
This PR fixes bugs reported after v0.0.8 where reference attributes were displaying numeric IDs instead of user-friendly names in Terraform plans, and the `display_name` field was causing unnecessary "(known after apply)" noise.

## Issues Fixed
Fixes issues reported by @jeanmarc77 after v0.0.8 release:
- Reference attributes showing IDs: `tenant = "My Tenant" -> "42"`
- display_name showing "(known after apply)" unnecessarily
- 18 acceptance test failures
- Various test reliability issues

## Changes Made

### 1. Removed display_name Field (Breaking Change - Low Impact)
- **Files**: 100 resource files + 1 utility file
- **Why**: The field was non-functional and confusing
  - Returned the resource's own Display value (e.g., "eth0" for interface named "eth0")
  - Did NOT show friendly names for referenced resources (tenant, vlan, cluster)
  - Caused "(known after apply)" noise in plans
- **Impact**: Low - Field never worked as intended, users shouldn't be referencing it
- **Migration**: Remove any `.display_name` references from configurations

### 2. Fixed Reference Attribute Display
- **Files**: internal/utils/state_helpers.go, 1 test file
- **Fix**: Updated `UpdateReferenceAttribute` to prefer name/slug over ID when state is Unknown
- **Result**: Reference attributes now preserve user-specified format
  - Before: `tenant = "My Tenant" -> "42"` ❌
  - After: `tenant = "My Tenant"` ✅
- **Scope**: Affects 47 resources with reference fields:
  - DCIM: Device, Interface, Cable, Rack, PowerPanel, PowerFeed, Module, etc.
  - IPAM: IPAddress, Prefix, VLAN, VRF, Aggregate, ASN, etc.
  - Virtualization: VirtualMachine, VMInterface, VirtualDisk, Cluster
  - Circuits: Circuit, CircuitTermination, CircuitGroupAssignment
  - Others: ContactAssignment, WirelessLAN, Tunnel, L2VPN, etc.

### 3. Fixed Acceptance Test Failures
Fixed 18 test failures in 3 groups:

**Group 1: display_name Schema Issues (3 tests)**
- CircuitGroup, ClusterGroup, ConfigContext
- Removed display_name field from schema

**Group 2: Reference Persistence (2 tests)**
- CircuitGroup: Added missing tenant handling
- ConfigContext: Fixed JSON normalization for data field

**Group 3: External Deletion Handling (3 tests)**
- ClusterGroup, Contact, ContactAssignment
- Fixed 404 error handling in cleanup functions

### 4. Improved Test Reliability
- **7 tests updated** with unique email addresses for parallel execution
- **MAC address handling** normalized (case-insensitive)
- **JSON data** properly normalized in ConfigContext tests

## Testing

### Manual Testing
- ✅ Created test Terraform stack with VM using name-based references
- ✅ Verified references display correctly: `cluster = "test-cluster-vm-ref"`
- ✅ Verified no drift after apply
- ✅ Verified updates work correctly

### Automated Testing
- ✅ `go build .` - Build successful
- ✅ `go vet ./...` - No issues found
- ✅ Unit tests - All passed (resources, datasources, utils)
- ✅ Acceptance tests - All 150+ consistency and reference preservation tests passed
- ✅ 18 previously failing tests now pass
- ✅ 7 reference preservation tests verify correct behavior

## Verification

### Reference Preservation Tests
Added/verified comprehensive test coverage:
1. `TestAccConsistency_VMInterface` - vlan reference
2. `TestAccConsistency_VirtualMachine_PlatformNamePersistence` - platform reference
3. `TestAccConsistency_Device_LiteralNames` - site, tenant, role, device_type
4. `TestAccConsistency_Prefix` - site, tenant, vlan
5. `TestAccConsistency_VLAN` - site, group, tenant, role
6. `TestAccConsistency_Cluster_LiteralNames` - type, group, site, tenant
7. `TestAccReferenceNamePersistence_IPAddress_TenantVRF` - tenant, vrf (NEW)

All tests use PlanOnly steps to verify no drift occurs.

## Statistics
- **109 files modified**
- **100 resources** - display_name removed
- **18 test failures** fixed
- **7 tests** improved for reliability
- **All tests passing** ✅

## Breaking Changes
- **display_name field removed from all resources**
  - Impact: Low (field was non-functional)
  - Migration: Remove any `.display_name` references
  - Benefit: Cleaner plans, less confusion

## Migration Guide
If you were referencing `.display_name` in your configurations or outputs:
```hcl
# Before (broken)
output "vm_display" {
  value = netbox_virtual_machine.vm.display_name
}

# After (use the name field)
output "vm_display" {
  value = netbox_virtual_machine.vm.name
}
```

## Related Work
This PR sets the foundation for future cleanup:
- Phase 5 (planned for v1.0.0): Remove duplicate `_id` fields
  - ~35 duplicate fields across ~20 resources
  - Will eliminate remaining plan noise from computed fields
  - See CLEANUP_PLAN.md for details

## Checklist
- [x] Code follows project conventions
- [x] All tests pass
- [x] CHANGELOG.md updated
- [x] Breaking changes documented
- [x] Migration guide provided
- [x] Manual testing completed
