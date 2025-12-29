# Changelog

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
