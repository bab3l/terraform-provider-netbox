package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressResource_basic(t *testing.T) {
	t.Parallel()

	ip := fmt.Sprintf("192.0.%d.%d/24", 100+acctest.RandIntRange(0, 50), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_basic(ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip),
				),
			},
		},
	})
}

func TestAccIPAddressResource_full(t *testing.T) {
	t.Parallel()

	ip := fmt.Sprintf("10.0.%d.%d/32", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_full(ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Test IP address"),
				),
			},
		},
	})
}

func TestAccIPAddressResource_update(t *testing.T) {
	t.Parallel()

	ip1 := fmt.Sprintf("172.16.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	ip2 := fmt.Sprintf("172.16.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_basic(ip1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip1),
				),
			},
			{
				Config: testAccIPAddressResourceConfig_full(ip2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Test IP address"),
				),
			},
		},
	})
}

func TestAccIPAddressResource_import(t *testing.T) {
	t.Parallel()

	ip := fmt.Sprintf("203.0.113.%d/32", acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_basic(ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ip_address.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIPAddressResource_importWithTags(t *testing.T) {
	t.Parallel()

	ip := fmt.Sprintf("203.0.113.%d/32", acctest.RandIntRange(1, 254))
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceImportConfig_full(ip, tenantName, tenantSlug, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip),
				),
			},
			{
				ResourceName:            "netbox_ip_address.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant", "tags"}, // Tags have import limitations
			},
		},
	})
}

func testAccIPAddressResourceImportConfig_full(ip, tenantName, tenantSlug, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# IP Address with tags (no custom fields support)
resource "netbox_ip_address" "test" {
  address = %q
  tenant     = netbox_tenant.test.id

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`, tenantName, tenantSlug, tag1, tag1Slug, tag2, tag2Slug, ip)
}

func TestAccIPAddressResource_IDPreservation(t *testing.T) {
	t.Parallel()

	ip := fmt.Sprintf("192.0.%d.%d/24", 200+acctest.RandIntRange(0, 50), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_basic(ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip),
				),
			},
		},
	})

}
func testAccIPAddressResourceConfig_basic(address string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %q
}
`, address)
}

func testAccIPAddressResourceConfig_full(address string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address     = %q
  status      = "active"
  dns_name    = "test.example.com"
  description = "Test IP address"
}
`, address)
}

func TestAccConsistency_IPAddress_LiteralNames(t *testing.T) {
	t.Parallel()

	address := fmt.Sprintf("10.200.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressConsistencyLiteralNamesConfig(address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
				),
			},
			{
				Config:   testAccIPAddressConsistencyLiteralNamesConfig(address),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
				),
			},
		},
	})
}

func TestAccIPAddressResource_externalDeletion(t *testing.T) {
	t.Parallel()

	ip := fmt.Sprintf("192.0.%d.%d/24", acctest.RandIntRange(100, 150), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_basic(ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamIpAddressesList(context.Background()).Address([]string{ip}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IP address for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamIpAddressesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IP address: %v", err)
					}
					t.Logf("Successfully externally deleted IP address with ID: %d", itemID)
				},
				Config: testAccIPAddressResourceConfig_basic(ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
				),
			},
		},
	})
}

func testAccIPAddressConsistencyLiteralNamesConfig(address string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %q
}
`, address)
}

func TestAccReferenceNamePersistence_IPAddress_TenantVRF(t *testing.T) {
	t.Parallel()

	tenantName := testutil.RandomName("tf-test-tenant-ip")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-ip")
	vrfName := testutil.RandomName("tf-test-vrf")
	vrfRD := testutil.RandomSlug("tf-test-vrf")
	address := fmt.Sprintf("10.200.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterVRFCleanup(vrfName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressReferenceConfig_tenantVRF(tenantName, tenantSlug, vrfName, vrfRD, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "tenant", tenantName),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "vrf", vrfName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccIPAddressReferenceConfig_tenantVRF(tenantName, tenantSlug, vrfName, vrfRD, address),
			},
		},
	})
}

