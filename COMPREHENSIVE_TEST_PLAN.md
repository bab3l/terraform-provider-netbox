# Comprehensive Test Coverage Plan for Optional Fields

## Overview
With all 20 optional field bugs fixed, we need comprehensive test coverage to ensure edge cases are handled properly and prevent regressions.

## Test Pattern Classifications

### Pattern A: Optional+Computed Fields (Always Set)
**Behavior**: Field always present in state with either user value or computed default
**Resources**: device.status, vm.status, vlan.status, role.weight, template fields

**Test Scenarios**:
1. **Field Not In Config** - Verify default value is set
2. **Field Removed From Config** - Verify default value maintained (no drift)
3. **Import Without Config** - Verify API value preserved
4. **Plan-Only After Import** - No changes proposed

### Pattern B: Optional Only Fields (Conditional)
**Behavior**: Present only when user specifies OR during import
**Resources**: vm_interface.mode, interface.enabled

**Test Scenarios**:
1. **Field Not In Config** - Verify field is null/absent
2. **Field Removed From Config** - Verify field becomes null
3. **Import With Field Set** - Preserve if not in config, conditional behavior
4. **Plan-Only Consistency** - No spurious changes

## Test Implementation Strategy

### Phase 1: High-Priority Resources (Batch 1)
Focus on most critical resources that could cause user impact:
- `device_resource.go` (status field)
- `virtual_machine_resource.go` (status field)
- `vlan_resource.go` (status field)
- `vm_interface_resource.go` (mode field)

### Phase 2: Network & Interface Resources (Batch 2)
- `interface_resource.go` (enabled field)
- `tunnel_resource.go` (status field)
- `interface_template_resource.go` (enabled field)

### Phase 3: Power & Template Resources (Batch 3)
- `power_feed_resource.go` (voltage, amperage fields)
- `power_outlet_template_resource.go` (label field)
- `power_port_template_resource.go` (label field)
- `front_port_template_resource.go` (label, color fields)
- `rear_port_template_resource.go` (label, color, positions fields)

### Phase 4: Miscellaneous Resources (Batch 4)
- `role_resource.go` (weight field)
- `rear_port_resource.go` (positions field)
- `journal_entry_resource.go` (kind field)

## Test Template Framework

### Template 1: Optional+Computed Field Tests
```go
func TestAcc{Resource}_OptionalComputedField_{FieldName}(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        Providers: testAccProviders,
        CheckDestroy: testAcc{Resource}Destroy,
        Steps: []resource.TestStep{
            // Step 1: Create without optional field
            {
                Config: testAcc{Resource}Config_withoutOptionalField(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "{field_name}", "{default_value}"),
                ),
            },
            // Step 2: Add optional field
            {
                Config: testAcc{Resource}Config_withOptionalField("{custom_value}"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "{field_name}", "{custom_value}"),
                ),
            },
            // Step 3: Remove optional field - should revert to default
            {
                Config: testAcc{Resource}Config_withoutOptionalField(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "{field_name}", "{default_value}"),
                ),
            },
        },
    })
}
```

### Template 2: Optional Only Field Tests
```go
func TestAcc{Resource}_OptionalOnlyField_{FieldName}(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        Providers: testAccProviders,
        CheckDestroy: testAcc{Resource}Destroy,
        Steps: []resource.TestStep{
            // Step 1: Create without optional field
            {
                Config: testAcc{Resource}Config_withoutOptionalField(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckNoResourceAttr(resourceName, "{field_name}"),
                ),
            },
            // Step 2: Add optional field
            {
                Config: testAcc{Resource}Config_withOptionalField("{custom_value}"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "{field_name}", "{custom_value}"),
                ),
            },
            // Step 3: Remove optional field - should become null
            {
                Config: testAcc{Resource}Config_withoutOptionalField(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckNoResourceAttr(resourceName, "{field_name}"),
                ),
            },
        },
    })
}
```

### Template 3: Import Behavior Tests
```go
func TestAcc{Resource}_ImportOptionalField_{FieldName}(t *testing.T) {
    resourceName := "netbox_{resource}.test"

    resource.Test(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        Providers: testAccProviders,
        CheckDestroy: testAcc{Resource}Destroy,
        Steps: []resource.TestStep{
            // Step 1: Create with field set
            {
                Config: testAcc{Resource}Config_withOptionalField("{custom_value}"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "{field_name}", "{custom_value}"),
                ),
            },
            // Step 2: Import resource
            {
                ResourceName:      resourceName,
                ImportState:       true,
                ImportStateVerify: true,
            },
            // Step 3: Plan-only to verify no drift
            {
                Config:   testAcc{Resource}Config_withOptionalField("{custom_value}"),
                PlanOnly: true,
            },
        },
    })
}
```

