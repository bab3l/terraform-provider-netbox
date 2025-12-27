# 🎯 Standardization Execution Checklist

**Date Started**: ________________
**Target Completion**: Week ____ (estimate based on hours)

---

## PHASE 1: FOUNDATION (9 hours) - DO THIS FIRST! ⭐

### S1.1: Pointer Helpers (2 hours)

**File**: `internal/utils/request_helpers.go`

- [ ] Create file `internal/utils/request_helpers.go`
- [ ] Add `StringPtr()` function
- [ ] Add `IntPtr()` function
- [ ] Add `BoolPtr()` function
- [ ] Add `Float64Ptr()` function
- [ ] Add unit tests for each
- [ ] `go build .` ✅ Compiles
- [ ] `go test ./...` ✅ Tests pass
- [ ] Commit: "Phase 1.1: Add pointer helper functions"
- [ ] Update STANDARDIZATION_TRACKING.md → S1.1 = ✅ DONE

**Estimated Time**: 2 hours
**Hours Spent**: ______

---

### S1.2: Expand ApplyDescription Helper (1 hour)

**Files**: Review existing ApplyDescription in `internal/utils/`

- [ ] Verify null checking works
- [ ] Verify unknown checking works
- [ ] Verify pointer conversion works
- [ ] Create/update helper documentation
- [ ] Add examples to docstring
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Phase 1.2: Document ApplyDescription helper"
- [ ] Update STANDARDIZATION_TRACKING.md → S1.2 = ✅ DONE

**Estimated Time**: 1 hour
**Hours Spent**: ______

---

### S1.3: Enum Conversion Helpers (3 hours)

**File**: `internal/utils/request_helpers.go` (extend)

- [ ] Identify enum types used in resources
- [ ] Create `ToDeviceStatus()` helper
- [ ] Create `ToDeviceRoleColor()` helper
- [ ] Create `ToCableType()` helper
- [ ] Create additional enum converters (estimate 3-4 total)
- [ ] Add error handling
- [ ] Add unit tests
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Phase 1.3: Add enum conversion helpers"
- [ ] Update STANDARDIZATION_TRACKING.md → S1.3 = ✅ DONE

**Estimated Time**: 3 hours
**Hours Spent**: ______

---

### S1.4: Reference Lookup Helpers (2 hours)

**File**: `internal/utils/request_helpers.go` (extend)

- [ ] Create `SetRequiredReference()` helper
- [ ] Add context parameter handling
- [ ] Add client parameter handling
- [ ] Add lookup function support
- [ ] Add diagnostics/error handling
- [ ] Add unit tests
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Phase 1.4: Add reference lookup helpers"
- [ ] Update STANDARDIZATION_TRACKING.md → S1.4 = ✅ DONE

**Estimated Time**: 2 hours
**Hours Spent**: ______

---

### S1.5: Optional Field Helper (1 hour)

**File**: `internal/utils/request_helpers.go` (extend)

- [ ] Create `ApplyOptionalField()` helper
- [ ] Add null/unknown check
- [ ] Add setter function support
- [ ] Add error handling
- [ ] Add unit tests
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Phase 1.5: Add optional field helper"
- [ ] Update STANDARDIZATION_TRACKING.md → S1.5 = ✅ DONE

**Estimated Time**: 1 hour
**Hours Spent**: ______

---

### PHASE 1 COMPLETION CHECKLIST

- [ ] All 5 sub-batches complete
- [ ] `go build .` succeeds
- [ ] `go test ./...` succeeds
- [ ] All helper functions documented
- [ ] STANDARDIZATION_TRACKING.md shows Phase 1 = ✅ DONE
- [ ] Create PR: "Phase 1: Request standardization foundation helpers"
- [ ] PR approved and merged

**Phase 1 Total Hours**: _______ / 9 hours
**Status**: ⏳ IN PROGRESS → 🚀 READY FOR S2

---

## PHASE 2: PORT RESOURCES (24 hours) - QUICK WIN! 🎯

### S2A: Console Port Resources (6 hours)

**Resources** (3):
1. `console_port_resource.go`
2. `console_port_template_resource.go`
3. `console_server_port_resource.go`

#### console_port_resource.go
- [ ] Open file in editor
- [ ] Identify all setter calls: SetDescription, SetLabel, SetType, SetSpeed, SetMarkConnected
- [ ] Replace each with helper calls
- [ ] Remove apiReq.SetXxx() patterns
- [ ] Add helper imports if needed
- [ ] Create: method for reusable helper calls
- [ ] Both Create() and Update() methods updated
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Refactor: Standardize console_port_resource to use helpers"

