package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterTypeDataSource_byID(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_id", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_cluster_type.by_id", "id"),
				),
			},
		},
	})
}

func TestAccClusterTypeDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("Public Cloud")
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
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_name", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_cluster_type.by_name", "id"),
				),
			},
		},
	})
}

func TestAccClusterTypeDataSource_bySlug(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_slug", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_slug", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_cluster_type.by_slug", "id"),
				),
			},
		},
	})
}

func TestAccClusterTypeDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cluster-type-id")
	slug := testutil.RandomSlug("cluster-type-id")

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
					resource.TestCheckResourceAttrSet("data.netbox_cluster_type.by_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_type.by_name", "slug", slug),
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

data "netbox_cluster_type" "by_id" {
  id = netbox_cluster_type.test.id
}

data "netbox_cluster_type" "by_name" {
  name = netbox_cluster_type.test.name
}

data "netbox_cluster_type" "by_slug" {
  slug = netbox_cluster_type.test.slug
}
`, name, slug)
}
