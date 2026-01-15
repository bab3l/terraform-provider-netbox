package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccEventRuleResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccEventRuleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	const testDescription = "Test Description"

	eventRuleName := testutil.RandomName("event-rule-remove")
	webhookName := testutil.RandomName("webhook-remove")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_withDescription(eventRuleName, webhookName, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", eventRuleName),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "description", testDescription),
				),
			},
			{
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", eventRuleName),
					resource.TestCheckNoResourceAttr("netbox_event_rule.test", "description"),
				),
			},
		},
	})
}

func testAccEventRuleResourceConfig_basic(eventRuleName, webhookName string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %[2]q
  payload_url = "https://example.com/webhook"
}

resource "netbox_event_rule" "test" {
  name               = %[1]q
  object_types       = ["dcim.device"]
  event_types        = ["object_created"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test.id
}
`, eventRuleName, webhookName)
}

func testAccEventRuleResourceConfig_withDescription(eventRuleName, webhookName, description string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %[2]q
  payload_url = "https://example.com/webhook"
}

resource "netbox_event_rule" "test" {
  name               = %[1]q
  description        = %[3]q
  object_types       = ["dcim.device"]
  event_types        = ["object_created"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test.id
}
`, eventRuleName, webhookName, description)
}
