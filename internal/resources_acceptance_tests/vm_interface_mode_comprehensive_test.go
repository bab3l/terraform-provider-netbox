package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccVMInterfaceResource_ModeOptionalField tests comprehensive scenarios for VM interface mode.
// This validates that Optional Only fields work correctly across all scenarios.
func TestAccVMInterfaceResource_ModeOptionalField(t *testing.T) {

	testutil.RunOptionalFieldTestSuite(t, testutil.OptionalFieldTestConfig{
		ResourceName:   "netbox_vm_interface",
		OptionalField:  "mode",
		FieldTestValue: "tagged",
		BaseConfig: func() string {
			return `
			resource "netbox_cluster_type" "test" {
				name = "test-cluster-type-vmif-mode"
				slug = "test-cluster-type-vmif-mode"
			}

			resource "netbox_cluster" "test" {
				name = "test-cluster-vmif-mode"
				type = netbox_cluster_type.test.id
			}

			resource "netbox_virtual_machine" "test" {
				name    = "test-vm-vmif-mode"
				cluster = netbox_cluster.test.id
			}

			resource "netbox_vm_interface" "test" {
				virtual_machine = netbox_virtual_machine.test.id
				name            = "eth0-mode-test"
				# mode field intentionally omitted - should be absent in state
			}
			`
		},
		WithFieldConfig: func(value string) string {
			return `
			resource "netbox_cluster_type" "test" {
				name = "test-cluster-type-vmif-mode"
				slug = "test-cluster-type-vmif-mode"
			}

			resource "netbox_cluster" "test" {
				name = "test-cluster-vmif-mode"
				type = netbox_cluster_type.test.id
			}

			resource "netbox_virtual_machine" "test" {
				name    = "test-vm-vmif-mode"
				cluster = netbox_cluster.test.id
			}

			resource "netbox_vm_interface" "test" {
				virtual_machine = netbox_virtual_machine.test.id
				name            = "eth0-mode-test"
				mode            = "` + value + `"
			}
			`
		},
	})
}
