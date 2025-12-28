package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTerminationResource_basic(t *testing.T) {

	t.Parallel()
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

	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	circuitTypeName := testutil.RandomName("tf-test-ct-full")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-full")
	circuitCID := testutil.RandomName("tf-test-circuit-full")
	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	description := testutil.RandomName("description")
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

func TestAccCircuitTerminationResource_update(t *testing.T) {
	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-upd")
	providerSlug := testutil.RandomSlug("tf-test-provider-upd")
	circuitTypeName := testutil.RandomName("tf-test-ct-upd")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-upd")
	circuitCID := testutil.RandomName("tf-test-circuit-upd")
	siteName := testutil.RandomName("tf-test-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-site-upd")

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
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, testutil.Description2, 10000000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", testutil.Description2),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
				),
			},
		},
	})
}

func TestAccCircuitTerminationResource_IDPreservation(t *testing.T) {

	t.Parallel()
	providerName := testutil.RandomName("ct-prov-id")
	providerSlug := testutil.RandomSlug("ct-prov-id")
	circuitTypeName := testutil.RandomName("ct-type-id")
	circuitTypeSlug := testutil.RandomSlug("ct-type-id")
	circuitCID := testutil.RandomName("ct-ckt-id")
	siteName := testutil.RandomName("ct-site-id")
	siteSlug := testutil.RandomSlug("ct-site-id")

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

	t.Parallel()
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

	t.Parallel()
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
  type = netbox_circuit_type.test.id
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

func TestAccCircuitTerminationResource_externalDeletion(t *testing.T) {
	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-ext-del")
	providerSlug := testutil.RandomSlug("tf-test-provider-ext-del")
	circuitTypeName := testutil.RandomName("tf-test-ct-ext-del")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-ext-del")
	circuitCID := testutil.RandomName("tf-test-circuit-ext-del")
	siteName := testutil.RandomName("tf-test-site-ext-del")
	siteSlug := testutil.RandomSlug("tf-test-site-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
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
  name = %q
  slug = %q
}
resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.cid
  term_side = "A"
  site      = netbox_site.test.slug
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List terminations filtered by site slug
					items, _, err := client.CircuitsAPI.CircuitsCircuitTerminationsList(context.Background()).Site([]string{siteSlug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit termination for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitTerminationsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit termination: %v", err)
					}
					t.Logf("Successfully externally deleted circuit termination with ID: %d", itemID)
				},
				Config: fmt.Sprintf(`
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
  name = %q
  slug = %q
}
resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.cid
  term_side = "A"
  site      = netbox_site.test.slug
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				ExpectNonEmptyPlan: true,
				RefreshState:       true,
			},
		},
	})
}
