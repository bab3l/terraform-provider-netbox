package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNResource_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)
	cleanup.RegisterL2VPNCleanup(name + "-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
				ResourceName:            "netbox_l2vpn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name"}, // display_name is computed and may differ after name changes
			},
		},
	})
}

func TestAccL2VPNResource_full(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

// NOTE: Custom field tests for l2vpn resource are in resources_acceptance_tests_customfields package

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

func TestAccL2VPNResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func TestAccConsistency_L2VPN_LiteralNames(t *testing.T) {
	t.Parallel()

	name := "test-l2vpn-lit"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				Config:   testAccL2VPNResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}

func TestAccL2VPNResource_update(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_updateInitial(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", testutil.Description1),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_updateModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", testutil.Description2),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_updateInitial(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = %q
}
`, name, name, testutil.Description1)
}

func testAccL2VPNResourceConfig_updateModified(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = %q
}
`, name, name, testutil.Description2)
}

func TestAccL2VPNResource_external_deletion(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnL2vpnsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find l2vpn for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnL2vpnsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete l2vpn: %v", err)
					}
					t.Logf("Successfully externally deleted l2vpn with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccL2VPNResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn-opt")
	tenantName := testutil.RandomName("tf-test-tenant-l2vpn")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckL2VPNDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create L2VPN with tenant and identifier
			{
				Config: testAccL2VPNResourceConfig_withTenant(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "identifier"),
				),
			},
			// Step 2: Remove tenant and identifier (should set them to null)
			{
				Config: testAccL2VPNResourceConfig_withoutTenant(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_l2vpn.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_l2vpn.test", "identifier"),
				),
			},
			// Step 3: Re-add tenant and identifier (verify they can be set again)
			{
				Config: testAccL2VPNResourceConfig_withTenant(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "identifier"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_withTenant(name, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_l2vpn" "test" {
  name       = %[1]q
  slug       = %[1]q
  type       = "vxlan"
  tenant     = netbox_tenant.test.id
  identifier = 12345
}
`, name, tenantName, tenantSlug)
}

func testAccL2VPNResourceConfig_withoutTenant(name, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_l2vpn" "test" {
  name = %[1]q
  slug = %[1]q
  type = "vxlan"
}
`, name, tenantName, tenantSlug)
}
