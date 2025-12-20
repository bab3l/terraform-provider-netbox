package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantResource_basic(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant")
	slug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccTenantResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-full")
	slug := testutil.RandomSlug("tf-test-tenant-full")
	description := "Test tenant with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant.test", "description", description),
				),
			},
		},
	})
}

func TestAccTenantResource_update(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-update")
	slug := testutil.RandomSlug("tf-test-tenant-upd")
	updatedName := testutil.RandomName("tf-test-tenant-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
				),
			},
			{
				Config: testAccTenantResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", updatedName),
				),
			},
		},
	})
}

// testAccTenantResourceConfig_basic returns a basic test configuration.
func testAccTenantResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

// testAccTenantResourceConfig_full returns a test configuration with all fields.
func testAccTenantResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccTenantResource_import(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-import")
	slug := testutil.RandomSlug("tf-test-tenant-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tenant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTenantResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_Tenant(t *testing.T) {
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTenantGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", tenantName),
					resource.TestCheckResourceAttr("netbox_tenant.test", "group", groupName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug),
			},
		},
	})
}

func testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
  group = netbox_tenant_group.test.name
}
`, groupName, groupSlug, tenantName, tenantSlug)
}
