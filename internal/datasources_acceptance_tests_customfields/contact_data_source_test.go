//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_contact_ds_cf")
	contactName := testutil.RandomName("tf-test-contact-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactDataSourceConfig_customFields(customFieldName, contactName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", contactName),
					// Note: Contact resource doesn't support custom_fields yet, so we just verify
					// the datasource can read the contact without errors. Custom fields would be
					// empty since they can't be set via the resource.
					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})
}

func testAccContactDataSourceConfig_customFields(customFieldName, contactName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["tenancy.contact"]
  type         = "text"
}

resource "netbox_contact" "test" {
  name = %q
}

data "netbox_contact" "test" {
  name = %q

  depends_on = [netbox_contact.test, netbox_custom_field.test]
}
`, customFieldName, contactName, contactName)
}
