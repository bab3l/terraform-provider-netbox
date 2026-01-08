//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccContactResource_TagsPreservation tests that tags are preserved
// when updating other fields on a contact. This addresses a critical bug where tags
// were being deleted when users updated unrelated fields.
func TestAccContactResource_TagsPreservation(t *testing.T) {
	contactName := testutil.RandomName("contact_preserve")
	groupName := testutil.RandomName("contact_group")
	groupSlug := testutil.RandomSlug("contact_group")
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with tags
				Config: testAccContactConfig_preservation_step1(
					contactName, groupName, groupSlug, tag1, tag1Slug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("netbox_contact.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_contact.test", "tags.#", "1"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning tags
				Config: testAccContactConfig_preservation_step2(
					contactName, groupName, groupSlug, tag1, tag1Slug, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("netbox_contact.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_contact.test", "tags.#", "0"),
				),
			},
			{
				// Step 3: Add tags back to verify they were preserved
				Config: testAccContactConfig_preservation_step1(
					contactName, groupName, groupSlug, tag1, tag1Slug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "tags.#", "1"),
				),
			},
		},
	})
}

func testAccContactConfig_preservation_step1(
	contactName, groupName, groupSlug, tag1, tag1Slug string,
) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tag" "tag1" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_contact" "test" {
  name        = %[1]q
  group       = netbox_contact_group.test.id
  description = "Initial description"

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    }
  ]
}
`, contactName, groupName, groupSlug, tag1, tag1Slug)
}

func testAccContactConfig_preservation_step2(
	contactName, groupName, groupSlug, tag1, tag1Slug, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tag" "tag1" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_contact" "test" {
  name        = %[1]q
  group       = netbox_contact_group.test.id
  description = %[6]q
  # tags intentionally omitted - should preserve existing values
}
`, contactName, groupName, groupSlug, tag1, tag1Slug, description)
}
