package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for cluster resource are in resources_acceptance_tests_customfields package

func TestAccClusterResource_basic(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
				),
			},
		},
	})
}

func TestAccClusterResource_full(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-full")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-full")
	clusterName := testutil.RandomName("tf-test-cluster-full")
	description := testutil.RandomName("description")
	comments := testutil.RandomName("comments")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, description, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "description", description),
					resource.TestCheckResourceAttr("netbox_cluster.test", "comments", comments),
				),
			},
		},
	})
}

func TestAccClusterResource_update(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-update")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-update")
	clusterName := testutil.RandomName("tf-test-cluster-update")
	updatedName := testutil.RandomName("tf-test-cluster-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterCleanup(updatedName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
				),
			},
			{
				Config: testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, updatedName, "Updated description", "Updated comments"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_cluster.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccClusterResource_import(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-import")
	clusterTypeSlug := clusterTypeName
	clusterName := testutil.RandomName("tf-test-cluster-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
				),
			},
			{
				ResourceName:            "netbox_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type"},
			},
		},
	})
}

// TestAccConsistency_Cluster_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_Cluster_LiteralNames(t *testing.T) {
	t.Parallel()

	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterGroupCleanup(groupSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConsistencyLiteralNamesConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
					resource.TestCheckResourceAttr("netbox_cluster.test", "type", clusterTypeSlug),
					resource.TestCheckResourceAttr("netbox_cluster.test", "group", groupSlug),
					resource.TestCheckResourceAttr("netbox_cluster.test", "site", siteName),
					resource.TestCheckResourceAttr("netbox_cluster.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccClusterConsistencyLiteralNamesConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccClusterResource_IDPreservation(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("clu-type-id")
	clusterTypeSlug := testutil.RandomSlug("clu-type-id")
	clusterName := testutil.RandomName("clu-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
				),
			},
		},
	})
}

func testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName)
}

func testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name        = %q
  type        = netbox_cluster_type.test.id
  status      = "active"
  description = %q
  comments    = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, description, comments)
}

func testAccClusterConsistencyLiteralNamesConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_cluster_group" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_site" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_tenant" "test" {
  name = "%[8]s"
  slug = "%[9]s"
}

resource "netbox_cluster" "test" {
  name = "%[1]s"
  # Use literal string names to mimic existing user state
  type = "%[3]s"
  group = "%[5]s"
  site = "%[6]s"
  tenant = "%[8]s"
  depends_on = [netbox_cluster_type.test, netbox_cluster_group.test, netbox_site.test, netbox_tenant.test]
}
`, clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug)
}

func TestAccClusterResource_externalDeletion(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List clusters filtered by name
					items, _, err := client.VirtualizationAPI.VirtualizationClustersList(context.Background()).NameIc([]string{clusterName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find cluster for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationClustersDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cluster: %v", err)
					}
					t.Logf("Successfully externally deleted cluster with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccClusterResource_removeOptionalFields tests that removing previously set group, tenant, and site fields correctly sets them to null.
// This addresses the bug where removing a nullable reference field from the configuration would not clear it in NetBox,
// causing "Provider produced inconsistent result after apply" errors.
func TestAccClusterResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-remove")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-remove")
	clusterName := testutil.RandomName("tf-test-cluster-remove")
	siteName := testutil.RandomName("test-site-cluster")
	siteSlug := testutil.GenerateSlug(siteName)
	tenantName := testutil.RandomName("test-tenant-cluster")
	tenantSlug := testutil.GenerateSlug(tenantName)
	groupName := testutil.RandomName("test-group-cluster")
	groupSlug := testutil.GenerateSlug(groupName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterClusterGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with group, tenant, and site
			{
				Config: testAccClusterResourceConfig_withAllFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "group"),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "site"),
				),
			},
			// Step 2: Remove group, tenant, and site - should set to null
			{
				Config: testAccClusterResourceConfig_withoutFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_cluster.test", "group"),
					resource.TestCheckNoResourceAttr("netbox_cluster.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_cluster.test", "site"),
				),
			},
			// Step 3: Re-add group, tenant, and site - should work without errors
			{
				Config: testAccClusterResourceConfig_withAllFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster.test", "status", "staging"),
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "site"),
				),
			},
		},
	})
}

func testAccClusterResourceConfig_withAllFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_site" "test" {
  name   = %[4]q
  slug   = %[5]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_cluster_group" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_cluster" "test" {
  name   = %[3]q
  type   = netbox_cluster_type.test.id
  status = "staging"
  group  = netbox_cluster_group.test.id
  tenant = netbox_tenant.test.id
  site   = netbox_site.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug)
}

func testAccClusterResourceConfig_withoutFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_site" "test" {
  name   = %[4]q
  slug   = %[5]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_cluster_group" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_cluster" "test" {
  name   = %[3]q
  type   = netbox_cluster_type.test.id
  status = "active"
  # group, tenant, and site removed - should set to null
}
`, clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, tenantName, tenantSlug, groupName, groupSlug)
}

func TestAccClusterResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-desc")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-desc")
	clusterName := testutil.RandomName("tf-test-cluster-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_cluster",
		BaseConfig: func() string {
			return testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName)
		},
		ConfigWithFields: func() string {
			return testAccClusterResourceConfig_full(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		CheckDestroy: testutil.CheckClusterDestroy,
	})
}
