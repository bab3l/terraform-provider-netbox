package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests.

func TestAccASNRangeResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-asn-range")
	slug := testutil.RandomSlug("tf-test-asn-range")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-asn-range-full")
	slug := testutil.RandomSlug("tf-test-asn-range-full")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "65000"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "65100"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Test ASN range with full options"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-asn-range-upd")
	slug := testutil.RandomSlug("tf-test-asn-range-upd")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
				),
			},
			{
				Config: testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64700"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-remove")
	slug := testutil.RandomSlug("tf-test-asn-range-remove")
	rirName := testutil.RandomName("tf-test-rir-remove")
	rirSlug := testutil.RandomSlug("tf-test-rir-remove")
	tenantName := testutil.RandomName("tf-test-tenant-remove")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-remove")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create ASN range with tenant
			{
				Config: testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Test ASN range with full options"),
				),
			},
			// Step 2: Remove tenant and verify it's actually removed
			{
				Config: testAccASNRangeResourceConfig_noTenant(name, slug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_asn_range.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Description after tenant removal"),
				),
			},
			// Step 3: Re-add tenant to verify it can be set again
			{
				Config: testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_external_deletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-ext-del")
	slug := testutil.RandomSlug("tf-test-asn-range-ext-del")
	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamAsnRangesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find ASN range for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAsnRangesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete ASN range: %v", err)
					}
					t.Logf("Successfully externally deleted ASN range with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccASNRangeResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-id")
	slug := testutil.RandomSlug("tf-test-asn-range-id")
	rirName := testutil.RandomName("tf-test-rir-id")
	rirSlug := testutil.RandomSlug("tf-test-rir-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
				),
			},
		},
	})

}

func testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
  rir   = netbox_rir.test.id
  start = "64512"
  end   = "64612"
}
`, rirName, rirSlug, name, slug)
}

// Test with ID references.
func testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "65000"
  end         = "65100"
  tenant      = netbox_tenant.test.id
  description = "Test ASN range with full options"
}
`, rirName, rirSlug, tenantName, tenantSlug, name, slug)
}

func testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "64512"
  end         = "64700"
  description = "Updated description"
}
`, rirName, rirSlug, name, slug)
}

func testAccASNRangeResourceConfig_noTenant(name, slug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "65000"
  end         = "65100"
  description = "Description after tenant removal"
  # tenant intentionally omitted - should be null in state
}
`, rirName, rirSlug, tenantName, tenantSlug, name, slug)
}

func TestAccASNRangeResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-asn-range")
	slug := testutil.RandomSlug("tf-test-asn-range")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},
			{
				ResourceName:            "netbox_asn_range.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir"},
			},
			{
				Config:   testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_ASNRange(t *testing.T) {
	t.Parallel()

	rangeName := testutil.RandomName("asn-range")
	rangeSlug := testutil.RandomSlug("asn-range")
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(rangeName)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

// Test with mixed references (id and name).
func testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_tenant" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_asn_range" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  rir = netbox_rir.test.id
  tenant = netbox_tenant.test.name
  start = 65000
  end = 65100
}
`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)
}

// TestAccConsistency_ASNRange_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ASNRange_LiteralNames(t *testing.T) {
	t.Parallel()

	rangeName := testutil.RandomName("asn-range")
	rangeSlug := testutil.RandomSlug("asn-range")
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(rangeName)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "rir", rirSlug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_tenant" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_asn_range" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  # Use literal string names to mimic existing user state
  rir = "%[4]s"
  tenant = "%[5]s"
  start = 65000
  end = 65100
  depends_on = [netbox_rir.test, netbox_tenant.test]
}
`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)
}

func TestAccASNRangeResource_removeDescription(t *testing.T) {
	t.Parallel()

	rangeName := testutil.RandomName("tf-test-asnrange-optional")
	rangeSlug := testutil.RandomSlug("tf-test-asnrange-optional")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(rangeSlug)
	cleanup.RegisterRIRCleanup(rirSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_asn_range",
		BaseConfig: func() string {
			return testAccASNRangeResourceConfig_basic(rangeName, rangeSlug, rirName, rirSlug)
		},
		ConfigWithFields: func() string {
			return testAccASNRangeResourceConfig_withDescription(
				rangeName,
				rangeSlug,
				rirName,
				rirSlug,
				"Test description",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
		},
		RequiredFields: map[string]string{
			"name": rangeName,
			"slug": rangeSlug,
		},
		CheckDestroy: testutil.CheckASNRangeDestroy,
	})
}

