package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestJournalEntryResource(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	if r == nil {
		t.Fatal("Expected non-nil Journal Entry resource")
	}
}

func TestJournalEntryResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}
	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes
	requiredAttrs := []string{"assigned_object_type", "assigned_object_id", "comments"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"id", "kind"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Optional attributes
	optionalAttrs := []string{"tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestJournalEntryResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}
	r.Metadata(context.Background(), metadataRequest, metadataResponse)
	expected := "netbox_journal_entry"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestJournalEntryResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource().(*resources.JournalEntryResource)

	// Test with nil provider data
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}
	r.Configure(context.Background(), configureRequest, configureResponse)
	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct client type
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwresource.ConfigureResponse{}
	r.Configure(context.Background(), configureRequest, configureResponse)
	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with incorrect provider data type
	configureRequest.ProviderData = invalidProviderData
	configureResponse = &fwresource.ConfigureResponse{}
	r.Configure(context.Background(), configureRequest, configureResponse)
	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccJournalEntryResource_basic(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-journal")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckJournalEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "assigned_object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "info"),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_full(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-journal-full")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckJournalEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_full(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "assigned_object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "# Important Note\n\nThis is a detailed journal entry with markdown."),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "warning"),
				),
			},
		},
	})
}

func TestAccJournalEntryResource_update(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-journal-update")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckJournalEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "info"),
				),
			},
			{
				Config: testAccJournalEntryResourceConfig_updated(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Updated journal entry content"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "success"),
				),
			},
		},
	})
}

func testAccJournalEntryResourceConfig_basic(siteName string) string {
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
`, siteName, testutil.GenerateSlug(siteName))
}

func testAccJournalEntryResourceConfig_full(siteName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "# Important Note\n\nThis is a detailed journal entry with markdown."
  kind                 = "warning"
}
`, siteName, testutil.GenerateSlug(siteName))
}

func testAccJournalEntryResourceConfig_updated(siteName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Updated journal entry content"
  kind                 = "success"
}
`, siteName, testutil.GenerateSlug(siteName))
}

func TestAccJournalEntryResource_import(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-journal")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(testutil.GenerateSlug(siteName))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckJournalEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryResourceConfig_basic(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_journal_entry.test", "id"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "assigned_object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Test journal entry"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "kind", "info"),
				),
			},
			{
				ResourceName:      "netbox_journal_entry.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
