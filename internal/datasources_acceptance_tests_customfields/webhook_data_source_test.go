//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_webhook_ds_cf")
	webhookName := testutil.RandomName("tf-test-webhook-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookDataSourceConfig_customFields(customFieldName, webhookName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_webhook.test", "name", webhookName),
					resource.TestCheckResourceAttr("data.netbox_webhook.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_webhook.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_webhook.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_webhook.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccWebhookDataSourceConfig_customFields(customFieldName, webhookName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["extras.webhook"]
  type         = "text"
}

resource "netbox_webhook" "test" {
  name        = %q
  payload_url = "https://example.com/webhook"
  http_method = "POST"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_webhook" "test" {
  name = netbox_webhook.test.name
}
`, customFieldName, webhookName)
}
