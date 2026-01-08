//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccEventRuleResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on an event rule. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create event rule with custom fields
// 2. Update event rule WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccEventRuleResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	ruleName := testutil.RandomName("tf-test-er-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create event rule WITH custom fields explicitly in config
				Config: testAccEventRuleConfig_preservation_step1(
					ruleName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", ruleName),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_event_rule.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_event_rule.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update enabled status WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccEventRuleConfig_preservation_step2(
					ruleName,
					cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", ruleName),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_event_rule.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_event_rule.test",
				ImportState:             true,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccEventRuleConfig_preservation_step1(
					ruleName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_event_rule.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_event_rule.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_event_rule.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccEventRuleConfig_preservation_step1(
	ruleName,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  object_types = ["extras.eventrule"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  object_types = ["extras.eventrule"]
  type         = "integer"
}

resource "netbox_event_rule" "test" {
  name           = %[1]q
  object_types   = ["dcim.device"]
  event_types    = ["object_created"]
  enabled        = true
  action_type    = "webhook"
  action_object_type = "extras.webhook"

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
`, ruleName, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccEventRuleConfig_preservation_step2(
	ruleName,
	cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  object_types = ["extras.eventrule"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  object_types = ["extras.eventrule"]
  type         = "integer"
}

resource "netbox_event_rule" "test" {
  name           = %[1]q
  object_types   = ["dcim.device"]
  event_types    = ["object_created"]
  enabled        = false  # Changed this field
  action_type    = "webhook"
  action_object_type = "extras.webhook"

  # custom_fields intentionally omitted
}
`, ruleName, cfTextName, cfIntName)
}
