//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackTypeResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	model := testutil.RandomName("rack_type")
	slug := testutil.RandomSlug("rack_type")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceImportConfig_full(model, slug, mfgName, mfgSlug, cfText, cfInteger, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_rack_type.test", "custom_fields.#", "2"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:      "netbox_rack_type.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccRackTypeResourceImportConfig_full(model, slug, mfgName, mfgSlug, cfText, cfInteger, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccRackTypeResourceImportConfig_full(model, slug, mfgName, mfgSlug, cfText, cfInteger, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.racktype"]
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
resource "netbox_rack_type" "test" {
  model        = %q
  slug         = %q
	manufacturer = netbox_manufacturer.test.name
  form_factor  = "4-post-frame"

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    }
  ]

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`,
		mfgName, mfgSlug,
		cfText,
		cfInteger,
		tag1, tag1Slug,
		tag2, tag2Slug,
		model, slug,
	)
}

// TestAccRackTypeResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a rack type.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccRackTypeResource_CustomFieldsPreservation(t *testing.T) {
	model := testutil.RandomName("rack_type")
	slug := testutil.RandomSlug("rack_type")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")

	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create rack type WITH custom fields
				Config: testAccRackTypeConfig_preservation_step1(model, slug, mfgName, mfgSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_rack_type.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_rack_type.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccRackTypeConfig_preservation_step2(model, slug, mfgName, mfgSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Add custom_fields back to verify they were preserved
				Config: testAccRackTypeConfig_preservation_step3(model, slug, mfgName, mfgSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_rack_type.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_rack_type.test", cfOwner, "text", "team-a"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccRackTypeConfig_preservation_step1(model, slug, mfgName, mfgSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_custom_field" "env" {
  name         = %[5]q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_custom_field" "owner" {
  name         = %[6]q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_rack_type" "test" {
  model        = %[1]q
  slug         = %[2]q
  manufacturer = netbox_manufacturer.test.slug
  form_factor  = "4-post-frame"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
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
`, model, slug, mfgName, mfgSlug, cfEnv, cfOwner)
}

func testAccRackTypeConfig_preservation_step2(model, slug, mfgName, mfgSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_custom_field" "env" {
  name         = %[5]q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_custom_field" "owner" {
  name         = %[6]q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_rack_type" "test" {
  model        = %[1]q
  slug         = %[2]q
  manufacturer = netbox_manufacturer.test.slug
  form_factor  = "4-post-frame"
  description  = "Updated description"
  # custom_fields intentionally omitted - testing preservation
}
`, model, slug, mfgName, mfgSlug, cfEnv, cfOwner)
}

func testAccRackTypeConfig_preservation_step3(model, slug, mfgName, mfgSlug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_custom_field" "env" {
  name         = %[5]q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_custom_field" "owner" {
  name         = %[6]q
  type         = "text"
  object_types = ["dcim.racktype"]
}

resource "netbox_rack_type" "test" {
  model        = %[1]q
  slug         = %[2]q
  manufacturer = netbox_manufacturer.test.slug
  form_factor  = "4-post-frame"
  description  = "Updated description"

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
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
`, model, slug, mfgName, mfgSlug, cfEnv, cfOwner)
}
