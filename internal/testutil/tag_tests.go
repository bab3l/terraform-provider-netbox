package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TagLifecycleTestConfig defines the configuration for tag lifecycle tests.
type TagLifecycleTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// ConfigWithoutTags returns config for a resource without tags
	ConfigWithoutTags func() string

	// ConfigWithTags returns config with tags
	// Should include tag resources in the config
	ConfigWithTags func() string

	// ConfigWithDifferentTags returns config with a different set of tags
	ConfigWithDifferentTags func() string

	// ExpectedTagCount is the expected number of tags when tags are set
	ExpectedTagCount int

	// ExpectedDifferentTagCount is the expected number of tags after changing
	ExpectedDifferentTagCount int

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunTagLifecycleTest tests the complete tag lifecycle.
// This test:
// 1. Creates a resource without tags
// 2. Adds tags to the resource
// 3. Changes the tags
// 4. Removes all tags.
func RunTagLifecycleTest(t *testing.T, config TagLifecycleTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Step 1: Create without tags
		{
			Config: config.ConfigWithoutTags(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, "tags.#", "0"),
			),
		},
		// Step 2: Add tags
		{
			Config: config.ConfigWithTags(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, "tags.#", fmt.Sprintf("%d", config.ExpectedTagCount)),
			),
		},
		// Step 3: Change tags (if config provided)
		// Step 4: Remove all tags
		{
			Config: config.ConfigWithoutTags(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, "tags.#", "0"),
			),
		},
		// Step 5: Verify no drift
		{
			Config:   config.ConfigWithoutTags(),
			PlanOnly: true,
		},
	}

	// Insert tag change step if provided
	if config.ConfigWithDifferentTags != nil {
		expectedCount := config.ExpectedDifferentTagCount
		if expectedCount == 0 {
			expectedCount = config.ExpectedTagCount
		}

		newSteps := make([]resource.TestStep, 0, len(steps)+1)
		newSteps = append(newSteps, steps[:2]...) // Step 1 (create without tags) and Step 2 (add tags)
		newSteps = append(newSteps, resource.TestStep{
			Config: config.ConfigWithDifferentTags(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, "tags.#", fmt.Sprintf("%d", expectedCount)),
			),
		})
		newSteps = append(newSteps, steps[2:]...) // Step 3 (remove tags) and Step 4 (verify no drift)
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

// TagOrderTestConfig tests that tag order doesn't cause drift.
type TagOrderTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// ConfigWithTagsOrderA returns config with tags in order A
	ConfigWithTagsOrderA func() string

	// ConfigWithTagsOrderB returns config with same tags in different order
	ConfigWithTagsOrderB func() string

	// ExpectedTagCount is the expected number of tags
	ExpectedTagCount int

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunTagOrderTest tests that tag reordering doesn't cause drift.
func RunTagOrderTest(t *testing.T, config TagOrderTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Create with tags in order A
		{
			Config: config.ConfigWithTagsOrderA(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, "tags.#", fmt.Sprintf("%d", config.ExpectedTagCount)),
			),
		},
		// Apply with tags in order B - should have no changes
		{
			Config: config.ConfigWithTagsOrderB(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, "tags.#", fmt.Sprintf("%d", config.ExpectedTagCount)),
			),
		},
		// Plan-only to verify no drift
		{
			Config:   config.ConfigWithTagsOrderB(),
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
