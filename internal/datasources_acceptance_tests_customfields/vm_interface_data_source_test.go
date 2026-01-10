//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVMInterfaceDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_vminterface_ds_cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.RandomSlug("tf-test-site-ds-cf")
	clusterName := testutil.RandomName("tf-test-cluster-ds-cf")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type-ds-cf")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-ds-cf")
	vmName := testutil.RandomName("tf-test-vm-ds-cf")
	ifaceName := "eth0"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceDataSourceConfig_customFields(customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_vm_interface.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccVMInterfaceDataSourceConfig_customFields(customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["virtualization.vminterface"]
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
}

resource "netbox_vm_interface" "test" {
  name            = %q
  virtual_machine = netbox_virtual_machine.test.name

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %q

  depends_on = [netbox_vm_interface.test]
}
`, customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, ifaceName)
}
