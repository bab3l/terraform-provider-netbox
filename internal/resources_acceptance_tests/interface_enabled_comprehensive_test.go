package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccInterfaceResource_EnabledComprehensive tests comprehensive scenarios for interface enabled field.
// This validates that Optional+Computed boolean fields work correctly across all scenarios.
func TestAccInterfaceResource_EnabledComprehensive(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		BaseConfig: func() string {
			return `
resource "netbox_site" "test" {
	name = "test-site-interface-enabled"
	slug = "test-site-interface-enabled"
}

resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-interface-enabled"
	slug = "test-manufacturer-interface-enabled"
}

resource "netbox_device_role" "test" {
	name = "test-device-role-interface-enabled"
	slug = "test-device-role-interface-enabled"
	color = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-interface-enabled"
	slug         = "test-device-type-interface-enabled"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	name        = "test-device-interface-enabled"
}

resource "netbox_interface" "test" {
	device = netbox_device.test.id
	name   = "eth0-enabled-test"
	type   = "1000base-t"
	# enabled field intentionally omitted - should get default true
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_site" "test" {
	name = "test-site-interface-enabled"
	slug = "test-site-interface-enabled"
}

resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-interface-enabled"
	slug = "test-manufacturer-interface-enabled"
}

resource "netbox_device_role" "test" {
	name = "test-device-role-interface-enabled"
	slug = "test-device-role-interface-enabled"
	color = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-interface-enabled"
	slug         = "test-device-type-interface-enabled"
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	name        = "test-device-interface-enabled"
}

resource "netbox_interface" "test" {
	device  = netbox_device.test.id
	name    = "eth0-enabled-test"
	type    = "1000base-t"
	enabled = ` + value + `
}
`
		},
	})
}
