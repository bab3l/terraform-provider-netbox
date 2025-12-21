package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderAccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "name", "Test Account"),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "account", "1234567890"),
				),
			},
		},
	})
}

const testAccProviderAccountDataSourceConfig = `
resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  account          = "1234567890"
  name             = "Test Account"
}

data "netbox_provider_account" "test" {
  id = netbox_provider_account.test.id
}
`
