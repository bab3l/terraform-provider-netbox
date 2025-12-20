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

func TestPlatformResource(t *testing.T) {
	t.Parallel()
	r := resources.NewPlatformResource()
	if r == nil {
		t.Fatal("Expected non-nil platform resource")
	}
}

func TestPlatformResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewPlatformResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}
	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "slug", "manufacturer"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	optionalAttrs := []string{"description"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}

	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestPlatformResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewPlatformResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}
	r.Metadata(context.Background(), metadataRequest, metadataResponse)
	expected := "netbox_platform"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestPlatformResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewPlatformResource().(*resources.PlatformResource)
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

func TestAccPlatformResource_basic(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	platformName := testutil.RandomName("tf-test-platform")
	platformSlug := testutil.RandomSlug("tf-test-plat")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-platform")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-plat")

	// Register cleanup for both resources to ensure they are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "manufacturer", manufacturerSlug),
				),
			},
		},
	})
}

func TestAccPlatformResource_full(t *testing.T) {
	// Generate unique names
	platformName := testutil.RandomName("tf-test-platform-full")
	platformSlug := testutil.RandomSlug("tf-test-plat-full")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-plat-full")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pf")
	description := "Test platform with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_full(platformName, platformSlug, manufacturerName, manufacturerSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "manufacturer", manufacturerSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "description", description),
				),
			},
		},
	})
}

func TestAccPlatformResource_update(t *testing.T) {
	// Generate unique names
	platformName := testutil.RandomName("tf-test-platform-update")
	platformSlug := testutil.RandomSlug("tf-test-plat-upd")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-plat-upd")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pu")
	updatedName := testutil.RandomName("tf-test-platform-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
				),
			},
			{
				Config: testAccPlatformResourceConfig_basic(updatedName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", updatedName),
				),
			},
		},
	})
}

// testAccPlatformResourceConfig_basic returns a basic test configuration with manufacturer.
func testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_manufacturer" "test_manufacturer" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_manufacturer.slug
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug)
}

// testAccPlatformResourceConfig_full returns a test configuration with all fields.
func testAccPlatformResourceConfig_full(platformName, platformSlug, manufacturerName, manufacturerSlug, description string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_manufacturer" "test_manufacturer" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_manufacturer.slug
  description  = %q
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug, description)
}

func TestAccPlatformResource_import(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	platformName := testutil.RandomName("tf-test-platform-import")
	platformSlug := testutil.RandomSlug("tf-test-plat-imp")
	manufacturerName := testutil.RandomName("tf-test-mfr-imp")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_import(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
				),
			},
			{
				ResourceName:            "netbox_platform.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})
}

func testAccPlatformResourceConfig_import(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_platform" "test" {
  name         = "%[1]s"
  slug         = "%[2]s"
  manufacturer = netbox_manufacturer.test.slug
}
`, platformName, platformSlug, manufacturerName, manufacturerSlug)
}

// TestAccConsistency_Platform_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_Platform_LiteralNames(t *testing.T) {
	t.Parallel()
	platformName := testutil.RandomName("platform")
	platformSlug := testutil.RandomSlug("platform")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformConsistencyLiteralNamesConfig(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "manufacturer", manufacturerSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccPlatformConsistencyLiteralNamesConfig(platformName, platformSlug, manufacturerName, manufacturerSlug),
			},
		},
	})
}

func testAccPlatformConsistencyLiteralNamesConfig(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_platform" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  # Use literal string name to mimic existing user state
  manufacturer = "%[4]s"

  depends_on = [netbox_manufacturer.test]
}

`, platformName, platformSlug, manufacturerName, manufacturerSlug)
}
