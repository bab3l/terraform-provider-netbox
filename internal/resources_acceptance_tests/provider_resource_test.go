package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider")
	slug := testutil.RandomSlug("tf-test-provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccProviderResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-full")
	slug := testutil.RandomSlug("tf-test-provider-full")
	description := testutil.RandomName("description")
	comments := testutil.RandomName("comments")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_full(name, slug, description, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_provider.test", "description", description),
					resource.TestCheckResourceAttr("netbox_provider.test", "comments", comments),
				),
			},
		},
	})
}

func TestAccProviderResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-update")
	slug := testutil.RandomSlug("tf-test-provider-update")
	updatedName := testutil.RandomName("tf-test-provider-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
				),
			},
			{
				Config: testAccProviderResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccProviderResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider")
	slug := testutil.RandomSlug("tf-test-provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccProviderResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccProviderResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-tags")
	slug := testutil.RandomSlug("tf-test-provider-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccProviderResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccProviderResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccProviderResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccProviderResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-tag-order")
	slug := testutil.RandomSlug("tf-test-provider-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccProviderResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_provider.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func testAccProviderResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccProviderResourceConfig_full(name, slug, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name        = %q
  slug        = %q
  description = %q
  comments    = %q
}
`, name, slug, description, comments)
}

func testAccProviderResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleNested
	case caseTag3:
		tagsConfig = tagsSingleNested
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

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccProviderResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleNestedReversed
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

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_Provider_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("provider")
	slug := testutil.RandomSlug("provider")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccProviderResourceConfig_basic(name, slug),
			},
		},
	})
}

func TestAccProviderResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-ext-del")
	slug := testutil.RandomSlug("provider-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List providers filtered by slug
					items, _, err := client.CircuitsAPI.CircuitsProvidersList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find provider for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsProvidersDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete provider: %v", err)
					}
					t.Logf("Successfully externally deleted provider with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProviderResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-provider-remove")
	slug := testutil.RandomSlug("tf-test-provider-remove")
	description := "Description"
	comments := "Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderResourceConfig_full(name, slug, description, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider.test", "description", description),
					resource.TestCheckResourceAttr("netbox_provider.test", "comments", comments),
				),
			},
			{
				Config: testAccProviderResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_provider.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_provider.test", "comments"),
				),
			},
		},
	})
}

func TestAccProviderResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_provider",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider" "test" {
  # name missing
  slug = "test-provider"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider" "test" {
  name = "Test Provider"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
