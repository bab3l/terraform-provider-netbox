//go:build customfields
// +build customfields

// Package resources_acceptance_tests_customfields contains acceptance tests for custom fields

// that require dedicated test runs to avoid conflicts with global custom field definitions.

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIPAddressResource_CustomFieldsPreservation tests that custom fields are preserved

// when updating other fields on an IP address.

//

// Filter-to-owned pattern:

// - Custom fields declared in config are managed by Terraform

// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform

// - Empty list explicitly clears all fields

func TestAccIPAddressResource_CustomFieldsPreservation(t *testing.T) {

	// Generate unique names

	ipAddress := "10.0.0.1/24"

	vrfName := testutil.RandomName("tf-test-vrf")

	vrfRD := testutil.RandomName("65000:100")

	// Custom field names

	cfEnvironment := testutil.RandomCustomFieldName("tf_environment")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create IP address WITH custom fields

				Config: testAccIPAddressConfig_withCustomFields(

					ipAddress, vrfName, vrfRD,

					cfEnvironment, cfOwner,

					"production", "team-a",
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ipAddress),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfOwner, "text", "team-a"),
				),
			},

			{

				// Step 2: Update description WITHOUT mentioning custom_fields in config

				// Custom fields should be preserved in NetBox (verified by import)

				Config: testAccIPAddressConfig_withoutCustomFields(

					ipAddress, vrfName, vrfRD,

					cfEnvironment, cfOwner,

					"Updated description",
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ipAddress),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Updated description"),

					// State shows 0 custom_fields (not in config = not owned)

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "0"),
				),
			},

			{

				// Step 3: Import to verify custom fields still exist in NetBox

				ResourceName: "netbox_ip_address.test",

				ImportState:       true,
				ImportStateKind:   resource.ImportCommandWithID,
				ImportStateVerify: false, // Can't verify - config has no custom_fields

			},

			{

				// Step 3a: Verify no drift after import

				Config: testAccIPAddressConfig_withoutCustomFields(

					ipAddress, vrfName, vrfRD,

					cfEnvironment, cfOwner,

					"Updated description",
				),

				PlanOnly: true,
			},

			{

				// Step 4: Add custom_fields back to config to verify they were preserved

				Config: testAccIPAddressConfig_withCustomFieldsAndDescription(

					ipAddress, vrfName, vrfRD,

					cfEnvironment, cfOwner,

					"production", "team-a",

					"Updated description", // Keep description from step 2

				),

				Check: resource.ComposeTestCheckFunc(

					// Custom fields should have their original values (preserved in NetBox)

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfOwner, "text", "team-a"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Updated description"),
				),
			},
		},
	})

}

// TestAccIPAddressResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern:

// - Only fields declared in config appear in state

// - Unowned fields are preserved in NetBox but invisible to Terraform

func TestAccIPAddressResource_CustomFieldsFilterToOwned(t *testing.T) {

	// Generate unique names

	ipAddress := "10.0.0.2/24"

	vrfName := testutil.RandomName("tf-test-vrf")

	vrfRD := testutil.RandomName("65000:200")

	// Custom field names

	cfEnv := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cfCostCenter := testutil.RandomCustomFieldName("tf_cost_center")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create with two custom fields (env and owner)

				Config: testAccIPAddressConfig_filterOwned_step1(

					ipAddress, vrfName, vrfRD,

					cfEnv, cfOwner, cfCostCenter,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfEnv, "text", "prod"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfOwner, "text", "team-a"),
				),
			},

			{

				// Step 2: Remove owner from config, keep only env

				// State should show only 1 field (env with updated value)

				Config: testAccIPAddressConfig_filterOwned_step2(

					ipAddress, vrfName, vrfRD,

					cfEnv, cfOwner, cfCostCenter,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfEnv, "text", "staging"),
				),
			},

			{

				// Step 3: Add cost_center to config, keep env

				// State should show 2 fields (env and cost_center)

				Config: testAccIPAddressConfig_filterOwned_step3(

					ipAddress, vrfName, vrfRD,

					cfEnv, cfOwner, cfCostCenter,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfCostCenter, "text", "123"),
				),
			},

			{

				// Step 4: Add owner back to config - should have preserved value

				Config: testAccIPAddressConfig_filterOwned_step4(

					ipAddress, vrfName, vrfRD,

					cfEnv, cfOwner, cfCostCenter,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "3"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfOwner, "text", "team-a"), // Preserved!

					testutil.CheckCustomFieldValue("netbox_ip_address.test", cfCostCenter, "text", "123"),
				),
			},
		},
	})

}

