//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachinesDataSource_queryWithCustomFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_vm_q_cf")
	siteName := testutil.RandomName("tf-test-site-q-cf")
	siteSlug := testutil.RandomSlug("tf-test-site-q-cf")
	clusterName := testutil.RandomName("tf-test-cluster-q-cf")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type-q-cf")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-q-cf")
	vmName := testutil.RandomName("tf-test-vm-q-cf")
	customFieldValue := testutil.RandomName("tf-test-cf-value")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(customFieldName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinesDataSourceConfig_withCustomFields(customFieldName, customFieldValue, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.0", vmName),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "virtual_machines.0.id", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.0.name", vmName),
				),
			},
		},
	})
}

func testAccVirtualMachinesDataSourceConfig_withCustomFields(customFieldName, customFieldValue, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["virtualization.virtualmachine"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.slug
  site = netbox_site.test.name
}

resource "netbox_virtual_machine" "test" {
  name    = %q
  cluster = netbox_cluster.test.name
  status  = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %q
    }
  ]
}

data "netbox_virtual_machines" "test" {
  filter {
		name   = "custom_field_value"
		values = ["${netbox_custom_field.test.name}=%s"]
  }

  depends_on = [netbox_virtual_machine.test]
}
`, customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, customFieldValue, customFieldValue)
}
