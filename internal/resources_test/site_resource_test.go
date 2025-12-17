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

func TestSiteResource(t *testing.T) {

	r := resources.NewSiteResource()

	if r == nil {

		t.Fatal("Site resource should not be nil")

	}

}

func TestSiteResourceSchema(t *testing.T) {

	ctx := context.Background()

	r := resources.NewSiteResource()

	schemaReq := fwresource.SchemaRequest{}

	schemaResp := &fwresource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {

		t.Fatalf("Site resource schema should not have errors: %v", schemaResp.Diagnostics.Errors())

	}

	attrs := schemaResp.Schema.Attributes

	requiredAttrs := []string{"id", "name", "slug"}

	for _, attr := range requiredAttrs {

		if _, ok := attrs[attr]; !ok {

			t.Errorf("Site resource schema should include %s attribute", attr)

		}

	}

	optionalAttrs := []string{"status", "description", "comments", "facility", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, ok := attrs[attr]; !ok {

			t.Errorf("Site resource schema should include %s attribute", attr)

		}

	}

}

func TestSiteResourceMetadata(t *testing.T) {

	ctx := context.Background()

	r := resources.NewSiteResource()

	metadataReq := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResp := &fwresource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_site"

	if metadataResp.TypeName != expectedTypeName {

		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)

	}

}

func TestSiteResourceConfigure(t *testing.T) {

	ctx := context.Background()

	r := resources.NewSiteResource().(*resources.SiteResource)

	configureReq := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResp := &fwresource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if configureResp.Diagnostics.HasError() {

		t.Error("Configure should not error with nil provider data")

	}

	client := &netbox.APIClient{}

	configureReq.ProviderData = client

	configureResp = &fwresource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if configureResp.Diagnostics.HasError() {

		t.Errorf("Configure should not error with correct provider data: %v", configureResp.Diagnostics.Errors())

	}

	configureReq.ProviderData = "invalid"

	configureResp = &fwresource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if !configureResp.Diagnostics.HasError() {

		t.Error("Configure should error with incorrect provider data type")

	}

}

func TestAccSiteResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-site")

	slug := testutil.RandomSlug("tf-test-site")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),

					resource.TestCheckResourceAttr("netbox_site.test", "name", name),

					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
				),
			},
		},
	})

}

func TestAccSiteResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-site-full")

	slug := testutil.RandomSlug("tf-test-site-full")

	description := "Test site with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),

					resource.TestCheckResourceAttr("netbox_site.test", "name", name),

					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_site.test", "description", description),
				),
			},
		},
	})

}

func TestAccSiteResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-site-update")

	slug := testutil.RandomSlug("tf-test-site-upd")

	updatedName := testutil.RandomName("tf-test-site-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),

					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
				),
			},

			{

				Config: testAccSiteResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),

					resource.TestCheckResourceAttr("netbox_site.test", "name", updatedName),
				),
			},
		},
	})

}

// testAccSiteResourceConfig_basic returns a basic test configuration.

func testAccSiteResourceConfig_basic(name, slug string) string {

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































































resource "netbox_site" "test" {































  name   = %q































  slug   = %q































  status = "active"































}































`, name, slug)

}

// testAccSiteResourceConfig_full returns a test configuration with all fields.

func testAccSiteResourceConfig_full(name, slug, description string) string {

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































































resource "netbox_site" "test" {































  name        = %q































  slug        = %q































  status      = "active"































  description = %q































}































`, name, slug, description)

}

func TestAccSiteResource_import(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-site-import")
	slug := testutil.RandomSlug("tf-test-site-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSiteResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}
