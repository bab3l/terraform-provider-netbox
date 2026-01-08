//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceImportConfig_full(prefix, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_aggregate.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccAggregateResourceImportConfig_full(prefix, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_aggregate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir", "custom_fields", "tags", "tenant"},
			},
			// Enhancement 1: Verify no changes after import
			{
				Config:   testAccAggregateResourceImportConfig_full(prefix, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

// TestAccAggregateResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on an aggregate.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccAggregateResource_CustomFieldsPreservation(t *testing.T) {
	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create aggregate WITH custom fields
				Config: testAccAggregateConfig_preservation_step1(prefix, rirName, rirSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccAggregateConfig_preservation_step2(prefix, rirName, rirSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_aggregate.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"rir", "custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccAggregateConfig_preservation_step3(prefix, rirName, rirSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", "Updated description"),
				),
			},
		},
	})
}

// TestAccAggregateResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
func TestAccAggregateResource_CustomFieldsFilterToOwned(t *testing.T) {
	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two fields
				Config: testAccAggregateConfig_filter_step1(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Remove owner, keep env with updated value
				Config: testAccAggregateConfig_filter_step2(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add cost_center
				Config: testAccAggregateConfig_filter_step3(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfCostCenter, "text", "CC123"),
				),
			},
			{
				// Step 4: Add owner back - should have preserved value
				Config: testAccAggregateConfig_filter_step4(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfOwner, "text", "team-a"),
					testutil.CheckCustomFieldValue("netbox_aggregate.test", cfCostCenter, "text", "CC123"),
				),
			},
		},
	})
}

// =============================================================================
// Helper Config Functions - Preservation Tests
// =============================================================================

func testAccAggregateConfig_preservation_step1(prefix, rirName, rirSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "environment" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.slug

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
`, prefix, rirName, rirSlug, cfEnv, cfOwner)
}

func testAccAggregateConfig_preservation_step2(prefix, rirName, rirSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "environment" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix      = %[1]q
  rir         = netbox_rir.test.slug
  description = "Updated description"
}
`, prefix, rirName, rirSlug, cfEnv, cfOwner)
}

func testAccAggregateConfig_preservation_step3(prefix, rirName, rirSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "environment" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix      = %[1]q
  rir         = netbox_rir.test.slug
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
`, prefix, rirName, rirSlug, cfEnv, cfOwner)
}

// =============================================================================
// Helper Config Functions - Filter-to-Owned Tests
// =============================================================================

func testAccAggregateConfig_filter_step1(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.slug

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
`, prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost)
}

func testAccAggregateConfig_filter_step2(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.slug

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]
}
`, prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost)
}

func testAccAggregateConfig_filter_step3(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.slug

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
`, prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost)
}

func testAccAggregateConfig_filter_step4(prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_custom_field" "env" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "owner" {
  name         = %[5]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cost" {
  name         = %[6]q
  type         = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.slug

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
`, prefix, rirName, rirSlug, cfEnv, cfOwner, cfCost)
}

func testAccAggregateResourceImportConfig_full(prefix, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["ipam.aggregate"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["ipam.aggregate"]
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

# Main Resource
resource "netbox_aggregate" "test" {
  prefix = %q
  rir    = netbox_rir.test.slug
  tenant = netbox_tenant.test.slug

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`,
		tenantName, tenantSlug,
		rirName, rirSlug,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		prefix,
	)
}
