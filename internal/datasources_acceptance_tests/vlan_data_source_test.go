package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANDataSource_byID(t *testing.T) {

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
				Config: testAccVLANDataSourceByIDConfig(vlanName, int(vlanVID)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVID)),
					resource.TestCheckResourceAttrSet("data.netbox_vlan.test", "id"),
				),
			},
		},
	})
}

func TestAccVLANDataSource_byVID(t *testing.T) {

	t.Parallel()

	vlanName := testutil.RandomName("vlan")
	vlanVID := int32(1001)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANDataSourceByVIDConfig(vlanName, int(vlanVID)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVID)),
					resource.TestCheckResourceAttrSet("data.netbox_vlan.test", "id"),
				),
			},
		},
	})
}

func TestAccVLANDataSource_byName(t *testing.T) {

	t.Parallel()

	vlanName := testutil.RandomName("vlan")
	vlanVID := int32(1002)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANDataSourceByNameConfig(vlanName, int(vlanVID)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVID)),
					resource.TestCheckResourceAttrSet("data.netbox_vlan.test", "id"),
				),
			},
		},
	})
}

func testAccVLANDataSourceByIDConfig(vlanName string, vlanVID int) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
	name = "%s"
	vid  = %d
}

data "netbox_vlan" "test" {
	id = netbox_vlan.test.id
}
`, vlanName, vlanVID)
}

func testAccVLANDataSourceByVIDConfig(vlanName string, vlanVID int) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
	name = "%s"
	vid  = %d
}

data "netbox_vlan" "test" {
	vid = netbox_vlan.test.vid
}
`, vlanName, vlanVID)
}

func testAccVLANDataSourceByNameConfig(vlanName string, vlanVID int) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
	name = "%s"
	vid  = %d
}

data "netbox_vlan" "test" {
	name = netbox_vlan.test.name
}
`, vlanName, vlanVID)
}
