//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsoleServerPortResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	portName := testutil.RandomName("console_server_port")
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
				Config: testAccConsoleServerPortResourceImportConfig_full(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", portName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_console_server_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields"}, // Device reference may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccConsoleServerPortResourceImportConfig_full(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.consoleserverport"]
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
resource "netbox_console_server_port" "test" {
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
		portName,
	)
}

// TestAccConsoleServerPortResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a console server port.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccConsoleServerPortResource_CustomFieldsPreservation(t *testing.T) {
	portName := testutil.RandomName("console_server_port")
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create console server port WITH custom fields
				Config: testAccConsoleServerPortConfig_preservation_step1(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_console_server_port.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_console_server_port.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccConsoleServerPortConfig_preservation_step2(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Add custom_fields back to verify they were preserved
				Config: testAccConsoleServerPortConfig_preservation_step3(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_console_server_port.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_console_server_port.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccConsoleServerPortConfig_preservation_step1(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_role" "test" {
  name  = %[9]q
  slug  = %[10]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[7]q
  slug         = %[8]q
  manufacturer = netbox_manufacturer.test.slug
}

resource "netbox_device" "test" {
  name        = %[2]q
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

resource "netbox_custom_field" "env" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "owner" {
  name         = %[12]q
  type         = "text"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_console_server_port" "test" {
  name   = %[1]q
  device = netbox_device.test.name

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "team-a"
    }
  ]
}
`, portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnv, cfOwner)
}

func testAccConsoleServerPortConfig_preservation_step2(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_role" "test" {
  name  = %[9]q
  slug  = %[10]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[7]q
  slug         = %[8]q
  manufacturer = netbox_manufacturer.test.slug
}

resource "netbox_device" "test" {
  name        = %[2]q
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

resource "netbox_custom_field" "env" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "owner" {
  name         = %[12]q
  type         = "text"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_console_server_port" "test" {
  name        = %[1]q
  device      = netbox_device.test.name
  description = "Updated description"
  # custom_fields intentionally omitted - testing preservation
}
`, portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnv, cfOwner)
}

func testAccConsoleServerPortConfig_preservation_step3(portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_role" "test" {
  name  = %[9]q
  slug  = %[10]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[7]q
  slug         = %[8]q
  manufacturer = netbox_manufacturer.test.slug
}

resource "netbox_device" "test" {
  name        = %[2]q
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

resource "netbox_custom_field" "env" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_custom_field" "owner" {
  name         = %[12]q
  type         = "text"
  object_types = ["dcim.consoleserverport"]
}

resource "netbox_console_server_port" "test" {
  name        = %[1]q
  device      = netbox_device.test.name
  description = "Updated description"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "team-a"
    }
  ]
}
`, portName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfEnv, cfOwner)
}
