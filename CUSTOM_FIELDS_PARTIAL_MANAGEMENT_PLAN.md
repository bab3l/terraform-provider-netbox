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

#### Leverage Existing Helper Functions

We already have excellent helper functions in `internal/utils/`:

**Request Helpers** (`request_helpers.go`):
- `ApplyCustomFields()` - Applies custom fields to API request
- `ApplyCommonFields()` - Applies description, comments, tags, and custom fields
- `CustomFieldsSetter` interface - Generic interface for request types

**State Helpers** (`state_helpers.go`):
- `PopulateCustomFieldsFromAPI()` - Converts API response to Terraform state
- `CustomFieldsFromAPI()` - Legacy version (backward compat)
- `IsSet()` - Checks if value is not null/unknown

**Data Conversion** (`common.go`):
- `CustomFieldModelsToMap()` - Converts Terraform models to API map
- `CustomFieldsToMap()` - Core conversion logic
- `MapToCustomFieldModels()` - Converts API map to Terraform models (preserves types)
- `BuildCustomFieldModelsFromAPI()` - Type inference for import scenarios

#### Phase 1: Enhance Existing Helpers (Critical Bug Fix)

**Step 1**: Add merge-aware function alongside existing `ApplyCustomFields`:

```go
// ApplyCustomFieldsWithMerge handles custom fields during Update operations with merge logic.
// Existing ApplyCustomFields remains for Create operations (no merge needed).
//
// Behavior:
//   - If plan has custom_fields: Merge plan + state, send merged map
//   - If plan is null/unknown: Send state values to preserve them
//   - If both null: Send empty map (safe default)
func ApplyCustomFieldsWithMerge[T CustomFieldsSetter](
    ctx context.Context,
    request T,
    planCustomFields types.Set,
    stateCustomFields types.Set,
    diags *diag.Diagnostics,
) {
    // User specified custom_fields in config - merge with existing state
    if IsSet(planCustomFields) {
        tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Merging plan with state")
        merged := MergeCustomFieldSets(ctx, planCustomFields, stateCustomFields, diags)
        if diags.HasError() {
            return
        }
        request.SetCustomFields(merged)
        return
    }

    // User omitted custom_fields - preserve ALL existing state values
    if IsSet(stateCustomFields) {
        tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Preserving state (plan is null)")
        // Reuse existing conversion helper
        var stateModels []CustomFieldModel
        cfDiags := stateCustomFields.ElementsAs(ctx, &stateModels, false)
        diags.Append(cfDiags...)
        if diags.HasError() {
            return
        }
        // Reuse existing helper to convert models to map
        stateMap := CustomFieldModelsToMap(stateModels)
        request.SetCustomFields(stateMap)
        return
    }

    // Both null - send empty map
    tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Both plan and state null, sending empty map")
    request.SetCustomFields(map[string]interface{}{})
}

// MergeCustomFieldSets merges plan custom fields with state custom fields.
// Uses existing conversion helpers to maintain consistency.
func MergeCustomFieldSets(
    ctx context.Context,
    plan types.Set,
    state types.Set,
    diags *diag.Diagnostics,
) map[string]interface{} {
    result := make(map[string]interface{})

    // Start with ALL state custom fields (preserves unmanaged fields)
    if IsSet(state) {
        var stateModels []CustomFieldModel
        stateDiags := state.ElementsAs(ctx, &stateModels, false)
        diags.Append(stateDiags...)
        if !diags.HasError() {
            // Reuse existing helper
            stateMap := CustomFieldModelsToMap(stateModels)
            for k, v := range stateMap {
                result[k] = v
            }
        }
    }

    // Overlay plan custom fields (these override state values)
    if IsSet(plan) {
        var planModels []CustomFieldModel
        planDiags := plan.ElementsAs(ctx, &planModels, false)
        diags.Append(planDiags...)
        if !diags.HasError() {
            // Process each plan field
            for _, cf := range planModels {
                name := cf.Name.ValueString()
                value := cf.Value.ValueString()

                // Empty or null value = remove this field
                if cf.Value.IsNull() || value == "" {
                    delete(result, name)
                    continue
                }

                // Convert value based on type (reuse existing logic pattern)
                cfType := cf.Type.ValueString()
                switch cfType {
                case "integer":
                    if intVal, err := strconv.Atoi(value); err == nil {
                        result[name] = intVal
                    } else {
                        result[name] = value // Fallback to string
                    }
                case "boolean":
                    result[name] = (value == "true")
                default:
                    result[name] = value
                }
            }
        }
    }

    return result
}
```

