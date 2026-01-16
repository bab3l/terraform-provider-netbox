# Acceptance Test Coverage Analysis

This document tracks the current state of acceptance test coverage for all resources and identifies gaps that need to be addressed.

## Coverage Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Test exists and follows naming convention |
| ⚠️ | Test exists but may need review/renaming |
| ❌ | Test missing |
| N/A | Not applicable to this resource |

---

## Resource Completion Checklist

Use this checklist to verify a resource's tests are complete before moving to the next resource.

### Gating Criteria (All Must Pass)

```markdown
## Resource: {resource_name}

### TIER 1: Core CRUD Tests
- [ ] `TestAcc{Resource}Resource_basic` exists and passes
- [ ] `TestAcc{Resource}Resource_full` exists and passes
- [ ] `TestAcc{Resource}Resource_update` exists and passes
- [ ] `TestAcc{Resource}Resource_import` exists and passes

### TIER 2: Reliability Tests
- [ ] `TestAcc{Resource}Resource_IDPreservation` exists and passes
- [ ] `TestAcc{Resource}Resource_externalDeletion` exists and passes (uses `RunExternalDeletionTest` helper)
- [ ] `TestAcc{Resource}Resource_removeOptionalFields` exists and passes (uses `TestRemoveOptionalFields` helper)

### TIER 3: Tag Tests (if resource supports tags)
- [ ] `TestAcc{Resource}Resource_tagLifecycle` exists and passes (uses `RunTagLifecycleTest` helper)
- [ ] `TestAcc{Resource}Resource_tagOrderInvariance` exists and passes (uses `RunTagOrderTest` helper)

### TIER 4: Quality Checks
- [ ] `TestAcc{Resource}Resource_validationErrors` exists (uses `RunMultiValidationErrorTest` helper)
- [ ] All test names follow camelCase convention (no underscores except after `Resource`)
- [ ] All config functions follow `testAcc{Resource}ResourceConfig_{variant}` naming
- [ ] Cleanup registration exists for all created resources
- [ ] All tests call `t.Parallel()`

### Helper Function Usage
- [ ] External deletion test uses `RunExternalDeletionTest()` or `RunExternalDeletionWithIDTest()`
- [ ] Tag tests use `RunTagLifecycleTest()` and `RunTagOrderTest()`
- [ ] Validation tests use `RunMultiValidationErrorTest()`
- [ ] Optional field tests use `TestRemoveOptionalFields()`
- [ ] Import tests use `RunImportTest()` or `RunSimpleImportTest()` (where applicable)

### Final Verification
- [ ] All tests pass locally: `go test -v -run "TestAcc{Resource}Resource" ./internal/resources_acceptance_tests/...`
- [ ] No naming convention violations
- [ ] Code formatted with `gofmt`
```

### Quick Validation Command

Run this to verify a resource's tests:

```bash
# Set environment
export TF_ACC=1
export NETBOX_SERVER_URL="http://localhost:8000"
export NETBOX_API_TOKEN="your-token"

# Run all tests for a specific resource (e.g., ip_address)
go test -v -run "TestAccIPAddressResource" ./internal/resources_acceptance_tests/...

# Run specific test type
go test -v -run "TestAccIPAddressResource_tag" ./internal/resources_acceptance_tests/...
```

---

## Test Coverage Matrix

### TIER 1: Core CRUD Tests

