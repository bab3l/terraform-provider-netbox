//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualChassisResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	virtualChassisName := testutil.RandomName("virtual_chassis")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualChassisCleanup(virtualChassisName)
	cleanup.RegisterTenantCleanup(tenantSlug)
	// Clean up custom fields and tags
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualChassisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceImportConfig_full(virtualChassisName, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", virtualChassisName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_virtual_chassis.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Custom fields have import limitations
			},
		},
	})
}

func testAccVirtualChassisResourceImportConfig_full(virtualChassisName, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["dcim.virtualchassis"]
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

# Virtual Chassis with comprehensive custom fields and tags
resource "netbox_virtual_chassis" "test" {
  name   = %q
  domain = "test-domain"

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this virtual chassis resource for testing purposes."
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
      value = "2023-01-15"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key": "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, virtualChassisName)
}

// TestAccVirtualChassisResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a virtual chassis. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create virtual chassis with custom fields
// 2. Update virtual chassis WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccVirtualChassisResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	virtualChassisName := testutil.RandomName("tf-test-vc-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualChassisCleanup(virtualChassisName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckVirtualChassisDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create virtual chassis WITH custom fields explicitly in config
				Config: testAccVirtualChassisConfig_preservation_step1(
					virtualChassisName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", virtualChassisName),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "initial-domain"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_virtual_chassis.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_virtual_chassis.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update domain WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccVirtualChassisConfig_preservation_step2(
					virtualChassisName,
					cfText, cfInteger, "updated-domain",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", virtualChassisName),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "updated-domain"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_virtual_chassis.test",
				ImportState:             true,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccVirtualChassisConfig_preservation_step1(
					virtualChassisName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_virtual_chassis.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_virtual_chassis.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccVirtualChassisConfig_preservation_step1(
	virtualChassisName,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  type         = "integer"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_virtual_chassis" "test" {
  name   = %[1]q
  domain = "initial-domain"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[4]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[5]d"
    }
  ]

  depends_on = [
    netbox_custom_field.text,
    netbox_custom_field.integer,
  ]
}
`,
		virtualChassisName,
		cfTextName, cfIntName, cfTextValue, cfIntValue,
	)
}

func testAccVirtualChassisConfig_preservation_step2(
	virtualChassisName,
	cfTextName, cfIntName, domain string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_custom_field" "integer" {
  name         = %[3]q
  type         = "integer"
  object_types = ["dcim.virtualchassis"]
}

resource "netbox_virtual_chassis" "test" {
  name   = %[1]q
  domain = %[4]q

  # NOTE: custom_fields is intentionally omitted to test preservation behavior
}
`,
		virtualChassisName,
		cfTextName, cfIntName, domain,
	)
}
