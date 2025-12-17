package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNDataSourceConfig(65001),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_asn.test", "asn", "65001"),
					resource.TestCheckResourceAttrSet("data.netbox_asn.test", "rir"),
				),
			},
		},
	})
}

func testAccASNDataSourceConfig(asn int) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "Test RIR"
  slug = "test-rir"
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
}

data "netbox_asn" "test" {
  id = netbox_asn.test.id
}
`, asn)
}
