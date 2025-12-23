package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNDataSource_basic(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	asnValue := acctest.RandIntRange(65000, 65999) // Generate ASN between 65000-65999

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
					// Test lookup by asn
					resource.TestCheckResourceAttr("data.netbox_asn.by_asn", "asn", fmt.Sprintf("%d", asnValue)),
					resource.TestCheckResourceAttrSet("data.netbox_asn.by_asn", "rir"),
					resource.TestCheckResourceAttrSet("data.netbox_asn.by_asn", "id"),
					// Test lookup by id
					resource.TestCheckResourceAttr("data.netbox_asn.by_id", "asn", fmt.Sprintf("%d", asnValue)),
					resource.TestCheckResourceAttrSet("data.netbox_asn.by_id", "rir"),
					resource.TestCheckResourceAttrSet("data.netbox_asn.by_id", "id"),
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

data "netbox_asn" "by_asn" {
  asn = netbox_asn.test.asn
}

data "netbox_asn" "by_id" {
  id = netbox_asn.test.id
}
`, rirName, rirSlug, asn)
}
