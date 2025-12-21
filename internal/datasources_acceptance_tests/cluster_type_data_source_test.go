package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterTypeDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeDataSourceConfig("Test Cluster Type"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "name", "Test Cluster Type"),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "slug", "test-cluster-type"),
				),
			},
		},
	})
}

func testAccClusterTypeDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%s"
  slug = "test-cluster-type"
}

data "netbox_cluster_type" "test" {
  id = netbox_cluster_type.test.id
}
`, name)
}
