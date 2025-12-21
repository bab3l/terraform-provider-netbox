package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig("Test Cluster"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cluster.test", "name", "Test Cluster"),
					resource.TestCheckResourceAttrSet("data.netbox_cluster.test", "type"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "%s"
  type = netbox_cluster_type.test.id
}

data "netbox_cluster" "test" {
  id = netbox_cluster.test.id
}
`, name)
}
