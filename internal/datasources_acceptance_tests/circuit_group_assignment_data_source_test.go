package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests

func TestAccCircuitGroupAssignmentDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	groupName := testutil.RandomName("tf-test-cga-ds-id-group")
	groupSlug := testutil.RandomSlug("tf-test-cga-ds-id-grp")
	providerName := testutil.RandomName("tf-test-cga-ds-id-prov")
	providerSlug := testutil.RandomSlug("tf-test-cga-ds-id-prov")
	circuitTypeName := testutil.RandomName("tf-test-cga-ds-id-type")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-ds-id-type")
	circuitCid := testutil.RandomSlug("tf-test-cga-ds-id-ckt")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerName)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentDataSourceConfig_byID(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "circuit_id"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentDataSource_byID(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	groupName := testutil.RandomName("tf-test-cga-ds-group")
	groupSlug := testutil.RandomSlug("tf-test-cga-ds-grp")
	providerName := testutil.RandomName("tf-test-cga-ds-prov")
	providerSlug := testutil.RandomSlug("tf-test-cga-ds-prov")
	circuitTypeName := testutil.RandomName("tf-test-cga-ds-type")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-ds-type")
	circuitCid := testutil.RandomSlug("tf-test-cga-ds-ckt")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerName)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentDataSourceConfig_byID(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "circuit_id"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentDataSource_withPriority(t *testing.T) {
	t.Parallel()

	// Generate unique names
	groupName := testutil.RandomName("tf-test-cga-ds-grp-pri")
	groupSlug := testutil.RandomSlug("tf-test-cga-ds-grp-pri")
	providerName := testutil.RandomName("tf-test-cga-ds-prov-pri")
	providerSlug := testutil.RandomSlug("tf-test-cga-ds-prov-pri")
	circuitTypeName := testutil.RandomName("tf-test-cga-ds-type-pri")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-ds-type-pri")
	circuitCid := testutil.RandomSlug("tf-test-cga-ds-ckt-pri")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerName)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentDataSourceConfig_withPriority(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, "primary",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group_assignment.test", "priority", "primary"),
				),
			},
		},
	})
}

func testAccCircuitGroupAssignmentDataSourceConfig_byID(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit_type" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_circuit" "test" {
  cid              = %[7]q
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
}

data "netbox_circuit_group_assignment" "test" {
  id = netbox_circuit_group_assignment.test.id
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid)
}

func testAccCircuitGroupAssignmentDataSourceConfig_withPriority(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, priority string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit_type" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_circuit" "test" {
  cid              = %[7]q
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
  priority   = %[8]q
}

data "netbox_circuit_group_assignment" "test" {
  id = netbox_circuit_group_assignment.test.id
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, priority)
}
