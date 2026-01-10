//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleTypeDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_module_type_ds_cf")
	manufacturerSlug := testutil.RandomName("tf-test-mfr-ds-cf")
	moduleModel := testutil.RandomName("tf-test-module-model-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeDataSourceConfig_customFields(customFieldName, manufacturerSlug, moduleModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_type.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_module_type.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_module_type.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_module_type.test", "custom_fields.0.value", "test-module-type-value"),
				),
			},
		},
	})
}

func testAccModuleTypeDataSourceConfig_customFields(customFieldName, manufacturerSlug, moduleModel string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  label        = "Test Custom Field"
  type         = "text"
  object_types = ["dcim.moduletype"]
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.slug
  model        = %[3]q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-module-type-value"
    }
  ]
}

data "netbox_module_type" "test" {
  id = netbox_module_type.test.id

  depends_on = [netbox_module_type.test]
}
`, customFieldName, manufacturerSlug, moduleModel)
}
