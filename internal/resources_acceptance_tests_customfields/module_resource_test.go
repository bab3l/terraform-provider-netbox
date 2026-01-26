//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	deviceName := testutil.RandomName("device")
	bayName := testutil.RandomName("module_bay")
	mtModel := testutil.RandomName("module_type")
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
		CheckDestroy:             testutil.CheckModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "device"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "module_bay"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "module_type"),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_module.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "2"),
				),
			},
			{
				Config:            testAccModuleResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:      "netbox_module.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccModuleResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccModuleResourceImportConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.module"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.module"]

  depends_on = [netbox_custom_field.cf_text]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.module"]

  depends_on = [netbox_custom_field.cf_longtext]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.module"]

  depends_on = [netbox_custom_field.cf_integer]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.module"]

  depends_on = [netbox_custom_field.cf_boolean]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.module"]

  depends_on = [netbox_custom_field.cf_date]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.module"]

  depends_on = [netbox_custom_field.cf_url]
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
resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id

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
		tenantName, tenantSlug,
		siteName, siteSlug,
		mfgName, mfgSlug,
		roleName, roleSlug,
		dtModel, dtSlug,
		deviceName,
		bayName,
		mtModel,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
	)
}

// TestAccModuleResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a module.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccModuleResource_CustomFieldsPreservation(t *testing.T) {
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	deviceName := testutil.RandomName("device")
	bayName := testutil.RandomName("module_bay")
	mtModel := testutil.RandomName("module_type")

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckModuleDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create module WITH custom fields
				Config: testAccModuleConfig_preservation_step1(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttr("netbox_module.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_module.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_module.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update comments WITHOUT mentioning custom_fields
				Config: testAccModuleConfig_preservation_step2(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "comments", "Updated comments"),
					resource.TestCheckResourceAttr("netbox_module.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Add custom_fields back to verify they were preserved
				Config: testAccModuleConfig_preservation_step3(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_module.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_module.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_module.test", "comments", "Updated comments"),
				),
			},
		},
	})
}

func testAccModuleConfig_preservation_step1(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[11]q
}

resource "netbox_custom_field" "env" {
  name         = %[12]q
  type         = "text"
  object_types = ["dcim.module"]
}

resource "netbox_custom_field" "owner" {
  name         = %[13]q
  type         = "text"
  object_types = ["dcim.module"]
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id

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
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnv, cfOwner)
}

func testAccModuleConfig_preservation_step2(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[11]q
}

resource "netbox_custom_field" "env" {
  name         = %[12]q
  type         = "text"
  object_types = ["dcim.module"]
}

resource "netbox_custom_field" "owner" {
  name         = %[13]q
  type         = "text"
  object_types = ["dcim.module"]
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  comments    = "Updated comments"
  # custom_fields intentionally omitted - testing preservation
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnv, cfOwner)
}

func testAccModuleConfig_preservation_step3(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[11]q
}

resource "netbox_custom_field" "env" {
  name         = %[12]q
  type         = "text"
  object_types = ["dcim.module"]
}

resource "netbox_custom_field" "owner" {
  name         = %[13]q
  type         = "text"
  object_types = ["dcim.module"]
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  comments    = "Updated comments"

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
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, cfEnv, cfOwner)
}
