# Request Standardization Implementation Plan

**Status**: 🚀 READY TO START - Detailed batches planned
**Phase**: Request Field Assignment Standardization
**Total Resources to Migrate**: ~85 resources
**Estimated Effort**: ~150 hours over 2-3 months
**Expected Impact**: Improved maintainability, reduced bugs, consistent error handling

---

## Phase 1: Foundation (PREREQUISITE - Do This First!)

### Phase 1 Goals
- Create comprehensive helper function suite
- Establish standardized patterns
- Build foundation for all migration work

### Phase 1 Tasks

#### 1.1: Create Pointer Helpers (2 hours)
```go
// Add to internal/utils/request_helpers.go

func StringPtr(v types.String) *string {
    if v.IsNull() || v.IsUnknown() {
        return nil
    }
    val := v.ValueString()
    return &val
}

func IntPtr(v types.Int64) *int32 {
    if v.IsNull() || v.IsUnknown() {
        return nil
    }
    val := int32(v.ValueInt64())
    return &val
}

func BoolPtr(v types.Bool) *bool {
    if v.IsNull() || v.IsUnknown() {
        return nil
    }
    val := v.ValueBool()
    return &val
}

func Float64Ptr(v types.String) *float64 {
    if v.IsNull() || v.IsUnknown() {
        return nil
    }
    floatVal, err := strconv.ParseFloat(v.ValueString(), 64)
    if err != nil {
        return nil
    }
    return &floatVal
}
```

#### 1.2: Expand ApplyDescription Helper (1 hour)
Ensure it handles:
- Null checking
- Unknown checking
- Pointer conversion
- Different request types

#### 1.3: Create Enum Conversion Helpers (3 hours)
```go
// For each enum type used in resources:
func ToDeviceStatus(v types.String) *netbox.WritableDeviceRequestStatus { ... }
func ToDeviceRoleColor(v types.String) *string { ... }
func ToCableType(v types.String) *netbox.WritableDeviceConnectorFormatValue { ... }
// etc.
```

#### 1.4: Create Reference Lookup Helpers (2 hours)
```go
func SetRequiredReference(
    ctx context.Context,
    client *netbox.APIClient,
    refValue types.String,
    lookupFunc func(context.Context, *netbox.APIClient, string) (*BriefRequest, diag.Diagnostics),
    setter func(*BriefRequest),
    diags *diag.Diagnostics,
) bool { ... }
```

#### 1.5: Create Optional Field Helper (1 hour)
```go
func ApplyOptionalField(value types.String, setter func(string) error) error {
    if value.IsNull() || value.IsUnknown() {
        return nil
    }
    return setter(value.ValueString())
}
```

**Phase 1 Subtotal: ~9 hours**
**Phase 1 Resources Needed: None (foundation only)**

---

## Phase 2: High-Priority Port Resources

### Rationale
- **15 resources** all using same non-standard `apiReq.SetField()` pattern
- All are port-related (high impact, focused domain)
- Easy to verify (consistent pattern)
- Unblocks teams working with interfaces/ports

### Phase 2A: Console Port Resources (3 resources)
**Effort**: ~6 hours

| Resource | Current Pattern | Target | Notes |
|----------|-----------------|--------|-------|
| console_port_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers | Setter for: Label, Type, Speed, Description, MarkConnected |
| console_port_template_resource.go | `apiReq.SetDeviceType()` | Helper pattern | Setter for: DeviceType, ModuleType, Label, Type |
| console_server_port_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers | Setter for: Label, Type, Speed, Description, MarkConnected |
| console_server_port_template_resource.go | `apiReq.SetType()` | Helper pattern | Setter for: DeviceType, ModuleType, Label, Type |

### Phase 2B: Front/Rear Port Resources (4 resources)
**Effort**: ~8 hours

| Resource | Current Pattern | Target |
|----------|-----------------|--------|
| front_port_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers |
| front_port_template_resource.go | `apiReq.SetDeviceType()` | Helper pattern |
| rear_port_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers |
| rear_port_template_resource.go | Mixed | Helper pattern |

