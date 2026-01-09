//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackRoleDataSource_customFields(t *testing.T) {
	roleName := testutil.RandomName("tf-test-rackrole-ds-cf")
	roleSlug := testutil.GenerateSlug(roleName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_rackrole_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(roleSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleDataSourceConfig_withCustomFields(roleName, roleSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "name", roleName),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccRackRoleDataSourceConfig_withCustomFields(roleName, roleSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.rackrole"]
  type         = "text"
}

resource "netbox_rack_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_rack_role" "test" {
  name = netbox_rack_role.test.name

  depends_on = [netbox_rack_role.test]
}
`, customFieldName, roleName, roleSlug)
}
