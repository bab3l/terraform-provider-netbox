package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccRearPortResource_Positions tests comprehensive scenarios for rear port positions field.
// This validates that Optional+Computed int32 fields with proper defaults work correctly.
func TestAccRearPortResource_Positions(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port",
		OptionalField:  "positions",
		DefaultValue:   "1",
		FieldTestValue: "4",
		BaseConfig: func() string {
			return `
resource "netbox_site" "test" {
	name = "test-site-rear-port"
	slug = "test-site-rear-port"
}

resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port"
	slug = "test-manufacturer-rear-port"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port"
	slug         = "test-device-type-rear-port"
}

resource "netbox_device_role" "test" {
	name = "test-role-rear-port"
	slug = "test-role-rear-port"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	name        = "test-device-rear-port"
	site        = netbox_site.test.id
	role        = netbox_device_role.test.id
}

resource "netbox_rear_port" "test" {
	device = netbox_device.test.id
	name   = "rear-port-positions-test"
	type   = "8p8c"
	# positions field intentionally omitted - should get default 1
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_site" "test" {
	name = "test-site-rear-port"
	slug = "test-site-rear-port"
}

resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port"
	slug = "test-manufacturer-rear-port"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port"
	slug         = "test-device-type-rear-port"
}

resource "netbox_device_role" "test" {
	name = "test-role-rear-port"
	slug = "test-role-rear-port"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	name        = "test-device-rear-port"
	site        = netbox_site.test.id
	role        = netbox_device_role.test.id
}

resource "netbox_rear_port" "test" {
	device    = netbox_device.test.id
	name      = "rear-port-positions-test"
	type      = "8p8c"
	positions = ` + value + `
}
`
		},
	})
}
