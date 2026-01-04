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

func TestAccCircuitGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group")
	slug := testutil.RandomSlug("tf-test-cg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cg-id")
	slug := testutil.RandomSlug("cg-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-full")
	slug := testutil.RandomSlug("tf-test-cg-full")
	description := testutil.Description1

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-upd")
	slug := testutil.RandomSlug("tf-test-cg-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
				),
			},
			{
				Config: testAccCircuitGroupResourceConfig_full(name, slug, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-imp")
	slug := testutil.RandomSlug("tf-test-cg-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
			},
			{
				ResourceName:      "netbox_circuit_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_CircuitGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-lit")
	slug := testutil.RandomSlug("tf-test-cg-lit")
	description := testutil.Description1

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", description),
				),
			},
			{
				Config:   testAccCircuitGroupConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
				),
			},
		},
	})
}

func testAccCircuitGroupConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}
`, name, slug, description)
}

func testAccCircuitGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccCircuitGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}
`, name, slug, description)
}

func TestAccCircuitGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-circuit-group-ext-del")
	slug := testutil.RandomSlug("circuit-group-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %q
  slug = %q
}
`, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List circuit groups filtered by slug
					items, _, err := client.CircuitsAPI.CircuitsCircuitGroupsList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit group: %v", err)
					}
					t.Logf("Successfully externally deleted circuit group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
