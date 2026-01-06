# Custom Fields Partial Management Implementation Plan

## Problem Statement

### Current Behavior (BROKEN)
The provider currently uses a `custom_fields = [...]` Set approach that exhibits critical data loss:

1. **Bug**: When a resource is managed in Terraform WITHOUT `custom_fields` in the config, updating ANY field causes ALL custom fields in NetBox to be **deleted**
2. **All-or-nothing**: Users must manage ALL custom fields or NONE - no partial management
3. **External management impossible**: Custom fields set via NetBox UI or other tools are lost on next Terraform apply

### Production Impact
Users report that running `terraform apply` to update unrelated fields (e.g., description) **deletes all custom fields** from their production resources. This is a critical data loss bug.

### Root Cause
In the Update path:
```go
func ApplyCustomFields(ctx context.Context, request T, customFields types.Set, diags *diag.Diagnostics) {
    if !IsSet(customFields) {  // customFields is null when not in config
        return  // Returns WITHOUT calling request.SetCustomFields()
    }
    // ... convert and set custom fields
    request.SetCustomFields(cfMap)
}
```

When `custom_fields` is omitted from config:
1. `customFields` is **null** → `ApplyCustomFields` returns early
2. The API request struct doesn't include `custom_fields` key (due to `omitempty`)
3. NetBox interprets missing key as "**clear all custom fields**"

## Design Goals

### Primary Goals
1. **Partial Management**: Users can specify SOME custom fields in Terraform, others are preserved
2. **No Data Loss**: Updating non-custom-field attributes never deletes custom fields
3. **Explicit Control**: Clear semantics for "manage this field" vs "ignore this field"
4. **Backward Compatible**: Existing configs continue to work

### User Stories

**Story 1: External Custom Field Management**
```hcl
# Custom fields are managed in NetBox UI, not Terraform
resource "netbox_device" "server" {
  name = "server-01"
  # custom_fields omitted - should preserve ALL existing values
}
```
✅ Expected: All custom fields in NetBox are preserved
❌ Current: All custom fields are DELETED on update

**Story 2: Partial Custom Field Management**
```hcl
# Manage ONLY the "environment" field, leave others alone
resource "netbox_device" "server" {
  name = "server-01"
  custom_fields = [
    {
      name  = "environment"
      type  = "text"
      value = "production"
    }
    # Other custom fields (e.g., "owner", "cost_center") exist in NetBox
    # but are NOT managed by Terraform - should be preserved
  ]
}
```
✅ Expected: "environment" is managed, others preserved
❌ Current: Only "environment" exists, others are DELETED

**Story 3: Explicit Field Removal**
```hcl
# Want to remove a specific custom field
resource "netbox_device" "server" {
  name = "server-01"
  custom_fields = [
    {
      name  = "environment"
      type  = "text"
      value = "" # or null - should remove this field
    }
  ]
}
```
✅ Expected: "environment" is removed, others preserved

## Design Options

### Option 1: Read-Merge-Write (CHOSEN)
**Approach**: During Update, always read current custom fields, merge with config values, send complete map

**Pros**:
- ✅ Preserves unmanaged fields automatically
- ✅ Works with current Set schema
- ✅ Minimal config changes
- ✅ Clear semantics: "fields in config are managed, others preserved"

**Cons**:
- Requires extra API read during Update (but we already do this for many resources)
- Slightly more complex Update logic

**Implementation**:
```go
func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state, plan DeviceResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

    // Build update request...

    // CUSTOM FIELDS: Read current, merge with plan
    if !plan.CustomFields.IsNull() && !plan.CustomFields.IsUnknown() {
        // User specified custom_fields in config - merge with existing
        cfMap := MergeCustomFields(ctx, r.client, deviceID, plan.CustomFields, &resp.Diagnostics)
        request.SetCustomFields(cfMap)
    } else {
        // User did NOT specify custom_fields - preserve all existing
        // Don't call SetCustomFields() - omit from API request
    }

    // Execute update...
}
```

### Option 2: Nested Block with Ignore Mode
**Approach**: Add `ignore_unmanaged = true` flag
```hcl
custom_fields {
  ignore_unmanaged = true  # Don't touch fields not listed here
  field {
    name = "environment"
    type = "text"
    value = "production"
  }
}
```

**Pros**:
- Explicit user control

**Cons**:
- Breaking schema change
- More complex config
- Unclear default behavior

