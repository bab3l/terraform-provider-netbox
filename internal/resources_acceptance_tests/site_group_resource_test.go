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
			{
				Config:   testAccSiteGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSiteGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group-full")
	slug := testutil.RandomSlug("tf-test-sg-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated site group description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site_group.test", "description", description),
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccSiteGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
			{
				Config: testAccSiteGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_site_group.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccSiteGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSiteGroupResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group-tags")
	slug := testutil.RandomSlug("tf-test-sg-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccSiteGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccSiteGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccSiteGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccSiteGroupResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-group-tag-order")
	slug := testutil.RandomSlug("tf-test-sg-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccSiteGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site_group.test", "tags.*", tag2Slug),
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
				Config:   testAccSiteGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccSiteGroupResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", updatedName),
				),
			},
			{
				Config:   testAccSiteGroupResourceConfig_basic(updatedName, slug),
				PlanOnly: true,
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
				Config:   testAccSiteGroupResourceConfig_import(name, slug),
				PlanOnly: true,
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
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site_group.test", "description", description),
				),
			},
			{
				Config:   testAccSiteGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
				),
			},
		},
	})
}

func testAccSiteGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[8]q
	object_types = ["dcim.sitegroup"]
	type         = "text"
}

resource "netbox_site_group" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "test_value"
		}
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccSiteGroupResourceConfig_fullUpdate(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[8]q
	object_types = ["dcim.sitegroup"]
	type         = "text"
}

resource "netbox_site_group" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "updated_value"
		}
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
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

func TestAccSiteGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_site_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_site_group" "test" {
  slug = "test-group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_site_group" "test" {
  name = "Test Group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccSiteGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_site_group" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccSiteGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_site_group" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}
