package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachinesDataSource_byNameFilter(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type-q")
	clusterTypeSlug := testutil.RandomSlug("cluster-type-q")
	clusterName := testutil.RandomName("cluster-q")
	vmName := testutil.RandomName("vm-q")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinesDataSourceConfig_byName(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.0", vmName),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "ids.0", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "virtual_machines.0.id", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.0.name", vmName),
				),
			},
		},
	})
}

func TestAccVirtualMachinesDataSource_byTagFilter(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type-q-tag")
	clusterTypeSlug := testutil.RandomSlug("cluster-type-q-tag")
	clusterName := testutil.RandomName("cluster-q-tag")
	vmName := testutil.RandomName("vm-q-tag")
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinesDataSourceConfig_byTag(clusterTypeName, clusterTypeSlug, clusterName, vmName, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.0", vmName),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "ids.0", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "virtual_machines.0.id", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.0.name", vmName),
				),
			},
		},
	})
}

func TestAccVirtualMachinesDataSource_byNameAndTagFilters(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type-q-multi")
	clusterTypeSlug := testutil.RandomSlug("cluster-type-q-multi")
	clusterName := testutil.RandomName("cluster-q-multi")
	vmName := testutil.RandomName("vm-q-multi")
	tagName := testutil.RandomName("tf-test-tag-multi")
	tagSlug := testutil.RandomSlug("tf-test-tag-multi")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinesDataSourceConfig_byNameAndTag(clusterTypeName, clusterTypeSlug, clusterName, vmName, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "names.0", vmName),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "ids.0", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_virtual_machines.test", "virtual_machines.0.id", "netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machines.test", "virtual_machines.0.name", vmName),
				),
			},
		},
	})
}

func testAccVirtualMachinesDataSourceConfig_byName(clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %q
  cluster = netbox_cluster.test.id
  status  = "active"
}

data "netbox_virtual_machines" "test" {
  filter {
    name   = "name"
    values = [netbox_virtual_machine.test.name]
  }
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}

func testAccVirtualMachinesDataSourceConfig_byTag(clusterTypeName, clusterTypeSlug, clusterName, vmName, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
	name = %q
	slug = %q
}

resource "netbox_cluster_type" "test" {
	name = %q
	slug = %q
}

resource "netbox_cluster" "test" {
	name = %q
	type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
	name    = %q
	cluster = netbox_cluster.test.id
	status  = "active"

	tags = [
		{
			name = netbox_tag.test.name
			slug = netbox_tag.test.slug
		}
	]
}

data "netbox_virtual_machines" "test" {
	filter {
		name   = "tag"
		values = [netbox_tag.test.slug]
	}

	depends_on = [netbox_virtual_machine.test]
}
`, tagName, tagSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}

func testAccVirtualMachinesDataSourceConfig_byNameAndTag(clusterTypeName, clusterTypeSlug, clusterName, vmName, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
	name = %q
	slug = %q
}

resource "netbox_cluster_type" "test" {
	name = %q
	slug = %q
}

resource "netbox_cluster" "test" {
	name = %q
	type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
	name    = %q
	cluster = netbox_cluster.test.id
	status  = "active"

	tags = [
		{
			name = netbox_tag.test.name
			slug = netbox_tag.test.slug
		}
	]
}

data "netbox_virtual_machines" "test" {
	filter {
		name   = "name"
		values = [netbox_virtual_machine.test.name]
	}

	filter {
		name   = "tag"
		values = [netbox_tag.test.slug]
	}
}
`, tagName, tagSlug, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}
