package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccVMInterfaceResource_ModeOptionalField tests comprehensive scenarios for VM interface mode.
// This validates that Optional Only fields work correctly across all scenarios.
func TestAccVMInterfaceResource_ModeOptionalField(t *testing.T) {
	// Generate unique names for this test run
	clusterTypeName := testutil.RandomName("tf-test-ct-vmif-mode")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-vmif-mode")
	clusterName := testutil.RandomName("tf-test-cluster-vmif-mode")
	vmName := testutil.RandomName("tf-test-vm-vmif-mode")
	interfaceName := testutil.RandomName("eth0-mode-test")

	testutil.RunOptionalFieldTestSuite(t, testutil.OptionalFieldTestConfig{
		ResourceName:   "netbox_vm_interface",
		OptionalField:  "mode",
		FieldTestValue: "tagged",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		BaseConfig: func() string {
			return `
			resource "netbox_cluster_type" "test" {
				name = "` + clusterTypeName + `"
				slug = "` + clusterTypeSlug + `"
			}

			resource "netbox_cluster" "test" {
				name = "` + clusterName + `"
				type = netbox_cluster_type.test.id
			}

			resource "netbox_virtual_machine" "test" {
				name    = "` + vmName + `"
				cluster = netbox_cluster.test.id
			}

			resource "netbox_vm_interface" "test" {
				virtual_machine = netbox_virtual_machine.test.id
				name            = "` + interfaceName + `"
				# mode field intentionally omitted - should be absent in state
			}
			`
		},
		WithFieldConfig: func(value string) string {
			return `
			resource "netbox_cluster_type" "test" {
				name = "` + clusterTypeName + `"
				slug = "` + clusterTypeSlug + `"
			}

			resource "netbox_cluster" "test" {
				name = "` + clusterName + `"
				type = netbox_cluster_type.test.id
			}

			resource "netbox_virtual_machine" "test" {
				name    = "` + vmName + `"
				cluster = netbox_cluster.test.id
			}

			resource "netbox_vm_interface" "test" {
				virtual_machine = netbox_virtual_machine.test.id
				name            = "` + interfaceName + `"
				mode            = "` + value + `"
			}
			`
		},
	})
}
