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

func TestAccConsoleServerPortTemplateResource_basic(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := testutil.RandomName("consoleserver")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_server_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_console_server_port_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

func TestAccConsoleServerPortTemplateResource_full(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := testutil.RandomName("consoleserver")

	label := testutil.RandomName("label")

	portType := "rj-45"

	description := testutil.RandomName("description")

	updatedName := testutil.RandomName("consoleserver-upd")

	updatedLabel := testutil.RandomName("label-upd")

	updatedDescription := testutil.RandomName("description-upd")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_server_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "description", description),
				),
			},

			{

				Config: testAccConsoleServerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccConsistency_ConsoleServerPortTemplate_LiteralNames(t *testing.T) {

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

				Config: testAccConsoleServerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", portName),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "device_type", deviceTypeSlug),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccConsoleServerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})

}

func testAccConsoleServerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {

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

resource "netbox_console_server_port_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

func testAccConsoleServerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string) string {

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

resource "netbox_console_server_port_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  label       = %q

  type        = %q

  description = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description)

}

func testAccConsoleServerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {

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

resource "netbox_console_server_port_template" "test" {

  # Use literal string slug to mimic existing user state

  device_type = %q

  name = %q

  type = "rj-45"

  depends_on = [netbox_device_type.test]

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, portName)

}
