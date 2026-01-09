//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_provider_ds_cf")
	providerName := testutil.RandomName("tf-test-provider-ds-cf")
	providerSlug := testutil.RandomSlug("tf-test-provider-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderDataSourceConfig_customFields(customFieldName, providerName, providerSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider.test", "name", providerName),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_provider.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccProviderDataSourceConfig_customFields(customFieldName, providerName, providerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.provider"]
  type         = "text"
}

resource "netbox_provider" "test" {
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

data "netbox_provider" "test" {
  name = %q

  depends_on = [netbox_provider.test]
}
`, customFieldName, providerName, providerSlug, providerName)
}
