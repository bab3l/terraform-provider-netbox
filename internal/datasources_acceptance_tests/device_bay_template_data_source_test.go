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

func TestAccDeviceBayTemplateDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-ds-id")
	manufacturerName := testutil.RandomName("tf-test-mfr-id")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-id")
	deviceTypeName := testutil.RandomName("tf-test-dt-id")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceBayTemplateCleanup(name)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceBayTemplateDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_bay_template.test", "name", name),
				),
			},
		},
	})
}

func TestAccDeviceBayTemplateDataSource_basic(t *testing.T) {

	t.Parallel()

	// Generate unique names

	name := testutil.RandomName("tf-test-dbt-ds")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})

}

func testAccDeviceBayTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

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

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}

resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

}

data "netbox_device_bay_template" "test" {

  id = netbox_device_bay_template.test.id

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}
func TestAccDeviceBayTemplateDataSource_byNameAndDeviceType(t *testing.T) {

	t.Parallel()

	// Generate unique names

	name := testutil.RandomName("tf-test-dbt-ds")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateDataSourceConfigByNameAndDeviceType(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})

}

func testAccDeviceBayTemplateDataSourceConfigByNameAndDeviceType(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

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

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}

resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

}

data "netbox_device_bay_template" "test" {

  name        = netbox_device_bay_template.test.name

  device_type = netbox_device_type.test.id

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}
