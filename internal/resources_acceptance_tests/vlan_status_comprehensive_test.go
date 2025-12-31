package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccVlanResource_StatusOptionalField tests comprehensive scenarios for VLAN status.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccVlanResource_StatusOptionalField(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_vlan",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "deprecated",
		BaseConfig: func() string {
			return `
resource "netbox_vlan" "test" {
	name = "test-vlan-status"
	vid  = 100
	# status field intentionally omitted - should get default "active"
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_vlan" "test" {
	name   = "test-vlan-status"
	vid    = 100
	status = "` + value + `"
}
`
		},
	})
}
