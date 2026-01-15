package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// OptionalFieldTestConfig generates test configurations for optional field scenarios.
type OptionalFieldTestConfig struct {
	ResourceName    string                         // e.g., "netbox_device"
	ResourceType    string                         // e.g., "device"
	OptionalField   string                         // e.g., "status"
	BaseConfig      func() string                  // Function that returns base config without optional field
	WithFieldConfig func(fieldValue string) string // Function that returns config with optional field
	FieldTestValue  string                         // Value to use when testing with field present
	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// GenerateOptionalFieldTests creates the standard set of optional field tests.
func GenerateOptionalFieldTests(t *testing.T, config OptionalFieldTestConfig) []resource.TestStep {
	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	return []resource.TestStep{
		// Test 1: Create without optional field - should not crash or show drift
		{
			Config: config.BaseConfig(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				// Optional field should not be set in state when not in config
				resource.TestCheckNoResourceAttr(resourceRef, config.OptionalField),
			),
		},
		// Test 2: Plan-only verification - no changes should be detected
		{
			PlanOnly: true,
			Config:   config.BaseConfig(),
		},
		// Test 3: Add optional field to existing resource
		{
			Config: config.WithFieldConfig(config.FieldTestValue),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, config.OptionalField, config.FieldTestValue),
			),
		},
		// Test 4: Remove optional field from config - should not crash
		{
			Config: config.BaseConfig(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				// Field should be removed/null when not in config
				resource.TestCheckNoResourceAttr(resourceRef, config.OptionalField),
			),
		},
		// Test 5: Final plan-only verification - no drift
		{
			PlanOnly: true,
			Config:   config.BaseConfig(),
		},
	}
}

// ImportOptionalFieldTest creates a test for import scenarios with optional fields.
func ImportOptionalFieldTest(t *testing.T, config OptionalFieldTestConfig) []resource.TestStep {
	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	return []resource.TestStep{
		// Create resource with optional field
		{
			Config: config.WithFieldConfig(config.FieldTestValue),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
			),
		},
		// Import resource
		{
			ResourceName:      resourceRef,
			ImportState:       true,
			ImportStateVerify: true,
		},
		// Apply base config (without optional field) - should not crash
		{
			Config: config.BaseConfig(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
			),
		},
	}
}

// OptionalFieldTestSuite runs all optional field tests for a resource.
func RunOptionalFieldTestSuite(t *testing.T, config OptionalFieldTestConfig) {
	t.Parallel()

	// Test suite for optional field handling
	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    GenerateOptionalFieldTests(t, config),
	}

	// Add CheckDestroy if provided
	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}

// OptionalFieldImportTestSuite runs import-specific tests.
func RunOptionalFieldImportTestSuite(t *testing.T, config OptionalFieldTestConfig) {
	t.Parallel()

	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    ImportOptionalFieldTest(t, config),
	}

	// Add CheckDestroy if provided
	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}

// MultiFieldOptionalTestConfig defines configuration for testing multiple optional fields together.
type MultiFieldOptionalTestConfig struct {
	ResourceName string // e.g., "netbox_aggregate"
	ResourceType string // e.g., "aggregate"

	// BaseConfig returns Terraform config without optional fields
	BaseConfig func() string

	// ConfigWithFields returns Terraform config with all optional fields populated
	ConfigWithFields func() string

	// OptionalFields maps field names to their expected values when set
	OptionalFields map[string]string

	// RequiredFields maps required field names to their expected values (for stability checks)
	RequiredFields map[string]string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// TestRemoveOptionalFields is a comprehensive test that verifies multiple optional fields
// can be removed from a resource configuration without causing "inconsistent result" errors.
//
// This test ensures that when optional fields are removed from Terraform configuration:
// 1. Fields are properly cleared in the API (no stale values)
// 2. State reflects the removal (fields are null/empty)
// 3. No "Provider produced inconsistent result after apply" errors occur
// 4. Fields can be re-added successfully
//
// Example usage:
//
//	config := testutil.MultiFieldOptionalTestConfig{
//	    ResourceName: "netbox_aggregate",
//	    BaseConfig: func() string {
//	        return `resource "netbox_aggregate" "test" {
//	            prefix = "10.0.0.0/8"
//	            rir = "rfc-1918"
//	        }`
//	    },
//	    ConfigWithFields: func() string {
//	        return `resource "netbox_aggregate" "test" {
//	            prefix = "10.0.0.0/8"
//	            rir = "rfc-1918"
//	            description = "Test description"
//	            comments = "Test comments"
//	        }`
//	    },
//	    OptionalFields: map[string]string{
//	        "description": "Test description",
//	        "comments": "Test comments",
//	    },
//	    RequiredFields: map[string]string{
//	        "prefix": "10.0.0.0/8",
//	    },
//	}
//	testutil.TestRemoveOptionalFields(t, config)
func TestRemoveOptionalFields(t *testing.T, config MultiFieldOptionalTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Step 1: Create resource with all optional fields populated
		{
			Config: config.ConfigWithFields(),
			Check:  buildMultiFieldChecks(resourceRef, config.OptionalFields, config.RequiredFields, true),
		},
		// Step 2: Remove optional fields and verify they are cleared (critical test)
		{
			Config: config.BaseConfig(),
			Check:  buildMultiFieldChecks(resourceRef, config.OptionalFields, config.RequiredFields, false),
		},
		// Step 3: Re-add optional fields to verify they can be set again
		{
			Config: config.ConfigWithFields(),
			Check:  buildMultiFieldChecks(resourceRef, config.OptionalFields, config.RequiredFields, true),
		},
	}

	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	}

	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}

// buildMultiFieldChecks creates test check functions for verifying multiple optional fields.
func buildMultiFieldChecks(resourceName string, optionalFields, requiredFields map[string]string, fieldsPresent bool) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc

	// Always verify required fields remain stable
	for field, value := range requiredFields {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, field, value))
	}

	if fieldsPresent {
		// Verify optional fields are set to expected values
		for field, value := range optionalFields {
			checks = append(checks, resource.TestCheckResourceAttr(resourceName, field, value))
		}
	} else {
		// Verify optional fields are removed/null (critical check)
		for field := range optionalFields {
			checks = append(checks, resource.TestCheckNoResourceAttr(resourceName, field))
		}
	}

	return resource.ComposeTestCheckFunc(checks...)
}
