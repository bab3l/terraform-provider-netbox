package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// OptionalComputedFieldTestConfig defines the configuration for testing Optional+Computed fields.
// These fields differ from Optional Only fields because they're always present in state
// with their default values, even when not specified in configuration.
type OptionalComputedFieldTestConfig struct {
	// Resource type name (e.g., "netbox_device")
	ResourceName string

	// The field being tested (e.g., "status")
	OptionalField string

	// The default value that should appear when field is not in config
	DefaultValue string

	// A test value different from the default to verify field updates work
	FieldTestValue string

	// BaseConfig returns Terraform config without the optional field
	BaseConfig func() string

	// WithFieldConfig returns Terraform config with the optional field set to the given value
	WithFieldConfig func(value string) string
}

// GenerateOptionalComputedFieldTests creates tests for Optional+Computed fields.
// These fields are always present in state with their default values, even when
// not specified in configuration.
func GenerateOptionalComputedFieldTests(t *testing.T, config OptionalComputedFieldTestConfig) []resource.TestStep {
	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	return []resource.TestStep{
		// Test 1: Create without optional field - should have default value
		{
			Config: config.BaseConfig(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				// Optional+Computed field should have default value when not in config
				resource.TestCheckResourceAttr(resourceRef, config.OptionalField, config.DefaultValue),
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
		// Test 4: Remove optional field from config - should revert to default
		{
			Config: config.BaseConfig(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				// Field should revert to default when removed from config
				resource.TestCheckResourceAttr(resourceRef, config.OptionalField, config.DefaultValue),
			),
		},
		// Test 5: Final plan-only verification - no drift
		{
			PlanOnly: true,
			Config:   config.BaseConfig(),
		},
	}
}

// ImportOptionalComputedFieldTest creates a test for import scenarios with Optional+Computed fields.
func ImportOptionalComputedFieldTest(t *testing.T, config OptionalComputedFieldTestConfig) []resource.TestStep {
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
	}
}

// RunOptionalComputedFieldTestSuite runs the complete test suite for an Optional+Computed field.
func RunOptionalComputedFieldTestSuite(t *testing.T, config OptionalComputedFieldTestConfig) {
	t.Helper()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    GenerateOptionalComputedFieldTests(t, config),
	})
}
