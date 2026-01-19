//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFHRPGroupResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on an FHRP group.
func TestAccFHRPGroupResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	groupName := testutil.RandomName("tf-test-fhrpgroup-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup("vrrp2", 1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckFHRPGroupDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create FHRP group WITH custom fields explicitly in config
				Config: testAccFHRPGroupConfig_preservation_step1(
					groupName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", "1"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_fhrp_group.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_fhrp_group.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				Config: testAccFHRPGroupConfig_preservation_step2(
					groupName,
					cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", "1"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:      "netbox_fhrp_group.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportCommandWithID,
				ImportStateVerify: false,
			},
			{
				// Step 3a: Verify no changes after import
				Config:   testAccFHRPGroupConfig_preservation_step2(groupName, cfText, cfInteger),
				PlanOnly: true,
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccFHRPGroupConfig_preservation_step1(
					groupName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_fhrp_group.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_fhrp_group.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccFHRPGroupConfig_preservation_step1(
	groupName,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  object_types = ["ipam.fhrpgroup"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  object_types = ["ipam.fhrpgroup"]
  type         = "integer"
}

resource "netbox_fhrp_group" "test" {
  group_id       = 1
  protocol       = "vrrp2"
  description    = "Initial description"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[4]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = %[5]d
    }
  ]
}
`, groupName, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccFHRPGroupConfig_preservation_step2(
	groupName,
	cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  object_types = ["ipam.fhrpgroup"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  object_types = ["ipam.fhrpgroup"]
  type         = "integer"
}

resource "netbox_fhrp_group" "test" {
  group_id       = 1
  protocol       = "vrrp2"
  description    = "Updated description"

  # custom_fields intentionally omitted
}
`, groupName, cfTextName, cfIntName)
}
