package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualChassisDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("vc-id")

	cleanup.RegisterVirtualChassisCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualChassisDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisDataSourceByNameConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_chassis.test", "name", name),
				),
			},
		},
	})
}

func TestAccVirtualChassisDataSource_byName(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("vc")

	cleanup.RegisterVirtualChassisCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualChassisDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisDataSourceByNameConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_chassis.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_chassis.test", "id"),
				),
			},
		},
	})
}

func TestAccVirtualChassisDataSource_byID(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("vc")

	cleanup.RegisterVirtualChassisCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualChassisDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisDataSourceByIDConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_chassis.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_chassis.test", "id"),
				),
			},
		},
	})
}

func testAccVirtualChassisDataSourceByNameConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
	name = "%s"
}

data "netbox_virtual_chassis" "test" {
	name = netbox_virtual_chassis.test.name
}
`, name)
}

func testAccVirtualChassisDataSourceByIDConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
	name = "%s"
}

data "netbox_virtual_chassis" "test" {
	id = netbox_virtual_chassis.test.id
}
`, name)
}
