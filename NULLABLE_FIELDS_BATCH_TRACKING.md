# Nullable Fields Bug Fix - Batch Tracking

## Current Status: Batch 1 (High Priority) - IN PROGRESS ⏳

**Last Updated**: 2026-01-09

---

## Batch 0: Foundation ✅ COMPLETED
**Status**: Complete
**Commit**: a516eb0 - Initial foundation with ASN fix and test pattern

### Completed
- [x] Create bug fix branch: `bugfix/nullable-field-removal`
- [x] Fix ASN resource (tenant, rir) with SetNil pattern
- [x] Create ASN test: `TestAccASNResource_removeOptionalFields`
- [x] Verify test passes with fix
- [x] Create comprehensive planning document: `NULLABLE_FIELDS_BUGFIX_PLAN.md`
- [x] Identify all 22 affected resources and 47 nullable fields

### Files Changed
- `internal/resources/asn_resource.go` - Added SetRirNil() and SetTenantNil()
- `internal/resources_acceptance_tests/asn_resource_test.go` - Added removeOptionalFields test

---

## Batch 1: High Priority - Tenant Fields ⏳ IN PROGRESS
**Target**: Resources with `tenant` field (most frequently used)
**Estimated Time**: 1-2 hours
**Status**: 1/7 complete

### Resources (7)
- [x] **asn_range** - Fields: tenant (1 field) ✅
  - [x] Code: Add SetTenantNil()
  - [x] Test: TestAccASNRangeResource_removeOptionalFields
  - [x] Verify: Build + test pass
  - Note: RIR is required, not nullable

- [ ] **circuit** - Fields: tenant (1 field)
  - [ ] Code: Add SetTenantNil()
  - [ ] Test: TestAccCircuit_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **ip_address** - Fields: vrf, tenant (2 fields)
  - [ ] Code: Add SetVrfNil() and SetTenantNil()
  - [ ] Test: TestAccIPAddress_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **ip_range** - Fields: vrf, tenant, role (3 fields)
  - [ ] Code: Add SetVrfNil(), SetTenantNil(), SetRoleNil()
  - [ ] Test: TestAccIPRange_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **route_target** - Fields: tenant (1 field)
  - [ ] Code: Add SetTenantNil()
  - [ ] Test: TestAccRouteTarget_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **vrf** - Fields: tenant (1 field)
  - [ ] Code: Add SetTenantNil()
  - [ ] Test: TestAccVRF_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **wireless_link** - Fields: tenant (1 field)
  - [ ] Code: Add SetTenantNil()
  - [ ] Test: TestAccWirelessLink_removeOptionalFields
  - [ ] Verify: Build + test pass

### Batch 1 Completion Checklist
- [ ] All 7 resources code complete
- [ ] All 7 tests passing
- [ ] Run full acceptance suite: `go test ./internal/resources_acceptance_tests/...`
- [ ] Commit: `fix(batch1): Handle nullable field removal for tenant-related resources`

**Notes**: Start with route_target or vrf (simplest - 1 field each) to build momentum.

---

## Batch 2: Infrastructure Resources ⏸️ PENDING
**Target**: Site-related and location resources
**Estimated Time**: 45-60 minutes
**Status**: Not Started

### Resources (4)
- [ ] **site** - Fields: tenant, region, group (3 fields)
  - [ ] Code: Add SetTenantNil(), SetRegionNil(), SetGroupNil()
  - [ ] Test: TestAccSite_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **location** - Fields: parent, tenant (2 fields)
  - [ ] Code: Add SetParentNil(), SetTenantNil()
  - [ ] Test: TestAccLocation_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **cluster** - Fields: group, tenant, site (3 fields)
  - [ ] Code: Add SetGroupNil(), SetTenantNil(), SetSiteNil()
  - [ ] Test: TestAccCluster_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **tenant** - Fields: group (1 field)
  - [ ] Code: Add SetGroupNil()
  - [ ] Test: TestAccTenant_removeOptionalFields
  - [ ] Verify: Build + test pass

### Batch 2 Completion Checklist
- [ ] All 4 resources code complete
- [ ] All 4 tests passing
- [ ] Run full acceptance suite
- [ ] Commit: `fix(batch2): Handle nullable field removal for infrastructure resources`

---

