package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualDiskDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-vdisk-ds")

	clusterTypeName := testutil.RandomName("tf-test-ct")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

	clusterName := testutil.RandomName("tf-test-cluster")

	vmName := testutil.RandomName("tf-test-vm")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualDiskCleanup(name)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVirtualDiskDestroy,

			testutil.CheckVirtualMachineDestroy,

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskDataSourceConfig(name, clusterTypeName, clusterTypeSlug, clusterName, vmName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "size", "100"),
				),
			},
		},
	})

}

func testAccVirtualDiskDataSourceConfig(name, clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_cluster_type" "test" {

  name = %q

  slug = %q

}

resource "netbox_cluster" "test" {

  name = %q

  type = netbox_cluster_type.test.id

}

resource "netbox_virtual_machine" "test" {

  name    = %q

  cluster = netbox_cluster.test.id

  # Ignore disk changes since Netbox auto-computes this from virtual_disks

  lifecycle {

    ignore_changes = [disk]

  }

}

resource "netbox_virtual_disk" "test" {

  virtual_machine = netbox_virtual_machine.test.id

  name            = %q

  size            = 100

}

data "netbox_virtual_disk" "test" {

  id = netbox_virtual_disk.test.id

}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName, name)

}

// ASN Range Data Source Tests
