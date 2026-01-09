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

func TestAccVRFDataSource_customFields(t *testing.T) {
	vrfName := testutil.RandomName("tf-test-vrf-ds-cf")
	customFieldName := testutil.RandomCustomFieldName("tf_test_vrf_ds_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(vrfName)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFDataSourceConfigWithCustomFields(vrfName, customFieldName, customFieldValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccVRFDataSourceConfigWithCustomFields(vrfName, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[2]q
  object_types = ["ipam.vrf"]
  type         = "text"
}

resource "netbox_vrf" "test" {
  name = %[1]q
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[3]q
    }
  ]
}

data "netbox_vrf" "test" {
  name = netbox_vrf.test.name
  depends_on = [netbox_vrf.test]
}
`, vrfName, customFieldName, customFieldValue)
}
