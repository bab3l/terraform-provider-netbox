package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider.test", "name", "Test Provider"),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "slug", "test-provider"),
				),
			},
		},
	})
}

const testAccProviderDataSourceConfig = `
resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

data "netbox_provider" "test" {
  id = netbox_provider.test.id
}
`
