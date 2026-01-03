# Import State Preservation Fix - Quick Start Guide

## What's This About?

The import functionality wasn't preserving optional fields from NetBox into Terraform state. This caused immediate "apply required" situations after import.

**Root Cause**: Helper functions check `if current.IsNull()` and skip populating. During import, state starts empty (all null), so fields never get populated.

**Solution**: Always populate from API response, regardless of current state.

## Current Status

✅ **Batch 1 COMPLETE**: Aggregate resource fixed and tested
- Removed problematic `else if data.Field.IsNull()` blocks
- All 7 acceptance tests passing
- Terraform integration test passing

## What You Need to Do

### Quick Overview (4 hours total)
1. **Phase 1**: Fix utility functions (30 min)
2. **Phase 2**: Fix 13 affected resources in 2 batches (2.5 hours)
3. **Phase 3**: Verify everything works (1 hour)
4. **Phase 4**: Clean up temporary files (30 min)

### Detailed Steps

#### Step 1: Fix Utility Functions (internal/utils/state_helpers.go)

Run these commands to apply the fix:

```bash
# Backup current file
cp internal/utils/state_helpers.go internal/utils/state_helpers.go.backup

# Open state_helpers.go and for each helper function, find and remove:
# } else if !current.IsNull() {
#     return types.XXXNull()
# }
# return current
#
# Replace with just:
# return types.XXXNull()

# Test
go test ./internal/utils/... -v

# If tests pass, delete the temporary fixed file
rm internal/utils/state_helpers_fixed.go
```

**Functions to fix**:
- UpdateReferenceAttribute
- StringFromAPI
- StringFromAPIPreserveEmpty
- NullableStringFromAPI
- Int64FromAPI
- Int64FromInt32API
- NullableInt64FromAPI
- Float64FromAPI
- NullableFloat64FromAPI

See IMPORT_FIX_ROLLOUT_PLAN.md Phase 1 for exact changes.

#### Step 2: Fix Resources (internal/resources/)

**Batch 2 - CustomFields Pattern (12 resources)**:

For each resource, find and remove this pattern:
```go
} else if data.CustomFields.IsNull() {
    // Keep null if it was null
} else {
    data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
}
```

Replace with:
```go
} else {
    data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
}
```

Resources:
1. device_role_resource.go
2. device_type_resource.go
3. fhrp_group_resource.go
4. journal_entry_resource.go
5. provider_resource.go
6. rack_role_resource.go
7. site_group_resource.go
8. tunnel_group_resource.go
9. virtual_machine_resource.go
10. vlan_resource.go
11. vlan_group_resource.go
12. vm_interface_resource.go

Test after fixing all:
```bash
TF_ACC=1 go test -v -run 'TestAcc(DeviceRole|DeviceType|FhrpGroup|JournalEntry|Provider|RackRole|SiteGroup|TunnelGroup|VirtualMachine|Vlan|VlanGroup|VmInterface)Resource' ./internal/resources_acceptance_tests/ -timeout 30m
```

**Batch 3 - Circuit Termination (1 resource)**:

circuit_termination_resource.go has 8 fields with the pattern:
- Site, ProviderNetwork, PortSpeed, UpstreamSpeed, XconnectID, PPInfo, Description, CustomFields

Apply the same fix to all 8.

Test:
```bash
TF_ACC=1 go test -v -run 'TestAccCircuitTerminationResource' ./internal/resources_acceptance_tests/ -timeout 15m
```

#### Step 3: Verify Everything Works

```bash
# Run all acceptance tests (optional - takes ~2 hours)
TF_ACC=1 go test -v ./internal/resources_acceptance_tests/ -timeout 120m

# Or just test affected resources
TF_ACC=1 go test -v -run 'TestAcc.*Resource_import' ./internal/resources_acceptance_tests/ -timeout 30m

# Test key resources with Terraform
.\scripts\run-terraform-tests.ps1 -TestDir "test\terraform\resources\aggregate"
.\scripts\run-terraform-tests.ps1 -TestDir "test\terraform\resources\virtual_machine"
.\scripts\run-terraform-tests.ps1 -TestDir "test\terraform\resources\circuit_termination"
```

#### Step 4: Clean Up

```bash
# Delete temporary files
rm internal/utils/state_helpers_fixed.go
rm internal/utils/state_helpers.go.backup
rm internal/examples/import_fix_pattern.go

# Commit
git add -A
git commit -m "fix: Import state preservation across all resources

- Updated utility functions to always populate from API
- Fixed 13 resources with explicit IsNull checks
- All affected resources tested and passing
- Removed temporary implementation files

Resolves import issues where optional fields were lost during import."
```

## Testing Checklist

After each batch, verify:
- [ ] Acceptance tests pass
- [ ] No new test failures
- [ ] Import tests don't show unexpected changes
- [ ] Reference fields can change format (ID→name/slug) - that's OK

## Common Issues

**Q: Tests fail with "attributes not equivalent" for reference fields**
A: Add the reference field to ImportStateVerifyIgnore. Format can change (ID→name/slug).

**Q: Should I set fields to Computed: true?**
A: No! Keep them Optional only. Computed causes other issues.

**Q: Import shows plan changes**
A: Make sure you're importing with config that matches the resource state, not minimal config.

## Files Reference

- **IMPORT_FIX_ROLLOUT_PLAN.md**: Complete detailed plan
- **IMPORT_FIX_IMPLEMENTATION.md**: What we did for aggregate (reference)
- **IMPORT_STATE_PRESERVATION_FIX.md**: Original analysis and patterns
- **state_helpers_fixed.go**: Reference for correct helper implementations

## Quick Commands

```bash
# Find all resources with the problematic pattern
grep -r "else if data\.\w\+\.IsNull()" internal/resources/*_resource.go

# Run tests for specific resource
TF_ACC=1 go test -v -run 'TestAccXXXResource' ./internal/resources_acceptance_tests/ -timeout 15m

# Run all import tests
TF_ACC=1 go test -v -run 'TestAcc.*import' ./internal/resources_acceptance_tests/ -timeout 30m
```

## Need Help?

See the detailed files:
1. IMPORT_FIX_ROLLOUT_PLAN.md - Complete step-by-step plan
2. IMPORT_FIX_IMPLEMENTATION.md - What was done for aggregate (example)
3. state_helpers_fixed.go - Reference implementations
