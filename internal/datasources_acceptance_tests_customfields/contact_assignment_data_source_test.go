//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactAssignmentDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_contact_assign_ds_cf")
	siteName := testutil.RandomName("tf-test-site-ca-cf")
	contactName := testutil.RandomName("tf-test-contact-ca-cf")
	roleName := testutil.RandomName("tf-test-role-ca-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentDataSourceConfig_customFields(customFieldName, siteName, contactName, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_assignment.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_contact_assignment.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_contact_assignment.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_contact_assignment.test", "custom_fields.0.value", "test-contact-assignment-value"),
				),
			},
		},
	})
}

func testAccContactAssignmentDataSourceConfig_customFields(customFieldName, siteName, contactName, roleName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  object_types = ["tenancy.contactassignment"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_contact" "test" {
  name = %[3]q
}

resource "netbox_contact_role" "test" {
  name = %[4]q
  slug = %[4]q
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-contact-assignment-value"
    }
  ]
}

data "netbox_contact_assignment" "test" {
  id = netbox_contact_assignment.test.id

  depends_on = [netbox_contact_assignment.test]
}
`, customFieldName, siteName, contactName, roleName)
}
