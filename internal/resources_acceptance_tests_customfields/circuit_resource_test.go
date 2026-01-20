//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	cid := testutil.RandomName("tf-test-circuit")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	providerName := providerSlug
	typeSlug := testutil.RandomSlug("tf-test-circuit-type")
	typeName := typeSlug
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate test data for all custom field types
	textValue := testutil.RandomName("text-value")
	longtextValue := testutil.RandomName("longtext-value") + "\nThis is a multiline text field for comprehensive testing."
	intValue := 42 // Fixed value for reproducibility
	boolValue := true
	dateValue := testutil.RandomDate()
	urlValue := testutil.RandomURL("test-url")
	jsonValue := testutil.RandomJSON()

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfJSON := testutil.RandomCustomFieldName("tf_json")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(typeSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitResourceImportConfig_full(
					cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
				),
			},
			{
				ResourceName:      "netbox_circuit.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			// Enhancement 1: Verify no changes after import
			{
				Config: testAccCircuitResourceImportConfig_full(
					cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccCircuitResourceImportConfig_full(
	cid, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug string,
	textValue, longtextValue string, intValue int, boolValue bool, dateValue, urlValue, jsonValue string,
	tag1, tag1Slug, tag2, tag2Slug string,
	cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON string,
) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Custom Fields for circuits.circuit object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["circuits.circuit"]
}

# Circuit with comprehensive custom fields and tags
resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
  tenant           = netbox_tenant.test.id

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = %q
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "%d"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "%t"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = %q
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = %q
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = %q
    },
  ]

  depends_on = [
    netbox_custom_field.test_text,
    netbox_custom_field.test_longtext,
    netbox_custom_field.test_integer,
    netbox_custom_field.test_boolean,
    netbox_custom_field.test_date,
    netbox_custom_field.test_url,
    netbox_custom_field.test_json,
  ]
}
`, providerName, providerSlug, typeName, typeSlug, tenantName, tenantSlug,
		tag1, tag1Slug, tag2, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		cid, textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
}

// TestAccCircuitResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a circuit.
func TestAccCircuitResource_CustomFieldsPreservation(t *testing.T) {
	cid := testutil.RandomName("tf-test-circuit")
	providerSlug := testutil.RandomSlug("tf-provider")
	typeSlug := testutil.RandomSlug("tf-type")

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create circuit WITH custom fields
				Config: testAccCircuitConfig_preservation_step1(cid, providerSlug, typeSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "cid", cid),
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccCircuitConfig_preservation_step2(cid, providerSlug, typeSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist
				ResourceName:            "netbox_circuit.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"circuit_provider", "type", "custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to verify preservation
				Config: testAccCircuitConfig_preservation_step3(cid, providerSlug, typeSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_circuit.test", "description", "Updated description"),
				),
			},
		},
	})
}

// TestAccCircuitResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
func TestAccCircuitResource_CustomFieldsFilterToOwned(t *testing.T) {
	cid := testutil.RandomName("tf-test-circuit")
	providerSlug := testutil.RandomSlug("tf-provider")
	typeSlug := testutil.RandomSlug("tf-type")

	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two fields
				Config: testAccCircuitConfig_filter_step1(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Remove owner, keep env with updated value
				Config: testAccCircuitConfig_filter_step2(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add cost_center
				Config: testAccCircuitConfig_filter_step3(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfCostCenter, "text", "CC123"),
				),
			},
			{
				// Step 4: Add owner back - should have preserved value
				Config: testAccCircuitConfig_filter_step4(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfOwner, "text", "team-a"),
					testutil.CheckCustomFieldValue("netbox_circuit.test", cfCostCenter, "text", "CC123"),
				),
			},
		},
	})
}

// Helper config functions for preservation tests
func testAccCircuitConfig_preservation_step1(cid, providerSlug, typeSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "environment" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id

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
`, cid, providerSlug, typeSlug, cfEnv, cfOwner)
}

func testAccCircuitConfig_preservation_step2(cid, providerSlug, typeSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "environment" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  description      = "Updated description"
}
`, cid, providerSlug, typeSlug, cfEnv, cfOwner)
}

func testAccCircuitConfig_preservation_step3(cid, providerSlug, typeSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "environment" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  description      = "Updated description"

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
`, cid, providerSlug, typeSlug, cfEnv, cfOwner)
}

// Helper config functions for filter-to-owned tests
func testAccCircuitConfig_filter_step1(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id

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
`, cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost)
}

func testAccCircuitConfig_filter_step2(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]
}
`, cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost)
}

func testAccCircuitConfig_filter_step3(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id

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
`, cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost)
}

func testAccCircuitConfig_filter_step4(cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["circuits.circuit"]
}

resource "netbox_circuit" "test" {
  cid              = %[1]q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id

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
`, cid, providerSlug, typeSlug, cfEnv, cfOwner, cfCost)
}
