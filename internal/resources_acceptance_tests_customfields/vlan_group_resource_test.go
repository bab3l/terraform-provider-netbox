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

func TestAccVLANGroupResource_CustomFieldsPreservation(t *testing.T) {
	vlanGroupName := "VLAN-Group-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	vlanGroupSlug := "vlan-group-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create VLAN Group WITH custom fields
				Config: testAccVLANGroupConfig_preservation_step1(vlanGroupName, vlanGroupSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", vlanGroupName),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vlan_group.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_vlan_group.test", cfOwner, "text", "admin"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccVLANGroupConfig_preservation_step2(vlanGroupName, vlanGroupSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "description", "Updated VLAN group"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_vlan_group.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccVLANGroupConfig_preservation_step3(vlanGroupName, vlanGroupSlug, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_vlan_group.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_vlan_group.test", cfOwner, "text", "admin"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "description", "Updated VLAN group"),
				),
			},
		},
	})
}

func TestAccVLANGroupResource_importWithCustomFieldsAndTags(t *testing.T) {
	vlanGroupName := "VLAN-Group-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	vlanGroupSlug := "vlan-group-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	cfEnvironment := testutil.RandomCustomFieldName("tf_env")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupConfig_importTest(vlanGroupName, vlanGroupSlug, cfEnvironment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", vlanGroupName),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:      "netbox_vlan_group.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccVLANGroupConfig_importTest(vlanGroupName, vlanGroupSlug, cfEnvironment),
				PlanOnly: true,
			},
		},
	})
}

// Test config generators

func testAccVLANGroupConfig_preservation_step1(name, slug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cf_env" {
  name        = %[3]q
  type        = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_custom_field" "cf_owner" {
  name        = %[4]q
  type        = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_vlan_group" "test" {
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

func testAccVLANGroupConfig_preservation_step2(name, slug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cf_env" {
  name        = %[3]q
  type        = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_custom_field" "cf_owner" {
  name        = %[4]q
  type        = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_vlan_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Updated VLAN group"
}
`, name, slug, cfEnv, cfOwner)
}

func testAccVLANGroupConfig_preservation_step3(name, slug, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cf_env" {
  name        = %[3]q
  type        = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_custom_field" "cf_owner" {
  name        = %[4]q
  type        = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_vlan_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = "Updated VLAN group"

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

func testAccVLANGroupConfig_importTest(name, slug, cfEnv string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name         = %[3]q
  type         = "text"
  object_types = ["ipam.vlangroup"]
}

resource "netbox_tag" "test1" {
  name = "test-tag-1-%[3]s"
  slug = "test-tag-1-%[3]s"
}

resource "netbox_tag" "test2" {
  name = "test-tag-2-%[3]s"
  slug = "test-tag-2-%[3]s"
}

resource "netbox_vlan_group" "test" {
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