#### console_port_template_resource.go
- [ ] Open file in editor
- [ ] Identify setter patterns
- [ ] Replace with helpers
- [ ] Both Create() and Update() methods
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Refactor: Standardize console_port_template_resource to use helpers"

#### console_server_port_resource.go
- [ ] Open file in editor
- [ ] Identify setter patterns
- [ ] Replace with helpers
- [ ] Both Create() and Update() methods
- [ ] `go build .` ✅ Compiles
- [ ] Commit: "Refactor: Standardize console_server_port_resource to use helpers"

**S2A Verification**:
- [ ] All 3 resources compile
- [ ] All 3 resources use helper pattern
- [ ] `go test ./...` passes
- [ ] Update STANDARDIZATION_TRACKING.md → S2A = ✅ DONE

**Estimated Time**: 6 hours
**Hours Spent**: ______

---

### S2B: Front/Rear Port Resources (8 hours)

**Resources** (4):
1. `front_port_resource.go`
2. `front_port_template_resource.go`
3. `rear_port_resource.go`
4. `rear_port_template_resource.go`

- [ ] front_port_resource.go → Replace setters with helpers
- [ ] front_port_template_resource.go → Replace with helpers
- [ ] rear_port_resource.go → Replace setters with helpers
- [ ] rear_port_template_resource.go → Replace with helpers
- [ ] All 4 compile: `go build .` ✅
- [ ] All tests pass: `go test ./...` ✅
- [ ] Update STANDARDIZATION_TRACKING.md → S2B = ✅ DONE

**Estimated Time**: 8 hours
**Hours Spent**: ______

---

### S2C: Power Port Resources (8 hours)

**Resources** (4):
1. `power_port_resource.go`
2. `power_port_template_resource.go`
3. `power_outlet_resource.go`
4. `power_outlet_template_resource.go`

- [ ] power_port_resource.go → Replace setters with helpers
- [ ] power_port_template_resource.go → Replace with helpers
- [ ] power_outlet_resource.go → Replace setters with helpers
- [ ] power_outlet_template_resource.go → Replace with helpers
- [ ] All 4 compile: `go build .` ✅
- [ ] All tests pass: `go test ./...` ✅
- [ ] Update STANDARDIZATION_TRACKING.md → S2C = ✅ DONE

**Estimated Time**: 8 hours
**Hours Spent**: ______

---

### S2D: Inventory Item (2 hours)

**Resources** (1):
1. `inventory_item_resource.go`

- [ ] Open file in editor
- [ ] Identify remaining direct assignments
- [ ] Replace with helper pattern
- [ ] Both Create() and Update() methods
- [ ] `go build .` ✅ Compiles
- [ ] `go test ./...` ✅ Passes
- [ ] Commit: "Refactor: Standardize inventory_item_resource to use helpers"
- [ ] Update STANDARDIZATION_TRACKING.md → S2D = ✅ DONE

**Estimated Time**: 2 hours
**Hours Spent**: ______

---

### PHASE 2 COMPLETION CHECKLIST

- [ ] All 4 sub-batches (S2A-S2D) complete
- [ ] 15 resources all use helper pattern
- [ ] `go build .` succeeds
- [ ] `go test ./...` succeeds
- [ ] All Phase 2 resources verified
- [ ] Create PR: "Phase 2: Standardize port resources (15 resources)"
- [ ] PR approved and merged

**Phase 2 Total Hours**: _______ / 24 hours
**Resources Standardized**: 15 / 85
**Status**: ⏳ IN PROGRESS → 🚀 READY FOR S3

---

## PHASE 3: TEMPLATES (20 hours)

### S3A: Device/Rack/Cluster Templates (10 hours)

**Resources** (5):
1. `device_bay_template_resource.go`
2. `interface_template_resource.go`
3. `module_bay_template_resource.go`
4. `inventory_item_template_resource.go`
5. `config_template_resource.go`

- [ ] All 5 templates identified
- [ ] Replace direct assignment with helpers
- [ ] All templates compile
- [ ] All tests pass
- [ ] Create PR: "Phase 3A: Standardize template resources (5 resources)"
- [ ] Update STANDARDIZATION_TRACKING.md → S3A = ✅ DONE

**Estimated Time**: 10 hours
**Hours Spent**: ______

---

### S3B: Power/Console Templates (10 hours)

