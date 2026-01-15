# Full Test Implementation Plan

## Overview

Implement comprehensive full configuration tests for the 6 remaining resources to achieve 100% full test coverage. Full tests verify that resources work correctly with all optional fields populated, ensuring complete feature coverage and preventing regressions.

**Current Status:** 93/99 complete (93.9%) - 6 resources remaining

## Test Pattern

```go
func TestAcc{Resource}Resource_full(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testutil.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
        CheckDestroy:             testAcc{Resource}ResourceCheckDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAcc{Resource}ResourceConfig_full(),
                Check: resource.ComposeTestCheckFunc(
                    // Verify all required fields
                    resource.TestCheckResourceAttr("netbox_{resource}.test", "required_field", "value"),
                    // Verify all optional fields
                    resource.TestCheckResourceAttr("netbox_{resource}.test", "optional_field", "value"),
                    // Verify computed fields
                    resource.TestCheckResourceAttrSet("netbox_{resource}.test", "computed_field"),
                ),
            },
        },
    })
}
```

## Remaining Resources (6)

Based on the coverage analysis, these 6 resources need full tests:

1. **circuit_termination** - Circuit endpoint configuration
2. **cluster_group** - Cluster grouping/organization
3. **contact_assignment** - Contact assignments to objects
4. **contact_group** - Contact group hierarchy
5. **contact_role** - Contact role definitions
6. **wireless_link** - Wireless connection between interfaces

**Note:** All 6 resources have comprehensive test coverage via other test types (_basic, _update, _import, _removeOptionalFields, _externalDeletion). Adding _full tests completes the standard test suite.

---

## Implementation Strategy

### Single Batch Approach

Process all 6 resources in one focused session:
- **Estimated Time**: 4-6 hours
- **Single Commit**: All changes together
- **Rationale**: Small number, related functionality, clean commit history

---

## Resource Analysis

### 1. CircuitTermination

**Purpose**: Represents A-side or Z-side endpoints of circuits

**Required Fields:**
- circuit (ID or CID)
- term_side (A or Z)

**Optional Fields:**
- port_speed (integer)
- upstream_speed (integer)
- xconnect_id (string)
- pp_info (string - patch panel info)
- description (string)
- mark_connected (boolean)
- cable (reference)
- tags (set)
- custom_fields (map)

**Full Test Focus:**
- All optional fields populated
- Valid port speeds
- Patch panel information
- Cable reference handling

**Estimated Time**: 45 minutes

---

### 2. ClusterGroup

**Purpose**: Organize virtualization clusters into hierarchical groups

**Required Fields:**
- name (string)
- slug (string)

**Optional Fields:**
- description (string)
- parent (reference to parent cluster group)
- tags (set)
- custom_fields (map)

**Full Test Focus:**
- Description field
- Parent group hierarchy (if applicable)
- Tags and custom fields

**Estimated Time**: 30 minutes

---

### 3. ContactAssignment

**Purpose**: Assign contacts to objects with specific roles

**Required Fields:**
- object_type (string - e.g., "dcim.device")
- object_id (integer)
- contact (reference)
- role (reference)

**Optional Fields:**
- priority (string - primary, secondary, tertiary, inactive)
- tags (set)
- custom_fields (map)

**Full Test Focus:**
- All priority values
- Contact and role references
- Object type validation

**Estimated Time**: 45 minutes

---

### 4. ContactGroup

**Purpose**: Hierarchical organization of contacts

**Required Fields:**
- name (string)
- slug (string)

**Optional Fields:**
- parent (reference to parent contact group)
- description (string)
- tags (set)
- custom_fields (map)

**Full Test Focus:**
- Description field
- Parent group hierarchy
- Tags and custom fields

**Estimated Time**: 30 minutes

---

### 5. ContactRole

**Purpose**: Define roles for contact assignments

**Required Fields:**
- name (string)
- slug (string)

**Optional Fields:**
- description (string)
- tags (set)
- custom_fields (map)

**Full Test Focus:**
- Description field
- Tags and custom fields
- Role usage in assignments

**Estimated Time**: 30 minutes

---

### 6. WirelessLink

**Purpose**: Wireless connection between two interfaces

**Required Fields:**
- interface_a (reference)
- interface_b (reference)

