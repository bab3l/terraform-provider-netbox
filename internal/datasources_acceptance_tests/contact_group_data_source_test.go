package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactGroupDataSource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("test-contact-group")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					// Check by_id lookup
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_id", "description", "Test Contact Group Description"),
					// Check by_name lookup
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_name", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_name", "description", "Test Contact Group Description"),
					// Check by_slug lookup
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_slug", "name", name),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_slug", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_contact_group.by_slug", "description", "Test Contact Group Description"),
					// Verify all lookups return same contact group
					resource.TestCheckResourceAttrPair("data.netbox_contact_group.by_id", "id", "data.netbox_contact_group.by_name", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_contact_group.by_id", "id", "data.netbox_contact_group.by_slug", "id"),
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

data "netbox_contact_group" "by_id" {
  id = netbox_contact_group.test.id
}

data "netbox_contact_group" "by_name" {
  name = netbox_contact_group.test.name
}

data "netbox_contact_group" "by_slug" {
  slug = netbox_contact_group.test.slug
}
`, name, slug)
}
