package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "slug"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", "e41e22"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-update")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", "e41e22"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_inventory_item_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccInventoryItemRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccInventoryItemRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccInventoryItemRoleResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name        = %q
  slug        = %q
  color       = "e41e22"
  description = "Test role description"
}
`, name, testutil.RandomSlug("role"))
}

func testAccInventoryItemRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[5]s"
  slug = %[5]q
}

resource "netbox_inventory_item_role" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccInventoryItemRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_inventory_item_role" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccInventoryItemRoleResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-tags")
	slug := testutil.RandomSlug("role")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-tag-order")
	slug := testutil.RandomSlug("role")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_inventory_item_role.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccConsistency_InventoryItemRole_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-lit")
	slug := testutil.RandomSlug("tf-test-inv-role-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", slug),
				),
			},
			{
				Config:   testAccInventoryItemRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-rem")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_inventory_item_role.test", "description"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-ext-del")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInventoryItemRolesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find inventory_item_role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInventoryItemRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete inventory_item_role: %v", err)
					}
					t.Logf("Successfully externally deleted inventory_item_role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccInventoryItemRoleResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_inventory_item_role",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_inventory_item_role" "test" {
  # name missing
  slug = "test-role"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_inventory_item_role" "test" {
  name = "Test Role"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
