package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(1, 50)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := acctest.RandIntRange(10, 100)
	endOctet := startOctet + 5

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfig(secondOctet, thirdOctet, startOctet, endOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "start_address", fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, startOctet)),
				),
			},
		},
	})
}

func TestAccIPRangeDataSource_basic(t *testing.T) {

	t.Parallel()

	secondOctet := acctest.RandIntRange(1, 50)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := acctest.RandIntRange(10, 100)
	endOctet := startOctet + 5

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfig(secondOctet, thirdOctet, startOctet, endOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "start_address", fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, startOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "end_address", fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, endOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "status", "active"),
				),
			},
		},
	})
}

func testAccIPRangeDataSourceConfig(secondOctet, thirdOctet, startOctet, endOctet int) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = "10.%d.%d.%d"
  end_address   = "10.%d.%d.%d"
  status        = "active"
}

data "netbox_ip_range" "test" {
  id = netbox_ip_range.test.id
}
`, secondOctet, thirdOctet, startOctet, secondOctet, thirdOctet, endOctet)
}

func TestAccIPRangeDataSource_byAddresses(t *testing.T) {

	t.Parallel()

	secondOctet := acctest.RandIntRange(51, 100)
	thirdOctet := acctest.RandIntRange(51, 100)
	startOctet := acctest.RandIntRange(10, 100)
	endOctet := startOctet + 5

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfigByAddresses(secondOctet, thirdOctet, startOctet, endOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "start_address", fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, startOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "end_address", fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, endOctet)),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "status", "active"),
				),
			},
		},
	})
}

func testAccIPRangeDataSourceConfigByAddresses(secondOctet, thirdOctet, startOctet, endOctet int) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = "10.%d.%d.%d"
  end_address   = "10.%d.%d.%d"
  status        = "active"
}

data "netbox_ip_range" "test" {
  start_address = netbox_ip_range.test.start_address
  end_address   = netbox_ip_range.test.end_address
}
`, secondOctet, thirdOctet, startOctet, secondOctet, thirdOctet, endOctet)
}
