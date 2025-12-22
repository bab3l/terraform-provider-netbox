package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupAssignmentResourceConfig_basic(

					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),
			},

			{

				ResourceName: "netbox_circuit_group_assignment.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"group_id", "circuit_id"},
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

  circuit_provider = netbox_provider.test.slug

  type             = netbox_circuit_type.test.slug

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

  circuit_provider = netbox_provider.test.slug

  type             = netbox_circuit_type.test.slug

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupAssignmentResourceConfigLiteralNames(

					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_group_assignment.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_group_assignment.test", "group_id", groupName),

					resource.TestCheckResourceAttr("netbox_circuit_group_assignment.test", "circuit_id", circuitCid),
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

  circuit_provider = netbox_provider.test.slug

  type             = netbox_circuit_type.test.slug

}

resource "netbox_circuit_group_assignment" "test" {

  group_id   = netbox_circuit_group.test.name

  circuit_id = netbox_circuit.test.cid

}

`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid)

}
