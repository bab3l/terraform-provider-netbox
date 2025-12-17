package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tag.test", "name", "Test Tag"),
					resource.TestCheckResourceAttr("data.netbox_tag.test", "slug", "test-tag"),
				),
			},
		},
	})
}

const testAccTagDataSourceConfig = `
resource "netbox_tag" "test" {
  name = "Test Tag"
  slug = "test-tag"
}

data "netbox_tag" "test" {
  id = netbox_tag.test.id
}
`
