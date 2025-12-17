package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTypeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeDataSourceConfig("Test Circuit Type"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "name", "Test Circuit Type"),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "slug", "test-circuit-type"),
				),
			},
		},
	})
}

func testAccCircuitTypeDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name = "%s"
  slug = "test-circuit-type"
}

data "netbox_circuit_type" "test" {
  id = netbox_circuit_type.test.id
}
`, name)
}
