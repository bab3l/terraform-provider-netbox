//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceTypeDataSource_customFields(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg-ds-cf")
	mfgSlug := testutil.GenerateSlug(mfgName)
	typeName := testutil.RandomName("tf-test-type-ds-cf")
	typeSlug := testutil.GenerateSlug(typeName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_devicetype_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(typeSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeDataSourceConfig_withCustomFields(mfgName, mfgSlug, typeName, typeSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_type.test", "model", typeName),
					resource.TestCheckResourceAttr("data.netbox_device_type.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_device_type.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_device_type.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_device_type.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccDeviceTypeDataSourceConfig_withCustomFields(mfgName, mfgSlug, typeName, typeSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.devicetype"]
  type         = "text"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_device_type" "test" {
  model = netbox_device_type.test.model

  depends_on = [netbox_device_type.test]
}
`, customFieldName, mfgName, mfgSlug, typeName, typeSlug)
}
