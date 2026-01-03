package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactGroupResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("test-contact-group")

	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactGroupResourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_contact_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},

			{

				Config: testAccContactGroupResourceConfig(name+"-updated", slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name+"-updated"),
				),
			},
		},
	})

}

func TestAccContactGroupResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-update")
	slug := testutil.GenerateSlug(name)
	updatedName := testutil.RandomName("test-contact-group-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
				),
			},
			{
				Config: testAccContactGroupResourceConfig(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccContactGroupResource_IDPreservation(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cg-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckContactGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccContactGroupResourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccConsistency_ContactGroup_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("test-contact-group-lit")

	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactGroupConsistencyLiteralNamesConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},

			{

				Config: testAccContactGroupConsistencyLiteralNamesConfig(name, slug),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
				),
			},
		},
	})

}

func testAccContactGroupConsistencyLiteralNamesConfig(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_group" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccContactGroupResourceConfig(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_group" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func TestAccContactGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyContactGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find contact_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyContactGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete contact_group: %v", err)
					}
					t.Logf("Successfully externally deleted contact_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccContactGroupResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

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
				Config:                  testAccContactGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_contact_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
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
