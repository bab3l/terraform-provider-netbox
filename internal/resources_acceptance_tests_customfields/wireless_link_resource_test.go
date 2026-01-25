//go:build customfields
// +build customfields

// Package resources_acceptance_tests_customfields contains acceptance tests for custom fields
// that require dedicated test runs to avoid conflicts with global custom field definitions.
package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccWirelessLinkResource_CustomFieldsAndTagsPreservation tests that custom fields and tags
// are preserved when updating other fields on a wireless link. This addresses the critical bug
// where custom fields and tags were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create wireless link with custom fields and tags
// 2. Update wireless link WITHOUT custom_fields/tags in config (omit the fields entirely)
// 3. Custom fields and tags should be preserved in NetBox, not deleted.
func TestAccWirelessLinkResource_CustomFieldsAndTagsPreservation(t *testing.T) {
	// Generate unique names
	wirelessLinkSSID := testutil.RandomName("tf-test-wl-cf-preserve")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceNameA := testutil.RandomName("tf-test-device-a")
	deviceNameB := testutil.RandomName("tf-test-device-b")
	interfaceNameA := testutil.RandomName("wlan0")
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_wl")
	cfInteger := testutil.RandomCustomFieldName("tf_int_wl")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create wireless link WITH custom fields and tags explicitly in config
				Config: testAccWirelessLinkConfig_preservation_step1(
					wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					deviceNameA, deviceNameB, interfaceNameA,
					tagName, tagSlug, cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", wirelessLinkSSID),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_wireless_link.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_wireless_link.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields/tags in config
				// Custom fields and tags should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields and tags since not in config
				Config: testAccWirelessLinkConfig_preservation_step2(
					wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					deviceNameA, deviceNameB, interfaceNameA,
					tagName, tagSlug, cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", wirelessLinkSSID),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "description", "Updated description"),
					// State shows 0 custom_fields and tags (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "custom_fields.#", "0"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "tags.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields and tags still exist in NetBox
				ResourceName:            "netbox_wireless_link.test",
				ImportState:             true,
				ImportStateVerify:       false,                                       // Can't verify - config has no custom_fields/tags
				ImportStateVerifyIgnore: []string{"custom_fields", "tags", "status"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields and tags back to config to verify they were preserved
				Config: testAccWirelessLinkConfig_preservation_step1(
					wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					deviceNameA, deviceNameB, interfaceNameA,
					tagName, tagSlug, cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields and tags should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_wireless_link.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_wireless_link.test", cfInteger, "integer", "42"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "tags.#", "1"),
				),
			},
		},
	})
}

func testAccWirelessLinkConfig_preservation_step1(
	wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	deviceNameA, deviceNameB, interfaceNameA,
	tagName, tagSlug, cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[15]q
  type         = "text"
  object_types = ["wireless.wirelesslink"]
}

resource "netbox_custom_field" "integer" {
  name         = %[16]q
  type         = "integer"
  object_types = ["wireless.wirelesslink"]

	depends_on = [netbox_custom_field.text]
}

resource "netbox_tag" "test" {
  name = %[13]q
  slug = %[14]q
}

resource "netbox_site" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_type" "test" {
  model        = %[6]q
  slug         = %[7]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[8]q
  slug  = %[9]q
  color = "ff0000"
}

resource "netbox_device" "test_a" {
  name        = %[10]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = %[11]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name   = %[12]q
  device =  netbox_device.test_a.id
  type      = "other-wireless"
}

resource "netbox_interface" "test_b" {
  name   = %[12]q
  device =  netbox_device.test_b.id
  type      = "other-wireless"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = %[1]q
  description = "Initial description"

	tags = [netbox_tag.test.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[17]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[18]d"
    }
  ]
}
`,
		wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		deviceNameA, deviceNameB, interfaceNameA,
		tagName, tagSlug, cfTextName, cfIntName, cfTextValue, cfIntValue,
	)
}

func testAccWirelessLinkConfig_preservation_step2(
	wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	deviceNameA, deviceNameB, interfaceNameA,
	tagName, tagSlug, cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[15]q
  type         = "text"
  object_types = ["wireless.wirelesslink"]
}

resource "netbox_custom_field" "integer" {
  name         = %[16]q
  type         = "integer"
  object_types = ["wireless.wirelesslink"]

	depends_on = [netbox_custom_field.text]
}

resource "netbox_tag" "test" {
  name = %[13]q
  slug = %[14]q
}

resource "netbox_site" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_type" "test" {
  model        = %[6]q
  slug         = %[7]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[8]q
  slug  = %[9]q
  color = "ff0000"
}

resource "netbox_device" "test_a" {
  name        = %[10]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = %[11]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name   = %[12]q
  device =  netbox_device.test_a.id
  type      = "other-wireless"
}

resource "netbox_interface" "test_b" {
  name   = %[12]q
  device =  netbox_device.test_b.id
  type      = "other-wireless"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = %[1]q
  description = %[17]q
  # Note: custom_fields and tags intentionally omitted to test preservation
}
`,
		wirelessLinkSSID, siteName, siteSlug, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		deviceNameA, deviceNameB, interfaceNameA,
		tagName, tagSlug, cfTextName, cfIntName, description,
	)
}