## Implementation Batches

### Batch 1: High-Priority Optional+Computed (4 tests)
- device_status_optional_computed_test.go
- virtual_machine_status_optional_computed_test.go
- vlan_status_optional_computed_test.go
- vm_interface_mode_optional_only_test.go (already exists, verify comprehensive)

### Batch 2: Network Resources (3 tests)
- interface_enabled_optional_computed_test.go (extend existing)
- tunnel_status_optional_computed_test.go
- interface_template_enabled_optional_computed_test.go

### Batch 3: Power & Template Resources (8 tests)
- power_feed_voltage_amperage_optional_computed_test.go
- power_outlet_template_label_optional_computed_test.go
- power_port_template_label_optional_computed_test.go
- front_port_template_label_color_optional_computed_test.go
- rear_port_template_label_color_positions_optional_computed_test.go

### Batch 4: Miscellaneous Resources (3 tests)
- role_weight_optional_computed_test.go
- rear_port_positions_optional_computed_test.go
- journal_entry_kind_optional_computed_test.go

## Quality Gates

### For Each Test:
- [ ] Covers field not in config scenario
- [ ] Covers field removed from config scenario
- [ ] Covers import behavior
- [ ] Verifies no unwanted drift
- [ ] Tests both default and custom values
- [ ] Uses proper TestCheckResourceAttr/TestCheckNoResourceAttr

### For Each Batch:
- [ ] All tests pass individually
- [ ] All tests pass when run together
- [ ] No interference between tests
- [ ] Proper cleanup after each test

## Success Metrics

- **100% Edge Case Coverage**: All identified edge cases have dedicated tests
- **Zero Flaky Tests**: All tests pass consistently
- **No Test Interference**: Tests can run in any order
- **Clear Failure Messages**: When tests fail, reason is immediately apparent
- **Performance**: Test suite completes in reasonable time (<30min total)

## Risk Mitigation

### High-Risk Scenarios:
1. **State Corruption**: Test field removal doesn't corrupt other fields
2. **Import Inconsistencies**: Imported resources behave identically to created ones
3. **Terraform Plan Drift**: No spurious changes after successful apply
4. **Default Value Conflicts**: Ensure defaults match between Create and Read operations

### Testing Strategy:
- Run tests individually first to isolate issues
- Run full batch to test interactions
- Include plan-only steps to catch drift
- Test with various field combinations

## Implementation Status

### ✅ Phase 1: Planning and Framework (COMPLETE)
- [x] Comprehensive test coverage plan developed
- [x] Test pattern templates created (Optional+Computed vs Optional Only)
- [x] Batch implementation strategy defined
- [x] Quality gates established

### 🔄 Phase 2: Batch 1 Implementation (IN PROGRESS)
- [x] Created device_status_optional_computed_test.go (needs testutil pattern fix)
- [x] Created virtual_machine_status_optional_computed_test.go (needs testutil pattern fix)
- [x] Created vlan_status_optional_computed_test.go (needs testutil pattern fix)
- [x] Created vm_interface_mode_optional_only_test.go (needs testutil pattern fix)
- [ ] Fix all test files to use correct testutil.OptionalFieldTestConfig pattern
- [ ] Validate Batch 1 tests pass individually
- [ ] Validate Batch 1 tests pass as group

### 🎯 Current Status: Test Framework Integration
**Issue**: Created comprehensive tests use manual pattern instead of existing `testutil.OptionalFieldTestConfig` framework
**Solution**: Refactor tests to use the established testutil pattern for consistency

### 📋 Discovered: Existing Test Framework
The codebase already has `testutil.OptionalFieldTestConfig` which provides:
- Automated test step generation
- Consistent patterns for Optional+Computed vs Optional Only fields
- Pre-built resource configurations
- Standard verification steps

### 🛠 Immediate Next Steps
1. **Refactor Batch 1 tests** to use testutil.OptionalFieldTestConfig
2. **Test individual resources** with new comprehensive coverage
3. **Validate no regressions** with existing basic tests
4. **Proceed with remaining batches** using correct pattern

### 📊 Overall Progress
- **Core Bugs Fixed**: 20/20 (100% COMPLETE)
- **Basic Test Coverage**: ✅ All resources validated
- **Comprehensive Test Coverage**: 4/20 created (needs framework integration)
- **Production Ready**: ✅ All fixes working in basic scenarios

---
**Status**: Framework integration needed before proceeding with comprehensive testing
**Next Action**: Use testutil.OptionalFieldTestConfig for all comprehensive tests
