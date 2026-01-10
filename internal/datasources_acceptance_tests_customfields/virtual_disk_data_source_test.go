//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualDiskDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_virtualdisk_ds_cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.RandomSlug("tf-test-site-ds-cf")
	clusterName := testutil.RandomName("tf-test-cluster-ds-cf")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type-ds-cf")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-ds-cf")
	vmName := testutil.RandomName("tf-test-vm-ds-cf")
	diskName := "disk0"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDiskDataSourceConfig_customFields(customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "name", diskName),
					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccVirtualDiskDataSourceConfig_customFields(customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["virtualization.virtualdisk"]
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
  disk    = 100
}

resource "netbox_virtual_disk" "test" {
  name            = %q
  virtual_machine = netbox_virtual_machine.test.name
  size            = 100

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_virtual_disk" "test" {
  name            = %q
  virtual_machine = netbox_virtual_machine.test.name

  depends_on = [netbox_virtual_disk.test]
}
`, customFieldName, siteName, siteSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName, diskName)
}
