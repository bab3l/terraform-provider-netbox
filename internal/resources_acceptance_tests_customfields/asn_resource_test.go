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

// TestAccASNResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on an ASN.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform

func TestAccASNResource_CustomFieldsPreservation(t *testing.T) {
	asn := int64(acctest.RandIntRange(64512, 65534))

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create ASN WITH custom fields
				Config: testAccASNConfig_preservation_step1(asn, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccASNConfig_preservation_step2(asn, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:      "netbox_asn.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportCommandWithID,
				ImportStateVerify: false,
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccASNConfig_preservation_step3(asn, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", "Updated description"),
				),
			},
		},
	})
}

// TestAccASNResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
func TestAccASNResource_CustomFieldsFilterToOwned(t *testing.T) {
	asn := int64(acctest.RandIntRange(64512, 65534))

	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two fields
				Config: testAccASNConfig_filter_step1(asn, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Remove owner, keep env with updated value
				Config: testAccASNConfig_filter_step2(asn, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add cost_center
				Config: testAccASNConfig_filter_step3(asn, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfCostCenter, "text", "CC123"),
				),
			},
			{
				// Step 4: Add owner back - should have preserved value
				Config: testAccASNConfig_filter_step4(asn, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfOwner, "text", "team-a"),
					testutil.CheckCustomFieldValue("netbox_asn.test", cfCostCenter, "text", "CC123"),
				),
			},
		},
	})
}

// =============================================================================
// Helper Config Functions - Preservation Tests
// =============================================================================

func testAccASNConfig_preservation_step1(asn int64, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "environment" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.slug

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
`, asn, cfEnv, cfOwner)
}

func testAccASNConfig_preservation_step2(asn int64, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "environment" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn         = %[1]d
  rir         = netbox_rir.test.slug
  description = "Updated description"
}
`, asn, cfEnv, cfOwner)
}

func testAccASNConfig_preservation_step3(asn int64, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "environment" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn         = %[1]d
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
`, asn, cfEnv, cfOwner)
}

// =============================================================================
// Helper Config Functions - Filter-to-Owned Tests
// =============================================================================

func testAccASNConfig_filter_step1(asn int64, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.slug

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
`, asn, cfEnv, cfOwner, cfCost)
}

func testAccASNConfig_filter_step2(asn int64, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.slug

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]
}
`, asn, cfEnv, cfOwner, cfCost)
}

func testAccASNConfig_filter_step3(asn int64, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.slug

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
`, asn, cfEnv, cfOwner, cfCost)
}

func testAccASNConfig_filter_step4(asn int64, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "ARIN"
  slug = "arin"
}

resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.asn"]
}

resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.slug

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
`, asn, cfEnv, cfOwner, cfCost)
}

func TestAccASNResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	asn := int64(acctest.RandIntRange(64512, 65534))
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceImportConfig_full(asn, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_asn.test", "tags.#", "2"),
				),
			},
			{
				Config:            testAccASNResourceImportConfig_full(asn, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:      "netbox_asn.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccASNResourceImportConfig_full(asn, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccASNResourceImportConfig_full(asn int64, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["ipam.asn"]
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
resource "netbox_asn" "test" {
  asn    = %d
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

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
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
		asn,
	)
}