**Optional Fields:**
- ssid (string, max 32 chars)
- status (string - connected, planned, decommissioning)
- tenant (reference)
- auth_type (string - open, wep, wpa-personal, wpa-enterprise)
- auth_cipher (string - auto, tkip, aes)
- auth_psk (string - pre-shared key)
- distance (float)
- distance_unit (string - km, m, mi, ft)
- description (string)
- comments (text)
- tags (set)
- custom_fields (map)

**Full Test Focus:**
- SSID validation
- Authentication configuration
- Distance with units
- Status values
- Tenant reference

**Estimated Time**: 60 minutes

---

## Test Configuration Patterns

### Full Configuration Template

Each resource needs a comprehensive configuration function:

```go
func testAcc{Resource}ResourceConfig_full() string {
    return `
# Create any prerequisite resources
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

# Full configuration with all optional fields
resource "netbox_{resource}" "test" {
  # Required fields
  required_field = "value"

  # Optional fields
  optional_field_1 = "value1"
  optional_field_2 = "value2"
  optional_bool    = true
  optional_int     = 100

  # Reference fields
  reference_field = netbox_other_resource.test.id

  # Description/comments
  description = "Full test configuration"
  comments    = "Testing all optional fields"

  # Tags
  tags = ["tag1", "tag2", "tag3"]

  # Custom fields (if applicable)
  custom_fields = {
    "field1" = "value1"
  }
}
`
}
```

### Check Functions

Verify all fields are properly set:

```go
Check: resource.ComposeTestCheckFunc(
    // Required fields
    resource.TestCheckResourceAttr("netbox_{resource}.test", "required_field", "value"),

    // Optional fields
    resource.TestCheckResourceAttr("netbox_{resource}.test", "optional_field_1", "value1"),
    resource.TestCheckResourceAttr("netbox_{resource}.test", "optional_field_2", "value2"),
    resource.TestCheckResourceAttr("netbox_{resource}.test", "optional_bool", "true"),
    resource.TestCheckResourceAttr("netbox_{resource}.test", "optional_int", "100"),

    // Reference fields (verify ID is set)
    resource.TestCheckResourceAttrPair(
        "netbox_{resource}.test", "reference_field",
        "netbox_other_resource.test", "id",
    ),

    // Description/comments
    resource.TestCheckResourceAttr("netbox_{resource}.test", "description", "Full test configuration"),

    // Tags (verify count and values)
    resource.TestCheckResourceAttr("netbox_{resource}.test", "tags.#", "3"),

    // Computed fields (verify they are set)
    resource.TestCheckResourceAttrSet("netbox_{resource}.test", "id"),
    resource.TestCheckResourceAttrSet("netbox_{resource}.test", "url"),
),
```

---

## Prerequisites and Dependencies

### Resource Dependencies

Some resources require prerequisite resources:

**CircuitTermination**:
- Circuit (circuit_resource)
- Potentially: Site, Provider, CircuitType

**ContactAssignment**:
- Contact (contact_resource)
- ContactRole (contact_role_resource)
- Object to assign to (e.g., Device, Site)

**ClusterGroup, ContactGroup, ContactRole**:
- Minimal dependencies (mostly self-contained)

**WirelessLink**:
- Two Interfaces (interface_resource)
- Device(s) for interfaces
- Site for devices
- DeviceType, DeviceRole, Manufacturer

### Test Isolation

Each test should:
- Create all prerequisites within the test
- Use unique names (with random suffixes)
- Clean up all created resources
- Not depend on external state

---

## Validation Criteria

For each test to be considered complete:

✅ Test creates resource with all optional fields populated
✅ Test verifies all required fields are correct
✅ Test verifies all optional fields are set
✅ Test verifies computed fields are populated
✅ Test verifies reference fields point to correct resources
✅ Test passes consistently (no flakiness)
✅ CheckDestroy function properly validates cleanup
✅ Code is formatted and passes linting

---

## Success Metrics

**Target**: 99/99 resources with full tests (100% coverage)

**Current Progress**: 93/99 (93.9%)

**Remaining**: 6 resources
- [ ] circuit_termination
- [ ] cluster_group
- [ ] contact_assignment
- [ ] contact_group
- [ ] contact_role
- [ ] wireless_link

**Definition of Done:**
- All 6 resources have _full tests
- All tests passing consistently
- Documentation updated
- Coverage analysis updated to show 100%

---

## Implementation Phases

### Phase 1: Simple Resources (3 resources)

**Resources**: cluster_group, contact_group, contact_role

**Rationale**: Simple schemas, minimal dependencies, quick wins

**Estimated Time**: 90 minutes

