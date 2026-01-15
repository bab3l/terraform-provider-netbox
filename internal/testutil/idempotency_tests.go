package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// IdempotencyTestConfig defines the configuration for idempotency tests.
// Idempotency tests verify that multiple applies produce no changes.
type IdempotencyTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// Config returns the Terraform configuration to test
	Config func() string

	// NumApplies is the number of times to apply the configuration (default: 3)
	NumApplies int

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc

	// AdditionalChecks are extra checks to run after each apply
	AdditionalChecks []resource.TestCheckFunc
}

// RunIdempotencyTest executes an idempotency test.
// This test:
// 1. Applies the configuration
// 2. Verifies no changes on subsequent plans
// 3. Repeats for the specified number of applies.
func RunIdempotencyTest(t *testing.T, config IdempotencyTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	numApplies := config.NumApplies
	if numApplies <= 0 {
		numApplies = 3
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceRef, "id"),
	}
	checks = append(checks, config.AdditionalChecks...)

	steps := make([]resource.TestStep, 0, numApplies*2)

	for i := 0; i < numApplies; i++ {
		// Apply step
		steps = append(steps, resource.TestStep{
			Config: config.Config(),
			Check:  resource.ComposeTestCheckFunc(checks...),
		})
		// Plan-only verification
		steps = append(steps, resource.TestStep{
			Config:   config.Config(),
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

// RefreshIdempotencyTestConfig tests idempotency after refresh.
type RefreshIdempotencyTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// Config returns the Terraform configuration to test
	Config func() string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunRefreshIdempotencyTest tests that refresh followed by plan shows no changes.
func RunRefreshIdempotencyTest(t *testing.T, config RefreshIdempotencyTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Initial apply
		{
			Config: config.Config(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
			),
		},
		// Refresh and verify no changes
		{
			Config:   config.Config(),
			PlanOnly: true,
			// RefreshState triggers a refresh before planning
			RefreshState: true,
		},
		// Another refresh cycle
		{
			Config:       config.Config(),
			PlanOnly:     true,
			RefreshState: true,
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
