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

func TestAccConsolePortTemplateResource_basic(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	const name = "console0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_console_port_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

func TestAccConsolePortTemplateResource_full(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	const name = "console0"

	label := "Console 0"

	portType := "rj-45"

	description := "Test console port template"

	updatedName := "console1"

	updatedLabel := "Console 1"

	updatedDescription := "Updated console port template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", description),
				),
			},

			{

				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccConsistency_ConsolePortTemplate(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("devicetype")

	deviceTypeSlug := testutil.RandomSlug("devicetype")

	const portName = "console0"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "device_type", deviceTypeName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})

}

func TestAccConsistency_ConsolePortTemplate_LiteralNames(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	portName := testutil.RandomName("port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "device_type", deviceTypeSlug),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})

}

func testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {

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

resource "netbox_console_port_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

func testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string) string {

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

resource "netbox_console_port_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  label       = %q

  type        = %q

  description = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description)

}

func testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}

resource "netbox_device_type" "test" {

  model = "%[3]s"

  slug = "%[4]s"

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_console_port_template" "test" {

  device_type = netbox_device_type.test.model

  name = "%[5]s"

  type = "rj-45"

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)

}

func testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  model        = %q

  slug         = %q

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_console_port_template" "test" {

  # Use literal string slug to mimic existing user state

  device_type = %q

  name = %q

  type = "rj-45"

  depends_on = [netbox_device_type.test]

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, portName)

}
