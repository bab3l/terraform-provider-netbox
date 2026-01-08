//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRoleResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a role.
func TestAccRoleResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	roleName := testutil.RandomName("tf-test-role-preserve")
	roleSlug := testutil.RandomSlug("tf-test-role-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRoleDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create role WITH custom fields explicitly in config
				Config: testAccRoleConfig_preservation_step1(
					roleName, roleSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", roleSlug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_role.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_role.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				Config: testAccRoleConfig_preservation_step2(
					roleName, roleSlug,
					cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", roleSlug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_role.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccRoleConfig_preservation_step1(
					roleName, roleSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_role.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_role.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccRoleConfig_preservation_step1(
	roleName, roleSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  content_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  content_types = ["ipam.role"]
  type         = "integer"
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Initial description"
  weight      = 1000

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[5]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = %[6]d
    }
  ]
}
`, roleName, roleSlug, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccRoleConfig_preservation_step2(
	roleName, roleSlug,
	cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  content_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  content_types = ["ipam.role"]
  type         = "integer"
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Updated description"
  weight      = 1000

  # custom_fields intentionally omitted
}
`, roleName, roleSlug, cfTextName, cfIntName)
}
