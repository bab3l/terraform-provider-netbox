package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJournalEntryResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
				),
			},
			{
				Config:   testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccJournalEntryResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-journal-site")
	siteSlug := testutil.RandomSlug("tf-test-journal-site")
	comments := "# Markdown header\n\nTest with markdown"
	updatedComments := "# Updated header\n\nUpdated markdown"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_full(siteName, siteSlug, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "info"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccJournalEntryResourceConfig_full(siteName, siteSlug, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
			{
				Config: testAccJournalEntryResourceConfig_fullUpdate(siteName, siteSlug, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", updatedComments),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "warning"),
				),
			},
			{
				Config:   testAccJournalEntryResourceConfig_fullUpdate(siteName, siteSlug, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccJournalEntryResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-journal-tags-site")
	siteSlug := testutil.RandomSlug("tf-test-journal-tags-site")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_tags(siteName, siteSlug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_tags(siteName, siteSlug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_tags(siteName, siteSlug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_tags(siteName, siteSlug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-journal-tag-order-site")
	siteSlug := testutil.RandomSlug("tf-test-journal-tag-order-site")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_tagsOrder(siteName, siteSlug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_tagsOrder(siteName, siteSlug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_journal_entry.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
				),
			},
			{
				Config:   testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				PlanOnly: true,
			},
			{
				Config: testAccJournalEntryResourceConfig_updated(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Updated journal entry"),
				),
			},
			{
				Config:   testAccJournalEntryResourceConfig_updated(siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccJournalEntryResource_import(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_journal_entry.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

func testAccJournalEntryResourceConfig_basic(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry"
}
`, siteName, siteSlug)
}

func testAccJournalEntryResourceConfig_full(siteName, siteSlug, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %q
	slug = %q
}

resource "netbox_tag" "tag2" {
	name = %q
	slug = %q
}

resource "netbox_site" "test" {
	name = %q
	slug = %q
}

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = %q
	kind                 = "info"

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, tagName1, tagSlug1, tagName2, tagSlug2, siteName, siteSlug, comments)
}

func testAccJournalEntryResourceConfig_fullUpdate(siteName, siteSlug, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %q
	slug = %q
}

resource "netbox_tag" "tag2" {
	name = %q
	slug = %q
}

resource "netbox_site" "test" {
	name = %q
	slug = %q
}

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = %q
	kind                 = "warning"

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, tagName1, tagSlug1, tagName2, tagSlug2, siteName, siteSlug, comments)
}

func testAccJournalEntryResourceConfig_updated(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Updated journal entry"
}
`, siteName, siteSlug)
}

func TestAccConsistency_JournalEntry_LiteralNames(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-lit")
	siteSlug := testutil.RandomSlug("tf-test-site-lit")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryConsistencyLiteralNamesConfig(siteName, siteSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccJournalEntryConsistencyLiteralNamesConfig(siteName, siteSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
			},
		},
	})
}

func testAccJournalEntryConsistencyLiteralNamesConfig(siteName, siteSlug, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_tag" "tag2" {
	name = %[5]q
	slug = %[6]q
}

resource "netbox_site" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = "Test journal entry"

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, siteName, siteSlug, tagName1, tagSlug1, tagName2, tagSlug2)
}

func TestAccJournalEntryResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-site-je-extdel")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find site by name to get its ID
					siteItems, _, err := client.DcimAPI.DcimSitesList(context.Background()).Name([]string{siteName}).Execute()
					if err != nil {
						t.Fatalf("Failed to list sites: %v", err)
					}
					if siteItems == nil || len(siteItems.Results) == 0 {
						t.Fatalf("Site not found with name: %s", siteName)
					}
					siteID := siteItems.Results[0].Id

					// Find journal entry by assigned_object_id (site ID)
					journalItems, _, err := client.ExtrasAPI.ExtrasJournalEntriesList(context.Background()).AssignedObjectId([]int32{siteID}).Execute()
					if err != nil {
						t.Fatalf("Failed to list journal entries: %v", err)
					}
					if journalItems == nil || len(journalItems.Results) == 0 {
						t.Fatalf("Journal entry not found for site ID: %d", siteID)
					}

					// Delete the journal entry
					itemID := journalItems.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasJournalEntriesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete journal entry: %v", err)
					}

					t.Logf("Successfully externally deleted journal entry with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccJournalEntryResource_Kind tests comprehensive scenarios for journal entry kind field.
// This validates that Optional+Computed string fields with proper defaults work correctly.
func TestAccJournalEntryResource_Kind(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-journal")
	siteSlug := testutil.RandomSlug("tf-test-site-journal")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_journal_entry",
		OptionalField:  "kind",
		DefaultValue:   "info",
		FieldTestValue: "warning",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		BaseConfig: func() string {
			return testAccJournalEntryResourceConfig_kindBase(siteName, siteSlug)
		},
		WithFieldConfig: func(value string) string {
			return testAccJournalEntryResourceConfig_kindWithField(siteName, siteSlug, value)
		},
	})
}

func testAccJournalEntryResourceConfig_kindBase(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry for kind field validation"
  # kind field intentionally omitted - should get default "info"
}
`, siteName, siteSlug)
}

func testAccJournalEntryResourceConfig_kindWithField(siteName, siteSlug, kindValue string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry for kind field validation"
  kind                 = %[3]q
}
`, siteName, siteSlug, kindValue)
}

// TestAccJournalEntryResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccJournalEntryResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-rem")
	siteSlug := testutil.RandomSlug("tf-test-site-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckJournalEntryDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_withKind(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "warning"),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_withoutKind(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "info"), // default value
				),
			},
		},
	})
}

func testAccJournalEntryResourceConfig_withKind(siteName, siteSlug string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry"
  kind                 = "warning"
}
`, siteName, siteSlug)
}

func testAccJournalEntryResourceConfig_withoutKind(siteName, siteSlug string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry"
}
`, siteName, siteSlug)
}

func TestAccJournalEntryResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_journal_entry",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_assigned_object_type": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_journal_entry" "test" {
  # assigned_object_type missing
  assigned_object_id = 1
  comments = "Test entry"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_assigned_object_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.device"
  # assigned_object_id missing
  comments = "Test entry"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccJournalEntryResourceConfig_tags(siteName, siteSlug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
resource "netbox_site" "test" {
	name = %[1]q
	slug = %[2]q
}

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

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = "Test journal entry with tags"
	%[6]s
}
`, siteName, siteSlug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccJournalEntryResourceConfig_tagsOrder(siteName, siteSlug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_tag" "tag1" {
	name = "Tag1-%[3]s"
	slug = %[3]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[4]s"
	slug = %[4]q
}

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = "Test journal entry with tags"
	%[5]s
}
`, siteName, siteSlug, tag1Slug, tag2Slug, tagsConfig)
}
