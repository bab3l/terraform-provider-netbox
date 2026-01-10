//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterGroupDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_clustergroup_ds_cf")
	clusterGroupName := testutil.RandomName("tf-test-clustergroup-ds-cf")
	clusterGroupSlug := testutil.RandomSlug("tf-test-clustergroup-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupDataSourceConfig_customFields(customFieldName, clusterGroupName, clusterGroupSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "name", clusterGroupName),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccClusterGroupDataSourceConfig_customFields(customFieldName, clusterGroupName, clusterGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["virtualization.clustergroup"]
  type         = "text"
}

resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_cluster_group" "test" {
  name = %q

  depends_on = [netbox_cluster_group.test]
}
`, customFieldName, clusterGroupName, clusterGroupSlug, clusterGroupName)
}
