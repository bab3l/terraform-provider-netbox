package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVMInterfaceDataSource_byID(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")
	interfaceName := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckVMInterfaceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceDataSourceByIDConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttrSet("data.netbox_vm_interface.test", "id"),
				),
			},
		},
	})
}

func TestAccVMInterfaceDataSource_byName(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")
	interfaceName := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckVMInterfaceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceDataSourceByNameConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttrSet("data.netbox_vm_interface.test", "id"),
				),
			},
		},
	})
}

func testAccVMInterfaceDataSourceByIDConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName string) string {
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

resource "netbox_vm_interface" "test" {
	name            = "%s"
	virtual_machine = netbox_virtual_machine.test.name
}

data "netbox_vm_interface" "test" {
	id = netbox_vm_interface.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName)
}

func testAccVMInterfaceDataSourceByNameConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName string) string {
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

resource "netbox_vm_interface" "test" {
	name            = "%s"
	virtual_machine = netbox_virtual_machine.test.name
}

data "netbox_vm_interface" "test" {
	name            = netbox_vm_interface.test.name
	virtual_machine = netbox_virtual_machine.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName)
}
