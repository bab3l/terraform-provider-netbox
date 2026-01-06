//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIPRangeResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on an IP Range.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccIPRangeResource_CustomFieldsPreservation(t *testing.T) {
	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfPurpose := testutil.RandomCustomFieldName("tf_purpose")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create IP Range WITH custom fields
				Config: testAccIPRangeConfig_preservation_step1(cfEnvironment, cfPurpose),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", "192.168.1.10"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", "192.168.1.20"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfPurpose, "text", "servers"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccIPRangeConfig_preservation_step2(cfEnvironment, cfPurpose),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", "Updated range"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_ip_range.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccIPRangeConfig_preservation_step3(cfEnvironment, cfPurpose),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfPurpose, "text", "servers"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", "Updated range"),
				),
			},
		},
	})
}

// TestAccIPRangeResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
func TestAccIPRangeResource_CustomFieldsFilterToOwned(t *testing.T) {
	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfPurpose := testutil.RandomCustomFieldName("tf_purpose")
	cfTeam := testutil.RandomCustomFieldName("tf_team")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two fields
				Config: testAccIPRangeConfig_filter_step1(cfEnv, cfPurpose, cfTeam),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfPurpose, "text", "servers"),
				),
			},
			{
				// Step 2: Remove purpose, keep env with updated value
				Config: testAccIPRangeConfig_filter_step2(cfEnv, cfPurpose, cfTeam),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add team
				Config: testAccIPRangeConfig_filter_step3(cfEnv, cfPurpose, cfTeam),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfTeam, "text", "network-team"),
				),
			},
			{
				// Step 4: Add purpose back - should have preserved value
				Config: testAccIPRangeConfig_filter_step4(cfEnv, cfPurpose, cfTeam),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfPurpose, "text", "servers"),
					testutil.CheckCustomFieldValue("netbox_ip_range.test", cfTeam, "text", "network-team"),
				),
			},
		},
	})
}

// =============================================================================
// Helper Config Functions - Preservation Tests
// =============================================================================

func testAccIPRangeConfig_preservation_step1(cfEnv, cfPurpose string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.purpose.name
      type  = "text"
      value = "servers"
    }
  ]
}
`, cfEnv, cfPurpose)
}

func testAccIPRangeConfig_preservation_step2(cfEnv, cfPurpose string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"
  description   = "Updated range"
}
`, cfEnv, cfPurpose)
}

func testAccIPRangeConfig_preservation_step3(cfEnv, cfPurpose string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"
  description   = "Updated range"

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.purpose.name
      type  = "text"
      value = "servers"
    }
  ]
}
`, cfEnv, cfPurpose)
}

// =============================================================================
// Helper Config Functions - Filter-to-Owned Tests
// =============================================================================

func testAccIPRangeConfig_filter_step1(cfEnv, cfPurpose, cfTeam string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "team" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "prod"
    },
    {
      name  = netbox_custom_field.purpose.name
      type  = "text"
      value = "servers"
    }
  ]
}
`, cfEnv, cfPurpose, cfTeam)
}

func testAccIPRangeConfig_filter_step2(cfEnv, cfPurpose, cfTeam string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "team" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]
}
`, cfEnv, cfPurpose, cfTeam)
}

func testAccIPRangeConfig_filter_step3(cfEnv, cfPurpose, cfTeam string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "team" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.team.name
      type  = "text"
      value = "network-team"
    }
  ]
}
`, cfEnv, cfPurpose, cfTeam)
}

func testAccIPRangeConfig_filter_step4(cfEnv, cfPurpose, cfTeam string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "purpose" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "team" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10"
  end_address   = "192.168.1.20"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.purpose.name
      type  = "text"
      value = "servers"
    },
    {
      name  = netbox_custom_field.team.name
      type  = "text"
      value = "network-team"
    }
  ]
}
`, cfEnv, cfPurpose, cfTeam)
}
func TestAccIPRangeResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	startAddress := "192.0.2.1"
	endAddress := "192.0.2.10"
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
				Config: testAccIPRangeResourceImportConfig_full(startAddress, endAddress, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccIPRangeResourceImportConfig_full(startAddress, endAddress, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_ip_range.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags", "tenant", "start_address", "end_address"},
			},
		},
	})
}

func testAccIPRangeResourceImportConfig_full(startAddress, endAddress, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["ipam.iprange"]
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
resource "netbox_ip_range" "test" {
  start_address = %q
  end_address   = %q
  tenant        = netbox_tenant.test.slug

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
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		startAddress, endAddress,
	)
}
