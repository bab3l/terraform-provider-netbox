//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeDataSource_customFields(t *testing.T) {
	octet2 := acctest.RandIntRange(100, 200)
	octet3 := acctest.RandIntRange(100, 200)
	octet4Start := acctest.RandIntRange(1, 50)
	startAddress := fmt.Sprintf("10.%d.%d.%d", octet2, octet3, octet4Start)
	endAddress := fmt.Sprintf("10.%d.%d.%d", octet2, octet3, octet4Start+50)
	customFieldName := testutil.RandomCustomFieldName("tf_test_iprange_ds_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeDataSourceConfigWithCustomFields(startAddress, endAddress, customFieldName, customFieldValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_ip_range.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccIPRangeDataSourceConfigWithCustomFields(startAddress, endAddress, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[3]q
  object_types = ["ipam.iprange"]
  type         = "text"
}

resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  status        = "active"
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[4]q
    }
  ]
}

data "netbox_ip_range" "test" {
  start_address = netbox_ip_range.test.start_address
  end_address   = netbox_ip_range.test.end_address
  depends_on = [netbox_ip_range.test]
}
`, startAddress, endAddress, customFieldName, customFieldValue)
}