**Order**:
1. ContactRole (simplest - just name, slug, description, tags)
2. ContactGroup (adds parent hierarchy)
3. ClusterGroup (similar to ContactGroup)

---

### Phase 2: Assignment Resources (1 resource)

**Resources**: contact_assignment

**Rationale**: Requires multiple prerequisites, more complex setup

**Estimated Time**: 45 minutes

**Prerequisites Needed**:
- Site (for device)
- Manufacturer, DeviceType, DeviceRole (for device)
- Device (assignment target)
- Contact
- ContactRole

---

### Phase 3: Complex Resources (2 resources)

**Resources**: circuit_termination, wireless_link

**Rationale**: Most complex schemas, most dependencies

**Estimated Time**: 2 hours

**circuit_termination Prerequisites**:
- Provider
- CircuitType
- Circuit
- Site (for cable/termination context)

**wireless_link Prerequisites**:
- Site
- Manufacturer, Platform, DeviceType, DeviceRole
- Two Devices
- Two Interfaces

---

## Testing Strategy

### Local Testing

Before committing, verify each test:

```powershell
# Test individual resource
$env:TF_ACC='1'
go test ./internal/resources_acceptance_tests/... `
  -run 'TestAcc{Resource}Resource_full' `
  -v -timeout 30m
```

### Batch Testing

Test all 6 new tests together:

```powershell
$env:TF_ACC='1'
go test ./internal/resources_acceptance_tests/... `
  -run 'TestAcc(CircuitTermination|ClusterGroup|ContactAssignment|ContactGroup|ContactRole|WirelessLink)Resource_full' `
  -v -timeout 30m -p 1
```

### Full Regression Testing

Ensure no existing tests were broken:

```powershell
$env:TF_ACC='1'
go test ./internal/resources_acceptance_tests/... -v -timeout 60m
```

---

## Documentation Updates

After completion, update:

1. **ACCEPTANCE_TEST_COVERAGE_ANALYSIS.md**
   - Update Full Test coverage to 100% (99/99)
   - Remove from gap analysis
   - Update executive summary

2. **This Document (FULL_TEST_IMPLEMENTATION_PLAN.md)**
   - Mark all resources as complete
   - Add final summary with execution results
   - Document any issues encountered

3. **README.md** (if applicable)
   - Update test coverage statistics

---

## Expected Challenges

### 1. Complex Prerequisites

**Issue**: Some resources require many prerequisite resources
**Solution**: Create helper functions for common prerequisite sets

```go
func createDevicePrerequisites(t *testing.T) string {
    return `
resource "netbox_site" "test" { ... }
resource "netbox_manufacturer" "test" { ... }
resource "netbox_device_type" "test" { ... }
resource "netbox_device_role" "test" { ... }
`
}
```

### 2. Field Value Validation

**Issue**: Some fields have specific validation rules
**Solution**: Use valid values from schema documentation

### 3. Reference Resolution

**Issue**: Some fields require valid references
**Solution**: Create referenced resources first, use TestCheckResourceAttrPair

### 4. Hierarchical Resources

**Issue**: Parent-child relationships (ContactGroup, ClusterGroup)
**Solution**: Create parent first, or test without parent for full test

---

## Timeline

**Start Date**: TBD
**Target Completion**: TBD
**Estimated Duration**: 4-6 hours (single focused session)

### Milestones

- [ ] Phase 1 Complete: Simple resources (ContactRole, ContactGroup, ClusterGroup)
- [ ] Phase 2 Complete: Assignment resource (ContactAssignment)
- [ ] Phase 3 Complete: Complex resources (CircuitTermination, WirelessLink)
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Final commit merged

---

## Related Work

**Prerequisites**:
- ✅ Validation tests (100% - 97/97 resources)
- ✅ Optional field tests (100% - 99/99 resources)
- ✅ Update tests (100% - 97/97 resources)
- ✅ Import tests (100% - 97/97 resources)

**Next Phase After Full Tests**:
- Edge case testing
- Consistency/LiteralNames tests
- Reference handling tests
- Performance testing

---

## Notes

- All 6 resources already have _basic tests, so we can reference those for minimal configuration
- Most resources already have _update and _removeOptionalFields tests, providing examples of valid optional field values
- Focus on completeness - verify EVERY optional field works correctly
- Full tests are valuable for regression prevention and feature validation

---

*Created: January 15, 2026*
*Status: Planning Phase*
*Priority: Medium*
*Target: 100% Full Test Coverage*
