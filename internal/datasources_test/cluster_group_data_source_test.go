package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterGroupDataSource_basic(t *testing.T) {
	name := testutil.RandomName("test-cluster-group")
	slug := testutil.GenerateSlug(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "description", "Test Cluster Group Description"),
				),
			},
		},
	})
}

func testAccClusterGroupDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
  name        = %q
  slug        = %q
  description = "Test Cluster Group Description"
}

data "netbox_cluster_group" "test" {
  id = netbox_cluster_group.test.id
}
`, name, slug)
}
