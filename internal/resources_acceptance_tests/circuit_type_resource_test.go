package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTypeResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type")
	slug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitTypeResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-full")
	slug := testutil.RandomSlug("tf-test-circuit-type-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_full(name, slug, description, testutil.Color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "color", testutil.Color),
				),
			},
		},
	})
}

func TestAccCircuitTypeResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-update")
	slug := testutil.RandomSlug("tf-test-circuit-type-update")
	updatedName := testutil.RandomName("tf-test-circuit-type-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
				),
			},
			{
				Config: testAccCircuitTypeResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccCircuitTypeResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type")
	slug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_circuit_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccCircuitTypeResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_CircuitType_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-lit")
	slug := testutil.RandomSlug("tf-test-circuit-type-lit")
	description := testutil.RandomName("description")
	color := "2196f3" //nolint:goconst // Blue color value used in multiple test files

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeConsistencyLiteralNamesConfig(name, slug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "color", color),
				),
			},
			{
				Config:   testAccCircuitTypeConsistencyLiteralNamesConfig(name, slug, description, color),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
				),
			},
		},
	})
}

func testAccCircuitTypeResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccCircuitTypeResourceConfig_full(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}
`, name, slug, description, color)
}

func testAccCircuitTypeConsistencyLiteralNamesConfig(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}
`, name, slug, description, color)
}

func TestAccCircuitTypeResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-ext-del")
	slug := testutil.RandomSlug("circuit-type-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List circuit types filtered by slug
					items, _, err := client.CircuitsAPI.CircuitsCircuitTypesList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit type for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitTypesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit type: %v", err)
					}
					t.Logf("Successfully externally deleted circuit type with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCircuitTypeResource_removeDescription(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-rem-desc")
	slug := testutil.RandomSlug("tf-test-circuit-type-rem-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_withDescription(name, slug, "Description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", "Description"),
				),
			},
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_circuit_type.test", "description"),
				),
			},
		},
	})
}

func testAccCircuitTypeResourceConfig_withDescription(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_type" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

// TestAccCircuitTypeResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccCircuitTypeResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-type-rem")
	slug := testutil.RandomSlug("tf-test-circuit-type-rem")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_full(name, slug, description, testutil.Color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "color", testutil.Color),
				),
			},
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
					resource.TestCheckNoResourceAttr("netbox_circuit_type.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_circuit_type.test", "color"),
				),
			},
		},
	})
}

func TestAccCircuitTypeResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_circuit_type",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_type" "test" {
  # name missing
  slug = "test-type"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_type" "test" {
  name = "Test Type"
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

// TestAccCircuitTypeResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccCircuitTypeResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-ct-tag")
	slug := testutil.RandomSlug("tf-ct-tag")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_circuit_type",
		ConfigWithoutTags: func() string {
			return testAccCircuitTypeResourceConfig_tagLifecycle(name, slug, "", "", "", "", "", "", "")
		},
		ConfigWithTags: func() string {
			return testAccCircuitTypeResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1,tag2", "", "")
		},
		ConfigWithDifferentTags: func() string {
			return testAccCircuitTypeResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag3", tag3Name, tag3Slug)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
	})
}

func testAccCircuitTypeResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagSet, tag3Name, tag3Slug string) string {
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
		case "tag1,tag2":
			tagsList = tagsDoubleSlug
		case "tag3":
			tagsList = tagsSingleSlug
		default:
			tagsList = tagsEmpty
		}
	} else {
		tagsList = tagsEmpty
	}

	return fmt.Sprintf(`
%s
resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
  %s
}
`, tagResources, name, slug, tagsList)
}

// TestAccCircuitTypeResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccCircuitTypeResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-ct-order")
	slug := testutil.RandomSlug("tf-ct-order")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_circuit_type",
		ConfigWithTagsOrderA: func() string {
			return testAccCircuitTypeResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccCircuitTypeResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
	})
}

func testAccCircuitTypeResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
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

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
  %s
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug, name, slug, tagsOrder)
}
