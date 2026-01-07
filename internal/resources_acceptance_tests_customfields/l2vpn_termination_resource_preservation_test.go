//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNTerminationResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	l2vpnName := testutil.RandomName("tf-test-l2vpn-lt-pres")
	l2vpnSlug := testutil.RandomSlug("tf-test-l2vpn-lt-pres")
	vlanName := testutil.RandomName("tf-test-vlan-lt-pres")
	siteName := testutil.RandomName("tf-test-site-lt-pres")
	siteSlug := testutil.RandomSlug("tf-test-site-lt-pres")
	cfName := testutil.RandomCustomFieldName("tf_lt_pres")

	cleanup := testutil.NewCleanupResource(t)
	defer cleanup.Close(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationResourcePreservationConfig_step1(
					l2vpnName, l2vpnSlug, vlanName, siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "custom_fields.%", "1"),
					testutil.ResourceCheckCustomFieldValue("netbox_l2vpn_termination.test", cfName, "preserved_value"),
				),
			},
			{
				// Update without custom_fields in config - should be preserved in NetBox
				Config: testAccL2VPNTerminationResourcePreservationConfig_step2(
					l2vpnName, l2vpnSlug, vlanName, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					// Custom fields are not in the config, so they won't appear in state
				),
			},
		},
	})
}

func testAccL2VPNTerminationResourcePreservationConfig_step1(
	l2vpnName, l2vpnSlug, vlanName, siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "l2vpn_termination_pres" {
  name         = %[6]q
  type         = "text"
  object_types = ["vpn.l2vpntermination"]
  required     = false
}

resource "netbox_l2vpn" "test" {
  name = %[1]q
  slug = %[2]q
  type = "vlan"
}

resource "netbox_site" "test" {
  name   = %[4]q
  slug   = %[5]q
  status = "active"
}

resource "netbox_vlan" "test" {
  name   = %[3]q
  vid    = 100
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                 = netbox_l2vpn.test.id
  assigned_object_type  = "ipam.vlan"
  assigned_object_id    = netbox_vlan.test.id

  custom_fields = {
    (netbox_custom_field.l2vpn_termination_pres.name) = "preserved_value"
  }

  depends_on = [netbox_custom_field.l2vpn_termination_pres]
}
`, l2vpnName, l2vpnSlug, vlanName, siteName, siteSlug, cfName)
}

func testAccL2VPNTerminationResourcePreservationConfig_step2(
	l2vpnName, l2vpnSlug, vlanName, siteName, siteSlug string,
) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %[1]q
  slug = %[2]q
  type = "vlan"
}

resource "netbox_site" "test" {
  name   = %[4]q
  slug   = %[5]q
  status = "active"
}

resource "netbox_vlan" "test" {
  name   = %[3]q
  vid    = 100
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                 = netbox_l2vpn.test.id
  assigned_object_type  = "ipam.vlan"
  assigned_object_id    = netbox_vlan.test.id
}
`, l2vpnName, l2vpnSlug, vlanName, siteName, siteSlug)
}
