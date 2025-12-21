package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANDataSource_basic(t *testing.T) {

	t.Parallel()

	vlanName := testutil.RandomName("vlan")
	vlanVID := int32(1000)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANDataSourceConfig(vlanName, int(vlanVID)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVID)),
					resource.TestCheckResourceAttrSet("data.netbox_vlan.test", "id"),
				),
			},
		},
	})
}

func testAccVLANDataSourceConfig(vlanName string, vlanVID int) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
	name = "%s"
	vid  = %d
}

data "netbox_vlan" "test" {
	name = netbox_vlan.test.name
	vid  = netbox_vlan.test.vid
}
`, vlanName, vlanVID)
}
