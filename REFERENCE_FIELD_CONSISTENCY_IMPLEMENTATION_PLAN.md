# Reference Field Display Consistency Implementation Plan

## Overview

This document outlines the implementation plan for resolving reference field display consistency issues in Terraform plans. The issue manifests as confusing plan output like `~ device_type = 'PowerEdge R640' -> '9'` where users input names but see IDs in the plan.

## Problem Statement

**Current Issue**: Terraform plans show inconsistent reference field values, displaying name → ID changes even when no actual change occurred. This happens because:

1. Users input human-readable names (e.g., "PowerEdge R640")
2. NetBox API returns objects with numeric IDs
3. Terraform sees the string change from name to ID as a modification
4. Plan shows confusing `'PowerEdge R640' -> '9'` even though both refer to the same object

**Goal**: Implement DiffSuppressFunc-based solution that maintains user-friendly name input while ensuring consistent plan display, following AWS provider patterns.

## Architecture Overview

### Implementation Strategy

**Approach**: Hybrid DiffSuppressFunc solution inspired by AWS provider patterns
- ✅ Keep user-friendly name/ID acceptance (superior to AWS approach)
- ✅ Add DiffSuppressFunc to suppress equivalent name ↔ ID diffs
- ✅ Use proven patterns from the most mature Terraform provider
- ✅ No complex plan modification - works within framework patterns

### Core Components

1. **Enhanced ReferenceAttribute Functions** - Extend existing `nbschema.ReferenceAttribute()`
2. **Generic DiffSuppressFunc** - Centralized logic for name ↔ ID equivalency checking
3. **Lookup Enhancement** - Extend `netboxlookup` package for reverse ID → name resolution
4. **Resource Type Registry** - Mapping system for resource type detection from attribute paths

## Detailed Implementation Plan

### Phase 1: Foundation Infrastructure

#### 1.1 Enhanced Schema Functions
**File**: `internal/schema/attributes.go`

```go
// New function with DiffSuppressFunc support
func ReferenceAttributeWithConsistentDisplay(targetResource string, description string) schema.StringAttribute {
    return schema.StringAttribute{
        MarkdownDescription: description,
        Optional:           true,
        DiffSuppressFunc:   suppressReferenceEquivalent,
        PlanModifiers: []planmodifier.String{
            stringplanmodifier.UseStateForUnknown(),
        },
    }
}

func RequiredReferenceAttributeWithConsistentDisplay(targetResource string, description string) schema.StringAttribute {
    return schema.StringAttribute{
        MarkdownDescription: description,
        Required:           true,
        DiffSuppressFunc:   suppressReferenceEquivalent,
        PlanModifiers: []planmodifier.String{
            stringplanmodifier.UseStateForUnknown(),
        },
    }
}
```

#### 1.2 DiffSuppressFunc Implementation
**File**: `internal/schema/diff_suppress.go` (new)

```go
// suppressReferenceEquivalent suppresses diffs when old and new values
// refer to the same NetBox object but use different representations
// (name vs ID vs slug)
func suppressReferenceEquivalent(k, old, new string, d *schema.ResourceData) bool {
    // If values are identical, no diff needed
    if old == new {
        return true
    }

    // Detect resource type from attribute path
    resourceType := getResourceTypeFromPath(k)
    if resourceType == "" {
        // Can't determine type, let Terraform show diff
        return false
    }

    // Check if values refer to the same object
    return areEquivalentReferences(old, new, resourceType)
}

// Helper functions for reference resolution and comparison
func areEquivalentReferences(val1, val2, resourceType string) bool {
    // Implementation will resolve both values to canonical IDs and compare
}

func getResourceTypeFromPath(attributePath string) string {
    // Map attribute paths to resource types for lookup
    // e.g., "device_type" -> "dcim.device-types"
}
```

#### 1.3 Lookup Enhancement
**File**: `internal/netboxlookup/reverse_lookup.go` (new)

```go
// ResolveToCanonicalID resolves any reference (name, slug, ID) to canonical ID
func ResolveToCanonicalID(ctx context.Context, client *netbox.APIClient, value string, resourceType string) (int32, error) {
    // Try parsing as ID first
    if id, err := strconv.ParseInt(value, 10, 32); err == nil {
        return int32(id), nil
    }

    // Use existing lookup functions to resolve name/slug to ID
    switch resourceType {
    case "dcim.device-types":
        return resolveDeviceTypeToID(ctx, client, value)
    case "dcim.device-roles":
        return resolveDeviceRoleToID(ctx, client, value)
    // ... other resource types
    }
}
```

