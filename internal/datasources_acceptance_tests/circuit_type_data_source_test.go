package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTypeDataSource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("circuit-type")
	slug := testutil.RandomSlug("circuit-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "slug", slug),
				),
			},
		},
	})
}

func testAccCircuitTypeDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_circuit_type" "test" {
  id = netbox_circuit_type.test.id
}
`, name, slug)
}