### Option 3: Lifecycle Ignore Changes
**Approach**: Use Terraform's `ignore_changes`
```hcl
resource "netbox_device" "server" {
  lifecycle {
    ignore_changes = [custom_fields]
  }
}
```

**Cons**:
- ❌ Doesn't help - still deletes fields on first apply
- ❌ All-or-nothing on Terraform side
- ❌ Doesn't solve the NetBox API deletion issue

## Chosen Design: Read-Merge-Write (Option 1)

### Implementation Strategy

#### Phase 1: Fix Update Path (Critical Bug Fix)
Modify `ApplyCustomFields` to support merge behavior:

```go
// ApplyCustomFieldsWithMerge handles custom fields in Update operations.
// If customFields is set in config:
//   - Reads current values from state (or API if needed)
//   - Merges plan values with current values
//   - Sends complete merged map to API
// If customFields is null (not in config):
//   - Reads current values from state
//   - Sends current values to preserve them
func ApplyCustomFieldsWithMerge[T CustomFieldsSetter](
    ctx context.Context,
    request T,
    planCustomFields types.Set,
    stateCustomFields types.Set,
    diags *diag.Diagnostics,
) {
    // If plan has custom_fields, merge with state
    if IsSet(planCustomFields) {
        // User specified custom_fields - manage these, preserve others
        merged := MergeCustomFieldSets(ctx, planCustomFields, stateCustomFields, diags)
        if diags.HasError() {
            return
        }
        request.SetCustomFields(merged)
        return
    }

    // If state has custom_fields but plan doesn't, preserve state values
    if IsSet(stateCustomFields) {
        // User did NOT specify custom_fields - preserve all existing
        stateMap := CustomFieldSetToMap(ctx, stateCustomFields, diags)
        if diags.HasError() {
            return
        }
        request.SetCustomFields(stateMap)
        return
    }

    // Neither plan nor state has custom_fields - send empty map to be safe
    request.SetCustomFields(map[string]interface{}{})
}

// MergeCustomFieldSets merges plan custom fields with state custom fields.
// Plan fields override state fields with same name.
// State fields not in plan are preserved.
func MergeCustomFieldSets(
    ctx context.Context,
    plan types.Set,
    state types.Set,
    diags *diag.Diagnostics,
) map[string]interface{} {
    result := make(map[string]interface{})

    // First, add all state custom fields
    if IsSet(state) {
        var stateFields []CustomFieldModel
        stateDiags := state.ElementsAs(ctx, &stateFields, false)
        diags.Append(stateDiags...)
        if diags.HasError() {
            return result
        }

        for _, cf := range stateFields {
            name := cf.Name.ValueString()
            value := ConvertCustomFieldValue(cf)
            result[name] = value
        }
    }

    // Then, overlay plan custom fields (these override state)
    if IsSet(plan) {
        var planFields []CustomFieldModel
        planDiags := plan.ElementsAs(ctx, &planFields, false)
        diags.Append(planDiags...)
        if diags.HasError() {
            return result
        }

        for _, cf := range planFields {
            name := cf.Name.ValueString()
            value := ConvertCustomFieldValue(cf)
            if value == nil || (cf.Value.IsNull() || cf.Value.ValueString() == "") {
                // Empty value = remove this field
                delete(result, name)
            } else {
                result[name] = value
            }
        }
    }

    return result
}
```

#### Phase 2: Update All Resources
Modify Update() methods to use new merge-aware function:

```go
// OLD (broken):
utils.ApplyCustomFields(ctx, &request, data.CustomFields, &resp.Diagnostics)

// NEW (fixed):
utils.ApplyCustomFieldsWithMerge(ctx, &request, data.CustomFields, state.CustomFields, &resp.Diagnostics)
```

### Semantic Behavior

| Scenario | Behavior |
|----------|----------|
| `custom_fields` omitted from config | All existing custom fields preserved on update |
| `custom_fields = []` | All custom fields removed |
| `custom_fields = [{ name="env", value="prod" }]` | "env" set to "prod", other fields preserved |
| `custom_fields = [{ name="env", value="" }]` | "env" removed, other fields preserved |

## Test Strategy

### Test 1: Custom Fields Preservation (Current Issue)
```go
// TestAccDeviceResource_CustomFieldsPreservation
// Step 1: Create with custom_fields = [field_a, field_b]
// Step 2: Update description WITHOUT custom_fields in config
// Expected: field_a and field_b still exist in NetBox
// Current: FAILS - fields are deleted
```

