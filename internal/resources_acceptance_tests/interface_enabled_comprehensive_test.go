package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccInterfaceResource_EnabledComprehensive tests comprehensive scenarios for interface enabled field.
// This validates that Optional+Computed boolean fields work correctly across all scenarios.
func TestAccInterfaceResource_EnabledComprehensive(t *testing.T) {
	// Generate unique names for this test run
	siteName := testutil.RandomName("tf-test-site-interface")
	siteSlug := testutil.RandomSlug("tf-test-site-interface")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-interface")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-interface")
	deviceRoleName := testutil.RandomName("tf-test-device-role-interface")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-interface")
	deviceTypeName := testutil.RandomName("tf-test-device-type-interface")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-interface")
	deviceName := testutil.RandomName("tf-test-device-interface")
	interfaceName := testutil.RandomName("eth0-enabled-test")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckInterfaceDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return `
resource "netbox_site" "test" {
	name = "` + siteName + `"
	slug = "` + siteSlug + `"
}

resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_role" "test" {
	name = "` + deviceRoleName + `"
	slug = "` + deviceRoleSlug + `"
	color = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	name        = "` + deviceName + `"
}

resource "netbox_interface" "test" {
	device = netbox_device.test.id
	name   = "` + interfaceName + `"
	type   = "1000base-t"
	# enabled field intentionally omitted - should get default true
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_site" "test" {
	name = "` + siteName + `"
	slug = "` + siteSlug + `"
}

resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_role" "test" {
	name = "` + deviceRoleName + `"
	slug = "` + deviceRoleSlug + `"
	color = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	name        = "` + deviceName + `"
}

resource "netbox_interface" "test" {
	device  = netbox_device.test.id
	name    = "` + interfaceName + `"
	type    = "1000base-t"
	enabled = ` + value + `
}
`
		},
	})
}
