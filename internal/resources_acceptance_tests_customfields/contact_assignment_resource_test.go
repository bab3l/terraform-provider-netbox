//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccContactAssignmentResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a contact assignment.
func TestAccContactAssignmentResource_CustomFieldsPreservation(t *testing.T) {
	contactName := testutil.RandomName("contact_preserve")
	groupName := testutil.RandomName("contact_group")
	groupSlug := testutil.RandomSlug("contact_group")
	roleName := testutil.RandomName("contact_role")
	roleSlug := testutil.RandomSlug("contact_role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with custom fields
				Config: testAccContactAssignmentConfig_preservation_step1(
					contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_contact_assignment.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_contact_assignment.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update priority WITHOUT mentioning custom_fields
				Config: testAccContactAssignmentConfig_preservation_step2(
					contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfText, cfInteger, "secondary",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "secondary"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Add custom_fields back to verify they were preserved
				Config: testAccContactAssignmentConfig_preservation_step1(
					contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_contact_assignment.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_contact_assignment.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccContactAssignmentConfig_preservation_step1(
	contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[8]q
  type         = "text"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "integer" {
  name         = %[9]q
  type         = "integer"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.text]
}

resource "netbox_contact_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_contact" "test" {
  name  = %[1]q
  group = netbox_contact_group.test.id
}

resource "netbox_contact_role" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_site" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[10]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[11]d"
    }
  ]
}
`, contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccContactAssignmentConfig_preservation_step2(
	contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfTextName, cfIntName, priority string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[8]q
  type         = "text"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "integer" {
  name         = %[9]q
  type         = "integer"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.text]
}

resource "netbox_contact_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_contact" "test" {
  name  = %[1]q
  group = netbox_contact_group.test.id
}

resource "netbox_contact_role" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_site" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = %[10]q
  # custom_fields intentionally omitted - should preserve existing values
}
`, contactName, groupName, groupSlug, roleName, roleSlug, siteName, siteSlug, cfTextName, cfIntName, priority)
}

func TestAccContactAssignmentResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	contactName := testutil.RandomName("contact")
	contactGroupName := testutil.RandomName("contact_group")
	contactGroupSlug := testutil.RandomSlug("contact_group")
	contactRoleName := testutil.RandomName("contact_role")
	contactRoleSlug := testutil.RandomSlug("contact_role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")

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
				Config: testAccContactAssignmentResourceImportConfig_full(contactName, contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "tags.#", "2"),
				),
			},
			{
				Config:            testAccContactAssignmentResourceImportConfig_full(contactName, contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:      "netbox_contact_assignment.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccContactAssignmentResourceImportConfig_full(contactName, contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccContactAssignmentResourceImportConfig_full(contactName, contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_contact_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_contact" "test" {
  name  = %q
  group = netbox_contact_group.test.id
}

resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.cf_text]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.cf_longtext]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.cf_integer]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.cf_boolean]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.cf_date]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["tenancy.contactassignment"]

  depends_on = [netbox_custom_field.cf_url]
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
resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id

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
		contactGroupName, contactGroupSlug,
		contactName,
		contactRoleName, contactRoleSlug,
		siteName, siteSlug,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
	)
}
