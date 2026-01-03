package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceTypeResource_basic(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("tf-test-device-type")

	slug := testutil.RandomSlug("tf-test-dt")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
				),
			},
			{
				// Test import
				ResourceName:            "netbox_device_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})

}

func TestAccDeviceTypeResource_IDPreservation(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("dt-id")
	slug := testutil.GenerateSlug(model)
	manufacturerName := testutil.RandomName("mfr-dt")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckDeviceTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_full(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-full")

	slug := testutil.RandomSlug("tf-test-dt-full")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "part_number", "TEST-PART-001"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "2"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "is_full_depth", "true"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "airflow", "front-to-rear"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Test device type with full options"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "comments", "Test comments for device type"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_update(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-update")

	slug := testutil.RandomSlug("tf-test-dt-upd")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	updatedModel := testutil.RandomName("tf-test-device-type-updated")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},

			{

				Config: testAccDeviceTypeResourceConfig_updated(updatedModel, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", updatedModel),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "4"),
				),
			},
		},
	})

}

func testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

}

`, manufacturerName, manufacturerSlug, model, slug)

}

func testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer      = netbox_manufacturer.test.id

  model             = %q

  slug              = %q

  part_number       = "TEST-PART-001"

  u_height          = 2

  is_full_depth     = true

  airflow           = "front-to-rear"

  description       = "Test device type with full options"

  comments          = "Test comments for device type"

}

`, manufacturerName, manufacturerSlug, model, slug)

}

func testAccDeviceTypeResourceConfig_updated(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

  u_height     = 4

  description  = "Updated description"

}

`, manufacturerName, manufacturerSlug, model, slug)

}

func TestAccConsistency_DeviceType_LiteralNames(t *testing.T) {
	t.Parallel()
	model := testutil.RandomName("tf-test-device-type-lit")
	slug := testutil.RandomSlug("tf-test-dt-lit")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-lit")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckDeviceTypeDestroy, testutil.CheckManufacturerDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", description),
				),
			},
			{
				Config:   testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  description  = %q
}
`, manufacturerName, manufacturerSlug, model, slug, description)
}

func TestAccDeviceTypeResource_externalDeletion(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("test-device-type-del")
	slug := testutil.GenerateSlug(model)
	manufacturerName := testutil.RandomName("test-manufacturer")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDeviceTypesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device_type for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDeviceTypesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device_type: %v", err)
					}
					t.Logf("Successfully externally deleted device_type with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccDeviceTypeResource_importWithCustomFieldsAndTags tests importing a device type
// with custom fields and tags to ensure all data is preserved during import.
func TestAccDeviceTypeResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

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
				ResourceName:            "netbox_device_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags", "manufacturer"},
				// The import should preserve basic fields
				Check: resource.ComposeTestCheckFunc(
					// Verify basic fields
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
				),
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
}

resource "netbox_custom_field" "dt_integer" {
  name         = %[15]q
  type         = "integer"
  object_types = ["dcim.devicetype"]
  required     = false
}

resource "netbox_custom_field" "dt_boolean" {
  name         = %[17]q
  type         = "boolean"
  object_types = ["dcim.devicetype"]
  required     = false
}

resource "netbox_custom_field" "dt_date" {
  name         = %[19]q
  type         = "date"
  object_types = ["dcim.devicetype"]
  required     = false
}

resource "netbox_custom_field" "dt_url" {
  name         = %[21]q
  type         = "url"
  object_types = ["dcim.devicetype"]
  required     = false
}

resource "netbox_custom_field" "dt_json" {
  name         = %[23]q
  type         = "json"
  object_types = ["dcim.devicetype"]
  required     = false
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

  tags = [
    {
      name = netbox_tag.dt_test1.name
      slug = netbox_tag.dt_test1.slug
    },
    {
      name = netbox_tag.dt_test2.name
      slug = netbox_tag.dt_test2.slug
    }
  ]

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
