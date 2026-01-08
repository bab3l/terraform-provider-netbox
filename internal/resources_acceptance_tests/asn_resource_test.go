package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNResource_basic(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	// Generate a random ASN in the private range (64512-65534)
	asn := int64(acctest.RandIntRange(64512, 65534))
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
			{
				ResourceName:            "netbox_asn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})
}

func TestAccASNResource_full(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate a random ASN in the private range (64512-65534)
	asn := int64(acctest.RandIntRange(64512, 65534))
	description := testutil.RandomName("description")
	updatedDescription := "Updated ASN description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, description, testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", description),
					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "tenant"),
				),
			},
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, updatedDescription, testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccASNResource_IDPreservation(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-id")
	rirSlug := testutil.RandomSlug("tf-test-rir-id")
	asn := int64(acctest.RandIntRange(64512, 65000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
		},
	})

}

func TestAccASNResource_update(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-update")
	rirSlug := testutil.RandomSlug("tf-test-rir-update")
	asn := int64(acctest.RandIntRange(64512, 65534))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_update(rirName, rirSlug, asn, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccASNResourceConfig_update(rirName, rirSlug, asn, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccASNResource_external_deletion(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")
	asn := int64(acctest.RandIntRange(64512, 65534))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// Safe conversion with bounds check
					if asn < 0 || asn > 4294967295 { // Max ASN value (32-bit)
						t.Fatalf("ASN value %d out of valid range", asn)
					}
					//nolint:gosec // Safe conversion - bounds checked above
					items, _, err := client.IpamAPI.IpamAsnsList(context.Background()).Asn([]int32{int32(asn)}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find asn for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAsnsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete asn: %v", err)
					}
					t.Logf("Successfully externally deleted asn with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccASNResourceConfig_basic(rirName, rirSlug string, asn int64) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
}
`, rirName, rirSlug, asn)
}

func testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug string, asn int64, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  tenant      = netbox_tenant.test.id
  description = %q
  comments    = %q
}
`, rirName, rirSlug, tenantName, tenantSlug, asn, description, comments)
}

func testAccASNResourceConfig_update(rirName, rirSlug string, asn int64, description string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  description = %q
}
`, rirName, rirSlug, asn, description)
}

// TestAccConsistency_ASN_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ASN_LiteralNames(t *testing.T) {
	t.Parallel()

	asn := int64(65100)
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "rir", rirSlug),
					resource.TestCheckResourceAttr("netbox_asn.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccASNConsistencyLiteralNamesConfig(asn int64, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_asn" "test" {
  asn = %[1]d
  # Use literal string names to mimic existing user state
  rir = "%[3]s"
  tenant = "%[4]s"
  depends_on = [netbox_rir.test, netbox_tenant.test]
}
`, asn, rirName, rirSlug, tenantName, tenantSlug)
}

// NOTE: Custom field tests for ASN resource are in resources_acceptance_tests_customfields package