**Resources** (5):
1. `power_port_template_resource.go` (from S2C)
2. `power_outlet_template_resource.go` (from S2C)
3. `console_port_template_resource.go` (from S2A)
4. `console_server_port_template_resource.go` (from S2A)
5. `rear_port_template_resource.go` (from S2B)

**Note**: These already handled in Phase 2, just verify

- [ ] All 5 verified as standardized
- [ ] All compile
- [ ] All tests pass
- [ ] Update STANDARDIZATION_TRACKING.md → S3B = ✅ DONE

**Estimated Time**: 10 hours (mostly verification)
**Hours Spent**: ______

---

### PHASE 3 COMPLETION CHECKLIST

- [ ] All template resources standardized (10 new + 5 from Phase 2)
- [ ] `go build .` succeeds
- [ ] `go test ./...` succeeds
- [ ] Create PR: "Phase 3: Standardize template resources (10 resources)"
- [ ] PR approved and merged

**Phase 3 Total Hours**: _______ / 20 hours
**Resources Standardized Cumulatively**: 25 / 85
**Status**: ⏳ IN PROGRESS → 🚀 READY FOR S4

---

## CONTINUING PHASES

### Phase 4: Special Cases (20 hours)
- S4A: config_context (4 hrs) ⚠️ Unique case
- S4B: Assignments (10 hrs)
- S4C-D: Custom Link & Interface (5 hrs)

### Phase 5: Core Resources (33 hours)
- S5A: Device Resources (12 hrs)
- S5B: Circuit Resources (12 hrs)
- S5C: Cable & Related (9 hrs)

### Phase 6: Medium Priority (72 hours)
- S6A-G: Remaining resources (can overlap with Phase 5)

### Phase 7: Touchpoint (15+ hours)
- Ongoing as code naturally touches files

---

## 📊 PROGRESS TRACKER

### Resources Standardized So Far

```
Phase 1: Foundation     ░░░░░░░░░░░░░░░░  0 resources (prerequisite)
Phase 2: Ports         ████░░░░░░░░░░░░ 15 resources
Phase 3: Templates     ██░░░░░░░░░░░░░░ 10 resources
Phase 4: Special       █░░░░░░░░░░░░░░░  8 resources
Phase 5: Core          █░░░░░░░░░░░░░░░ 10 resources
Phase 6: Medium        ████████░░░░░░░░ 40+ resources
Phase 7: Remaining     ██░░░░░░░░░░░░░░ 10-15 resources
                       ────────────────────────────────
Total Planned:         85 resources (100%)
```

### Cumulative Hours

```
After Phase 1:   9 hours   (Foundation)
After Phase 2:  33 hours   (9 + 24)
After Phase 3:  53 hours   (33 + 20)
After Phase 4:  73 hours   (53 + 20)
After Phase 5: 106 hours   (73 + 33)
After Phase 6: 178 hours   (106 + 72)
After Phase 7: 193 hours   (178 + 15)
```

---

## 🎓 HELPFUL COMMANDS

```powershell
# Build check
go build .

# Run tests
go test ./... -v

# Find setter patterns in file
grep -n "\.Set[A-Z]" <resource_file>

# Find direct assignment patterns
grep -n "\. [A-Z][a-zA-Z]* = &" <resource_file>

# Check which files use old patterns
grep -r "\.SetDescription\|\.SetComments" internal/resources/

# After changes, verify build
go build . && echo "✅ Build successful"

# Run specific resource tests
go test ./internal/resources -run TestResource -v
```

---

## ✅ FINAL CHECKLIST

### When Everything is Complete:
- [ ] Phase 1-7 all complete
- [ ] 85+ resources standardized
- [ ] `go build .` succeeds
- [ ] `go test ./...` succeeds
- [ ] All PRs merged
- [ ] Documentation updated
- [ ] CONTRIBUTING.md updated with patterns
- [ ] Team trained on new patterns
- [ ] Code review checklist updated

### Success Metrics Verified:
- [ ] 0 direct assignment patterns remaining (Pattern A)
- [ ] 0 setter methods remaining (Pattern B)
- [ ] 100% helper pattern usage (Pattern C)
- [ ] All resources follow same pattern
- [ ] Consistent error handling
- [ ] ~200-300 lines of boilerplate removed

---

**Start Date**: _______________
**Phase 1 Complete**: _______________
**Phase 2 Complete**: _______________
**Phase 5 Complete**: _______________
**All Phases Complete**: _______________

**Total Time Invested**: _______ hours
**Actual vs Estimated**: _______ hours (variance)

**Notes**:
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
