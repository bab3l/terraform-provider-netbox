package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccDeviceResource_StatusOptionalField tests comprehensive scenarios for the device status optional field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccDeviceResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_device",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		BaseConfig: func() string {
			return `
			resource "netbox_site" "test" {
				name = "test-site-device-status"
				slug = "test-site-device-status"
			}

			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer-device-status"
				slug = "test-manufacturer-device-status"
			}

			resource "netbox_device_role" "test" {
				name = "test-device-role-status"
				slug = "test-device-role-status"
				color = "aa1409"
				vm_role = false
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type-status"
				slug         = "test-device-type-status"
			}

			resource "netbox_device" "test" {
				device_type = netbox_device_type.test.id
				role        = netbox_device_role.test.id
				site        = netbox_site.test.id
				name        = "test-device-status"
				# status field intentionally omitted - should get default "active"
			}
			`
		},
		WithFieldConfig: func(value string) string {
			return `
			resource "netbox_site" "test" {
				name = "test-site-device-status"
				slug = "test-site-device-status"
			}

			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer-device-status"
				slug = "test-manufacturer-device-status"
			}

			resource "netbox_device_role" "test" {
				name = "test-device-role-status"
				slug = "test-device-role-status"
				color = "aa1409"
				vm_role = false
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type-status"
				slug         = "test-device-type-status"
			}

			resource "netbox_device" "test" {
				device_type = netbox_device_type.test.id
				role        = netbox_device_role.test.id
				site        = netbox_site.test.id
				name        = "test-device-status"
				status      = "` + value + `"
			}
			`
		},
	})
}
