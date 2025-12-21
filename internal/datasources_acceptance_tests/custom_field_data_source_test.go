package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFieldDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldDataSourceConfig("test_custom_field"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "name", "test_custom_field"),
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "type", "text"),
				),
			},
		},
	})
}

func testAccCustomFieldDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = "%s"
  type         = "text"
  object_types = ["dcim.site"]
}

data "netbox_custom_field" "test" {
  id = netbox_custom_field.test.id
}
`, name)
}
