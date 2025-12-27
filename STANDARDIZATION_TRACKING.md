# Request Standardization - Progress Tracking

**Last Updated**: 2025-12-27
**Phase**: Phase 2A Complete → Ready for Phase 2B (Front/Rear Ports)
**Overall Progress**: 3/85 resources standardized (4%) | Phase 1: ✅ COMPLETE | Phase 2A: ✅ COMPLETE

---

## Batch Status Summary

### Phase 1: Foundation (PREREQUISITE)

| ID | Task | Status | Hours | Estimated Completion |
|----|------|--------|-------|----------------------|
| S1.1 | Create pointer helpers (StringPtr, IntPtr, etc) | ✅ DONE | 2 | 2025-12-27 |
| S1.2 | Expand ApplyDescription helper | ✅ DONE | 1 | 2025-12-27 |
| S1.3 | Create enum conversion helpers | ✅ DONE | 3 | 2025-12-27 |
| S1.4 | Create reference lookup helpers | ✅ DONE | 2 | 2025-12-27 |
| S1.5 | Create optional field helper | ✅ DONE | 1 | 2025-12-27 |
| **PHASE 1 TOTAL** | **Foundation Complete** | **✅ COMPLETE** | **9** | **2025-12-27** |

**Blocker Status**: ✅ UNBLOCKED - All helpers validated and tested

---

### Phase 2: Port Resources (15 resources - 24 hours)

#### Batch S2A: Console Port Resources (3 resources, 6 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| console_port_resource.go | ✅ DONE | Removed redundant SetDescription, using ApplyDescription + ApplyMetadataFields | 2 |
| console_port_template_resource.go | ✅ DONE | Already clean - no redundancy found | 2 |
| console_server_port_resource.go | ✅ DONE | Removed redundant SetDescription, using ApplyDescription + ApplyMetadataFields | 2 |
| **S2A SUBTOTAL** | **✅ COMPLETE** | **3/3 resources done** | **6** |

#### Batch S2B: Front/Rear Port Resources (4 resources, 8 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| front_port_resource.go | ⏳ NOT STARTED | SetDescription and related setters | 2 |
| front_port_template_resource.go | ⏳ NOT STARTED | SetDeviceType pattern | 2 |
| rear_port_resource.go | ⏳ NOT STARTED | SetDescription and related setters | 2 |
| rear_port_template_resource.go | ⏳ NOT STARTED | Mixed pattern migration | 2 |
| **S2B SUBTOTAL** | **⏳ READY AFTER S1** | | **8** |

#### Batch S2C: Power Port Resources (4 resources, 8 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| power_port_resource.go | ⏳ NOT STARTED | SetDescription pattern | 2 |
| power_port_template_resource.go | ⏳ NOT STARTED | Direct field → helpers | 2 |
| power_outlet_resource.go | ⏳ NOT STARTED | SetDescription pattern | 2 |
| power_outlet_template_resource.go | ⏳ NOT STARTED | Direct field → helpers | 2 |
| **S2C SUBTOTAL** | **⏳ READY AFTER S1** | | **8** |

#### Batch S2D: Inventory Item (1 resource, 2 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| inventory_item_resource.go | ⏳ NOT STARTED | Already partially refactored, finish remaining | 2 |
| **S2D SUBTOTAL** | **⏳ READY AFTER S1** | | **2** |

| **PHASE 2 TOTAL** | **⏳ READY AFTER S1** | **15 resources** | **24** |

---

### Phase 3: Template Resources (10 resources - 20 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| device_bay_template_resource.go | ⏳ NOT STARTED | Description only | 1 |
| interface_template_resource.go | ⏳ NOT STARTED | Description only | 1 |
| module_bay_template_resource.go | ⏳ NOT STARTED | Description only | 1 |
| inventory_item_template_resource.go | ⏳ NOT STARTED | Description only | 1 |
| config_template_resource.go | ⏳ NOT STARTED | Description only | 1 |
| power_port_template_resource.go | ⏳ READY (Phase 2C) | - | - |
| power_outlet_template_resource.go | ⏳ READY (Phase 2C) | - | - |
| console_port_template_resource.go | ⏳ READY (Phase 2A) | - | - |
| console_server_port_template_resource.go | ⏳ READY (Phase 2A) | - | - |
| rear_port_template_resource.go | ⏳ READY (Phase 2B) | - | - |
| **PHASE 3 TOTAL** | **⏳ READY AFTER S1-S6** | **10 resources** | **20** |

---

### Phase 4: Special Cases (8 resources - 20 hours)

#### Phase 4A: config_context (1 resource, 4 hours)

| Resource | Status | Issue | Hours |
|----------|--------|-------|-------|
| config_context_resource.go | ⏳ NOT STARTED | 🔴 CRITICAL - Unique `setToStringSlice()` pattern | 4 |
| **4A SUBTOTAL** | **⏳ READY AFTER S1** | **BLOCKER** | **4** |

#### Phase 4B: Assignment Resources (5 resources, 10 hours)

