package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterDataSource_byID(t *testing.T) {

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
					resource.TestCheckResourceAttr("data.netbox_cluster.by_id", "name", clusterName),
					resource.TestCheckResourceAttrSet("data.netbox_cluster.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_cluster.by_id", "type"),
				),
			},
		},
	})
}

func TestAccClusterDataSource_byName(t *testing.T) {

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
					resource.TestCheckResourceAttr("data.netbox_cluster.by_name", "name", clusterName),
					resource.TestCheckResourceAttrSet("data.netbox_cluster.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_cluster.by_name", "type"),
				),
			},
		},
	})
}

func TestAccClusterDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type-id")
	clusterTypeSlug := testutil.RandomSlug("cluster-type-id")
	clusterName := testutil.RandomName("cluster-id")

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
					resource.TestCheckResourceAttrSet("data.netbox_cluster.by_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_cluster.by_name", "name", clusterName),
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

data "netbox_cluster" "by_id" {
  id = netbox_cluster.test.id
}

data "netbox_cluster" "by_name" {
  name = netbox_cluster.test.name
}
`, typeName, typeSlug, name)
}
