package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitResource_basic(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("tf-test-circuit")
	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	typeName := testutil.RandomName("tf-test-circuit-type")
	typeSlug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceConfig_basic(cid, providerName, providerSlug, typeName, typeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttr("netbox_circuit.test", "circuit_provider", providerSlug),
					resource.TestCheckResourceAttr("netbox_circuit.test", "type", typeSlug),
				),
			},
		},
	})
}

func TestAccCircuitResource_full(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("tf-test-circuit-full")
	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	typeName := testutil.RandomName("tf-test-circuit-type-full")
	typeSlug := testutil.RandomSlug("tf-test-circuit-type-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceConfig_full(cid, providerName, providerSlug, typeName, typeSlug, testutil.Description1, testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttr("netbox_circuit.test", "circuit_provider", providerSlug),
					resource.TestCheckResourceAttr("netbox_circuit.test", "type", typeSlug),
					resource.TestCheckResourceAttr("netbox_circuit.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "description", testutil.Description1),
					resource.TestCheckResourceAttr("netbox_circuit.test", "comments", testutil.Comments),
					resource.TestCheckResourceAttr("netbox_circuit.test", "commit_rate", "10000"),
				),
			},
		},
	})
}

func TestAccCircuitResource_update(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("tf-test-circuit-update")
	providerName := testutil.RandomName("tf-test-provider-update")
	providerSlug := testutil.RandomSlug("tf-test-provider-update")
	typeName := testutil.RandomName("tf-test-circuit-type-update")
	typeSlug := testutil.RandomSlug("tf-test-circuit-type-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceConfig_basic(cid, providerName, providerSlug, typeName, typeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
				),
			},
			{
				Config: testAccCircuitResourceConfig_withDescription(cid, providerName, providerSlug, typeName, typeSlug, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttr("netbox_circuit.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccCircuitResource_IDPreservation(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("ckt-id")
	providerName := testutil.RandomName("prov-id")
	providerSlug := testutil.RandomSlug("prov-id")
	typeName := testutil.RandomName("type-id")
	typeSlug := testutil.RandomSlug("type-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceConfig_basic(cid, providerName, providerSlug, typeName, typeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttr("netbox_circuit.test", "circuit_provider", providerSlug),
					resource.TestCheckResourceAttr("netbox_circuit.test", "type", typeSlug),
				),
			},
		},
	})
}

func TestAccCircuitResource_import(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("tf-test-circuit")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	providerName := providerSlug
	typeSlug := testutil.RandomSlug("tf-test-circuit-type")
	typeName := typeSlug

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceConfig_basic(cid, providerName, providerSlug, typeName, typeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
				),
			},
			{
				ResourceName:            "netbox_circuit.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit_provider", "type"},
			},
		},
	})
}

// NOTE: Custom field tests for circuit resource are in resources_acceptance_tests_customfields package

func TestAccConsistency_Circuit(t *testing.T) {
	t.Parallel()
	cid := testutil.RandomName("cid")
	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	typeName := testutil.RandomName("type")
	typeSlug := testutil.RandomSlug("type")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitConsistencyConfig(cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "circuit_provider"),
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "type"),
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccCircuitConsistencyConfig(cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccConsistency_Circuit_LiteralNames(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("cid")
	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	typeName := testutil.RandomName("type")
	typeSlug := testutil.RandomSlug("type")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitConsistencyLiteralNamesConfig(cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttr("netbox_circuit.test", "circuit_provider", providerSlug),
					resource.TestCheckResourceAttr("netbox_circuit.test", "type", typeSlug),
					resource.TestCheckResourceAttr("netbox_circuit.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccCircuitConsistencyLiteralNamesConfig(cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccCircuitResourceConfig_basic(cid, providerName, providerSlug, typeName, typeSlug string) string {
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
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}
`, providerName, providerSlug, typeName, typeSlug, cid)
}

func testAccCircuitResourceConfig_full(cid, providerName, providerSlug, typeName, typeSlug, description, comments string) string {
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
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
  status           = "active"
  description      = %q
  comments         = %q
  commit_rate      = 10000
}
`, providerName, providerSlug, typeName, typeSlug, cid, description, comments)
}

func testAccCircuitResourceConfig_withDescription(cid, providerName, providerSlug, typeName, typeSlug, description string) string {
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
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
  description      = %q
}

`, providerName, providerSlug, typeName, typeSlug, cid, description)

}

func testAccCircuitConsistencyConfig(cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_circuit_type" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_tenant" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_circuit" "test" {
  cid = "%[1]s"
  circuit_provider = netbox_provider.test.slug
  type = netbox_circuit_type.test.slug
  tenant = netbox_tenant.test.name
}
`, cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug)
}

func testAccCircuitConsistencyLiteralNamesConfig(cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_circuit_type" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_tenant" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_circuit" "test" {
  cid = "%[1]s"
  # Use literal string names to mimic existing user state
  circuit_provider = "%[3]s"
  type = "%[5]s"
  tenant = "%[6]s"
  depends_on = [netbox_provider.test, netbox_circuit_type.test, netbox_tenant.test]
}
`, cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug)
}

func TestAccCircuitResource_externalDeletion(t *testing.T) {
	t.Parallel()

	cid := testutil.RandomName("tf-test-circuit-ext-del")
	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("provider")
	typeName := testutil.RandomName("tf-test-circuit-type")
	typeSlug := testutil.RandomSlug("circuit-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceConfig_basic(cid, providerName, providerSlug, typeName, typeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List circuits filtered by CID
					items, _, err := client.CircuitsAPI.CircuitsCircuitsList(context.Background()).CidIc([]string{cid}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit: %v", err)
					}
					t.Logf("Successfully externally deleted circuit with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
