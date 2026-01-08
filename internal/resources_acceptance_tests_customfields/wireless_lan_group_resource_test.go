//go:build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccWirelessLANGroupResource_CustomFieldsPreservation tests that custom fields
// are preserved when updating other resource attributes without custom_fields in config.
func TestAccWirelessLANGroupResource_CustomFieldsPreservation(t *testing.T) {
	randomName := fmt.Sprintf("test-wlan-group-%s", acctest.RandString(8))
	randomSlug := fmt.Sprintf("test-wlan-group-%s", acctest.RandString(8))
	cf1Name := fmt.Sprintf("test_cf1_%s", acctest.RandString(8))
	cf2Name := fmt.Sprintf("test_cf2_%s", acctest.RandString(8))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with 2 custom fields
			{
				Config: testAccWirelessLANGroupResourceConfig_withCustomFields(randomName, randomSlug, "Initial description", cf1Name, cf2Name, "value1", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "custom_fields.#", "2"),
				),
			},
			// Step 2: Update description WITHOUT custom_fields in config
			{
				Config: testAccWirelessLANGroupResourceConfig_noCustomFields(randomName, randomSlug, "Updated description", cf1Name, cf2Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", "Updated description"),
					resource.TestCheckNoResourceAttr("netbox_wireless_lan_group.test", "custom_fields.#"),
				),
			},
			// Step 3: Re-add custom_fields to config - verify both fields still exist
			{
				Config: testAccWirelessLANGroupResourceConfig_withCustomFields(randomName, randomSlug, "Updated description", cf1Name, cf2Name, "value1", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "custom_fields.#", "2"),
				),
			},
		},
	})
}

// Test helper config generators

func testAccWirelessLANGroupResourceConfig_base(cf1Name, cf2Name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test1" {
  name         = "%s"
  object_types = ["wireless.wirelesslangroup"]
  type         = "text"
}

resource "netbox_custom_field" "test2" {
  name         = "%s"
  object_types = ["wireless.wirelesslangroup"]
  type         = "text"
}
`, cf1Name, cf2Name)
}

func testAccWirelessLANGroupResourceConfig_withCustomFields(name, slug, description, cf1Name, cf2Name, val1, val2 string) string {
	return testAccWirelessLANGroupResourceConfig_base(cf1Name, cf2Name) + fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name        = "%s"
  slug        = "%s"
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
`, name, slug, description, val1, val2)
}

func testAccWirelessLANGroupResourceConfig_noCustomFields(name, slug, description, cf1Name, cf2Name string) string {
	return testAccWirelessLANGroupResourceConfig_base(cf1Name, cf2Name) + fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name        = "%s"
  slug        = "%s"
  description = "%s"
  # custom_fields intentionally omitted - should preserve existing values in NetBox
}
`, name, slug, description)
}
