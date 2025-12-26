package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNDataSource_IDPreservation(t *testing.T) {
	cleanup := testutil.NewCleanupResource(t)
	name := acctest.RandomWithPrefix("test-l2vpn-ds-id")
	cleanup.RegisterL2VPNCleanup(name)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckL2VPNDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", name),
				),
			},
		},
	})
}

func TestAccL2VPNDataSource_byID(t *testing.T) {

	cleanup := testutil.NewCleanupResource(t)

	name := acctest.RandomWithPrefix("test-l2vpn-ds")

	cleanup.RegisterL2VPNCleanup(name)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckL2VPNDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func TestAccL2VPNDataSource_byName(t *testing.T) {

	cleanup := testutil.NewCleanupResource(t)

	name := acctest.RandomWithPrefix("test-l2vpn-ds")

	cleanup.RegisterL2VPNCleanup(name)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckL2VPNDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_byName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func testAccL2VPNDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

data "netbox_l2vpn" "test" {
  id = netbox_l2vpn.test.id
}
`, name, name)
}

func testAccL2VPNDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

data "netbox_l2vpn" "test" {
  name = netbox_l2vpn.test.name

  depends_on = [netbox_l2vpn.test]
}
`, name, name)
}

func TestAccL2VPNDataSource_bySlug(t *testing.T) {

	cleanup := testutil.NewCleanupResource(t)

	name := acctest.RandomWithPrefix("test-l2vpn-ds")

	cleanup.RegisterL2VPNCleanup(name)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckL2VPNDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_bySlug(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func testAccL2VPNDataSourceConfig_bySlug(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

data "netbox_l2vpn" "test" {
  slug = netbox_l2vpn.test.slug

  depends_on = [netbox_l2vpn.test]
}
`, name, name)
}
