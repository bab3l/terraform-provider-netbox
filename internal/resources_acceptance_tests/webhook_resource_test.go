package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource_basic(t *testing.T) {

	t.Parallel()
	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-webhook")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
			// ImportState testing
			{
				ResourceName:            "netbox_webhook.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
			// Update and Read testing
			{
				Config: testAccWebhookResource(randomName, "https://example.com/webhook2"),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_webhook.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_webhook.test", "payload_url", "https://example.com/webhook2"),
				),
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
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
