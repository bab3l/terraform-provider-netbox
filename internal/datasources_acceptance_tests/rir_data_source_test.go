package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRIRDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRIRDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rir.test", "name", "Test RIR"),
					resource.TestCheckResourceAttr("data.netbox_rir.test", "slug", "test-rir"),
				),
			},
		},
	})
}

const testAccRIRDataSourceConfig = `
resource "netbox_rir" "test" {
  name = "Test RIR"
  slug = "test-rir"
}

data "netbox_rir" "test" {
  id = netbox_rir.test.id
}
`
