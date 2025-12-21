package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANGroupDataSource_basic(t *testing.T) {

	t.Parallel()

	vlanGroupName := testutil.RandomName("vlan-group")
	vlanGroupSlug := testutil.RandomSlug("vlan-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(vlanGroupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupDataSourceConfig(vlanGroupName, vlanGroupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan_group.test", "name", vlanGroupName),
					resource.TestCheckResourceAttr("data.netbox_vlan_group.test", "slug", vlanGroupSlug),
					resource.TestCheckResourceAttrSet("data.netbox_vlan_group.test", "id"),
				),
			},
		},
	})
}

func testAccVLANGroupDataSourceConfig(vlanGroupName, vlanGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_vlan_group" "test" {
	name = "%s"
	slug = "%s"
}

data "netbox_vlan_group" "test" {
	name = netbox_vlan_group.test.name
}
`, vlanGroupName, vlanGroupSlug)
}