| Resource | `_basic` | `_full` | `_update` | `_import` |
|----------|----------|---------|-----------|-----------|
| aggregate | ✅ | ✅ | ✅ | ✅ |
| asn | ✅ | ✅ | ✅ | ✅ |
| asn_range | ✅ | ✅ | ✅ | ✅ |
| cable | ✅ | ✅ | ✅ | ⚠️ |
| circuit | ✅ | ✅ | ✅ | ✅ |
| circuit_group | ✅ | ✅ | ✅ | ✅ |
| circuit_group_assignment | ✅ | ✅ | ✅ | ❌ |
| circuit_termination | ✅ | ✅ | ✅ | ❌ |
| circuit_type | ✅ | ✅ | ✅ | ✅ |
| cluster | ✅ | ✅ | ✅ | ✅ |
| cluster_group | ✅ | ✅ | ✅ | ✅ |
| cluster_type | ✅ | ✅ | ✅ | ✅ |
| config_context | ✅ | ✅ | ✅ | ✅ |
| config_template | ✅ | ✅ | ✅ | ✅ |
| console_port | ✅ | ✅ | ✅ | ❌ |
| console_port_template | ✅ | ✅ | ✅ | ❌ |
| console_server_port | ✅ | ✅ | ✅ | ❌ |
| console_server_port_template | ✅ | ✅ | ✅ | ❌ |
| contact | ✅ | ✅ | ✅ | ✅ |
| contact_assignment | ✅ | ✅ | ✅ | ❌ |
| contact_group | ✅ | ✅ | ✅ | ✅ |
| contact_role | ✅ | ✅ | ✅ | ✅ |
| custom_link | ✅ | ✅ | ✅ | ❌ |
| device | ✅ | ✅ | ✅ | ✅ |
| device_bay | ✅ | ✅ | ✅ | ❌ |
| device_bay_template | ✅ | ✅ | ✅ | ❌ |
| device_role | ✅ | ✅ | ✅ | ✅ |
| device_type | ✅ | ✅ | ✅ | ✅ |
| event_rule | ✅ | ✅ | ✅ | ❌ |
| export_template | ✅ | ✅ | ✅ | ❌ |
| fhrp_group | ✅ | ✅ | ✅ | ❌ |
| fhrp_group_assignment | ✅ | ✅ | ✅ | ❌ |
| front_port | ✅ | ✅ | ✅ | ❌ |
| front_port_template | ✅ | ✅ | ✅ | ❌ |
| ike_policy | ✅ | ✅ | ✅ | ❌ |
| ike_proposal | ✅ | ✅ | ✅ | ❌ |
| interface | ✅ | ✅ | ✅ | ✅ |
| interface_template | ✅ | ✅ | ✅ | ❌ |
| inventory_item | ✅ | ✅ | ✅ | ❌ |
| inventory_item_role | ✅ | ✅ | ✅ | ❌ |
| inventory_item_template | ✅ | ✅ | ✅ |
| ip_address | ✅ | ✅ | ✅ | ✅ | **COMPLETE ✅** |
| ip_range | ✅ | ✅ | ✅ | ❌ |
| ipsec_policy | ✅ | ✅ | ✅ | ❌ |
| ipsec_profile | ✅ | ✅ | ✅ | ❌ |
| ipsec_proposal | ✅ | ✅ | ✅ | ❌ |
| journal_entry | ✅ | ✅ | ✅ | ❌ |
| l2vpn | ✅ | ✅ | ✅ | ✅ |
| l2vpn_termination | ✅ | ✅ | ✅ | ❌ |
| location | ✅ | ✅ | ✅ | ✅ |
| manufacturer | ✅ | ✅ | ✅ | ✅ |
| module | ✅ | ✅ | ✅ | ❌ |
| module_bay | ✅ | ✅ | ✅ | ❌ |
| module_bay_template | ✅ | ✅ | ✅ | ❌ |
| module_type | ✅ | ✅ | ✅ | ❌ |
| notification_group | ✅ | ✅ | ✅ | ✅ |
| platform | ✅ | ✅ | ✅ | ✅ |
| power_feed | ✅ | ✅ | ✅ | ❌ |
| power_outlet | ✅ | ✅ | ✅ | ❌ |
| power_outlet_template | ✅ | ✅ | ✅ | ❌ |
| power_panel | ✅ | ✅ | ✅ | ❌ |
| power_port | ✅ | ✅ | ✅ | ❌ |
| power_port_template | ✅ | ✅ | ✅ | ❌ |
| prefix | ✅ | ✅ | ✅ | ✅ | **COMPLETE ✅** |
| provider | ✅ | ✅ | ✅ | ✅ |
| provider_account | ✅ | ✅ | ✅ | ✅ |
| provider_network | ✅ | ✅ | ✅ | ✅ |
| rack | ✅ | ✅ | ✅ | ✅ |
| rack_reservation | ✅ | ✅ | ✅ | ❌ |
| rack_role | ✅ | ✅ | ✅ | ✅ |
| rack_type | ✅ | ✅ | ✅ | ❌ |
| rear_port | ✅ | ✅ | ✅ | ❌ |
| rear_port_template | ✅ | ✅ | ✅ | ❌ |
| region | ✅ | ✅ | ✅ | ✅ |
| rir | ✅ | ✅ | ✅ | ❌ |
| role | ✅ | ✅ | ✅ | ❌ |
| route_target | ✅ | ✅ | ✅ | ✅ |
| service | ✅ | ✅ | ✅ | ❌ |
| service_template | ✅ | ✅ | ✅ | ❌ |
| site | ✅ | ✅ | ✅ | ✅ |
| site_group | ✅ | ✅ | ✅ | ✅ |
| tag | ✅ | ✅ | ✅ | ❌ |
| tenant | ✅ | ✅ | ✅ | ✅ |
| tenant_group | ✅ | ✅ | ✅ | ✅ |
| tunnel | ✅ | ✅ | ✅ | ✅ |
| tunnel_group | ✅ | ✅ | ✅ | ✅ |
| tunnel_termination | ✅ | ✅ | ✅ | ✅ |
| virtual_chassis | ✅ | ✅ | ✅ | ❌ |
| virtual_device_context | ✅ | ✅ | ✅ | ❌ |
| virtual_disk | ✅ | ✅ | ✅ | ✅ |
| virtual_machine | ✅ | ✅ | ✅ | ❌ |
| vlan | ✅ | ✅ | ✅ | ✅ |
| vlan_group | ✅ | ✅ | ✅ | ✅ |
| vm_interface | ✅ | ✅ | ✅ | ✅ |
| vrf | ✅ | ✅ | ✅ | ✅ |
| webhook | ✅ | ✅ | ✅ | ❌ |
| wireless_lan | ✅ | ✅ | ✅ | ❌ |
| wireless_lan_group | ✅ | ✅ | ✅ | ❌ |
| wireless_link | ✅ | ✅ | ✅ | ❌ |