| Resource | Status | Issue | Hours |
|----------|--------|-------|-------|
| circuit_group_assignment_resource.go | ⏳ NOT STARTED | AdditionalProperties pattern | 2 |
| contact_assignment_resource.go | ⏳ NOT STARTED | AdditionalProperties pattern | 2 |
| fhrp_group_assignment_resource.go | ⏳ NOT STARTED | AdditionalProperties pattern | 2 |
| tunnel_termination_resource.go | ⏳ NOT STARTED | AdditionalProperties pattern | 2 |
| tunnel_group_resource.go | ⏳ NOT STARTED | AdditionalProperties pattern | 2 |
| **4B SUBTOTAL** | **⏳ READY AFTER S1-S2** | | **10** |

#### Phase 4C-D: Custom Link & Interface (2 resources, 5 hours)

| Resource | Status | Issue | Hours |
|----------|--------|-------|-------|
| custom_link_resource.go | ⏳ NOT STARTED | All direct fields: Enabled, Weight, GroupName, etc | 2 |
| interface_resource.go | ⏳ NOT STARTED | Mixed patterns within resource | 3 |
| **4C-D SUBTOTAL** | **⏳ READY AFTER S1-S2** | | **5** |

| **PHASE 4 TOTAL** | **⏳ READY AFTER S1-S2** | **8 resources** | **20** |

---

### Phase 5: Core Device & Circuit (10 resources - 33 hours)

| Resource | Status | Fields | Hours |
|----------|--------|--------|-------|
| device_resource.go | ⏳ NOT STARTED | Serial, Face, Status, Airflow, etc (heavy) | 4 |
| device_type_resource.go | ⏳ NOT STARTED | PartNumber, UHeight, ExcludeFromUtilization, etc | 4 |
| device_role_resource.go | ⏳ NOT STARTED | Color, VmRole | 2 |
| circuit_resource.go | ⏳ NOT STARTED | Description, Comments, Tags, CustomFields | 3 |
| circuit_type_resource.go | ⏳ NOT STARTED | Description, Tags, CustomFields | 3 |
| circuit_group_resource.go | ⏳ NOT STARTED | Mixed pointer patterns | 3 |
| circuit_termination_resource.go | ⏳ NOT STARTED | Description, Tags, CustomFields | 3 |
| cable_resource.go | ⏳ NOT STARTED | Type, Status, Label, Color, LengthUnit | 4 |
| cluster_resource.go | ⏳ NOT STARTED | Status, Description, Comments, Tags, CustomFields | 3 |
| event_rule_resource.go | ⏳ NOT STARTED | Enabled, ActionType, Description, Tags, CustomFields | 3 |
| **PHASE 5 TOTAL** | **⏳ READY AFTER S1-S8** | **10 resources** | **33** |

---

### Phase 6: Medium Priority (40+ resources - 72+ hours)

#### Phase 6A: Location & Site (4 resources, 8 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| location_resource.go | ⏳ NOT STARTED | Description, Parent, Status, Facility, Tags, CustomFields | 2 |
| site_resource.go | ⏳ NOT STARTED | Description, Comments, Tags, CustomFields | 2 |
| site_group_resource.go | ⏳ NOT STARTED | Description, Tags, CustomFields | 2 |
| region_resource.go | ✅ DONE | Already refactored | - |
| **6A SUBTOTAL** | **⏳ READY AFTER S1** | | **8** |

#### Phase 6B: Contact Resources (4 resources, 8 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| contact_resource.go | ⏳ NOT STARTED | Title, Phone, Email, Address, Link, Description, Comments, Tags | 2 |
| contact_group_resource.go | ⏳ NOT STARTED | Description, Tags, CustomFields | 2 |
| contact_role_resource.go | ⏳ NOT STARTED | Description, Tags, CustomFields | 2 |
| contact_assignment_resource.go | ⏳ READY (Phase 4B) | AdditionalProperties pattern | - |
| **6B SUBTOTAL** | **⏳ READY AFTER S1** | | **8** |

#### Phase 6C: Cluster Resources (3 resources, 6 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| cluster_group_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| cluster_type_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| fhrp_group_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| **6C SUBTOTAL** | **⏳ READY AFTER S1** | | **6** |

#### Phase 6D: VPN Resources (5 resources, 10 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| ipsec_policy_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| ipsec_profile_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| ipsec_proposal_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| ike_policy_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| ike_proposal_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| **6D SUBTOTAL** | **⏳ READY AFTER S1** | | **10** |

#### Phase 6E: Routing & Network (5 resources, 10 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| route_target_resource.go | ✅ DONE | Already refactored | - |
| ip_range_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| l2vpn_resource.go | ✅ DONE | Already refactored | - |
| tunnel_resource.go | ✅ DONE | Already refactored | - |
| wireless_lan_resource.go | ✅ DONE | Already refactored | - |
| **6E SUBTOTAL** | **⏳ READY AFTER S1** | | **10** |

#### Phase 6F: Template & Config (5 resources, 10 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| export_template_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| custom_field_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| custom_field_choice_set_resource.go | ⏳ NOT STARTED | Direct field | 2 |
| journal_entry_resource.go | ✅ DONE | Already refactored | - |
| event_rule_resource.go | ⏳ READY (Phase 5) | - | - |
| **6F SUBTOTAL** | **⏳ READY AFTER S1** | | **10** |

