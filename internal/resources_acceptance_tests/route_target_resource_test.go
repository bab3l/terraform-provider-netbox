package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteTargetResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("65000:100")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("65000:200")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRouteTargetDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Test route target with full options"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "comments", "Test comments for route target"),
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "tenant"),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRouteTargetResource_update(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("65000:300")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccRouteTargetResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Updated description"),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_updated(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRouteTargetResource_import(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("65000:100")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_route_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRouteTargetResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:999")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamRouteTargetsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find route target for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamRouteTargetsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete route target: %v", err)
					}
					t.Logf("Successfully externally deleted route target with ID: %d", itemID)
				},
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_RouteTarget_LiteralNames(t *testing.T) {

	t.Parallel()
	rtName := testutil.RandomName("65000:100")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", rtName),
					resource.TestCheckResourceAttr("netbox_route_target.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccRouteTargetResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:400")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccRouteTargetResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "test" {
  name = %q
}
`, name)
}

func testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_route_target" "test" {
  name        = %q
  description = "Test route target with full options"
  comments    = "Test comments for route target"
  tenant      = netbox_tenant.test.id
}
`, tenantName, tenantSlug, name)
}

func testAccRouteTargetResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "test" {
  name        = %q
  description = "Updated description"
}
`, name)
}

func testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_route_target" "test" {
  name = "%[1]s"
  # Use literal string name to mimic existing user state
  tenant = "%[2]s"

  depends_on = [netbox_tenant.test]
}
`, rtName, tenantName, tenantSlug)
}