## Batch 3: VLAN/Prefix Resources ⏸️ PENDING
**Target**: Networking resources with multiple nullable fields
**Estimated Time**: 45-60 minutes
**Status**: Not Started

### Resources (3)
- [ ] **prefix** - Fields: site, vrf, tenant, vlan, role (5 fields)
  - [ ] Code: Add SetSiteNil(), SetVrfNil(), SetTenantNil(), SetVlanNil(), SetRoleNil()
  - [ ] Test: TestAccPrefix_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **vlan** - Fields: site, group, tenant, role (4 fields)
  - [ ] Code: Add SetSiteNil(), SetGroupNil(), SetTenantNil(), SetRoleNil()
  - [ ] Test: TestAccVLAN_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **vm_interface** - Fields: untagged_vlan, vrf (2 fields)
  - [ ] Code: Add SetUntaggedVlanNil(), SetVrfNil()
  - [ ] Test: TestAccVMInterface_removeOptionalFields
  - [ ] Verify: Build + test pass

### Batch 3 Completion Checklist
- [ ] All 3 resources code complete
- [ ] All 3 tests passing
- [ ] Run full acceptance suite
- [ ] Commit: `fix(batch3): Handle nullable field removal for VLAN/prefix resources`

**Notes**: Prefix has 5 fields - most complex resource in this batch.

---

## Batch 4: Device/Rack Resources ⏸️ PENDING
**Target**: Physical infrastructure resources
**Estimated Time**: 1 hour
**Status**: Not Started

### Resources (4)
- [ ] **rack** - Fields: location, tenant, role, rack_type (4 fields)
  - [ ] Code: Add SetLocationNil(), SetTenantNil(), SetRoleNil(), SetRackTypeNil()
  - [ ] Test: TestAccRack_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **device_bay** - Fields: installed_device (1 field)
  - [ ] Code: Add SetInstalledDeviceNil()
  - [ ] Test: TestAccDeviceBay_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **platform** - Fields: manufacturer (1 field)
  - [ ] Code: Add SetManufacturerNil()
  - [ ] Test: TestAccPlatform_removeOptionalFields
  - [ ] Verify: Build + test pass

- [ ] **virtual_machine** - Fields: site, cluster, role, tenant, platform (5 fields)
  - [ ] Code: Add SetSiteNil(), SetClusterNil(), SetRoleNil(), SetTenantNil(), SetPlatformNil()
  - [ ] Test: TestAccVirtualMachine_removeOptionalFields
  - [ ] Verify: Build + test pass

### Batch 4 Completion Checklist
- [ ] All 4 resources code complete
- [ ] All 4 tests passing
- [ ] Run full acceptance suite
- [ ] Commit: `fix(batch4): Handle nullable field removal for device/rack resources`

**Notes**: Virtual machine has 5 fields - tied with prefix for most complex.

---

## Batch 5: Cleanup & Partial Fixes ⏸️ PENDING
**Target**: Fix partially implemented resources
**Estimated Time**: 30-45 minutes
**Status**: Not Started

### Resources (3)
- [ ] **cable** - Fields: tenant (1 field)
  - Issue: Uses NewNullable(nil) only in Update
  - [ ] Code: Add NewNullable(nil) to Create function
  - [ ] Test: Verify TestAccCable tests cover removal
  - [ ] Verify: Build + test pass

- [ ] **circuit_group** - Fields: tenant (1 field)
  - Issue: Uses NewNullable(nil) only in Update
  - [ ] Code: Add NewNullable(nil) to Create function
  - [ ] Test: Verify TestAccCircuitGroup tests cover removal
  - [ ] Verify: Build + test pass

- [ ] **l2vpn** - Fields: tenant, identifier (2 fields)
  - Issue: identifier has SetNil, tenant uses NewNullable(nil)
  - [ ] Code: Standardize tenant to use SetTenantNil()
  - [ ] Test: Verify TestAccL2VPN tests cover removal
  - [ ] Verify: Build + test pass

### Batch 5 Completion Checklist
- [ ] All 3 resources code complete
- [ ] All 3 tests passing
- [ ] Run full acceptance suite
- [ ] Commit: `fix(batch5): Standardize nullable field handling for remaining resources`

**Notes**: These are minor fixes to resources that already partially handle the pattern.

