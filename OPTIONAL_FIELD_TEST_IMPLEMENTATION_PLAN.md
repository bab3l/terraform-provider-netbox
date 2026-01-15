# Optional Field Test Implementation Plan

## Overview

Implement comprehensive optional field tests for remaining resources to verify that optional fields can be safely removed from configurations and properly cleared in Netbox. These tests ensure Terraform correctly detects and applies changes when optional fields are removed.

**Current Status:** 5 resources remaining (97.9% complete - 97/99 resources with tests)

## Test Pattern

```go
func TestAcc{Resource}Resource_removeOptionalFields(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testutil.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAcc{Resource}ResourceConfig_full(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("netbox_{resource}.test", "optional_field", "value"),
                ),
            },
            {
                Config: testAcc{Resource}ResourceConfig_basic(),
                Check: resource.ComposeTestCheckFunc(
                    testutil.TestCheckNoResourceAttr("netbox_{resource}.test", "optional_field"),
                ),
            },
        },
    })
}
```

## Current Status

### Completed Resources (97/99)

All major resource categories are complete:
- ✅ Core Infrastructure (Site, Device, Rack, etc.)
- ✅ DCIM Components (Ports, Bays, Templates)
- ✅ IPAM (IPAddress, Prefix, VLAN, VRF, etc.)
- ✅ Virtualization (VirtualMachine, Cluster, etc.)
- ✅ Circuits & Providers
- ✅ VPN & Security (IKE, IPSec, Tunnels)
- ✅ Tenancy & Contacts
- ✅ Templates & Configuration
- ✅ Wireless Infrastructure

### Skipped Resources (2/99)

1. **FHRPGroupAssignment** - Legitimately skipped (no removable optional fields)
   - Only has: group, interface_type, interface_id, priority
   - All fields are either required or should not be removed

2. **L2VPNTermination** - Skipped due to provider bug
   - Provider bug: Removing tags causes state inconsistency
   - Issue documented in test comments
   - Will be addressed when provider bug is fixed

### Remaining Work (5 resources)

These resources need optional field tests added:

1. **custom_field** - Test removing: description, weight, filter_logic, etc.
2. **custom_field_choice_set** - Test removing: description, extra_choices, etc.
3. **ike_proposal** - Test removing: description, authentication_algorithm, sa_lifetime
4. **ipsec_policy** - Test removing: description, pfs_group
5. **virtual_device_context** - Test removing: description, tenant, etc.

## Implementation Phases

### Phase 1: IPSec/VPN Resources (3 resources)

**Priority**: High - Recently added resources, need complete test coverage

**Resources:**
- ike_proposal
- ipsec_policy
- ipsec_proposal (if missing)

**Optional Fields to Test:**
- description (common to all)
- authentication_algorithm (IKE)
- sa_lifetime (IKE)
- pfs_group (IPSec Policy)

**Estimated Time**: 2 hours

**Pattern:**
```go
func TestAccIKEProposalResource_removeOptionalFields(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testutil.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccIKEProposalResourceConfig_withOptionalFields(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("netbox_ike_proposal.test", "description", "Test Description"),
                    resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_algorithm", "hmac-sha256"),
                ),
            },
            {
                Config: testAccIKEProposalResourceConfig_basic(),
                Check: resource.ComposeTestCheckFunc(
                    testutil.TestCheckNoResourceAttr("netbox_ike_proposal.test", "description"),
                    testutil.TestCheckNoResourceAttr("netbox_ike_proposal.test", "authentication_algorithm"),
                ),
            },
        },
    })
}
```

---

### Phase 2: Custom Field Resources (2 resources)

**Priority**: Medium - Core Netbox extensibility features

**Resources:**
- custom_field
- custom_field_choice_set

**Optional Fields to Test:**

**custom_field:**
- description
- group_name
- weight
- filter_logic
- default value fields

**custom_field_choice_set:**
- description
- base_choices (if optional)
- extra_choices

**Estimated Time**: 2 hours

**Complexity Notes:**
- custom_field has many optional fields - may need separate tests for different field types
- May need to test removing fields specific to certain field types (string, integer, etc.)

---

### Phase 3: Virtual Device Context (1 resource)

