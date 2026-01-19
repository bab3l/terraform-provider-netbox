//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccPowerPanelResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a power panel. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create power panel with custom fields
// 2. Update power panel WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccPowerPanelResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-pp-site-preserve")
	siteSlug := testutil.RandomSlug("tf-test-pp-site-preserve")
	powerPanelName := testutil.RandomName("tf-test-pp-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPowerPanelCleanup(powerPanelName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckPowerPanelDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create power panel WITH custom fields explicitly in config
				Config: testAccPowerPanelConfig_preservation_step1(
					siteName, siteSlug, powerPanelName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", powerPanelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_power_panel.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_power_panel.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update power panel WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccPowerPanelConfig_preservation_step2(
					siteName, siteSlug, powerPanelName,
					cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", powerPanelName),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_power_panel.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_power_panel.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,                             // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccPowerPanelConfig_preservation_step1(
					siteName, siteSlug, powerPanelName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_power_panel.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_power_panel.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_power_panel.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccPowerPanelConfig_preservation_step1(
	siteName, siteSlug, powerPanelName,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_custom_field" "text" {
  name         = %[4]q
  type         = "text"
  object_types = ["dcim.powerpanel"]
}

resource "netbox_custom_field" "integer" {
  name         = %[5]q
  type         = "integer"
  object_types = ["dcim.powerpanel"]
}

resource "netbox_power_panel" "test" {
  name = %[3]q
  site = netbox_site.test.id

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
		siteName, siteSlug, powerPanelName,
		cfTextName, cfIntName, cfTextValue, cfIntValue,
	)
}

func testAccPowerPanelConfig_preservation_step2(
	siteName, siteSlug, powerPanelName,
	cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_custom_field" "text" {
  name         = %[4]q
  type         = "text"
  object_types = ["dcim.powerpanel"]
}

resource "netbox_custom_field" "integer" {
  name         = %[5]q
  type         = "integer"
  object_types = ["dcim.powerpanel"]
}

resource "netbox_power_panel" "test" {
  name = %[3]q
  site = netbox_site.test.id

  # NOTE: custom_fields is intentionally omitted to test preservation behavior
}
`,
		siteName, siteSlug, powerPanelName,
		cfTextName, cfIntName,
	)
}
