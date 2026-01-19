//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccTunnelResource_importWithCustomFieldsAndTags tests that custom fields and tags
// are correctly imported for a tunnel resource.
func TestAccTunnelResource_importWithCustomFieldsAndTags(t *testing.T) {
	tunnelName := testutil.RandomName("tf-test-tunnel-import")
	cfText := testutil.RandomCustomFieldName("tf_tunnel_text")
	tagName := testutil.RandomName("tf-test-tunnel-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckTunnelDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelConfig_importWithCustomFieldsAndTags(
					tunnelName, cfText, tagName, "import-value",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", tunnelName),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfText, "text", "import-value"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "1"),
				),
			},
			{
				ResourceName:      "netbox_tunnel.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config: testAccTunnelConfig_importWithCustomFieldsAndTags(
					tunnelName, cfText, tagName, "import-value",
				),
				PlanOnly: true,
			},
		},
	})
}

// TestAccTunnelResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a tunnel.
func TestAccTunnelResource_CustomFieldsPreservation(t *testing.T) {
	tunnelName := testutil.RandomName("tf-test-tunnel-preserve")
	cfText := testutil.RandomCustomFieldName("tf_tunnel_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_tunnel_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckTunnelDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create tunnel WITH custom fields explicitly in config
				Config: testAccTunnelConfig_preservation_step1(
					tunnelName, cfText, cfInteger, "initial value", 100,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", tunnelName),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfInteger, "integer", "100"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				Config: testAccTunnelConfig_preservation_step2(
					tunnelName, cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", tunnelName),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", "Updated description"),
					// Filter-to-owned: custom_fields omitted from config, so state shows 0
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_tunnel.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccTunnelConfig_preservation_step1(
					tunnelName, cfText, cfInteger, "initial value", 100,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfInteger, "integer", "100"),
				),
			},
		},
	})
}

// TestAccTunnelResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
// where state only shows custom fields declared in config, but all fields are preserved in NetBox.
func TestAccTunnelResource_CustomFieldsFilterToOwned(t *testing.T) {
	tunnelName := testutil.RandomName("tf-test-tunnel-filter")
	cfManaged := testutil.RandomCustomFieldName("tf_tunnel_managed")
	cfUnmanaged := testutil.RandomCustomFieldName("tf_tunnel_unmanaged")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckTunnelDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with both fields in config
				Config: testAccTunnelConfig_filterToOwned_step1(
					tunnelName, cfManaged, cfUnmanaged, "managed-value", "unmanaged-value",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfManaged, "text", "managed-value"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfUnmanaged, "text", "unmanaged-value"),
				),
			},
			{
				// Step 2: Remove unmanaged field from config (but preserve in NetBox)
				Config: testAccTunnelConfig_filterToOwned_step2(
					tunnelName, cfManaged, cfUnmanaged, "managed-value",
				),
				Check: resource.ComposeTestCheckFunc(
					// State only shows the managed field
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfManaged, "text", "managed-value"),
				),
			},
			{
				// Step 3: Import to verify both fields exist in NetBox
				ResourceName:            "netbox_tunnel.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Update managed field, unmanaged field still preserved
				Config: testAccTunnelConfig_filterToOwned_step2(
					tunnelName, cfManaged, cfUnmanaged, "new-managed-value",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfManaged, "text", "new-managed-value"),
				),
			},
			{
				// Step 5: Add unmanaged field back to config to verify it was preserved
				Config: testAccTunnelConfig_filterToOwned_step1(
					tunnelName, cfManaged, cfUnmanaged, "new-managed-value", "unmanaged-value",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfManaged, "text", "new-managed-value"),
					testutil.CheckCustomFieldValue("netbox_tunnel.test", cfUnmanaged, "text", "unmanaged-value"),
				),
			},
		},
	})
}

// Config helpers

func testAccTunnelConfig_importWithCustomFieldsAndTags(
	tunnelName, cfName, tagName, cfValue string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[2]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_tag" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[4]q
    }
  ]

	tags = [netbox_tag.test.slug]
}
`, tunnelName, cfName, tagName, cfValue)
}

func testAccTunnelConfig_preservation_step1(
	tunnelName, cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  object_types = ["vpn.tunnel"]
  type         = "integer"
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"
  description   = "Initial description"

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
`, tunnelName, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccTunnelConfig_preservation_step2(
	tunnelName, cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  object_types = ["vpn.tunnel"]
  type         = "integer"
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"
  description   = "Updated description"

  # custom_fields intentionally omitted - should preserve in NetBox
}
`, tunnelName, cfTextName, cfIntName)
}

func testAccTunnelConfig_filterToOwned_step1(
	tunnelName, cfManaged, cfUnmanaged, managedValue, unmanagedValue string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "managed" {
  name         = %[2]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_custom_field" "unmanaged" {
  name         = %[3]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"

  custom_fields = [
    {
      name  = netbox_custom_field.managed.name
      type  = "text"
      value = %[4]q
    },
    {
      name  = netbox_custom_field.unmanaged.name
      type  = "text"
      value = %[5]q
    }
  ]
}
`, tunnelName, cfManaged, cfUnmanaged, managedValue, unmanagedValue)
}

func testAccTunnelConfig_filterToOwned_step2(
	tunnelName, cfManaged, cfUnmanaged, managedValue string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "managed" {
  name         = %[2]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_custom_field" "unmanaged" {
  name         = %[3]q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"

  # Only manage the "managed" custom field
  # The "unmanaged" field should be preserved in NetBox but not in state
  custom_fields = [
    {
      name  = netbox_custom_field.managed.name
      type  = "text"
      value = %[4]q
    }
  ]
}
`, tunnelName, cfManaged, cfUnmanaged, managedValue)
}
