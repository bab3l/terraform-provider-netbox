package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

func TestAccModuleBayResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-mbay-ext")
	siteSlug := testutil.RandomSlug("tf-test-site-mbay-ext")
	mfgName := testutil.RandomName("tf-test-mfg-mbay-ext")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-mbay-ext")
	dtModel := testutil.RandomName("tf-test-dt-mbay-ext")
	dtSlug := testutil.RandomSlug("tf-test-dt-mbay-ext")
	roleName := testutil.RandomName("tf-test-role-mbay-ext")
	roleSlug := testutil.RandomSlug("tf-test-role-mbay-ext")
	deviceName := testutil.RandomName("tf-test-device-mbay-ext")
	bayName := testutil.RandomName("tf-test-mbay-ext")
	position := testutil.RandomName("pos")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterModuleBayCleanup(bayName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_module_bay",
		BaseConfig: func() string {
			return testAccModuleBayResourceConfig_basic(
				siteName,
				siteSlug,
				mfgName,
				mfgSlug,
				dtModel,
				dtSlug,
				roleName,
				roleSlug,
				deviceName,
				bayName,
			)
		},
		ConfigWithFields: func() string {
			return testAccModuleBayResourceConfig_withPosition(
				siteName,
				siteSlug,
				mfgName,
				mfgSlug,
				dtModel,
				dtSlug,
				roleName,
				roleSlug,
				deviceName,
				bayName,
				position,
			)
		},
		OptionalFields: map[string]string{
			"position": position,
		},
		RequiredFields: map[string]string{
			"name": bayName,
		},
		CheckDestroy: testutil.CheckModuleBayDestroy,
	})
}

func testAccModuleBayResourceConfig_withPosition(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, position string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device   = netbox_device.test.id
  name     = %q
  position = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, position)
}
