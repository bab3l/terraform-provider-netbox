package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
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

func TestAccTunnelResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnTunnelsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tunnel for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnTunnelsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tunnel: %v", err)
					}
					t.Logf("Successfully externally deleted tunnel with ID: %d", itemID)
				},
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
				),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

// TestAccTunnelResource_StatusComprehensive tests comprehensive scenarios for tunnel status field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccTunnelResource_StatusComprehensive(t *testing.T) {
	t.Parallel()

	// Generate unique names for this test run
	tunnelName := testutil.RandomName("tf-test-tunnel-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_tunnel",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTunnelDestroy,
			testutil.CheckTunnelGroupDestroy,
		),
		BaseConfig: func() string {
			return testAccTunnelResourceConfig_statusBase(tunnelName)
		},
		WithFieldConfig: func(value string) string {
			return testAccTunnelResourceConfig_statusWithField(tunnelName, value)
		},
	})
}

func testAccTunnelResourceConfig_statusBase(name string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
	name          = %[1]q
	encapsulation = "gre"
	# status field intentionally omitted - should get default "active"
}
`, name)
}

func testAccTunnelResourceConfig_statusWithField(name, status string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
	name          = %[1]q
	encapsulation = "gre"
	status        = %[2]q
}
`, name, status)
}

func TestAccTunnelResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-rem")
	groupName := testutil.RandomName("tf-test-tunnel-group")
	groupSlug := testutil.RandomSlug("tf-test-tunnel-group")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTunnelGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_tunnel",
		BaseConfig: func() string {
			return testAccTunnelResourceConfig_removeOptionalFields_base(
				name, groupName, groupSlug, tenantName, tenantSlug,
			)
		},
		ConfigWithFields: func() string {
			return testAccTunnelResourceConfig_removeOptionalFields_withFields(
				name, groupName, groupSlug, tenantName, tenantSlug,
			)
		},
		OptionalFields: map[string]string{
			"description": "Test Description",
			"comments":    "Test Comments",
			"tunnel_id":   "100",
			// Note: status has a default value and cannot be truly cleared
			// Note: ipsec_profile requires ipsec encapsulation type
		},
		RequiredFields: map[string]string{
			"name": name,
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
	})
}

func testAccTunnelResourceConfig_removeOptionalFields_base(name, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"
}
`, name, groupName, groupSlug, tenantName, tenantSlug)
}

func testAccTunnelResourceConfig_removeOptionalFields_withFields(
	name, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"
  description   = "Test Description"
  comments      = "Test Comments"
  tunnel_id     = 100
  group         = netbox_tunnel_group.test.id
  tenant        = netbox_tenant.test.id
}
`, name, groupName, groupSlug, tenantName, tenantSlug)
}
