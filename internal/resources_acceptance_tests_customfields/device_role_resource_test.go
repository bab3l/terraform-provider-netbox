//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceRoleResource_importWithCustomFieldsAndTags tests importing a device role
// with custom fields and tags to ensure all data is preserved during import.
func TestAccDeviceRoleResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	roleName := testutil.RandomName("tf-test-role-import")
	roleSlug := testutil.RandomSlug("tf-test-role-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-role-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-role-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-role-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-role-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values for different data types
	cfText := testutil.RandomCustomFieldName("tf_role_text")
	cfTextValue := testutil.RandomName("role-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_role_longtext")
	cfLongtextValue := fmt.Sprintf("Role description with details: %s", testutil.RandomName("role-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_role_integer")
	cfIntegerValue := 42
	cfBoolean := testutil.RandomCustomFieldName("tf_role_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_role_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_role_url")
	cfURLValue := testutil.RandomURL("role")
	cfJSON := testutil.RandomCustomFieldName("tf_role_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the device role with all custom fields and tags
				Config: testAccDeviceRoleResourceImportConfig_full(
					roleName, roleSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_device_role.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the device role and verify all fields are preserved
				ResourceName:      "netbox_device_role.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
				// The import should preserve basic fields
				Check: resource.ComposeTestCheckFunc(
					// Verify basic fields
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", roleName),
				),
			},
			{
				Config: testAccDeviceRoleResourceImportConfig_full(
					roleName, roleSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccDeviceRoleResourceImportConfig_full(
	roleName, roleSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "role_test1" {
  name  = %[3]q
  slug  = %[4]q
  color = %[5]q
}

resource "netbox_tag" "role_test2" {
  name  = %[6]q
  slug  = %[7]q
  color = %[8]q
}

# Create custom fields for dcim.devicerole
resource "netbox_custom_field" "role_text" {
  name         = %[9]q
  type         = "text"
  object_types = ["dcim.devicerole"]
  required     = false
}

resource "netbox_custom_field" "role_longtext" {
  name         = %[11]q
  type         = "longtext"
  object_types = ["dcim.devicerole"]
  required     = false
}

resource "netbox_custom_field" "role_integer" {
  name         = %[13]q
  type         = "integer"
  object_types = ["dcim.devicerole"]
  required     = false
}

resource "netbox_custom_field" "role_boolean" {
  name         = %[15]q
  type         = "boolean"
  object_types = ["dcim.devicerole"]
  required     = false
}

resource "netbox_custom_field" "role_date" {
  name         = %[17]q
  type         = "date"
  object_types = ["dcim.devicerole"]
  required     = false
}

resource "netbox_custom_field" "role_url" {
  name         = %[19]q
  type         = "url"
  object_types = ["dcim.devicerole"]
  required     = false
}

resource "netbox_custom_field" "role_json" {
  name         = %[21]q
  type         = "json"
  object_types = ["dcim.devicerole"]
  required     = false
}

# Create device role with all custom fields and tags
resource "netbox_device_role" "test" {
  name = %[1]q
  slug = %[2]q

	depends_on = [
		netbox_custom_field.role_text,
		netbox_custom_field.role_longtext,
		netbox_custom_field.role_integer,
		netbox_custom_field.role_boolean,
		netbox_custom_field.role_date,
		netbox_custom_field.role_url,
		netbox_custom_field.role_json,
	]

  tags = [netbox_tag.role_test1.slug, netbox_tag.role_test2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.role_text.name
      type  = "text"
      value = %[10]q
    },
    {
      name  = netbox_custom_field.role_longtext.name
      type  = "longtext"
      value = %[12]q
    },
    {
      name  = netbox_custom_field.role_integer.name
      type  = "integer"
      value = "%[14]d"
    },
    {
      name  = netbox_custom_field.role_boolean.name
      type  = "boolean"
      value = "%[16]t"
    },
    {
      name  = netbox_custom_field.role_date.name
      type  = "date"
      value = %[18]q
    },
    {
      name  = netbox_custom_field.role_url.name
      type  = "url"
      value = %[20]q
    },
    {
      name  = netbox_custom_field.role_json.name
      type  = "json"
      value = %[22]q
    }
  ]
}
`, roleName, roleSlug, tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
