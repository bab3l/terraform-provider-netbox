package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRearPortTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateDataSourceConfig("test-rear-port-template"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_port_template.test", "name", "test-rear-port-template"),
					resource.TestCheckResourceAttr("data.netbox_rear_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttrSet("data.netbox_rear_port_template.test", "device_type"),
				),
			},
		},
	})
}

func testAccRearPortTemplateDataSourceConfig(name string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "8p8c"
}

data "netbox_rear_port_template" "test" {
  id = netbox_rear_port_template.test.id
}
`, name)
}
