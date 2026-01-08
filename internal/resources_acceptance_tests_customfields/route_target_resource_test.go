//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRouteTargetResource_CustomFieldsPreservation tests that custom fields are preserved

// when updating other fields on a Route Target.

//

// Filter-to-owned pattern:

// - Custom fields declared in config are managed by Terraform

// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform

func TestAccRouteTargetResource_CustomFieldsPreservation(t *testing.T) {

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")

	cfDescription := testutil.RandomCustomFieldName("tf_desc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create Route Target WITH custom fields

				Config: testAccRouteTargetConfig_preservation_step1(cfEnvironment, cfDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "name", "65000:100"),

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfDescription, "text", "main-vrf"),
				),
			},

			{

				// Step 2: Update name WITHOUT mentioning custom_fields

				Config: testAccRouteTargetConfig_preservation_step2(cfEnvironment, cfDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "name", "65000:101"),

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "0"),
				),
			},

			{

				// Step 3: Import to verify custom fields still exist in NetBox

				ResourceName: "netbox_route_target.test",

				ImportState: true,

				ImportStateVerify: false,

				ImportStateVerifyIgnore: []string{"custom_fields"},
			},

			{

				// Step 4: Add custom_fields back to verify they were preserved

				Config: testAccRouteTargetConfig_preservation_step3(cfEnvironment, cfDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfDescription, "text", "main-vrf"),

					resource.TestCheckResourceAttr("netbox_route_target.test", "name", "65000:101"),
				),
			},
		},
	})

}

// TestAccRouteTargetResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern

func TestAccRouteTargetResource_CustomFieldsFilterToOwned(t *testing.T) {

	cfEnv := testutil.RandomCustomFieldName("tf_env")

	cfDesc := testutil.RandomCustomFieldName("tf_desc")

	cfTeam := testutil.RandomCustomFieldName("tf_team")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create with two fields

				Config: testAccRouteTargetConfig_filter_step1(cfEnv, cfDesc, cfTeam),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfEnv, "text", "prod"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfDesc, "text", "route-target"),
				),
			},

			{

				// Step 2: Remove description, keep env with updated value

				Config: testAccRouteTargetConfig_filter_step2(cfEnv, cfDesc, cfTeam),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfEnv, "text", "staging"),
				),
			},

			{

				// Step 3: Add team

				Config: testAccRouteTargetConfig_filter_step3(cfEnv, cfDesc, cfTeam),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfTeam, "text", "network-ops"),
				),
			},

			{

				// Step 4: Add description back - should have preserved value

				Config: testAccRouteTargetConfig_filter_step4(cfEnv, cfDesc, cfTeam),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "3"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfDesc, "text", "route-target"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfTeam, "text", "network-ops"),
				),
			},
		},
	})

}

// TestAccRouteTargetResource_importWithCustomFields tests importing a Route Target with custom fields

func TestAccRouteTargetResource_importWithCustomFields(t *testing.T) {

	cfField := testutil.RandomCustomFieldName("tf_field")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRouteTargetConfig_importTest(cfField),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "name", "65000:200"),

					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_route_target.test", cfField, "text", "test-value"),
				),
			},

			{

				Config: testAccRouteTargetConfig_importTest(cfField),

				PlanOnly: true,
			},

			{

				ResourceName: "netbox_route_target.test",

				ImportState: true,

				ImportStateVerify: false,

				ImportStateVerifyIgnore: []string{"custom_fields"},
			},

			{

				Config: testAccRouteTargetConfig_importTest(cfField),

				PlanOnly: true,
			},
		},
	})

}

// =============================================================================

// Helper Config Functions - Preservation Tests

// =============================================================================

func testAccRouteTargetConfig_preservation_step1(cfEnv, cfDesc string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "description" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:100"



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.description.name

      type  = "text"

      value = "main-vrf"

    }

  ]

}

`, cfEnv, cfDesc)

}

func testAccRouteTargetConfig_preservation_step2(cfEnv, cfDesc string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "description" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:101"

}

`, cfEnv, cfDesc)

}

func testAccRouteTargetConfig_preservation_step3(cfEnv, cfDesc string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "description" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:101"



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.description.name

      type  = "text"

      value = "main-vrf"

    }

  ]

}

`, cfEnv, cfDesc)

}

// =============================================================================

// Helper Config Functions - Filter-to-Owned Tests

// =============================================================================

func testAccRouteTargetConfig_filter_step1(cfEnv, cfDesc, cfTeam string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "desc" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "team" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:210"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "prod"

    },

    {

      name  = netbox_custom_field.desc.name

      type  = "text"

      value = "route-target"

    }

  ]

}

`, cfEnv, cfDesc, cfTeam)

}

func testAccRouteTargetConfig_filter_step2(cfEnv, cfDesc, cfTeam string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "desc" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "team" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:210"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    }

  ]

}

`, cfEnv, cfDesc, cfTeam)

}

func testAccRouteTargetConfig_filter_step3(cfEnv, cfDesc, cfTeam string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "desc" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "team" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:210"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.team.name

      type  = "text"

      value = "network-ops"

    }

  ]

}

`, cfEnv, cfDesc, cfTeam)

}

func testAccRouteTargetConfig_filter_step4(cfEnv, cfDesc, cfTeam string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "desc" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_custom_field" "team" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:210"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.desc.name

      type  = "text"

      value = "route-target"

    },

    {

      name  = netbox_custom_field.team.name

      type  = "text"

      value = "network-ops"

    }

  ]

}

`, cfEnv, cfDesc, cfTeam)

}

// =============================================================================

// Helper Config Functions - Import Test

// =============================================================================

func testAccRouteTargetConfig_importTest(cfField string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "test_field" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.routetarget"]

}



resource "netbox_route_target" "test" {

  name = "65000:200"



  custom_fields = [

    {

      name  = netbox_custom_field.test_field.name

      type  = "text"

      value = "test-value"

    }

  ]

}

`, cfField)

}
