package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeDataSource_basic(t *testing.T) {

	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "start_address", "10.0.0.10/24"),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "end_address", "10.0.0.20/24"),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "status", "active"),
				),
			},
		},
	})
}

const testAccIPRangeDataSourceConfig = `
resource "netbox_ip_range" "test" {
  start_address = "10.0.0.10/24"
  end_address   = "10.0.0.20/24"
  status        = "active"
}

data "netbox_ip_range" "test" {
  id = netbox_ip_range.test.id
}
`
