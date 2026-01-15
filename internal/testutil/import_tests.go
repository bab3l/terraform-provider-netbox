package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ImportTestConfig defines the configuration for standardized import tests.
type ImportTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// Config returns the Terraform configuration that creates the resource
	Config func() string

	// ImportStateIDFunc optionally generates a custom import ID.
	// If nil, the resource's "id" attribute is used.
	ImportStateIDFunc func(state *terraform.InstanceState) (string, error)

	// ImportStateVerifyIgnore lists attributes to ignore during import verification.
	// Use this for attributes that differ between config and imported state
	// (e.g., password fields, computed-only fields that aren't in API response).
	ImportStateVerifyIgnore []string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc

	// AdditionalChecks are extra checks to run after import
	AdditionalChecks []resource.TestCheckFunc
}

// RunImportTest executes a standardized import test.
// This test:
// 1. Creates the resource
// 2. Imports the resource
// 3. Verifies the imported state matches the original
// 4. Applies the config again to verify no drift.
func RunImportTest(t *testing.T, config ImportTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceRef, "id"),
	}
	checks = append(checks, config.AdditionalChecks...)

	steps := []resource.TestStep{
		// Step 1: Create the resource
		{
			Config: config.Config(),
			Check:  resource.ComposeTestCheckFunc(checks...),
		},
		// Step 2: Import the resource
		{
			ResourceName:            resourceRef,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: config.ImportStateVerifyIgnore,
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				if config.ImportStateIDFunc != nil {
					rs, ok := s.RootModule().Resources[resourceRef]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceRef)
					}
					return config.ImportStateIDFunc(rs.Primary)
				}
				// Default: use the resource ID
				rs, ok := s.RootModule().Resources[resourceRef]
				if !ok {
					return "", fmt.Errorf("resource %s not found in state", resourceRef)
				}
				return rs.Primary.ID, nil
			},
		},
		// Step 3: Apply config again - should have no changes
		{
			Config:   config.Config(),
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

// SimpleImportTestConfig is a simplified version for resources with straightforward import.
type SimpleImportTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// Config returns the Terraform configuration that creates the resource
	Config func() string

	// ImportStateVerifyIgnore lists attributes to ignore during import verification.
	ImportStateVerifyIgnore []string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunSimpleImportTest executes a basic import test without custom ID handling.
func RunSimpleImportTest(t *testing.T, config SimpleImportTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Step 1: Create the resource
		{
			Config: config.Config(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
			),
		},
		// Step 2: Import the resource
		{
			ResourceName:            resourceRef,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: config.ImportStateVerifyIgnore,
		},
		// Step 3: Apply config again - should have no changes
		{
			Config:   config.Config(),
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
