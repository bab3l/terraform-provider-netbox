package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackRoleResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	rackRoleName := testutil.RandomName("tf-test-rack-role")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
				),
			},
		},
	})
}

func TestAccRackRoleResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	rackRoleName := testutil.RandomName("tf-test-rack-role-full")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-f")
	description := testutil.RandomName("description")
	color := testutil.ColorOrange

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", color),
				),
			},
		},
	})
}

func TestAccRackRoleResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	rackRoleName := testutil.RandomName("tf-test-rack-role-upd")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-u")
	updatedDescription := testutil.Description2

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
				),
			},
			{
				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, updatedDescription, "00bcd4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", "00bcd4"),
				),
			},
		},
	})
}

func TestAccRackRoleResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names
	rackRoleName := testutil.RandomName("tf-test-rack-role-imp")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-i")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
			},
			{
				ResourceName:      "netbox_rack_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRackRoleResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	rackRoleName := testutil.RandomName("tf-test-rack-role-tags")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_tags(rackRoleName, rackRoleSlug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRackRoleResourceConfig_tags(rackRoleName, rackRoleSlug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRackRoleResourceConfig_tags(rackRoleName, rackRoleSlug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccRackRoleResourceConfig_tags(rackRoleName, rackRoleSlug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRackRoleResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	rackRoleName := testutil.RandomName("tf-test-rack-role-tag-order")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_tagsOrder(rackRoleName, rackRoleSlug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRackRoleResourceConfig_tagsOrder(rackRoleName, rackRoleSlug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack_role.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func TestAccConsistency_RackRole(t *testing.T) {
	t.Parallel()

	rackRoleName := testutil.RandomName("rack-role")
	rackRoleSlug := testutil.RandomSlug("rack-role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleConsistencyConfig(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRackRoleConsistencyConfig(rackRoleName, rackRoleSlug),
			},
		},
	})
}

func testAccRackRoleResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccRackRoleResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_rack_role" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRackRoleResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_rack_role" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}

func testAccRackRoleResourceConfig_full(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  color       = %[4]q
}
`, name, slug, description, color)
}

func TestAccConsistency_RackRole_LiteralNames(t *testing.T) {
	t.Parallel()

	rackRoleName := testutil.RandomName("tf-test-rack-role-lit")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-lit")
	description := testutil.RandomName("description")
	color := "4caf50"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", color),
				),
			},
			{
				Config:   testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
				),
			},
		},
	})
}

func testAccRackRoleConsistencyConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccRackRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	rackRoleName := testutil.RandomName("tf-test-rack-role-extdel")
	rackRoleSlug := testutil.RandomSlug("tf-test-rr-ed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					roles, _, err := client.DcimAPI.DcimRackRolesList(context.Background()).Slug([]string{rackRoleSlug}).Execute()
					if err != nil || roles == nil || len(roles.Results) == 0 {
						t.Fatalf("Failed to find rack role for external deletion: %v", err)
					}
					roleID := roles.Results[0].Id
					_, err = client.DcimAPI.DcimRackRolesDestroy(context.Background(), roleID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete rack role: %v", err)
					}
					t.Logf("Successfully externally deleted rack role with ID: %d", roleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccRackRoleResource_removeOptionalFields tests that the description field
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccRackRoleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	roleName := testutil.RandomName("tf-test-rackrole-rem")
	roleSlug := testutil.RandomSlug("tf-test-rr-rem")
	const testDescription = "Test Description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_withDescription(roleName, roleSlug, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", testDescription),
				),
			},
			{
				Config: testAccRackRoleConsistencyConfig(roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", roleName),
					resource.TestCheckNoResourceAttr("netbox_rack_role.test", "description"),
				),
			},
		},
	})
}

func testAccRackRoleResourceConfig_withDescription(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccRackRoleResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_rack_role",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_rack_role" "test" {
  slug = "test-role"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_rack_role" "test" {
  name = "Test Role"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
