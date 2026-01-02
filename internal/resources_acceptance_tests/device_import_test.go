package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceResource_importWithCustomFieldsAndTags tests importing a pre-existing device
// that has various custom field types and tags properly imports all data.
func TestAccDeviceResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-import")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfTextValue := testutil.RandomName("text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfLongtextValue := fmt.Sprintf("This is a longer text field with multiple words: %s", testutil.RandomName("longtext"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_integer")
	cfIntegerValue := 12345
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfURLValue := testutil.RandomURL("device")
	cfJSON := testutil.RandomCustomFieldName("tf_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the device with all custom fields and tags
				Config: testAccDeviceResourceImportConfig_full(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the device and verify all fields are preserved
				ResourceName:      "netbox_device.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The import should preserve all custom fields and tags
				Check: resource.ComposeTestCheckFunc(
					// Verify basic fields
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),

					// Verify tags are imported
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "2"),

					// Verify custom fields are imported
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "7"),

					// Verify specific custom field values
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.0.name", cfText),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.0.value", cfTextValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.1.name", cfLongtext),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.1.type", "longtext"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.1.value", cfLongtextValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.2.name", cfIntegerName),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.2.type", "integer"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.2.value", fmt.Sprintf("%d", cfIntegerValue)),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.3.name", cfBoolean),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.3.type", "boolean"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.3.value", fmt.Sprintf("%t", cfBooleanValue)),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.4.name", cfDate),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.4.type", "date"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.4.value", cfDateValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.5.name", cfURL),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.5.type", "url"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.5.value", cfURLValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.6.name", cfJSON),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.6.type", "json"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.6.value", cfJSONValue),
				),
			},
		},
	})
}

func testAccDeviceResourceImportConfig_full(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "test1" {
  name  = %[10]q
  slug  = %[11]q
  color = %[12]q
}

resource "netbox_tag" "test2" {
  name  = %[13]q
  slug  = %[14]q
  color = %[15]q
}

# Create custom fields
resource "netbox_custom_field" "text" {
  name         = %[16]q
  type         = "text"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "longtext" {
  name         = %[18]q
  type         = "longtext"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "integer" {
  name         = %[20]q
  type         = "integer"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "boolean" {
  name         = %[22]q
  type         = "boolean"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "date" {
  name         = %[24]q
  type         = "date"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "url" {
  name         = %[26]q
  type         = "url"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "json" {
  name         = %[28]q
  type         = "json"
  object_types = ["dcim.device"]
  required     = false
}

# Create dependencies
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

# Create device with all custom fields and tags
resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"

  tags = [
    {
      name = netbox_tag.test1.name
      slug = netbox_tag.test1.slug
    },
    {
      name = netbox_tag.test2.name
      slug = netbox_tag.test2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[17]q
    },
    {
      name  = netbox_custom_field.longtext.name
      type  = "longtext"
      value = %[19]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[21]d"
    },
    {
      name  = netbox_custom_field.boolean.name
      type  = "boolean"
      value = "%[23]t"
    },
    {
      name  = netbox_custom_field.date.name
      type  = "date"
      value = %[25]q
    },
    {
      name  = netbox_custom_field.url.name
      type  = "url"
      value = %[27]q
    },
    {
      name  = netbox_custom_field.json.name
      type  = "json"
      value = %[29]q
    }
  ]
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName,
		tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
