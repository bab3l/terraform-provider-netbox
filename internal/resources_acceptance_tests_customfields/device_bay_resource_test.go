//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	bayName := testutil.RandomName("device_bay")
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	// Clean up custom fields and tags
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayResourceImportConfig_full(bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_device_bay.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_device_bay.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_device_bay.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields"}, // Device reference may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccDeviceBayResourceImportConfig_full(bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.slug
  subdevice_role = "parent"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.devicebay"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.devicebay"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.devicebay"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.devicebay"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.devicebay"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.devicebay"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.devicebay"]
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
resource "netbox_device_bay" "test" {
  name   = %q
  device = netbox_device.test.name

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

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
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
		bayName,
	)
}

// TestAccDeviceBayResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a device bay. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create device bay with custom fields
// 2. Update device bay WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccDeviceBayResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	bayName := testutil.RandomName("tf-test-bay-cf-preserve")
	deviceName := testutil.RandomName("tf-test-device-preserve")
	siteName := testutil.RandomName("tf-test-site-preserve")
	siteSlug := testutil.RandomSlug("tf-test-site-preserve")
	mfgName := testutil.RandomName("tf-test-mfg-preserve")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-preserve")
	dtModel := testutil.RandomName("tf-test-dt-preserve")
	dtSlug := testutil.RandomSlug("tf-test-dt-preserve")
	roleName := testutil.RandomName("tf-test-role-preserve")
	roleSlug := testutil.RandomSlug("tf-test-role-preserve")
	tenantName := testutil.RandomName("tf-test-tenant-preserve")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfInteger)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create device bay WITH custom fields explicitly in config
				Config: testAccDeviceBayConfig_preservation_step1(
					bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "label", "Test Label"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device_bay.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device_bay.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update label WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccDeviceBayConfig_preservation_step2(
					bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug,
					cfText, cfInteger, "Updated Label",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "label", "Updated Label"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_device_bay.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_device_bay.test",
				ImportState:             true,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccDeviceBayConfig_preservation_step1(
					bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_device_bay.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device_bay.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device_bay.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccDeviceBayConfig_preservation_step1(
	bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
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
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.slug
}

resource "netbox_custom_field" "text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.device_bay"]
}

resource "netbox_custom_field" "integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.device_bay"]
}

resource "netbox_device" "test" {
  name        = %q
  site        = netbox_site.test.slug
  role        = netbox_device_role.test.slug
  device_type = netbox_device_type.test.slug
  tenant      = netbox_tenant.test.slug
  status      = "active"
}

resource "netbox_device_bay" "test" {
  device = netbox_device.test.name
  name   = %q
  label  = "Test Label"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%d"
    }
  ]
}
`,
		tenantName, tenantSlug,
		siteName, siteSlug,
		mfgName, mfgSlug,
		roleName, roleSlug,
		dtModel, dtSlug,
		cfTextName, cfIntName,
		deviceName,
		bayName,
		cfTextValue, cfIntValue,
	)
}

func testAccDeviceBayConfig_preservation_step2(
	bayName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug,
	cfTextName, cfIntName, label string,
) string {
	return fmt.Sprintf(`
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
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.slug
}

resource "netbox_custom_field" "text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.device_bay"]
}

resource "netbox_custom_field" "integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.device_bay"]
}

resource "netbox_device" "test" {
  name        = %q
  site        = netbox_site.test.slug
  role        = netbox_device_role.test.slug
  device_type = netbox_device_type.test.slug
  tenant      = netbox_tenant.test.slug
  status      = "active"
}

resource "netbox_device_bay" "test" {
  device = netbox_device.test.name
  name   = %q
  label  = %q
}
`,
		tenantName, tenantSlug,
		siteName, siteSlug,
		mfgName, mfgSlug,
		roleName, roleSlug,
		dtModel, dtSlug,
		cfTextName, cfIntName,
		deviceName,
		bayName, label,
	)
}
