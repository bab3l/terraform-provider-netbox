//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRearPortResource_CustomFieldsPreservation verifies that custom fields
// set outside of Terraform are preserved when updating other resource fields.
func TestAccRearPortResource_CustomFieldsPreservation(t *testing.T) {
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	dtModel := testutil.RandomName("dt")
	dtSlug := testutil.RandomSlug("dt")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	deviceName := testutil.RandomName("device")
	rearPortName := testutil.RandomName("rearport")
	cfText := testutil.RandomCustomFieldName("cf_text")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterCustomFieldCleanup(cfText)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortResourceConfig_withCustomField(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, "initial_value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccRearPortResourceConfig_withoutCustomFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "description", "updated description"),
					// Custom fields omitted from config - should be preserved as empty in state
					resource.TestCheckResourceAttr("netbox_rear_port.test", "custom_fields.#", "0"),
				),
			},
		},
	})
}

func TestAccRearPortResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	dtModel := testutil.RandomName("dt")
	dtSlug := testutil.RandomSlug("dt")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	deviceName := testutil.RandomName("device")
	rearPortName := testutil.RandomName("rearport")

	// Custom field names
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_rear_port.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_rear_port.test", "tags.#", "2"),
				),
			},
			{
				Config:            testAccRearPortResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:      "netbox_rear_port.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccRearPortResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccRearPortResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name        = %q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.rearport"]

  depends_on = [netbox_custom_field.cf_text]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["dcim.rearport"]

  depends_on = [netbox_custom_field.cf_longtext]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.rearport"]

  depends_on = [netbox_custom_field.cf_integer]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["dcim.rearport"]

  depends_on = [netbox_custom_field.cf_boolean]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["dcim.rearport"]

  depends_on = [netbox_custom_field.cf_date]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["dcim.rearport"]

  depends_on = [netbox_custom_field.cf_url]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["dcim.rearport"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Main Resource
resource "netbox_rear_port" "test" {
  device      = netbox_device.test.id
  name        = %q
  type        = "lc"
  positions   = 4

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
    }
  ]

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`,
		siteName, siteSlug,
		mfgName, mfgSlug,
		dtModel, dtSlug,
		roleName, roleSlug,
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
		rearPortName,
	)
}

// Helper functions for custom field tests
func testAccRearPortResourceConfig_withCustomField(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, cfValue string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name        = %q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.rearport"]
}

resource "netbox_rear_port" "test" {
  device    = netbox_device.test.id
  name      = %q
  type      = "lc"
  positions = 4

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = %q
    }
  ]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, cfText, rearPortName, cfValue)
}

func testAccRearPortResourceConfig_withoutCustomFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, cfText, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name        = %q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.rearport"]
}

resource "netbox_rear_port" "test" {
  device      = netbox_device.test.id
  name        = %q
  type        = "lc"
  positions   = 4
  description = %q
  # custom_fields intentionally omitted - should preserve existing values in NetBox
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, cfText, rearPortName, description)
}
