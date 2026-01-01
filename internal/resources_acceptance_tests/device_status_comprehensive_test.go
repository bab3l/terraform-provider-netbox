package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccDeviceResource_StatusOptionalField tests comprehensive scenarios for the device status optional field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccDeviceResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	// Generate unique names for this test run
	siteName := testutil.RandomName("tf-test-site-device-status")
	siteSlug := testutil.RandomSlug("tf-test-site-device-status")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-device-status")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-device-status")
	deviceRoleName := testutil.RandomName("tf-test-device-role-status")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-status")
	deviceTypeName := testutil.RandomName("tf-test-device-type-status")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-status")
	deviceName := testutil.RandomName("tf-test-device-status")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_device",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		CheckDestroy: testutil.ComposeCheckDestroy(
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
				# status field intentionally omitted - should get default "active"
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
				status      = "` + value + `"
			}
			`
		},
	})
}
