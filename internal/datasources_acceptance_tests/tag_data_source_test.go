package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	tagName := testutil.RandomName("tag")
	tagSlug := testutil.RandomSlug("tag")

	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "slug", tagSlug),
				),
			},
		},
	})
}

func TestAccTagDataSource_byName(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	tagName := fmt.Sprintf("Public Cloud %s", testutil.RandomName("tag"))
	tagSlug := testutil.RandomSlug("tag")

	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfigByName(tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "slug", tagSlug),
				),
			},
		},
	})
}

func TestAccTagDataSource_bySlug(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	tagName := testutil.RandomName("tag")
	tagSlug := testutil.RandomSlug("tag")

	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfigBySlug(tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "slug", tagSlug),
				),
			},
		},
	})
}

func TestAccTagDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	cleanup := testutil.NewCleanupResource(t)
	tagName := testutil.RandomName("tag-ds-id")
	tagSlug := testutil.RandomSlug("tag-ds-id")
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "slug", tagSlug),
				),
			},
		},
	})
}

func testAccTagDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_tag" "test" {
  id = netbox_tag.test.id
}
`, name, slug)
}

func testAccTagDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_tag" "test" {
  name = netbox_tag.test.name
}
`, name, slug)
}

func testAccTagDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_tag" "test" {
  slug = netbox_tag.test.slug
}
`, name, slug)
}
