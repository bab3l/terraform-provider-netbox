package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for contact group resource are in resources_acceptance_tests_customfields package

func TestAccContactGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-full")
	slug := testutil.GenerateSlug(name)
	description := "Full test contact group with all optional fields"
	parentName := testutil.RandomName("test-cg-parent")
	parentSlug := testutil.GenerateSlug(parentName)
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)
	cleanup.RegisterContactGroupCleanup(parentSlug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig_full(name, slug, description, parentName, parentSlug, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "description", description),
					resource.TestCheckResourceAttrPair(
						"netbox_contact_group.test", "parent",
						"netbox_contact_group.parent", "id",
					),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_contact_group.test", "tags.*", tagSlug),
				),
			},
		},
	})
}

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

func TestAccContactGroupResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-tags")
	slug := testutil.GenerateSlug(name)
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_contact_group",
		ConfigWithoutTags: func() string {
			return testAccContactGroupResourceConfig_tagLifecycle(name, slug, "", "", "", "", "", "", "")
		},
		ConfigWithTags: func() string {
			return testAccContactGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1,tag2", "", "")
		},
		ConfigWithDifferentTags: func() string {
			return testAccContactGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag3", tag3Name, tag3Slug)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
	})
}

func testAccContactGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagSet, tag3Name, tag3Slug string) string {
	tagResources := ""
	tagsList := ""

	if tag1Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}
`, tag1Name, tag1Slug)
	}
	if tag2Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}
`, tag2Name, tag2Slug)
	}
	if tag3Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag3" {
  name = %q
  slug = %q
}
`, tag3Name, tag3Slug)
	}

	if tagSet != "" {
		switch tagSet {
		case caseTag1Tag2:
			tagsList = tagsDoubleSlug
		case caseTag3:
			tagsList = tagsSingleSlug
		default:
			tagsList = tagsEmpty
		}
	} else {
		tagsList = tagsEmpty
	}

	return fmt.Sprintf(`
%s

resource "netbox_contact_group" "test" {
  name = %q
  slug = %q
%s
}
`, tagResources, name, slug, tagsList)
}

func TestAccContactGroupResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-tag-order")
	slug := testutil.GenerateSlug(name)
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_contact_group",
		ConfigWithTagsOrderA: func() string {
			return testAccContactGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccContactGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
	})
}

func testAccContactGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	tagsOrder := tagsDoubleSlug
	if !tag1First {
		tagsOrder = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

resource "netbox_contact_group" "test" {
  name = %q
  slug = %q
  %s
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug, name, slug, tagsOrder)
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

func testAccContactGroupResourceConfig_full(name, slug, description, parentName, parentSlug, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_contact_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_contact_group" "test" {
  name        = %q
  slug        = %q
  description = %q
  parent      = netbox_contact_group.parent.id
	tags = [netbox_tag.test.slug]
}
`, tagName, tagSlug, parentName, parentSlug, name, slug, description)
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

func TestAccContactGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-group-rem")
	slug := testutil.GenerateSlug(name)
	parentName := testutil.RandomName("test-cg-parent")
	parentSlug := testutil.GenerateSlug(parentName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(slug)
	cleanup.RegisterContactGroupCleanup(parentSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig_withParent(name, slug, parentName, parentSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_contact_group.test", "parent"),
				),
			},
			{
				Config: testAccContactGroupResourceConfig_detached(name, slug, parentName, parentSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_contact_group.test", "parent"),
				),
			},
		},
	})
}

func testAccContactGroupResourceConfig_withParent(name, slug, parentName, parentSlug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_contact_group" "test" {
  name   = %q
  slug   = %q
  parent = netbox_contact_group.parent.id
}
`, parentName, parentSlug, name, slug)
}

func testAccContactGroupResourceConfig_detached(name, slug, parentName, parentSlug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_contact_group" "test" {
  name   = %q
  slug   = %q
}
`, parentName, parentSlug, name, slug)
}

func TestAccContactGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_contact_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact_group" "test" {
  # name missing
  slug = "test-group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact_group" "test" {
  name = "Test Group"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
