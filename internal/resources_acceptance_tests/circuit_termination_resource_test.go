package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
			{
				Config:             testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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
	tagName := testutil.RandomName("tf-test-tag-upd")
	tagSlug := testutil.RandomSlug("tf-test-tag-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				Config: testAccCircuitTerminationResourceConfig_withTags(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, testutil.Description2, tagName, tagSlug, 10000000, 5000000, "XCON-UPDATE", "PP-UPDATE", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", testutil.Description2),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "upstream_speed", "5000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "xconnect_id", "XCON-UPDATE"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "pp_info", "PP-UPDATE"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "true"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "tags.#", "1"),
				),
			},
		},
	})
}

func testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, tagName, tagSlug string) string {
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

resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_termination" "test" {
  circuit        = netbox_circuit.test.id
  term_side      = "A"
  site           = netbox_site.test.id
  port_speed     = 10000000
  upstream_speed = 5000000
  xconnect_id    = "XCON-FULL-123"
  pp_info        = "Patch Panel 1, Port 24"
  mark_connected = true
  description    = %q
	tags           = [netbox_tag.test.slug]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tagName, tagSlug, description)
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

func testAccCircuitTerminationResourceConfig_withTags(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, tagName, tagSlug string, portSpeed, upstreamSpeed int, xconnectID, ppInfo string, markConnected bool) string {
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

resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_termination" "test" {
  circuit        = netbox_circuit.test.id
  term_side      = "A"
  site           = netbox_site.test.id
  port_speed     = %d
  upstream_speed = %d
  xconnect_id    = %q
  pp_info        = %q
  mark_connected = %t
  description    = %q
	tags           = [netbox_tag.test.slug]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tagName, tagSlug, portSpeed, upstreamSpeed, xconnectID, ppInfo, markConnected, description)
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
			{
				Config:   testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

// NOTE: Custom field tests for circuit_termination resource are in resources_acceptance_tests_customfields package

func TestAccCircuitTerminationResource_full(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	circuitTypeName := testutil.RandomName("tf-test-ct-full")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-full")
	circuitCID := testutil.RandomName("tf-test-circuit-full")
	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	description := "Full test with all optional fields"
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "upstream_speed", "5000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "xconnect_id", "XCON-FULL-123"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "pp_info", "Patch Panel 1, Port 24"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "true"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_circuit_termination.test", "tags.*", tagSlug),
				),
			},
		},
	})
}

// NOTE: Custom field tests for circuit_termination resource are in resources_acceptance_tests_customfields package

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
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
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCircuitTerminationResource_removeDescription(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-prov-desc")
	providerSlug := testutil.RandomSlug("tf-test-prov-desc")
	circuitTypeName := testutil.RandomName("tf-test-ct-desc")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-desc")
	circuitCID := testutil.RandomName("tf-test-circ-desc")
	siteName := testutil.RandomName("tf-test-site-desc")
	siteSlug := testutil.RandomSlug("tf-test-site-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_withDescription(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, "Description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", "Description"),
				),
			},
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_circuit_termination.test", "description"),
				),
			},
		},
	})
}

func testAccCircuitTerminationResourceConfig_withDescription(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description string) string {
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
  description = %q
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description)
}

func TestAccCircuitTerminationResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-opt")
	providerSlug := testutil.RandomSlug("tf-test-provider-opt")
	circuitTypeName := testutil.RandomName("tf-test-ct-opt")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-opt")
	circuitCID := testutil.RandomName("tf-test-circuit-opt")
	siteName := testutil.RandomName("tf-test-site-opt")
	siteSlug := testutil.RandomSlug("tf-test-site-opt")
	site2Name := testutil.RandomName("tf-test-site2-opt")
	site2Slug := testutil.RandomSlug("tf-test-site2-opt")
	providerNetworkName := testutil.RandomName("tf-test-pn-opt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterSiteCleanup(site2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit" "test" {
  cid              = %[5]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = %[6]q
  slug   = %[7]q
  status = "active"
}

resource "netbox_site" "test2" {
  name   = %[8]q
  slug   = %[9]q
  status = "active"
}

resource "netbox_provider_network" "test" {
  name             = %[10]q
  circuit_provider = netbox_provider.test.id
}

resource "netbox_circuit_termination" "test" {
  circuit          = netbox_circuit.test.id
  term_side        = "A"
  provider_network = netbox_provider_network.test.name
  port_speed       = 1000000
  upstream_speed   = 512000
  xconnect_id      = "XCON-123"
  pp_info          = "PP1-Port5"
  mark_connected   = true
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, site2Name, site2Slug, providerNetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "provider_network"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "1000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "upstream_speed", "512000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "xconnect_id", "XCON-123"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "pp_info", "PP1-Port5"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit" "test" {
  cid              = %[5]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = %[6]q
  slug   = %[7]q
  status = "active"
}

resource "netbox_site" "test2" {
  name   = %[8]q
  slug   = %[9]q
  status = "active"
}

resource "netbox_provider_network" "test" {
  name             = %[10]q
  circuit_provider = netbox_provider.test.id
}

resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.id
  term_side = "A"
  site      = netbox_site.test.id
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, site2Name, site2Slug, providerNetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// site must remain since circuit_termination requires either site or provider_network
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "site"),
					// provider_network should be removed
					resource.TestCheckNoResourceAttr("netbox_circuit_termination.test", "provider_network"),
					// All other optional fields should be removed
					resource.TestCheckNoResourceAttr("netbox_circuit_termination.test", "port_speed"),
					resource.TestCheckNoResourceAttr("netbox_circuit_termination.test", "upstream_speed"),
					resource.TestCheckNoResourceAttr("netbox_circuit_termination.test", "xconnect_id"),
					resource.TestCheckNoResourceAttr("netbox_circuit_termination.test", "pp_info"),
					// mark_connected is Computed with default false, so check for default
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "false"),
				),
			},
		},
	})
}

func TestAccCircuitTerminationResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_circuit_termination",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_circuit": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_termination" "test" {
  # circuit missing
  term_side = "A"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_term_side": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_circuit_type" "test" {
  name = "Test Type"
  slug = "test-type"
}

resource "netbox_circuit" "test" {
  cid = "TEST-001"
  circuit_provider = netbox_provider.test.id
  type = netbox_circuit_type.test.id
}

resource "netbox_circuit_termination" "test" {
  circuit = netbox_circuit.test.id
  # term_side missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

// =============================================================================
// STANDARDIZED TAG TESTS (using helpers)
// =============================================================================

// TestAccCircuitTerminationResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccCircuitTerminationResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-prov-tag")
	providerSlug := testutil.RandomSlug("tf-prov-tag")
	circuitTypeName := testutil.RandomName("tf-ct-tag")
	circuitTypeSlug := testutil.RandomSlug("tf-ct-tag")
	circuitCID := testutil.RandomName("tf-ckt-tag")
	siteName := testutil.RandomName("tf-site-tag")
	siteSlug := testutil.RandomSlug("tf-site-tag")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_circuit_termination",
		ConfigWithoutTags: func() string {
			return testAccCircuitTerminationResourceConfig_tagLifecycle(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, "", "", "", "", "", "", "")
		},
		ConfigWithTags: func() string {
			return testAccCircuitTerminationResourceConfig_tagLifecycle(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1,tag2", "", "")
		},
		ConfigWithDifferentTags: func() string {
			return testAccCircuitTerminationResourceConfig_tagLifecycle(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag3", tag3Name, tag3Slug)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
	})
}

func testAccCircuitTerminationResourceConfig_tagLifecycle(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagSet, tag3Name, tag3Slug string) string {
	tagResources := ""
	tagsList := ""

	if tag1Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}
`, tag1Name, tag1Slug)
	}
	if tag2Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}
`, tag2Name, tag2Slug)
	}
	if tag3Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag3" {
  name = %q
  slug = %q
}
`, tag3Name, tag3Slug)
	}

	switch tagSet {
	case caseTag1Tag2:
		tagsList = tagsDoubleSlug
	case caseTag3:
		tagsList = tagsSingleSlug
	default:
		tagsList = tagsEmpty
	}

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
  name = %q
  slug = %q
}
%s
resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.id
  term_side = "A"
  site      = netbox_site.test.id
  %s
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tagResources, tagsList)
}

// TestAccCircuitTerminationResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccCircuitTerminationResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-prov-order")
	providerSlug := testutil.RandomSlug("tf-prov-order")
	circuitTypeName := testutil.RandomName("tf-ct-order")
	circuitTypeSlug := testutil.RandomSlug("tf-ct-order")
	circuitCID := testutil.RandomName("tf-ckt-order")
	siteName := testutil.RandomName("tf-site-order")
	siteSlug := testutil.RandomSlug("tf-site-order")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_circuit_termination",
		ConfigWithTagsOrderA: func() string {
			return testAccCircuitTerminationResourceConfig_tagOrder(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccCircuitTerminationResourceConfig_tagOrder(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
	})
}

func testAccCircuitTerminationResourceConfig_tagOrder(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	tagsOrder := tagsDoubleSlug
	if !tag1First {
		tagsOrder = tagsDoubleSlugReversed
	}

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
  name = %q
  slug = %q
}

resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.id
  term_side = "A"
  site      = netbox_site.test.id
  %s
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagsOrder)
}
