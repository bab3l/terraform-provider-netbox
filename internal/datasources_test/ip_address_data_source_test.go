package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "address", "10.0.0.1/24"),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "status", "active"),
				),
			},
		},
	})
}

const testAccIPAddressDataSourceConfig = `
resource "netbox_ip_address" "test" {
  address = "10.0.0.1/24"
  status  = "active"
}

data "netbox_ip_address" "test" {
  address = netbox_ip_address.test.address
}
`
