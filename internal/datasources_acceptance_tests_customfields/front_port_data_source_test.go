//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFrontPortDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_frontport_ds_cf")
	portName := testutil.RandomName("tf-test-frontport-ds-cf")
	rearPortName := testutil.RandomName("tf-test-rearport-ds-cf")
	deviceName := testutil.RandomName("tf-test-device-ds-cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortDataSourceConfig_customFields(customFieldName, portName, rearPortName, deviceName, siteName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "name", portName),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccFrontPortDataSourceConfig_customFields(customFieldName, portName, rearPortName, deviceName, siteName string) string {
	manufacturerSlug := testutil.RandomName("test-mfg")
	roleSlug := testutil.RandomName("test-role")
	modelSlug := testutil.RandomName("test-model")

	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.frontport"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  name      = %q
  device    = netbox_device.test.id
  type      = "8p8c"
  positions = 1
}

resource "netbox_front_port" "test" {
  name   = %q
  device = netbox_device.test.id
  type   = "8p8c"
  rear_port = netbox_rear_port.test.id
  rear_port_position = 1

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_front_port" "test" {
  device_id = netbox_device.test.id
  name      = %q

  depends_on = [netbox_front_port.test]
}
`, customFieldName, siteName, siteName, roleSlug, roleSlug, modelSlug, modelSlug, manufacturerSlug, manufacturerSlug, deviceName, rearPortName, portName, portName)
}
