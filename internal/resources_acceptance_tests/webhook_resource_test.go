package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource_basic(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-webhook")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWebhookResource(randomName, "https://example.com/webhook1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_webhook.test", "payload_url", "https://example.com/webhook1"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_method", "POST"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "ssl_verification", "true"),
					resource.TestCheckResourceAttrSet("netbox_webhook.test", "id"),
				),
			},
			// PlanOnly after Create
			{
				Config:   testAccWebhookResource(randomName, "https://example.com/webhook1"),
				PlanOnly: true,
			},
			// ImportState testing
			{
				ResourceName:            "netbox_webhook.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
			// PlanOnly after Import
			{
				Config:   testAccWebhookResource(randomName, "https://example.com/webhook1"),
				PlanOnly: true,
			},
			// Update and Read testing
			{
				Config: testAccWebhookResource(randomName, "https://example.com/webhook2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_webhook.test", "payload_url", "https://example.com/webhook2"),
				),
			},
			// PlanOnly after Update
			{
				Config:   testAccWebhookResource(randomName, "https://example.com/webhook2"),
				PlanOnly: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccWebhookResource_full(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-webhook-full")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: testAccWebhookResourceFull(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_webhook.test", "payload_url", "https://example.com/webhook"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_method", "PUT"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_content_type", "application/xml"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "description", "Test webhook description"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "additional_headers", "X-Custom-Header: test-value"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "ssl_verification", "false"),
					resource.TestCheckResourceAttrSet("netbox_webhook.test", "id"),
				),
			},
			// PlanOnly after Create
			{
				Config:   testAccWebhookResourceFull(randomName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccWebhookResource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-webhook-id")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookResource(randomName, "https://example.com/webhook"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_webhook.test", "id"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "name", randomName),
				),
			},
			// PlanOnly after Create
			{
				Config:   testAccWebhookResource(randomName, "https://example.com/webhook"),
				PlanOnly: true,
			},
		},
	})
}

func testAccWebhookResource(name, payloadURL string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %[1]q
  payload_url = %[2]q
}
`, name, payloadURL)
}

func testAccWebhookResourceFull(name string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name               = %[1]q
  payload_url        = "https://example.com/webhook"
  http_method        = "PUT"
  http_content_type  = "application/xml"
  description        = "Test webhook description"
  additional_headers = "X-Custom-Header: test-value"
  ssl_verification   = false
}
`, name)
}

func TestAccConsistency_Webhook_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("webhook")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_webhook.test", "id"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "name", name),
				),
			},
			{
				Config:   testAccWebhookConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_webhook.test", "id"),
				),
			},
		},
	})
}

func testAccWebhookConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %q
  payload_url = "https://example.com/webhook"
  http_method = "POST"
}
`, name)
}

func TestAccWebhookResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-webhook-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "description", testutil.Description1),
				),
			},
			// PlanOnly after Create
			{
				Config:   testAccWebhookResourceConfig_withDescription(name, testutil.Description1),
				PlanOnly: true,
			},
			{
				Config: testAccWebhookResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "description", testutil.Description2),
				),
			},
			// PlanOnly after Update
			{
				Config:   testAccWebhookResourceConfig_withDescription(name, testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func testAccWebhookResourceConfig_withDescription(name string, description string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = %[1]q
  payload_url = "http://example.com/webhook"
  description = %[2]q
}
`, name, description)
}

func TestAccWebhookResource_removeOptionalFields(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-webhook-remove")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with optional fields
			{
				Config: testAccWebhookResourceWithOptionalFields(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "description", "Test webhook description"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "additional_headers", "X-Custom-Header: test-value"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "body_template", "{ \"foo\": \"bar\" }"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "secret", "mysecretkey"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "ca_file_path", "/path/to/ca.crt"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_content_type", "application/xml"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_method", "PUT"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "ssl_verification", "true"),
				),
			},
			// Update to remove optional fields
			{
				Config: testAccWebhookResource(randomName, "https://example.com/webhook"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_webhook.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_webhook.test", "additional_headers"),
					resource.TestCheckNoResourceAttr("netbox_webhook.test", "body_template"),
					resource.TestCheckNoResourceAttr("netbox_webhook.test", "secret"),
					resource.TestCheckNoResourceAttr("netbox_webhook.test", "ca_file_path"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_content_type", "application/json"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "http_method", "POST"),
					resource.TestCheckResourceAttr("netbox_webhook.test", "ssl_verification", "true"),
				),
			},
		},
	})
}

func testAccWebhookResourceWithOptionalFields(name string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name               = %[1]q
  payload_url        = "https://example.com/webhook"
  http_method        = "PUT"
  http_content_type  = "application/xml"
  ca_file_path       = "/path/to/ca.crt"
  description        = "Test webhook description"
  additional_headers = "X-Custom-Header: test-value"
  body_template      = "{ \"foo\": \"bar\" }"
  secret             = "mysecretkey"
  ssl_verification   = true
}
`, name)
}

func TestAccWebhookResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-webhook-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookResource(name, "http://example.com/webhook"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_webhook.test", "id"),
				),
			},
			// PlanOnly after Create
			{
				Config:   testAccWebhookResource(name, "http://example.com/webhook"),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find webhook by name
					items, _, err := client.ExtrasAPI.ExtrasWebhooksList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list webhooks: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Webhook not found with name: %s", name)
					}

					// Delete the webhook
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasWebhooksDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete webhook: %v", err)
					}

					t.Logf("Successfully externally deleted webhook with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWebhookResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_webhook",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_webhook" "test" {
  # name missing
  payload_url = "https://example.com/webhook"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_payload_url": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_webhook" "test" {
  name = "Test Webhook"
  # payload_url missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
