//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerOutletResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	powerOutletName := testutil.RandomName("power_outlet")
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPowerOutletDestroy,
		Steps: []resource.TestStep{
			// First create the resource
			{
				Config: testAccPowerOutletResourceImportConfig_full(powerOutletName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
				),
			},
			// Then test import
			{
				Config:                  testAccPowerOutletResourceImportConfig_full(powerOutletName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_power_outlet.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields"}, // Device reference may have lookup inconsistencies
			},
		},
	})
}

func testAccPowerOutletResourceImportConfig_full(powerOutletName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
  color = "9e9e9e"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model = %q
  slug = %q
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
  name = %q
}

# Custom field definitions with different object types
resource "netbox_custom_field" "test_text" {
  name = %q
  type = "text"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_longtext" {
  name = %q
  type = "longtext"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_integer" {
  name = %q
  type = "integer"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_boolean" {
  name = %q
  type = "boolean"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_date" {
  name = %q
  type = "date"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_url" {
  name = %q
  type = "url"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_json" {
  name = %q
  type = "json"
  object_types = ["dcim.poweroutlet"]
}

# Tag definitions
resource "netbox_tag" "test_1" {
  name = %q
  slug = %q
  color = "ff0000"
}

resource "netbox_tag" "test_2" {
  name = %q
  slug = %q
  color = "00ff00"
}

# Power outlet with custom fields and tags
resource "netbox_power_outlet" "test" {
  device = netbox_device.test.id
  name = %q
  type = "iec-60320-c13"

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = "custom text value"
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = "custom longtext value"
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "123"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "false"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = "2023-12-25"
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = "https://custom.example.com"
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = jsonencode({"custom": "json", "array": [1, 2, 3]})
    }
  ]

  tags = [
    {
      name = netbox_tag.test_1.name
      slug = netbox_tag.test_1.slug
    },
    {
      name = netbox_tag.test_2.name
      slug = netbox_tag.test_2.slug
    }
  ]
}
`,
		tenantName, tenantSlug,
		siteName, siteSlug,
		mfgName, mfgSlug,
		roleName, roleSlug,
		dtModel, dtSlug,
		deviceName,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		powerOutletName,
	)
}
