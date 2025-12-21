package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterTypeDataSource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cluster-type")
	slug := testutil.RandomSlug("cluster-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.test", "slug", slug),
				),
			},
		},
	})
}

func testAccClusterTypeDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_cluster_type" "test" {
  id = netbox_cluster_type.test.id
}
`, name, slug)
}
