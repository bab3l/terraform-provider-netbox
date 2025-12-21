package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNDataSource_basic(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	asnValue := 65000

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRIRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASNDataSourceConfig(rirName, rirSlug, asnValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_asn.test", "asn", fmt.Sprintf("%d", asnValue)),
					resource.TestCheckResourceAttrSet("data.netbox_asn.test", "rir"),
				),
			},
		},
	})
}

func testAccASNDataSourceConfig(rirName, rirSlug string, asn int) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
}

data "netbox_asn" "test" {
  id = netbox_asn.test.id
}
`, rirName, rirSlug, asn)
}
