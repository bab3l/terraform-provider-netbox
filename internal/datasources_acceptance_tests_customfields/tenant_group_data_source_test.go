//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantGroupDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_tenantgroup_ds_cf")
	tenantGroupName := testutil.RandomName("tf-test-tenantgroup-ds-cf")
	tenantGroupSlug := testutil.RandomSlug("tf-test-tenantgroup-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupDataSourceConfig_customFields(customFieldName, tenantGroupName, tenantGroupSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", tenantGroupName),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccTenantGroupDataSourceConfig_customFields(customFieldName, tenantGroupName, tenantGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["tenancy.tenantgroup"]
  type         = "text"
}

resource "netbox_tenant_group" "test" {
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

data "netbox_tenant_group" "test" {
  name = %q

  depends_on = [netbox_tenant_group.test]
}
`, customFieldName, tenantGroupName, tenantGroupSlug, tenantGroupName)
}
