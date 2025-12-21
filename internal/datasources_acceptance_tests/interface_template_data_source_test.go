package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccInterfaceTemplateDataSourcePrereqs creates prerequisites for interface template data source tests.

func testAccInterfaceTemplateDataSourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

}

resource "netbox_interface_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  type        = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType)

}

// testAccInterfaceTemplateDataSourceByID looks up an interface template by ID.

func testAccInterfaceTemplateDataSourceByID(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType string) string {

	return testAccInterfaceTemplateDataSourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType) + `

data "netbox_interface_template" "test" {

  id = netbox_interface_template.test.id

}

`

}

// testAccInterfaceTemplateDataSourceByName looks up an interface template by name and device type.

func testAccInterfaceTemplateDataSourceByName(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType string) string {

	return testAccInterfaceTemplateDataSourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType) + fmt.Sprintf(`

data "netbox_interface_template" "test" {

  name        = %q

  device_type = netbox_device_type.test.id

  depends_on = [netbox_interface_template.test]

}

`, templateName)

}

func TestAccInterfaceTemplateDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	manufacturerName := testutil.RandomName("mfr-ds")

	manufacturerSlug := testutil.RandomSlug("mfr-ds")

	deviceTypeName := testutil.RandomName("dt-ds")

	deviceTypeSlug := testutil.RandomSlug("dt-ds")

	templateName := testutil.RandomName("eth")

	templateType := "1000base-t"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceTemplateDataSourceByID(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "name", templateName),

					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "type", templateType),

					resource.TestCheckResourceAttrSet("data.netbox_interface_template.test", "id"),
				),
			},
		},
	})

}

func TestAccInterfaceTemplateDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	manufacturerName := testutil.RandomName("mfr-ds")

	manufacturerSlug := testutil.RandomSlug("mfr-ds")

	deviceTypeName := testutil.RandomName("dt-ds")

	deviceTypeSlug := testutil.RandomSlug("dt-ds")

	templateName := testutil.RandomName("eth")

	templateType := "1000base-t"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceTemplateDataSourceByName(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "name", templateName),

					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "type", templateType),

					resource.TestCheckResourceAttrSet("data.netbox_interface_template.test", "id"),
				),
			},
		},
	})

}
