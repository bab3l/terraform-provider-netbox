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

func TestRackResource(t *testing.T) {
	r := resources.NewRackResource()
	if r == nil {
		t.Fatal("Rack resource should not be nil")
	}
}

func TestRackResourceSchema(t *testing.T) {
	ctx := context.Background()
	r := resources.NewRackResource()

	schemaReq := fwresource.SchemaRequest{}
	schemaResp := &fwresource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Rack resource schema should not have errors: %v", schemaResp.Diagnostics.Errors())
	}

	attrs := schemaResp.Schema.Attributes
	requiredAttrs := []string{"id", "name", "site"}
	for _, attr := range requiredAttrs {
		if _, ok := attrs[attr]; !ok {
			t.Errorf("Rack resource schema should include %s attribute", attr)
		}
	}

	optionalAttrs := []string{"location", "tenant", "status", "role", "rack_type", "serial", "asset_tag",
		"form_factor", "width", "u_height", "starting_unit", "weight", "max_weight", "weight_unit",
		"desc_units", "outer_width", "outer_depth", "outer_unit", "mounting_depth", "airflow",
		"description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, ok := attrs[attr]; !ok {
			t.Errorf("Rack resource schema should include %s attribute", attr)
		}
	}
}

func TestRackResourceMetadata(t *testing.T) {
	ctx := context.Background()
	r := resources.NewRackResource()

	metadataReq := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResp := &fwresource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_rack"
	if metadataResp.TypeName != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)
	}
}

func TestRackResourceConfigure(t *testing.T) {
	ctx := context.Background()
	r := resources.NewRackResource().(*resources.RackResource)

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

	configureReq.ProviderData = invalidProviderData
	configureResp = &fwresource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if !configureResp.Diagnostics.HasError() {
		t.Error("Configure should error with incorrect provider data type")
	}
}

func TestAccRackResource_basic(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	siteName := testutil.RandomName("tf-test-rack-site")
	siteSlug := testutil.RandomSlug("tf-test-rack-site")
	rackName := testutil.RandomName("tf-test-rack")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})
}

func TestAccRackResource_full(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-full")
	siteSlug := testutil.RandomSlug("tf-test-rack-s-full")
	rackName := testutil.RandomName("tf-test-rack-full")
	description := "Test rack with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_full(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_rack.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack.test", "u_height", "42"),
					resource.TestCheckResourceAttr("netbox_rack.test", "width", "19"),
				),
			},
		},
	})
}

func TestAccRackResource_update(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-rack-s-upd")
	rackName := testutil.RandomName("tf-test-rack-update")
	updatedName := testutil.RandomName("tf-test-rack-upd2")

	// Register cleanup (use original name for initial cleanup, register updated name too)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterRackCleanup(updatedName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccRackResource_withLocation(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-loc")
	siteSlug := testutil.RandomSlug("tf-test-rack-s-loc")
	locationName := testutil.RandomName("tf-test-rack-location")
	locationSlug := testutil.RandomSlug("tf-test-rack-loc")
	rackName := testutil.RandomName("tf-test-rack-with-loc")

	// Register cleanup (rack first, then location, then site due to dependency)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterLocationCleanup(locationSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "location", "netbox_location.test", "id"),
				),
			},
		},
	})
}

// testAccRackResourceConfig_basic returns a basic test configuration.
func testAccRackResourceConfig_basic(siteName, siteSlug, rackName string) string {
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

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.id
}
`, siteName, siteSlug, rackName)
}

// testAccRackResourceConfig_full returns a test configuration with all fields.
func testAccRackResourceConfig_full(siteName, siteSlug, rackName, description string) string {
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

resource "netbox_rack" "test" {
  name        = %q
  site        = netbox_site.test.id
  status      = "active"
  u_height    = 42
  width       = 19
  description = %q
}
`, siteName, siteSlug, rackName, description)
}

// testAccRackResourceConfig_withLocation returns a test configuration with location.
func testAccRackResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, rackName string) string {
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

resource "netbox_location" "test" {
  name = %q
  slug = %q
  site = netbox_site.test.id
}

resource "netbox_rack" "test" {
  name     = %q
  site     = netbox_site.test.id
  location = netbox_location.test.id
}
`, siteName, siteSlug, locationName, locationSlug, rackName)
}
