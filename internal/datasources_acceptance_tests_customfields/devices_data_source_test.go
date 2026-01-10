//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevicesDataSource_queryWithCustomFields(t *testing.T) {
	deviceName := testutil.RandomName("tf-test-device-q-cf")
	siteName := testutil.RandomName("tf-test-site-q-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	roleName := testutil.RandomName("tf-test-role-q-cf")
	roleSlug := testutil.GenerateSlug(roleName)
	mfgName := testutil.RandomName("tf-test-mfg-q-cf")
	mfgSlug := testutil.GenerateSlug(mfgName)
	typeName := testutil.RandomName("tf-test-type-q-cf")
	typeSlug := testutil.GenerateSlug(typeName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_device_q_cf")
	customFieldValue := testutil.RandomName("tf-test-cf-value")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(typeSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesDataSourceConfig_withCustomFields(deviceName, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, customFieldName, customFieldValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_devices.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.0", deviceName),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "devices.0.id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.0.name", deviceName),
				),
			},
		},
	})
}

func testAccDevicesDataSourceConfig_withCustomFields(deviceName, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.device"]
  type         = "text"
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name    = %q
  slug    = %q
  color   = "ff0000"
  vm_role = false
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %q
    }
  ]
}

data "netbox_devices" "test" {
  filter {
		name   = "custom_field_value"
		values = ["${netbox_custom_field.test.name}=%s"]
  }

  depends_on = [netbox_device.test]
}
`, customFieldName, siteName, siteSlug, mfgName, mfgSlug, typeName, typeSlug, roleName, roleSlug, deviceName, customFieldValue, customFieldValue)
}
