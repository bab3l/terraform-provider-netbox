package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVRFDataSource_basic(t *testing.T) {

	t.Parallel()

	vrfName := testutil.RandomName("vrf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFDataSourceConfig(vrfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "name", vrfName),
					resource.TestCheckResourceAttrSet("data.netbox_vrf.test", "id"),
				),
			},
		},
	})
}

func testAccVRFDataSourceConfig(vrfName string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
	name = "%s"
}

data "netbox_vrf" "test" {
	name = netbox_vrf.test.name
}
`, vrfName)
}
