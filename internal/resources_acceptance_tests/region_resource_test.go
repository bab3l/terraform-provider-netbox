package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegionResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region")
	slug := testutil.RandomSlug("tf-test-region")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccRegionResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-full")
	slug := testutil.RandomSlug("tf-test-region-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_region.test", "description", description),
				),
			},
		},
	})
}

func TestAccRegionResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-update")
	slug := testutil.RandomSlug("tf-test-region-upd")
	updatedName := testutil.RandomName("tf-test-region-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
				),
			},
			{
				Config: testAccRegionResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccRegionResource_withParent(t *testing.T) {
	t.Parallel()

	parentName := testutil.RandomName("tf-test-region-parent")
	parentSlug := testutil.RandomSlug("tf-test-region-prnt")
	childName := testutil.RandomName("tf-test-region-child")
	childSlug := testutil.RandomSlug("tf-test-region-chld")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(childSlug)
	cleanup.RegisterRegionCleanup(parentSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.parent", "id"),
					resource.TestCheckResourceAttr("netbox_region.parent", "name", parentName),
					resource.TestCheckResourceAttrSet("netbox_region.child", "id"),
					resource.TestCheckResourceAttr("netbox_region.child", "name", childName),
					resource.TestCheckResourceAttrPair("netbox_region.child", "parent", "netbox_region.parent", "id"),
				),
			},
		},
	})
}

func TestAccRegionResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-import")
	slug := testutil.RandomSlug("tf-test-region-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_region.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-tags")
	slug := testutil.RandomSlug("tf-test-region-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRegionResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRegionResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccRegionResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRegionResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-tag-order")
	slug := testutil.RandomSlug("tf-test-region-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRegionResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_region.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func testAccRegionResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_region" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccRegionResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_region" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_region" "parent" {
  name = %q
  slug = %q
}

resource "netbox_region" "child" {
  name   = %q
  slug   = %q
  parent = netbox_region.parent.id
}
`, parentName, parentSlug, childName, childSlug)
}

func testAccRegionResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_region" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRegionResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_region" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_Region_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-lit")
	slug := testutil.RandomSlug("tf-test-region-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_region.test", "description", description),
				),
			},
			{
				Config:   testAccRegionResourceConfig_full(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
				),
			},
		},
	})
}

func testAccRegionResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccRegionResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-del")
	slug := testutil.RandomSlug("tf-test-region-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimRegionsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find region for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimRegionsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete region: %v", err)
					}
					t.Logf("Successfully externally deleted region with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRegionResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	childName := testutil.RandomName("tf-test-region-child")
	childSlug := testutil.RandomSlug("tf-test-region-child")
	parentName := testutil.RandomName("tf-test-region-parent")
	parentSlug := testutil.RandomSlug("tf-test-region-parent")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(childSlug)
	cleanup.RegisterRegionCleanup(parentSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.child", "name", childName),
					resource.TestCheckResourceAttrSet("netbox_region.child", "parent"),
				),
			},
			{
				Config: testAccRegionResourceConfig_childOnly(parentName, parentSlug, childName, childSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.child", "name", childName),
					resource.TestCheckNoResourceAttr("netbox_region.child", "parent"),
				),
			},
		},
	})
}

func testAccRegionResourceConfig_childOnly(parentName, parentSlug, childName, childSlug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "parent" {
  name = %q
  slug = %q
}

resource "netbox_region" "child" {
  name = %q
  slug = %q
}
`, parentName, parentSlug, childName, childSlug)
}

func TestAccRegionResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_region",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_region" "test" {
  slug = "test-region"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_region" "test" {
  name = "Test Region"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
