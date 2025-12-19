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

func TestPowerPanelResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPanelResource()

	if r == nil {

		t.Fatal("Expected non-nil PowerPanel resource")

	}

}

func TestPowerPanelResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPanelResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"site", "name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"location", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestPowerPanelResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPanelResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_power_panel"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestPowerPanelResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPanelResource().(*resources.PowerPanelResource)

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

func TestAccPowerPanelResource_basic(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	panelName := testutil.RandomName("tf-test-panel")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccPowerPanelResourceConfig_basic(siteName, siteSlug, panelName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
				),
			},

			{

				ResourceName: "netbox_power_panel.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"site"},
			},
		},
	})

}

func TestAccPowerPanelResource_full(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site-full")

	siteSlug := testutil.RandomSlug("tf-test-site-full")

	panelName := testutil.RandomName("tf-test-panel-full")

	description := "Test power panel with all fields"

	updatedDescription := "Updated power panel description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccPowerPanelResourceConfig_full(siteName, siteSlug, panelName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),

					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", description),

					resource.TestCheckResourceAttr("netbox_power_panel.test", "comments", "Test comments"),
				),
			},

			{

				Config: testAccPowerPanelResourceConfig_full(siteName, siteSlug, panelName, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccPowerPanelResourceConfig_basic(siteName, siteSlug, panelName string) string {

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



`, siteName, siteSlug, panelName)

}

func testAccPowerPanelResourceConfig_full(siteName, siteSlug, panelName, description string) string {

	return fmt.Sprintf(`



resource "netbox_site" "test" {

  name   = %q

  slug   = %q

  status = "active"

}



resource "netbox_power_panel" "test" {

  site        = netbox_site.test.id

  name        = %q

  description = %q

  comments    = "Test comments"

}



`, siteName, siteSlug, panelName, description)

}

// TestAccConsistency_PowerPanel_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_PowerPanel_LiteralNames(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	panelName := testutil.RandomName("power-panel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelConsistencyLiteralNamesConfig(siteName, siteSlug, panelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "site", siteName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccPowerPanelConsistencyLiteralNamesConfig(siteName, siteSlug, panelName),
			},
		},
	})
}

func testAccPowerPanelConsistencyLiteralNamesConfig(siteName, siteSlug, panelName string) string {
	return fmt.Sprintf(`

resource "netbox_site" "test" {
  name   = "%[1]s"
  slug   = "%[2]s"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "%[3]s"
  # Use literal string name to mimic existing user state
  site = "%[1]s"

  depends_on = [netbox_site.test]
}

`, siteName, siteSlug, panelName)
}