### Phase 2: Testing Infrastructure

#### 2.1 Unit Test Coverage
**Files**:
- `internal/schema/diff_suppress_test.go`
- `internal/netboxlookup/reverse_lookup_test.go`

**Coverage Areas**:
- DiffSuppressFunc logic with various input combinations
- Resource type detection from attribute paths
- Reverse lookup functionality for all supported resource types
- Edge cases: empty values, invalid IDs, API errors

#### 2.2 Acceptance Test Framework
**File**: `internal/acctest/reference_consistency_test_helpers.go` (new)

```go
// Helper functions for testing reference field consistency
func TestReferenceFieldConsistency(t *testing.T, resourceName string, referenceFields map[string]string) {
    // Test that name input doesn't show as change in subsequent plans
}
```

### Phase 3: Resource Migration (Batched Implementation)

## Implementation Batches

### Batch 1: Core Infrastructure (Foundation) - 2-3 days
**Scope**: Essential infrastructure without resource changes
- ✅ Enhanced schema functions with DiffSuppressFunc
- ✅ Core DiffSuppressFunc implementation
- ✅ Resource type registry and path mapping
- ✅ Reverse lookup infrastructure
- ✅ Unit test framework for all components
- ✅ Basic acceptance test helpers

**Success Criteria**:
- [ ] All unit tests pass for diff suppression logic
- [ ] Resource type detection works for 5+ test cases
- [ ] Reverse lookup resolves names to IDs for 3+ resource types
- [ ] No regression in existing functionality

**Files Modified**:
- `internal/schema/attributes.go` (new functions)
- `internal/schema/diff_suppress.go` (new file)
- `internal/netboxlookup/reverse_lookup.go` (new file)
- `internal/acctest/reference_consistency_helpers.go` (new file)

### Batch 2: High-Impact Resources - 3-4 days
**Scope**: Most commonly used resources with reference fields
- ✅ `device_resource.go` (device_type, role, site, location, rack, tenant, platform)
- ✅ `virtual_machine_resource.go` (site, cluster, role, tenant, platform)
- ✅ `site_resource.go` (region, group, tenant)
- ✅ `rack_resource.go` (site, location, tenant, role, rack_type)

**Success Criteria**:
- [ ] Acceptance tests pass for all reference fields in scope
- [ ] Plan consistency verified: no name → ID changes shown
- [ ] User can input names, slugs, or IDs interchangeably
- [ ] All existing functionality preserved

**User Scenarios Tested**:
- Create device with device_type name → subsequent plan shows no changes
- Update device from name to ID → no spurious diff shown
- Mixed reference types (some names, some IDs) work correctly

### Batch 3: Network Resources - 2-3 days
**Scope**: Network-focused resources
- ✅ `ip_address_resource.go` (vrf, tenant)
- ✅ `prefix_resource.go` (site, vrf, tenant, vlan)
- ✅ `vlan_resource.go` (site, tenant, group)
- ✅ `aggregate_resource.go` (rir, tenant)

**Success Criteria**:
- [ ] Network resource reference fields show consistent behavior
- [ ] IP address and prefix tenant references work correctly
- [ ] VLAN site references maintain consistency

### Batch 4: Organization Resources - 2-3 days
**Scope**: Organizational hierarchy resources
- ✅ `location_resource.go` (site, parent, tenant)
- ✅ `region_resource.go` (parent)
- ✅ `site_group_resource.go` (parent)
- ✅ `tenant_resource.go` (group)
- ✅ `tenant_group_resource.go` (parent)
- ✅ `contact_resource.go` (group)

**Success Criteria**:
- [ ] Hierarchical relationships work with consistent reference display
- [ ] Parent-child relationships maintain plan consistency

### Batch 5: Device Components - 3-4 days
**Scope**: Device-related component resources
- ✅ `power_port_resource.go`, `power_outlet_resource.go`
- ✅ `console_port_resource.go`, `console_server_port_resource.go`
- ✅ `interface_resource.go`, `front_port_resource.go`, `rear_port_resource.go`
- ✅ Template resources for above components

