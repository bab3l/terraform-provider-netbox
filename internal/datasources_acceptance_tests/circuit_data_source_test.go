package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitDataSource_byID(t *testing.T) {

	t.Parallel()

	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	circuitTypeName := testutil.RandomName("circuit-type")
	circuitTypeSlug := testutil.RandomSlug("circuit-type")
	cid := testutil.RandomName("circuit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckProviderDestroy,
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitDataSourceConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit.by_id", "cid", cid),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_id", "circuit_provider"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_id", "type"),
				),
			},
		},
	})
}

func TestAccCircuitDataSource_byCID(t *testing.T) {

	t.Parallel()

	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	circuitTypeName := testutil.RandomName("circuit-type")
	circuitTypeSlug := testutil.RandomSlug("circuit-type")
	cid := testutil.RandomName("circuit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckProviderDestroy,
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitDataSourceConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit.by_cid", "cid", cid),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_cid", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_cid", "circuit_provider"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_cid", "type"),
				),
			},
		},
	})
}

func TestAccCircuitDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-id")
	providerSlug := testutil.RandomSlug("tf-test-provider-id")
	circuitTypeName := testutil.RandomName("tf-test-circuit-type-id")
	circuitTypeSlug := testutil.RandomSlug("tf-test-circuit-type-id")
	cid := testutil.RandomName("tf-test-circuit-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckProviderDestroy,
			testutil.CheckCircuitTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitDataSourceConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns ID correctly
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_cid", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit.by_cid", "cid", cid),
					// Verify reference attributes are set
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_cid", "circuit_provider"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.by_cid", "type"),
				),
			},
		},
	})
}

func testAccCircuitDataSourceConfig(providerName, providerSlug, typeName, typeSlug, cid string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_circuit_type" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_circuit" "test" {
  cid              = "%s"
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
}

data "netbox_circuit" "by_cid" {
  cid = netbox_circuit.test.cid
}

data "netbox_circuit" "by_id" {
  id = netbox_circuit.test.id
}
`, providerName, providerSlug, typeName, typeSlug, cid)
}
