//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	vlanName := testutil.RandomName("vlan")
	vid := testutil.RandomVID()
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Generate random names once for the entire test
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceImportConfig_full(vlanName, int(vid), tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", int(vid))),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_vlan.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:      "netbox_vlan.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccVLANResourceImportConfig_full(vlanName, int(vid), tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccVLANResourceImportConfig_full(vlanName string, vid int, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {

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
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["ipam.vlan"]
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

# VLAN with comprehensive custom fields and tags
resource "netbox_vlan" "test" {
  name   = %q
  vid    = %d
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
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this VLAN resource for testing purposes."
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

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, vlanName, vid)
}

// TestAccVLANResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a VLAN.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccVLANResource_CustomFieldsPreservation(t *testing.T) {
	vlanName := testutil.RandomName("vlan")
	vid := testutil.RandomVID()

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create VLAN WITH custom fields
				Config: testAccVLANConfig_preservation_step1(vlanName, int(vid), cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", int(vid))),
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccVLANConfig_preservation_step2(vlanName, int(vid), cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_vlan.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccVLANConfig_preservation_step3(vlanName, int(vid), cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", "Updated description"),
				),
			},
		},
	})
}

// TestAccVLANResource_CustomFieldsFilterToOwned tests the filter-to-owned pattern
func TestAccVLANResource_CustomFieldsFilterToOwned(t *testing.T) {
	vlanName := testutil.RandomName("vlan")
	vid := testutil.RandomVID()

	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two fields
				Config: testAccVLANConfig_filter_step1(vlanName, int(vid), cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Remove owner, keep env with updated value
				Config: testAccVLANConfig_filter_step2(vlanName, int(vid), cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add cost_center
				Config: testAccVLANConfig_filter_step3(vlanName, int(vid), cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfCostCenter, "text", "CC123"),
				),
			},
			{
				// Step 4: Add owner back - should have preserved value
				Config: testAccVLANConfig_filter_step4(vlanName, int(vid), cfEnv, cfOwner, cfCostCenter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfOwner, "text", "team-a"),
					testutil.CheckCustomFieldValue("netbox_vlan.test", cfCostCenter, "text", "CC123"),
				),
			},
		},
	})
}

// Helper config functions for preservation tests
func testAccVLANConfig_preservation_step1(name string, vid int, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name = %[3]q
  vid  = %[4]d

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
`, cfEnv, cfOwner, name, vid)
}

func testAccVLANConfig_preservation_step2(name string, vid int, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name        = %[3]q
  vid         = %[4]d
  description = "Updated description"
}
`, cfEnv, cfOwner, name, vid)
}

func testAccVLANConfig_preservation_step3(name string, vid int, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name        = %[3]q
  vid         = %[4]d
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
`, cfEnv, cfOwner, name, vid)
}

// Helper config functions for filter-to-owned tests
func testAccVLANConfig_filter_step1(name string, vid int, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cost" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name = %[4]q
  vid  = %[5]d

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
`, cfEnv, cfOwner, cfCost, name, vid)
}

func testAccVLANConfig_filter_step2(name string, vid int, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cost" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name = %[4]q
  vid  = %[5]d

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]
}
`, cfEnv, cfOwner, cfCost, name, vid)
}

func testAccVLANConfig_filter_step3(name string, vid int, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cost" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name = %[4]q
  vid  = %[5]d

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
`, cfEnv, cfOwner, cfCost, name, vid)
}

func testAccVLANConfig_filter_step4(name string, vid int, cfEnv, cfOwner, cfCost string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cost" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_vlan" "test" {
  name = %[4]q
  vid  = %[5]d

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
`, cfEnv, cfOwner, cfCost, name, vid)
}
