# Standardization Batches at a Glance

**14 Batches | 7 Phases | ~190 Hours | 85+ Resources**

---

## 📊 Visual Batch Layout

```
PHASE 1: FOUNDATION (9 hours) ⭐ PREREQUISITE
│
├─ S1.1: Pointer Helpers (StringPtr, IntPtr, etc)          [2 hrs]
├─ S1.2: ApplyDescription Expansion                        [1 hr]
├─ S1.3: Enum Conversion Helpers                           [3 hrs]
├─ S1.4: Reference Lookup Helpers                          [2 hrs]
└─ S1.5: Optional Field Helpers                            [1 hr]
         ✓ Blocks all other phases
         ✓ Ready to start immediately

PHASE 2: PORT RESOURCES (24 hours) 🎯 HIGH IMPACT
│
├─ S2A: Console Ports (3 resources)                        [6 hrs]
│   └─ console_port, console_port_template, console_server_port
│
├─ S2B: Front/Rear Ports (4 resources)                     [8 hrs]
│   └─ front_port, front_port_template, rear_port, rear_port_template
│
├─ S2C: Power Ports (4 resources)                          [8 hrs]
│   └─ power_port, power_port_template, power_outlet, power_outlet_template
│
└─ S2D: Inventory Item (1 resource)                        [2 hrs]
        └─ inventory_item
        ✓ 15 resources total
        ✓ All use setter pattern (easy migration)

PHASE 3: TEMPLATES (20 hours)
│
├─ S3A: Device/Rack/Cluster Templates (5 resources)        [10 hrs]
└─ S3B: Power/Console Templates (5 resources)              [10 hrs]
        ✓ 10 resources total
        ✓ Description-focused pattern

PHASE 4: SPECIAL CASES (20 hours) 🔴 CRITICAL PATH
│
├─ S4A: config_context (1 resource) ⚠️ UNIQUE OUTLIER     [4 hrs]
│       └─ setToStringSlice() custom pattern
│
├─ S4B: Assignment Resources (5 resources)                 [10 hrs]
│       └─ circuit_group_assignment, contact_assignment, etc
│
└─ S4C-D: Custom Link & Interface (2 resources)            [5 hrs]
         └─ custom_link, interface

PHASE 5: CORE DEVICE/CIRCUIT (33 hours) 🎯 HEAVY HITTERS
│
├─ S5A: Device Resources (3 resources)                     [12 hrs]
│       └─ device, device_type, device_role
│
├─ S5B: Circuit Resources (4 resources)                    [12 hrs]
│       └─ circuit, circuit_type, circuit_group, circuit_termination
│
└─ S5C: Cable & Related (3 resources)                      [9 hrs]
        └─ cable, cluster, event_rule

PHASE 6: MEDIUM PRIORITY (72+ hours) 📋 FLEXIBLE TIMING
│
├─ S6A: Location & Site (4 resources)                      [8 hrs]
├─ S6B: Contact Resources (4 resources)                    [8 hrs]
├─ S6C: Cluster Resources (3 resources)                    [6 hrs]
├─ S6D: VPN Resources (5 resources)                        [10 hrs]
├─ S6E: Routing & Network (5 resources)                    [10 hrs]
├─ S6F: Template & Config (5 resources)                    [10 hrs]
└─ S6G: Miscellaneous (8+ resources)                       [16 hrs]
        ✓ 40+ resources total
        ✓ Can overlap with Phase 5
        ✓ Can be done incrementally

PHASE 7: REMAINING (15+ hours) 🟢 LOW PRIORITY
│
└─ S7+: Remaining Resources (10-15 resources)              [15+ hrs]
        ✓ Touchpoint-based (as code naturally touches)
        ✓ No deadline
```

---

## ⏱️ Effort Timeline

```
Phase 1:    ████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  9 hours    (BLOCKING)
Phase 2:    ████████████░░░░░░░░░░░░░░░░░░░░░░░░ 24 hours    (High Priority)
Phase 3:    ██████████░░░░░░░░░░░░░░░░░░░░░░░░░░ 20 hours    (High Priority)
Phase 4:    ██████████░░░░░░░░░░░░░░░░░░░░░░░░░░ 20 hours    (Special Cases)
Phase 5:    ████████████████░░░░░░░░░░░░░░░░░░░░ 33 hours    (Heavy Hitters)
Phase 6:    ████████████████████████░░░░░░░░░░░░ 72 hours    (Flexible)
Phase 7:    ███░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 15 hours    (Ongoing)
────────────────────────────────────────────────────────────────────────
Total:      ████████████████████████████████████ 193 hours   (2-3 months)
```

---

## 🎯 Recommended Phasing

### PHASE 1: FOUNDATION (Week 1)
```
Mon-Tue: S1.1 + S1.2 (Pointer helpers)
Wed-Thu: S1.3 + S1.4 (Enum/Reference helpers)
Fri:     S1.5 + Code review + Integration
```
**Output**: Reusable helper suite for all subsequent batches

### PHASE 2: PORT RESOURCES (Week 2-3)
```
Week 2:
  Mon-Tue: S2A (Console Ports - 6 hrs)
  Wed-Fri: S2B (Front/Rear Ports - 8 hrs)

Week 3:
  Mon-Wed: S2C (Power Ports - 8 hrs)
  Thu-Fri: S2D (Inventory Item - 2 hrs) + PR
```
**Output**: 15 consistent port resources, highest impact quick win

