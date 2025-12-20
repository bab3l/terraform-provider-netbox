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

func TestAccManufacturerResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-manufacturer")

	slug := testutil.RandomSlug("tf-test-mfr")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

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

	name := testutil.RandomName("tf-test-manufacturer-full")

	slug := testutil.RandomSlug("tf-test-mfr-full")

	description := "Test manufacturer with all fields"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

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

	name := testutil.RandomName("tf-test-manufacturer-update")

	slug := testutil.RandomSlug("tf-test-mfr-upd")

	updatedName := testutil.RandomName("tf-test-manufacturer-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

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

	name := testutil.RandomName("tf-test-manufacturer-import")

	slug := testutil.RandomSlug("tf-test-mfr-imp")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccManufacturerResourceConfig_import(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_manufacturer.test",

				ImportState: true,

				ImportStateVerify: true,
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
