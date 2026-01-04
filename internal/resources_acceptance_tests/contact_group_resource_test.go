package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for contact group resource are in resources_acceptance_tests_customfields package

func TestAccContactGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_contact_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactGroupResourceConfig(name+"-updated", slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name+"-updated"),
				),
			},
		},
	})
}

func TestAccContactGroupResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-update")
	slug := testutil.GenerateSlug(name)
	updatedName := testutil.RandomName("test-contact-group-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
				),
			},
			{
				Config: testAccContactGroupResourceConfig(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccContactGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cg-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccConsistency_ContactGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-lit")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccContactGroupResourceConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
				),
			},
		},
	})
}

func testAccContactGroupResourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccContactGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-del")
	slug := testutil.GenerateSlug(name)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyContactGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find contact_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyContactGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete contact_group: %v", err)
					}
					t.Logf("Successfully externally deleted contact_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
