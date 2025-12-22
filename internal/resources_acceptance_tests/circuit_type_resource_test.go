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

func TestAccCircuitTypeResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type")

	slug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccCircuitTypeResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-full")

	slug := testutil.RandomSlug("tf-test-circuit-type-full")

	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_full(name, slug, description, testutil.Color),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", description),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "color", testutil.Color),
				),
			},
		},
	})

}

func TestAccCircuitTypeResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-update")

	slug := testutil.RandomSlug("tf-test-circuit-type-update")

	updatedName := testutil.RandomName("tf-test-circuit-type-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
				),
			},

			{

				Config: testAccCircuitTypeResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", updatedName),
				),
			},
		},
	})

}

func TestAccCircuitTypeResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type")

	slug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_circuit_type.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConsistency_CircuitType_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-circuit-type-lit")
	slug := testutil.RandomSlug("tf-test-circuit-type-lit")
	description := testutil.RandomName("description")
	color := "2196f3"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeConsistencyLiteralNamesConfig(name, slug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "color", color),
				),
			},
			{
				Config:   testAccCircuitTypeConsistencyLiteralNamesConfig(name, slug, description, color),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
				),
			},
		},
	})
}

func testAccCircuitTypeResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccCircuitTypeResourceConfig_full(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}
`, name, slug, description, color)
}

func testAccCircuitTypeConsistencyLiteralNamesConfig(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}
`, name, slug, description, color)
}
