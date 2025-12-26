package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressDataSource_basic(t *testing.T) {

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
				Config: testAccIPAddressDataSourceConfig(ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "address", ipAddress),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccIPAddressDataSource_IDPreservation(t *testing.T) {
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
				Config: testAccIPAddressDataSourceConfig(ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "address", ipAddress),
				),
			},
		},
	})
}

func testAccIPAddressDataSourceConfig(ipAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = "%s"
  status  = "active"
}

data "netbox_ip_address" "test" {
  address = netbox_ip_address.test.address
}
`, ipAddress)
}

func TestAccIPAddressDataSource_byID(t *testing.T) {

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
				Config: testAccIPAddressDataSourceConfigByID(ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "address", ipAddress),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "status", "active"),
				),
			},
		},
	})
}

func testAccIPAddressDataSourceConfigByID(ipAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = "%s"
  status  = "active"
}

data "netbox_ip_address" "test" {
  id = netbox_ip_address.test.id
}
`, ipAddress)
}