// TestAccIPAddressResource_ImportWithCustomFields tests import behavior

func TestAccIPAddressResource_ImportWithCustomFields(t *testing.T) {

	// Generate unique names

	ipAddress := "10.0.0.3/24"

	vrfName := testutil.RandomName("tf-test-vrf")

	vrfRD := testutil.RandomName("65000:300")

	// Custom field names

	cfText := testutil.RandomCustomFieldName("tf_text")

	cfTextValue := testutil.RandomName("text-value")

	cfInteger := testutil.RandomCustomFieldName("tf_integer")

	cfIntegerValue := 12345

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create the IP address with custom fields

				Config: testAccIPAddressConfig_import(

					ipAddress, vrfName, vrfRD,

					cfText, cfTextValue, cfInteger, cfIntegerValue,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ipAddress),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "custom_fields.#", "2"),
				),
			},

			{

				// Step 2: Import the IP address

				ResourceName: "netbox_ip_address.test",

				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},

			{

				// Step 3: Verify no drift after import

				Config: testAccIPAddressConfig_import(

					ipAddress, vrfName, vrfRD,

					cfText, cfTextValue, cfInteger, cfIntegerValue,
				),

				PlanOnly: true,
			},
		},
	})

}

// =============================================================================

// Helper Config Functions

// =============================================================================

func testAccIPAddressConfig_base(vrfName, vrfRD string, customFieldDefs string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[1]q

  rd   = %[2]q

}



%s

`, vrfName, vrfRD, customFieldDefs)

}

func testAccIPAddressConfig_withCustomFields(

	ipAddress, vrfName, vrfRD,

	cfEnvironment, cfOwner,

	cfEnvValue, cfOwnerValue string,

) string {

	return testAccIPAddressConfig_withCustomFieldsAndDescription(

		ipAddress, vrfName, vrfRD,

		cfEnvironment, cfOwner,

		cfEnvValue, cfOwnerValue,

		"", // No description

	)

}

func testAccIPAddressConfig_withCustomFieldsAndDescription(

	ipAddress, vrfName, vrfRD,

	cfEnvironment, cfOwner,

	cfEnvValue, cfOwnerValue, description string,

) string {

	customFieldDefs := fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}

`, cfEnvironment, cfOwner)

	base := testAccIPAddressConfig_base(vrfName, vrfRD, customFieldDefs)

	descriptionField := ""

	if description != "" {

		descriptionField = fmt.Sprintf("\n  description = %q", description)

	}

	return base + fmt.Sprintf(`

resource "netbox_ip_address" "test" {

  address = %[1]q

  vrf     = netbox_vrf.test.id

  status  = "active"%[4]s



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = %[2]q

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = %[3]q

    }

  ]

}

`, ipAddress, cfEnvValue, cfOwnerValue, descriptionField)

}

func testAccIPAddressConfig_withoutCustomFields(

	ipAddress, vrfName, vrfRD,

	cfEnvironment, cfOwner,

	description string,

) string {

	customFieldDefs := fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}

`, cfEnvironment, cfOwner)

	base := testAccIPAddressConfig_base(vrfName, vrfRD, customFieldDefs)

	return base + fmt.Sprintf(`

resource "netbox_ip_address" "test" {

  address     = %[1]q

  vrf         = netbox_vrf.test.id

  status      = "active"

  description = %[2]q



  # custom_fields intentionally omitted - should be preserved in NetBox

  depends_on = [

    netbox_custom_field.environment,

    netbox_custom_field.owner

  ]

}

