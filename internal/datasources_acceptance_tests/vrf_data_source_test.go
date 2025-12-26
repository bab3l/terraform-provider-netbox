package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVRFDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	vrfName := testutil.RandomName("vrf-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFDataSourceByIDConfig(vrfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "name", vrfName),
				),
			},
		},
	})
}

func TestAccVRFDataSource_byID(t *testing.T) {

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
				Config: testAccVRFDataSourceByIDConfig(vrfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "name", vrfName),
					resource.TestCheckResourceAttrSet("data.netbox_vrf.test", "id"),
				),
			},
		},
	})
}

func TestAccVRFDataSource_byName(t *testing.T) {

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
				Config: testAccVRFDataSourceByNameConfig(vrfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "name", vrfName),
					resource.TestCheckResourceAttrSet("data.netbox_vrf.test", "id"),
				),
			},
		},
	})
}

func testAccVRFDataSourceByIDConfig(vrfName string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
	name = "%s"
}

data "netbox_vrf" "test" {
	id = netbox_vrf.test.id
}
`, vrfName)
}

func testAccVRFDataSourceByNameConfig(vrfName string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
	name = "%s"
}

data "netbox_vrf" "test" {
	name = netbox_vrf.test.name
}
`, vrfName)
}
