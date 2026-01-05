package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManufacturerResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-manufacturer")
	slug := testutil.RandomSlug("tf-test-mfr")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccManufacturerResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-manufacturer-full")
	slug := testutil.RandomSlug("tf-test-mfr-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "description", description),
				),
			},
		},
	})
}

func TestAccManufacturerResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-manufacturer-update")
	slug := testutil.RandomSlug("tf-test-mfr-upd")
	updatedName := testutil.RandomName("tf-test-manufacturer-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
				),
			},
			{
				Config: testAccManufacturerResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccManufacturerResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-manufacturer-import")
	slug := testutil.RandomSlug("tf-test-mfr-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
				),
			},
			{
				ResourceName:            "netbox_manufacturer.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name"}, // display_name is computed and may differ after name changes
			},
		},
	})
}

func TestAccConsistency_Manufacturer_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-manufacturer-lit")
	slug := testutil.RandomSlug("tf-test-mfr-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "description", description),
				),
			},
			{
				Config:   testAccManufacturerConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
				),
			},
		},
	})
}

func testAccManufacturerConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccManufacturerResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-manufacturer-id")
	slug := testutil.RandomSlug("tf-test-manufacturer-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
				),
			},
		},
	})
}

func testAccManufacturerResourceConfig_basic(name, slug string) string {
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

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccManufacturerResourceConfig_full(name, slug, description string) string {
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

resource "netbox_manufacturer" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func testAccManufacturerResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccManufacturerResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-manufacturer-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimManufacturersList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find manufacturer for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimManufacturersDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete manufacturer: %v", err)
					}
					t.Logf("Successfully externally deleted manufacturer with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
