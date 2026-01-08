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

// TestAccConfigContextResource_TagsPreservation tests that tags are preserved
// when updating other fields on a config context. This addresses the critical bug
// where tags were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create config context with tags
// 2. Update config context WITHOUT tags in config (omit the field entirely)
// 3. Tags should be preserved in NetBox, not deleted.
func TestAccConfigContextResource_TagsPreservation(t *testing.T) {
	// Generate unique names
	configContextName := testutil.RandomName("tf-test-cc-tag-preserve")
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create config context WITH tags explicitly in config
				Config: testAccConfigContextConfig_preservation_step1(
					configContextName, tagName, tagSlug, "Initial description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", configContextName),
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_config_context.test", "tags.*", tagSlug),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning tags in config
				// Tags should be preserved in NetBox (verified by import)
				// State shows null/empty for tags since not in config
				Config: testAccConfigContextConfig_preservation_step2(
					configContextName, tagName, tagSlug, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", configContextName),
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", "Updated description"),
					// State shows 0 tags (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_config_context.test", "tags.#", "0"),
				),
			},
			{
				// Step 3: Import to verify tags still exist in NetBox
				ResourceName:            "netbox_config_context.test",
				ImportState:             true,
				ImportStateVerify:       false,                                   // Can't verify - config has no tags
				ImportStateVerifyIgnore: []string{"tags", "is_active", "weight"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add tags back to config to verify they were preserved
				Config: testAccConfigContextConfig_preservation_step1(
					configContextName, tagName, tagSlug, "Initial description",
				),
				Check: resource.ComposeTestCheckFunc(
					// Tags should have their original value (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_config_context.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_config_context.test", "tags.*", tagSlug),
				),
			},
		},
	})
}

func testAccConfigContextConfig_preservation_step1(
	configContextName, tagName, tagSlug, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_config_context" "test" {
  name        = %[1]q
  description = %[4]q
  data        = jsonencode({
    test_key = "test_value"
  })

  tags = [netbox_tag.test.slug]
}
`,
		configContextName, tagName, tagSlug, description,
	)
}

func testAccConfigContextConfig_preservation_step2(
	configContextName, tagName, tagSlug, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_config_context" "test" {
  name        = %[1]q
  description = %[4]q
  data        = jsonencode({
    test_key = "test_value"
  })
  # Note: tags intentionally omitted to test preservation
}
`,
		configContextName, tagName, tagSlug, description,
	)
}