func testAccASNRangeResourceConfig_withDescription(name, slug, rirName, rirSlug, description string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_asn_range" "test" {
  name        = %[1]q
  slug        = %[2]q
  rir         = netbox_rir.test.id
  start       = 65000
  end         = 65100
  description = %[5]q
}
`, name, slug, rirName, rirSlug, description)
}

// NOTE: Custom field tests for ASN Range resource are in resources_acceptance_tests_customfields package

func TestAccASNRangeResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_asn_range",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_rir" "test" {
  name = "test-rir"
  slug = "test-rir"
}

resource "netbox_asn_range" "test" {
  # name missing
  slug  = "test-range"
  rir   = netbox_rir.test.id
  start = "65000"
  end   = "65100"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_rir" "test" {
  name = "test-rir"
  slug = "test-rir"
}

resource "netbox_asn_range" "test" {
  name  = "Test Range"
  # slug missing
  rir   = netbox_rir.test.id
  start = "65000"
  end   = "65100"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_rir": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_asn_range" "test" {
  name  = "Test Range"
  slug  = "test-range"
  # rir missing
  start = "65000"
  end   = "65100"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_start": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_rir" "test" {
  name = "test-rir"
  slug = "test-rir"
}

resource "netbox_asn_range" "test" {
  name = "Test Range"
  slug = "test-range"
  rir  = netbox_rir.test.id
  # start missing
  end  = "65100"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_end": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_rir" "test" {
  name = "test-rir"
  slug = "test-rir"
}

resource "netbox_asn_range" "test" {
  name  = "Test Range"
  slug  = "test-range"
  rir   = netbox_rir.test.id
  start = "65000"
  # end missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_rir_reference": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_asn_range" "test" {
  name  = "Test Range"
  slug  = "test-range"
  rir   = "nonexistent-rir"
  start = "65000"
  end   = "65100"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

func TestAccASNRangeResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-tag")
	slug := testutil.RandomSlug("tf-test-asn-range-tag")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_asn_range",
		ConfigWithoutTags: func() string {
			return testAccASNRangeResourceConfig_tagLifecycle(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccASNRangeResourceConfig_tagLifecycle(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccASNRangeResourceConfig_tagLifecycle(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckASNRangeDestroy,
	})
}

func TestAccASNRangeResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-tagord")
	slug := testutil.RandomSlug("tf-test-asn-range-tagord")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_asn_range",
		ConfigWithTagsOrderA: func() string {
			return testAccASNRangeResourceConfig_tagOrder(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2_tag3")
		},
		ConfigWithTagsOrderB: func() string {
			return testAccASNRangeResourceConfig_tagOrder(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag3_tag2_tag1")
		},
		ExpectedTagCount: 3,
		CheckDestroy:     testutil.CheckASNRangeDestroy,
	})
}

// Configuration functions for tag tests.

func testAccASNRangeResourceConfig_tagLifecycle(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_tag" "tag1" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag2" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_tag" "tag3" {
  name = %[9]q
  slug = %[10]q
}
`, name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn_range" "test" {
  name  = %[1]q
  slug  = %[2]q
  rir   = netbox_rir.test.id
  start = "65200"
  end   = "65300"
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
`, name, slug)
	case "tag2_tag3":
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn_range" "test" {
  name  = %[1]q
  slug  = %[2]q
  rir   = netbox_rir.test.id
  start = "65200"
  end   = "65300"
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
`, name, slug)
	default: // "none"
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn_range" "test" {
  name  = %[1]q
  slug  = %[2]q
  rir   = netbox_rir.test.id
  start = "65200"
  end   = "65300"
  tags  = []
}
`, name, slug)
	}
}

func testAccASNRangeResourceConfig_tagOrder(name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagOrder string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_tag" "tag1" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag2" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_tag" "tag3" {
  name = %[9]q
  slug = %[10]q
}

resource "netbox_asn_range" "test" {
  name  = %[1]q
  slug  = %[2]q
  rir   = netbox_rir.test.id
  start = "65301"
  end   = "65401"
`, name, slug, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	var tagConfig string
	switch tagOrder {
	case "tag1_tag2_tag3":
		tagConfig = `
  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    }
  ]`
	case "tag3_tag2_tag1":
		tagConfig = `
  tags = [
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    }
  ]`
	}

	return baseConfig + tagConfig + "\n}\n"
}
