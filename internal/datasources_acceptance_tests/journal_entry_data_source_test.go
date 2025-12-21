package datasources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestJournalEntryDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewJournalEntryDataSource()

	if d == nil {

		t.Fatal("Expected non-nil Journal Entry data source")

	}

}

func TestJournalEntryDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewJournalEntryDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Required attribute for lookup

	if _, exists := schemaResponse.Schema.Attributes["id"]; !exists {

		t.Error("Expected lookup attribute 'id' to exist in schema")

	}

	// Computed attributes

	computedAttrs := []string{"assigned_object_type", "assigned_object_id", "kind", "comments"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

}

func TestJournalEntryDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewJournalEntryDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_journal_entry"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestJournalEntryDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewJournalEntryDataSource().(*datasources.JournalEntryDataSource)

	// Test with nil provider data

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct client type

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with incorrect provider data type

	configureRequest.ProviderData = testutil.InvalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

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
