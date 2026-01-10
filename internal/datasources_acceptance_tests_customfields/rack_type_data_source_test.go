//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackTypeDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_rack_type_ds_cf")
	rackTypeSlug := testutil.RandomName("tf-test-rack-type-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeDataSourceConfig_customFields(customFieldName, rackTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "custom_fields.0.value", "test-rack-type-value"),
				),
			},
		},
	})
}

func testAccRackTypeDataSourceConfig_customFields(customFieldName, rackTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  label        = "Test Custom Field"
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer-rt-ds-cf"
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.slug
  slug         = %[2]q
  model        = "Test Rack Type"
  form_factor  = "4-post-cabinet"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-rack-type-value"
    }
  ]
}

data "netbox_rack_type" "test" {
  slug = netbox_rack_type.test.slug

  depends_on = [netbox_rack_type.test]
}
`, customFieldName, rackTypeSlug)
}
