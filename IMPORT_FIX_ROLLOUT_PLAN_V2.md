# Import State Preservation Fix - Complete Rollout Plan

## Executive Summary

**Problem:** Import doesn't preserve optional fields because helper functions check `IsNull()` and skip population when state is empty (which it always is during import).

**Solution:** Remove `IsNull()` checks from helpers + add PlanOnly validation to ALL import tests.

**Scope:**
- Fix 9 utility helper functions
- Fix explicit mapping logic in 14 resources
- Enhance **~85-90 import tests** across ALL resources (not just 40)

**Timeline:** ~15-18 hours total work across 12+ batches

---

## Discovery: Actual Test Coverage

Initial grep search found 40 resources with `importWithCustomFieldsAndTags` tests. However, complete search reveals:

- **42 comprehensive tests** (`importWithCustomFieldsAndTags`) - includes all 7 custom field types + tags
- **45 basic tests** (`import`) - minimal field validation
- **Total: ~85-90 unique import tests** needing PlanOnly step enhancement

### Resources with BOTH Test Types (need PlanOnly in both):
- aggregate, asn_range, cable, circuit, circuit_termination
- cluster, cluster_type, inventory_item, inventory_item_role, interface
- ip_range, location, rack, site, tenant
- virtual_disk, vlan, vm_interface, vrf

### Resources with ONLY Comprehensive Tests:
- asn, contact_assignment, contact_group, contact_role, console_port
- console_server_port, device_bay, device_role, device_type, front_port
- l2vpn, module, module_bay, power_feed, power_outlet
- power_port, rear_port, tenant_group, virtual_chassis, virtual_device_context
- vlan (also has basic), etc.

### Resources with ONLY Basic Tests:
- circuit_group, circuit_group_assignment, circuit_type, fhrp_group
- ike_policy, ike_proposal, ip_address (also has importWithTags)
- ipsec_policy, ipsec_profile, ipsec_proposal, inventory_item_template
- journal_entry, manufacturer, platform, prefix (also has importWithTags)
- provider, provider_account, provider_network, rack_role
- region, route_target, site_group, tunnel, tunnel_group
- tunnel_termination, vlan_group

---

## Phase 1: Fix Utility Functions (30 minutes)

### Backup & Update `internal/utils/state_helpers.go`

```bash
cp internal/utils/state_helpers.go internal/utils/state_helpers.go.backup
```

Update these 9 functions to remove `IsNull()` checks:

1. **UpdateReferenceAttribute** - Remove `if current.IsNull() { return current }`
2. **StringFromAPI** - Remove IsNull check, always return `types.StringNull()` when API value is nil
3. **StringFromAPIPreserveEmpty** - Same as StringFromAPI
4. **NullableStringFromAPI** - Remove IsNull check
5. **Int64FromAPI** - Remove IsNull check
6. **Int64FromInt32API** - Remove IsNull check
7. **Float64FromAPI** - Remove IsNull check
8. **NullableInt64FromAPI** - Remove IsNull check
9. **NullableFloat64FromAPI** - Remove IsNull check

### Validation
```bash
go test ./internal/utils/... -v
```

---

## Phase 2: Resource Fixes & Test Enhancements

### Part A: Fix Explicit Mapping Logic (14 resources with IsNull checks)

#### Batch 1: ✅ COMPLETE - Aggregate (reference implementation)
- [x] aggregate_resource.go - Removed 4 IsNull checks
- [x] Added PlanOnly to comprehensive import test
- [x] All tests passing

#### Batch 2: IPAM Resources (1 hour)
**Resources:** asn, asn_range, ip_range, tenant

For each resource:
1. Remove `else if data.CustomFields.IsNull()` block from mapResponseToModel
2. Add PlanOnly step to `importWithCustomFieldsAndTags` test
3. Add PlanOnly step to basic `import` test (if exists)

```bash
# Test command
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(Asn|AsnRange|IpRange|Tenant)Resource_import' -v
```

#### Batch 3: Device/Rack Roles & Types (1.5 hours)
**Resources:** device_role, device_type, fhrp_group, rack_role, vm_interface

