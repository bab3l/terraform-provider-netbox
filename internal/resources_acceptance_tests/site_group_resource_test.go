package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group")
	slug := testutil.RandomSlug("tf-test-sg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
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
			{
				Config:   testAccSiteGroupResourceConfig_import(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSiteGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group-id")
	slug := testutil.RandomSlug("tf-test-sg-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
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

func testAccSiteGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_SiteGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group-lit")
	slug := testutil.RandomSlug("tf-test-site-group-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
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
			{
				Config:   testAccSiteGroupResourceConfig_full(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
				),
			},
		},
	})
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

func TestAccSiteGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group-del")
	slug := testutil.RandomSlug("tf-test-site-group-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimSiteGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find site_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimSiteGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete site_group: %v", err)
					}
					t.Logf("Successfully externally deleted site_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSiteGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-sg-rem")
	slug := testutil.RandomSlug("tf-test-sg-rem")
	parentName := testutil.RandomName("tf-test-sg-parent")
	parentSlug := testutil.RandomSlug("tf-test-sg-parent")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)
	cleanup.RegisterSiteGroupCleanup(parentSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_withParent(name, slug, parentName, parentSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "parent"),
				),
			},
			{
				Config: testAccSiteGroupResourceConfig_detached(name, slug, parentName, parentSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_site_group.test", "parent"),
				),
			},
		},
	})
}

func testAccSiteGroupResourceConfig_withParent(name, slug, parentName, parentSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_site_group" "test" {
  name   = %q
  slug   = %q
  parent = netbox_site_group.parent.id
}
`, parentName, parentSlug, name, slug)
}

func testAccSiteGroupResourceConfig_detached(name, slug, parentName, parentSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_site_group" "test" {
  name   = %q
  slug   = %q
}
`, parentName, parentSlug, name, slug)
}
