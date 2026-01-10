//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerFeedDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_powerfeed_ds_cf")
	feedName := testutil.RandomName("tf-test-powerfeed-ds-cf")
	panelName := testutil.RandomName("tf-test-panel-ds-cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedDataSourceConfig_customFields(customFieldName, feedName, panelName, siteName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccPowerFeedDataSourceConfig_customFields(customFieldName, feedName, panelName, siteName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.powerfeed"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_power_panel" "test" {
  name = %q
  site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
  name         = %q
  power_panel  = netbox_power_panel.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = %q

  depends_on = [netbox_power_feed.test]
}
`, customFieldName, siteName, siteName, panelName, feedName, feedName)
}
