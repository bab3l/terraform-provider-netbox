package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// UpdateTestConfig defines the configuration for resource update tests.
type UpdateTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// InitialConfig returns the initial Terraform configuration
	InitialConfig func() string

	// UpdatedConfig returns the updated Terraform configuration
	UpdatedConfig func() string

	// InitialChecks are checks to run after initial creation
	InitialChecks []resource.TestCheckFunc

	// UpdatedChecks are checks to run after the update
	UpdatedChecks []resource.TestCheckFunc

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunUpdateTest executes a resource update test.
// This test:
// 1. Creates the resource with initial config
// 2. Updates the resource with new config
// 3. Verifies no drift after update.
func RunUpdateTest(t *testing.T, config UpdateTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	// Build initial checks
	initialChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceRef, "id"),
	}
	initialChecks = append(initialChecks, config.InitialChecks...)

	// Build updated checks
	updatedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceRef, "id"),
	}
	updatedChecks = append(updatedChecks, config.UpdatedChecks...)

	steps := []resource.TestStep{
		// Step 1: Create initial resource
		{
			Config: config.InitialConfig(),
			Check:  resource.ComposeTestCheckFunc(initialChecks...),
		},
		// Step 2: Update the resource
		{
			Config: config.UpdatedConfig(),
			Check:  resource.ComposeTestCheckFunc(updatedChecks...),
		},
		// Step 3: Verify no drift
		{
			Config:   config.UpdatedConfig(),
			PlanOnly: true,
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

// MultiStepUpdateTestConfig tests multiple sequential updates.
type MultiStepUpdateTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// Steps defines the sequence of configurations and checks
	Steps []UpdateStep

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// UpdateStep represents a single update step in a multi-step update test.
type UpdateStep struct {
	// Name is a descriptive name for this step (for debugging)
	Name string

	// Config returns the Terraform configuration for this step
	Config func() string

	// Checks are verifications to run after applying this config
	Checks []resource.TestCheckFunc
}

// RunMultiStepUpdateTest executes a test with multiple update steps.
func RunMultiStepUpdateTest(t *testing.T, config MultiStepUpdateTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	var steps []resource.TestStep

	for i, step := range config.Steps {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttrSet(resourceRef, "id"),
		}
		checks = append(checks, step.Checks...)

		steps = append(steps, resource.TestStep{
			Config: step.Config(),
			Check:  resource.ComposeTestCheckFunc(checks...),
		})

		// Add plan-only verification after each step (except the last)
		if i < len(config.Steps)-1 {
			steps = append(steps, resource.TestStep{
				Config:   step.Config(),
				PlanOnly: true,
			})
		}
	}

	// Final drift check
	if len(config.Steps) > 0 {
		lastStep := config.Steps[len(config.Steps)-1]
		steps = append(steps, resource.TestStep{
			Config:   lastStep.Config(),
			PlanOnly: true,
		})
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

// FieldUpdateTestConfig tests updating a specific field.
type FieldUpdateTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// FieldName is the name of the field being updated
	FieldName string

	// BaseConfig returns config without the field or with default value
	BaseConfig func() string

	// ConfigWithValue returns config with the field set to the given value
	ConfigWithValue func(value string) string

	// InitialValue is the initial value to set
	InitialValue string

	// UpdatedValue is the value to update to
	UpdatedValue string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunFieldUpdateTest tests updating a specific field through various values.
func RunFieldUpdateTest(t *testing.T, config FieldUpdateTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Create with initial value
		{
			Config: config.ConfigWithValue(config.InitialValue),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, config.FieldName, config.InitialValue),
			),
		},
		// Update to new value
		{
			Config: config.ConfigWithValue(config.UpdatedValue),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, config.FieldName, config.UpdatedValue),
			),
		},
		// Verify no drift
		{
			Config:   config.ConfigWithValue(config.UpdatedValue),
			PlanOnly: true,
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
