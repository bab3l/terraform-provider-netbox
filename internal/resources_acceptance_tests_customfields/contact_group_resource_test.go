//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccContactGroupResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a contact group. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
func TestAccContactGroupResource_CustomFieldsPreservation(t *testing.T) {
	groupName := testutil.RandomName("contact_group_preserve")
	groupSlug := testutil.RandomSlug("contact_group_preserve")
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with custom fields
				Config: testAccContactGroupConfig_preservation_step1(
					groupName, groupSlug, cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", groupName),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_contact_group.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_contact_group.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccContactGroupConfig_preservation_step2(
					groupName, groupSlug, cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", groupName),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Add custom_fields back to verify they were preserved
				Config: testAccContactGroupConfig_preservation_step1(
					groupName, groupSlug, cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_contact_group.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_contact_group.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccContactGroupConfig_preservation_step1(
	groupName, groupSlug, cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  type         = "text"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  type         = "integer"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_contact_group" "test" {
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
      value = "%[6]d"
    }
  ]
}
`, groupName, groupSlug, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccContactGroupConfig_preservation_step2(
	groupName, groupSlug, cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  type         = "text"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  type         = "integer"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_contact_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[5]q
  # custom_fields intentionally omitted - should preserve existing values
}
`, groupName, groupSlug, cfTextName, cfIntName, description)
}

func TestAccContactGroupResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	groupName := testutil.RandomName("contact_group")
	groupSlug := testutil.RandomSlug("contact_group")

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
				Config: testAccContactGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", groupName),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", groupSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_contact_group.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_contact_group.test", "tags.#", "2"),
				),
			},
			{
				Config:            testAccContactGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:      "netbox_contact_group.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccContactGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccContactGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["tenancy.contactgroup"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["tenancy.contactgroup"]
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
resource "netbox_contact_group" "test" {
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
		groupName, groupSlug,
	)
}
