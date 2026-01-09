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

func TestAccClusterTypeDataSource_customFields(t *testing.T) {
	clusterTypeName := testutil.RandomName("tf-test-clustertype-ds-cf")
	clusterTypeSlug := testutil.GenerateSlug(clusterTypeName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_clustertype_ds_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeDataSourceConfigWithCustomFields(clusterTypeName, clusterTypeSlug, customFieldName, customFieldValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccClusterTypeDataSourceConfigWithCustomFields(clusterTypeName, clusterTypeSlug, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[3]q
  object_types = ["virtualization.clustertype"]
  type         = "text"
}

resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[4]q
    }
  ]
}

data "netbox_cluster_type" "test" {
  name = netbox_cluster_type.test.name
  depends_on = [netbox_cluster_type.test]
}
`, clusterTypeName, clusterTypeSlug, customFieldName, customFieldValue)
}
