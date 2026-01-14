package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	wirelessInterfaceNameA = "wlan0"
	wirelessInterfaceNameB = "wlan1"
)

func TestAccWirelessLinkResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-wireless")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-wireless")
	interfaceNameA := wirelessInterfaceNameA
	interfaceNameB := wirelessInterfaceNameB

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", "Test SSID"),
				),
			},
			{
				ResourceName:            "netbox_wireless_link.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"interface_a", "interface_b"},
			},
		},
	})
}

func TestAccWirelessLinkResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-wl-id")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("device-wl-id")
	interfaceNameA := wirelessInterfaceNameA
	interfaceNameB := wirelessInterfaceNameB

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_link.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", "Test SSID"),
				),
			},
		},
	})
}

func testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string) string {
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeName := testutil.RandomName("dtype")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  model = %[7]q
  slug  = %[8]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%[9]s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%[9]s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %[10]q
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = %[11]q
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "Test SSID"
  status      = "connected"
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB)
}

func TestAccConsistency_WirelessLink_LiteralNames(t *testing.T) {
	t.Parallel()
	ssid := testutil.RandomName("tf-test-ssid-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkConsistencyLiteralNamesConfig(ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_link.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", ssid),
				),
			},
		},
	})
}

func testAccWirelessLinkConsistencyLiteralNamesConfig(ssid string) string {
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeName := testutil.RandomName("dtype")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.GenerateSlug(siteName)

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%[10]s"
  slug   = "%[11]s"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_device_role" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_device_type" "test" {
  model        = "%[7]s"
  slug         = "%[8]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name        = "%[9]s-a"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = "%[9]s-b"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name   = "wlan0"
  device = netbox_device.test_a.id
  type   = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name   = "wlan1"
  device = netbox_device.test_b.id
  type   = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "%[1]s"
  status      = "connected"
}
`, ssid, siteName, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug)
}

func TestAccWirelessLinkResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("test-site-wireless-upd")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-wireless-upd")
	interfaceNameA := wirelessInterfaceNameA
	interfaceNameB := wirelessInterfaceNameB

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkResourceConfig_withDescription(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccWirelessLinkResourceConfig_withDescription(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func testAccWirelessLinkResourceConfig_withDescription(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string, description string) string {
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeName := testutil.RandomName("dtype")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  model = %[7]q
  slug  = %[8]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%[9]s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%[9]s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %[10]q
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = %[11]q
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "Test SSID"
  status      = "connected"
  description = %[12]q
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB, description)
}

func TestAccWirelessLinkResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("test-site-wireless-extdel")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-wireless-extdel")
	interfaceNameA := wirelessInterfaceNameA
	interfaceNameB := wirelessInterfaceNameB
	ssid := testutil.RandomName("tf-test-ssid-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkResourceConfig_forExternalDeletion(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB, ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_link.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find wireless link by SSID
					items, _, err := client.WirelessAPI.WirelessWirelessLinksList(context.Background()).Ssid([]string{ssid}).Execute()
					if err != nil {
						t.Fatalf("Failed to list wireless links: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Wireless link not found with SSID: %s", ssid)
					}

					// Delete the wireless link
					itemID := items.Results[0].Id
					_, err = client.WirelessAPI.WirelessWirelessLinksDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete wireless link: %v", err)
					}

					t.Logf("Successfully externally deleted wireless link with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccWirelessLinkResourceConfig_forExternalDeletion(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string, ssid string) string {
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeName := testutil.RandomName("dtype")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  model = %[7]q
  slug  = %[8]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%[9]s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%[9]s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %[10]q
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = %[11]q
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = %[12]q
  status      = "connected"
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB, ssid)
}

// TestAccWirelessLinkResource_removeOptionalFields tests that removing previously set optional fields correctly sets them to null.
// This addresses the bug where removing nullable fields from the configuration would not clear them in NetBox,
// causing "Provider produced inconsistent result after apply" errors.
func TestAccWirelessLinkResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("wl-rem")

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_wireless_link",
		BaseConfig: func() string {
			return testAccWirelessLinkResourceConfig_removeOptionalFields_base(name)
		},
		ConfigWithFields: func() string {
			return testAccWirelessLinkResourceConfig_removeOptionalFields_withFields(name)
		},
		OptionalFields: map[string]string{
			"auth_type":     "wpa-personal",
			"auth_cipher":   "aes",
			"auth_psk":      "secret-key-value",
			"description":   "Test description",
			"comments":      "Test comments",
			"distance":      "10.5",
			"distance_unit": "km",
		},
		RequiredFields: map[string]string{
			"ssid":   "Test SSID",
			"status": "connected",
		},
	})
}

func testAccWirelessLinkResourceConfig_removeOptionalFields_base(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[1]s-site"
}

resource "netbox_tenant" "test" {
  name = "%[1]s-tenant"
  slug = "%[1]s-tenant"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[1]s-mfr"
}

resource "netbox_device_role" "test" {
  name  = "%[1]s-role"
  slug  = "%[1]s-role"
  color = "ff5733"
}

resource "netbox_device_type" "test" {
  model = "%[1]s-dtype"
  slug  = "%[1]s-dtype"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%[1]s-dev-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%[1]s-dev-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = "wlan0"
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = "wlan1"
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "Test SSID"
  status      = "connected"
}
`, name)
}

func testAccWirelessLinkResourceConfig_removeOptionalFields_withFields(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[1]s-site"
}

resource "netbox_tenant" "test" {
  name = "%[1]s-tenant"
  slug = "%[1]s-tenant"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[1]s-mfr"
}

resource "netbox_device_role" "test" {
  name  = "%[1]s-role"
  slug  = "%[1]s-role"
  color = "ff5733"
}

resource "netbox_device_type" "test" {
  model = "%[1]s-dtype"
  slug  = "%[1]s-dtype"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%[1]s-dev-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%[1]s-dev-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = "wlan0"
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = "wlan1"
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "Test SSID"
  status      = "connected"
  tenant      = netbox_tenant.test.id
  auth_type   = "wpa-personal"
  auth_cipher = "aes"
  auth_psk    = "secret-key-value"
  description = "Test description"
  comments    = "Test comments"
  distance    = 10.5
  distance_unit = "km"
}
`, name)
}
