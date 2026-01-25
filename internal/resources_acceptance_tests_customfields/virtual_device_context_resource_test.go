//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualDeviceContextResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	vdcName := testutil.RandomName("vdc")
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
	// Note: VDC cleanup function requires ID (int32), but we only have name - skipping VDC cleanup
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
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceImportConfig_full(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:      "netbox_virtual_device_context.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccVirtualDeviceContextResourceImportConfig_full(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccVirtualDeviceContextResourceImportConfig_full(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  color = "9e9e9e"
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

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["dcim.virtualdevicecontext"]
  depends_on = [netbox_custom_field.cf_text]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.virtualdevicecontext"]
  depends_on = [netbox_custom_field.cf_longtext]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["dcim.virtualdevicecontext"]
  depends_on = [netbox_custom_field.cf_integer]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["dcim.virtualdevicecontext"]
  depends_on = [netbox_custom_field.cf_boolean]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["dcim.virtualdevicecontext"]
  depends_on = [netbox_custom_field.cf_date]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["dcim.virtualdevicecontext"]
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

# Virtual Device Context with comprehensive custom fields and tags
resource "netbox_virtual_device_context" "test" {
  name   = %q
  device = netbox_device.test.id
  status = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this virtual device context resource for testing purposes."
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
      value = "2023-01-15"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key": "value"})
    }
  ]

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, tenantName, tenantSlug, siteName, siteSlug, mfgName, mfgSlug, roleName, roleSlug, dtModel, dtSlug, deviceName, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, vdcName)
}

// TestAccVirtualDeviceContextResource_CustomFieldsPreservation tests that custom fields are preserved
// when they are not specified in the configuration after creation.
func TestAccVirtualDeviceContextResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	vdcName := testutil.RandomName("vdc")
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfInteger)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with custom fields
			{
				Config: testAccVirtualDeviceContextResourceConfig_withCustomFields(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.#", "2"),
				),
			},
			// Step 2: Update without custom_fields in config - should preserve in NetBox
			{
				Config: testAccVirtualDeviceContextResourceConfig_withoutCustomFields(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					// Custom fields should be null/empty in state (filter-to-owned)
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Add custom fields back - should see them again
			{
				Config: testAccVirtualDeviceContextResourceConfig_withCustomFields(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.#", "2"),
				),
			},
		},
	})
}

func testAccVirtualDeviceContextResourceConfig_withCustomFields(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfText, cfInteger string) string {
	return fmt.Sprintf(`
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

resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_virtual_device_context" "test" {
  name   = %q
  device = netbox_device.test.id
  status = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    }
  ]
}
`, siteName, siteSlug, mfgName, mfgSlug, roleName, roleSlug, dtModel, dtSlug, deviceName, cfText, cfInteger, vdcName)
}

func testAccVirtualDeviceContextResourceConfig_withoutCustomFields(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, cfText, cfInteger string) string {
	return fmt.Sprintf(`
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

resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_virtual_device_context" "test" {
  name   = %q
  device = netbox_device.test.id
  status = "active"

  # custom_fields intentionally omitted - should preserve in NetBox
}
`, siteName, siteSlug, mfgName, mfgSlug, roleName, roleSlug, dtModel, dtSlug, deviceName, cfText, cfInteger, vdcName)
}