### Phase 2C: Power Port Resources (4 resources)
**Effort**: ~8 hours

| Resource | Current Pattern | Target |
|----------|-----------------|--------|
| power_port_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers |
| power_port_template_resource.go | Direct field | Standardize to helpers |
| power_outlet_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers |
| power_outlet_template_resource.go | Direct field | Standardize to helpers |

### Phase 2D: Inventory Item Resource (1 resource)
**Effort**: ~2 hours

| Resource | Current Pattern | Target |
|----------|-----------------|--------|
| inventory_item_resource.go | `apiReq.SetDescription()` | `utils.ApplyDescription()` + helpers |

**Phase 2 Subtotal: ~24 hours**
**Phase 2 Total Resources: 15 resources**

---

## Phase 3: Template Resources Standardization

### Rationale
- 8 remaining template resources with direct field assignments
- Currently inconsistent with Phase 2 refactored templates
- Lower complexity than device resources

### Phase 3A: Port/Interface Templates (Already in Phase 2)
**Status**: Handled in Phase 2

### Phase 3B: Device/Rack/Cluster Templates (5 resources)
**Effort**: ~10 hours

| Resource | Current Pattern | Target | Notes |
|----------|-----------------|--------|-------|
| device_bay_template_resource.go | Direct field | Helper pattern | Description only |
| interface_template_resource.go | Direct field | Helper pattern | Description only |
| module_bay_template_resource.go | Direct field | Helper pattern | Description only |
| rear_port_template_resource.go | Direct field | Helper pattern | Description only (in Phase 2) |
| inventory_item_template_resource.go | Direct field | Helper pattern | Description only |

### Phase 3C: Power/Console Templates (5 resources)
**Effort**: ~10 hours

| Resource | Current Pattern | Target |
|----------|-----------------|--------|
| power_port_template_resource.go | Direct field | Helper pattern |
| power_outlet_template_resource.go | Direct field | Helper pattern |
| console_port_template_resource.go | Setter (Phase 2) | ✅ Done in Phase 2 |
| console_server_port_template_resource.go | Setter (Phase 2) | ✅ Done in Phase 2 |
| config_template_resource.go | Direct field | Helper pattern |

**Phase 3 Subtotal: ~20 hours**
**Phase 3 Total Resources: 10 resources**

---

## Phase 4: Special Cases & Unique Patterns

### Rationale
- Resources with non-standard patterns needing custom solutions
- Lower volume but higher complexity

### Phase 4A: config_context_resource (CRITICAL - UNIQUE)
**Effort**: ~4 hours
**Status**: 🔴 BLOCKER - Most unique pattern

| Resource | Current Pattern | Issue | Solution |
|----------|-----------------|-------|----------|
| config_context_resource.go | `request.Tags = setToStringSlice()` | ONLY resource using custom Tag conversion | Create `ApplyConfigContextTags()` helper |

**Why This Matters:**
- ConfigContextRequest uses `[]string` for Tags (not `[]NestedTagRequest`)
- Currently uses unique `setToStringSlice()` function
- Solution: Create specialized helper for this request type

```go
func ApplyConfigContextTags(ctx context.Context,
    request *netbox.ConfigContextRequest,
    tags types.Set,
    diags *diag.Diagnostics) {
    if tags.IsNull() || tags.IsUnknown() {
        return
    }
    // Convert tags to []string format for ConfigContext
    tagStrings := setToStringSlice(ctx, tags)
    request.Tags = tagStrings
}
```

### Phase 4B: Assignment Resources (5 resources)
**Effort**: ~10 hours
**Pattern**: AdditionalProperties workaround for references

