package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderDataSource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("provider")
	slug := testutil.RandomSlug("provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccProviderDataSource_bySlug(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("provider")
	slug := testutil.RandomSlug("provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderDataSourceConfigBySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccProviderDataSource_byName(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("provider")
	slug := testutil.RandomSlug("provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "slug", slug),
				),
			},
		},
	})
}

func testAccProviderDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_provider" "test" {
  id = netbox_provider.test.id
}
`, name, slug)
}

func testAccProviderDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_provider" "test" {
  slug = netbox_provider.test.slug
}
`, name, slug)
}

func testAccProviderDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_provider" "test" {
  name = netbox_provider.test.name
}
`, name, slug)
}
