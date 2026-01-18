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

func TestAccProviderResource_CustomFieldsPreservation(t *testing.T) {
	providerName := "Provider-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	providerSlug := "provider-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create Provider WITH custom fields
				Config: testAccProviderConfig_preservation_step1(providerName, providerSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "name", providerName),
					resource.TestCheckResourceAttr("netbox_provider.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_provider.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_provider.test", cfOwner, "text", "admin"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccProviderConfig_preservation_step2(providerName, providerSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "description", "Updated provider"),
					resource.TestCheckResourceAttr("netbox_provider.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_provider.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccProviderConfig_preservation_step3(providerName, providerSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_provider.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_provider.test", cfOwner, "text", "admin"),
					resource.TestCheckResourceAttr("netbox_provider.test", "description", "Updated provider"),
				),
			},
		},
	})
}

func TestAccProviderResource_importWithCustomFieldsAndTags(t *testing.T) {
	providerName := "Provider-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	providerSlug := "provider-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	cfEnvironment := testutil.RandomCustomFieldName("tf_env")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig_importTest(providerName, providerSlug, cfEnvironment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "name", providerName),
					resource.TestCheckResourceAttr("netbox_provider.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_provider.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
		},
	})
}

// Test config generators

func testAccProviderConfig_preservation_step1(name, slug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cf_env" {
  name        = %[3]q
  type        = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_custom_field" "cf_owner" {
  name        = %[4]q
  type        = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q

  custom_fields = [
    {
      name  = netbox_custom_field.cf_env.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.cf_owner.name
      type  = "text"
      value = "admin"
    }
  ]
}
`, name, slug, cfEnv, cfOwner)
}

func testAccProviderConfig_preservation_step2(name, slug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cf_env" {
  name        = %[3]q
  type        = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_custom_field" "cf_owner" {
  name        = %[4]q
  type        = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_provider" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Updated provider"
}
`, name, slug, cfEnv, cfOwner)
}

func testAccProviderConfig_preservation_step3(name, slug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cf_env" {
  name        = %[3]q
  type        = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_custom_field" "cf_owner" {
  name        = %[4]q
  type        = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_provider" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Updated provider"

  custom_fields = [
    {
      name  = netbox_custom_field.cf_env.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.cf_owner.name
      type  = "text"
      value = "admin"
    }
  ]
}
`, name, slug, cfEnv, cfOwner)
}

func testAccProviderConfig_importTest(name, slug, cfEnv string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[3]q
  type         = "text"
  object_types = ["circuits.provider"]
}

resource "netbox_tag" "test1" {
  name = "test-tag-1-%[3]s"
  slug = "test-tag-1-%[3]s"
}

resource "netbox_tag" "test2" {
  name = "test-tag-2-%[3]s"
  slug = "test-tag-2-%[3]s"
}

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    }
  ]

	tags = [netbox_tag.test1.slug, netbox_tag.test2.slug]
}
`, name, slug, cfEnv)
}