| Resource | Current Pattern | Issue | Solution |
|----------|-----------------|-------|----------|
| circuit_group_assignment_resource.go | AdditionalProperties | Reference not as direct field | Create reference helper |
| contact_assignment_resource.go | AdditionalProperties | Reference not as direct field | Create reference helper |
| fhrp_group_assignment_resource.go | AdditionalProperties | Reference not as direct field | Create reference helper |
| tunnel_termination_resource.go | AdditionalProperties | Reference not as direct field | Create reference helper |
| tunnel_group_resource.go | AdditionalProperties | Reference not as direct field | Create reference helper |

**Solution:**
```go
func SetAdditionalPropertyReference(
    request interface{},
    key string,
    id int32) {
    addlProps := request.AdditionalProperties
    if addlProps == nil {
        addlProps = make(map[string]interface{})
        request.AdditionalProperties = addlProps
    }
    addlProps[key] = int(id)
}
```

### Phase 4C: Custom Link Resource (1 resource)
**Effort**: ~2 hours
**Pattern**: All direct field assignments

| Resource | Current Pattern | Target |
|----------|-----------------|--------|
| custom_link_resource.go | Direct: Enabled, Weight, GroupName, ButtonClass, NewWindow | Standardize with helpers |

### Phase 4D: Interface Resource (1 resource)
**Effort**: ~3 hours
**Pattern**: Mixed patterns within same resource

| Resource | Current Pattern | Issue |
|----------|-----------------|-------|
| interface_resource.go | Mixed direct + some references | Inconsistent within resource |

**Phase 4 Subtotal: ~20 hours**
**Phase 4 Total Resources: 8 resources**

---

## Phase 5: Core Device & Circuit Resources

### Rationale
- High-impact resources (widely used)
- Moderate complexity (multiple fields)
- Medium volume (10-15 resources)

### Phase 5A: Device Resources (3 resources)
**Effort**: ~12 hours

| Resource | Current Pattern | Fields | Notes |
|----------|-----------------|--------|-------|
| device_resource.go | Direct field | Serial, Face, Status, Airflow, etc. | Heavy usage, many fields |
| device_type_resource.go | Direct field | PartNumber, UHeight, ExcludeFromUtilization, etc. | Enum + numeric fields |
| device_role_resource.go | Direct field | Color, VmRole | Simple but used widely |

### Phase 5B: Circuit Resources (4 resources)
**Effort**: ~12 hours

| Resource | Current Pattern | Fields |
|----------|-----------------|--------|
| circuit_resource.go | Direct field | Description, Comments, Tags, CustomFields |
| circuit_type_resource.go | Direct field | Description, Tags, CustomFields |
| circuit_group_resource.go | Mixed (`netbox.PtrString()`) | Description, Tags, CustomFields |
| circuit_termination_resource.go | Direct field | Description, Tags, CustomFields |

### Phase 5C: Cable & Related (3 resources)
**Effort**: ~9 hours

| Resource | Current Pattern | Fields |
|----------|-----------------|--------|
| cable_resource.go | Direct field | Type, Status, Label, Color, LengthUnit |
| cluster_resource.go | Direct field | Status, Description, Comments, Tags, CustomFields |
| event_rule_resource.go | Direct field | Enabled, ActionType, Description, Tags, CustomFields |

**Phase 5 Subtotal: ~33 hours**
**Phase 5 Total Resources: 10 resources**

---

## Phase 6: Medium Priority Direct Assignment Resources

### Rationale
- Large volume (40+ resources)
- Similar patterns (mostly direct assignment)
- Can be spread over time

### Phase 6A: Location & Site Hierarchy (4 resources)
**Effort**: ~8 hours

| Resource | Current Pattern | Key Fields |
|----------|-----------------|-----------|
| location_resource.go | Direct field + mixed | Description, Parent, Status, Facility, Tags, CustomFields |
| site_resource.go | Direct field | Description, Comments, Tags, CustomFields |
| site_group_resource.go | Direct field | Description, Tags, CustomFields |
| region_resource.go | ✅ Already refactored | Already using helpers |

### Phase 6B: Contact Resources (4 resources)
**Effort**: ~8 hours