### TIER 2: Reliability Tests

| Resource | `_IDPreservation` | `_externalDeletion` | `_removeOptionalFields` |
|----------|-------------------|---------------------|-------------------------|
| aggregate | ✅ | ✅ | ✅ |
| asn | ✅ | ✅ | ✅ |
| asn_range | ✅ | ✅ | ✅ |
| cable | ✅ | ✅ | ✅ |
| circuit | ✅ | ✅ | ✅ |
| circuit_group | ✅ | ✅ | ✅ |
| circuit_group_assignment | ✅ | ✅ | ❌ |
| circuit_termination | ✅ | ✅ | ✅ |
| circuit_type | ✅ | ✅ | ✅ |
| cluster | ✅ | ✅ | ✅ |
| cluster_group | ✅ | ✅ | ✅ |
| cluster_type | ✅ | ✅ | ✅ |
| config_context | ✅ | ✅ | ✅ |
| config_template | ✅ | ✅ | ✅ |
| console_port | ✅ | ✅ | ✅ |
| console_port_template | ✅ | ✅ | ✅ |
| console_server_port | ✅ | ✅ | ✅ |
| console_server_port_template | ✅ | ✅ | ✅ |
| contact | ✅ | ✅ | ✅ |
| contact_assignment | ✅ | ✅ | ❌ |
| contact_group | ✅ | ✅ | ✅ |
| contact_role | ✅ | ✅ | ✅ |
| custom_link | ✅ | ✅ | ✅ |
| device | ✅ | ✅ | ✅ |
| device_bay | ✅ | ✅ | ✅ |
| device_bay_template | ✅ | ✅ | ✅ |
| device_role | ✅ | ✅ | ✅ |
| device_type | ✅ | ✅ | ✅ |
| event_rule | ✅ | ✅ | ✅ |
| export_template | ✅ | ✅ | ✅ |
| fhrp_group | ✅ | ✅ | ✅ |
| fhrp_group_assignment | ✅ | ✅ | ❌ |
| front_port | ✅ | ✅ | ✅ |
| front_port_template | ✅ | ✅ | ✅ |
| ike_policy | ✅ | ✅ | ✅ |
| ike_proposal | ✅ | ✅ | ✅ |
| interface | ✅ | ✅ | ✅ |
| interface_template | ✅ | ✅ | ✅ |
| inventory_item | ✅ | ✅ | ✅ |
| inventory_item_role | ✅ | ✅ | ✅ |
| inventory_item_template | ✅ | ✅ | ✅ |
| ip_address | ✅ | ✅ | ✅ | **COMPLETE ✅** |
| ip_range | ✅ | ✅ | ✅ |
| ipsec_policy | ✅ | ✅ | ✅ |
| ipsec_profile | ✅ | ✅ | ✅ |
| ipsec_proposal | ✅ | ✅ | ✅ |
| journal_entry | ✅ | ✅ | ❌ |
| l2vpn | ✅ | ✅ | ✅ |
| l2vpn_termination | ✅ | ⚠️ | ✅ |
| location | ✅ | ✅ | ✅ |
| manufacturer | ✅ | ✅ | ✅ |
| module | ✅ | ⚠️ | ✅ |
| module_bay | ✅ | ⚠️ | ✅ |
| module_bay_template | ✅ | ⚠️ | ✅ |
| module_type | ✅ | ⚠️ | ✅ |
| notification_group | ✅ | ✅ | ✅ |
| platform | ✅ | ✅ | ✅ |
| power_feed | ✅ | ✅ | ✅ |
| power_outlet | ✅ | ✅ | ✅ |
| power_outlet_template | ✅ | ✅ | ✅ |
| power_panel | ✅ | ✅ | ✅ |
| power_port | ✅ | ✅ | ✅ |
| power_port_template | ✅ | ✅ | ✅ |
| prefix | ✅ | ✅ | ✅ | **COMPLETE ✅** |
| provider | ✅ | ✅ | ✅ |
| provider_account | ✅ | ✅ | ✅ |
| provider_network | ✅ | ✅ | ✅ |
| rack | ✅ | ✅ | ⚠️ |
| rack_reservation | ✅ | ✅ | ✅ |
| rack_role | ✅ | ✅ | ✅ |
| rack_type | ✅ | ✅ | ✅ |
| rear_port | ✅ | ✅ | ✅ |
| rear_port_template | ✅ | ✅ | ✅ |
| region | ✅ | ✅ | ✅ |
| rir | ✅ | ⚠️ | ✅ |
| role | ✅ | ✅ | ✅ |
| route_target | ✅ | ✅ | ✅ |
| service | ✅ | ⚠️ | ✅ |
| service_template | ✅ | ⚠️ | ✅ |
| site | ✅ | ✅ | ✅ |
| site_group | ✅ | ✅ | ✅ |
| tag | ✅ | ✅ | ✅ |
| tenant | ✅ | ✅ | ✅ |
| tenant_group | ✅ | ✅ | ✅ |
| tunnel | ✅ | ✅ | ✅ |
| tunnel_group | ✅ | ✅ | ✅ |
| tunnel_termination | ✅ | ✅ | ✅ |
| virtual_chassis | ✅ | ✅ | ✅ |
| virtual_device_context | ✅ | ✅ | ✅ |
| virtual_disk | ✅ | ✅ | ✅ |
| virtual_machine | ✅ | ✅ | ✅ |
| vlan | ✅ | ✅ | ✅ |
| vlan_group | ✅ | ✅ | ✅ |
| vm_interface | ✅ | ⚠️ | ✅ |
| vrf | ✅ | ✅ | ✅ |
| webhook | ✅ | ✅ | ✅ |
| wireless_lan | ✅ | ✅ | ✅ |
| wireless_lan_group | ✅ | ✅ | ✅ |
| wireless_link | ✅ | ✅ | ✅ |