### PHASE 3-4: TEMPLATES & SPECIAL (Week 4)
```
Mon-Fri: S3 (Templates - 20 hrs) + S4A (config_context - 4 hrs)
Bonus:   S4B-D (Assignment/Special - 15 hrs) if time
```
**Output**: 10 templates + critical config_context resolved

### PHASE 5: CORE RESOURCES (Week 5-6)
```
Week 5: S5A + S5B (Device/Circuit - 24 hrs)
Week 6: S5C + S6A start (Cable/Location - 10 hrs)
```
**Output**: 10 most-used resources standardized

### PHASE 6-7: REMAINING (Week 7+)
```
Ongoing: S6B-G + S7 (50+ resources)
Pace:    2-3 batches per week
Approach: Touchpoint-based, no deadline
```
**Output**: 100% standardization complete

---

## 📊 Resource Count Breakdown

```
Phase 1: - (Foundation only)
Phase 2: ████████████ 15 resources
Phase 3: ██████████ 10 resources
Phase 4: ████████ 8 resources
Phase 5: ██████████ 10 resources
Phase 6: ████████████████████████████████████ 40+ resources
Phase 7: ███████████ 10-15 resources
         ────────────────────────────
Total:   85-90 resources

Progress after Phase 1-2: 25/85 (29%)
Progress after Phase 1-5: 43/85 (51%)
Progress after Phase 1-6: 83/85 (98%)
```

---

## 🔄 Dependency Chain

```
START → S1 (Foundation)
         ↓
      S2 (Ports) ──────┐
         ↓              ├──→ S3 (Templates)
      S4A (config)──┐   │
         ↓          ├──→ S4B-D (Assignments)
      S4B-D────────┤
         ↓          └──→ S5 (Device/Circuit)
      S5 ──────────────→ S6 (Medium Priority)
         ↓
      S6 ──────────────→ S7 (Remaining)
         ↓
      COMPLETE ✅

Legend:
S = Batch
→ = depends on
```

---

## ✅ Success Checklist

### Per Batch
- [ ] All target resources identified
- [ ] Helper functions created (for Phase 1)
- [ ] Direct field assignments removed
- [ ] Setter methods replaced with helpers
- [ ] `go build .` succeeds
- [ ] Existing tests pass
- [ ] Code review approved
- [ ] Update STANDARDIZATION_TRACKING.md

### Overall
- [ ] 85+ resources standardized
- [ ] 3 patterns → 1 pattern
- [ ] No functional changes
- [ ] Documentation updated
- [ ] Team trained

---

## 🚀 Quick Start Commands

```powershell
# Create feature branch for Phase 1
git checkout -b refactor/request-standardization-phase1

# Start work on S1.1
# Edit: internal/utils/request_helpers.go
# Add pointer helpers: StringPtr, IntPtr, BoolPtr, Float64Ptr

# Verify
go build .
go test ./...

# When complete
git commit -m "Phase 1: Request standardization foundation helpers"
git push

# Update tracking
# Edit: STANDARDIZATION_TRACKING.md
# Mark S1.1 ✅ DONE, move to S1.2
```

---

## 📌 Key Dates

| Milestone | Target | Status |
|-----------|--------|--------|
| Phase 1 Complete | This week | ⏳ Not started |
| Phase 2 Complete | Next week | ⏳ Blocked by P1 |
| Phase 3-4 Complete | Week 3-4 | ⏳ Blocked by P1-2 |
| Phase 5 Complete | Week 5-6 | ⏳ Blocked by P4 |
| 50% Overall (Phase 6 start) | Week 7 | ⏳ Blocked by P5 |
| 100% Complete | Week 9-12 | ⏳ Depends on velocity |

---

## 💡 Pro Tips

1. **Start Phase 1**: Do foundation work first (blocks nothing, enables everything)
2. **Phase 2 is Quick Win**: 24 hours for 15 resources, consistent pattern
3. **Parallel Work**: Can do Phase 6 while Phase 5 is in progress
4. **Touchpoint Approach**: Phase 7 resources handled as code naturally touches them
5. **Commit Often**: Commit after each sub-batch (S1.1, S1.2, etc)
6. **Keep Tests Running**: `go test ./...` after each batch
7. **Update Tracking**: Keep STANDARDIZATION_TRACKING.md current

---

## 🎓 Pattern Reference

### Pattern A (Current - 68%)
```go
// BEFORE: Direct field assignment
desc := data.Description.ValueString()
request.Description = &desc
```

### Pattern B (Current - 15%)
```go
// BEFORE: Setter methods
apiReq.SetDescription(data.Description.ValueString())
```

### Pattern C (Target - 100%)
```go
// AFTER: Helper functions
utils.ApplyDescription(request, data.Description)
```

All three patterns do the same thing. Pattern C is clearest and most maintainable.

---

**Created**: 2025-12-27
**Status**: Ready to Execute
**Next**: Start Phase 1 → Read STANDARDIZATION_EXECUTIVE_SUMMARY.md