Mapping fixes + test enhancements for all import tests

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(DeviceRole|DeviceType|FhrpGroup|RackRole|VmInterface)Resource_import' -v
```

#### Batch 4: Network/VM Resources (1 hour)
**Resources:** journal_entry, provider, site_group, tunnel_group, virtual_machine, vlan, vlan_group

Mapping fixes + test enhancements

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(JournalEntry|Provider|SiteGroup|TunnelGroup|VirtualMachine|Vlan|VlanGroup)Resource_import' -v
```

#### Batch 5: Circuit Termination (30 minutes)
**Resources:** circuit_termination (8 IsNull checks total)

Remove IsNull checks for: Site, ProviderNetwork, PortSpeed, UpstreamSpeed, XconnectID, PPInfo, Description, CustomFields

Add PlanOnly to both basic and comprehensive import tests

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAccCircuitTerminationResource_import' -v
```

---

### Part B: Test Enhancement ONLY (No Mapping Changes)

These resources don't have explicit IsNull checks in mapping - just need PlanOnly test enhancements.

#### Batch 6: Tenancy & Contacts (1.5 hours)
**Resources:** contact_assignment, contact_group, contact_role, inventory_item_role, tenant_group

Add PlanOnly to comprehensive tests

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(ContactAssignment|ContactGroup|ContactRole|InventoryItemRole|TenantGroup)Resource_import' -v
```

#### Batch 7: Virtualization (2 hours)
**Resources:** cluster, cluster_group, cluster_type, virtual_chassis, virtual_device_context, virtual_disk, vrf

Add PlanOnly to both basic and comprehensive tests where applicable

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(Cluster|ClusterGroup|ClusterType|VirtualChassis|VirtualDeviceContext|VirtualDisk|Vrf)Resource_import' -v
```

#### Batch 8: Device Components Part 1 (2 hours)
**Resources:** console_port, console_server_port, device_bay, front_port, interface, inventory_item, module, module_bay

Add PlanOnly to all import tests (basic + comprehensive)

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(ConsolePort|ConsoleServerPort|DeviceBay|FrontPort|Interface|InventoryItem|Module|ModuleBay)Resource_import' -v
```

#### Batch 9: Device Components Part 2 (2 hours)
**Resources:** power_feed, power_outlet, power_port, location, rear_port, rack

Add PlanOnly to all import tests

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(PowerFeed|PowerOutlet|PowerPort|Location|RearPort|Rack)Resource_import' -v
```

#### Batch 10: Infrastructure (2 hours)
**Resources:** cable, circuit, l2vpn, site

Add PlanOnly to all import tests

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(Cable|Circuit|L2vpn|Site)Resource_import' -v
```

#### Batch 11: VPN & Security (1.5 hours)
**Resources:** ike_policy, ike_proposal, ipsec_policy, ipsec_profile, ipsec_proposal, tunnel, tunnel_group, tunnel_termination

Add PlanOnly to basic import tests

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(Ike|Ipsec|Tunnel).*Resource_import' -v
```

#### Batch 12: Supporting Resources (1.5 hours)
**Resources:** circuit_group, circuit_group_assignment, circuit_type, ip_address, inventory_item_template, manufacturer, platform, prefix, provider_account, provider_network, rack_role, region, route_target, vlan_group

Add PlanOnly to basic import tests (some have importWithTags variants)

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(CircuitGroup|CircuitType|IpAddress|InventoryItemTemplate|Manufacturer|Platform|Prefix|Provider|RackRole|Region|RouteTarget|VlanGroup)Resource_import' -v
```

#### Batch 13: Device & VM (special handling) (1 hour)
**Resources:** device, virtual_machine

These have dedicated import test files. Add PlanOnly to comprehensive tests.

```bash
TF_ACC=1 go test ./internal/resources_acceptance_tests \
  -run 'TestAcc(Device|VirtualMachine)Resource_import' -v
```

---

## Phase 3: Verification (1 hour)

1. **Full acceptance test suite:**
   ```bash
   TF_ACC=1 go test ./internal/resources_acceptance_tests/... -v -timeout 120m
   ```

