package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for cluster group resource are in resources_acceptance_tests_customfields package

func TestAccClusterGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group-full")
	slug := testutil.RandomSlug("tf-test-cluster-group-full")
	description := "Full test cluster group with all optional fields"
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_full(name, slug, description, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "description", description),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "tags.0.name", tagName),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "tags.0.slug", tagSlug),
				),
			},
		},
	})
}

func TestAccClusterGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group")
	slug := testutil.RandomSlug("tf-test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_cluster_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_ClusterGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group-lit")
	slug := testutil.RandomSlug("tf-test-cluster-group-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccClusterGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
				),
			},
		},
	})
}

func TestAccClusterGroupResource_update(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	name := testutil.RandomName("tf-test-cluster-group")
	slug := testutil.RandomSlug("tf-test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				Config: testAccClusterGroupResourceConfig_basic(name+"-updated", slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccClusterGroupResourceConfig_full(name, slug, description, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster_group" "test" {
  name        = %q
  slug        = %q
  description = %q
  tags = [
    {
      name = netbox_tag.test.name
      slug = netbox_tag.test.slug
    }
  ]
}
`, tagName, tagSlug, name, slug, description)
}

func testAccClusterGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccClusterGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-cluster-group-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VirtualizationAPI.VirtualizationClusterGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find cluster_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationClusterGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cluster_group: %v", err)
					}
					t.Logf("Successfully externally deleted cluster_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccClusterGroupResource_removeDescription(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group-desc")
	slug := testutil.RandomSlug("tf-test-cluster-group-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_cluster_group",
		BaseConfig: func() string {
			return testAccClusterGroupResourceConfig_basic(name, slug)
		},
		ConfigWithFields: func() string {
			return testAccClusterGroupResourceConfig_withDescription(
				name,
				slug,
				"Test description",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
		},
		CheckDestroy: testutil.CheckClusterGroupDestroy,
	})
}

func testAccClusterGroupResourceConfig_withDescription(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccClusterGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_cluster_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_cluster_group" "test" {
  # name missing
  slug = "test-cluster-group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_cluster_group" "test" {
  name = "Test Cluster Group"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

// =============================================================================
// STANDARDIZED TAG TESTS (using helpers)
// =============================================================================

// TestAccClusterGroupResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccClusterGroupResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-cluster-group-tag")
	slug := testutil.RandomSlug("tf-cluster-group-tag")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_cluster_group",
		ConfigWithoutTags: func() string {
			return testAccClusterGroupResourceConfig_tagLifecycle(name, slug, "", "", "", "", "", "")
		},
		ConfigWithTags: func() string {
			return testAccClusterGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1,tag2", "")
		},
		ConfigWithDifferentTags: func() string {
			return testAccClusterGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag3", tag3Name+"|"+tag3Slug)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
	})
}

func testAccClusterGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagSet, tag3Info string) string {
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
	if tag3Info != "" {
		parts := tag3Info
		idx := 0
		for i, ch := range parts {
			if ch == '|' {
				idx = i
				break
			}
		}
		tag3Name := parts[:idx]
		tag3Slug := parts[idx+1:]
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
			tagsList = tagsDoubleNested
		case caseTag3:
			tagsList = tagsSingleNested
		default:
			tagsList = tagsEmpty
		}
	} else {
		tagsList = tagsEmpty
	}

	return fmt.Sprintf(`
%s
resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q
  %s
}
`, tagResources, name, slug, tagsList)
}

// TestAccClusterGroupResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccClusterGroupResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-cluster-group-order")
	slug := testutil.RandomSlug("tf-cluster-group-order")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_cluster_group",
		ConfigWithTagsOrderA: func() string {
			return testAccClusterGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccClusterGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
	})
}

func testAccClusterGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	tagsOrder := tagsDoubleNested
	if !tag1First {
		tagsOrder = tagsDoubleNestedReversed
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

resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q
  %s
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug, name, slug, tagsOrder)
}