#### Phase 6G: Miscellaneous (8 resources, 16 hours)

| Resource | Status | Notes | Hours |
|----------|--------|-------|-------|
| asn_range_resource.go | ⏳ NOT STARTED | Description, Tags, CustomFields | 2 |
| device_bay_resource.go | ✅ DONE | Already refactored | - |
| manufacturer_resource.go | ✅ DONE | Using ApplyDescription | - |
| platform_resource.go | ✅ DONE | Using ApplyDescription | - |
| rack_type_resource.go | ⏳ NOT STARTED | Description only | 1 |
| rir_resource.go | ✅ DONE | Using ApplyMetadataFields | - |
| + 10 others | ⏳ NOT STARTED | Various patterns | 12 |
| **6G SUBTOTAL** | **⏳ READY AFTER S1** | | **16** |

| **PHASE 6 TOTAL** | **⏳ READY AFTER S1** | **40+ resources** | **72+** |

---

### Phase 7: Remaining/Touchpoint (10-15 resources - 15+ hours)

| Status | Count | Hours | Notes |
|--------|-------|-------|-------|
| ⏳ NOT STARTED | 10-15 | 15+ | Ongoing, handled as development touchpoints |

---

## Overall Progress Dashboard

### Completion by Resource Count

```
████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0/85 Resources (0%)

Legend:
████ Completed
░░░░ Not Started
```

### Completion by Effort Hours

```
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0/190 Hours (0%)

Legend:
████ Completed
░░░░ Not Started
```

### Batch Execution Timeline

```
Phase 1: S1           [9 hrs, READY]
                            ↓
Phase 2: S2-S6      [24 hrs, BLOCKED by S1]
                            ↓
Phase 3: S7         [20 hrs, BLOCKED by S2]
                            ↓
Phase 4: S8-S9      [15 hrs, BLOCKED by S2]
                            ↓
Phase 5: S10-S12    [33 hrs, BLOCKED by S4]
                            ↓
Phase 6: S13+       [72 hrs, can overlap with Phase 5]
                            ↓
Phase 7: S14+       [15+ hrs, ongoing touchpoint]

Total Linear Time: ~190 hours
With Parallelization: ~100+ hours (Phase 6 can overlap)
```

---

## Key Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Resources Standardized | 85 | 0 | 0% |
| Foundation Complete | 1/1 | 0/1 | ⏳ Not Started |
| Phase 2 (Port Resources) Complete | 15 | 0 | 0% |
| Phase 3 (Templates) Complete | 10 | 0 | 0% |
| Phase 4 (Special Cases) Complete | 8 | 0 | 0% |
| Phase 5 (Device/Circuit) Complete | 10 | 0 | 0% |
| Phase 6 (Medium Priority) Complete | 40+ | 0 | 0% |
| **TOTAL HOURS INVESTED** | ~190 | ~0 | 0% |

---

## Dependency Map

```
                    S1 (Foundation)
                          |
                __________|__________
               |                     |
              S2               S8-S9  S4A
           (Ports)         (Special) (config)
               |                |
               └────────┬───────┘
                        |
                       S7         S10-S12
                    (Templates)  (Device)
                        |
                        └─────┬─────┘
                              |
                           S13+
                      (Medium Priority)
```

---

## Quick Start Instructions

### When S1 is Complete:
```
1. Checkout feature branch: git checkout -b refactor/request-standardization-s2
2. Start with S2A (Console Ports): 3 resources, 6 hours
3. Apply helper pattern from S1 foundation
4. Build verification: go build .
5. Create PR with 3 resources
```

### Resource Checklist Template
```
- [ ] Resource file open
- [ ] Identify all field assignments (grep for "request.X = ")
- [ ] Create helper calls (utils.ApplyX pattern)
- [ ] Remove direct assignments
- [ ] Check Build: go build .
- [ ] Run tests: go test ./...
- [ ] Commit changes
```

---

## Notes & Considerations

### What Changes:
- ❌ Removes direct field assignment patterns
- ✅ Adds helper function calls
- ✅ Centralizes pointer handling
- ✅ Standardizes error handling for references

### What Stays Same:
- ✅ Resource behavior (no functional change)
- ✅ Create/Update method signatures
- ✅ Reference resolution logic (separate helpers)
- ✅ Test suite

### Risk Mitigation:
- Build verification after each batch
- Systematic one-batch-at-a-time approach
- Reference existing helper patterns
- Keep PRs focused (one batch per PR)

---

## Status Legend

| Status | Meaning |
|--------|---------|
| ✅ DONE | Completed and verified |
| 🚀 IN PROGRESS | Currently being worked on |
| ⏳ NOT STARTED | Planned but not begun |
| ⏳ READY AFTER X | Blocked until dependency X completes |
| 🔴 BLOCKER | Critical issue needing resolution |
| 🔄 NEEDS REVIEW | Completed, awaiting approval |

---

**Last Updated**: 2025-12-27
**Next Review**: After S1 completion
**Questions/Blockers**: See REQUEST_STANDARDIZATION_IMPLEMENTATION_PLAN.md for details
