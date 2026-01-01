package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccVirtualMachineResource_StatusOptionalField tests comprehensive scenarios for virtual machine status.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccVirtualMachineResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	// Generate unique names for this test run
	clusterTypeName := testutil.RandomName("tf-test-ct-vm-status")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-vm-status")
	clusterName := testutil.RandomName("tf-test-cluster-vm-status")
	vmName := testutil.RandomName("tf-test-vm-status")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_virtual_machine",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "staged", CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterGroupDestroy,
			testutil.CheckRoleDestroy,
		), BaseConfig: func() string {
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
				# status field intentionally omitted - should get default "active"
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
				status  = "` + value + `"
			}
			`
		},
	})
}
