//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCableDataSource_customFields(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	roleName := testutil.RandomName("tf-test-role-ds-cf")
	roleSlug := testutil.GenerateSlug(roleName)
	mfgName := testutil.RandomName("tf-test-mfg-ds-cf")
	mfgSlug := testutil.GenerateSlug(mfgName)
	typeName := testutil.RandomName("tf-test-type-ds-cf")
	typeSlug := testutil.GenerateSlug(typeName)
	deviceName1 := testutil.RandomName("tf-test-device1-ds-cf")
	deviceName2 := testutil.RandomName("tf-test-device2-ds-cf")
	customFieldName := testutil.RandomCustomFieldName("tf_test_cable_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(typeSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableDataSourceConfig_withCustomFields(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName1, deviceName2, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_cable.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccCableDataSourceConfig_withCustomFields(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName1, deviceName2, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.cable"]
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

resource "netbox_device" "test1" {
  name        = %q
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name
}

resource "netbox_device" "test2" {
  name        = %q
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name
}

resource "netbox_interface" "test1" {
  name   = "eth0"
  device = netbox_device.test1.name
  type   = "1000base-t"
}

resource "netbox_interface" "test2" {
  name   = "eth0"
  device = netbox_device.test2.name
  type   = "1000base-t"
}

resource "netbox_cable" "test" {
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test1.id
    }
  ]
  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test2.id
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_cable" "test" {
  id = netbox_cable.test.id

  depends_on = [netbox_cable.test]
}
`, customFieldName, siteName, siteSlug, mfgName, mfgSlug, typeName, typeSlug, roleName, roleSlug, deviceName1, deviceName2)
}
