# Nullable Reference Fields Bug Fix Plan

## Issue Summary

**Production Bug**: Provider produces inconsistent result after apply when removing optional nullable reference fields from resources.

**Root Cause**: Resources don't explicitly set nullable fields to `null` when they are removed from configuration. Instead, they omit the field from API requests, causing the API to preserve the existing value. This violates Terraform's expectation that a planned null value results in null after apply.

**Example Error**:
```
Error: Provider produced inconsistent result after apply

When applying changes to netbox_asn.test, provider produced an unexpected new value:
.tenant: was null, but now cty.StringVal("RTX rhou")
```

## Technical Solution

### Pattern: Explicit Nil Handling

When a nullable reference field is null in the Terraform config, we must explicitly set it to null in the API request:

**BEFORE (Incorrect)**:
```go
// Handle Tenant (optional)
if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
    tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
    // ...
    request.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
}
// If null, field is omitted → API preserves existing value
```

**AFTER (Correct)**:
```go
// Handle Tenant (optional)
if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
    tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
    // ...
    request.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
} else if data.Tenant.IsNull() {
    // Explicitly set to null to clear the field
    request.SetTenantNil()
}
```

**Alternative Pattern** (equivalent, used in some resources):
```go
if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
    // ... lookup tenant
    request.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)
} else {
    request.Tenant = *netbox.NewNullableBriefTenantRequest(nil)
}
```

## Test Coverage Gap

### Missing Test Pattern

Our acceptance tests did not cover **removing** optional nullable reference fields:
- ✅ Tests create resources with fields
- ✅ Tests create resources without fields
- ✅ Tests update field values
- ❌ Tests don't remove previously set optional fields

### Required Test Pattern

```go
func TestAccResourceName_removeOptionalFields(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            // Step 1: Create with nullable field set
            {
                Config: testConfigWithField(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttrSet("resource.test", "nullable_field"),
                ),
            },
            // Step 2: Remove nullable field (set to null)
            {
                Config: testConfigWithoutField(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckNoResourceAttr("resource.test", "nullable_field"),
                ),
            },
            // Step 3: Re-add field to verify it can be set again
            {
                Config: testConfigWithField(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttrSet("resource.test", "nullable_field"),
                ),
            },
        },
    })
}
```

## Affected Resources Analysis

### Summary Statistics
- **Total resources with nullable fields**: 22
- **Fully fixed**: 1 (asn_resource)
- **Partially fixed**: 3 (cable, circuit_group, l2vpn)
- **Need fixes**: 18 resources
- **Total nullable fields to fix**: 47 fields

### Resources by Status

#### ✅ Fully Fixed (1)
| Resource | Fields | Status |
|----------|--------|--------|
| asn | tenant, rir | ✅ Complete (SetTenantNil, SetRirNil) |

#### ⚠️ Partially Fixed (3)
| Resource | Fields | Issue |
|----------|--------|-------|
| cable | tenant | Uses NewNullable(nil) in Update only |
| circuit_group | tenant | Uses NewNullable(nil) in Update only |
| l2vpn | tenant, identifier | identifier has SetNil, tenant uses NewNullable(nil) |

#### ❌ Need Fixes (18)

| Resource | Nullable Fields | Field Count |
|----------|----------------|-------------|
| asn_range | tenant, rir | 2 |
| circuit | tenant | 1 |
| cluster | group, tenant, site | 3 |
| device_bay | installed_device | 1 |
| ip_address | vrf, tenant | 2 |
| ip_range | vrf, tenant, role | 3 |
| location | parent, tenant | 2 |
| platform | manufacturer | 1 |
| prefix | site, vrf, tenant, vlan, role | 5 |
| rack | location, tenant, role, rack_type | 4 |
| route_target | tenant | 1 |
| site | tenant, region, group | 3 |
| tenant | group | 1 |
| virtual_machine | site, cluster, role, tenant, platform | 5 |
| vlan | site, group, tenant, role | 4 |
| vm_interface | untagged_vlan, vrf | 2 |
| vrf | tenant | 1 |
| wireless_link | tenant | 1 |

## Implementation Batches

### Batch 1: High Priority (Most Common Fields)
**Resources with `tenant` field** (most frequently used)

| # | Resource | Fields | Test Required |
|---|----------|--------|---------------|
| 1 | asn_range | tenant, rir | ✅ |
| 2 | circuit | tenant | ✅ |
| 3 | ip_address | vrf, tenant | ✅ |
| 4 | ip_range | vrf, tenant, role | ✅ |
| 5 | route_target | tenant | ✅ |
| 6 | vrf | tenant | ✅ |
| 7 | wireless_link | tenant | ✅ |

