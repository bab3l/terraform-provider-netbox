package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccRoleResource_Weight tests comprehensive scenarios for role weight field.
// This validates that Optional+Computed int64 fields with proper defaults work correctly.
func TestAccRoleResource_Weight(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_role",
		OptionalField:  "weight",
		DefaultValue:   "1000",
		FieldTestValue: "2000",
		BaseConfig: func() string {
			return `
resource "netbox_role" "test" {
	name = "role-weight-test"
	slug = "role-weight-test"
	# weight field intentionally omitted - should get default 1000
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_role" "test" {
	name = "role-weight-test"
	slug = "role-weight-test"
	weight = ` + value + `
}
`
		},
	})
}
