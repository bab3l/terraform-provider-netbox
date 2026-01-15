package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEventRuleResource_basic(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("tf-test-eventrule-basic")
	webhookName := testutil.RandomName("tf-test-webhook-basic")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckEventRuleDestroy,
			testutil.CheckWebhookDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_event_rule.test", "id"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", eventRuleName),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "action_type", "webhook"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "action_object_type", "extras.webhook"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "object_types.#", "1"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "event_types.#", "1"),
				),
			},
		},
	})
}

func TestAccEventRuleResource_full(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("tf-test-eventrule-full")
	webhookName := testutil.RandomName("tf-test-webhook-full")
	description := "Test event rule with all fields"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckEventRuleDestroy,
			testutil.CheckWebhookDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_fullWithoutConditions(eventRuleName, webhookName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_event_rule.test", "id"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", eventRuleName),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "description", description),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "action_type", "webhook"),
				),
			},
		},
	})
}

func TestAccEventRuleResource_update(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("tf-test-eventrule-update")
	updatedName := testutil.RandomName("tf-test-eventrule-updated")
	webhookName := testutil.RandomName("tf-test-webhook-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterEventRuleCleanup(updatedName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckEventRuleDestroy,
			testutil.CheckWebhookDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", eventRuleName),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "enabled", "true"),
				),
			},
			{
				Config: testAccEventRuleResourceConfig_updated(updatedName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_event_rule.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("netbox_event_rule.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccEventRuleResource_import(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("tf-test-eventrule-import")
	webhookName := testutil.RandomName("tf-test-webhook-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckEventRuleDestroy,
			testutil.CheckWebhookDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_event_rule.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_event_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccEventRuleResource_IDPreservation(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("tf-test-eventrule-id")
	webhookName := testutil.RandomName("tf-test-webhook-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckEventRuleDestroy,
			testutil.CheckWebhookDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_event_rule.test", "id"),
				),
			},
		},
	})
}

func TestAccEventRuleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	eventRuleName := testutil.RandomName("tf-test-eventrule-ext-del")
	webhookName := testutil.RandomName("tf-test-webhook-ext-del")
	var eventRuleID string

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterEventRuleCleanup(eventRuleName)
	cleanup.RegisterWebhookCleanup(webhookName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckEventRuleDestroy,
			testutil.CheckWebhookDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_event_rule.test", "id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["netbox_event_rule.test"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						eventRuleID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					if eventRuleID == "" {
						t.Fatal("event rule ID not captured from previous step")
					}
					// Delete the event rule externally via API
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("failed to get API client: %s", err)
					}
					id, err := strconv.Atoi(eventRuleID)
					if err != nil {
						t.Fatalf("failed to convert ID to int: %s", err)
					}
					//nolint:gosec // ID from NetBox API is always a valid positive integer
					_, err = client.ExtrasAPI.ExtrasEventRulesDestroy(context.Background(), int32(id)).Execute()
					if err != nil {
						t.Fatalf("failed to delete event rule externally: %s", err)
					}
				},
				Config: testAccEventRuleResourceConfig_basic(eventRuleName, webhookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_event_rule.test", "id"),
				),
			},
		},
	})
}

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

func testAccEventRuleResourceConfig_fullWithoutConditions(eventRuleName, webhookName, description string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %[2]q
  payload_url = "https://example.com/webhook"
}

resource "netbox_event_rule" "test" {
  name               = %[1]q
  description        = %[3]q
  object_types       = ["dcim.device", "dcim.site"]
  event_types        = ["object_created", "object_updated"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test.id
  enabled            = true
}
`, eventRuleName, webhookName, description)
}

func testAccEventRuleResourceConfig_updated(eventRuleName, webhookName string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %[2]q
  payload_url = "https://example.com/webhook"
}

resource "netbox_event_rule" "test" {
  name               = %[1]q
  description        = "Updated description"
  object_types       = ["dcim.device"]
  event_types        = ["object_created"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test.id
  enabled            = false
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
