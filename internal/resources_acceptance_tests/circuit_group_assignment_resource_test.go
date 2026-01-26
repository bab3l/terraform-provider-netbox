package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitGroupAssignmentResource_basic(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-group")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp")
	providerName := testutil.RandomName("tf-test-cga-provider")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov")
	circuitTypeName := testutil.RandomName("tf-test-cga-type")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_basic(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "circuit_id"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentResource_full(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-grp-full")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp-full")
	providerName := testutil.RandomName("tf-test-cga-prov-full")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-full")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-full")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-full")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-full")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_withPriority(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, "primary"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group_assignment.test", "priority", "primary"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentResource_withPriority(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-grp-pri")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp-pri")
	providerName := testutil.RandomName("tf-test-cga-prov-pri")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-pri")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-pri")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-pri")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-pri")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_withPriority(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, "primary"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group_assignment.test", "priority", "primary"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentResource_update(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-grp-upd")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp-upd")
	providerName := testutil.RandomName("tf-test-cga-prov-upd")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-upd")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-upd")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-upd")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-upd")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_basic(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
				),
			},
			{
				Config: testAccCircuitGroupAssignmentResourceConfig_withPriority(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, "secondary"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_group_assignment.test", "priority", "secondary"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentResource_import(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-grp-imp")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp-imp")
	providerName := testutil.RandomName("tf-test-cga-prov-imp")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-imp")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-imp")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-imp")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-imp")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_basic(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
			},
			{
				ResourceName:            "netbox_circuit_group_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "circuit_id"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_circuit_group_assignment.test", "group_id"),
					testutil.ReferenceFieldCheck("netbox_circuit_group_assignment.test", "circuit_id"),
				),
			},
			{
				Config:   testAccCircuitGroupAssignmentResourceConfig_basic(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCircuitGroupAssignmentResource_IDPreservation(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("cga-id")
	groupSlug := testutil.RandomSlug("cga-id")
	providerName := testutil.RandomName("cga-prov-id")
	providerSlug := testutil.RandomSlug("cga-prov-id")
	circuitTypeName := testutil.RandomName("cga-type-id")
	circuitTypeSlug := testutil.RandomSlug("cga-type-id")
	circuitCid := testutil.RandomSlug("cga-ckt-id")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_basic(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "circuit_id"),
				),
			},
		},
	})
}

func testAccCircuitGroupAssignmentResourceConfig_basic(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid string) string {
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
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid)
}

func testAccCircuitGroupAssignmentResourceConfig_withPriority(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, priority string) string {
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
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
  priority   = %[8]q
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, priority)
}

func TestAccConsistency_CircuitGroupAssignment_LiteralNames(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-grp-lit")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp-lit")
	providerName := testutil.RandomName("tf-test-cga-prov-lit")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-lit")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-lit")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-lit")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentResourceConfigLiteralNames(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_circuit_group_assignment.test", "group_id", "netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_circuit_group_assignment.test", "circuit_id", "netbox_circuit.test", "id"),
				),
			},
			{
				PlanOnly: true,
				Config: testAccCircuitGroupAssignmentResourceConfigLiteralNames(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
			},
		},
	})
}

func testAccCircuitGroupAssignmentResourceConfigLiteralNames(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid string) string {
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
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_circuit_group_assignment" "test" {
	group_id   = netbox_circuit_group.test.id
	circuit_id = netbox_circuit.test.id
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid)
}

func TestAccCircuitGroupAssignmentResource_externalDeletion(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-group-ext-del")
	groupSlug := testutil.RandomSlug("tf-test-cga-grp-ext-del")
	providerName := testutil.RandomName("tf-test-cga-provider-ext-del")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-ext-del")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-ext-del")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-ext-del")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerName)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentResourceConfigLiteralNames(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List assignments filtered by circuit CID
					items, _, err := client.CircuitsAPI.CircuitsCircuitGroupAssignmentsList(context.Background()).Circuit([]string{circuitCid}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit group assignment for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitGroupAssignmentsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit group assignment: %v", err)
					}
					t.Logf("Successfully externally deleted circuit group assignment with ID: %d", itemID)
				},

				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccCircuitGroupAssignmentResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccCircuitGroupAssignmentResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	groupName := testutil.RandomName("tf-test-cga-rem")
	groupSlug := testutil.RandomSlug("tf-test-cga-rem")
	providerName := testutil.RandomName("tf-test-cga-prov-rem")
	providerSlug := testutil.RandomSlug("tf-test-cga-prov-rem")
	circuitTypeName := testutil.RandomName("tf-test-cga-type-rem")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-type-rem")
	circuitCid := testutil.RandomSlug("tf-test-cga-ckt-rem")

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
				Config: testAccCircuitGroupAssignmentResourceConfig_withPriority(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, "primary"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "circuit_id"),
					resource.TestCheckResourceAttr("netbox_circuit_group_assignment.test", "priority", "primary"),
				),
			},
			{
				Config: testAccCircuitGroupAssignmentResourceConfig_basic(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "circuit_id"),
					resource.TestCheckNoResourceAttr("netbox_circuit_group_assignment.test", "priority"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_circuit_group_assignment",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_group_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_group_assignment" "test" {
  # group_id missing
  circuit_id = "1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_circuit_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_group_assignment" "test" {
  group_id = "1"
  # circuit_id missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
