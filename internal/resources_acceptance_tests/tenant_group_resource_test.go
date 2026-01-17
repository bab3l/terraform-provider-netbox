package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for tenant group resource are in resources_acceptance_tests_customfields package

func TestAccTenantGroupResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-group")
	slug := testutil.RandomSlug("tf-test-tg")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-group-full")
	slug := testutil.RandomSlug("tf-test-tg-full")
	description := "Test tenant group with all fields"
	updatedDescription := "Updated tenant group description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
			{
				Config: testAccTenantGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-tags")
	slug := testutil.RandomSlug("tf-test-tg-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccTenantGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccTenantGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccTenantGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccTenantGroupResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-tag-order")
	slug := testutil.RandomSlug("tf-test-tg-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccTenantGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_tenant_group.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func TestAccTenantGroupResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-group-update")
	slug := testutil.RandomSlug("tf-test-tg-upd")
	updatedName := testutil.RandomName("tf-test-tenant-group-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccTenantGroupResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", updatedName),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(updatedName, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-group-import")
	slug := testutil.RandomSlug("tf-test-tenant-group-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_import(name, slug),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_tenant_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTenantGroupResourceConfig_import(name, slug),
				PlanOnly: true,
			},
		},
	})
}

// testAccTenantGroupResourceConfig_basic returns a basic test configuration.
func testAccTenantGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_TenantGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-lit")
	slug := testutil.RandomSlug("tf-test-tenant-group-lit")
	description := testutil.RandomName("description")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
				),
			},
		},
	})
}

// testAccTenantGroupResourceConfig_full returns a test configuration with all fields.
func testAccTenantGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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
	object_types = ["tenancy.tenantgroup"]
	type         = "text"
}

resource "netbox_tenant_group" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q

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
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccTenantGroupResourceConfig_fullUpdate(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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
	object_types = ["tenancy.tenantgroup"]
	type         = "text"
}

resource "netbox_tenant_group" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q

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
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccTenantGroupResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccTenantGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-del")
	slug := testutil.RandomSlug("tf-test-tenant-group-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyTenantGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tenant_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyTenantGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tenant_group: %v", err)
					}
					t.Logf("Successfully externally deleted tenant_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTenantGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tg-remove")
	slug := testutil.RandomSlug("tf-test-tg-remove")
	parentName := testutil.RandomName("tf-test-tg-parent")
	parentSlug := testutil.RandomSlug("tf-test-tg-parent")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)
	cleanup.RegisterTenantGroupCleanup(parentSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_withParent(name, slug, parentName, parentSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "parent"),
				),
			},
			{
				Config: testAccTenantGroupResourceConfig_detached(name, slug, parentName, parentSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_tenant_group.test", "parent"),
				),
			},
		},
	})
}

func testAccTenantGroupResourceConfig_withParent(name, slug, parentName, parentSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_tenant_group" "test" {
  name   = %q
  slug   = %q
  parent = netbox_tenant_group.parent.id
}
`, parentName, parentSlug, name, slug)
}

func testAccTenantGroupResourceConfig_detached(name, slug, parentName, parentSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_tenant_group" "test" {
  name   = %q
  slug   = %q
}
`, parentName, parentSlug, name, slug)
}

func TestAccTenantGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_tenant_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_tenant_group" "test" {
  # name missing
  slug = "test-tenant-group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_tenant_group" "test" {
  name = "Test Tenant Group"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccTenantGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_tenant_group" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccTenantGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_tenant_group" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}