---

## Final Verification ⏸️ PENDING
**Status**: Not Started

### Pre-Release Checklist
- [ ] All 5 batches completed (22 resources total)
- [ ] Full acceptance test suite passes: `go test -v ./internal/resources_acceptance_tests/... -timeout 120m`
- [ ] Build passes: `go build .`
- [ ] Run pre-commit hooks: `make lint` or `.venv/Scripts/pre-commit.exe run --all-files`
- [ ] Review all commits in branch
- [ ] Update CHANGELOG.md with v0.0.14 entry
- [ ] Create PR with comprehensive description

### PR Description Template
```markdown
# Fix: Provider produced inconsistent result when removing nullable fields

## Problem
Production bug where removing optional nullable reference fields (tenant, site, rir, etc.) from resources caused "Provider produced inconsistent result after apply" errors.

## Root Cause
Resources omitted nullable fields from API requests when null in config, causing API to preserve existing values instead of clearing them.

## Solution
Explicitly set nullable fields to null using SetFieldNil() when removed from configuration.

## Changes
- Fixed 22 resources to handle nullable field removal correctly
- Added 22 new acceptance tests: TestAccXxx_removeOptionalFields
- Added comprehensive planning document: NULLABLE_FIELDS_BUGFIX_PLAN.md

## Testing
- All new tests pass (22 new tests)
- All existing acceptance tests pass
- Manual verification with production case

## Resources Fixed
Batch 1: asn, asn_range, circuit, ip_address, ip_range, route_target, vrf, wireless_link
Batch 2: site, location, cluster, tenant
Batch 3: prefix, vlan, vm_interface
Batch 4: rack, device_bay, platform, virtual_machine
Batch 5: cable, circuit_group, l2vpn (cleanup)

Closes #XX (if issue exists)
```

### Release Checklist
- [ ] PR approved and merged to main
- [ ] Tag v0.0.14: `git tag -a v0.0.14 -m "Release v0.0.14: Fix nullable field removal bug"`
- [ ] Push tag: `git push origin v0.0.14`
- [ ] Create GitHub release with CHANGELOG excerpt
- [ ] Notify user who reported the bug

---

## Progress Summary

### Overall Status
- **Total Resources**: 22
- **Completed**: 1 (asn - foundation)
- **Remaining**: 21
- **Total Fields**: 47 nullable fields

### Batch Status
| Batch | Resources | Status | Time Estimate |
|-------|-----------|--------|---------------|
| 0 - Foundation | 1 | ✅ Complete | - |
| 1 - High Priority | 7 | ⏳ Ready | 1-2 hours |
| 2 - Infrastructure | 4 | ⏸️ Pending | 45-60 min |
| 3 - VLAN/Prefix | 3 | ⏸️ Pending | 45-60 min |
| 4 - Device/Rack | 4 | ⏸️ Pending | 1 hour |
| 5 - Cleanup | 3 | ⏸️ Pending | 30-45 min |
| Final Verification | - | ⏸️ Pending | 1 hour |

**Total Remaining Time**: ~5-6 hours

---

## Next Steps

1. **Start Batch 1**: Begin with simplest resources (route_target or vrf - 1 field each)
2. **Momentum Building**: Complete 1-2 resources to validate pattern
3. **Batch Completion**: Finish all Batch 1 resources
4. **Test & Commit**: Run tests and commit Batch 1
5. **Repeat**: Continue with Batches 2-5

---

## Notes & Lessons Learned

### Batch 0
- ASN resource serves as perfect reference implementation
- Test pattern validated: create with field → remove field → re-add field
- `TestCheckNoResourceAttr` is key for verifying null values
- Both `SetFieldNil()` and `NewNullable(nil)` patterns work equivalently

### Tips for Implementation
- Copy ASN resource code as template
- Each resource typically takes 10-15 minutes (code + test)
- Run individual tests frequently: `go test -run TestAccXxx_removeOptionalFields`
- Keep test configs simple - focus on the nullable field being tested
- Use testutil helpers for random names/slugs

### Common Gotchas
- Remember to add cleanup registration in tests
- Ensure both Create AND Update functions handle the null case
- Some fields may not have SetFieldNil() - use NewNullable(nil) fallback
- Test step count matters for error messages - always use descriptive configs
