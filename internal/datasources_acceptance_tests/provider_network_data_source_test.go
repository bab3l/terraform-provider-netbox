package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderNetworkDataSource_basic(t *testing.T) {

	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "name", "Test Provider Network"),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "service_id", "12345"),
				),
			},
		},
	})
}

const testAccProviderNetworkDataSourceConfig = `
resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = "Test Provider Network"
  service_id       = "12345"
}

data "netbox_provider_network" "test" {
  id = netbox_provider_network.test.id
}
`