**Key Design Decision**: Keep existing `ApplyCustomFields()` unchanged for Create operations, add new `ApplyCustomFieldsWithMerge()` for Update operations. This minimizes risk and maintains backward compatibility.

#### Phase 2: Update Resource Update() Methods

The pattern for updating resources becomes very simple - just switch from `ApplyCustomFields` to `ApplyCustomFieldsWithMerge`:

**Before (broken)**:
```go
func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data DeviceResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    // ... build request ...

    // OLD: Doesn't preserve unmanaged fields
    utils.ApplyCustomFields(ctx, &request, data.CustomFields, &resp.Diagnostics)

    // ... execute update ...
}
```

**After (fixed)**:
```go
func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var state, plan DeviceResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)  // ADD: Read state
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    // ... build request ...

    // NEW: Preserves unmanaged fields by merging plan + state
    utils.ApplyCustomFieldsWithMerge(ctx, &request, plan.CustomFields, state.CustomFields, &resp.Diagnostics)

    // ... execute update ...
}
```

**Key Changes**:
1. Read both `state` and `plan` (not just plan)
2. Use `ApplyCustomFieldsWithMerge` instead of `ApplyCustomFields`
3. Pass both `plan.CustomFields` and `state.CustomFields`

**Resources Using `ApplyCommonFields`**:

Some resources use the composite helper `ApplyCommonFields` which internally calls `ApplyCustomFields`. We'll need a merge-aware version:

```go
// Add new composite helper to request_helpers.go
func ApplyCommonFieldsWithMerge[T FullCommonFieldsSetter](
    ctx context.Context,
    request T,
    description, comments types.String,
    planTags, stateTags types.Set,
    planCustomFields, stateCustomFields types.Set,
    diags *diag.Diagnostics,
) {
    ApplyDescriptiveFields(request, description, comments)
    ApplyTags(ctx, request, planTags, diags)  // Tags don't need merge (different semantics)
    if diags.HasError() {
        return
    }
    ApplyCustomFieldsWithMerge(ctx, request, planCustomFields, stateCustomFields, diags)
}
```

Then in resources:
```go
// OLD:
utils.ApplyCommonFields(ctx, &request, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)

// NEW:
utils.ApplyCommonFieldsWithMerge(ctx, &request, plan.Description, plan.Comments, plan.Tags, state.Tags, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
```

**Note**: Tags typically use "replace all" semantics (unlike custom fields), so we don't merge them. Users specify all tags they want.

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
**Files**: 2 files
**Estimated Time**: 2-3 hours

#### Files to Modify:
1. `internal/utils/request_helpers.go`:
   - Add `ApplyCustomFieldsWithMerge()` function
   - Add `MergeCustomFieldSets()` helper
   - Add `ApplyCommonFieldsWithMerge()` composite helper
   - Keep existing functions unchanged (backward compat)

