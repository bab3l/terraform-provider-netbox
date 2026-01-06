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

// TestAccVRFResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a VRF.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccVRFResource_CustomFieldsPreservation(t *testing.T) {
	vrfName := "vrf-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create VRF WITH custom fields
				Config: testAccVRFConfig_preservation_step1(vrfName, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", vrfName),
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfOwner, "text", "network-team"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccVRFConfig_preservation_step2(vrfName, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "Updated VRF"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_vrf.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccVRFConfig_preservation_step3(vrfName, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfOwner, "text", "network-team"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "Updated VRF"),
				),
			},
		},
	})
}

// TestAccVRFResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
func TestAccVRFResource_CustomFieldsFilterToOwned(t *testing.T) {
	vrfName := "vrf-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two fields
				Config: testAccVRFConfig_filter_step1(vrfName, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfOwner, "text", "network-team"),
				),
			},
			{
				// Step 2: Remove owner, keep env with updated value
				Config: testAccVRFConfig_filter_step2(vrfName, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add cost_center
				Config: testAccVRFConfig_filter_step3(vrfName, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfCostCenter, "text", "COST123"),
				),
			},
			{
				// Step 4: Add owner back - should have preserved value
				Config: testAccVRFConfig_filter_step4(vrfName, cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfOwner, "text", "network-team"),
					testutil.CheckCustomFieldValue("netbox_vrf.test", cfCostCenter, "text", "COST123"),
				),
			},
		},
	})
}

// =============================================================================
// Helper Config Functions - Preservation Tests
// =============================================================================

func testAccVRFConfig_preservation_step1(vrfName, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "network-team"
    }
  ]
}
`, vrfName, cfEnv, cfOwner)
}

func testAccVRFConfig_preservation_step2(vrfName, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name        = %[1]q
  description = "Updated VRF"
}
`, vrfName, cfEnv, cfOwner)
}

func testAccVRFConfig_preservation_step3(vrfName, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name        = %[1]q
  description = "Updated VRF"

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "network-team"
    }
  ]
}
`, vrfName, cfEnv, cfOwner)
}

// =============================================================================
// Helper Config Functions - Filter-to-Owned Tests
// =============================================================================

func testAccVRFConfig_filter_step1(vrfName, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "prod"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "network-team"
    }
  ]
}
`, vrfName, cfEnv, cfOwner, cfCost)
}

func testAccVRFConfig_filter_step2(vrfName, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]
}
`, vrfName, cfEnv, cfOwner, cfCost)
}

func testAccVRFConfig_filter_step3(vrfName, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.cost.name
      type  = "text"
      value = "COST123"
    }
  ]
}
`, vrfName, cfEnv, cfOwner, cfCost)
}

func testAccVRFConfig_filter_step4(vrfName, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "owner" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cost" {
  name         = %[4]q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_vrf" "test" {
  name = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "network-team"
    },
    {
      name  = netbox_custom_field.cost.name
      type  = "text"
      value = "COST123"
    }
  ]
}
`, vrfName, cfEnv, cfOwner, cfCost)
}

func TestAccVRFResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	vrfName := testutil.RandomName("vrf")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(vrfName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceImportConfig_full(vrfName, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", vrfName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_vrf.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_vrf.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccVRFResourceImportConfig_full(vrfName, tenantName, tenantSlug),
				PlanOnly: true,
			},
			{
				ResourceName:            "netbox_vrf.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant", "custom_fields"}, // Tenant may have lookup inconsistencies, custom fields have import limitations
			},
			{
				Config:   testAccVRFResourceImportConfig_full(vrfName, tenantName, tenantSlug),
				PlanOnly: true,
			},
		},
	})
}

func testAccVRFResourceImportConfig_full(vrfName, tenantName, tenantSlug string) string {
	// Custom field names with underscore format - deterministic based on vrfName
	cfText := fmt.Sprintf("cf_text_%s", strings.ReplaceAll(vrfName, "-", "_"))
	cfLongtext := fmt.Sprintf("cf_longtext_%s", strings.ReplaceAll(vrfName, "-", "_"))
	cfInteger := fmt.Sprintf("cf_integer_%s", strings.ReplaceAll(vrfName, "-", "_"))
	cfBoolean := fmt.Sprintf("cf_boolean_%s", strings.ReplaceAll(vrfName, "-", "_"))
	cfDate := fmt.Sprintf("cf_date_%s", strings.ReplaceAll(vrfName, "-", "_"))
	cfUrl := fmt.Sprintf("cf_url_%s", strings.ReplaceAll(vrfName, "-", "_"))
	cfJson := fmt.Sprintf("cf_json_%s", strings.ReplaceAll(vrfName, "-", "_"))

	// Tag names - deterministic based on vrfName (no random generation)
	tag1 := fmt.Sprintf("tag1_%s", vrfName)
	tag1Slug := fmt.Sprintf("tag1_%s", strings.ReplaceAll(vrfName, "-", "_"))
	tag2 := fmt.Sprintf("tag2_%s", vrfName)
	tag2Slug := fmt.Sprintf("tag2_%s", strings.ReplaceAll(vrfName, "-", "_"))

	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["ipam.vrf"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["ipam.vrf"]
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

# VRF with comprehensive custom fields and tags
resource "netbox_vrf" "test" {
  name   = %q
  tenant = netbox_tenant.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this VRF resource for testing purposes."
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
      value = "2023-01-15"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key": "value"})
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
`, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, vrfName)
}
