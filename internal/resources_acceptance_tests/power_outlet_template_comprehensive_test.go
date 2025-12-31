package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccPowerOutletTemplateResource_Label tests comprehensive scenarios for power outlet template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccPowerOutletTemplateResource_Label(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_outlet_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "Outlet-01",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-outlet-template"
	slug = "test-manufacturer-outlet-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-outlet-template"
	slug         = "test-device-type-outlet-template"
}

resource "netbox_power_outlet_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "power-outlet-template-label-test"
	type        = "iec-60320-c13"
	# label field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-outlet-template"
	slug = "test-manufacturer-outlet-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-outlet-template"
	slug         = "test-device-type-outlet-template"
}

resource "netbox_power_outlet_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "power-outlet-template-label-test"
	type        = "iec-60320-c13"
	label       = "` + value + `"
}
`
		},
	})
}
