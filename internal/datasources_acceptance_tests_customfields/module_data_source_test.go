//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleDataSource_customFields(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	roleName := testutil.RandomName("tf-test-role-ds-cf")
	roleSlug := testutil.GenerateSlug(roleName)
	mfgName := testutil.RandomName("tf-test-mfg-ds-cf")
	mfgSlug := testutil.GenerateSlug(mfgName)
	typeName := testutil.RandomName("tf-test-type-ds-cf")
	typeSlug := testutil.GenerateSlug(typeName)
	deviceName := testutil.RandomName("tf-test-device-ds-cf")
	moduleTypeName := testutil.RandomName("tf-test-modtype-ds-cf")
	moduleTypeSlug := testutil.GenerateSlug(moduleTypeName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_module_ds_cf")

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
				Config: testAccModuleDataSourceConfig_withCustomFields(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName, moduleTypeName, moduleTypeSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_module.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_module.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_module.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_module.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_module.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccModuleDataSourceConfig_withCustomFields(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, typeName, typeSlug, deviceName, moduleTypeName, moduleTypeSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.module"]
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
}

resource "netbox_module_type" "test" {
  model        = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.name
  name   = "Bay 1"
}

resource "netbox_module" "test" {
  device      = netbox_device.test.name
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.model
  status      = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_module" "test" {
  id = netbox_module.test.id

  depends_on = [netbox_module.test]
}
`, customFieldName, siteName, siteSlug, mfgName, mfgSlug, typeName, typeSlug, roleName, roleSlug, deviceName, moduleTypeName)
}
