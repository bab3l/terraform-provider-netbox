package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNotificationGroupResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccNotificationGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	const testDescription = "Test Description"

	name := testutil.RandomName("notification-group-remove")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_withDescription(name, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "description", testDescription),
				),
			},
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_notification_group.test", "description"),
				),
			},
		},
	})
}

func testAccNotificationGroupResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name = %[1]q
}
`, name)
}

func testAccNotificationGroupResourceConfig_withDescription(name, description string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
