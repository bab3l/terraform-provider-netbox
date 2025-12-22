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

func TestAccSiteGroupResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-group")
	slug := testutil.RandomSlug("tf-test-sg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccSiteGroupResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-group-full")
	slug := testutil.RandomSlug("tf-test-sg-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccSiteGroupResource_update(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-group-update")
	slug := testutil.RandomSlug("tf-test-sg-upd")
	updatedName := testutil.RandomName("tf-test-site-group-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
				),
			},
			{
				Config: testAccSiteGroupResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccSiteGroupResource_import(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-group")
	slug := testutil.RandomSlug("tf-test-sg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_site_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSiteGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccSiteGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func testAccSiteGroupResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}
