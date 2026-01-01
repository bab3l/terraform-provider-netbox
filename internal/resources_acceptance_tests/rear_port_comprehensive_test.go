package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccRearPortResource_Positions tests comprehensive scenarios for rear port positions field.
// This validates that Optional+Computed int32 fields with proper defaults work correctly.
func TestAccRearPortResource_Positions(t *testing.T) {
	// Generate unique names for this test run
	siteName := testutil.RandomName("tf-test-site-rear-port")
	siteSlug := testutil.RandomSlug("tf-test-site-rear-port")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-rear-port")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-rear-port")
	deviceTypeName := testutil.RandomName("tf-test-device-type-rear-port")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-rear-port")
	deviceRoleName := testutil.RandomName("tf-test-role-rear-port")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-rear-port")
	deviceName := testutil.RandomName("tf-test-device-rear-port")
	rearPortName := testutil.RandomName("rear-port-positions-test")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port",
		OptionalField:  "positions",
		DefaultValue:   "1",
		FieldTestValue: "4",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRearPortDestroy,
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

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_device_role" "test" {
	name = "` + deviceRoleName + `"
	slug = "` + deviceRoleSlug + `"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + deviceName + `"
	site        = netbox_site.test.id
	role        = netbox_device_role.test.id
}

resource "netbox_rear_port" "test" {
	device = netbox_device.test.id
	name   = "` + rearPortName + `"
	type   = "8p8c"
	# positions field intentionally omitted - should get default 1
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

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_device_role" "test" {
	name = "` + deviceRoleName + `"
	slug = "` + deviceRoleSlug + `"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + deviceName + `"
	site        = netbox_site.test.id
	role        = netbox_device_role.test.id
}

resource "netbox_rear_port" "test" {
	device    = netbox_device.test.id
	name      = "` + rearPortName + `"
	type      = "8p8c"
	positions = ` + value + `
}
`
		},
	})
}
