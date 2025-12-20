package resources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccVirtualMachineResource_customFieldsWithDigit tests that custom field names

// can start with digits (e.g., "4me_name"), matching Netbox's actual validation.

func TestAccVirtualMachineResource_customFieldsWithDigit(t *testing.T) {

	t.Skip("Requires custom fields to be pre-configured in Netbox")

	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")

	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	clusterName := testutil.RandomName("cluster")

	vmName := testutil.RandomName("vm")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualMachineCustomFieldsWithDigitConfig(

					clusterTypeName, clusterTypeSlug, clusterName, vmName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.#", "3"),
				),
			},
		},
	})

}

func testAccVirtualMachineCustomFieldsWithDigitConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name = %[1]q

  slug = %[2]q

}



resource "netbox_cluster" "test" {

  name = %[3]q

  type = netbox_cluster_type.test.id

}



resource "netbox_virtual_machine" "test" {

  name    = %[4]q

  cluster = netbox_cluster.test.id

  status  = "active"



  # Test custom field names that start with digits

  custom_fields = [

    {

      name  = "4me_name"

      type  = "text"

      value = "test-value-1"

    },

    {

      name  = "2factor_enabled"

      type  = "boolean"

      value = "true"

    },

    {

      name  = "normal_field"

      type  = "text"

      value = "test-value-2"

    }

  ]

}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName)

}
