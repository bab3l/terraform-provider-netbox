# Changelog

## v0.0.3 (2025-12-16)

### Bug Fixes

*   **Custom Fields:** Fixed a panic that occurred when custom fields contained non-string values (e.g., `float64` from JSON unmarshalling). The provider now safely handles different types and converts them to strings where appropriate.
*   **Data Sources:** Fixed an issue where `display_url` was incorrectly treated as a required field in some data sources, causing errors when reading resources where this field was missing or null.

## v0.0.2 (2025-12-15)

*   Initial release with support for Netbox v4.1.11.
