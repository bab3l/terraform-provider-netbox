package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccVirtualMachineResource_StatusOptionalField tests comprehensive scenarios for virtual machine status.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccVirtualMachineResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_virtual_machine",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "staged",
		BaseConfig: func() string {
			return `
			resource "netbox_cluster_type" "test" {
				name = "test-cluster-type-vm-status"
				slug = "test-cluster-type-vm-status"
			}

			resource "netbox_cluster" "test" {
				name = "test-cluster-vm-status"
				type = netbox_cluster_type.test.id
			}

			resource "netbox_virtual_machine" "test" {
				name    = "test-vm-status"
				cluster = netbox_cluster.test.id
				# status field intentionally omitted - should get default "active"
			}
			`
		},
		WithFieldConfig: func(value string) string {
			return `
			resource "netbox_cluster_type" "test" {
				name = "test-cluster-type-vm-status"
				slug = "test-cluster-type-vm-status"
			}

			resource "netbox_cluster" "test" {
				name = "test-cluster-vm-status"
				type = netbox_cluster_type.test.id
			}

			resource "netbox_virtual_machine" "test" {
				name    = "test-vm-status"
				cluster = netbox_cluster.test.id
				status  = "` + value + `"
			}
			`
		},
	})
}
