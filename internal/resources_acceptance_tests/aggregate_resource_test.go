package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateResource_basic(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
				),
			},
			{
				ResourceName:            "netbox_aggregate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})
}

func TestAccAggregateResource_update(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-update")
	rirSlug := testutil.RandomSlug("tf-test-rir-update")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_withDescription(rirName, rirSlug, prefix, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccAggregateResourceConfig_withDescription(rirName, rirSlug, prefix, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccAggregateResource_full(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-full")
	rirSlug := testutil.RandomSlug("tf-test-rir-full")
	tenantName := testutil.RandomName("tf-test-tenant-full")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-full")
	prefix := testutil.RandomIPv4Prefix()
	description := testutil.RandomName("description")
	updatedDescription := "Updated aggregate description"
	dateAdded := "2024-01-15"
	updatedDate := "2024-06-20"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, testutil.Comments, dateAdded),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", testutil.Comments),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "date_added", dateAdded),
				),
			},
			{
				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, updatedDescription, testutil.Comments, updatedDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "date_added", updatedDate),
				),
			},
		},
	})
}

func TestAccAggregateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-id")
	rirSlug := testutil.RandomSlug("tf-test-rir-id")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
				),
			},
		},
	})

}

func TestAccAggregateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-optional")
	rirSlug := testutil.RandomSlug("tf-test-rir-optional")
	tenantName := testutil.RandomName("tf-test-tenant-optional")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-optional")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_aggregate",
		BaseConfig: func() string {
			return testAccAggregateResourceConfig_withTenant(rirName, rirSlug, tenantName, tenantSlug, prefix)
		},
		ConfigWithFields: func() string {
			return testAccAggregateResourceConfig_full(
				rirName, rirSlug, tenantName, tenantSlug,
				prefix,
				"Test description",
				"Test comments",
				"2024-01-15",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
			"date_added":  "2024-01-15",
			// Note: tenant is not included as it requires TestCheckResourceAttrSet verification
			// which the test helper doesn't support for ID-based references
		},
		RequiredFields: map[string]string{
			"prefix": prefix,
		},
	})
}

func testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix = %q
  rir    = netbox_rir.test.id
}
`, rirName, rirSlug, prefix)
}

func testAccAggregateResourceConfig_withDescription(rirName, rirSlug, prefix, description string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix      = %q
  rir         = netbox_rir.test.id
  description = %q
}
`, rirName, rirSlug, prefix, description)
}

func testAccAggregateResourceConfig_withTenant(rirName, rirSlug, tenantName, tenantSlug, prefix string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix = %q
  rir    = netbox_rir.test.id
}
`, rirName, rirSlug, tenantName, tenantSlug, prefix)
}

func testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix      = %q
  rir         = netbox_rir.test.id
  tenant      = netbox_tenant.test.id
  description = %q
  comments    = %q
  date_added  = %q
}
`, rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded)
}

func TestAccConsistency_Aggregate(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_aggregate" "test" {
  prefix = "%[1]s"
  rir = netbox_rir.test.id
  tenant = netbox_tenant.test.name
}
`, prefix, rirName, rirSlug, tenantName, tenantSlug)
}

// TestAccConsistency_Aggregate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_Aggregate_LiteralNames(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_aggregate" "test" {
  prefix = "%[1]s"
  # Use literal string names to mimic existing user state
  rir = "%[3]s"
  tenant = "%[4]s"
  depends_on = [netbox_rir.test, netbox_tenant.test]
}

`, prefix, rirName, rirSlug, tenantName, tenantSlug)

}

func TestAccAggregateResource_externalDeletion(t *testing.T) {
	t.Parallel()
	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterAggregateCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamAggregatesList(context.Background()).Prefix(prefix).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find aggregate for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAggregatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete aggregate: %v", err)
					}
					t.Logf("Successfully externally deleted aggregate with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// NOTE: Custom field tests for aggregate resource are in resources_acceptance_tests_customfields package

// Enhancement 2: Test import with full optional fields
// This verifies that import correctly populates all fields from the API.
func TestAccAggregateResource_importPreservesOptionalFields(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	description := testutil.RandomName("description")
	comments := testutil.RandomName("comments")
	dateAdded := "2024-01-01"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create resource with ALL optional fields
			{
				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "date_added", dateAdded),
				),
			},
			// Step 2: Import with SAME full config - verifies import preserves all fields
			{
				Config:            testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded),
				ResourceName:      "netbox_aggregate.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Note: tenant format may change (ID->name/slug), rir is always lookup
				ImportStateVerifyIgnore: []string{"rir", "tenant"},
			},
			// Step 3: Verify no changes after import
			{
				Config:   testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded),
				PlanOnly: true,
			},
		},
	})
}
