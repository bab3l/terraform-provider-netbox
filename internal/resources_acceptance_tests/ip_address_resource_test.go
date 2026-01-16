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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ip)

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ip)

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ip1)
	cleanup.RegisterIPAddressCleanup(ip2)

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ip)

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
	cleanup.RegisterIPAddressCleanup(ip)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressResourceConfig_basic(address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
				),
			},
			{
				Config:   testAccIPAddressResourceConfig_basic(address),
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ip)

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
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
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

// TestAccIPAddressResource_removeOptionalFields tests that removing previously set VRF and tenant fields correctly sets them to null.
// This addresses the bug where removing a nullable reference field from the configuration would not clear it in NetBox,
// causing "Provider produced inconsistent result after apply" errors.
func TestAccIPAddressResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	vrfName := testutil.RandomName("test-vrf-remove")
	vrfRD := testutil.RandomName("65000:999")
	tenantName := testutil.RandomName("test-tenant-remove")
	tenantSlug := testutil.GenerateSlug(tenantName)
	address := fmt.Sprintf("10.250.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterVRFCleanup(vrfRD)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with VRF, tenant, and other optional fields
			{
				Config: testAccIPAddressResourceConfig_withVRFAndTenant(vrfName, vrfRD, tenantName, tenantSlug, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "role", "loopback"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "reserved"),
				),
			},
			// Step 2: Remove all optional fields - should set to null or defaults
			{
				Config: testAccIPAddressResourceConfig_withoutVRFAndTenant(vrfName, vrfRD, tenantName, tenantSlug, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_ip_address.test", "vrf"),
					resource.TestCheckNoResourceAttr("netbox_ip_address.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_ip_address.test", "dns_name"),
					resource.TestCheckNoResourceAttr("netbox_ip_address.test", "role"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"), // Should revert to default
				),
			},
			// Step 3: Re-add all optional fields - should work without errors
			{
				Config: testAccIPAddressResourceConfig_withVRFAndTenant(vrfName, vrfRD, tenantName, tenantSlug, address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "role", "loopback"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "reserved"),
				),
			},
		},
	})
}

func testAccIPAddressResourceConfig_withVRFAndTenant(vrfName, vrfRD, tenantName, tenantSlug, address string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %[1]q
  rd   = %[2]q
}

resource "netbox_tenant" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_ip_address" "test" {
  address  = %[5]q
  vrf      = netbox_vrf.test.id
  tenant   = netbox_tenant.test.id
  status   = "reserved"
  role     = "loopback"
  dns_name = "test.example.com"
}
`, vrfName, vrfRD, tenantName, tenantSlug, address)
}

func testAccIPAddressResourceConfig_withoutVRFAndTenant(vrfName, vrfRD, tenantName, tenantSlug, address string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %[1]q
  rd   = %[2]q
}

resource "netbox_tenant" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_ip_address" "test" {
  address = %[5]q
  # All optional fields removed - should set to null or default
}
`, vrfName, vrfRD, tenantName, tenantSlug, address)
}

func TestAccIPAddressResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	address := fmt.Sprintf("198.51.100.%d/32", acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_ip_address",
		BaseConfig: func() string {
			return testAccIPAddressResourceConfig_basic(address)
		},
		ConfigWithFields: func() string {
			return testAccIPAddressResourceConfig_withDescriptionAndComments(
				address,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"address": address,
		},
	})
}

func testAccIPAddressResourceConfig_withDescriptionAndComments(address, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address     = %[1]q
  status      = "active"
  description = %[2]q
  comments    = %[3]q
}
`, address, description, comments)
}
func TestAccIPAddressResource_validationErrors(t *testing.T) {
	t.Parallel()

	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_ip_address",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_address": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  status = "active"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_ip_format": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  address = "not-an-ip-address"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidFormat,
			},
			"missing_prefix_length": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  address = "192.168.1.1"
}
`
				},
				ExpectedError: testutil.ErrPatternInconsistent,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  address = "192.168.1.1/24"
  status  = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_role": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  address = "192.168.1.1/24"
  role    = "invalid_role"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_vrf_reference": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  address = "192.168.1.1/24"
  vrf     = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_ip_address" "test" {
  address = "192.168.1.1/24"
  tenant  = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

// =============================================================================
// STANDARDIZED TAG TESTS (using helpers)
// =============================================================================

// TestAccIPAddressResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccIPAddressResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	address := fmt.Sprintf("10.50.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_ip_address",
		ConfigWithoutTags: func() string {
			return testAccIPAddressResourceConfig_tagLifecycle(address, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccIPAddressResourceConfig_tagLifecycle(address, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccIPAddressResourceConfig_tagLifecycle(address, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckIPAddressDestroy,
	})
}

func testAccIPAddressResourceConfig_tagLifecycle(address, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = %[5]q
  slug = %[6]q
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %[1]q
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
`, address)
	case "tag2_tag3":
		return baseConfig + fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %[1]q
  tags = [
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    }
  ]
}
`, address)
	default: // "none"
		return baseConfig + fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %[1]q
  tags   = []
}
`, address)
	}
}

// TestAccIPAddressResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccIPAddressResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	address := fmt.Sprintf("10.51.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(address)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_ip_address",
		ConfigWithTagsOrderA: func() string {
			return testAccIPAddressResourceConfig_tagOrder(address, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccIPAddressResourceConfig_tagOrder(address, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckIPAddressDestroy,
	})
}

func testAccIPAddressResourceConfig_tagOrder(address, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = %[3]q
  slug = %[4]q
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug)

	if tag1First {
		return baseConfig + fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %[1]q
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
`, address)
	}

	return baseConfig + fmt.Sprintf(`
resource "netbox_ip_address" "test" {
  address = %[1]q
  tags = [
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    }
  ]
}
`, address)
}