**Success Criteria**:
- [ ] Device component references maintain consistency
- [ ] Template references work correctly
- [ ] Device associations show stable plan output

### Batch 6: Specialized Resources - 2-3 days
**Scope**: Remaining specialized resources
- ✅ `circuit_resource.go`, `circuit_termination_resource.go`
- ✅ `cable_resource.go`, `module_resource.go`
- ✅ `inventory_item_resource.go`, `service_resource.go`
- ✅ Other miscellaneous resources with reference fields

**Success Criteria**:
- [ ] All remaining resources maintain reference field consistency
- [ ] Edge case resources work correctly
- [ ] Complete coverage across provider

## Testing Strategy

### Unit Testing Approach
**Test Categories**:
1. **DiffSuppressFunc Logic**: Various input combinations (name/ID/slug, same object)
2. **Resource Type Detection**: Attribute path → resource type mapping
3. **Reverse Lookup**: Name/slug resolution to canonical IDs
4. **Edge Cases**: Invalid inputs, API errors, empty values

### Acceptance Testing Approach
**Test Scenarios** (per batch):
1. **Plan Consistency**: Name input → no change in subsequent plan
2. **Mixed References**: Combination of names, slugs, IDs work together
3. **Update Scenarios**: Change from name to ID shows no diff
4. **Error Handling**: Invalid references produce meaningful errors

### Integration Testing
**Focus Areas**:
- End-to-end workflows with multiple reference types
- Complex resource hierarchies (location → site → region)
- Resource dependencies and ordering

## Risk Assessment & Mitigation

### Technical Risks
1. **Performance Impact**: DiffSuppressFunc adds API calls
   - *Mitigation*: Cache lookups, optimize API usage patterns
2. **API Rate Limits**: Increased lookup calls may hit limits
   - *Mitigation*: Implement intelligent caching and batching
3. **Backwards Compatibility**: Changes might break existing configurations
   - *Mitigation*: Thorough testing with existing state files

### Implementation Risks
1. **Complexity Creep**: Resource type detection becomes unwieldy
   - *Mitigation*: Keep mapping simple, use consistent patterns
2. **Test Coverage**: Missing edge cases cause runtime issues
   - *Mitigation*: Comprehensive test matrix, real-world scenario testing

## Success Metrics

### User Experience Metrics
- ✅ No confusing name → ID changes in plan output
- ✅ Users can input names/slugs/IDs interchangeably
- ✅ Plan consistency across apply cycles
- ✅ Meaningful error messages for invalid references

### Technical Metrics
- ✅ No regression in existing functionality
- ✅ Test coverage >90% for new components
- ✅ Performance impact <10% on plan operations
- ✅ Acceptance test pass rate 100% for modified resources

## Development Timeline

**Total Estimated Duration**: 14-18 days

**Phase Breakdown**:
- Phase 1 (Infrastructure): 2-3 days
- Phase 2 (Testing Setup): 1-2 days
- Phase 3 (Resource Migration): 11-13 days across 6 batches

**Parallel Work Opportunities**:
- Unit tests can be developed alongside infrastructure
- Resource batches can be developed in parallel after Batch 1
- Documentation can be updated incrementally

## Completion Criteria

### Batch Completion Criteria
Each batch must meet these criteria before proceeding:

1. **Functionality**: All reference fields work with name/slug/ID input
2. **Plan Consistency**: No spurious name → ID changes in plan output
3. **Test Coverage**: Unit tests + acceptance tests pass
4. **Backwards Compatibility**: Existing configurations continue working
5. **Error Handling**: Invalid references produce helpful error messages

### Overall Project Completion
1. **Full Provider Coverage**: All resources with reference fields updated
2. **User Experience**: Zero confusing plan output for reference fields
3. **Performance**: No significant degradation in provider performance
4. **Documentation**: Updated examples and troubleshooting guide
5. **Community Validation**: Beta testing with key users successful

---

## Next Steps

1. **Create Feature Branch**: `feature/reference-field-display-consistency` ✅
2. **Begin Batch 1**: Infrastructure development with test-driven approach
3. **Continuous Integration**: Set up automated testing for each batch
4. **User Feedback Loop**: Early testing with common use cases

This implementation plan provides a structured, testable approach to resolving the reference field consistency issue while maintaining backward compatibility and following proven patterns from the AWS provider.
