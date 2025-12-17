package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTerminationDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationDataSourceConfig("TEST-CIRCUIT-TERM"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "term_side", "A"),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "port_speed", "1000"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "circuit"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_termination.test", "site"),
				),
			},
		},
	})
}

func testAccCircuitTerminationDataSourceConfig(cid string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

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

resource "netbox_circuit_termination" "test" {
  circuit    = netbox_circuit.test.id
  term_side  = "A"
  site       = netbox_site.test.id
  port_speed = 1000
}

data "netbox_circuit_termination" "test" {
  id = netbox_circuit_termination.test.id
}
`, cid)
}
