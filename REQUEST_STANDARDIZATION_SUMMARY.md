# Quick Summary: Request Standardization Issues

## TL;DR

**Problem:** Three different patterns exist for assigning fields to request objects across 102 resources.

**Impact:** Code maintainability, consistency, error handling

---

## The Three Patterns

### Pattern 1: Direct Field Assignment (70% of resources)
```go
request.Description = &desc
request.Tags = tags
```
**Examples:** config_context, custom_link, device_type, location

### Pattern 2: Setter Methods (15% of resources)
```go
apiReq.SetDescription(value)
apiReq.SetTags(tags)
```
**Examples:** console_port, front_port, inventory_item

### Pattern 3: Helper Functions (10% of resources - NEWEST)
```go
utils.ApplyDescription(request, data.Description)
utils.ApplyMetadataFields(ctx, request, data.Tags, data.CustomFields, diags)
```
**Examples:** All Batch 1-13 refactored resources

---

## Additional Standardization Issues

### Issue 1: Pointer Handling (3 variants)
- `&variable` (most common)
- `utils.StringPtr(value)` (helper)
- `netbox.PtrString(value)` (go-netbox)

### Issue 2: Optional Field Checking
- `if !value.IsNull() { ... }`
- `if !value.IsNull() && !value.IsUnknown() { ... }`
- Helper handles automatically

### Issue 3: Enum Conversion (varies per field)
- String passthrough
- `netbox.ValueType(string)` conversion
- No conversion needed

### Issue 4: Reference Lookups (3 approaches)
- Standard lookup with error checking
- Direct assignment
- AdditionalProperties workaround (config_context only!)

### Issue 5: Request Construction (3 methods)
- Struct literal: `netbox.Request{Name: "...", ...}`
- Constructor: `netbox.NewRequest(...)`
- Builder-like: `netbox.NewWritableRequest(...)`

### Issue 6: Unique Patterns (Should be flagged)
- **config_context_resource.go** - Only resource with `setToStringSlice()` for Tags
- **custom_link_resource.go** - Direct field assignments for all fields
- **interface_resource.go** - Mix of patterns within same resource

---

## Specifics by Resource Type

### Port Resources (HIGH INCONSISTENCY)
- console_port, console_server_port, front_port, rear_port
- power_port, power_outlet
- Plus all `*_template` variants
- **Pattern:** All use `apiReq.SetField()` setter methods
- **Issue:** ONLY these use setters; rest use direct assignment

### Template Resources
- All 8+ `*_template` resources
- **Pattern:** Mix of SetField() and direct assignment
- **Issue:** Inconsistent within templates group

### Assignment Resources
- circuit_group_assignment, contact_assignment, fhrp_group_assignment
- **Pattern:** Direct assignment + AdditionalProperties for complex refs
- **Issue:** AdditionalProperties seems like workaround

### Config Context (UNIQUE)
- **Pattern:** `request.Tags = setToStringSlice(ctx, data.Tags)`
- **Issue:** ONLY resource with this custom conversion function
- **Why:** Go-netbox ConfigContextRequest uses `[]string` for Tags, not `[]NestedTagRequest`

---

## Recommended Solution

**Adopt Helper Functions** (already partially done):

```go
// Phase 1: Create standardized helpers
utils.StringPtr(v)                              // Replaces &variable
utils.IntPtr(v)                                 // New
utils.BoolPtr(v)                                // New
utils.Float64Ptr(v)                             // New

// Already exist or expanded:
utils.ApplyDescription(req, data.Description)
utils.ApplyComments(req, data.Comments)
utils.ApplyTags(ctx, req, data.Tags, diags)
utils.ApplyCustomFields(ctx, req, data.CustomFields, diags)
utils.ApplyMetadataFields(ctx, req, tags, cf, diags)
utils.ApplyCommonFields(ctx, req, desc, comments, tags, cf, diags)

// Phase 2: Migrate all resources to use these helpers
```

---

## Implementation Effort

| Phase | Scope | Effort | Impact | Priority |
|-------|-------|--------|--------|----------|
| 1 | Create helpers | 40 hrs | Foundation | HIGH |
| 2a | Port/Console resources | 20 hrs | 15 resources | MEDIUM |
| 2b | Template resources | 20 hrs | 8 resources | MEDIUM |
| 2c | Device resources | 20 hrs | 10 resources | MEDIUM |
| 3 | Remaining resources | 50 hrs | 60 resources | LOW (can spread) |
| **TOTAL** | **Full standardization** | **~150 hours** | **Complete consistency** | |

---

## Key Metrics

- **Resources with direct assignment:** ~70
- **Resources with setter methods:** ~15
- **Resources with helpers:** ~17
- **Resources with mixed patterns:** ~10 (flag these!)
- **Unique pattern resources:** 3-5 (config_context, custom_link, interface, etc.)

---

## Immediate Actions

1. ✅ **Documentation** - This analysis identifies all issues
2. 📋 **Code Review Guide** - Add standardization checks
3. 🔧 **Helper Expansion** - Create Phase 1 helpers
4. 📊 **Tracking** - Add to project plan for gradual migration

---

## Key Takeaways

1. **You were correct** - Go-netbox does NOT have discrepancies; our code DOES
2. **Pattern adoption is inconsistent** - Likely evolved over time as different devs joined
3. **Low-hanging fruit** - Port resources (15) all use same non-standard pattern
4. **Already started** - Recent batches (13) successfully use helper pattern
5. **Should be completed** - Helps with Batch 1-13 goals; standardizes rest
