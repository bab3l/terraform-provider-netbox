//go:build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccWirelessLANResource_CustomFieldsPreservation tests that custom fields
// are preserved when updating other resource attributes without custom_fields in config.
// This is the core bug fix test - ensuring no data loss on updates.
func TestAccWirelessLANResource_CustomFieldsPreservation(t *testing.T) {
	randomSSID := fmt.Sprintf("test-ssid-%s", acctest.RandString(8))
	cf1Name := fmt.Sprintf("test_cf1_%s", acctest.RandString(8))
	cf2Name := fmt.Sprintf("test_cf2_%s", acctest.RandString(8))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with 2 custom fields
			{
				Config: testAccWirelessLANResourceConfig_withCustomFields(randomSSID, "Initial description", cf1Name, cf2Name, "value1", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", randomSSID),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "custom_fields.#", "2"),
				),
			},
			// Step 2: Update description WITHOUT custom_fields in config
			// BUG FIX TEST: Custom fields should be preserved in NetBox, not visible in state
			{
				Config: testAccWirelessLANResourceConfig_noCustomFields(randomSSID, "Updated description", cf1Name, cf2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", randomSSID),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", "Updated description"),
					// Custom fields omitted from config, so not in state
					resource.TestCheckNoResourceAttr("netbox_wireless_lan.test", "custom_fields.#"),
				),
			},
			// Step 3: Re-add custom_fields to config - verify both fields still exist
			{
				Config: testAccWirelessLANResourceConfig_withCustomFields(randomSSID, "Updated description", cf1Name, cf2Name, "value1", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", randomSSID),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "custom_fields.#", "2"),
				),
			},
		},
	})
}

// Test helper config generators

func testAccWirelessLANResourceConfig_base(cf1Name, cf2Name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test1" {
  name         = "%s"
  object_types = ["wireless.wirelesslan"]
  type         = "text"
}

resource "netbox_custom_field" "test2" {
  name         = "%s"
  object_types = ["wireless.wirelesslan"]
  type         = "text"
}
`, cf1Name, cf2Name)
}

func testAccWirelessLANResourceConfig_withCustomFields(ssid, description, cf1Name, cf2Name, val1, val2 string) string {
	return testAccWirelessLANResourceConfig_base(cf1Name, cf2Name) + fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
  ssid        = "%s"
  description = "%s"

  custom_fields = [
    {
      name  = netbox_custom_field.test1.name
      type  = "text"
      value = "%s"
    },
    {
      name  = netbox_custom_field.test2.name
      type  = "text"
      value = "%s"
    }
  ]
}
`, ssid, description, val1, val2)
}

func testAccWirelessLANResourceConfig_noCustomFields(ssid, description, cf1Name, cf2Name string) string {
	return testAccWirelessLANResourceConfig_base(cf1Name, cf2Name) + fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
  ssid        = "%s"
  description = "%s"
  # custom_fields intentionally omitted - should preserve existing values in NetBox
}
`, ssid, description)
}
