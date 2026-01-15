package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// HierarchicalTestConfig defines the configuration for testing hierarchical resources.
// Hierarchical resources have parent-child relationships (e.g., Region, Location, ContactGroup).
type HierarchicalTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_region")
	ResourceName string

	// ParentField is the name of the parent attribute (e.g., "parent")
	ParentField string

	// ConfigWithoutParent returns config for a resource without a parent
	ConfigWithoutParent func() string

	// ConfigWithParent returns config for a resource with a parent
	// The parent resource should be created in the same config
	ConfigWithParent func() string

	// ConfigWithDifferentParent returns config moving the resource to a different parent
	ConfigWithDifferentParent func() string

	// ExpectedParentValue is the expected value of the parent attribute when set
	ExpectedParentValue string

	// ExpectedNewParentValue is the expected value after changing parents
	ExpectedNewParentValue string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunHierarchicalTest executes tests for hierarchical resource relationships.
// This test:
// 1. Creates a resource without a parent
// 2. Adds a parent relationship
// 3. Changes the parent to a different resource
// 4. Removes the parent relationship.
func RunHierarchicalTest(t *testing.T, config HierarchicalTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Step 1: Create without parent
		{
			Config: config.ConfigWithoutParent(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckNoResourceAttr(resourceRef, config.ParentField),
			),
		},
		// Step 2: Add parent
		{
			Config: config.ConfigWithParent(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttrSet(resourceRef, config.ParentField),
			),
		},
		// Step 3: Change to different parent (if config provided)
		// Step 4: Remove parent
		{
			Config: config.ConfigWithoutParent(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckNoResourceAttr(resourceRef, config.ParentField),
			),
		},
		// Step 5: Verify no drift
		{
			Config:   config.ConfigWithoutParent(),
			PlanOnly: true,
		},
	}

	// Insert parent change step if provided
	if config.ConfigWithDifferentParent != nil {
		newSteps := make([]resource.TestStep, 0, len(steps)+1)
		newSteps = append(newSteps, steps[:3]...)
		newSteps = append(newSteps, resource.TestStep{
			Config: config.ConfigWithDifferentParent(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttrSet(resourceRef, config.ParentField),
			),
		})
		newSteps = append(newSteps, steps[3:]...)
		steps = newSteps
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

// NestedHierarchyTestConfig tests deeply nested hierarchies.
type NestedHierarchyTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_region")
	ResourceName string

	// ParentField is the name of the parent attribute
	ParentField string

	// ConfigWithNestedHierarchy returns config with multiple levels of nesting
	// (e.g., grandparent -> parent -> child)
	ConfigWithNestedHierarchy func() string

	// Depth is the expected nesting depth (e.g., 3 for grandparent->parent->child)
	Depth int

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunNestedHierarchyTest tests that deeply nested hierarchies work correctly.
func RunNestedHierarchyTest(t *testing.T, config NestedHierarchyTestConfig) {
	t.Helper()

	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config.ConfigWithNestedHierarchy(),
				Check: resource.ComposeTestCheckFunc(
					// Verify the deepest child resource exists
					resource.TestCheckResourceAttrSet(fmt.Sprintf("%s.child", config.ResourceName), "id"),
				),
			},
			// Plan-only to verify stability
			{
				Config:   config.ConfigWithNestedHierarchy(),
				PlanOnly: true,
			},
		},
	}

	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}
