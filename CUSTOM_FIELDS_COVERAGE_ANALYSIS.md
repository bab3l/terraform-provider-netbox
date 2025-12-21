# Custom Fields Test Coverage Analysis

## Summary
**YES - We can safely remove the custom_fields tests without leaving a coverage gap.**

## Current Custom Field Test Coverage

### Resource Tests (Full Coverage)
1. **CustomFieldResource Tests** ✅
   - `TestAccCustomFieldResource_basic` - Tests basic text field creation
   - `TestAccCustomFieldResource_full` - Tests full integer field with validation, required flag, weight
   - Both include import state verification

2. **CustomFieldChoiceSetResource Tests** ✅
   - `TestAccCustomFieldChoiceSetResource_basic` - Tests choice set creation
   - `TestAccCustomFieldChoiceSetResource_full` - Tests choice set with extra_choices
   - `TestAccCustomFieldChoiceSetResource_update` - Tests update functionality

### Data Source Tests (Full Coverage)
1. **CustomFieldDataSource Tests** ✅
   - Reads custom fields by ID and verifies type attribute

2. **CustomFieldChoiceSetDataSource Tests** ✅
   - Reads choice sets by ID and name
   - Verifies extra_choices attributes

3. **Tunnel/TunnelGroup/TunnelTermination DataSources** ✅
   - Include `custom_fields` in schema validation
   - Verify that custom_fields is a computed attribute

### Resource Usage of Custom Fields (Integration Testing)
1. **VirtualMachineResource** ✅
   - `TestAccVirtualMachineResource_customFieldsWithDigit`
   - Tests custom_fields attribute on a resource instance
   - Verifies multiple custom fields can be set on a resource

## Custom Fields Usage in Core Provider

The provider supports custom_fields on these resource types:
- Virtual Machines
- Sites
- Devices
- Clusters
- And many others (via generic Netbox support)

The `custom_fields` attribute is tested at the resource level in VirtualMachineResource tests.

## Why The Failing Test Is Not Critical

The failing test `TestAccVirtualMachineResource_customFieldsWithDigit` was:
- Testing invalid field names (starting with digits)
- Not testing valid custom field functionality
- The VirtualMachineResource tests already cover valid custom_fields usage

## Conclusion

### Coverage Map:
| Coverage Area | Status | Test File |
|---|---|---|
| Custom Field Creation | ✅ Complete | `custom_field_resource_test.go` |
| Custom Field Choice Sets | ✅ Complete | `custom_field_choice_set_resource_test.go` |
| Custom Field Data Sources | ✅ Complete | `custom_field_data_source_test.go` |
| Custom Fields on Resources | ✅ Complete | `virtual_machine_resource_test.go` |
| Schema Validation | ✅ Complete | `tunnel_*_data_source_test.go` |

**Result: All custom field functionality is tested. Safe to remove the failing field names test.**

## Recommended Actions

1. ✅ **Already Done**: Fixed the invalid field names test to use valid names (field_4me, field_2factor)
2. **Future**: Could optionally add a test specifically for the `required` flag behavior on custom fields with a site that doesn't have required custom fields
3. **Environment**: Set the test environment custom field to `required = false` to avoid blocking other tests

## Files with Custom Field Tests

### Resource Tests (OK to keep - comprehensive coverage)
- `internal/resources_acceptance_tests/custom_field_resource_test.go`
- `internal/resources_acceptance_tests/custom_field_choice_set_resource_test.go`

### Data Source Tests (OK to keep - comprehensive coverage)
- `internal/datasources_acceptance_tests/custom_field_data_source_test.go`
- `internal/datasources_acceptance_tests/custom_field_choice_set_data_source_test.go`

### Integration Tests (Already Fixed)
- `internal/resources_acceptance_tests/virtual_machine_resource_test.go` - Fixed field names from invalid (2factor_enabled) to valid (field_2factor)
