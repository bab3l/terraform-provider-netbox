package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ReferenceChangeTestConfig defines the configuration for testing reference attribute changes.
// Reference attributes are foreign key relationships to other resources.
type ReferenceChangeTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// ReferenceField is the name of the reference attribute (e.g., "site", "tenant")
	ReferenceField string

	// ConfigWithReference returns config with the reference set to a value
	ConfigWithReference func() string

	// ConfigWithDifferentReference returns config with the reference changed to a different value
	ConfigWithDifferentReference func() string

	// ConfigWithoutReference returns config with the reference removed (if optional)
	// Set to nil if the reference is required
	ConfigWithoutReference func() string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunReferenceChangeTest tests changing reference attributes.
// This test:
// 1. Creates a resource with a reference
// 2. Changes the reference to a different value
// 3. Optionally removes the reference (if optional)
// 4. Verifies no drift after changes.
func RunReferenceChangeTest(t *testing.T, config ReferenceChangeTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		// Step 1: Create with initial reference
		{
			Config: config.ConfigWithReference(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttrSet(resourceRef, config.ReferenceField),
			),
		},
		// Step 2: Change to different reference
		{
			Config: config.ConfigWithDifferentReference(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttrSet(resourceRef, config.ReferenceField),
			),
		},
		// Step 3: Verify no drift
		{
			Config:   config.ConfigWithDifferentReference(),
			PlanOnly: true,
		},
	}

	// Add optional reference removal test
	if config.ConfigWithoutReference != nil {
		steps = append(steps,
			resource.TestStep{
				Config: config.ConfigWithoutReference(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRef, "id"),
					resource.TestCheckNoResourceAttr(resourceRef, config.ReferenceField),
				),
			},
			resource.TestStep{
				Config:   config.ConfigWithoutReference(),
				PlanOnly: true,
			},
		)
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

// MultiReferenceTestConfig tests resources with multiple reference attributes.
type MultiReferenceTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// ReferenceFields lists all reference attributes to test
	ReferenceFields []string

	// ConfigAllReferences returns config with all references set
	ConfigAllReferences func() string

	// ConfigNoOptionalReferences returns config with only required references
	ConfigNoOptionalReferences func() string

	// RequiredReferences lists which references are required (cannot be removed)
	RequiredReferences []string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunMultiReferenceTest tests resources with multiple reference attributes.
func RunMultiReferenceTest(t *testing.T, config MultiReferenceTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	// Build checks for all references present
	var allRefChecks []resource.TestCheckFunc
	allRefChecks = append(allRefChecks, resource.TestCheckResourceAttrSet(resourceRef, "id"))
	for _, field := range config.ReferenceFields {
		allRefChecks = append(allRefChecks, resource.TestCheckResourceAttrSet(resourceRef, field))
	}

	// Build checks for only required references
	var requiredRefChecks []resource.TestCheckFunc
	requiredRefChecks = append(requiredRefChecks, resource.TestCheckResourceAttrSet(resourceRef, "id"))

	isRequired := make(map[string]bool)
	for _, field := range config.RequiredReferences {
		isRequired[field] = true
	}

	for _, field := range config.ReferenceFields {
		if isRequired[field] {
			requiredRefChecks = append(requiredRefChecks, resource.TestCheckResourceAttrSet(resourceRef, field))
		} else {
			requiredRefChecks = append(requiredRefChecks, resource.TestCheckNoResourceAttr(resourceRef, field))
		}
	}

	steps := []resource.TestStep{
		// Step 1: Create with all references
		{
			Config: config.ConfigAllReferences(),
			Check:  resource.ComposeTestCheckFunc(allRefChecks...),
		},
		// Step 2: Remove optional references
		{
			Config: config.ConfigNoOptionalReferences(),
			Check:  resource.ComposeTestCheckFunc(requiredRefChecks...),
		},
		// Step 3: Re-add all references
		{
			Config: config.ConfigAllReferences(),
			Check:  resource.ComposeTestCheckFunc(allRefChecks...),
		},
		// Step 4: Verify no drift
		{
			Config:   config.ConfigAllReferences(),
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
