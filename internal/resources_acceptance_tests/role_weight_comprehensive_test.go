package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccRoleResource_Weight tests comprehensive scenarios for role weight field.
// This validates that Optional+Computed int64 fields with proper defaults work correctly.
func TestAccRoleResource_Weight(t *testing.T) {
	// Generate unique names for this test run
	roleName := testutil.RandomName("tf-test-role-weight")
	roleSlug := testutil.RandomSlug("tf-test-role-weight")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_role",
		OptionalField:  "weight",
		DefaultValue:   "1000",
		FieldTestValue: "2000",
		CheckDestroy:   testutil.CheckRoleDestroy,
		BaseConfig: func() string {
			return `
resource "netbox_role" "test" {
	name = "` + roleName + `"
	slug = "` + roleSlug + `"
	# weight field intentionally omitted - should get default 1000
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_role" "test" {
	name = "` + roleName + `"
	slug = "` + roleSlug + `"
	weight = ` + value + `
}
`
		},
	})
}
