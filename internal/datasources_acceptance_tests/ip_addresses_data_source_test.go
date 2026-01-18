package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressesDataSource_byAddressFilter(t *testing.T) {
	t.Parallel()

	ipAddress := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ipAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressesDataSourceConfig_byAddress(ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.0", ipAddress),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ids.0", "netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ip_addresses.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ip_addresses.0.id", "netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ip_addresses.0.address", ipAddress),
				),
			},
		},
	})
}

func TestAccIPAddressesDataSource_byTagFilter(t *testing.T) {
	t.Parallel()

	ipAddress := testutil.RandomIPv4Prefix()
	tagName := testutil.RandomName("tf-test-tag-ip-q")
	tagSlug := testutil.RandomSlug("tf-test-tag-ip-q")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ipAddress)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPAddressDestroy,
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressesDataSourceConfig_byTag(ipAddress, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.0", ipAddress),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ids.0", "netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ip_addresses.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ip_addresses.0.id", "netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ip_addresses.0.address", ipAddress),
				),
			},
		},
	})
}

func TestAccIPAddressesDataSource_byAddressAndStatusFilters(t *testing.T) {
	t.Parallel()

	ipAddress := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ipAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressesDataSourceConfig_byAddressAndStatus(ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.0", ipAddress),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ids.0", "netbox_ip_address.test", "id"),
				),
			},
		},
	})
}

func testAccIPAddressesDataSourceConfig_byAddress(ipAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %q
  status  = "active"
}

data "netbox_ip_addresses" "test" {
  filter {
    name   = "address"
    values = [netbox_ip_address.test.address]
  }
}
`, ipAddress)
}

func testAccIPAddressesDataSourceConfig_byTag(ipAddress, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_ip_address" "test" {
  address = %q
  status  = "active"

	tags = [netbox_tag.test.slug]
}

data "netbox_ip_addresses" "test" {
  filter {
    name   = "tag"
    values = [netbox_tag.test.slug]
  }

  depends_on = [netbox_ip_address.test]
}
`, tagName, tagSlug, ipAddress)
}

func testAccIPAddressesDataSourceConfig_byAddressAndStatus(ipAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %q
  status  = "active"
}

data "netbox_ip_addresses" "test" {
  filter {
    name   = "status"
    values = ["active"]
  }

  filter {
    name   = "address"
    values = [netbox_ip_address.test.address]
  }
}
`, ipAddress)
}
