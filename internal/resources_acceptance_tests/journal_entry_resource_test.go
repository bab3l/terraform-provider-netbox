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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-journal")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "# Markdown header\n\nTest with markdown"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "info"),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.GenerateSlug(siteName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_updated(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Updated journal entry"),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_import(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.GenerateSlug(siteName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

func TestAccJournalEntryResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-id")
	siteSlug := testutil.RandomSlug("site-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
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

func testAccJournalEntryResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "# Markdown header\n\nTest with markdown"
  kind                 = "info"
}
`, name, testutil.RandomSlug("site"))
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryConsistencyLiteralNamesConfig(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
			},
			{
				Config:   testAccJournalEntryConsistencyLiteralNamesConfig(siteName, siteSlug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
			},
		},
	})
}

func testAccJournalEntryConsistencyLiteralNamesConfig(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry"
}
`, siteName, siteSlug)
}

func TestAccJournalEntryResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-site-je-extdel")
	siteSlug := testutil.GenerateSlug(siteName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
