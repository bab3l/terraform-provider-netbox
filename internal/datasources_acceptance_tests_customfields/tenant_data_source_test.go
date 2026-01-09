//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_tenant_ds_cf")
	tenantName := testutil.RandomName("tf-test-tenant-ds-cf")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantDataSourceConfig_customFields(customFieldName, tenantName, tenantSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "name", tenantName),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccTenantDataSourceConfig_customFields(customFieldName, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["tenancy.tenant"]
  type         = "text"
}

resource "netbox_tenant" "test" {
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

data "netbox_tenant" "test" {
  name = %q

  depends_on = [netbox_tenant.test]
}
`, customFieldName, tenantName, tenantSlug, tenantName)
}