### TIER 3: Tag Tests (Resources with Tags)

| Resource | `_tagLifecycle` | `_tagOrderInvariance` |
|----------|-----------------|----------------------|
| aggregate | ❌ | ❌ |
| asn | ❌ | ❌ |
| asn_range | ❌ | ❌ |
| cable | ❌ | ❌ |
| circuit | ❌ | ❌ |
| circuit_termination | ❌ | ❌ |
| circuit_type | ❌ | ❌ |
| cluster | ❌ | ❌ |
| cluster_group | ❌ | ❌ |
| cluster_type | ❌ | ❌ |
| config_context | ❌ | ❌ |
| config_template | ❌ | ❌ |
| console_port | ❌ | ❌ |
| console_port_template | ❌ | ❌ |
| console_server_port | ❌ | ❌ |
| console_server_port_template | ❌ | ❌ |
| contact | ❌ | ❌ |
| contact_group | ❌ | ❌ |
| contact_role | ❌ | ❌ |
| device | ❌ | ❌ |
| device_bay | ❌ | ❌ |
| device_bay_template | ❌ | ❌ |
| device_role | ❌ | ❌ |
| device_type | ❌ | ❌ |
| fhrp_group | ❌ | ❌ |
| front_port | ❌ | ❌ |
| front_port_template | ❌ | ❌ |
| interface | ❌ | ❌ |
| interface_template | ❌ | ❌ |
| inventory_item | ❌ | ❌ |
| inventory_item_role | ❌ | ❌ |
| inventory_item_template | ❌ | ❌ |
| ip_address | ⚠️ | ✅ | **COMPLETE ✅** (See IP_ADDRESS_CHECKLIST.md) |
| ip_range | ❌ | ❌ |
| l2vpn | ❌ | ❌ |
| l2vpn_termination | ❌ | ❌ |
| location | ❌ | ❌ |
| manufacturer | ❌ | ❌ |
| module | ❌ | ❌ |
| module_bay | ❌ | ❌ |
| module_bay_template | ❌ | ❌ |
| module_type | ❌ | ❌ |
| platform | ❌ | ❌ |
| power_feed | ❌ | ❌ |
| power_outlet | ❌ | ❌ |
| power_outlet_template | ❌ | ❌ |
| power_panel | ❌ | ❌ |
| power_port | ❌ | ❌ |
| power_port_template | ❌ | ❌ |
| prefix | ✅ | ✅ | **COMPLETE ✅** (See PREFIX_CHECKLIST.md) |
| provider | ❌ | ❌ |
| provider_account | ❌ | ❌ |
| provider_network | ❌ | ❌ |
| rack | ❌ | ❌ |
| rack_reservation | ❌ | ❌ |
| rack_role | ❌ | ❌ |
| rack_type | ❌ | ❌ |
| rear_port | ❌ | ❌ |
| rear_port_template | ❌ | ❌ |
| region | ❌ | ❌ |
| rir | ❌ | ❌ |
| role | ❌ | ❌ |
| route_target | ❌ | ❌ |
| service | ❌ | ❌ |
| service_template | ❌ | ❌ |
| site | ❌ | ❌ |
| site_group | ❌ | ❌ |
| tenant | ❌ | ❌ |
| tenant_group | ❌ | ❌ |
| tunnel | ❌ | ❌ |
| tunnel_group | ❌ | ❌ |
| tunnel_termination | ❌ | ❌ |
| virtual_chassis | ❌ | ❌ |
| virtual_device_context | ❌ | ❌ |
| virtual_disk | ❌ | ❌ |
| virtual_machine | ❌ | ❌ |
| vlan | ❌ | ❌ |
| vlan_group | ❌ | ❌ |
| vm_interface | ❌ | ❌ |
| vrf | ❌ | ❌ |
| webhook | ❌ | ❌ |
| wireless_lan | ❌ | ❌ |
| wireless_lan_group | ❌ | ❌ |
| wireless_link | ❌ | ❌ |

