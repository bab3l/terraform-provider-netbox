//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_service_ds_cf")
	serviceName := testutil.RandomName("tf-test-service-ds-cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	deviceName := testutil.RandomName("tf-test-device-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig_customFields(customFieldName, serviceName, siteName, deviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("data.netbox_service.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_service.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_service.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_service.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig_customFields(customFieldName, serviceName, siteName, deviceName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  object_types = ["ipam.service"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_device_role" "test" {
  name  = "test-role-%[3]s"
  slug  = "test-role-%[3]s"
  color = "ff0000"
}

resource "netbox_manufacturer" "test" {
  name = "test-manufacturer-%[3]s"
  slug = "test-manufacturer-%[3]s"
}

resource "netbox_device_type" "test" {
  model        = "test-model-%[3]s"
  slug         = "test-model-%[3]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %[4]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  device   = netbox_device.test.id
  name     = %[2]q
  protocol = "tcp"
  ports    = [80]

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_service" "test" {
  device = netbox_device.test.id
  name   = %[2]q

  depends_on = [netbox_service.test]
}
`, customFieldName, serviceName, siteName, deviceName)
}
