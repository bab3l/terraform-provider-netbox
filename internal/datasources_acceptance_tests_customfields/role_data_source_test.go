//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_role_ds_cf")
	roleName := testutil.RandomName("tf-test-role-ds-cf")
	roleSlug := testutil.RandomSlug("tf-test-role-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig_customFields(customFieldName, roleName, roleSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_role.test", "name", roleName),
					resource.TestCheckResourceAttr("data.netbox_role.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_role.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_role.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_role.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfig_customFields(customFieldName, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_role" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_role" "test" {
  slug = netbox_role.test.slug

  depends_on = [netbox_role.test]
}
`, customFieldName, roleName, roleSlug)
}
