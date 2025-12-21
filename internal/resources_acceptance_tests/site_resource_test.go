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

func TestAccSiteResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site")
	slug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccSiteResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-full")
	slug := testutil.RandomSlug("tf-test-site-full")
	description := "Test site with all fields"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_site.test", "description", description),
				),
			},
		},
	})
}

func TestAccSiteResource_update(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-update")
	slug := testutil.RandomSlug("tf-test-site-upd")
	updatedName := testutil.RandomName("tf-test-site-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
				),
			},
			{
				Config: testAccSiteResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccSiteResource_import(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-import")
	slug := testutil.RandomSlug("tf-test-site-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_Site(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	regionName := testutil.RandomName("region")
	regionSlug := testutil.RandomSlug("region")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "region", regionName),
					resource.TestCheckResourceAttr("netbox_site.test", "group", groupName),
					resource.TestCheckResourceAttr("netbox_site.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccSiteResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name        = %q
  slug        = %q
  status      = "active"
  description = %q
}
`, name, slug, description)
}

func testAccSiteResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_site_group" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_tenant" "test" {
  name = "%[7]s"
  slug = "%[8]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  region = netbox_region.test.name
  group = netbox_site_group.test.name
  tenant = netbox_tenant.test.name
}
`, siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug)
}
