//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderAccountDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_provideraccount_ds_cf")
	providerName := testutil.RandomName("tf-test-provider-ds-cf")
	providerSlug := testutil.RandomSlug("tf-test-provider-ds-cf")
	accountNumber := testutil.RandomName("account-")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountDataSourceConfig_customFields(customFieldName, providerName, providerSlug, accountNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "account", accountNumber),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccProviderAccountDataSourceConfig_customFields(customFieldName, providerName, providerSlug, accountNumber string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.provideraccount"]
  type         = "text"
}

resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  account     = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_provider_account" "test" {
  account = %q

  depends_on = [netbox_provider_account.test]
}
`, customFieldName, providerName, providerSlug, accountNumber, accountNumber)
}
