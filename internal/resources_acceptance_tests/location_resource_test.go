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

func TestAccLocationResource_basic(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site")

	siteSlug := testutil.RandomSlug("tf-test-loc-site")

	name := testutil.RandomName("tf-test-location")

	slug := testutil.RandomSlug("tf-test-location")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterLocationCleanup(slug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),

					resource.TestCheckResourceAttr("netbox_location.test", "name", name),

					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),

					resource.TestCheckResourceAttrPair("netbox_location.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})

}

func TestAccLocationResource_full(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-full")

	siteSlug := testutil.RandomSlug("tf-test-loc-s-full")

	name := testutil.RandomName("tf-test-location-full")

	slug := testutil.RandomSlug("tf-test-loc-full")

	description := testutil.RandomName("description")

	facility := "Building A"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterLocationCleanup(slug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccLocationResourceConfig_full(siteName, siteSlug, name, slug, description, facility),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),

					resource.TestCheckResourceAttr("netbox_location.test", "name", name),

					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_location.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_location.test", "description", description),

					resource.TestCheckResourceAttr("netbox_location.test", "facility", facility),
				),
			},
		},
	})

}

func TestAccLocationResource_import(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site")

	siteSlug := testutil.RandomSlug("tf-test-loc-site")

	name := testutil.RandomName("tf-test-location")

	slug := testutil.RandomSlug("tf-test-location")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterLocationCleanup(slug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccLocationResourceConfig_import(siteName, siteSlug, name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),

					resource.TestCheckResourceAttr("netbox_location.test", "name", name),

					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),

					resource.TestCheckResourceAttrPair("netbox_location.test", "site", "netbox_site.test", "id"),
				),
			},

			{

				ResourceName: "netbox_location.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"site"},
			},
		},
	})

}

func TestAccLocationResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site")
	siteSlug := testutil.RandomSlug("tf-test-loc-site")
	name := testutil.RandomName("tf-test-location")
	slug := testutil.RandomSlug("tf-test-location")
	updatedName := testutil.RandomName("tf-test-location-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
				),
			},
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccConsistency_Location_LiteralNames(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("tf-test-loc-site-lit")
	siteSlug := testutil.RandomSlug("tf-test-loc-site-lit")
	name := testutil.RandomName("tf-test-location-lit")
	slug := testutil.RandomSlug("tf-test-location-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationConsistencyLiteralNamesConfig(siteName, siteSlug, name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_location.test", "description", description),
				),
			},
			{
				Config:   testAccLocationConsistencyLiteralNamesConfig(siteName, siteSlug, name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
				),
			},
		},
	})
}

func TestAccLocationResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-id")
	siteSlug := testutil.RandomSlug("tf-test-site-id")
	locationName := testutil.RandomName("tf-test-location-id")
	locationSlug := testutil.RandomSlug("tf-test-loc-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterLocationCleanup(locationSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, locationName, locationSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", locationName),
					resource.TestCheckResourceAttrSet("netbox_location.test", "site"),
				),
			},
		},
	})
}

func testAccLocationConsistencyLiteralNamesConfig(siteName, siteSlug, name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_location" "test" {
  name        = %q
  slug        = %q
  site        = netbox_site.test.id
  description = %q
}
`, siteName, siteSlug, name, slug, description)
}

func testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source  = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_site" "test" {

  name   = %q

  slug   = %q

  status = "active"

}

resource "netbox_location" "test" {

  name = %q

  slug = %q

  site = netbox_site.test.id

}

`, siteName, siteSlug, name, slug)

}

func testAccLocationResourceConfig_full(siteName, siteSlug, name, slug, description, facility string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source  = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_site" "test" {

  name   = %q

  slug   = %q

  status = "active"

}

resource "netbox_location" "test" {

  name        = %q

  slug        = %q

  site        = netbox_site.test.id

  status      = "active"

  description = %q

  facility    = %q

}

`, siteName, siteSlug, name, slug, description, facility)

}

func testAccLocationResourceConfig_import(siteName, siteSlug, name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_location" "test" {
  name = %q
  slug = %q
  site = netbox_site.test.id
}
`, siteName, siteSlug, name, slug)
}

func TestAccLocationResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-del")
	siteSlug := testutil.RandomSlug("tf-test-loc-site-del")
	name := testutil.RandomName("tf-test-location-del")
	slug := testutil.RandomSlug("tf-test-location-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimLocationsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find location for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimLocationsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete location: %v", err)
					}
					t.Logf("Successfully externally deleted location with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