**Estimated Effort**: 1-2 hours (7 resources × 10-15 min each)

### Batch 2: Infrastructure Resources
**Site-related and location resources**

| # | Resource | Fields | Test Required |
|---|----------|--------|---------------|
| 1 | site | tenant, region, group | ✅ |
| 2 | location | parent, tenant | ✅ |
| 3 | cluster | group, tenant, site | ✅ |
| 4 | tenant | group | ✅ |

**Estimated Effort**: 45-60 minutes (4 resources)

### Batch 3: VLAN/Prefix Resources
**Networking resources with multiple nullable fields**

| # | Resource | Fields | Test Required |
|---|----------|--------|---------------|
| 1 | prefix | site, vrf, tenant, vlan, role | ✅ |
| 2 | vlan | site, group, tenant, role | ✅ |
| 3 | vm_interface | untagged_vlan, vrf | ✅ |

**Estimated Effort**: 45-60 minutes (3 resources, more complex)

### Batch 4: Device/Rack Resources
**Physical infrastructure resources**

| # | Resource | Fields | Test Required |
|---|----------|--------|---------------|
| 1 | rack | location, tenant, role, rack_type | ✅ |
| 2 | device_bay | installed_device | ✅ |
| 3 | platform | manufacturer | ✅ |
| 4 | virtual_machine | site, cluster, role, tenant, platform | ✅ |

**Estimated Effort**: 1 hour (4 resources)

### Batch 5: Cleanup & Partial Fixes
**Fix partially implemented resources**

| # | Resource | Fields | Issue | Fix Required |
|---|----------|--------|-------|--------------|
| 1 | cable | tenant | NewNullable(nil) only in Update | Add to Create |
| 2 | circuit_group | tenant | NewNullable(nil) only in Update | Add to Create |
| 3 | l2vpn | tenant | Uses NewNullable(nil) | Standardize to SetNil |

**Estimated Effort**: 30-45 minutes (3 resources, minor fixes)

## Implementation Checklist

### For Each Resource

#### Code Changes
- [ ] Locate buildRequest function (or equivalent: buildXxxRequest, Create/Update request building)
- [ ] For each nullable reference field:
  - [ ] Add `else if data.Field.IsNull()` branch after the existing field check
  - [ ] Call `request.SetFieldNil()` in the else branch
  - [ ] Verify SetFieldNil() method exists in go-netbox model (fallback to NewNullable(nil) if not)
- [ ] Apply same pattern to BOTH Create and Update functions
- [ ] Verify no other request building locations need updates

#### Test Changes
- [ ] Create `TestAccResourceName_removeOptionalFields` test
- [ ] Test structure:
  - [ ] Step 1: Create resource WITH all nullable fields populated
  - [ ] Step 2: Update to REMOVE nullable fields (omit from config)
  - [ ] Step 3: Verify fields are null using `TestCheckNoResourceAttr`
  - [ ] Step 4: Re-add fields to verify they can be set again
- [ ] Add helper config functions:
  - [ ] `testAccResourceConfig_withFields()` - includes all nullable fields
  - [ ] `testAccResourceConfig_withoutFields()` - omits nullable fields
- [ ] Ensure test uses unique resource IDs (random names/slugs)
- [ ] Add proper cleanup registration

#### Verification
- [ ] Build passes: `go build .`
- [ ] New test passes: `go test -v -run TestAccResourceName_removeOptionalFields`
- [ ] All existing tests still pass: `go test ./internal/resources_acceptance_tests/...`
- [ ] Check test covers all nullable fields in the resource

### Progress Tracking Template

```markdown
## Batch X Progress

- [ ] Resource 1: code_changes + tests
- [ ] Resource 2: code_changes + tests
- [ ] Resource 3: code_changes + tests
- [ ] Batch verification: all tests passing
- [ ] Commit: "fix(batchX): Handle nullable field removal for X resources"
```

## Testing Strategy

### Unit Test Coverage
Not applicable - this is integration-level behavior requiring actual API calls.

### Acceptance Test Coverage
**Required for each resource**:
1. New test: `TestAccResourceName_removeOptionalFields`
2. Verify existing tests don't break
3. Run full acceptance test suite before final merge

