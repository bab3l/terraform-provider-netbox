package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitDataSourceConfig("TEST-CIRCUIT-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit.test", "cid", "TEST-CIRCUIT-1"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.test", "circuit_provider"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit.test", "type"),
				),
			},
		},
	})
}

func testAccCircuitDataSourceConfig(cid string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_circuit_type" "test" {
  name = "Test Circuit Type"
  slug = "test-circuit-type"
}

resource "netbox_circuit" "test" {
  cid              = "%s"
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
}

data "netbox_circuit" "test" {
  id = netbox_circuit.test.id
}
`, cid)
}