**Priority**: Low - Less commonly used feature

**Resources:**
- virtual_device_context

**Optional Fields to Test:**
- description
- tenant
- comments
- tags

**Estimated Time**: 1 hour

---

## Test Validation Criteria

For each test to be considered complete:

✅ Test creates resource with optional fields populated
✅ Test verifies optional fields are present in state
✅ Test removes optional fields from configuration
✅ Test verifies fields are absent using `TestCheckNoResourceAttr`
✅ Test passes consistently (no flakiness)
✅ Code is formatted and passes linting

## Success Metrics

**Target**: 99/99 resources with optional field tests (100% coverage)

- [x] Phase 1: IPSec/VPN resources (3/3)
- [ ] Phase 2: Custom Field resources (0/2)
- [ ] Phase 3: Virtual Device Context (0/1)

**Definition of Done:**
- All 5 remaining resources have removeOptionalFields tests
- All tests passing consistently
- Documentation updated
- Coverage analysis updated to show 100% (99/99)

## Resource Organization

### Batch Strategy

**Single Batch Approach** (Recommended):
- Process all 5 remaining resources in one focused session
- Total estimated time: 5 hours
- Single commit with all changes

**Rationale:**
- Small number of remaining resources
- Related functionality (IPSec/VPN group together)
- Faster completion
- Clean commit history

---

## Test Framework Integration

### Using Existing Helper Functions

The `testutil.TestCheckNoResourceAttr` helper is already implemented and widely used:

```go
// TestCheckNoResourceAttr checks that an attribute does not exist in state
func TestCheckNoResourceAttr(resourceName, attributeName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("Not found: %s", resourceName)
        }
        if val, ok := rs.Primary.Attributes[attributeName]; ok && val != "" {
            return fmt.Errorf("Attribute %s still exists with value: %s", attributeName, val)
        }
        return nil
    }
}
```

### Test Configuration Patterns

Each test typically needs two configuration functions:

1. **_full()** - Configuration with optional fields populated
2. **_basic()** - Minimal configuration without optional fields

Example:
```go
func testAccIKEProposalResourceConfig_withOptionalFields() string {
    return `
resource "netbox_ike_proposal" "test" {
  name                     = "test-proposal"
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  group                    = 14
  description             = "Test Description"
  authentication_algorithm = "hmac-sha256"
  sa_lifetime             = 3600
}
`
}

func testAccIKEProposalResourceConfig_basic() string {
    return `
resource "netbox_ike_proposal" "test" {
  name                  = "test-proposal"
  authentication_method = "preshared-keys"
  encryption_algorithm  = "aes-256-cbc"
  group                 = 14
}
`
}
```

---

## Expected Challenges

### 1. Complex Optional Fields

**Issue**: Some resources have many optional fields
**Solution**: Focus on 2-3 representative optional fields per test

### 2. Dependent Optional Fields

**Issue**: Some optional fields may depend on others
**Solution**: Test removal of independent fields, document dependencies

### 3. API Behavior Variations

**Issue**: Different Netbox APIs may handle null/empty differently
**Solution**: Use TestCheckNoResourceAttr which works regardless of API behavior

---

## Documentation Updates

After completion, update:

1. **ACCEPTANCE_TEST_COVERAGE_ANALYSIS.md**
   - Update removeOptionalFields coverage to 100% (99/99)
   - Remove from gap analysis

2. **README.md** (if applicable)
   - Update test coverage statistics

3. **CHANGELOG.md**
   - Add entry for optional field test completion

---

## Timeline

**Start Date**: TBD
**Target Completion**: TBD
**Estimated Duration**: 1 day (5 hours focused work)

### Milestones

- [ ] Phase 1 Complete: IPSec/VPN resources
- [ ] Phase 2 Complete: Custom Field resources
- [ ] Phase 3 Complete: Virtual Device Context
- [ ] Documentation updated
- [ ] All tests passing
- [ ] Final commit merged

---

## Related Work

**Prerequisite**: Validation tests (✅ Complete - 97/97 resources)
**Next Phase**: Edge case testing, reference handling, or hierarchical tests

---

*Created: January 15, 2026*
*Status: Planning Phase*
*Priority: Medium*
