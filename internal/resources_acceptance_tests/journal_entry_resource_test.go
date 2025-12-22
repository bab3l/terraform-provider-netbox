package resources_acceptance_tests

import (
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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccJournalEntryResourceConfig_basic(siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
				),
			},

			{

				ResourceName: "netbox_journal_entry.test",

				ImportState: true,

				ImportStateVerify: true,
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
