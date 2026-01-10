//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualDeviceContextDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_vdc_ds_cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.RandomSlug("tf-test-site-ds-cf")
	mfgName := testutil.RandomName("tf-test-mfg-ds-cf")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-ds-cf")
	typeName := testutil.RandomName("tf-test-type-ds-cf")
	typeSlug := testutil.RandomSlug("tf-test-type-ds-cf")
	roleName := testutil.RandomName("tf-test-role-ds-cf")
	roleSlug := testutil.RandomSlug("tf-test-role-ds-cf")
	deviceName := testutil.RandomName("tf-test-device-ds-cf")
	vdcName := testutil.RandomName("tf-test-vdc-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextDataSourceConfig_customFields(customFieldName, siteName, siteSlug, mfgName, mfgSlug, typeName, typeSlug, roleName, roleSlug, deviceName, vdcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccVirtualDeviceContextDataSourceConfig_customFields(customFieldName, siteName, siteSlug, mfgName, mfgSlug, typeName, typeSlug, roleName, roleSlug, deviceName, vdcName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.virtualdevicecontext"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.name
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

resource "netbox_virtual_device_context" "test" {
  name   = %q
  device = netbox_device.test.name
  status = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_virtual_device_context" "test" {
  id = netbox_virtual_device_context.test.id

  depends_on = [netbox_virtual_device_context.test]
}
`, customFieldName, siteName, siteSlug, mfgName, mfgSlug, typeName, typeSlug, roleName, roleSlug, deviceName, vdcName)
}
