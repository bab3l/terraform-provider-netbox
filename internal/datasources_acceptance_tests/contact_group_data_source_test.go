package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactGroupDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-ds-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_contact_group.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "slug", slug),
				),
			},
		},
	})
}

func TestAccContactGroupDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupDataSourceConfigByID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "description", "Test Contact Group Description"),
				),
			},
		},
	})
}

func TestAccContactGroupDataSource_byName(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("Public Cloud %s", testutil.RandomName("test-contact-group"))
	slug := testutil.RandomSlug("test-contact-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_name", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_name", "description", "Test Contact Group Description"),
				),
			},
		},
	})
}

func TestAccContactGroupDataSource_bySlug(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupDataSourceConfigBySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_slug", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_slug", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_slug", "description", "Test Contact Group Description"),
				),
			},
		},
	})
}

func testAccContactGroupDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Group Description"
}

data "netbox_contact_group" "by_id" {
  id = netbox_contact_group.test.id
}
`, name, slug)
}

func testAccContactGroupDataSourceConfigByID(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Group Description"
}

data "netbox_contact_group" "by_id" {
  id = netbox_contact_group.test.id
}
`, name, slug)
}

func testAccContactGroupDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Group Description"
}

data "netbox_contact_group" "by_name" {
  name = netbox_contact_group.test.name
}
`, name, slug)
}

func testAccContactGroupDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Group Description"
}

data "netbox_contact_group" "by_slug" {
  slug = netbox_contact_group.test.slug
}
`, name, slug)
}
