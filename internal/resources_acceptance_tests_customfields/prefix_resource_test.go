//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccPrefixResource_CustomFieldsPreservation tests that custom fields are preserved

// when updating other fields on a prefix.

//

// Filter-to-owned pattern:

// - Custom fields declared in config are managed by Terraform

// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform

func TestAccPrefixResource_CustomFieldsPreservation(t *testing.T) {

	prefix := "10.0.0.0/24"

	vrfName := testutil.RandomName("tf-test-vrf")

	vrfRD := testutil.RandomName("65000:100")

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create prefix WITH custom fields

				Config: testAccPrefixConfig_preservation_step1(prefix, vrfName, vrfRD, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfOwner, "text", "team-a"),
				),
			},

			{

				// Step 2: Update description WITHOUT mentioning custom_fields

				Config: testAccPrefixConfig_preservation_step2(prefix, vrfName, vrfRD, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "0"),
				),
			},

			{

				// Step 3: Import to verify custom fields still exist in NetBox

				ResourceName: "netbox_prefix.test",

				ImportState:     true,
				ImportStateKind: resource.ImportCommandWithID,

				ImportStateVerify: false,

				ImportStateVerifyIgnore: []string{"vrf", "custom_fields"},
			},

			{

				// Step 4: Add custom_fields back to verify they were preserved

				Config: testAccPrefixConfig_preservation_step3(prefix, vrfName, vrfRD, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfOwner, "text", "team-a"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "description", "Updated description"),
				),
			},
		},
	})

}

// TestAccPrefixResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern

func TestAccPrefixResource_CustomFieldsFilterToOwned(t *testing.T) {

	prefix := "10.1.0.0/24"

	vrfName := testutil.RandomName("tf-test-vrf")

	vrfRD := testutil.RandomName("65000:200")

	cfEnv := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create with two fields

				Config: testAccPrefixConfig_filter_step1(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfEnv, "text", "prod"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfOwner, "text", "team-a"),
				),
			},

			{

				// Step 2: Remove owner, keep env with updated value

				Config: testAccPrefixConfig_filter_step2(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfEnv, "text", "staging"),
				),
			},

			{

				// Step 3: Add cost_center

				Config: testAccPrefixConfig_filter_step3(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfCostCenter, "text", "CC123"),
				),
			},

			{

				// Step 4: Add owner back - should have preserved value

				Config: testAccPrefixConfig_filter_step4(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "3"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfOwner, "text", "team-a"),

					testutil.CheckCustomFieldValue("netbox_prefix.test", cfCostCenter, "text", "CC123"),
				),
			},
		},
	})

}

// TestAccPrefixResource_ImportWithCustomFields tests import behavior

func TestAccPrefixResource_ImportWithCustomFields(t *testing.T) {

	prefix := "10.2.0.0/24"

	vrfName := testutil.RandomName("tf-test-vrf")

	vrfRD := testutil.RandomName("65000:300")

	cfText := testutil.RandomCustomFieldName("tf_text")

	cfInteger := testutil.RandomCustomFieldName("tf_integer")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create the prefix with custom fields

				Config: testAccPrefixConfig_import(prefix, vrfName, vrfRD, cfText, cfInteger),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_prefix.test", "custom_fields.#", "2"),
				),
			},

			{

				// Step 2: Import the prefix

				Config: testAccPrefixConfig_import(prefix, vrfName, vrfRD, cfText, cfInteger),

				ResourceName: "netbox_prefix.test",

				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,

				ImportStateVerify: false,
			},

			{

				// Step 3: Verify no drift after import

				Config: testAccPrefixConfig_import(prefix, vrfName, vrfRD, cfText, cfInteger),

				PlanOnly: true,
			},
		},
	})

}

// =============================================================================

// Helper Config Functions - Preservation Tests

// =============================================================================

func testAccPrefixConfig_preservation_step1(prefix, vrfName, vrfRD, cfEnv, cfOwner string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "environment" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix = %[1]q

  vrf    = netbox_vrf.test.id

  status = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "team-a"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner)

}

func testAccPrefixConfig_preservation_step2(prefix, vrfName, vrfRD, cfEnv, cfOwner string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "environment" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix      = %[1]q

  vrf         = netbox_vrf.test.id

  status      = "active"

  description = "Updated description"

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner)

}

func testAccPrefixConfig_preservation_step3(prefix, vrfName, vrfRD, cfEnv, cfOwner string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "environment" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix      = %[1]q

  vrf         = netbox_vrf.test.id

  status      = "active"

  description = "Updated description"



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "team-a"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner)

}

// =============================================================================

// Helper Config Functions - Filter-to-Owned Tests

// =============================================================================

func testAccPrefixConfig_filter_step1(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "env" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "cost" {

  name         = %[6]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix = %[1]q

  vrf    = netbox_vrf.test.id

  status = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "prod"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "team-a"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost)

}

func testAccPrefixConfig_filter_step2(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "env" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "cost" {

  name         = %[6]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix = %[1]q

  vrf    = netbox_vrf.test.id

  status = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost)

}

func testAccPrefixConfig_filter_step3(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "env" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "cost" {

  name         = %[6]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix = %[1]q

  vrf    = netbox_vrf.test.id

  status = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.cost.name

      type  = "text"

      value = "CC123"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost)

}

func testAccPrefixConfig_filter_step4(prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "env" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "owner" {

  name         = %[5]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "cost" {

  name         = %[6]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix = %[1]q

  vrf    = netbox_vrf.test.id

  status = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "team-a"

    },

    {

      name  = netbox_custom_field.cost.name

      type  = "text"

      value = "CC123"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfEnv, cfOwner, cfCost)

}

// =============================================================================

// Helper Config Functions - Import Tests

// =============================================================================

func testAccPrefixConfig_import(prefix, vrfName, vrfRD, cfText, cfInteger string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "text" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.prefix"]

}



resource "netbox_custom_field" "integer" {

  name         = %[5]q

  type         = "integer"

  object_types = ["ipam.prefix"]

}



resource "netbox_prefix" "test" {

  prefix = %[1]q

  vrf    = netbox_vrf.test.id

  status = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.text.name

      type  = "text"

      value = "test-value"

    },

    {

      name  = netbox_custom_field.integer.name

      type  = "integer"

      value = "42"

    }

  ]

}

`, prefix, vrfName, vrfRD, cfText, cfInteger)

}
