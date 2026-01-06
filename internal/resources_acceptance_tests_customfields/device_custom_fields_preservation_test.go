//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

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

// TestAccDeviceResource_CustomFieldsPartialManagement tests that users can manage
// only specific custom fields in Terraform while others are preserved.
func TestAccDeviceResource_CustomFieldsPartialManagement(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-partial")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfEnv := testutil.RandomCustomFieldName("tf_environment")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost_center")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with environment="prod" and owner="team-a"
				Config: testAccDeviceResourceConfig_partialManagement_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update environment to "staging", remove owner from config
				// Owner should be preserved (partial management)
				Config: testAccDeviceResourceConfig_partialManagement_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfOwner, "text", "team-a"), // Preserved!
				),
			},
			{
				// Step 3: Add cost_center, owner and environment should still be present
				Config: testAccDeviceResourceConfig_partialManagement_step3(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfOwner, "text", "team-a"), // Still preserved!
					testutil.CheckCustomFieldValue("netbox_device.test", cfCostCenter, "text", "123"),
				),
			},
		},
	})
}

// TestAccDeviceResource_CustomFieldsExplicitRemoval tests that users can
// explicitly remove specific custom fields by setting value to empty string.
func TestAccDeviceResource_CustomFieldsExplicitRemoval(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-remove")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfFieldA := testutil.RandomCustomFieldName("tf_field_a")
	cfFieldB := testutil.RandomCustomFieldName("tf_field_b")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two custom fields
				Config: testAccDeviceResourceConfig_explicitRemoval_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfFieldA, "text", "value1"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfFieldB, "text", "value2"),
				),
			},
			{
				// Step 2: Explicitly remove field_a with empty value, field_b should be preserved
				Config: testAccDeviceResourceConfig_explicitRemoval_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					// After removal, only field_b should remain
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfFieldB, "text", "value2"),
				),
			},
		},
	})
}

// TestAccDeviceResource_CustomFieldsCompleteRemoval tests that users can
// remove ALL custom fields by setting custom_fields to an empty list.
func TestAccDeviceResource_CustomFieldsCompleteRemoval(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-clear")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfFieldA := testutil.RandomCustomFieldName("tf_field_a")
	cfFieldB := testutil.RandomCustomFieldName("tf_field_b")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two custom fields
				Config: testAccDeviceResourceConfig_completeRemoval_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
				),
			},
			{
				// Step 2: Clear all custom fields with empty list
				Config: testAccDeviceResourceConfig_completeRemoval_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					// All custom fields should be removed
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
		},
	})
}

// Helper configs for partial management test
func testAccDeviceResourceConfig_partialManagement_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[12]q
  type         = "text"
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

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "prod"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "team-a"
    }
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfEnv, cfOwner, cfCostCenter,
	)
}

func testAccDeviceResourceConfig_partialManagement_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[12]q
  type         = "text"
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

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
    # NOTE: owner removed from config but should be preserved in NetBox
  ]

  depends_on = [netbox_custom_field.owner]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfEnv, cfOwner, cfCostCenter,
	)
}

func testAccDeviceResourceConfig_partialManagement_step3(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[12]q
  type         = "text"
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

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.cost_center.name
      type  = "text"
      value = "123"
    }
    # NOTE: owner still not in config but should be preserved
  ]

  depends_on = [netbox_custom_field.owner]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfEnv, cfOwner, cfCostCenter,
	)
}

// Helper configs for explicit removal test
func testAccDeviceResourceConfig_explicitRemoval_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[11]q
  type         = "text"
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

  custom_fields = [
    {
      name  = netbox_custom_field.field_a.name
      type  = "text"
      value = "value1"
    },
    {
      name  = netbox_custom_field.field_b.name
      type  = "text"
      value = "value2"
    }
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfFieldA, cfFieldB,
	)
}

func testAccDeviceResourceConfig_explicitRemoval_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[11]q
  type         = "text"
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

  custom_fields = [
    {
      name  = netbox_custom_field.field_a.name
      type  = "text"
      value = "" # Empty value removes field_a
    }
    # field_b not in config, should be preserved
  ]

  depends_on = [netbox_custom_field.field_b]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfFieldA, cfFieldB,
	)
}

// Helper configs for complete removal test
func testAccDeviceResourceConfig_completeRemoval_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[11]q
  type         = "text"
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

  custom_fields = [
    {
      name  = netbox_custom_field.field_a.name
      type  = "text"
      value = "value1"
    },
    {
      name  = netbox_custom_field.field_b.name
      type  = "text"
      value = "value2"
    }
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfFieldA, cfFieldB,
	)
}

func testAccDeviceResourceConfig_completeRemoval_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[11]q
  type         = "text"
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

  custom_fields = [] # Empty list removes all custom fields

  depends_on = [
    netbox_custom_field.field_a,
    netbox_custom_field.field_b
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfFieldA, cfFieldB,
	)
}
