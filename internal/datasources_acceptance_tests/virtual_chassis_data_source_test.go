package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualChassisDataSource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("vc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_chassis.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_chassis.test", "id"),
				),
			},
		},
	})
}

func testAccVirtualChassisDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
	name = "%s"
}

data "netbox_virtual_chassis" "test" {
	name = netbox_virtual_chassis.test.name
}
`, name)
}
