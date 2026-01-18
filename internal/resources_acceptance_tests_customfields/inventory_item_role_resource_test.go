//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	roleName := testutil.RandomName("inventory_item_role")
	roleSlug := testutil.RandomSlug("inventory_item_role")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", roleSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_inventory_item_role.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
			{
				Config:   testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.inventoryitemrole"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Main Resource
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
    }
  ]

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		roleName, roleSlug,
	)
}

// TestAccInventoryItemRoleResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on an inventory item role.
func TestAccInventoryItemRoleResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	roleName := testutil.RandomName("tf-test-iir-preserve")
	roleSlug := testutil.RandomSlug("tf-test-iir-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckInventoryItemRoleDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create inventory item role WITH custom fields explicitly in config
				Config: testAccInventoryItemRoleConfig_preservation_step1(
					roleName, roleSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", roleSlug),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_inventory_item_role.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_inventory_item_role.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				Config: testAccInventoryItemRoleConfig_preservation_step2(
					roleName, roleSlug,
					cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", roleSlug),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_inventory_item_role.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccInventoryItemRoleConfig_preservation_step1(
					roleName, roleSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_inventory_item_role.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_inventory_item_role.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccInventoryItemRoleConfig_preservation_step1(
	roleName, roleSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  object_types = ["dcim.inventoryitemrole"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  object_types = ["dcim.inventoryitemrole"]
  type         = "integer"
}

resource "netbox_inventory_item_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Initial description"

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

func testAccInventoryItemRoleConfig_preservation_step2(
	roleName, roleSlug,
	cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  object_types = ["dcim.inventoryitemrole"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  object_types = ["dcim.inventoryitemrole"]
  type         = "integer"
}

resource "netbox_inventory_item_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Updated description"

  # custom_fields intentionally omitted
}
`, roleName, roleSlug, cfTextName, cfIntName)
}