| Resource | Current Pattern | Key Fields |
|----------|-----------------|-----------|
| contact_resource.go | Direct field | Title, Phone, Email, Address, Link, Description, Comments, Tags |
| contact_group_resource.go | Direct field | Description, Tags, CustomFields |
| contact_role_resource.go | Direct field | Description, Tags, CustomFields |
| contact_assignment_resource.go | Direct field | Priority, Tags, CustomFields |

### Phase 6C: Cluster Resources (5 resources)
**Effort**: ~10 hours

| Resource | Current Pattern |
|----------|-----------------|
| cluster_group_resource.go | Direct field |
| cluster_type_resource.go | Direct field |
| cluster_resource.go | Direct field (Phase 5) |
| virtual_device_context_resource.go | ✅ Already refactored |
| fhrp_group_resource.go | Direct field |

### Phase 6D: VPN Resources (5 resources)
**Effort**: ~10 hours

| Resource | Current Pattern |
|----------|-----------------|
| ipsec_policy_resource.go | Direct field |
| ipsec_profile_resource.go | Direct field |
| ipsec_proposal_resource.go | Direct field |
| ike_policy_resource.go | Direct field |
| ike_proposal_resource.go | Direct field |

### Phase 6E: Routing & Network Resources (5 resources)
**Effort**: ~10 hours

| Resource | Current Pattern |
|----------|-----------------|
| route_target_resource.go | ✅ Already refactored |
| ip_range_resource.go | Direct field |
| l2vpn_resource.go | ✅ Already refactored |
| tunnel_resource.go | ✅ Already refactored |
| wireless_lan_resource.go | ✅ Already refactored |

### Phase 6F: Template & Configuration Resources (5 resources)
**Effort**: ~10 hours

| Resource | Current Pattern |
|----------|-----------------|
| export_template_resource.go | Direct field |
| custom_field_resource.go | Direct field |
| custom_field_choice_set_resource.go | Direct field |
| journal_entry_resource.go | ✅ Already refactored |
| event_rule_resource.go | Direct field (Phase 5) |

### Phase 6G: Miscellaneous (8 resources)
**Effort**: ~16 hours

| Resource | Current Pattern | Notes |
|----------|-----------------|-------|
| asn_range_resource.go | Direct field | Description, Tags, CustomFields |
| device_bay_resource.go | ✅ Already refactored | Already using helpers |
| device_role_resource.go | Direct field (Phase 5) | Already scheduled |
| device_role_resource.go | Already covered | See Phase 5 |
| manufacturer_resource.go | ✅ Already refactored | Using ApplyDescription |
| platform_resource.go | ✅ Already refactored | Using ApplyDescription |
| rack_type_resource.go | Direct field | Description only |
| rir_resource.go | ✅ Already refactored | Using ApplyMetadataFields |

**Phase 6 Subtotal: ~72 hours**
**Phase 6 Total Resources: 40+ resources**

---

## Phase 7: Remaining Direct Assignment Resources

### Rationale
- Lowest priority
- Can be completed incrementally
- Spread over time as touchpoints occur

**Estimated scope**: 10-15 remaining resources
**Estimated effort**: 15-20 hours
**Timeline**: Ongoing (touchpoint-based)

---

## Batch Execution Sequence

### Recommended Order (for maximum efficiency)

