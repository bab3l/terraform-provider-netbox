package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVMInterfaceDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_cluster_type" "test" {
					name = "test-cluster-type"
					slug = "test-cluster-type"
				}

				resource "netbox_cluster" "test" {
					name = "test-cluster"
					type = netbox_cluster_type.test.id
				}

				resource "netbox_virtual_machine" "test" {
					name    = "test-vm"
					cluster = netbox_cluster.test.id
					status  = "active"
				}

				resource "netbox_vm_interface" "test" {
					name            = "test-interface"
					virtual_machine = netbox_virtual_machine.test.name
				}

				data "netbox_vm_interface" "test" {
					name            = netbox_vm_interface.test.name
					virtual_machine = netbox_virtual_machine.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "name", "test-interface"),
					resource.TestCheckResourceAttrSet("data.netbox_vm_interface.test", "id"),
				),
			},
		},
	})
}
