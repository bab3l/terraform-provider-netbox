package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateDataSource_basic(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRIRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateDataSourceConfig(rirName, rirSlug, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("data.netbox_aggregate.test", "rir"),
				),
			},
		},
	})
}

func testAccAggregateDataSourceConfig(rirName, rirSlug, prefix string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_aggregate" "test" {
  prefix = "%s"
  rir    = netbox_rir.test.id
}

data "netbox_aggregate" "test" {
  id = netbox_aggregate.test.id
}
`, rirName, rirSlug, prefix)
}
