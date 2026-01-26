package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANResource_basic(t *testing.T) {
	t.Parallel()

	ssid := testutil.RandomName("tf-test-ssid")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_wireless_lan.test", "vlan"),
					testutil.ReferenceFieldCheck("netbox_wireless_lan.test", "tenant"),
				),
			},
			{
				Config:             testAccWirelessLANResourceConfig_basic(ssid),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

func TestAccWirelessLANResource_IDPreservation(t *testing.T) {
	t.Parallel()

	ssid := testutil.RandomName("tf-test-ssid-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_basic(ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
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

func TestAccWirelessLANResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	ssid := testutil.RandomName("tf-test-ssid-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_withDescription(ssid, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccWirelessLANResourceConfig_withDescription(ssid, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func testAccWirelessLANResourceConfig_withDescription(ssid string, description string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
  ssid        = %[1]q
  description = %[2]q
}
`, ssid, description)
}

func TestAccWirelessLANResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	ssid := testutil.RandomName("tf-test-ssid-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_basic(ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find wireless LAN by SSID
					items, _, err := client.WirelessAPI.WirelessWirelessLansList(context.Background()).Ssid([]string{ssid}).Execute()
					if err != nil {
						t.Fatalf("Failed to list wireless LANs: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Wireless LAN not found with SSID: %s", ssid)
					}

					// Delete the wireless LAN
					itemID := items.Results[0].Id
					_, err = client.WirelessAPI.WirelessWirelessLansDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete wireless LAN: %v", err)
					}

					t.Logf("Successfully externally deleted wireless LAN with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANConsistencyConfig(wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "group"),
					testutil.ReferenceFieldCheck("netbox_wireless_lan.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccWirelessLANConsistencyConfig(wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug),
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
	tenant = netbox_tenant.test.id
}
`, wlanName, ssid, groupName, groupSlug, tenantName, tenantSlug)
}

func TestAccWirelessLANResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	ssid := testutil.RandomName("tf-test-wlan-rem")
	const testDescription = "Test Description"
	const testComments = "Test Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterWirelessLANCleanup(ssid)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_wireless_lan",
		BaseConfig: func() string {
			return testAccWirelessLANResourceConfig_basic(ssid)
		},
		ConfigWithFields: func() string {
			return testAccWirelessLANResourceConfig_removeOptionalFields_withFields(ssid)
		},
		OptionalFields: map[string]string{
			"auth_cipher": "aes",
			"auth_psk":    "test_psk_12345678",
			"auth_type":   "wpa-personal",
			"description": testDescription,
			"comments":    testComments,
		},
		RequiredFields: map[string]string{
			"ssid": ssid,
		},
		CheckDestroy: testutil.CheckWirelessLANDestroy,
	})
}

func testAccWirelessLANResourceConfig_removeOptionalFields_withFields(ssid string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
  ssid         = %[1]q
  auth_cipher  = "aes"
  auth_psk     = "test_psk_12345678"
  auth_type    = "wpa-personal"
  description  = "Test Description"
  comments     = "Test Comments"
}
`, ssid)
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

func TestAccWirelessLANResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_wireless_lan",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_ssid": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_wireless_lan" "test" {
  # ssid missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