### Test 2: Partial Management
```go
// TestAccDeviceResource_CustomFieldsPartialManagement
// Step 1: Create with custom_fields = [field_a, field_b]
// Step 2: Update with custom_fields = [field_a] (removed field_b from config)
// Expected: field_a managed, field_b preserved in NetBox
// Step 3: Update with custom_fields = [field_a, field_c]
// Expected: field_a managed, field_b preserved, field_c added
```

### Test 3: Explicit Removal
```go
// TestAccDeviceResource_CustomFieldsExplicitRemoval
// Step 1: Create with custom_fields = [field_a, field_b]
// Step 2: Update with custom_fields = [{ name="field_a", value="" }]
// Expected: field_a removed, field_b preserved
```

### Test 4: Complete Removal
```go
// TestAccDeviceResource_CustomFieldsCompleteRemoval
// Step 1: Create with custom_fields = [field_a, field_b]
// Step 2: Update with custom_fields = []
// Expected: All custom fields removed
```

## Implementation Batches

### Batch 1: Core Utilities (Foundation)
**Priority**: CRITICAL
**Files**: 1 file
**Estimated Time**: 2 hours

1. Add new merge functions to `internal/utils/request_helpers.go`:
   - `ApplyCustomFieldsWithMerge()`
   - `MergeCustomFieldSets()`
   - `CustomFieldSetToMap()` (may already exist)

2. Add unit tests for merge logic:
   - `internal/utils/request_helpers_test.go`
   - Test all merge scenarios
   - Test empty/null handling

**Success Criteria**:
- [ ] `MergeCustomFieldSets()` correctly merges plan + state
- [ ] Empty value removes field from result
- [ ] Null plan preserves all state fields
- [ ] Unit tests pass

### Batch 2: Device Resource (Pilot Implementation)
**Priority**: CRITICAL
**Files**: 2 files
**Estimated Time**: 3 hours

1. Update `internal/resources/device_resource.go`:
   - Modify `Update()` to use `ApplyCustomFieldsWithMerge()`
   - Ensure state is read before building request

2. Create comprehensive acceptance tests:
   - `internal/resources_acceptance_tests/device_custom_fields_preservation_test.go`
   - All 4 test scenarios above
   - Verify against live NetBox instance

**Success Criteria**:
- [ ] All 4 acceptance tests pass
- [ ] No data loss when updating other fields
- [ ] Partial management works as expected
- [ ] Explicit removal works

### Batch 3: High-Priority Resources (IPAM Core)
**Priority**: HIGH
**Files**: 15-20 files
**Estimated Time**: 4-6 hours

Update resources that commonly use custom fields:
- `ip_address_resource.go`
- `prefix_resource.go`
- `vlan_resource.go`
- `vrf_resource.go`
- `aggregate_resource.go`
- `asn_resource.go`
- `asn_range_resource.go`

**Pattern**:
```go
// Get state for merge
var state DeviceResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

// ... build request ...

// Apply custom fields with merge
utils.ApplyCustomFieldsWithMerge(ctx, &request, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
```

**Success Criteria**:
- [ ] All resources use `ApplyCustomFieldsWithMerge()`
- [ ] State is read in Update() for context
- [ ] At least 2 acceptance tests per resource

### Batch 4: DCIM Resources
**Priority**: HIGH
**Files**: 20-25 files
**Estimated Time**: 6-8 hours

- Site, Location, Region, Rack resources
- Device Type, Module Type
- Interface, Console Port, Power Port, etc.
- Device Bay, Inventory Item

### Batch 5: Circuits & VPN Resources
**Priority**: MEDIUM
**Files**: 15-20 files
**Estimated Time**: 4-6 hours

- Circuit, Provider, Termination
- Tunnel, VPN resources
- L2VPN

### Batch 6: Virtualization & Tenancy
**Priority**: MEDIUM
**Files**: 10-15 files
**Estimated Time**: 3-4 hours

- Virtual Machine, Cluster
- Tenant, Tenant Group
- Contact resources

### Batch 7: Templates & Extras
**Priority**: LOW
**Files**: 10-15 files
**Estimated Time**: 3-4 hours

- Component Templates
- Config Context, Config Template
- Custom Link, Webhook, etc.

### Batch 8: Comprehensive Testing
**Priority**: CRITICAL
**Files**: Test files
**Estimated Time**: 4-6 hours

