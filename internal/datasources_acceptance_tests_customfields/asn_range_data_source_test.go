//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNRangeDataSource_customFields(t *testing.T) {
	rangeName := testutil.RandomName("tf-test-range-ds-cf")
	rangeSlug := testutil.GenerateSlug(rangeName)
	rirName := testutil.RandomName("tf-test-rir-ds-cf")
	rirSlug := testutil.GenerateSlug(rirName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_range_ds_cf")
	customFieldValue := "range-datasource-test"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNRangeCleanup(rangeSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeDataSourceConfig_withCustomFields(rangeName, rangeSlug, rirName, rirSlug, customFieldName, customFieldValue),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccASNRangeDataSourceConfig_withCustomFields(rangeName, rangeSlug, rirName, rirSlug, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["ipam.asnrange"]
  type         = "text"
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
	rir   = netbox_rir.test.id
  start = "64512"
  end   = "64520"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %q
    }
  ]
}

data "netbox_asn_range" "test" {
  slug = netbox_asn_range.test.slug

  depends_on = [netbox_asn_range.test]
}
`, rirName, rirSlug, customFieldName, rangeName, rangeSlug, customFieldValue)
}
