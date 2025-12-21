package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_prefix.test", "prefix", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("data.netbox_prefix.test", "status", "active"),
				),
			},
		},
	})
}

const testAccPrefixDataSourceConfig = `
resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/24"
  status = "active"
}

data "netbox_prefix" "test" {
  prefix = netbox_prefix.test.prefix
}
`
