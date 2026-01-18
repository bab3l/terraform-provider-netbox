package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixesDataSource_byPrefixFilter(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixesDataSourceConfig_byPrefix(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.0", prefix),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "ids.0", "netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "prefixes.0.id", "netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "prefixes.0.prefix", prefix),
				),
			},
		},
	})
}

func TestAccPrefixesDataSource_byTagFilter(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	tagName := testutil.RandomName("tf-test-tag-prefix-q")
	tagSlug := testutil.RandomSlug("tf-test-tag-prefix-q")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPrefixDestroy,
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixesDataSourceConfig_byTag(prefix, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.0", prefix),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "ids.0", "netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "prefixes.0.id", "netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "prefixes.0.prefix", prefix),
				),
			},
		},
	})
}

func TestAccPrefixesDataSource_byPrefixAndStatusFilters(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixesDataSourceConfig_byPrefixAndStatus(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.0", prefix),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "ids.0", "netbox_prefix.test", "id"),
				),
			},
		},
	})
}

func testAccPrefixesDataSourceConfig_byPrefix(prefix string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %q
  status = "active"
}

data "netbox_prefixes" "test" {
  filter {
    name   = "prefix"
    values = [netbox_prefix.test.prefix]
  }
}
`, prefix)
}

func testAccPrefixesDataSourceConfig_byTag(prefix, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_prefix" "test" {
  prefix = %q
  status = "active"

	tags = [netbox_tag.test.slug]
}

data "netbox_prefixes" "test" {
  filter {
    name   = "tag"
    values = [netbox_tag.test.slug]
  }

  depends_on = [netbox_prefix.test]
}
`, tagName, tagSlug, prefix)
}

func testAccPrefixesDataSourceConfig_byPrefixAndStatus(prefix string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %q
  status = "active"
}

data "netbox_prefixes" "test" {
  filter {
    name   = "status"
    values = ["active"]
  }

  filter {
    name   = "prefix"
    values = [netbox_prefix.test.prefix]
  }
}
`, prefix)
}
