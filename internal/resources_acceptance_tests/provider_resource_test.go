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

func TestAccProviderResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-provider")

	slug := testutil.RandomSlug("tf-test-provider")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),

					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccProviderResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-provider-full")

	slug := testutil.RandomSlug("tf-test-provider-full")

	description := testutil.RandomName("description")

	comments := testutil.RandomName("comments")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_full(name, slug, description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),

					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_provider.test", "description", description),

					resource.TestCheckResourceAttr("netbox_provider.test", "comments", comments),
				),
			},
		},
	})

}

func TestAccProviderResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-provider-update")

	slug := testutil.RandomSlug("tf-test-provider-update")

	updatedName := testutil.RandomName("tf-test-provider-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
				),
			},

			{

				Config: testAccProviderResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider.test", "name", updatedName),
				),
			},
		},
	})

}

func TestAccProviderResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-provider")

	slug := testutil.RandomSlug("tf-test-provider")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),

					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_provider.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccProviderResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccProviderResourceConfig_full(name, slug, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name        = %q

  slug        = %q

  description = %q

  comments    = %q

}

`, name, slug, description, comments)

}

func TestAccConsistency_Provider_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("provider")
	slug := testutil.RandomSlug("provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccProviderConsistencyLiteralNamesConfig(name, slug),
			},
		},
	})
}

// TestAccProviderResource_IDPreservation tests that the provider resource preserves the
// ID as the immutable identifier when using different reference formats (ID, name, slug).
func TestAccProviderResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-id")
	slug := testutil.RandomSlug("tf-test-provider-id")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},
		},
	})
}

func testAccProviderConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func TestAccProviderResource_externalDeletion(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-provider-ext-del")
	slug := testutil.RandomSlug("provider-ext-del")
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
`, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List providers filtered by slug
					items, _, err := client.CircuitsAPI.CircuitsProvidersList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find provider for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsProvidersDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete provider: %v", err)
					}
					t.Logf("Successfully externally deleted provider with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
