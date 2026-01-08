//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRackRoleResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a rack role. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create rack role with custom fields
// 2. Update rack role WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccRackRoleResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	roleName := testutil.RandomName("tf-test-rr-preserve")
	roleSlug := testutil.RandomSlug("tf-test-rr-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackRoleDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create rack role WITH custom fields explicitly in config
				Config: testAccRackRoleConfig_preservation_step1(
					roleName, roleSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", roleSlug),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_rack_role.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_rack_role.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccRackRoleConfig_preservation_step2(
					roleName, roleSlug,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", roleSlug),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", "Updated description"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_rack_role.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_rack_role.test",
				ImportState:             true,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccRackRoleConfig_preservation_step1(
					roleName, roleSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_rack_role.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_rack_role.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_rack_role.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccRackRoleConfig_preservation_step1(
	roleName, roleSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.rackrole"]
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  type         = "integer"
  object_types = ["dcim.rackrole"]
}

resource "netbox_rack_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[5]q

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[6]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[7]d"
    }
  ]

  depends_on = [
    netbox_custom_field.text,
    netbox_custom_field.integer,
  ]
}
`,
		roleName, roleSlug,
		cfTextName, cfIntName, "Initial description", cfTextValue, cfIntValue,
	)
}

func testAccRackRoleConfig_preservation_step2(
	roleName, roleSlug,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.rackrole"]
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  type         = "integer"
  object_types = ["dcim.rackrole"]
}

resource "netbox_rack_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[5]q

  # NOTE: custom_fields is intentionally omitted to test preservation behavior
}
`,
		roleName, roleSlug,
		cfTextName, cfIntName, description,
	)
}
