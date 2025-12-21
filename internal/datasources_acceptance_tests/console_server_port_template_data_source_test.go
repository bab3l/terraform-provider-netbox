package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsoleServerPortTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsoleServerPortTemplateDataSourceConfig("console-server-port-0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.test", "name", "console-server-port-0"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.test", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port_template.test", "device_type"),
				),
			},
		},
	})
}

func testAccConsoleServerPortTemplateDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "test-device-type"
}

resource "netbox_console_server_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "de-9"
}

data "netbox_console_server_port_template" "test" {
  id = netbox_console_server_port_template.test.id
}
`, name)
}
