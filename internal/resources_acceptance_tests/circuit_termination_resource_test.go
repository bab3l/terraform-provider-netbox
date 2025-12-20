package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTerminationResource_basic(t *testing.T) {
	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	circuitTypeName := testutil.RandomName("tf-test-ct")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")
	circuitCID := testutil.RandomName("tf-test-circuit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				ResourceName:            "netbox_circuit_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit", "site"},
			},
		},
	})
}

func TestAccCircuitTerminationResource_full(t *testing.T) {
	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	circuitTypeName := testutil.RandomName("tf-test-ct-full")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-full")
	circuitCID := testutil.RandomName("tf-test-circuit-full")
	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	description := "Test circuit termination with all fields"
	updatedDescription := "Updated circuit termination description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, 1000000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "1000000"),
				),
			},
			{
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, updatedDescription, 10000000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
				),
			},
		},
	})
}

func testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.id
  term_side = "A"
  site      = netbox_site.test.id
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug)
}

func testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description string, portSpeed int) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_circuit_termination" "test" {
  circuit     = netbox_circuit.test.id
  term_side   = "A"
  site        = netbox_site.test.id
  port_speed  = %d
  description = %q
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, portSpeed, description)
}

func TestAccCircuitTerminationResource_import(t *testing.T) {
	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	circuitTypeName := testutil.RandomName("tf-test-ct")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")
	circuitCID := testutil.RandomName("tf-test-circuit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				ResourceName:            "netbox_circuit_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit", "site"},
			},
		},
	})
}

// TestAccConsistency_CircuitTermination_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_CircuitTermination_LiteralNames(t *testing.T) {
	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	circuitTypeName := testutil.RandomName("circuit-type")
	circuitTypeSlug := testutil.RandomSlug("circuit-type")
	circuitCid := testutil.RandomName("CID")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(circuitCid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationConsistencyLiteralNamesConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "circuit", circuitCid),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "site", siteName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccCircuitTerminationConsistencyLiteralNamesConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug),
			},
		},
	})
}

func testAccCircuitTerminationConsistencyLiteralNamesConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_circuit_type" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_circuit" "test" {
  cid = "%[5]s"
  circuit_provider = netbox_provider.test.id
  type = netbox_circuit_type.test.slug
}

resource "netbox_site" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_circuit_termination" "test" {
  # Use literal string names to mimic existing user state
  circuit = "%[5]s"
  term_side = "A"
  site = "%[6]s"
  depends_on = [netbox_circuit.test, netbox_site.test]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug)
}
