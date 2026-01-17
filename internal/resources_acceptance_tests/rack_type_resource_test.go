package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackTypeResource_basic(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	model := testutil.RandomName("tf-test-rack-type")
	slug := testutil.RandomSlug("tf-test-rack-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "manufacturer"),
				),
			},
			{
				ResourceName:            "netbox_rack_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})
}

func TestAccRackTypeResource_full(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	model := testutil.RandomName("tf-test-rack-type-full")
	slug := testutil.RandomSlug("tf-test-rack-type-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated rack type description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, description, 42, 19),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "42"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "width", "19"),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, updatedDescription, 48, 19),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "48"),
				),
			},
		},
	})
}

func TestAccRackTypeResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-tags")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-tags")
	model := testutil.RandomName("tf-test-rack-type-tags")
	slug := testutil.RandomSlug("tf-test-rack-type-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_tags(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_tags(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_tags(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_tags(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRackTypeResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-tag-order")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-tag-order")
	model := testutil.RandomName("tf-test-rack-type-tag-order")
	slug := testutil.RandomSlug("tf-test-rack-type-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_tagsOrder(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_tagsOrder(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_type.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccConsistency_RackType_LiteralNames(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	model := testutil.RandomName("rack-type")
	slug := testutil.RandomSlug("rack-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "manufacturer", mfgName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug),
			},
		},
	})
}

func testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  form_factor  = "4-post-cabinet"
}
`, mfgName, mfgSlug, model, slug)
}

func testAccRackTypeResourceConfig_tags(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
resource "netbox_manufacturer" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_tag" "tag1" {
	name = "Tag1-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[6]s"
	slug = %[6]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[7]s"
	slug = %[7]q
}

resource "netbox_rack_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[3]q
	slug         = %[4]q
	form_factor  = "4-post-cabinet"
	%[8]s
}
`, mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRackTypeResourceConfig_tagsOrder(mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_tag" "tag1" {
	name = "Tag1-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[6]s"
	slug = %[6]q
}

resource "netbox_rack_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[3]q
	slug         = %[4]q
	form_factor  = "4-post-cabinet"
	%[7]s
}
`, mfgName, mfgSlug, model, slug, tag1Slug, tag2Slug, tagsConfig)
}

func testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, description string, uHeight, width int) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  description  = %q
  u_height     = %d
  width        = %d
  form_factor  = "4-post-cabinet"
}
`, mfgName, mfgSlug, model, slug, description, uHeight, width)
}

func testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_type" "test" {
  # Use literal string name to mimic existing user state
  manufacturer = %q
  model        = %q
  slug         = %q
  u_height     = 42
  width        = 19
  form_factor  = "4-post-cabinet"
  depends_on = [netbox_manufacturer.test]
}
`, mfgName, mfgSlug, mfgName, model, slug)
}

func TestAccRackTypeResource_update(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-upd")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-upd")
	model := testutil.RandomName("tf-test-rack-type-upd")
	slug := testutil.RandomSlug("tf-test-rt-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, testutil.Description2, 48, 19),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", testutil.Description2),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "48"),
				),
			},
		},
	})
}

func TestAccRackTypeResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	mfgName := testutil.RandomName("tf-test-mfg-extdel")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-ed")
	model := testutil.RandomName("tf-test-rack-type-extdel")
	slug := testutil.RandomSlug("tf-test-rt-ed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					types, _, err := client.DcimAPI.DcimRackTypesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || types == nil || len(types.Results) == 0 {
						t.Fatalf("Failed to find rack type for external deletion: %v", err)
					}
					typeID := types.Results[0].Id
					_, err = client.DcimAPI.DcimRackTypesDestroy(context.Background(), typeID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete rack type: %v", err)
					}
					t.Logf("Successfully externally deleted rack type with ID: %d", typeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRackTypeResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-optional")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-opt")
	model := testutil.RandomName("tf-test-rack-type-opt")
	slug := testutil.RandomSlug("tf-test-rack-type-opt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_rack_type",
		BaseConfig: func() string {
			return testAccRackTypeResourceConfig_removeOptionalFields_base(mfgName, mfgSlug, model, slug)
		},
		ConfigWithFields: func() string {
			return testAccRackTypeResourceConfig_removeOptionalFields_withFields(mfgName, mfgSlug, model, slug)
		},
		OptionalFields: map[string]string{
			"mounting_depth": "30",
			"outer_depth":    "1000",
			"outer_unit":     "mm",
			"outer_width":    "600",
		},
		RequiredFields: map[string]string{
			"model": model,
			"slug":  slug,
		},
		CheckDestroy: nil, // No CheckRackTypeDestroy function available
	})
}

func testAccRackTypeResourceConfig_removeOptionalFields_base(mfgName, mfgSlug, model, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
  form_factor  = "4-post-frame"
  weight_unit  = "kg"
}
`, mfgName, mfgSlug, model, slug)
}

func testAccRackTypeResourceConfig_removeOptionalFields_withFields(mfgName, mfgSlug, model, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_rack_type" "test" {
  manufacturer   = netbox_manufacturer.test.id
  model          = %[3]q
  slug           = %[4]q
  form_factor    = "4-post-frame"
  weight_unit    = "kg"
  mounting_depth = 30
  outer_depth    = 1000
  outer_unit     = "mm"
  outer_width    = 600
}
`, mfgName, mfgSlug, model, slug)
}

func TestAccRackTypeResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_rack_type",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_manufacturer": {
				Config: func() string {
					return `
resource "netbox_rack_type" "test" {
  model = "Test Model"
  slug  = "test-model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_model": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  slug         = "test-model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_manufacturer_reference": {
				Config: func() string {
					return `
resource "netbox_rack_type" "test" {
  manufacturer = "99999"
  model        = "Test Model"
  slug         = "test-model"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
