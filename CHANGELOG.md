# Changelog

## v0.0.10 (TBD)

### Features & Improvements

#### Terraform Integration Test Suite
*   **Complete Phase 10: All 204 Terraform Integration Tests Passing**
    - 206 test invocations across 204 unique test directories
    - 100% pass rate (206/206 successful)
    - ~3 hours of comprehensive end-to-end validation
    - Tests cover all 101 resource types and 103 data source types

**Test Coverage by Category:**
- DCIM: 80+ tests (devices, racks, interfaces, ports, components, templates)
- IPAM: 40+ tests (IP addresses, prefixes, VLANs, ASNs, aggregates)
- Virtualization: 20+ tests (VMs, clusters, disks, interfaces)
- Circuits: 20+ tests (circuits, providers, terminations)
- VPN & Security: 20+ tests (tunnels, IPSec, IKE policies)
- Wireless & Tenancy: 15+ tests
- Extras: 15+ tests (webhooks, events, templates, configs)

**Test Organization:**
- Batch 10.1-10.10: 204 tests organized by dependency and functionality
- Each batch achieves 100% pass rate
- Average test execution: 35-50 seconds (includes setup, apply, verify, cleanup)

#### Test Infrastructure Improvements
*   **Enhanced Resource Cleanup**
    - Added pre-cleanup for aggregates to prevent overlapping aggregate errors
    - Added pre-cleanup for IP ranges to prevent overlapping range errors
    - Intelligent multi-iteration cleanup handling dependencies
    - Comprehensive deletion order mapping for 70+ resource types
*   **Configuration Standardization**
    - All test configurations include proper terraform block
    - Consistent required_providers configuration across all tests
    - Proper dependency ordering for resource and data source tests

#### Bug Fixes in Test Configurations
*   **Fixed Port Template Issues (Batch 10.5)**
    - Corrected rear_port reference in front_port_template (ID → NAME)
    - Fixed missing terraform blocks in device component templates
    - Port template tests now fully passing
*   **Fixed Data Source Schema Issues (Batch 10.6)**
    - Corrected front_port data source outputs to match actual schema
    - Removed unsupported attribute references
    - Aligned test patterns with other successful data sources
*   **Fixed IP Data Overlaps (Batch 10.7)**
    - Added aggregate and IP range pre-cleanup to prevent conflicts
    - Resolves "400 Bad Request" from overlapping network ranges
    - Enables reliable test execution across multiple runs

### Breaking Changes

#### Removed Duplicate `_id` Computed Fields
*   **Removed duplicate ID fields from 24 resources**
    - These fields duplicated information already available in primary reference fields
    - Primary fields (tenant, cluster, site, etc.) work for all references
    - Removed fields were causing unnecessary plan noise with "(known after apply)"
    - State migration is automatic (Terraform drops removed computed fields)

**Resources affected and fields removed:**
- `netbox_device`: device_type_id, role_id, tenant_id, platform_id, site_id, location_id, rack_id (7 fields)
- `netbox_virtual_machine`: site_id, cluster_id, role_id, tenant_id, platform_id (5 fields)
- `netbox_rack`: site_id, location_id, tenant_id, role_id, rack_type_id (5 fields)
- `netbox_rack_type`: manufacturer_id (1 field)
- `netbox_vlan`: site_id, tenant_id (2 fields)
- `netbox_vrf`: tenant_id (1 field)
- `netbox_route_target`: tenant_id (1 field)
- `netbox_site_group`: parent_id (1 field)
- `netbox_tenant_group`: parent_id (1 field)
- `netbox_tenant`: group_id (1 field)
- `netbox_region`: parent_id (1 field)
- `netbox_location`: parent_id (1 field)
- `netbox_platform`: manufacturer_id (1 field)

**Migration Guide:**

If you were referencing `_id` fields in your configuration:

```hcl
# Before (v0.0.9)
resource "netbox_ip_address" "example" {
  tenant = netbox_virtual_machine.vm.tenant_id  # ❌ No longer available
}

# After (v0.0.10)
resource "netbox_ip_address" "example" {
  tenant = netbox_virtual_machine.vm.tenant     # ✅ Use primary field
  # OR reference the source directly:
  tenant = netbox_tenant.my_tenant.id
}
```

**Primary reference fields work with name, slug, or ID:**
```hcl
resource "netbox_virtual_machine" "vm" {
  tenant = "tenant-name"              # ✅ Works
  tenant = "tenant-slug"              # ✅ Works
  tenant = netbox_tenant.test.id     # ✅ Works
  tenant = netbox_tenant.test.name   # ✅ Works
}
```

**Benefits of this change:**
- Cleaner plans: No more `(known after apply)` noise for ID fields
- Simpler state: Less duplication
- Easier maintenance: Fewer fields to manage
- Better UX: One consistent way to reference resources

### Technical Details

#### Files Modified
*   **13 resource files** - Removed 24 duplicate computed ID fields from model structs, schemas, and state mapping
*   **1 test file** - Updated platform_resource_test.go to remove manufacturer_id from expected computed fields
*   **1 test script** - Enhanced run-terraform-tests.ps1 with aggregate and IP range pre-cleanup functions
*   **204 test configurations** - Updated with terraform blocks and proper provider configuration

#### Impact Analysis
*   **High**: Configurations explicitly referencing `._id` fields must be updated
*   **Medium**: Plans become cleaner without `(known after apply)` noise
*   **Low**: State migration is fully automatic

