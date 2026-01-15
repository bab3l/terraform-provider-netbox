package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasic(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTagResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")
	color := testutil.ColorOrange
	description := testutil.RandomName("description")
	updatedName := testutil.RandomName("tag-updated")
	updatedSlug := testutil.RandomSlug("tag-updated")
	updatedColor := "2196f3"
	updatedDescription := "Updated test tag description"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceFull(name, slug, color, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tag.test", "color", color),
					resource.TestCheckResourceAttr("netbox_tag.test", "description", description),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccTagResourceFull(name, slug, color, description),
				PlanOnly: true,
			},
			{
				Config: testAccTagResourceFull(updatedName, updatedSlug, updatedColor, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tag.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", updatedSlug),
					resource.TestCheckResourceAttr("netbox_tag.test", "color", updatedColor),
					resource.TestCheckResourceAttr("netbox_tag.test", "description", updatedDescription),
				),
			},
			{
				// Verify no changes after update
				Config:   testAccTagResourceFull(updatedName, updatedSlug, updatedColor, updatedDescription),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTagResource_withObjectTypes(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceWithObjectTypes(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tag.test", "object_types.#", "2"),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccTagResourceWithObjectTypes(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTagResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag-id")
	slug := testutil.RandomSlug("tag-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccTagResourceBasic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTagResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")

	config := testutil.MultiFieldOptionalTestConfig{
		ResourceType: "netbox_tag",
		ResourceName: "netbox_tag",
		ConfigWithFields: func() string {
			return testAccTagResourceFull(name, slug, "2196f3", "Test description")
		},
		BaseConfig: func() string {
			return testAccTagResourceBasic(name, slug)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
		},
	}

	testutil.TestRemoveOptionalFields(t, config)
}

func testAccTagResourceBasic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_Tag_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag-lit")
	slug := testutil.RandomSlug("tag-lit")
	color := testutil.ColorOrange
	description := testutil.RandomName("description")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagConsistencyLiteralNamesConfig(name, slug, color, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tag.test", "color", color),
					resource.TestCheckResourceAttr("netbox_tag.test", "description", description),
				),
			},
			{
				Config:   testAccTagConsistencyLiteralNamesConfig(name, slug, color, description),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
				),
			},
		},
	})
}

func testAccTagConsistencyLiteralNamesConfig(name, slug, color, description string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name        = %q
  slug        = %q
  color       = %q
  description = %q
}
`, name, slug, color, description)
}
func testAccTagResourceFull(name, slug, color, description string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name        = %q
  slug        = %q
  color       = %q
  description = %q
}
`, name, slug, color, description)
}

func testAccTagResourceWithObjectTypes(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name         = %q
  slug         = %q
  object_types = ["dcim.device", "dcim.site"]
}
`, name, slug)
}

func TestAccTagResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag-del")
	slug := testutil.RandomSlug("tag-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasic(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.ExtrasAPI.ExtrasTagsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tag for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasTagsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tag: %v", err)
					}
					t.Logf("Successfully externally deleted tag with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccTagResource_removeOptionalFieldsObjectTypes tests that the object_types field
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccTagResource_removeOptionalFieldsObjectTypes(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tag-rem")
	slug := testutil.RandomSlug("tf-test-tag-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTagCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceConfig_withObjectTypes(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tag.test", "object_types.#", "2"),
				),
			},
			{
				Config: testAccTagResourceConfig_withoutObjectTypes(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_tag.test", "object_types"),
				),
			},
		},
	})
}

func testAccTagResourceConfig_withObjectTypes(name, slug string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_tag" "test" {
  name         = %[1]q
  slug         = %[2]q
  object_types = ["dcim.device", "dcim.site"]
}
`, name, slug)
}

func testAccTagResourceConfig_withoutObjectTypes(name, slug string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_tag" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}
