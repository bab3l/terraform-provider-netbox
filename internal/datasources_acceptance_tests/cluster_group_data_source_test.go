package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterGroupDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cg-ds-id-group")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterGroupDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "description", "Test Cluster Group Description"),
				),
			},
		},
	})
}

func TestAccClusterGroupDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-cluster-group")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterGroupDestroy,
		),
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

func TestAccClusterGroupDataSource_byName(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("Public Cloud %s", testutil.RandomName("test-cluster-group"))
	slug := testutil.RandomSlug("test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterGroupDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccClusterGroupDataSource_bySlug(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-cluster-group")
	slug := testutil.RandomSlug("test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterGroupDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupDataSourceConfigBySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_cluster_group.test", "slug", slug),
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

func testAccClusterGroupDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
	name        = %q
	slug        = %q
	description = "Test Cluster Group Description"
}

data "netbox_cluster_group" "test" {
	name = netbox_cluster_group.test.name
}
`, name, slug)
}

func testAccClusterGroupDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
	name        = %q
	slug        = %q
	description = "Test Cluster Group Description"
}

data "netbox_cluster_group" "test" {
	slug = netbox_cluster_group.test.slug
}
`, name, slug)
}
