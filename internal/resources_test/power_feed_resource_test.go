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

func TestPowerFeedResource(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerFeedResource()
	if r == nil {
		t.Fatal("Expected non-nil PowerFeed resource")
	}
}

func TestPowerFeedResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerFeedResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"power_panel", "name"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	computedAttrs := []string{"id", "status", "type", "supply", "phase", "voltage", "amperage", "max_utilization"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	optionalAttrs := []string{"rack", "mark_connected", "description", "tenant", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestPowerFeedResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerFeedResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_power_feed"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestPowerFeedResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerFeedResource().(*resources.PowerFeedResource)

	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	configureRequest.ProviderData = invalidProviderData
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccPowerFeedResource_basic(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	panelName := testutil.RandomName("tf-test-panel")
	feedName := testutil.RandomName("tf-test-feed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "status", "active"),
				),
			},
			{
				ResourceName:      "netbox_power_feed.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPowerFeedResource_full(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	panelName := testutil.RandomName("tf-test-panel-full")
	feedName := testutil.RandomName("tf-test-feed-full")
	description := "Test power feed with all fields"
	updatedDescription := "Updated power feed description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceConfig_full(siteName, siteSlug, panelName, feedName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "type", "primary"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "supply", "ac"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "phase", "single-phase"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "voltage", "240"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "amperage", "30"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "description", description),
				),
			},
			{
				Config: testAccPowerFeedResourceConfig_full(siteName, siteSlug, panelName, feedName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_feed.test", "description", updatedDescription),
				),
			},
		},
	})
}

func testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %q
}

resource "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = %q
}
`, siteName, siteSlug, panelName, feedName)
}

func testAccPowerFeedResourceConfig_full(siteName, siteSlug, panelName, feedName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %q
}

resource "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = %q
  status      = "active"
  type        = "primary"
  supply      = "ac"
  phase       = "single-phase"
  voltage     = 240
  amperage    = 30
  description = %q
}
`, siteName, siteSlug, panelName, feedName, description)
}
