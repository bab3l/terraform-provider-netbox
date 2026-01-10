//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNTerminationDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_l2vpn_term_ds_cf")
	l2vpnName := testutil.RandomName("tf-test-l2vpn-ds-cf")
	siteName := testutil.RandomName("tf-test-site-l2vpn-cf")
	vlanName := testutil.RandomName("tf-test-vlan-l2vpn-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationDataSourceConfig_customFields(customFieldName, l2vpnName, siteName, vlanName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn_termination.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_l2vpn_termination.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_l2vpn_termination.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_l2vpn_termination.test", "custom_fields.0.value", "test-l2vpn-term-value"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationDataSourceConfig_customFields(customFieldName, l2vpnName, siteName, vlanName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  object_types = ["vpn.l2vpntermination"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_vlan" "test" {
  name = %[4]q
  vid  = 100
  site = netbox_site.test.slug
}

resource "netbox_l2vpn" "test" {
  name = %[2]q
  slug = %[2]q
  type = "vpws"
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-l2vpn-term-value"
    }
  ]
}

data "netbox_l2vpn_termination" "test" {
  id = netbox_l2vpn_termination.test.id

  depends_on = [netbox_l2vpn_termination.test]
}
`, customFieldName, l2vpnName, siteName, vlanName)
}