2. **Spot-check key resources manually:**
   - aggregate (reference implementation)
   - circuit_termination (most complex mapping)
   - device (large resource)
   - virtual_machine (VM context)

3. **Verify no regressions** in existing CRUD tests

---

## Phase 4: Cleanup (30 minutes)

```bash
# Delete temporary files
rm internal/utils/state_helpers_fixed.go
rm internal/utils/state_helpers.go.backup
rm internal/examples/import_fix_pattern.go

# Update documentation
# Update CONTRIBUTING.md with correct patterns
git add -A
git commit -m "chore: Clean up import fix temporary files and update docs"
```

---

## Implementation Checklist

### Phase 1: Utilities (30 minutes)
- [ ] Backup state_helpers.go
- [ ] Update 9 helper functions (remove IsNull checks)
- [ ] Run utility tests: `go test ./internal/utils/...`
- [ ] Delete state_helpers_fixed.go

### Phase 2: Resource Updates (13 batches, ~16 hours)

**Part A: Mapping Fixes** (5 batches, ~4.5 hours)
- [x] Batch 1: aggregate ✅ COMPLETE
- [ ] Batch 2: IPAM (4 resources) - 1 hour
- [ ] Batch 3: Device/Rack (5 resources) - 1.5 hours
- [ ] Batch 4: Network/VM (7 resources) - 1 hour
- [ ] Batch 5: Circuit termination (1 resource) - 30 minutes

**Part B: Test Enhancements Only** (8 batches, ~11.5 hours)
- [ ] Batch 6: Tenancy/Contacts (5 resources) - 1.5 hours
- [ ] Batch 7: Virtualization (7 resources) - 2 hours
- [ ] Batch 8: Device Components Part 1 (8 resources) - 2 hours
- [ ] Batch 9: Device Components Part 2 (6 resources) - 2 hours
- [ ] Batch 10: Infrastructure (4 resources) - 2 hours
- [ ] Batch 11: VPN/Security (8 resources) - 1.5 hours
- [ ] Batch 12: Supporting (14 resources) - 1.5 hours
- [ ] Batch 13: Device/VM special (2 resources) - 1 hour

### Phase 3: Verification (1 hour)
- [ ] Run full acceptance test suite
- [ ] Manual spot-checks (4 resources)
- [ ] Verify no regressions

### Phase 4: Cleanup (30 minutes)
- [ ] Delete temporary files
- [ ] Update CONTRIBUTING.md
- [ ] Final commit

---

## Success Criteria

1. ✅ All 9 helper functions always populate from API (no IsNull checks)
2. ✅ All 14 resources with explicit IsNull checks fixed
3. ✅ **ALL ~85-90 import tests** have PlanOnly validation step
4. ✅ All acceptance tests pass
5. ✅ No temporary/duplicate files remain
6. ✅ Documentation updated

---

## Timeline Summary

| Phase | Work | Duration |
|-------|------|----------|
| Phase 1 | Fix utilities | 30 minutes |
| Phase 2A | Fix mappings (5 batches) | 4.5 hours |
| Phase 2B | Enhance tests (8 batches) | 11.5 hours |
| Phase 3 | Verification | 1 hour |
| Phase 4 | Cleanup | 30 minutes |
| **TOTAL** | **13 batches** | **~18 hours** |

**Already complete:** Batch 1 (aggregate) - 1 hour

**Remaining:** ~17 hours over 12 batches

---

## Key Differences from Original Plan

1. **Scope correction:** 85-90 tests (not 40) need PlanOnly steps
2. **More batches:** 13 batches (not 9) for manageable work chunks
3. **Dual coverage:** Many resources have both basic + comprehensive tests - both need updates
4. **Timeline:** 18 hours (not 10.5) due to larger scope
5. **Organization:** Grouped by resource type for easier context switching

---

## Notes

- **ImportStateVerifyIgnore:** Reference fields that change format (ID→name/slug) need to be in ImportStateVerifyIgnore - this is expected behavior
- **Test Pattern:** create resource → import → **PlanOnly step** (no changes expected)
- **Commit strategy:** One commit per batch for easy rollback
- **Testing:** Run batch tests after each batch before moving to next
