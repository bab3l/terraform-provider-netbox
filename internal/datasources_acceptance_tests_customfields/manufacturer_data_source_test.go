//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManufacturerDataSource_customFields(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg-ds-cf")
	mfgSlug := testutil.GenerateSlug(mfgName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_mfg_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerDataSourceConfig_withCustomFields(mfgName, mfgSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", mfgName),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccManufacturerDataSourceConfig_withCustomFields(mfgName, mfgSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.manufacturer"]
  type         = "text"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_manufacturer" "test" {
  name = netbox_manufacturer.test.name

  depends_on = [netbox_manufacturer.test]
}
`, customFieldName, mfgName, mfgSlug)
}
