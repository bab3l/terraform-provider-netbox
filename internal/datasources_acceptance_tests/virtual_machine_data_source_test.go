package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachineDataSource_basic(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineDataSourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.test", "status", "active"),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_machine.test", "id"),
				),
			},
		},
	})
}

func testAccVirtualMachineDataSourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_cluster" "test" {
	name = "%s"
	type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
	name    = "%s"
	cluster = netbox_cluster.test.id
	status  = "active"
}

data "netbox_virtual_machine" "test" {
	name = netbox_virtual_machine.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}