`, ipAddress, description)

}

// Filter-to-owned test configs

func testAccIPAddressConfig_filterOwned_step1(

	ipAddress, vrfName, vrfRD,

	cfEnv, cfOwner, cfCostCenter string,

) string {

	customFieldDefs := fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "cost_center" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}

`, cfEnv, cfOwner, cfCostCenter)

	base := testAccIPAddressConfig_base(vrfName, vrfRD, customFieldDefs)

	return base + fmt.Sprintf(`

resource "netbox_ip_address" "test" {

  address = %[1]q

  vrf     = netbox_vrf.test.id

  status  = "active"



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

`, ipAddress)

}

func testAccIPAddressConfig_filterOwned_step2(

	ipAddress, vrfName, vrfRD,

	cfEnv, cfOwner, cfCostCenter string,

) string {

	customFieldDefs := fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "cost_center" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}

`, cfEnv, cfOwner, cfCostCenter)

	base := testAccIPAddressConfig_base(vrfName, vrfRD, customFieldDefs)

	return base + fmt.Sprintf(`

resource "netbox_ip_address" "test" {

  address = %[1]q

  vrf     = netbox_vrf.test.id

  status  = "active"



  # owner removed from config - should be preserved in NetBox

  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    }

  ]



  depends_on = [netbox_custom_field.owner]

}

`, ipAddress)

}

func testAccIPAddressConfig_filterOwned_step3(

	ipAddress, vrfName, vrfRD,

	cfEnv, cfOwner, cfCostCenter string,

) string {

	customFieldDefs := fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "cost_center" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}

`, cfEnv, cfOwner, cfCostCenter)

	base := testAccIPAddressConfig_base(vrfName, vrfRD, customFieldDefs)

	return base + fmt.Sprintf(`

resource "netbox_ip_address" "test" {

  address = %[1]q

  vrf     = netbox_vrf.test.id

  status  = "active"



  # owner still not in config, but add cost_center

  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.cost_center.name

      type  = "text"

      value = "123"

    }

  ]



  depends_on = [netbox_custom_field.owner]

}

`, ipAddress)

}

func testAccIPAddressConfig_filterOwned_step4(

	ipAddress, vrfName, vrfRD,

	cfEnv, cfOwner, cfCostCenter string,

) string {

	customFieldDefs := fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}



resource "netbox_custom_field" "cost_center" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

}

`, cfEnv, cfOwner, cfCostCenter)

	base := testAccIPAddressConfig_base(vrfName, vrfRD, customFieldDefs)

	return base + fmt.Sprintf(`

resource "netbox_ip_address" "test" {

  address = %[1]q

  vrf     = netbox_vrf.test.id

  status  = "active"



  # Add owner back - should have preserved value from step 1

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

      name  = netbox_custom_field.cost_center.name

      type  = "text"

      value = "123"

    }

  ]

}

`, ipAddress)

}

// Import test config

func testAccIPAddressConfig_import(

	ipAddress, vrfName, vrfRD,

	cfText, cfTextValue, cfInteger string, cfIntegerValue int,

) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = %[2]q

  rd   = %[3]q

}



resource "netbox_custom_field" "text" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.ipaddress"]

  required     = false

}



resource "netbox_custom_field" "integer" {

  name         = %[5]q

  type         = "integer"

  object_types = ["ipam.ipaddress"]

  required     = false

}



resource "netbox_ip_address" "test" {

  address = %[1]q

  vrf     = netbox_vrf.test.id

  status  = "active"



  custom_fields = [

    {

      name  = netbox_custom_field.text.name

      type  = "text"

      value = %[6]q

    },

    {

      name  = netbox_custom_field.integer.name

      type  = "integer"

      value = "%[7]d"

    }

  ]

}

`, ipAddress, vrfName, vrfRD, cfText, cfInteger, cfTextValue, cfIntegerValue)

}
