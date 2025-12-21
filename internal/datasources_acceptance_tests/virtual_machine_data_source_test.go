package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachineDataSource_basic(t *testing.T) {

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

				data "netbox_virtual_machine" "test" {
					name = netbox_virtual_machine.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.test", "name", "test-vm"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.test", "status", "active"),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_machine.test", "id"),
				),
			},
		},
	})
}
