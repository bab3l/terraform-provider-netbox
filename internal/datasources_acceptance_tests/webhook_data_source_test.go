package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccWebhookDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("test-webhook-ds-id")

	cleanup.RegisterWebhookCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckWebhookDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccWebhookDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_webhook.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_webhook.test", "name", randomName),
				),
			},
		},
	})

}

func TestAccWebhookDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("test-webhook-ds")

	cleanup.RegisterWebhookCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckWebhookDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccWebhookDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_webhook.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_webhook.test", "id"),
				),
			},
		},
	})

}

func TestAccWebhookDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("test-webhook-ds")

	cleanup.RegisterWebhookCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckWebhookDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccWebhookDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_webhook.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_webhook.test", "id"),
				),
			},
		},
	})

}

func testAccWebhookDataSourceByID(name string) string {

	return fmt.Sprintf(`

resource "netbox_webhook" "test" {

  name        = %[1]q

  payload_url = "https://example.com/webhook"

}

data "netbox_webhook" "test" {

  id = netbox_webhook.test.id

}

`, name)

}

func testAccWebhookDataSourceByName(name string) string {

	return fmt.Sprintf(`

resource "netbox_webhook" "test" {

  name        = %[1]q

  payload_url = "https://example.com/webhook"

}

data "netbox_webhook" "test" {

  name = netbox_webhook.test.name

}

`, name)

}
