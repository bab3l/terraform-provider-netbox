package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJournalEntryDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)
	siteName := testutil.RandomName("tf-test-site-journal-ds-id")
	cleanup.RegisterSiteCleanup(testutil.GenerateSlug(siteName))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckJournalEntryDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryDataSourceConfig_byID(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_journal_entry.test", "id"),
				),
			},
		},
	})
}

func TestAccJournalEntryDataSource_byID(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-site-journal-ds")

	cleanup.RegisterSiteCleanup(testutil.GenerateSlug(siteName))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckJournalEntryDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccJournalEntryDataSourceConfig_byID(siteName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_journal_entry.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "assigned_object_type", "dcim.site"),

					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "comments", "Test journal entry for data source"),

					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "kind", "info"),
				),
			},
		},
	})

}

func testAccJournalEntryDataSourceConfig_byID(siteName string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %q

  slug = %q

}

resource "netbox_journal_entry" "test" {

  assigned_object_type = "dcim.site"

  assigned_object_id   = netbox_site.test.id

  comments             = "Test journal entry for data source"

}

data "netbox_journal_entry" "test" {

  id = netbox_journal_entry.test.id

}

`, siteName, testutil.GenerateSlug(siteName))

}
