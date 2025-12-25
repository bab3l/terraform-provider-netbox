package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANResource_basic(t *testing.T) {

	t.Parallel()
	ssid := testutil.RandomName("tf-test-ssid")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_basic(ssid),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "status", "active"),
				),
			},
			{
				ResourceName:      "netbox_wireless_lan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWirelessLANResource_full(t *testing.T) {

	t.Parallel()
	ssid := testutil.RandomName("tf-test-ssid-full")
	groupName := testutil.RandomName("tf-test-wlan-group")
	groupSlug := testutil.RandomSlug("tf-test-wlan-group")
	description := "Test wireless LAN with all fields"
	updatedDescription := "Updated wireless LAN description"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_full(ssid, groupName, groupSlug, description, "active"),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", description),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "status", "active"),
				),
			},
			{
				Config: testAccWirelessLANResourceConfig_full(ssid, groupName, groupSlug, updatedDescription, "disabled"),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "status", "disabled"),
				),
			},
		},
	})
}

func testAccWirelessLANResourceConfig_basic(ssid string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
  ssid = %q
}
`, ssid)
}

func testAccWirelessLANResourceConfig_full(ssid, groupName, groupSlug, description, status string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_wireless_lan" "test" {
  ssid        = %q
  group       = netbox_wireless_lan_group.test.id
  description = %q
  status      = %q
}
`, groupName, groupSlug, ssid, description, status)
}

func TestAccConsistency_WirelessLAN(t *testing.T) {

	t.Parallel()

	wlanName := testutil.RandomName("wlan")
	ssid := testutil.RandomName("ssid")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANConsistencyConfig(wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "group"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,

				Config: testAccWirelessLANConsistencyConfig(wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccWirelessLANConsistencyConfig(wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_tenant" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_wireless_lan" "test" {
  ssid = "%[2]s"
  group = netbox_wireless_lan_group.test.id
  tenant = netbox_tenant.test.name
}
`, wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug)
}

func TestAccConsistency_WirelessLAN_LiteralNames(t *testing.T) {
	t.Parallel()
	ssid := testutil.RandomName("tf-test-ssid-lit")
	groupName := testutil.RandomName("tf-test-wlan-group-lit")
	groupSlug := testutil.RandomSlug("tf-test-wlan-group-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANConsistencyLiteralNamesConfig(ssid, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
				),
			},
			{
				Config:   testAccWirelessLANConsistencyLiteralNamesConfig(ssid, groupName, groupSlug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
				),
			},
		},
	})
}

func testAccWirelessLANConsistencyLiteralNamesConfig(ssid, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_wireless_lan" "test" {
  ssid  = %q
  group = netbox_wireless_lan_group.test.id
}
`, groupName, groupSlug, ssid)
}
