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

func TestAccRegionResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-region")

	slug := testutil.RandomSlug("tf-test-region")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),

					resource.TestCheckResourceAttr("netbox_region.test", "name", name),

					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccRegionResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-region-full")

	slug := testutil.RandomSlug("tf-test-region-full")

	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),

					resource.TestCheckResourceAttr("netbox_region.test", "name", name),

					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_region.test", "description", description),
				),
			},
		},
	})

}

func TestAccRegionResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-region-update")

	slug := testutil.RandomSlug("tf-test-region-upd")

	updatedName := testutil.RandomName("tf-test-region-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),

					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
				),
			},

			{

				Config: testAccRegionResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),

					resource.TestCheckResourceAttr("netbox_region.test", "name", updatedName),
				),
			},
		},
	})

}

func TestAccRegionResource_withParent(t *testing.T) {

	t.Parallel()

	parentName := testutil.RandomName("tf-test-region-parent")

	parentSlug := testutil.RandomSlug("tf-test-region-prnt")

	childName := testutil.RandomName("tf-test-region-child")

	childSlug := testutil.RandomSlug("tf-test-region-chld")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(childSlug)

	cleanup.RegisterRegionCleanup(parentSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.parent", "id"),

					resource.TestCheckResourceAttr("netbox_region.parent", "name", parentName),

					resource.TestCheckResourceAttrSet("netbox_region.child", "id"),

					resource.TestCheckResourceAttr("netbox_region.child", "name", childName),

					resource.TestCheckResourceAttrPair("netbox_region.child", "parent", "netbox_region.parent", "id"),
				),
			},
		},
	})

}

func TestAccRegionResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-region-import")

	slug := testutil.RandomSlug("tf-test-region-imp")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionResourceConfig_import(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_region.test", "name", name),

					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_region.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccRegionResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_region" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccRegionResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_region" "test" {

  name        = %q

  slug        = %q

  description = %q

}

`, name, slug, description)

}

func testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_region" "parent" {

  name = %q

  slug = %q

}

resource "netbox_region" "child" {

  name   = %q

  slug   = %q

  parent = netbox_region.parent.id

}

`, parentName, parentSlug, childName, childSlug)

}

func TestAccConsistency_Region_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-region-lit")

	slug := testutil.RandomSlug("tf-test-region-lit")

	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionConsistencyLiteralNamesConfig(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),

					resource.TestCheckResourceAttr("netbox_region.test", "name", name),

					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_region.test", "description", description),
				),
			},

			{

				Config: testAccRegionConsistencyLiteralNamesConfig(name, slug, description),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
				),
			},
		},
	})

}

func testAccRegionConsistencyLiteralNamesConfig(name, slug, description string) string {

	return fmt.Sprintf(`

resource "netbox_region" "test" {

  name        = %q

  slug        = %q

  description = %q

}

`, name, slug, description)

}

func testAccRegionResourceConfig_import(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_region" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}
