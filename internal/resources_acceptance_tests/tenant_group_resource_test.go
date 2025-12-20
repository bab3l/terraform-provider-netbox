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

func TestAccTenantGroupResource_basic(t *testing.T) {
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
		},
	})
}

func TestAccTenantGroupResource_update(t *testing.T) {
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
				Config: testAccTenantGroupResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccTenantGroupResource_import(t *testing.T) {
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

// testAccTenantGroupResourceConfig_full returns a test configuration with all fields.
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
