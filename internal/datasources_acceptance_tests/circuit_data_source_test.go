package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitDataSource_basic(t *testing.T) {

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
					resource.TestCheckResourceAttr("data.netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.test", "circuit_provider"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.test", "type"),
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

data "netbox_circuit" "test" {
  id = netbox_circuit.test.id
}
`, providerName, providerSlug, typeName, typeSlug, cid)
}