| Batch | Phase | Resources | Effort | Priority | Blockers |
|-------|-------|-----------|--------|----------|----------|
| **S1** | Phase 1 | Foundation | 9 hrs | 🔴 CRITICAL | None - DO FIRST |
| **S2** | Phase 4A | config_context | 4 hrs | 🔴 HIGH | S1 must complete |
| **S3** | Phase 2A | Console Ports | 6 hrs | 🟡 HIGH | S1, S2 |
| **S4** | Phase 2B | Front/Rear Ports | 8 hrs | 🟡 HIGH | S1, S2 |
| **S5** | Phase 2C | Power Ports | 8 hrs | 🟡 HIGH | S1, S2 |
| **S6** | Phase 2D | Inventory Item | 2 hrs | 🟡 HIGH | S1, S2 |
| **S7** | Phase 3 | Templates | 20 hrs | 🟡 HIGH | S1-S6 |
| **S8** | Phase 4B | Assignments | 10 hrs | 🟠 MEDIUM | S1, S7 |
| **S9** | Phase 4C-D | Custom Link, Interface | 5 hrs | 🟠 MEDIUM | S1, S7 |
| **S10** | Phase 5A | Device Resources | 12 hrs | 🟠 MEDIUM | S1, S7-S9 |
| **S11** | Phase 5B | Circuit Resources | 12 hrs | 🟠 MEDIUM | S1, S7-S9 |
| **S12** | Phase 5C | Cable & Related | 9 hrs | 🟠 MEDIUM | S1, S7-S9 |
| **S13+** | Phase 6 | Medium Priority | 72 hrs | 🟢 LOW | S1-S12 (can overlap) |
| **S14+** | Phase 7 | Remaining | 15+ hrs | 🟢 LOW | Ongoing (touchpoint) |

---

## Resource Readiness Matrix

### Ready Immediately (Phase 1 Complete)
- ✅ Foundation helpers
- ✅ Pointer utilities
- ✅ ApplyDescription expansion

### Ready After Phase 1+2 (Port Standardization)
- 15 port resources
- 1 inventory_item resource

### Ready After Phase 3 (Templates)
- 8 template resources

### Ready After Phase 4 (Special Cases)
- config_context (custom Tag handling)
- 5 assignment resources
- custom_link, interface resources

### Ready After Phase 5 (Core Resources)
- 10 device/circuit/cable resources

### Ongoing (Phase 6+)
- 40+ medium-priority resources (can overlap with later phases)

---

## Success Criteria

### Per Batch
- ✅ All target resources migrate to helper functions
- ✅ No breaking changes to resource behavior
- ✅ All tests pass (existing + new)
- ✅ Build verification: `go build .` succeeds
- ✅ Code review approval

### Overall
- ✅ 85+ resources standardized
- ✅ Consistent error handling across all resources
- ✅ Maintainability improved
- ✅ Documentation updated
- ✅ 2-3% code reduction achieved

---

## Effort Summary

| Phase | Effort | Resources | Status |
|-------|--------|-----------|--------|
| Phase 1: Foundation | 9 hrs | - | Not Started |
| Phase 2: Port Resources | 24 hrs | 15 | Blocked by Phase 1 |
| Phase 3: Templates | 20 hrs | 10 | Blocked by Phase 1-2 |
| Phase 4: Special Cases | 20 hrs | 8 | Blocked by Phase 1 |
| Phase 5: Core Device | 33 hrs | 10 | Blocked by Phase 1 |
| Phase 6: Medium Priority | 72 hrs | 40+ | Blocked by Phase 1 (can overlap) |
| Phase 7: Remaining | 15+ hrs | 10-15 | Ongoing (touchpoint) |
| **TOTAL** | **~190 hours** | **~100 resources** | **Ready to plan** |

---

## Risk Mitigation

### Risk: Breaking Changes in go-netbox
**Likelihood:** Low
**Mitigation:** Comprehensive tests after each batch, staged rollout

### Risk: Increased Helper Complexity
**Likelihood:** Medium
**Mitigation:** Keep helpers focused on single patterns, document thoroughly

### Risk: Merge Conflicts
**Likelihood:** Medium-High (if working in parallel)
**Mitigation:** Coordinate batches, use feature branches, regular syncs

### Risk: Developer Learning Curve
**Likelihood:** Low
**Mitigation:** Create clear documentation, code examples, PR templates

---

## Documentation & Tracking

### Files to Update
- `CONTRIBUTING.md` - Add standardization guidelines
- `INTERFACE_HELPERS_PLAN.md` - Add standardization section
- `.githooks/` - Add linting rules (future)
- PR templates - Add standardization checklist

### Tracking Updates
- Update status after each batch
- Maintain resource inventory
- Track effort vs. actual
- Document any blockers
