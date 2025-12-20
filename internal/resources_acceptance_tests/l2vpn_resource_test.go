package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Updated description"),
				),
			},
			{
				ResourceName:      "netbox_l2vpn.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccL2VPNResource_full(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vpls"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "identifier", "12345"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Test L2VPN"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "comments", "Test comments"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}
`, name, name)
}

func testAccL2VPNResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = "Updated description"
}
`, name+"-updated", name)
}

func testAccL2VPNResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vpls"
  identifier  = 12345
  description = "Test L2VPN"
  comments    = "Test comments"
}
`, name, name)
}
