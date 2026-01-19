//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccServiceResource_CustomFieldsPreservation tests that custom fields are preserved

// when updating other fields on a Service.

//

// Filter-to-owned pattern:

func TestAccServiceResource_CustomFieldsPreservation(t *testing.T) {

	serviceName := "svc-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create Service WITH custom fields

				Config: testAccServiceConfig_preservation_step1(serviceName, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfOwner, "text", "platform-team"),
				),
			},

			{

				// Step 2: Update description WITHOUT mentioning custom_fields

				Config: testAccServiceConfig_preservation_step2(serviceName, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "description", "Updated service"),

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "0"),
				),
			},

			{

				// Step 3: Import to verify custom fields still exist in NetBox

				ResourceName: "netbox_service.test",

				ImportState: true,

				ImportStateKind: resource.ImportCommandWithID,

				ImportStateVerify: false,

				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},

			{

				// Step 4: Add custom_fields back to verify they were preserved

				Config: testAccServiceConfig_preservation_step3(serviceName, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfOwner, "text", "platform-team"),

					resource.TestCheckResourceAttr("netbox_service.test", "description", "Updated service"),
				),
			},
		},
	})

}

// TestAccServiceResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern

func TestAccServiceResource_CustomFieldsFilterToOwned(t *testing.T) {

	serviceName := "svc-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	cfEnv := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create with two fields

				Config: testAccServiceConfig_filter_step1(serviceName, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfEnv, "text", "prod"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfOwner, "text", "platform-team"),
				),
			},

			{

				// Step 2: Remove owner, keep env with updated value

				Config: testAccServiceConfig_filter_step2(serviceName, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfEnv, "text", "staging"),
				),
			},

			{

				// Step 3: Add cost_center

				Config: testAccServiceConfig_filter_step3(serviceName, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfCostCenter, "text", "ENG-001"),
				),
			},

			{

				// Step 4: Add owner back - should have preserved value

				Config: testAccServiceConfig_filter_step4(serviceName, cfEnv, cfOwner, cfCostCenter),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "3"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfOwner, "text", "platform-team"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfCostCenter, "text", "ENG-001"),
				),
			},
		},
	})

}

// TestAccServiceResource_importWithCustomFields tests importing a Service with custom fields

func TestAccServiceResource_importWithCustomFields(t *testing.T) {

	cfField := testutil.RandomCustomFieldName("tf_field")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccServiceConfig_importTest(cfField),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_service.test", "name", "http"),

					resource.TestCheckResourceAttr("netbox_service.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_service.test", cfField, "text", "test-value"),
				),
			},

			{

				Config:            testAccServiceConfig_importTest(cfField),
				ResourceName:      "netbox_service.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},

			{

				Config: testAccServiceConfig_importTest(cfField),

				PlanOnly: true,
			},
		},
	})

}

// =============================================================================

// Helper Config Functions

// =============================================================================

func testAccServiceConfig_preservation_step2_helper() string {

	return `

resource "netbox_site" "test" {

  name = "test-site"

  slug = "test-site"

}



resource "netbox_manufacturer" "test" {

  name = "test-mfg"

  slug = "test-mfg"

}



resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = "test-model"

  slug         = "test-model"

}



resource "netbox_device_role" "test" {

  name  = "test-role"

  slug  = "test-role"

  color = "aa1409"

}



resource "netbox_device" "test" {

  name        = "test-device"

  device_type = netbox_device_type.test.id

  role        = netbox_device_role.test.id

  site        = netbox_site.test.id

}

`

}

// =============================================================================

// Test Config Functions - Preservation Tests

// =============================================================================

func testAccServiceConfig_preservation_step1(serviceName, cfEnv, cfOwner string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device   = netbox_device.test.name

  name     = %[1]q

  protocol = "tcp"

  ports    = [8080]



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "platform-team"

    }

  ]

}

`, serviceName, cfEnv, cfOwner)

}

func testAccServiceConfig_preservation_step2(serviceName, cfEnv, cfOwner string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device      = netbox_device.test.id

  name        = %[1]q

  protocol    = "tcp"

  description = "Updated service"

  ports       = [8080]

}

`, serviceName, cfEnv, cfOwner)

}

func testAccServiceConfig_preservation_step3(serviceName, cfEnv, cfOwner string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device      = netbox_device.test.id

  name        = %[1]q

  protocol    = "tcp"

  description = "Updated service"

  ports       = [8080]



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "platform-team"

    }

  ]

}

`, serviceName, cfEnv, cfOwner)

}

// =============================================================================

// Test Config Functions - Filter-to-Owned Tests

// =============================================================================

func testAccServiceConfig_filter_step1(serviceName, cfEnv, cfOwner, cfCost string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "cost" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device   = netbox_device.test.name

  name     = %[1]q

  protocol = "tcp"

  ports    = [8080]



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "prod"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "platform-team"

    }

  ]

}

`, serviceName, cfEnv, cfOwner, cfCost)

}

func testAccServiceConfig_filter_step2(serviceName, cfEnv, cfOwner, cfCost string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "cost" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device   = netbox_device.test.name

  name     = %[1]q

  protocol = "tcp"

  ports    = [8080]



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    }

  ]

}

`, serviceName, cfEnv, cfOwner, cfCost)

}

func testAccServiceConfig_filter_step3(serviceName, cfEnv, cfOwner, cfCost string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "cost" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device   = netbox_device.test.id

  name     = %[1]q

  protocol = "tcp"

  ports    = [8080]



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.cost.name

      type  = "text"

      value = "ENG-001"

    }

  ]

}

`, serviceName, cfEnv, cfOwner, cfCost)

}

func testAccServiceConfig_filter_step4(serviceName, cfEnv, cfOwner, cfCost string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "owner" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_custom_field" "cost" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device   = netbox_device.test.id

  name     = %[1]q

  protocol = "tcp"

  ports    = [8080]



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "platform-team"

    },

    {

      name  = netbox_custom_field.cost.name

      type  = "text"

      value = "ENG-001"

    }

  ]

}

`, serviceName, cfEnv, cfOwner, cfCost)

}

// =============================================================================

// Test Config Functions - Import Test

// =============================================================================

func testAccServiceConfig_importTest(cfField string) string {

	return testAccServiceConfig_preservation_step2_helper() + fmt.Sprintf(`

resource "netbox_custom_field" "test_field" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.service"]

}



resource "netbox_service" "test" {

  device   = netbox_device.test.name

  name     = "http"

  protocol = "tcp"

  ports    = [80]



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
