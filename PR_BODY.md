## Description

This PR fixes two failing datasource acceptance tests in the IDPreservation test suite:

### Issues Fixed

1. **TestAccConsoleServerPortDataSource_IDPreservation**
   - **Error**: Data source reference mismatch - test referenced non-existent `data.netbox_console_server_port.test`
   - **Fix**: Updated test to reference the correct datasource created by config: `data.netbox_console_server_port.by_id`

2. **TestAccModuleBayDataSource_IDPreservation**
   - **Error**: Attribute name mismatch - test expected hardcoded name "Module Bay 1" but received random device name
   - **Fix**: Updated test to use the `moduleBayName` variable generated in the test setup, ensuring the test expects the actual name that will be created

### Testing

- ✅ Resource acceptance tests: 99/99 passed (4758 seconds)
- ✅ Datasource acceptance tests: Both fixes validated
- ✅ All pre-commit checks passing

### Files Changed

- `internal/datasources_acceptance_tests/console_server_port_data_source_test.go`
- `internal/datasources_acceptance_tests/module_bay_data_source_test.go`

This completes the ID preservation test coverage for all 203 Terraform provider resources and datasources (99 resources + 104 datasources).