2. `internal/utils/request_helpers_test.go` (create if doesn't exist):
   - Unit tests for `MergeCustomFieldSets()`
   - Test plan + state merge scenarios
   - Test empty value removal
   - Test null plan preserves state
   - Test type conversions (integer, boolean, text)

#### Implementation Details:

**request_helpers.go additions**:
```go
// Add after existing ApplyCustomFields function:

// ApplyCustomFieldsWithMerge handles custom fields during Update operations.
// Leverages existing helper functions (CustomFieldModelsToMap, etc.)
// [Full implementation from Phase 1 section above]
func ApplyCustomFieldsWithMerge[T CustomFieldsSetter](...) { /* ... */ }

// MergeCustomFieldSets merges plan with state custom fields.
// [Full implementation from Phase 1 section above]
func MergeCustomFieldSets(...) map[string]interface{} { /* ... */ }

// ApplyCommonFieldsWithMerge - merge-aware version of ApplyCommonFields
func ApplyCommonFieldsWithMerge[T FullCommonFieldsSetter](...) { /* ... */ }
```

**Success Criteria**:
- [ ] `ApplyCustomFieldsWithMerge()` function added
- [ ] `MergeCustomFieldSets()` correctly merges plan + state
- [ ] Empty value removes field from result
- [ ] Null plan preserves all state fields
- [ ] Reuses existing helpers (CustomFieldModelsToMap, etc.)
- [ ] Unit tests pass with 100% coverage of merge scenarios
- [ ] No changes to existing functions (backward compatible)

### Batch 2: Device Resource (Pilot Implementation)
**Priority**: CRITICAL
**Files**: 2 files
**Estimated Time**: 3-4 hours

#### Files to Modify:
1. `internal/resources/device_resource.go`:
   - Update `Update()` method to read both state and plan
   - Replace `ApplyCustomFields` with `ApplyCustomFieldsWithMerge`
   - Pass both `plan.CustomFields` and `state.CustomFields`

2. `internal/resources_acceptance_tests/device_custom_fields_preservation_test.go`:
   - Already created with comprehensive test scenarios
   - Fix any schema issues (object_types vs content_types)
   - Add 3 additional test functions

#### Implementation Pattern:

**device_resource.go Update() method**:
```go
func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Read BOTH state and plan (critical for merge)
    var state, plan DeviceResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // ... parse ID, lookup references, build request ...

    // Apply common fields with merge-aware helper
    utils.ApplyCommonFieldsWithMerge(
        ctx, &deviceRequest,
        plan.Description, plan.Comments,
        plan.Tags, state.Tags,  // Tags use replace-all semantics
        plan.CustomFields, state.CustomFields,  // Custom fields use merge semantics
        &resp.Diagnostics,
    )

    // ... execute update, map response to state ...
}
```

**Alternative if resource doesn't use ApplyCommonFields**:
```go
// Apply fields individually
utils.ApplyDescription(&deviceRequest, plan.Description)
utils.ApplyComments(&deviceRequest, plan.Comments)
utils.ApplyTags(ctx, &deviceRequest, plan.Tags, &resp.Diagnostics)
utils.ApplyCustomFieldsWithMerge(ctx, &deviceRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
```

#### Test Scenarios:

**Test 1: Custom Fields Preservation** (already created):
- Step 1: Create with custom_fields = [field_a, field_b]
- Step 2: Update description WITHOUT custom_fields in config
- Expected: field_a and field_b still exist in NetBox
- Current: FAILS - fields are deleted

**Test 2: Partial Management** (add):
```go
func TestAccDeviceResource_CustomFieldsPartialManagement(t *testing.T) {
    // Step 1: Create with custom_fields = [env="prod", owner="team-a"]
    // Step 2: Update with custom_fields = [env="staging"] (removed owner from config)
    // Expected: env="staging", owner="team-a" (preserved)
    // Step 3: Update with custom_fields = [env="staging", cost_center="123"]
    // Expected: env="staging", owner="team-a", cost_center="123"
}
```

**Test 3: Explicit Removal** (add):
```go
func TestAccDeviceResource_CustomFieldsExplicitRemoval(t *testing.T) {
    // Step 1: Create with custom_fields = [field_a="value1", field_b="value2"]
    // Step 2: Update with custom_fields = [{ name="field_a", value="" }]
    // Expected: field_a removed, field_b preserved
}
```

**Test 4: Complete Removal** (add):
```go
func TestAccDeviceResource_CustomFieldsCompleteRemoval(t *testing.T) {
    // Step 1: Create with custom_fields = [field_a, field_b]
    // Step 2: Update with custom_fields = []
    // Expected: All custom fields removed
}
```

**Success Criteria**:
- [ ] Device resource Update() uses `ApplyCustomFieldsWithMerge()`
- [ ] State is read in Update() for merge context
- [ ] All 4 acceptance tests pass against live NetBox
- [ ] No data loss when updating other fields
- [ ] Partial management works as expected
- [ ] Explicit removal works
- [ ] Build succeeds with no errors

### Batch 3: High-Priority Resources (IPAM Core)
**Priority**: HIGH
**Files**: 15-20 files
**Estimated Time**: 4-6 hours

Update resources that commonly use custom fields in production:

#### Resources to Update:
- `ip_address_resource.go`
- `prefix_resource.go`
- `vlan_resource.go`
- `vrf_resource.go`
- `aggregate_resource.go`
- `asn_resource.go`
- `asn_range_resource.go`
- `ip_range_resource.go`
- `l2vpn_resource.go`
- `vlan_group_resource.go`
- `rir_resource.go`
- `circuit_resource.go`
- `circuit_termination_resource.go`
- `provider_resource.go`

#### Implementation Pattern:

**Step 1**: Check how resource currently applies custom fields:

**Pattern A: Uses `ApplyCommonFields`**:
```go
// OLD:
utils.ApplyCommonFields(ctx, &request, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)

// NEW:
var state, plan IPAddressResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

utils.ApplyCommonFieldsWithMerge(ctx, &request,
    plan.Description, plan.Comments,
    plan.Tags, state.Tags,
    plan.CustomFields, state.CustomFields,
    &resp.Diagnostics)
```

**Pattern B: Uses `ApplyCustomFields` directly**:
```go
// OLD:
utils.ApplyCustomFields(ctx, &request, data.CustomFields, &resp.Diagnostics)

// NEW:
var state, plan PrefixResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

utils.ApplyCustomFieldsWithMerge(ctx, &request, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
```

**Step 2**: Update the Update() method:
1. Change `var data ModelType` to `var state, plan ModelType`
2. Add `req.State.Get(ctx, &state)...`
3. Change `req.Plan.Get(ctx, &data)...` to `req.Plan.Get(ctx, &plan)...`
4. Update all references from `data.Field` to `plan.Field`
5. Replace `ApplyCustomFields` with `ApplyCustomFieldsWithMerge`
6. Pass both `plan.CustomFields` and `state.CustomFields`

**Step 3**: Add basic acceptance test for each resource:
```go
// TestAcc<Resource>Resource_CustomFieldsPreservation
// Minimal test: Create with custom fields, update without, verify preserved
```

#### Automation Helper Script:

Create PowerShell script to identify which pattern each resource uses:

```powershell
# scripts/identify-custom-field-patterns.ps1
$resources = Get-ChildItem "internal/resources/*_resource.go"
foreach ($file in $resources) {
    $content = Get-Content $file -Raw
    if ($content -match "ApplyCommonFields") {
        Write-Host "$($file.Name): Uses ApplyCommonFields (Pattern A)"
    } elseif ($content -match "ApplyCustomFields") {
        Write-Host "$($file.Name): Uses ApplyCustomFields (Pattern B)"
    } else {
        Write-Host "$($file.Name): No custom fields helper"
    }
}
```

**Success Criteria**:
- [ ] All IPAM core resources use merge-aware helpers
- [ ] State is read in Update() for all resources
- [ ] At least 1 acceptance test per resource verifies preservation
- [ ] grep -r "ApplyCustomFields\[" shows zero matches in updated files (all use WithMerge version)
- [ ] All existing acceptance tests still pass
- [ ] Build succeeds with no errors

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

### Batch 8: Update Examples
**Priority**: HIGH
**Files**: 20-30 example files
**Estimated Time**: 3-4 hours

#### Goals:
Update example configurations to demonstrate partial custom field management and document best practices.

#### Files to Update:

**examples/resources/** - One example per resource type:
- `netbox_device/resource.tf` - Show partial CF management
- `netbox_ip_address/resource.tf` - Show CF preservation
- `netbox_prefix/resource.tf` - Show explicit removal
- `netbox_vlan/resource.tf`
- `netbox_vm/resource.tf`
- And 15+ other resource examples

#### Example Patterns to Add:

**Pattern 1: Partial Management (Recommended)**:
```hcl
# examples/resources/netbox_device/resource.tf
resource "netbox_device" "example" {
  name        = "router-01"
  device_type = netbox_device_type.example.id
  role        = netbox_device_role.example.id
  site        = netbox_site.example.id

  # Manage only specific custom fields
  # Other custom fields set in NetBox UI are preserved
  custom_fields = [
    {
      name  = "environment"
      type  = "text"
      value = "production"
    },
    {
      name  = "managed_by_terraform"
      type  = "boolean"
      value = "true"
    }
  ]
}
```

**Pattern 2: External Management**:
```hcl
# examples/resources/netbox_device/external_custom_fields.tf
# Custom fields are managed outside Terraform (e.g., NetBox UI, automation)
# Terraform manages only the device configuration
resource "netbox_device" "example" {
  name        = "router-02"
  device_type = netbox_device_type.example.id
  role        = netbox_device_role.example.id
  site        = netbox_site.example.id

  # custom_fields intentionally omitted
  # All custom fields set externally are preserved
}
```

**Pattern 3: Explicit Removal**:
```hcl
# examples/resources/netbox_device/remove_custom_fields.tf
resource "netbox_device" "example" {
  name        = "router-03"
  device_type = netbox_device_type.example.id
  role        = netbox_device_role.example.id
  site        = netbox_site.example.id

  # Remove specific custom field by setting value to empty
  custom_fields = [
    {
      name  = "old_field"
      type  = "text"
      value = ""  # This removes the field
    }
  ]
}
```

**Pattern 4: Complete Removal**:
```hcl
# examples/resources/netbox_device/clear_all_custom_fields.tf
resource "netbox_device" "example" {
  name        = "router-04"
  device_type = netbox_device_type.example.id
  role        = netbox_device_role.example.id
  site        = netbox_site.example.id

  # Empty list removes ALL custom fields
  custom_fields = []
}
```

#### New Example Files to Create:

**examples/guides/custom_fields_management.md**:
```markdown
# Custom Fields Management Guide

## Overview
Terraform provider for NetBox supports partial custom field management,
allowing you to manage some fields while preserving others set externally.

## Behavior
- Fields in config: Managed by Terraform
- Fields NOT in config: Preserved (not touched)
- Empty value: Explicitly removes that field
- Empty list: Removes ALL fields

## Use Cases
[detailed examples...]
```

**Success Criteria**:
- [ ] All resource examples updated with custom_fields patterns
- [ ] New guide document created
- [ ] Examples demonstrate all 4 patterns
- [ ] Code comments explain partial management behavior
- [ ] Examples use realistic field names and values

### Batch 9: Regenerate Documentation
**Priority**: HIGH
**Files**: Generated docs
**Estimated Time**: 1-2 hours

#### Documentation Generation:

**Step 1**: Run documentation generator:
```bash
# Uses terraform-plugin-docs tool
make docs
# or
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate
```

**Step 2**: Verify generated documentation:

**docs/resources/** - Check all resource docs:
- Custom fields schema is documented
- Examples are included from examples/
- Behavior is clearly explained

**docs/guides/** - New guide should appear:
- `custom_fields_management.md` rendered correctly

**Step 3**: Manual updates to docs/index.md:

Add section about custom fields to provider documentation:
```markdown
## Custom Fields Management

The NetBox provider supports partial custom field management:

- **Partial Management**: Specify only the custom fields you want to manage
- **Preservation**: Fields not in your config are preserved
- **External Management**: Omit custom_fields to preserve all external values

See the [Custom Fields Management Guide](guides/custom_fields_management) for details.
```

**Step 4**: Update CHANGELOG.md:

Add comprehensive entry:
```markdown
## v0.0.13 (2026-01-XX)

### Features

#### Custom Fields Partial Management
- **BREAKING**: Custom fields now use merge semantics instead of replace-all
- Fields specified in config are managed by Terraform
- Fields NOT in config are preserved (not deleted)
- Empty value explicitly removes a field
- Empty list removes all fields

**Migration Impact**:
- **Low Impact**: This is a bug fix for critical data loss issue
- **No Config Changes Required**: Existing configs work better after fix
- **Benefit**: No more data loss when updating resources

#### Before (v0.0.12 and earlier - BROKEN):
```hcl
resource "netbox_device" "server" {
  name = "server-01"
  # custom_fields omitted
  description = "update description"
}
# ❌ BUG: Updates DELETE all custom fields in NetBox
```

#### After (v0.0.13 - FIXED):
```hcl
resource "netbox_device" "server" {
  name = "server-01"
  # custom_fields omitted
  description = "update description"
}
# ✅ Updates preserve all custom fields in NetBox
```

See migration guide for details.

### Bug Fixes
- Fixed critical data loss bug where updating resources without custom_fields
  in config would delete ALL custom fields in NetBox
- Custom fields are now properly preserved during updates
```

**Success Criteria**:
- [ ] `make docs` runs successfully
- [ ] All resource docs regenerated with latest examples
- [ ] Guide documentation rendered correctly
- [ ] CHANGELOG.md updated with comprehensive entry
- [ ] Provider index.md mentions custom fields feature
- [ ] No broken links in documentation

### Batch 10: Comprehensive Testing
**Priority**: CRITICAL
**Files**: Test files
**Estimated Time**: 4-6 hours

1. Run full acceptance test suite
2. Test import scenarios with custom fields
3. Test state migration edge cases
4. Performance testing with many custom fields
5. Test all 4 custom field patterns from examples
6. Verify documentation examples work

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
1. ✅ Create feature branch: `feature/custom-fields-partial-management` (DONE)
2. Implement Batch 1 (core utilities) with tests
3. Implement Batch 2 (device resource) with comprehensive tests
4. If tests pass, continue with Batch 3 (IPAM resources)
5. Update examples (Batch 8) and regenerate documentation (Batch 9)
6. Release v0.0.13 with emergency fix for critical resources
7. Complete remaining batches for v0.0.14

## Release Strategy

### Emergency Release: v0.0.13 (2-3 days)
**Critical Bug Fix - Data Loss Prevention**

#### Scope:
- **Batch 1**: Core utilities (2-3 hours)
- **Batch 2**: Device resource pilot (3-4 hours)
- **Batch 3**: IPAM core resources (4-6 hours)
- **Batch 8**: Update critical examples (2 hours for device, ip_address, prefix only)
- **Batch 9**: Documentation generation (1-2 hours)

**Total Time**: 12-18 hours (2-3 days with testing)

#### Critical Resources Fixed:
- netbox_device (most commonly used)
- netbox_ip_address (high custom field usage)
- netbox_prefix (high custom field usage)
- netbox_vlan
- netbox_vrf
- netbox_virtual_machine
- netbox_cluster

#### Examples Updated (Emergency):
- `examples/resources/netbox_device/resource.tf` - All 4 patterns
- `examples/resources/netbox_ip_address/resource.tf` - Preservation pattern
- `examples/resources/netbox_prefix/resource.tf` - Explicit removal pattern

#### Documentation (Emergency):
- Update critical examples (device, ip_address, prefix)
- Regenerate docs with `make docs` / `tfplugindocs generate`
- Update CHANGELOG.md with bug fix and breaking change notice
- Create `examples/guides/custom_fields_management.md` with quick start guide

#### Testing:
- Run acceptance tests for fixed resources
- Verify custom field preservation in manual tests
- Test all 4 patterns (partial, external, explicit removal, complete removal)
- Ensure no regression in existing tests

**Release Criteria**:
- All Batch 1-3 tests passing
- No regression in existing tests
- Critical examples updated and working
- Documentation generated and reviewed
- CHANGELOG.md updated
- Migration guide in place

---

### Complete Release: v0.0.14 (1-2 weeks)
**Full Rollout - All Resources + Complete Documentation**

#### Scope:
- **Batch 4**: DCIM remaining (4-5 hours)
- **Batch 5**: Virtualization (2-3 hours)
- **Batch 6**: Circuits, Tenancy, Users (2-3 hours)
- **Batch 7**: Extras, Wireless, VPN (3-4 hours)
- **Batch 8**: Complete all remaining examples (2-3 hours additional)
- **Batch 10**: Comprehensive testing (4-6 hours)

**Total Additional Time**: 18-25 hours

#### All Resources Fixed:
- Remaining 90+ resources with custom fields support
- All data sources updated (if applicable)
- All examples updated
- Full documentation suite

#### Examples Updated (Complete):
- **All resource examples** in `examples/resources/` (30-40 files)
- Each example shows appropriate custom field pattern
- Comments explain behavior clearly
- Realistic field names and values

#### Documentation (Complete):
- **All resource docs** regenerated in `docs/resources/`
- **Complete guide**: `examples/guides/custom_fields_management.md`
  * Overview of partial management feature
  * All 4 patterns with detailed examples
  * Migration guide for v0.0.12 users
  * Best practices and common scenarios
  * Troubleshooting section
- **Provider docs** updated in `docs/index.md` with custom fields section
- **CHANGELOG.md** comprehensive entry with before/after examples

#### Testing:
- Full acceptance test suite (120m timeout)
- Import scenarios with custom fields
- State migration edge cases
- Performance testing with many custom fields
- All 4 custom field patterns validated for all resources
- Documentation examples validated (actually run in test environment)

**Release Criteria**:
- All tests passing (100+ resources)
- All examples working and validated
- Complete documentation generated and reviewed
- Performance benchmarks met (< 5% slower on Update)
- No broken links in documentation
- CHANGELOG.md comprehensive
- Migration guide complete
