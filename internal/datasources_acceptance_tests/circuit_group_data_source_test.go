package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitGroupDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-cg-ds-id-name")
	slug := testutil.RandomSlug("tf-test-cg-ds-id-name")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_byID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupDataSource_byID(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-cg-ds-id")
	slug := testutil.RandomSlug("tf-test-cg-ds-id")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_byID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupDataSource_byName(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := fmt.Sprintf("Public Cloud %s", testutil.RandomName("tf-test-cg-ds-name"))
	slug := testutil.RandomSlug("tf-test-cg-ds-name")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_byName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupDataSource_bySlug(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-cg-ds-slug")
	slug := testutil.RandomSlug("tf-test-cg-ds-slug")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_bySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccCircuitGroupDataSourceConfig_byID(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_circuit_group" "test" {
  id = netbox_circuit_group.test.id
}
`, name, slug)
}

func testAccCircuitGroupDataSourceConfig_byName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_circuit_group" "test" {
  name = netbox_circuit_group.test.name
}
`, name, slug)
}

func testAccCircuitGroupDataSourceConfig_bySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_circuit_group" "test" {
  slug = netbox_circuit_group.test.slug
}
`, name, slug)
}
