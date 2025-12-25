package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeDataSource_basic(t *testing.T) {

	t.Parallel()

	startOctet := acctest.RandIntRange(10, 100)
	endOctet := startOctet + 5

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfig(startOctet, endOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "start_address", fmt.Sprintf("10.0.1.%d/32", startOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "end_address", fmt.Sprintf("10.0.1.%d/32", endOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "status", "active"),
				),
			},
		},
	})
}

func testAccIPRangeDataSourceConfig(startOctet, endOctet int) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = "10.0.1.%d"
  end_address   = "10.0.1.%d"
  status        = "active"
}

data "netbox_ip_range" "test" {
  id = netbox_ip_range.test.id
}
`, startOctet, endOctet)
}

func TestAccIPRangeDataSource_byAddresses(t *testing.T) {

	t.Parallel()

	startOctet := acctest.RandIntRange(110, 200)
	endOctet := startOctet + 5

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfigByAddresses(startOctet, endOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "start_address", fmt.Sprintf("10.0.0.%d/24", startOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "end_address", fmt.Sprintf("10.0.0.%d/24", endOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "status", "active"),
				),
			},
		},
	})
}

func testAccIPRangeDataSourceConfigByAddresses(startOctet, endOctet int) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = "10.0.0.%d/24"
  end_address   = "10.0.0.%d/24"
  status        = "active"
}

data "netbox_ip_range" "test" {
  start_address = netbox_ip_range.test.start_address
  end_address   = netbox_ip_range.test.end_address
}
`, startOctet, endOctet)
}