### Manual Testing (Optional)
Can be performed against production-like Netbox instance:
1. Create resource with nullable field
2. Remove field from config
3. Run `terraform plan` - should show field changing to null
4. Run `terraform apply` - should succeed without "inconsistent result" error
5. Verify field is null in Netbox UI

## Release Plan

### Version Target
**v0.0.14** (Bug fix release after v0.0.13)

### Release Type
**Patch Release** - Critical bug fix

### Changelog Entry
```markdown
## [0.0.14] - 2026-01-XX

### Fixed
- **CRITICAL**: Fixed "Provider produced inconsistent result after apply" errors when removing optional nullable reference fields (tenant, site, rir, vrf, etc.) from resources. Resources now explicitly set fields to null when removed from configuration instead of omitting them from API requests. (#XX)
  - Affected resources: asn, asn_range, cable, circuit, circuit_group, cluster, device_bay, ip_address, ip_range, l2vpn, location, platform, prefix, rack, route_target, site, tenant, virtual_machine, vlan, vm_interface, vrf, wireless_link
  - Added comprehensive acceptance tests for nullable field removal across all affected resources
```

### Branch Strategy
1. Work branch: `bugfix/nullable-field-removal` ✅ Created
2. Complete all batches in feature branch
3. Run full test suite
4. Create PR to main
5. Merge and tag v0.0.14
6. Release notes in GitHub

### Communication
- Document in PR description the scope of the bug and testing performed
- Emphasize this is a critical fix for production stability
- Note that this adds extensive test coverage for a previously untested scenario

## Time Estimates

### Total Implementation Time
- Batch 1 (High Priority): 1-2 hours
- Batch 2 (Infrastructure): 45-60 minutes
- Batch 3 (VLAN/Prefix): 45-60 minutes
- Batch 4 (Device/Rack): 1 hour
- Batch 5 (Cleanup): 30-45 minutes
- Testing & Verification: 1 hour
- Documentation & PR: 30 minutes

**Total: 5-7 hours** of focused development time

### Recommended Approach
1. ✅ Complete ASN resource (done - serves as reference)
2. Complete Batch 1 (high priority, most common)
3. Run acceptance tests for Batch 1
4. Commit Batch 1
5. Repeat for Batches 2-5
6. Final full test suite run
7. Create PR and release

## Reference Implementation

### ASN Resource (Completed)
- Code: `internal/resources/asn_resource.go` lines 319-342
- Test: `internal/resources_acceptance_tests/asn_resource_test.go` `TestAccASNResource_removeOptionalFields`
- Commit: bugfix/nullable-field-removal branch

### Example Code Pattern
See [asn_resource.go](internal/resources/asn_resource.go#L319-L342) for the complete implementation pattern.

### Example Test Pattern
See [asn_resource_test.go](internal/resources_acceptance_tests/asn_resource_test.go#L147-L194) for the complete test pattern.

## Notes

- Some go-netbox models may not have `SetFieldNil()` methods. In those cases, use `NewNullable*Request(nil)` pattern as fallback.
- The `SetFieldNil()` pattern is preferred as it's more explicit and self-documenting.
- Cable and circuit_group already use `NewNullable(nil)` in Update - verify this is also applied to Create functions for consistency.
- L2VPN has mixed patterns - standardize to SetNil for consistency.

## Risk Assessment

### Low Risk Changes
- Each fix is isolated to a single resource
- Pattern is well-established (ASN resource proves it works)
- Comprehensive test coverage added
- No changes to schema or data structures

### Medium Risk
- Large number of resources affected (22 total)
- Risk of missing edge cases in individual resources

### Mitigation
- Incremental batch approach allows for early detection of issues
- Comprehensive acceptance tests verify behavior
- Reference implementation (ASN) validates pattern
- Can release batches incrementally if needed

## Success Criteria

### Definition of Done
- [ ] All 22 resources handle nullable field removal correctly
- [ ] All 22 resources have `TestAccXxx_removeOptionalFields` tests
- [ ] All acceptance tests pass (existing + new)
- [ ] Build completes without errors
- [ ] Documentation updated (this plan, CHANGELOG)
- [ ] PR created with comprehensive description
- [ ] Code review completed
- [ ] Merged to main
- [ ] Tagged as v0.0.14
- [ ] Release notes published

### Verification
Production validation with user's reported case:
1. User's ASN resource with tenant field
2. Remove tenant from config
3. Apply succeeds without "inconsistent result" error
4. Tenant is null in state and API
