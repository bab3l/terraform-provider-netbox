package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventRuleDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventRuleDataSourceConfig("Test Event Rule"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "name", "Test Event Rule"),
					resource.TestCheckResourceAttr("data.netbox_event_rule.test", "action_type", "webhook"),
				),
			},
		},
	})
}

func testAccEventRuleDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_webhook" "test" {
  name        = "Test Webhook"
  payload_url = "http://example.com/webhook"
}

resource "netbox_event_rule" "test" {
  name               = "%s"
  object_types       = ["dcim.site"]
  event_types        = ["object_created"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test.id
}

data "netbox_event_rule" "test" {
  id = netbox_event_rule.test.id
}
`, name)
}
