package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTypeDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("ct-ds-id-type")
	slug := testutil.RandomSlug("ct-ds-id-type")

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_type.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_id", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitTypeDataSource_byID(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("circuit-type")
	slug := testutil.RandomSlug("circuit-type")

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_id", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_type.by_id", "id"),
				),
			},
		},
	})
}

func TestAccCircuitTypeDataSource_byName(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("circuit-type")
	slug := testutil.RandomSlug("circuit-type")

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_name", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_type.by_name", "id"),
				),
			},
		},
	})
}

func TestAccCircuitTypeDataSource_bySlug(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("circuit-type")
	slug := testutil.RandomSlug("circuit-type")

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_slug", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.by_slug", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_type.by_slug", "id"),
				),
			},
		},
	})
}

func testAccCircuitTypeDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_circuit_type" "by_id" {
  id = netbox_circuit_type.test.id
}

data "netbox_circuit_type" "by_name" {
  name = netbox_circuit_type.test.name
}

data "netbox_circuit_type" "by_slug" {
  slug = netbox_circuit_type.test.slug
}
`, name, slug)
}
