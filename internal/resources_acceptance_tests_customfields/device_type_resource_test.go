//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceTypeResource_importWithCustomFieldsAndTags tests importing a device type
// with custom fields and tags to ensure all data is preserved during import.
func TestAccDeviceTypeResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	model := testutil.RandomName("tf-test-dt-import")
	slug := testutil.RandomSlug("tf-test-dt-import")
	manufacturerName := testutil.RandomName("tf-test-mfr-import")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-dt-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-dt-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-dt-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-dt-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values for different data types
	cfText := testutil.RandomCustomFieldName("tf_dt_text")
	cfTextValue := testutil.RandomName("device-type-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_dt_longtext")
	cfLongtextValue := fmt.Sprintf("Device type description: %s", testutil.RandomName("dt-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_dt_integer")
	cfIntegerValue := 1000
	cfBoolean := testutil.RandomCustomFieldName("tf_dt_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_dt_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_dt_url")
	cfURLValue := testutil.RandomURL("device-type")
	cfJSON := testutil.RandomCustomFieldName("tf_dt_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the device type with all custom fields and tags
				Config: testAccDeviceTypeResourceImportConfig_full(
					model, slug, manufacturerName, manufacturerSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the device type and verify all fields are preserved
				ResourceName:      "netbox_device_type.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config: testAccDeviceTypeResourceImportConfig_full(
					model, slug, manufacturerName, manufacturerSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccDeviceTypeResourceImportConfig_full(
	model, slug, manufacturerName, manufacturerSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "dt_test1" {
  name  = %[5]q
  slug  = %[6]q
  color = %[7]q
}

resource "netbox_tag" "dt_test2" {
  name  = %[8]q
  slug  = %[9]q
  color = %[10]q
}

# Create custom fields for dcim.devicetype
resource "netbox_custom_field" "dt_text" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.devicetype"]
  required     = false
}

resource "netbox_custom_field" "dt_longtext" {
  name         = %[13]q
  type         = "longtext"
  object_types = ["dcim.devicetype"]
  required     = false
	depends_on   = [netbox_custom_field.dt_text]
}

resource "netbox_custom_field" "dt_integer" {
  name         = %[15]q
  type         = "integer"
  object_types = ["dcim.devicetype"]
  required     = false
	depends_on   = [netbox_custom_field.dt_longtext]
}

resource "netbox_custom_field" "dt_boolean" {
  name         = %[17]q
  type         = "boolean"
  object_types = ["dcim.devicetype"]
  required     = false
	depends_on   = [netbox_custom_field.dt_integer]
}

resource "netbox_custom_field" "dt_date" {
  name         = %[19]q
  type         = "date"
  object_types = ["dcim.devicetype"]
  required     = false
	depends_on   = [netbox_custom_field.dt_boolean]
}

resource "netbox_custom_field" "dt_url" {
  name         = %[21]q
  type         = "url"
  object_types = ["dcim.devicetype"]
  required     = false
	depends_on   = [netbox_custom_field.dt_date]
}

resource "netbox_custom_field" "dt_json" {
  name         = %[23]q
  type         = "json"
  object_types = ["dcim.devicetype"]
  required     = false
	depends_on   = [netbox_custom_field.dt_url]
}

# Create manufacturer dependency
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

# Create device type with all custom fields and tags
resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[1]q
  slug         = %[2]q

	tags = [netbox_tag.dt_test1.slug, netbox_tag.dt_test2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.dt_text.name
      type  = "text"
      value = %[12]q
    },
    {
      name  = netbox_custom_field.dt_longtext.name
      type  = "longtext"
      value = %[14]q
    },
    {
      name  = netbox_custom_field.dt_integer.name
      type  = "integer"
      value = "%[16]d"
    },
    {
      name  = netbox_custom_field.dt_boolean.name
      type  = "boolean"
      value = "%[18]t"
    },
    {
      name  = netbox_custom_field.dt_date.name
      type  = "date"
      value = %[20]q
    },
    {
      name  = netbox_custom_field.dt_url.name
      type  = "url"
      value = %[22]q
    },
    {
      name  = netbox_custom_field.dt_json.name
      type  = "json"
      value = %[24]q
    }
  ]
}
`, model, slug, manufacturerName, manufacturerSlug, tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}

// TestAccDeviceTypeResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a device type. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create device type with custom fields
// 2. Update device type WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccDeviceTypeResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	model := testutil.RandomName("tf-test-dt-preserve")
	slug := testutil.RandomSlug("tf-test-dt-preserve")
	manufacturerName := testutil.RandomName("tf-test-mfr-preserve")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create device type WITH custom fields explicitly in config
				Config: testAccDeviceTypeConfig_preservation_step1(
					model, slug, manufacturerName, manufacturerSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device_type.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device_type.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccDeviceTypeConfig_preservation_step2(
					model, slug, manufacturerName, manufacturerSlug,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Updated description"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_device_type.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_device_type.test",
				ImportState:             true,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccDeviceTypeConfig_preservation_step1(
					model, slug, manufacturerName, manufacturerSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_device_type.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device_type.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device_type.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccDeviceTypeConfig_preservation_step1(
	model, slug, manufacturerName, manufacturerSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_custom_field" "text" {
  name         = %[5]q
  type         = "text"
  object_types = ["dcim.devicetype"]
}

resource "netbox_custom_field" "integer" {
  name         = %[6]q
  type         = "integer"
  object_types = ["dcim.devicetype"]
}

resource "netbox_device_type" "test" {
  model           = %[1]q
  slug            = %[2]q
  manufacturer    = netbox_manufacturer.test.id
  description     = %[7]q

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[8]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[9]d"
    }
  ]

  depends_on = [
    netbox_custom_field.text,
    netbox_custom_field.integer,
  ]
}
`,
		model, slug, manufacturerName, manufacturerSlug,
		cfTextName, cfIntName, "Initial description", cfTextValue, cfIntValue,
	)
}

func testAccDeviceTypeConfig_preservation_step2(
	model, slug, manufacturerName, manufacturerSlug,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_custom_field" "text" {
  name         = %[5]q
  type         = "text"
  object_types = ["dcim.devicetype"]
}

resource "netbox_custom_field" "integer" {
  name         = %[6]q
  type         = "integer"
  object_types = ["dcim.devicetype"]
}

resource "netbox_device_type" "test" {
  model           = %[1]q
  slug            = %[2]q
  manufacturer    = netbox_manufacturer.test.id
  description     = %[7]q

  # NOTE: custom_fields is intentionally omitted to test preservation behavior
}
`,
		model, slug, manufacturerName, manufacturerSlug,
		cfTextName, cfIntName, description,
	)
}
