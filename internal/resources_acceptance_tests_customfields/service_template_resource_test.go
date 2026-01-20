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

// TestAccServiceTemplateResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a ServiceTemplate.
//
// Filter-to-owned pattern:
// - Custom fields declared in config are managed by Terraform
// - Custom fields NOT in config are preserved in NetBox but invisible to Terraform
func TestAccServiceTemplateResource_CustomFieldsPreservation(t *testing.T) {
	serviceTemplateName := "svctempl-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	cfEnvironment := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create ServiceTemplate WITH custom fields
				Config: testAccServiceTemplateConfig_preservation_step1(serviceTemplateName, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", serviceTemplateName),
					resource.TestCheckResourceAttr("netbox_service_template.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_service_template.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_service_template.test", cfOwner, "text", "platform-team"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields
				Config: testAccServiceTemplateConfig_preservation_step2(serviceTemplateName, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Updated service template"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_service_template.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
			{
				// Step 4: Add custom_fields back to verify they were preserved
				Config: testAccServiceTemplateConfig_preservation_step3(serviceTemplateName, cfEnvironment, cfOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_service_template.test", cfEnvironment, "text", "production"),
					testutil.CheckCustomFieldValue("netbox_service_template.test", cfOwner, "text", "platform-team"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Updated service template"),
				),
			},
		},
	})
}

// Step 1: Create with custom fields
func testAccServiceTemplateConfig_preservation_step1(name, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name          = %[2]q
  type          = "text"
  object_types  = ["ipam.servicetemplate"]
}

resource "netbox_custom_field" "owner" {
  name          = %[3]q
  type          = "text"
  object_types  = ["ipam.servicetemplate"]
}

resource "netbox_service_template" "test" {
  name     = %[1]q
  protocol = "tcp"
  ports    = [80, 443]

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "platform-team"
    }
  ]
}
`, name, cfEnv, cfOwner)
}

// Step 2: Update without custom_fields (they should be preserved)
func testAccServiceTemplateConfig_preservation_step2(name, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name          = %[2]q
  type          = "text"
  object_types  = ["ipam.servicetemplate"]
}

resource "netbox_custom_field" "owner" {
  name          = %[3]q
  type          = "text"
  object_types  = ["ipam.servicetemplate"]
}

resource "netbox_service_template" "test" {
  name        = %[1]q
  protocol    = "tcp"
  ports       = [80, 443]
  description = "Updated service template"

  # custom_fields omitted - should preserve existing values
}
`, name, cfEnv, cfOwner)
}

// Step 3/4: Add custom_fields back to show they were preserved
func testAccServiceTemplateConfig_preservation_step3(name, cfEnv, cfOwner string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "environment" {
  name          = %[2]q
  type          = "text"
  object_types  = ["ipam.servicetemplate"]
}

resource "netbox_custom_field" "owner" {
  name          = %[3]q
  type          = "text"
  object_types  = ["ipam.servicetemplate"]
}

resource "netbox_service_template" "test" {
  name        = %[1]q
  protocol    = "tcp"
  ports       = [80, 443]
  description = "Updated service template"

  custom_fields = [
    {
      name  = netbox_custom_field.environment.name
      type  = "text"
      value = "production"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "platform-team"
    }
  ]
}
`, name, cfEnv, cfOwner)
}
