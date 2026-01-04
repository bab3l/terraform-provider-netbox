package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_contact_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactRoleResourceConfig(name+"-updated", slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name+"-updated"),
				),
			},
		},
	})
}

func TestAccContactRoleResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-full")
	slug := testutil.GenerateSlug(name)
	description := testutil.RandomName("description")
	updatedDescription := "Updated contact role description"
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")
	customFieldName := testutil.RandomCustomFieldName("tf_test_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig_full(name, slug, description, tagName, tagSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccContactRoleResourceConfig_full(name, slug, updatedDescription, tagName, tagSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccContactRoleResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-update")
	slug := testutil.GenerateSlug(name)
	updatedName := testutil.RandomName("test-contact-role-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
				),
			},
			{
				Config: testAccContactRoleResourceConfig(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccContactRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cr-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccConsistency_ContactRole_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-lit")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
			{
				Config:   testAccContactRoleConsistencyLiteralNamesConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
				),
			},
		},
	})
}

func testAccContactRoleConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccContactRoleResourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccContactRoleResourceConfig_full(name, slug, description, tagName, tagSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["tenancy.contactrole"]
  type         = "text"
}

resource "netbox_contact_role" "test" {
  name        = %q
  slug        = %q
  description = %q
  tags = [
    {
      name = netbox_tag.test.name
      slug = netbox_tag.test.slug
    }
  ]
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-value"
    }
  ]
}
`, tagName, tagSlug, customFieldName, name, slug, description)
}

func TestAccContactRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyContactRolesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find contact_role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyContactRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete contact_role: %v", err)
					}
					t.Logf("Successfully externally deleted contact_role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccContactRoleResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	roleName := testutil.RandomName("contact_role")
	roleSlug := testutil.RandomSlug("contact_role")

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
				Config: testAccContactRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", roleSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_contact_role.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_contact_role.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccContactRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_contact_role.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
			{
				Config:   testAccContactRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccContactRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["tenancy.contactrole"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["tenancy.contactrole"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["tenancy.contactrole"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["tenancy.contactrole"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["tenancy.contactrole"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["tenancy.contactrole"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["tenancy.contactrole"]
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
resource "netbox_contact_role" "test" {
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
		roleName, roleSlug,
	)
}
