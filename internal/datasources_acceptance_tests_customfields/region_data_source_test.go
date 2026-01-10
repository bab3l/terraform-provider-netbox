//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegionDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_region_ds_cf")
	regionSlug := testutil.RandomName("tf-test-region-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfig_customFields(customFieldName, regionSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_region.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_region.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_region.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_region.test", "custom_fields.0.value", "test-region-value"),
				),
			},
		},
	})
}

func testAccRegionDataSourceConfig_customFields(customFieldName, regionSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  label        = "Test Custom Field"
  type         = "text"
  object_types = ["dcim.region"]
}

resource "netbox_region" "test" {
  name = "Test Region"
  slug = %[2]q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-region-value"
    }
  ]
}

data "netbox_region" "test" {
  slug = netbox_region.test.slug

  depends_on = [netbox_region.test]
}
`, customFieldName, regionSlug)
}
