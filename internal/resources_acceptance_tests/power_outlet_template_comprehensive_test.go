package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccPowerOutletTemplateResource_Label tests comprehensive scenarios for power outlet template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccPowerOutletTemplateResource_Label(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-pwr-out-tpl")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pwr-out-tpl")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-pwr-out-tpl")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-pwr-out-tpl")
	powerOutletTemplateName := testutil.RandomName("tf-test-pwr-out-tpl")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_outlet_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "Outlet-01",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_power_outlet_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + powerOutletTemplateName + `"
	type        = "iec-60320-c13"
	# label field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_power_outlet_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + powerOutletTemplateName + `"
	type        = "iec-60320-c13"
	label       = "` + value + `"
}
`
		},
	})
}
