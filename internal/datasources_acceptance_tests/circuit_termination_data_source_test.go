package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTerminationDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("ct-ds-id-site")
	siteSlug := testutil.RandomSlug("ct-ds-id-site")
	providerName := testutil.RandomName("ct-ds-id-provider")
	providerSlug := testutil.RandomSlug("ct-ds-id-provider")
	circuitTypeName := testutil.RandomName("ct-ds-id-circuit-type")
	circuitTypeSlug := testutil.RandomSlug("ct-ds-id-circuit-type")
	cid := testutil.RandomName("ct-ds-id-circuit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(cid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckProviderDestroy,
			testutil.CheckCircuitTypeDestroy,
			testutil.CheckCircuitDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationDataSourceConfig(siteName, siteSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "term_side", "A"),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "port_speed", "1000"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "circuit"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "site"),
				),
			},
		},
	})
}

func TestAccCircuitTerminationDataSource_byID(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	circuitTypeName := testutil.RandomName("circuit-type")
	circuitTypeSlug := testutil.RandomSlug("circuit-type")
	cid := testutil.RandomName("circuit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(cid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckProviderDestroy,
			testutil.CheckCircuitTypeDestroy,
			testutil.CheckCircuitDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationDataSourceConfig(siteName, siteSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "term_side", "A"),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "port_speed", "1000"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "circuit"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "site"),
				),
			},
		},
	})
}

func testAccCircuitTerminationDataSourceConfig(siteName, siteSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

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

resource "netbox_circuit_termination" "test" {
  circuit    = netbox_circuit.test.id
  term_side  = "A"
  site       = netbox_site.test.id
  port_speed = 1000
}

data "netbox_circuit_termination" "test" {
  id = netbox_circuit_termination.test.id
}
`, siteName, siteSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, cid)
}