1. Run full acceptance test suite
2. Test import scenarios with custom fields
3. Test state migration edge cases
4. Performance testing with many custom fields

### Batch 9: Documentation
**Priority**: HIGH
**Files**: Documentation files
**Estimated Time**: 2-3 hours

1. Update CHANGELOG.md with breaking changes notice
2. Update examples/ with custom field patterns
3. Add migration guide for users
4. Document partial management feature

## Migration Guide for Users

### For Users Experiencing Data Loss

**Immediate Workaround**:
```hcl
# Option 1: Add all custom fields to config (tedious but safe)
resource "netbox_device" "server" {
  name = "server-01"
  custom_fields = [
    { name = "field1", type = "text", value = "value1" },
    { name = "field2", type = "text", value = "value2" },
    # ... all fields ...
  ]
}

# Option 2: Use lifecycle ignore_changes (prevents all CF updates)
resource "netbox_device" "server" {
  name = "server-01"
  lifecycle {
    ignore_changes = [custom_fields]
  }
}
```

### After Fix Is Deployed

**Recommended Pattern** (v0.0.13+):
```hcl
# Manage only specific fields, others preserved
resource "netbox_device" "server" {
  name = "server-01"

  # Only manage these specific custom fields
  # Other custom fields in NetBox are preserved
  custom_fields = [
    {
      name  = "environment"
      type  = "text"
      value = "production"
    }
  ]
}
```

## Rollout Plan

### Version 0.0.13 (Emergency Fix)
**Target**: 2-3 days
- Batch 1: Core utilities
- Batch 2: Device resource (pilot)
- Batch 3: High-priority IPAM resources
- Release with clear documentation

### Version 0.0.14 (Complete Fix)
**Target**: 1-2 weeks
- Batches 4-7: All remaining resources
- Batch 8: Comprehensive testing
- Batch 9: Full documentation

## Risk Assessment

### Risks
1. **Breaking Change**: Behavior changes from "replace all" to "merge"
   - **Mitigation**: This is a bug fix, new behavior is correct
   - **Mitigation**: Clear CHANGELOG and migration guide

2. **Performance**: Extra state reads during Update
   - **Mitigation**: State reads are fast (in-memory)
   - **Mitigation**: Only affects Update path, not Create/Read

3. **State Inconsistency**: Merge logic bugs
   - **Mitigation**: Comprehensive unit tests
   - **Mitigation**: Acceptance tests against real NetBox

### Benefits
- ✅ **Fixes critical data loss bug**
- ✅ **Enables partial custom field management**
- ✅ **Preserves backward compatibility** (mostly)
- ✅ **Aligns with Terraform best practices**

## Success Metrics

### Definition of Done
- [ ] All 100+ resources updated to use merge logic
- [ ] Zero acceptance test failures
- [ ] All 4 custom field test scenarios pass for pilot resources
- [ ] No reported data loss incidents
- [ ] Documentation complete

### Performance Metrics
- Update operations with custom fields: < 5% slower than before
- No additional API calls during Update (use state, not API read)
- Memory usage unchanged

## Alternative: Quick Band-Aid Fix

If full implementation takes too long, implement minimal fix:

```go
// In ApplyCustomFields, change behavior during Update:
func ApplyCustomFields[T CustomFieldsSetter](ctx context.Context, request T, customFields types.Set, diags *diag.Diagnostics) {
    if !IsSet(customFields) {
        // BAND-AID: Send empty map instead of omitting field
        // This tells NetBox "no changes to custom fields"
        request.SetCustomFields(map[string]interface{}{})
        return
    }
    // ... rest of function unchanged ...
}
```

**Pros**:
- 1-line fix
- Prevents data loss immediately

**Cons**:
- Doesn't enable partial management
- May have unintended consequences
- Not the right long-term solution

## Conclusion

This plan provides a comprehensive path to fixing the critical custom fields data loss bug while enabling the valuable partial management feature. The batched approach allows for incremental delivery with the most critical resources fixed first.

**Recommended Next Steps**:
1. Create feature branch: `feature/custom-fields-partial-management`
2. Implement Batch 1 (core utilities) with tests
3. Implement Batch 2 (device resource) with comprehensive tests
4. If tests pass, continue with Batches 3-7
5. Release v0.0.13 with emergency fix for critical resources
6. Complete remaining batches for v0.0.14
