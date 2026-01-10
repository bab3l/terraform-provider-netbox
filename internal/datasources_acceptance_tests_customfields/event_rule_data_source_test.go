//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventRuleDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_event_rule_ds_cf")
	eventRuleName := testutil.RandomName("tf-test-event-rule-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleDataSourceConfig_customFields(customFieldName, eventRuleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "name", eventRuleName),
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccEventRuleDataSourceConfig_customFields(customFieldName, eventRuleName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["extras.eventrule"]
  type         = "text"
}

resource "netbox_webhook" "test" {
  name        = "test-webhook-event-rule-ds-cf"
  payload_url = "https://example.com/webhook"
  http_method = "POST"
}

resource "netbox_event_rule" "test" {
  name               = %q
  object_types       = ["dcim.device"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = tostring(netbox_webhook.test.id)
  event_types        = ["object_created"]
  enabled            = true

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_event_rule" "test" {
  id = netbox_event_rule.test.id
}
`, customFieldName, eventRuleName)
}
