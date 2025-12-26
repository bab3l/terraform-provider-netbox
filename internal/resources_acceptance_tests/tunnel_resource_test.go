package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelResource_basic(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tunnel")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
				),
			},
		},
	})
}

func TestAccTunnelResource_IDPreservation(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tnl-id")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
				),
			},
		},
	})
}

func TestAccTunnelResource_full(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-full")
	description := testutil.RandomName("description")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_full(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "wireguard"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", description),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tunnel_id", "12345"),
				),
			},
		},
	})
}

func TestAccTunnelResource_update(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-upd")
	updatedDescription := testutil.Description2

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),
				),
			},
			{
				Config: testAccTunnelResourceConfig_full(name, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccTunnelResource_import(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
			},
			{
				ResourceName:      "netbox_tunnel.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_Tunnel_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tunnel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTunnelConsistencyLiteralNamesConfig(name),
			},
		},
	})
}

func testAccTunnelResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  status        = "active"
  encapsulation = "gre"
}
`, name)
}

func testAccTunnelResourceConfig_full(name, description string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  status        = "planned"
  encapsulation = "wireguard"
  description   = %[2]q
  tunnel_id     = 12345
}
`, name, description)
}

func testAccTunnelConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  status        = "active"
  encapsulation = "gre"
}
`, name)
}
