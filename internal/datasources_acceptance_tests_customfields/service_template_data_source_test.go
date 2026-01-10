//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceTemplateDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_svctmpl_ds_cf")
	templateName := testutil.RandomName("tf-test-svctmpl-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateDataSourceConfig_customFields(customFieldName, templateName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "name", templateName),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccServiceTemplateDataSourceConfig_customFields(customFieldName, templateName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["ipam.servicetemplate"]
  type         = "text"
}

resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [443]

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_service_template" "test" {
  name = netbox_service_template.test.name

  depends_on = [netbox_service_template.test]
}
`, customFieldName, templateName)
}
