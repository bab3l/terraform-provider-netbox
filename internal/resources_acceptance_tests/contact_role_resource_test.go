package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for contact role resource are in resources_acceptance_tests_customfields package

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
