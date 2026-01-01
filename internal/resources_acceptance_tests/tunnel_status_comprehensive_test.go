package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccTunnelResource_StatusComprehensive tests comprehensive scenarios for tunnel status field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccTunnelResource_StatusComprehensive(t *testing.T) {
	// Generate unique names for this test run
	tunnelName := testutil.RandomName("tf-test-tunnel-status")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_tunnel",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTunnelDestroy,
			testutil.CheckTunnelGroupDestroy,
		),
		BaseConfig: func() string {
			return `
resource "netbox_tunnel" "test" {
	name       = "` + tunnelName + `"
	encapsulation = "gre"
	# status field intentionally omitted - should get default "active"
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_tunnel" "test" {
	name       = "` + tunnelName + `"
	encapsulation = "gre"
	status     = "` + value + `"
}
`
		},
	})
}
