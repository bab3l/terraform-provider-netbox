package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccPowerPortTemplateResource_Label tests comprehensive scenarios for power port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccPowerPortTemplateResource_Label(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_port_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "Port-01",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-power-port-template"
	slug = "test-manufacturer-power-port-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-power-port-template"
	slug         = "test-device-type-power-port-template"
}

resource "netbox_power_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "power-port-template-label-test"
	type        = "iec-60320-c14"
	# label field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-power-port-template"
	slug = "test-manufacturer-power-port-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-power-port-template"
	slug         = "test-device-type-power-port-template"
}

resource "netbox_power_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "power-port-template-label-test"
	type        = "iec-60320-c14"
	label       = "` + value + `"
}
`
		},
	})
}
