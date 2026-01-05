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

// NOTE: Custom field tests for tenant group resource are in resources_acceptance_tests_customfields package

func TestAccTenantGroupResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-group")
	slug := testutil.RandomSlug("tf-test-tg")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccTenantGroupResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-group-full")
	slug := testutil.RandomSlug("tf-test-tg-full")
	description := "Test tenant group with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_full(name, slug, description),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-group-update")
	slug := testutil.RandomSlug("tf-test-tg-upd")
	updatedName := testutil.RandomName("tf-test-tenant-group-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccTenantGroupResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", updatedName),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(updatedName, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-group-import")
	slug := testutil.RandomSlug("tf-test-tenant-group-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tenant_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTenantGroupResourceConfig_import(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-id")
	slug := testutil.RandomSlug("tf-test-tg-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

// testAccTenantGroupResourceConfig_basic returns a basic test configuration.
func testAccTenantGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_TenantGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-lit")
	slug := testutil.RandomSlug("tf-test-tenant-group-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
				),
			},
			{
				Config:   testAccTenantGroupConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
				),
			},
		},
	})
}

func testAccTenantGroupConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
} // testAccTenantGroupResourceConfig_full returns a test configuration with all fields.
func testAccTenantGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func testAccTenantGroupResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccTenantGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-del")
	slug := testutil.RandomSlug("tf-test-tenant-group-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyTenantGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tenant_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyTenantGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tenant_group: %v", err)
					}
					t.Logf("Successfully externally deleted tenant_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
