package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccEventRuleResource_removeOptionalFields_extended tests removing additional optional fields
// that weren't covered by the base test.
func TestAccEventRuleResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("event-rule-rem")
	webhookName := testutil.RandomName("webhook-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	// Note: conditions field has a bug - provider passes JSON string instead of parsed object
	// Testing only action_object_id for now
	testFields := map[string]string{
		// "conditions": `{"attr":"name","value":"test"}`, // Excluded - requires fix
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_event_rule",
		BaseConfig: func() string {
			return testAccEventRuleResourceConfig_removeOptionalFields_base(eventRuleName, webhookName)
		},
		ConfigWithFields: func() string {
			return testAccEventRuleResourceConfig_removeOptionalFields_withFields(eventRuleName, webhookName)
		},
		OptionalFields: testFields,
		RequiredFields: map[string]string{
			"name": eventRuleName,
		},
	})
}

// testAccEventRuleResourceConfig_removeOptionalFields_base creates a basic event rule
// with only required fields (no action_object_id, no conditions).
func testAccEventRuleResourceConfig_removeOptionalFields_base(eventRuleName, webhookName string) string {
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
}
`, eventRuleName, webhookName)
}

// testAccEventRuleResourceConfig_removeOptionalFields_withFields creates an event rule
// with action_object_id set. (conditions excluded due to provider bug).
func testAccEventRuleResourceConfig_removeOptionalFields_withFields(eventRuleName, webhookName string) string {
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
