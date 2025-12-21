package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateDataSourceConfig("10.0.0.0/8"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "prefix", "10.0.0.0/8"),
					resource.TestCheckResourceAttrSet("data.netbox_aggregate.test", "rir"),
				),
			},
		},
	})
}

func testAccAggregateDataSourceConfig(prefix string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "Test RIR"
  slug = "test-rir"
}

resource "netbox_aggregate" "test" {
  prefix = "%s"
  rir    = netbox_rir.test.id
}

data "netbox_aggregate" "test" {
  id = netbox_aggregate.test.id
}
`, prefix)
}
