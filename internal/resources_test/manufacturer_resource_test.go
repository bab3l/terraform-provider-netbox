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

func TestManufacturerResource(t *testing.T) {

	t.Parallel()

	r := resources.NewManufacturerResource()

	if r == nil {

		t.Fatal("Expected non-nil manufacturer resource")

	}

}

func TestManufacturerResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewManufacturerResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"name", "slug"}

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

}

func TestManufacturerResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewManufacturerResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_manufacturer"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestManufacturerResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewManufacturerResource().(*resources.ManufacturerResource)

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

func TestAccManufacturerResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-manufacturer")

	slug := testutil.RandomSlug("tf-test-mfr")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccManufacturerResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccManufacturerResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-manufacturer-full")

	slug := testutil.RandomSlug("tf-test-mfr-full")

	description := "Test manufacturer with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccManufacturerResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "description", description),
				),
			},
		},
	})

}

func TestAccManufacturerResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-manufacturer-update")

	slug := testutil.RandomSlug("tf-test-mfr-upd")

	updatedName := testutil.RandomName("tf-test-manufacturer-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccManufacturerResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
				),
			},

			{

				Config: testAccManufacturerResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),

					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", updatedName),
				),
			},
		},
	})

}

// testAccManufacturerResourceConfig_basic returns a basic test configuration.

func testAccManufacturerResourceConfig_basic(name, slug string) string {

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































































resource "netbox_manufacturer" "test" {































  name = %q































  slug = %q































}































`, name, slug)

}

// testAccManufacturerResourceConfig_full returns a test configuration with all fields.

func testAccManufacturerResourceConfig_full(name, slug, description string) string {

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































































resource "netbox_manufacturer" "test" {































  name        = %q































  slug        = %q































  description = %q































}































`, name, slug, description)

}

func TestAccManufacturerResource_import(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-manufacturer-import")
	slug := testutil.RandomSlug("tf-test-mfr-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_manufacturer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccManufacturerResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}
