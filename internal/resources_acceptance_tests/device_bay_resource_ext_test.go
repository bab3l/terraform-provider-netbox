package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

func TestAccDeviceBayResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-db-ext")
	siteSlug := testutil.RandomSlug("tf-test-site-db-ext")
	mfgName := testutil.RandomName("tf-test-mfg-db-ext")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-db-ext")
	dtModel := testutil.RandomName("tf-test-dt-db-ext")
	dtSlug := testutil.RandomSlug("tf-test-dt-db-ext")
	roleName := testutil.RandomName("tf-test-role-db-ext")
	roleSlug := testutil.RandomSlug("tf-test-role-db-ext")
	parentDeviceName := testutil.RandomName("tf-test-parent-db-ext")
	childDeviceName := testutil.RandomName("tf-test-child-db-ext")
	bayName := testutil.RandomName("tf-test-bay-db-ext")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(parentDeviceName)
	cleanup.RegisterDeviceCleanup(childDeviceName)

	testFields := map[string]string{
		"label": "Bay Label",
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_device_bay",
		BaseConfig: func() string {
			return testAccDeviceBayResourceConfig_minimal(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName)
		},
		ConfigWithFields: func() string {
			return testAccDeviceBayResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName, testFields)
		},
		OptionalFields: testFields,
		RequiredFields: map[string]string{
			"name": bayName,
		},
	})
}

func testAccDeviceBayResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName string, fields map[string]string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device" "parent" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device" "child" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device_bay" "test" {
  device = netbox_device.parent.id
  name   = %q
  label  = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName, fields["label"])
}
