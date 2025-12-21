package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterDataSource_basic(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cluster.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("data.netbox_cluster.test", "type"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig(typeName, typeSlug, name string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_cluster" "test" {
  name = "%s"
  type = netbox_cluster_type.test.id
}

data "netbox_cluster" "test" {
  id = netbox_cluster.test.id
}
`, typeName, typeSlug, name)
}
