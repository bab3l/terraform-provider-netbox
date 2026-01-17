package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role")
	slug := testutil.RandomSlug("tf-test-role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRoleResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-update")
	updatedName := testutil.RandomName("tf-test-role-updated")
	slug := testutil.RandomSlug("tf-test-role-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_forUpdate(name, slug, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
					resource.TestCheckResourceAttr("netbox_role.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccRoleResourceConfig_forUpdate(updatedName, slug, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "2000"),
					resource.TestCheckResourceAttr("netbox_role.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccRoleResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-full")
	slug := testutil.RandomSlug("tf-test-role-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated IPAM role description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_full(name, slug, description, 100, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "100"),
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_full(name, slug, description, 100, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
			{
				Config: testAccRoleResourceConfig_fullUpdate(name, slug, updatedDescription, 200, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "200"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_fullUpdate(name, slug, updatedDescription, 200, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRoleResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-tags")
	slug := testutil.RandomSlug("tf-test-role-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRoleResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-tag-order")
	slug := testutil.RandomSlug("tf-test-role-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func testAccRoleResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccRoleResourceConfig_forUpdate(name, slug, description string) string {
	// Toggle weight based on description
	weight := 1000
	if description == testutil.Description2 {
		weight = 2000
	}

	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  weight      = %[4]d
}
`, name, slug, description, weight)
}

func testAccRoleResourceConfig_full(name, slug, description string, weight int, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	cfName := fmt.Sprintf("test_field_%s", strings.ReplaceAll(slug, "-", "_"))
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag2" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_custom_field" "test_field" {
  name         = %[9]q
  object_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  weight      = %[4]d

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "test_value"
    }
  ]
}
`, name, slug, description, weight, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccRoleResourceConfig_fullUpdate(name, slug, description string, weight int, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	cfName := fmt.Sprintf("test_field_%s", strings.ReplaceAll(slug, "-", "_"))
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag2" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_custom_field" "test_field" {
  name         = %[9]q
  object_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  weight      = %[4]d

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "updated_value"
    }
  ]
}
`, name, slug, description, weight, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_role" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_role" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_Role_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-lit")
	slug := testutil.RandomSlug("tf-test-role-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", description),
				),
			},
			{
				Config:   testAccRoleConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
				),
			},
		},
	})
}

func testAccRoleConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-del")
	slug := testutil.RandomSlug("tf-test-role-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamRolesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete role: %v", err)
					}
					t.Logf("Successfully externally deleted role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccRoleResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccRoleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-rem")
	slug := testutil.RandomSlug("tf-test-role-rem")
	description := testutil.RandomName("description")
	tagSlug := testutil.RandomSlug("tf-test-tag")
	const weight = 2000

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_withOptionalFields(name, slug, description, tagSlug, weight),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "2000"),
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "1"),
				),
			},
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckNoResourceAttr("netbox_role.test", "description"),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
					resource.TestCheckNoResourceAttr("netbox_role.test", "tags"),
				),
			},
		},
	})
}

func testAccRoleResourceConfig_withOptionalFields(name, slug, description, tagSlug string, weight int) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %[4]q
  slug = %[4]q
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  weight      = %[5]d
  tags = [
    {
      name = netbox_tag.test.name
      slug = netbox_tag.test.slug
    }
  ]
}
`, name, slug, description, tagSlug, weight)
}

// TestAccRoleResource_Weight tests comprehensive scenarios for role weight field.
// This validates that Optional+Computed int64 fields with proper defaults work correctly.
func TestAccRoleResource_Weight(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-weight")
	slug := testutil.RandomSlug("tf-test-role-weight")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccRoleResourceConfig_withWeight(name, slug, 2000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "weight", "2000"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_withWeight(name, slug, 2000),
				PlanOnly: true,
			},
		},
	})
}

func testAccRoleResourceConfig_withWeight(name, slug string, weight int) string {
	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name   = %q
  slug   = %q
  weight = %d
}
`, name, slug, weight)
}

func TestAccRoleResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_role",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_role" "test" {
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

resource "netbox_role" "test" {
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