---

## Naming Convention Issues

The following tests have naming inconsistencies that should be addressed:

### External Deletion Tests (should be `_externalDeletion`)
- `l2vpn_termination`: `_external_deletion` → `_externalDeletion`
- `module`: `_external_deletion` → `_externalDeletion`
- `module_bay`: `_external_deletion` → `_externalDeletion`
- `module_bay_template`: `_external_deletion` → `_externalDeletion`
- `module_type`: `_external_deletion` → `_externalDeletion`
- `rir`: `_external_deletion` → `_externalDeletion`
- `service`: `_external_deletion` → `_externalDeletion`
- `service_template`: `_external_deletion` → `_externalDeletion`
- `vm_interface`: `_external_deletion` → `_externalDeletion`

---

## Priority Work Items

### Phase 1: Fix Tag Removal Bug (COMPLETED)
- [x] Fix `TagsToNestedTagRequests` to return empty slice instead of nil
- [x] Add `TestAccIPAddressResource_tagRemoval` test
- [x] Add `TestAccIPAddressResource_tagOrderInvariance` test

### Phase 2: Standardize IP Address Tests as Reference Implementation
- [x] `_tagRemoval` (add → remove → verify)
- [x] `_createWithTags` (create with tags)
- [x] `_modifyTags` (change tags)
- [x] `_tagOrderInvariance` (order doesn't matter)
- [ ] Rename to `_tagLifecycle` (consolidate into single comprehensive test)
- [ ] Refactor to use `RunTagLifecycleTest` and `RunTagOrderTest` helpers

### Phase 3: Migrate Existing Tests to Use Helper Functions
Audit and refactor existing tests to use standardized helpers from `internal/testutil/`:

#### Available Test Helpers

| Helper Function | Purpose | File |
|-----------------|---------|------|
| `RunTagLifecycleTest()` | Complete tag add/modify/remove cycle | `tag_tests.go` |
| `RunTagOrderTest()` | Tag reordering doesn't cause drift | `tag_tests.go` |
| `RunExternalDeletionTest()` | Handle resource deleted outside TF | `external_deletion_tests.go` |
| `RunExternalDeletionWithIDTest()` | External deletion with ID tracking | `external_deletion_tests.go` |
| `RunImportTest()` | Standard import testing | `import_tests.go` |
| `RunSimpleImportTest()` | Simplified import testing | `import_tests.go` |
| `RunUpdateTest()` | Single-step update testing | `update_tests.go` |
| `RunMultiStepUpdateTest()` | Multi-step update scenarios | `update_tests.go` |
| `RunFieldUpdateTest()` | Individual field update testing | `update_tests.go` |
| `TestRemoveOptionalFields()` | Remove optional fields from config | `optional_field_tests.go` |
| `RunOptionalFieldTestSuite()` | Comprehensive optional field tests | `optional_field_tests.go` |
| `RunOptionalComputedFieldTestSuite()` | Optional computed field tests | `optional_computed_field_tests.go` |
| `RunValidationErrorTest()` | Single validation error test | `validation_tests.go` |
| `RunMultiValidationErrorTest()` | Multiple validation error cases | `validation_tests.go` |
| `RunReferenceChangeTest()` | Reference field change testing | `reference_tests.go` |
| `RunMultiReferenceTest()` | Multi-reference field testing | `reference_tests.go` |
| `RunIdempotencyTest()` | Idempotency verification | `idempotency_tests.go` |
| `RunRefreshIdempotencyTest()` | Refresh idempotency testing | `idempotency_tests.go` |
| `RunHierarchicalTest()` | Parent-child relationship tests | `hierarchical_tests.go` |
| `RunNestedHierarchyTest()` | Multi-level hierarchy tests | `hierarchical_tests.go` |

#### Migration Priority
1. **Tag tests** - Use `RunTagLifecycleTest()` and `RunTagOrderTest()`
2. **External deletion tests** - Use `RunExternalDeletionTest()`
3. **Import tests** - Use `RunImportTest()` or `RunSimpleImportTest()`
4. **Validation tests** - Use `RunMultiValidationErrorTest()`
5. **Update tests** - Use `RunUpdateTest()` where applicable
6. **Optional field tests** - Use `TestRemoveOptionalFields()`

### Phase 4: Apply Tag Tests to All Resources with Tags
Priority order (most commonly used resources first):
1. `virtual_machine`
2. `device`
3. `prefix`
4. `vlan`
5. `site`
6. `interface`
7. `cluster`
8. ... (remaining resources)

**Implementation approach:**
- Use `RunTagLifecycleTest()` helper for all tag lifecycle tests
- Use `RunTagOrderTest()` helper for tag order invariance tests
- Create reusable config generator patterns for each resource

### Phase 5: Fix Naming Inconsistencies
- Rename `_external_deletion` tests to `_externalDeletion`
- Standardize config function naming

### Phase 6: Add Missing Import Tests
Resources needing `_import` tests (42 resources):
- Use `RunImportTest()` or `RunSimpleImportTest()` helpers
- circuit_group_assignment
- circuit_termination
- console_port
- console_port_template
- console_server_port
- console_server_port_template
- contact_assignment
- custom_link
- device_bay
- device_bay_template
- event_rule
- export_template
- fhrp_group
- fhrp_group_assignment
- front_port
- front_port_template
- ike_policy
- ike_proposal
- interface_template
- inventory_item
- inventory_item_role
- inventory_item_template
- ip_range
- ipsec_policy
- ipsec_profile
- ipsec_proposal
- journal_entry
- l2vpn_termination
- module
- module_bay
- module_bay_template
- module_type
- power_feed
- power_outlet
- power_outlet_template
- power_panel
- power_port
- power_port_template
- rack_reservation
- rack_type
- rear_port
- rear_port_template
- rir
- role
- service
- service_template
- tag
- virtual_chassis
- virtual_device_context
- virtual_machine
- webhook
- wireless_lan
- wireless_lan_group
- wireless_link

---

## Helper Function Usage Tracking

Track which resources are using standardized helper functions vs custom implementations.

### Legend
- ✅ Uses helper function
- ❌ Custom implementation (needs migration)
- N/A Not applicable

| Resource | Tag Helpers | External Del. | Import | Validation | Optional Fields |
|----------|-------------|---------------|--------|------------|-----------------|
| ip_address | ❌ | ❌ | ❌ | ✅ | ✅ |
| *other resources need audit* | | | | | |

**TODO:** Complete audit of all 86 resources to track helper function usage.

---

## Statistics

### Current Coverage Summary

| Test Category | Implemented | Total Required | Coverage % |
|---------------|-------------|----------------|------------|
| `_basic` | 86 | 86 | 100% |
| `_full` | 86 | 86 | 100% |
| `_update` | 86 | 86 | 100% |
| `_import` | 44 | 86 | 51% |
| `_IDPreservation` | 86 | 86 | 100% |
| `_externalDeletion` | 86 | 86 | 100% |
| `_removeOptionalFields` | 82 | 86 | 95% |
| `_tagLifecycle` | 1 | 73 | 1% |
| `_tagOrderInvariance` | 1 | 73 | 1% |

### Tag Test Gap
- **73 resources** support tags
- **Only 1 resource** (ip_address) has tag lifecycle tests
- **72 resources** need tag tests added

---

*Last updated: January 2026*