#### Validation
*   **All 204 Terraform integration tests passing** (100% pass rate)
*   **All resource and data source types validated**
*   **Cross-resource references verified**
*   **State management and drift detection confirmed**

## v0.0.9 (2025-12-29)

### Bug Fixes

#### Fixed display_name Field Issues
*   **Removed display_name field from all 100 resources**
    - The field was showing "(known after apply)" in plans unnecessarily
    - It returned the resource's own Display value, not referenced resource names
    - Added no user value while adding plan noise and complexity
    - Removed from: all DCIM, IPAM, Circuits, Virtualization, Tenancy, and Extras resources

#### Fixed Reference Attribute Plan Display
*   **Reference attributes now preserve user-specified format**
    - Previously: `tenant = "My Tenant" -> "42"` (unwanted ID conversion)
    - Now: `tenant = "My Tenant"` (preserves name/slug as specified)
    - Affects 47 resources with reference fields (tenant, cluster, site, vlan, etc.)
    - Fix: Updated `UpdateReferenceAttribute` in state_helpers.go to prefer name/slug over ID
*   **Comprehensive test coverage**
    - 7 reference preservation tests verify no drift occurs
    - All 150+ consistency tests passing
    - Manual testing confirms correct behavior in real Terraform workflows

#### Fixed Acceptance Test Failures
*   **Fixed 18 test failures in 3 groups:**
    - Group 1: display_name schema issues (CircuitGroup, ClusterGroup, ConfigContext) - 3 tests
    - Group 2: Reference persistence (CircuitGroup, ConfigContext) - 2 tests
    - Group 3: External deletion handling (ClusterGroup, Contact, ContactAssignment) - 3 tests
*   **Improved test reliability**
    - 7 tests updated with unique email addresses for parallel execution
    - MAC address handling normalized (case-insensitive)
    - JSON data properly normalized in ConfigContext
    - 404 errors properly handled in external deletion tests

### Technical Details

#### Files Modified
*   **100 resource files** - Removed display_name field (model, schema, state mapping)
*   **1 utility file** - Enhanced UpdateReferenceAttribute function
*   **1 test file** - Added IPAddress reference preservation test
*   **7 test files** - Fixed bugs and improved test reliability

#### Breaking Changes
*   **display_name field removed** - Users should not be referencing this computed field
    - Impact: Low (field was non-functional and confusing)
    - Migration: Remove any references to `.display_name` in configurations
    - The field never worked as intended and caused plan noise

### Statistics
*   **109 files modified**
*   **All unit tests passing** (resources, datasources, utils)
*   **All 150+ consistency tests passing**
*   **18 test failures fixed**
*   **7 tests improved for reliability**

## v0.0.8 (2025-12-29)

### Major Improvements

#### Resources
*   **Complete Test Coverage:** Added external deletion tests for all 99 resources (100% coverage)
    - Tests verify graceful handling when resources are deleted outside Terraform
    - Resources are cleanly removed from state without errors
    - Comprehensive coverage across all resource types
*   **Code Refactoring:** Extracted common patterns into reusable utility functions
    - ~2,000 lines of code removed through refactoring
    - Improved maintainability and consistency
    - Better error handling patterns

#### Datasources
*   **Enhanced Error Handling:** Added 404 error handling to 21 datasources (Phase 1)
    - cable, cable_termination, circuit_group_assignment, circuit_termination
    - config_template, contact_assignment, event_rule
    - fhrp_group_assignment, interface_template, inventory_item
    - inventory_item_role, inventory_item_template, journal_entry
    - l2vpn_termination, module_bay_template, notification_group
    - rack_reservation, virtual_device_context, wireless_lan
    - wireless_lan_group, wireless_link
    - Clear, actionable error messages when resources don't exist
*   **Improved Test Coverage:** Enhanced test coverage for 32 datasources (Phase 2)
    - Split combined tests into granular, focused test functions
    - 81 new tests added (79 passed, 2 skipped)
    - Easier debugging with specific test functions for each lookup method
    - Consistent naming conventions (byID, byName, bySlug, IDPreservation)

#### Utilities & Infrastructure
*   **New Utility Functions:** Added comprehensive helper libraries
    - `internal/utils/state_helpers.go` - State management utilities (296 lines)
    - `internal/utils/request_helpers.go` - API request helpers (178 lines)
    - `internal/schema/attributes.go` - Common schema attributes (90 lines)
    - `internal/utils/state_helpers_test.go` - Comprehensive tests (444 lines)
*   **Improved Code Quality:**
    - Reduced code duplication across resources
    - Consistent error handling patterns
    - Better type safety with utility functions
    - 100% test coverage for utilities

### Statistics
*   **252 files modified** with improved quality
*   **~2,000 lines removed** through refactoring
*   **1,008 lines added** in utilities and tests
*   **11,580 insertions, 13,600 deletions** (net reduction)
*   **Test pass rate:** 97.5% (79/81 datasource tests passed)

## v0.0.3 (2025-12-16)

### Bug Fixes

*   **Custom Fields:** Fixed a panic that occurred when custom fields contained non-string values (e.g., `float64` from JSON unmarshalling). The provider now safely handles different types and converts them to strings where appropriate.
*   **Data Sources:** Fixed an issue where `display_url` was incorrectly treated as a required field in some data sources, causing errors when reading resources where this field was missing or null.

## v0.0.2 (2025-12-15)

*   Initial release with support for Netbox v4.1.11.
