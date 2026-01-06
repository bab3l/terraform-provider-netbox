package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a device. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create device with custom fields in NetBox API (not via Terraform)
// 2. Import device into Terraform WITHOUT custom_fields in config
// 3. Update device (change description, etc.)
// 4. Custom fields should be preserved, not deleted.
func TestAccDeviceResource_CustomFieldsPreservation(t *testing.T) {
	t.Parallel()

	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-preserve")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names - must be consistent across steps
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create device WITH custom fields explicitly in config
				Config: testAccDeviceResourceConfig_withCustomFields(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					// Verify the custom field values
					testutil.CheckCustomFieldValue("netbox_device.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// This simulates the real-world scenario where a user manages device configs
				// but not custom fields (which may be managed externally or manually)
				Config: testAccDeviceResourceConfig_withoutCustomFields(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Updated description"),
					// CRITICAL: Custom fields should still exist in NetBox even though not in config
					// This test will FAIL if the bug exists - custom fields will be deleted
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 3: Update description again to verify fields remain stable
				Config: testAccDeviceResourceConfig_withoutCustomFields(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfText, cfInteger, "Second update",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Second update"),
					// Custom fields should STILL be present
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
				),
			},
		},
	})
}

func testAccDeviceResourceConfig_withCustomFields(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "integer" {
  name         = %[11]q
  type         = "integer"
  object_types = ["dcim.device"]
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model        = %[4]q
  slug         = %[5]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "ff0000"
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  description = "Initial description"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[12]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[13]d"
    }
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfTextName, cfIntName, cfTextValue, cfIntValue,
	)
}

func testAccDeviceResourceConfig_withoutCustomFields(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "integer" {
  name         = %[11]q
  type         = "integer"
  object_types = ["dcim.device"]
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model        = %[4]q
  slug         = %[5]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "ff0000"
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  description = %[12]q

  # NOTE: custom_fields intentionally omitted to test preservation
  # In real-world usage, custom fields might be managed outside Terraform

  # Keep dependencies alive
  depends_on = [
    netbox_custom_field.text,
    netbox_custom_field.integer
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfTextName, cfIntName, description,
	)
}
