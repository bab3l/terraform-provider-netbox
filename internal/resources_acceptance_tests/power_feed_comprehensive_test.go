package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccPowerFeedResource_VoltageAmperage tests comprehensive scenarios for power feed voltage and amperage fields.
// This validates that Optional+Computed numeric fields work correctly across all scenarios.
func TestAccPowerFeedResource_VoltageAmperage(t *testing.T) {

	// Test voltage field
	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_feed",
		OptionalField:  "voltage",
		DefaultValue:   "120",
		FieldTestValue: "240",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPowerFeedDestroy,
			testutil.CheckPowerPanelDestroy,
			testutil.CheckLocationDestroy,
			testutil.CheckSiteDestroy,
		),
		BaseConfig: func() string {
			return `
resource "netbox_site" "test" {
	name = "test-site-power-feed"
	slug = "test-site-power-feed"
}

resource "netbox_power_panel" "test" {
	name = "test-panel-power-feed"
	site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
	name         = "test-power-feed-voltage"
	power_panel  = netbox_power_panel.test.id
	# voltage field intentionally omitted - should get default 120
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_site" "test" {
	name = "test-site-power-feed"
	slug = "test-site-power-feed"
}

resource "netbox_power_panel" "test" {
	name = "test-panel-power-feed"
	site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
	name         = "test-power-feed-voltage"
	power_panel  = netbox_power_panel.test.id
	voltage      = ` + value + `
}
`
		},
	})
}

// TestAccPowerFeedResource_Amperage tests the amperage field separately.
func TestAccPowerFeedResource_Amperage(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_feed",
		OptionalField:  "amperage",
		DefaultValue:   "20",
		FieldTestValue: "30",
		BaseConfig: func() string {
			return `
resource "netbox_site" "test" {
	name = "test-site-power-feed-amp"
	slug = "test-site-power-feed-amp"
}

resource "netbox_power_panel" "test" {
	name = "test-panel-power-feed-amp"
	site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
	name         = "test-power-feed-amperage"
	power_panel  = netbox_power_panel.test.id
	# amperage field intentionally omitted - should get default 20
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_site" "test" {
	name = "test-site-power-feed-amp"
	slug = "test-site-power-feed-amp"
}

resource "netbox_power_panel" "test" {
	name = "test-panel-power-feed-amp"
	site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
	name         = "test-power-feed-amperage"
	power_panel  = netbox_power_panel.test.id
	amperage     = ` + value + `
}
`
		},
	})
}
