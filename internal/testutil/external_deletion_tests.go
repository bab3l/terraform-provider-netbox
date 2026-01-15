package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ExternalDeletionTestConfig defines the configuration for external deletion tests.
// These tests verify that Terraform properly handles resources deleted outside of Terraform.
type ExternalDeletionTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// Config returns the Terraform configuration that creates the resource
	Config func() string

	// DeleteFunc is called to delete the resource externally (via API)
	// It receives the resource ID and should return an error if deletion fails
	DeleteFunc func(ctx context.Context, id string) error

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunExternalDeletionTest executes an external deletion test.
// This test:
// 1. Creates the resource via Terraform
// 2. Deletes the resource externally (via API)
// 3. Runs terraform plan/apply - should detect and recreate.
func RunExternalDeletionTest(t *testing.T, config ExternalDeletionTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)
	var resourceID string

	steps := []resource.TestStep{
		// Step 1: Create the resource
		{
			Config: config.Config(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				// Capture the resource ID for external deletion
				func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources[resourceRef]
					if !ok {
						return fmt.Errorf("resource %s not found in state", resourceRef)
					}
					resourceID = rs.Primary.ID
					return nil
				},
			),
		},
		// Step 2: Delete externally and verify Terraform recreates
		{
			PreConfig: func() {
				if resourceID == "" {
					t.Fatal("resource ID not captured from previous step")
				}
				ctx := context.Background()
				if err := config.DeleteFunc(ctx, resourceID); err != nil {
					t.Fatalf("failed to delete resource externally: %s", err)
				}
			},
			Config: config.Config(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
			),
		},
		// Step 3: Verify no drift after recreation
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

// ExternalDeletionWithIDTestConfig is a variant that captures ID differently.
type ExternalDeletionWithIDTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// Config returns the Terraform configuration that creates the resource
	Config func() string

	// GetIDFromState extracts the API ID from Terraform state
	// Default behavior is to use Primary.ID
	GetIDFromState func(rs *terraform.ResourceState) (string, error)

	// DeleteFunc is called to delete the resource externally (via API)
	DeleteFunc func(ctx context.Context, id string) error

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunExternalDeletionWithIDTest is similar to RunExternalDeletionTest but allows
// custom ID extraction from state.
func RunExternalDeletionWithIDTest(t *testing.T, config ExternalDeletionWithIDTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)
	var resourceID string

	getID := config.GetIDFromState
	if getID == nil {
		getID = func(rs *terraform.ResourceState) (string, error) {
			return rs.Primary.ID, nil
		}
	}

	steps := []resource.TestStep{
		// Step 1: Create the resource
		{
			Config: config.Config(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources[resourceRef]
					if !ok {
						return fmt.Errorf("resource %s not found in state", resourceRef)
					}
					id, err := getID(rs)
					if err != nil {
						return err
					}
					resourceID = id
					return nil
				},
			),
		},
		// Step 2: Delete externally and verify Terraform recreates
		{
			PreConfig: func() {
				if resourceID == "" {
					t.Fatal("resource ID not captured from previous step")
				}
				ctx := context.Background()
				if err := config.DeleteFunc(ctx, resourceID); err != nil {
					t.Fatalf("failed to delete resource externally: %s", err)
				}
			},
			Config: config.Config(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
			),
		},
		// Step 3: Verify no drift after recreation
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
