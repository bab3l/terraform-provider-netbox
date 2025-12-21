package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationGroupDataSource_basic(t *testing.T) {

	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_notification_group.test", "name", "Test Notification Group"),
				),
			},
		},
	})
}

const testAccNotificationGroupDataSourceConfig = `
resource "netbox_notification_group" "test" {
  name = "Test Notification Group"
}

data "netbox_notification_group" "test" {
  id = netbox_notification_group.test.id
}
`
