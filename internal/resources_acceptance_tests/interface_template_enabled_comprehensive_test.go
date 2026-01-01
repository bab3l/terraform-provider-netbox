package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccInterfaceTemplateResource_EnabledComprehensive tests comprehensive scenarios for interface template enabled field.
// This validates that Optional+Computed boolean fields work correctly across all scenarios.
func TestAccInterfaceTemplateResource_EnabledComprehensive(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-interface-template-enabled"
	slug = "test-manufacturer-interface-template-enabled"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-interface-template-enabled"
	slug         = "test-device-type-interface-template-enabled"
}

resource "netbox_interface_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "eth0-template-enabled-test"
	type        = "1000base-t"
	# enabled field intentionally omitted - should get default true
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-interface-template-enabled"
	slug = "test-manufacturer-interface-template-enabled"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-interface-template-enabled"
	slug         = "test-device-type-interface-template-enabled"
}

resource "netbox_interface_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "eth0-template-enabled-test"
	type        = "1000base-t"
	enabled     = ` + value + `
}
`
		},
	})
}
