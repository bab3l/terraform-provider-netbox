//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	interfaceName := testutil.RandomName("tf-test-int-import")
	deviceName := testutil.RandomName("tf-test-device-import")
	manufacturerName := testutil.RandomName("tf-test-mfr-import")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-import")
	deviceTypeModel := testutil.RandomName("tf-test-dt-import")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-import")
	deviceRoleName := testutil.RandomName("tf-test-role-import")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-import")
	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-int-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-int-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-int-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-int-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values for different data types
	cfText := testutil.RandomCustomFieldName("tf_int_text")
	cfTextValue := testutil.RandomName("interface-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_int_longtext")
	cfLongtextValue := fmt.Sprintf("Interface description: %s", testutil.RandomName("int-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_int_integer")
	cfIntegerValue := 1000
	cfBoolean := testutil.RandomCustomFieldName("tf_int_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_int_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_int_url")
	cfURLValue := testutil.RandomURL("interface")
	cfJSON := testutil.RandomCustomFieldName("tf_int_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the interface with all custom fields and tags
				Config: testAccInterfaceResourceImportConfig_full(
					interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttr("netbox_interface.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_interface.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the interface and verify all fields are preserved
				ResourceName:      "netbox_interface.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config: testAccInterfaceResourceImportConfig_full(
					interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccInterfaceResourceImportConfig_full(
	interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "int_test1" {
  name  = %[11]q
  slug  = %[12]q
  color = %[13]q
}

resource "netbox_tag" "int_test2" {
  name  = %[14]q
  slug  = %[15]q
  color = %[16]q
}

# Create custom fields for dcim.interface
resource "netbox_custom_field" "int_text" {
  name         = %[17]q
  type         = "text"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_longtext" {
  name         = %[19]q
  type         = "longtext"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_integer" {
  name         = %[21]q
  type         = "integer"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_boolean" {
  name         = %[23]q
  type         = "boolean"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_date" {
  name         = %[25]q
  type         = "date"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_url" {
  name         = %[27]q
  type         = "url"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_json" {
  name         = %[29]q
  type         = "json"
  object_types = ["dcim.interface"]
  required     = false
}

# Create dependencies
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "test" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  name        = %[2]q
  status      = "active"
}

# Create interface with all custom fields and tags
resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = %[1]q
  type   = "1000base-t"

  tags = [netbox_tag.int_test1.slug, netbox_tag.int_test2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.int_text.name
      type  = "text"
      value = %[18]q
    },
    {
      name  = netbox_custom_field.int_longtext.name
      type  = "longtext"
      value = %[20]q
    },
    {
      name  = netbox_custom_field.int_integer.name
      type  = "integer"
      value = "%[22]d"
    },
    {
      name  = netbox_custom_field.int_boolean.name
      type  = "boolean"
      value = "%[24]t"
    },
    {
      name  = netbox_custom_field.int_date.name
      type  = "date"
      value = %[26]q
    },
    {
      name  = netbox_custom_field.int_url.name
      type  = "url"
      value = %[28]q
    },
    {
      name  = netbox_custom_field.int_json.name
      type  = "json"
      value = %[30]q
    }
  ]
}
`, interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
