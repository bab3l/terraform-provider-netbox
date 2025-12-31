package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDevice_StatusNotInConfig verifies that when status is not specified in config,
// but the device has a status in Netbox, there is no drift.
// This tests the bug where Terraform would show: status = null -> "active".
func TestAccDevice_StatusNotInConfig(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("tf-test-mfr-status")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-status")
	deviceTypeName := testutil.RandomName("tf-test-dt-status")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-status")
	siteName := testutil.RandomName("tf-test-site-status")
	siteSlug := testutil.RandomSlug("tf-test-site-status")
	deviceName := testutil.RandomName("tf-test-device-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create device without status in config
			{
				Config: testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					// status should not be set in state when not in config
					resource.TestCheckNoResourceAttr("netbox_device.test", "status"),
				),
			},
			// Step 2: Refresh state and verify no drift (status should remain unset)
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					// After refresh, status should still not be in state
					resource.TestCheckNoResourceAttr("netbox_device.test", "status"),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, deviceName),
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
	siteName := testutil.RandomName("tf-test-site-add-rem")
	siteSlug := testutil.RandomSlug("tf-test-site-add-rem")
	deviceName := testutil.RandomName("tf-test-device-add-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create device without status
			{
				Config: testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
				),
			},
			// Step 2: Add status to config
			{
				Config: testAccDeviceConfig_withStatus(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, deviceName, "planned"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "status", "planned"),
				),
			},
			// Step 3: Remove status from config - should not crash
			{
				Config: testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					// status should be removed/null when not in config
					resource.TestCheckNoResourceAttr("netbox_device.test", "status"),
				),
			},
			// Step 4: Plan only - verify no drift
			{
				PlanOnly: true,
				Config:   testAccDeviceConfig_statusNotInConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, deviceName),
			},
		},
	})
}

func testAccDeviceConfig_statusNotInConfig(mfrName, mfrSlug, dtName, dtSlug, siteName, siteSlug, deviceName string) string {
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

resource "netbox_site" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device" "test" {
  name        = %[7]q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
}
`, mfrName, mfrSlug, dtName, dtSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceConfig_withStatus(mfrName, mfrSlug, dtName, dtSlug, siteName, siteSlug, deviceName, status string) string {
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

resource "netbox_site" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device" "test" {
  name        = %[7]q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  status      = %[8]q
}
`, mfrName, mfrSlug, dtName, dtSlug, siteName, siteSlug, deviceName, status)
}