func testAccIPAddressReferenceConfig_tenantVRF(tenantName, tenantSlug, vrfName, vrfRD, address string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_vrf" "test" {
  name   = %[3]q
  rd     = %[4]q
  tenant = netbox_tenant.test.id
}

resource "netbox_ip_address" "test" {
  address = %[5]q
  tenant  = netbox_tenant.test.name
  vrf     = netbox_vrf.test.name
}
`, tenantName, tenantSlug, vrfName, vrfRD, address)
}

// TestAccIPAddress_TenantNameNotID verifies that when a tenant is specified by name,
// the state stores the name (not the numeric ID) and remains consistent after refresh.
func TestAccIPAddress_TenantNameNotID(t *testing.T) {
	t.Parallel()

	tenantName := testutil.RandomName("tf-test-tenant-namenotid")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-namenotid")
	address := fmt.Sprintf("10.201.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with tenant name
			{
				Config: testAccIPAddressConfig_tenantByName(tenantName, tenantSlug, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
					// Verify tenant is stored as NAME, not numeric ID
					resource.TestCheckResourceAttr("netbox_ip_address.test", "tenant", tenantName),
				),
			},
			// Step 2: Refresh state and verify no drift (tenant should still be name)
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					// After refresh, tenant should still be the name, not a number
					resource.TestCheckResourceAttr("netbox_ip_address.test", "tenant", tenantName),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccIPAddressConfig_tenantByName(tenantName, tenantSlug, address),
			},
		},
	})
}

func testAccIPAddressConfig_tenantByName(tenantName, tenantSlug, address string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_ip_address" "test" {
  address = %[3]q
  tenant  = netbox_tenant.test.name
}
`, tenantName, tenantSlug, address)
}

// TestAccIPAddress_TenantByID verifies that when a tenant is specified by ID,
// the state stores the ID correctly.
func TestAccIPAddress_TenantByID(t *testing.T) {
	t.Parallel()

	tenantName := testutil.RandomName("tf-test-tenant-byid")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-byid")
	address := fmt.Sprintf("10.202.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with tenant ID
			{
				Config: testAccIPAddressConfig_tenantByID(tenantName, tenantSlug, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
					// When specified by ID, tenant should be stored as ID
					resource.TestCheckResourceAttrPair("netbox_ip_address.test", "tenant", "netbox_tenant.test", "id"),
				),
			},
			// Step 2: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccIPAddressConfig_tenantByID(tenantName, tenantSlug, address),
			},
		},
	})
}

func testAccIPAddressConfig_tenantByID(tenantName, tenantSlug, address string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_ip_address" "test" {
  address = %[3]q
  tenant  = netbox_tenant.test.id
}
`, tenantName, tenantSlug, address)
}

// TestAccIPAddress_TenantReferenceByID_StoresID verifies that when a tenant is referenced
// by ID in config (like netbox_tenant.test.id), the state stores the ID consistently.
func TestAccIPAddress_TenantReferenceByID_StoresID(t *testing.T) {
	t.Parallel()

	tenantName := testutil.RandomName("tf-test-tenant-ref")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-ref")
	address := fmt.Sprintf("10.203.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with tenant referenced by ID
			{
				Config: testAccIPAddressConfig_tenantByID(tenantName, tenantSlug, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
					// When config uses .id, state should store the ID
					resource.TestCheckResourceAttrPair("netbox_ip_address.test", "tenant", "netbox_tenant.test", "id"),
				),
			},
			// Step 2: Refresh state - should remain consistent
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					// After refresh, tenant should still be the ID (no drift)
					resource.TestCheckResourceAttrPair("netbox_ip_address.test", "tenant", "netbox_tenant.test", "id"),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccIPAddressConfig_tenantByID(tenantName, tenantSlug, address),
			},
		},
	})
}
