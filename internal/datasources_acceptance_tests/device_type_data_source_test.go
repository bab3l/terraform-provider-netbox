package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceTypeDataSource_basic(t *testing.T) {

	// Generate unique names

	model := testutil.RandomName("tf-test-devicetype-ds")

	slug := testutil.RandomSlug("tf-test-dt-ds")

	manufacturerName := testutil.RandomName("tf-test-mfr-ds")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceTypeCleanup(slug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeDataSourceConfig(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "manufacturer", manufacturerSlug),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "u_height", "1"),
				),
			},
		},
	})

}

func testAccDeviceTypeDataSourceConfig(model, slug, manufacturerName, manufacturerSlug string) string {

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

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.slug

  model        = %q

  slug         = %q

}

data "netbox_device_type" "test" {

  slug = netbox_device_type.test.slug

}

`, manufacturerName, manufacturerSlug, model, slug)

}

// Route Target Data Source Tests
