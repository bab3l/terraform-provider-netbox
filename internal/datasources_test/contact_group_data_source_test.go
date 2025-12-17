package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactGroupDataSource_basic(t *testing.T) {
	name := testutil.RandomName("test-contact-group")
	slug := testutil.GenerateSlug(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.test", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.test", "description", "Test Contact Group Description"),
				),
			},
		},
	})
}

func testAccContactGroupDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name        = %q
  slug        = %q
  description = "Test Contact Group Description"
}

data "netbox_contact_group" "test" {
  id = netbox_contact_group.test.id
}
`, name, slug)
}
