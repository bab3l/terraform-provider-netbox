package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	roleName := testutil.RandomName("test-role-ds")
	roleSlug := testutil.GenerateSlug(roleName)

	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRoleDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig(roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_role.test", "name", roleName),
					resource.TestCheckResourceAttr("data.netbox_role.test", "slug", roleSlug),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_role" "test" {
  id = netbox_role.test.id
}
`, name, slug)
}
