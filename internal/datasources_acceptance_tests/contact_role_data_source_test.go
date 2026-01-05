package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactRoleDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-ds-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactRoleDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccContactRoleDataSource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactRoleDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "description", "Test Contact Role Description"),
				),
			},
		},
	})
}

func testAccContactRoleDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_role" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Role Description"
}

data "netbox_contact_role" "test" {
  id = netbox_contact_role.test.id
}
`, name, slug)
}
func TestAccContactRoleDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactRoleDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccContactRoleDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_role" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Role Description"
}

data "netbox_contact_role" "test" {
  name = netbox_contact_role.test.name
}
`, name, slug)
}

func TestAccContactRoleDataSource_bySlug(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactRoleDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleDataSourceConfigBySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccContactRoleDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_role" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Role Description"
}

data "netbox_contact_role" "test" {
  slug = netbox_contact_role.test.slug
}
`, name, slug)
}
