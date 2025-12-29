package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationGroupDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notif-grp")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_notification_group.test", "name", name),
				),
			},
		},
	})
}

func TestAccNotificationGroupDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("notif-grp-ds-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_notification_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_notification_group.test", "name", name),
				),
			},
		},
	})
}

func testAccNotificationGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name = %[1]q
}

data "netbox_notification_group" "test" {
  id = netbox_notification_group.test.id
}
`, name)
}
