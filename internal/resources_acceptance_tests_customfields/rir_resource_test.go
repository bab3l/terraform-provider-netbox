//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRIRResource_CustomFieldsPreservation tests that custom fields are preserved

// when updating other fields on a RIR.

//

// Filter-to-owned pattern:

// - Custom fields declared in config are managed by Terraform

// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform

func TestAccRIRResource_CustomFieldsPreservation(t *testing.T) {

	rirName := "RIR-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	rirSlug := "rir-" + strings.ToLower(acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum))

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create RIR WITH custom fields

				Config: testAccRIRConfig_preservation_step1(rirName, rirSlug, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "name", rirName),

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfOwner, "text", "admin"),
				),
			},

			{

				// Step 2: Update description WITHOUT mentioning custom_fields

				Config: testAccRIRConfig_preservation_step2(rirName, rirSlug, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "description", "Updated RIR"),

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "0"),
				),
			},

			{

				// Step 3: Import to verify custom fields still exist in NetBox

				ResourceName: "netbox_rir.test",

				ImportState:     true,
				ImportStateKind: resource.ImportCommandWithID,

				ImportStateVerify: false,

				ImportStateVerifyIgnore: []string{"custom_fields"},
			},

			{

				// Step 4: Add custom_fields back to verify they were preserved

				Config: testAccRIRConfig_preservation_step3(rirName, rirSlug, cfEnvironment, cfOwner),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfEnvironment, "text", "production"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfOwner, "text", "admin"),

					resource.TestCheckResourceAttr("netbox_rir.test", "description", "Updated RIR"),
				),
			},
		},
	})

}

// TestAccRIRResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern

func TestAccRIRResource_CustomFieldsFilterToOwned(t *testing.T) {

	cfEnv := testutil.RandomCustomFieldName("tf_env")

	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cfType := testutil.RandomCustomFieldName("tf_type")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Step 1: Create with two fields

				Config: testAccRIRConfig_filter_step1(cfEnv, cfOwner, cfType),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfEnv, "text", "prod"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfOwner, "text", "admin"),
				),
			},

			{

				// Step 2: Remove owner, keep env with updated value

				Config: testAccRIRConfig_filter_step2(cfEnv, cfOwner, cfType),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfEnv, "text", "staging"),
				),
			},

			{

				// Step 3: Add type

				Config: testAccRIRConfig_filter_step3(cfEnv, cfOwner, cfType),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "2"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfType, "text", "regional"),
				),
			},

			{

				// Step 4: Add owner back - should have preserved value

				Config: testAccRIRConfig_filter_step4(cfEnv, cfOwner, cfType),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "3"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfEnv, "text", "staging"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfOwner, "text", "admin"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfType, "text", "regional"),
				),
			},
		},
	})

}

// TestAccRIRResource_importWithCustomFields tests importing a RIR with custom fields

func TestAccRIRResource_importWithCustomFields(t *testing.T) {

	cfField := testutil.RandomCustomFieldName("tf_field")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRIRConfig_importTest(cfField),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "name", "Import Test RIR"),

					resource.TestCheckResourceAttr("netbox_rir.test", "custom_fields.#", "1"),

					testutil.CheckCustomFieldValue("netbox_rir.test", cfField, "text", "test-value"),
				),
			},

			{

				Config: testAccRIRConfig_importTest(cfField),

				PlanOnly: true,
			},

			{

				ResourceName: "netbox_rir.test",

				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,

				ImportStateVerify: false,

				ImportStateVerifyIgnore: []string{"custom_fields"},
			},

			{

				Config: testAccRIRConfig_importTest(cfField),

				PlanOnly: true,
			},
		},
	})

}

// =============================================================================

// Helper Config Functions - Preservation Tests

// =============================================================================

func testAccRIRConfig_preservation_step1(name, slug, cfEnv, cfOwner string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name = %[1]q

  slug = %[2]q



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "admin"

    }

  ]

}

`, name, slug, cfEnv, cfOwner)

}

func testAccRIRConfig_preservation_step2(name, slug, cfEnv, cfOwner string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name        = %[1]q

  slug        = %[2]q

  description = "Updated RIR"

}

`, name, slug, cfEnv, cfOwner)

}

func testAccRIRConfig_preservation_step3(name, slug, cfEnv, cfOwner string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "environment" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[4]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name        = %[1]q

  slug        = %[2]q

  description = "Updated RIR"



  custom_fields = [

    {

      name  = netbox_custom_field.environment.name

      type  = "text"

      value = "production"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "admin"

    }

  ]

}

`, name, slug, cfEnv, cfOwner)

}

// =============================================================================

// Helper Config Functions - Filter-to-Owned Tests

// =============================================================================

func testAccRIRConfig_filter_step1(cfEnv, cfOwner, cfType string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "rir_type" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name = "Filter Test RIR"

  slug = "filter-test-rir"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "prod"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "admin"

    }

  ]

}

`, cfEnv, cfOwner, cfType)

}

func testAccRIRConfig_filter_step2(cfEnv, cfOwner, cfType string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "rir_type" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name = "Filter Test RIR"

  slug = "filter-test-rir"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    }

  ]

}

`, cfEnv, cfOwner, cfType)

}

func testAccRIRConfig_filter_step3(cfEnv, cfOwner, cfType string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "rir_type" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name = "Filter Test RIR"

  slug = "filter-test-rir"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.rir_type.name

      type  = "text"

      value = "regional"

    }

  ]

}

`, cfEnv, cfOwner, cfType)

}

func testAccRIRConfig_filter_step4(cfEnv, cfOwner, cfType string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "env" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "owner" {

  name         = %[2]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_custom_field" "rir_type" {

  name         = %[3]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name = "Filter Test RIR"

  slug = "filter-test-rir"



  custom_fields = [

    {

      name  = netbox_custom_field.env.name

      type  = "text"

      value = "staging"

    },

    {

      name  = netbox_custom_field.owner.name

      type  = "text"

      value = "admin"

    },

    {

      name  = netbox_custom_field.rir_type.name

      type  = "text"

      value = "regional"

    }

  ]

}

`, cfEnv, cfOwner, cfType)

}

// =============================================================================

// Helper Config Functions - Import Test

// =============================================================================

func testAccRIRConfig_importTest(cfField string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "test_field" {

  name         = %[1]q

  type         = "text"

  object_types = ["ipam.rir"]

}



resource "netbox_rir" "test" {

  name = "Import Test RIR"

  slug = "import-test-rir"



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
