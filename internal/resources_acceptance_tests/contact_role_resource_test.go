package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for contact role resource are in resources_acceptance_tests_customfields package

func TestAccContactRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_contact_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactRoleResourceConfig(name+"-updated", slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name+"-updated"),
				),
			},
		},
	})
}

func TestAccContactRoleResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-full")
	slug := testutil.GenerateSlug(name)
	description := "Full test contact role with all optional fields"
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckContactRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig_full(name, slug, description, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_contact_role.test", "tags.*", tagSlug),
				),
			},
		},
	})
}

func TestAccContactRoleResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-tags")
	slug := testutil.GenerateSlug(name)
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_contact_role",
		ConfigWithoutTags: func() string {
			return testAccContactRoleResourceConfig_tagLifecycle(name, slug, "", "", "", "", "", "", "")
		},
		ConfigWithTags: func() string {
			return testAccContactRoleResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1,tag2", "", "")
		},
		ConfigWithDifferentTags: func() string {
			return testAccContactRoleResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag3", tag3Name, tag3Slug)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
	})
}

func testAccContactRoleResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagSet, tag3Name, tag3Slug string) string {
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

resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
%s
}
`, tagResources, name, slug, tagsList)
}

func TestAccContactRoleResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-tag-order")
	slug := testutil.GenerateSlug(name)
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_contact_role",
		ConfigWithTagsOrderA: func() string {
			return testAccContactRoleResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccContactRoleResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
	})
}

func testAccContactRoleResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
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

resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
  %s
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug, name, slug, tagsOrder)
}

func TestAccContactRoleResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-update")
	slug := testutil.GenerateSlug(name)
	updatedName := testutil.RandomName("test-contact-role-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
				),
			},
			{
				Config: testAccContactRoleResourceConfig(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccConsistency_ContactRole_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-lit")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
			{
				Config:   testAccContactRoleResourceConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
				),
			},
		},
	})
}

func testAccContactRoleResourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccContactRoleResourceConfig_full(name, slug, description, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_contact_role" "test" {
  name        = %q
  slug        = %q
  description = %q
	tags = [netbox_tag.test.slug]
}
`, tagName, tagSlug, name, slug, description)
}

func testAccContactRoleResourceConfig_withTags(name, slug, description, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_contact_role" "test" {
  name        = %q
  slug        = %q
  description = %q
	tags = [netbox_tag.test.slug]
}
`, tagName, tagSlug, name, slug, description)
}

func TestAccContactRoleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-rem")
	slug := testutil.GenerateSlug(name)
	tagName := testutil.RandomName("test-cr-tag")
	tagSlug := testutil.GenerateSlug(tagName)
	description := "Test description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig_withTags(name, slug, description, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "tags.#", "1"),
				),
			},
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_contact_role.test", "description"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccContactRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-contact-role-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactRoleResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyContactRolesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find contact_role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyContactRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete contact_role: %v", err)
					}
					t.Logf("Successfully externally deleted contact_role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccContactRoleResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_contact_role",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact_role" "test" {
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

resource "netbox_contact_role" "test" {
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
