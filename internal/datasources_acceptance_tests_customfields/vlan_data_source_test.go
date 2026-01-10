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

func TestAccVLANDataSource_customFields(t *testing.T) {
	vlanName := testutil.RandomName("tf-test-vlan-ds-cf")
	vid := int32(acctest.RandIntRange(100, 4000))
	customFieldName := testutil.RandomCustomFieldName("tf_test_vlan_ds_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANDataSourceConfigWithCustomFields(vlanName, vid, customFieldName, customFieldValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccVLANDataSourceConfigWithCustomFields(vlanName string, vid int32, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[3]q
  object_types = ["ipam.vlan"]
  type         = "text"
}

resource "netbox_vlan" "test" {
  name   = %[1]q
  vid    = %[2]d
  status = "active"
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[4]q
    }
  ]
}

data "netbox_vlan" "test" {
  vid = netbox_vlan.test.vid
  depends_on = [netbox_vlan.test]
}
`, vlanName, vid, customFieldName, customFieldValue)
}
