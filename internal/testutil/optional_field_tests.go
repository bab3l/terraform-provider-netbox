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
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    GenerateOptionalFieldTests(t, config),
	})
}

// OptionalFieldImportTestSuite runs import-specific tests.
func RunOptionalFieldImportTestSuite(t *testing.T, config OptionalFieldTestConfig) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    ImportOptionalFieldTest(t, config),
	})
}
