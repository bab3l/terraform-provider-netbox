//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterDataSource_customFields(t *testing.T) {
	clusterTypeName := testutil.RandomName("tf-test-clustertype-ds-cf")
	clusterTypeSlug := testutil.GenerateSlug(clusterTypeName)
	clusterName := testutil.RandomName("tf-test-cluster-ds-cf")
	customFieldName := testutil.RandomCustomFieldName("tf_test_cluster_ds_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfigWithCustomFields(clusterTypeName, clusterTypeSlug, clusterName, customFieldName, customFieldValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_cluster.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cluster.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_cluster.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_cluster.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfigWithCustomFields(clusterTypeName, clusterTypeSlug, clusterName, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_custom_field" "test" {
  name         = %[4]q
  object_types = ["virtualization.cluster"]
  type         = "text"
}

resource "netbox_cluster" "test" {
  name = %[3]q
  type = netbox_cluster_type.test.id
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[5]q
    }
  ]
}

data "netbox_cluster" "test" {
  name = netbox_cluster.test.name
  depends_on = [netbox_cluster.test]
}
`, clusterTypeName, clusterTypeSlug, clusterName, customFieldName, customFieldValue)
}
