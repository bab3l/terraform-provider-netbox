package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDevice_StatusNotInConfig verifies that the status field is not set in state
// when not specified in config, avoiding unwanted drift.
func TestAccDevice_StatusNotInConfig(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("tf-test-mfr-status")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-status")
	deviceTypeName := testutil.RandomName("tf-test-dt-status")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-status")
	deviceRoleName := testutil.RandomName("tf-test-dr-status")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-status")
	siteName := testutil.RandomName("tf-test-site-status")
	siteSlug := testutil.RandomSlug("tf-test-site-status")
	deviceName := testutil.RandomName("tf-test-device-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create device without status in config
			{
				Config: testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					// Status should NOT be present in state when not specified in config
					resource.TestCheckNoResourceAttr("netbox_device.test", "status"),
				),
			},
			// Step 2: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName),
			},
		},
	})
}

// TestAccDevice_StatusAddedThenRemoved verifies that status can be added to existing device
// and then removed without causing crashes or unwanted drift.
func TestAccDevice_StatusAddedThenRemoved(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("tf-test-mfr-add-rem")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-add-rem")
	deviceTypeName := testutil.RandomName("tf-test-dt-add-rem")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-add-rem")
	deviceRoleName := testutil.RandomName("tf-test-dr-add-rem")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-add-rem")
	siteName := testutil.RandomName("tf-test-site-add-rem")
	siteSlug := testutil.RandomSlug("tf-test-site-add-rem")
	deviceName := testutil.RandomName("tf-test-device-add-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create device without status
			{
				Config: testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckNoResourceAttr("netbox_device.test", "status"),
				),
			},
			// Step 2: Add status to existing device
			{
				Config: testAccDeviceConfig_withStatus(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, "active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
				),
			},
			// Step 3: Remove status from config - should not show in state
			{
				Config: testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					// Status should NOT be present in state when removed from config
					resource.TestCheckNoResourceAttr("netbox_device.test", "status"),
				),
			},
			// Step 4: Final plan-only verification - no changes
			{
				PlanOnly: true,
				Config:   testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName),
			},
		},
	})
}

func testAccDeviceConfig_statusNotInConfig(mfrName, mfrSlug, dtName, dtSlug, drName, drSlug, siteName, siteSlug, deviceName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}
`, mfrName, mfrSlug, dtName, dtSlug, drName, drSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceConfig_withStatus(mfrName, mfrSlug, dtName, dtSlug, drName, drSlug, siteName, siteSlug, deviceName, status string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = %[10]q
}
`, mfrName, mfrSlug, dtName, dtSlug, drName, drSlug, siteName, siteSlug, deviceName, status)
}
