package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsolePortTemplateDataSource_basic(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.GenerateSlug(mfgName)

	deviceTypeModel := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	portTemplateName := testutil.RandomName("console-port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateDataSourceConfig(mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, portTemplateName),
				Check: resource.ComposeTestCheckFunc(
					// Check by_id lookup
					resource.TestCheckResourceAttr("data.netbox_console_port_template.by_id", "name", portTemplateName),
					resource.TestCheckResourceAttr("data.netbox_console_port_template.by_id", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port_template.by_id", "device_type"),
					// Check by_device_type_and_name lookup
					resource.TestCheckResourceAttr("data.netbox_console_port_template.by_device_type_and_name", "name", portTemplateName),
					resource.TestCheckResourceAttr("data.netbox_console_port_template.by_device_type_and_name", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port_template.by_device_type_and_name", "device_type"),
					// Verify both lookups return same console port template
					resource.TestCheckResourceAttrPair("data.netbox_console_port_template.by_id", "id", "data.netbox_console_port_template.by_device_type_and_name", "id"),
				),
			},
		},
	})
}

func testAccConsolePortTemplateDataSourceConfig(mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, portTemplateName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_console_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %[5]q
  type        = "de-9"
}

data "netbox_console_port_template" "by_id" {
  id = netbox_console_port_template.test.id
}

data "netbox_console_port_template" "by_device_type_and_name" {
  device_type = netbox_device_type.test.id
  name        = netbox_console_port_template.test.name
}
`, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, portTemplateName)
}
