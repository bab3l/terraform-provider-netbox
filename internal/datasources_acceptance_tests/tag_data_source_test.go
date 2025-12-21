package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	tagName := testutil.RandomName("tag")
	tagSlug := testutil.RandomSlug("tag")

	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "slug", tagSlug),
				),
			},
		},
	})
}

func testAccTagDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_tag" "test" {
  id = netbox_tag.test.id
}
`, name, slug)
}
