//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesDataSource_queryWithCustomFields(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-ifaces-q-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	roleName := testutil.RandomName("tf-test-role-ifaces-q-cf")
	roleSlug := testutil.GenerateSlug(roleName)
	mfgName := testutil.RandomName("tf-test-mfg-ifaces-q-cf")
	mfgSlug := testutil.GenerateSlug(mfgName)
	typeName := testutil.RandomName("tf-test-type-ifaces-q-cf")
	typeSlug := testutil.GenerateSlug(typeName)
	deviceName := testutil.RandomName("tf-test-device-ifaces-q-cf")
	ifaceName := testutil.RandomName("tf-test-iface-ifaces-q-cf")
	customFieldName := testutil.RandomCustomFieldName("tf_test_ifaces_q_cf")
	customFieldValue := "datasource-test-value"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(typeSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(ifaceName, deviceName)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfacesDataSourceConfig_withCustomFields(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName, ifaceName, customFieldName, customFieldValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.0", ifaceName),
					resource.TestCheckResourceAttrPair("data.netbox_interfaces.test", "ids.0", "netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "interfaces.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_interfaces.test", "interfaces.0.id", "netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "interfaces.0.name", ifaceName),
				),
			},
		},
	})
}

func testAccInterfacesDataSourceConfig_withCustomFields(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName, ifaceName, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[11]q
  object_types = ["dcim.interface"]
  type         = "text"
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  model        = %[7]q
  slug         = %[8]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name    = %[3]q
  slug    = %[4]q
  color   = "ff0000"
  vm_role = false
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name
}

resource "netbox_interface" "test" {
  name   = %[10]q
  device = netbox_device.test.name
  type   = "1000base-t"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[12]q
    }
  ]
}

data "netbox_interfaces" "test" {
  filter {
    name   = "custom_field_value"
    values = ["${netbox_custom_field.test.name}=%[12]s"]
  }

  depends_on = [netbox_interface.test]
}
`, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName, ifaceName, customFieldName, customFieldValue)
}
