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

func TestLocationResource(t *testing.T) {
	r := resources.NewLocationResource()
	if r == nil {
		t.Fatal("Location resource should not be nil")
	}
}

func TestLocationResourceSchema(t *testing.T) {
	ctx := context.Background()
	r := resources.NewLocationResource()

	schemaReq := fwresource.SchemaRequest{}
	schemaResp := &fwresource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Location resource schema should not have errors: %v", schemaResp.Diagnostics.Errors())
	}

	attrs := schemaResp.Schema.Attributes
	requiredAttrs := []string{"id", "name", "slug", "site"}
	for _, attr := range requiredAttrs {
		if _, ok := attrs[attr]; !ok {
			t.Errorf("Location resource schema should include %s attribute", attr)
		}
	}

	optionalAttrs := []string{"parent", "status", "tenant", "facility", "description", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, ok := attrs[attr]; !ok {
			t.Errorf("Location resource schema should include %s attribute", attr)
		}
	}
}

func TestLocationResourceMetadata(t *testing.T) {
	ctx := context.Background()
	r := resources.NewLocationResource()

	metadataReq := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResp := &fwresource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_location"
	if metadataResp.TypeName != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)
	}
}

func TestLocationResourceConfigure(t *testing.T) {
	ctx := context.Background()
	r := resources.NewLocationResource().(*resources.LocationResource)

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

func TestAccLocationResource_basic(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	siteName := testutil.RandomName("tf-test-loc-site")
	siteSlug := testutil.RandomSlug("tf-test-loc-site")
	name := testutil.RandomName("tf-test-location")
	slug := testutil.RandomSlug("tf-test-location")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttrPair("netbox_location.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})
}

func TestAccLocationResource_full(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-loc-site-full")
	siteSlug := testutil.RandomSlug("tf-test-loc-s-full")
	name := testutil.RandomName("tf-test-location-full")
	slug := testutil.RandomSlug("tf-test-loc-full")
	description := "Test location with all fields"
	facility := "Building A"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_full(siteName, siteSlug, name, slug, description, facility),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_location.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_location.test", "description", description),
					resource.TestCheckResourceAttr("netbox_location.test", "facility", facility),
				),
			},
		},
	})
}

func TestAccLocationResource_update(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-loc-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-loc-s-upd")
	name := testutil.RandomName("tf-test-location-update")
	slug := testutil.RandomSlug("tf-test-loc-upd")
	updatedName := testutil.RandomName("tf-test-location-upd2")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
				),
			},
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccLocationResource_withParent(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-loc-site-prnt")
	siteSlug := testutil.RandomSlug("tf-test-loc-s-prnt")
	parentName := testutil.RandomName("tf-test-loc-parent")
	parentSlug := testutil.RandomSlug("tf-test-loc-prnt")
	childName := testutil.RandomName("tf-test-loc-child")
	childSlug := testutil.RandomSlug("tf-test-loc-chld")

	// Register cleanup (child first, then parent, then site due to dependency)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(childSlug)
	cleanup.RegisterLocationCleanup(parentSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_withParent(siteName, siteSlug, parentName, parentSlug, childName, childSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.parent", "id"),
					resource.TestCheckResourceAttr("netbox_location.parent", "name", parentName),
					resource.TestCheckResourceAttrSet("netbox_location.child", "id"),
					resource.TestCheckResourceAttr("netbox_location.child", "name", childName),
					resource.TestCheckResourceAttrPair("netbox_location.child", "parent", "netbox_location.parent", "id"),
				),
			},
		},
	})
}

// testAccLocationResourceConfig_basic returns a basic test configuration
func testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug string) string {
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
`, siteName, siteSlug, name, slug)
}

// testAccLocationResourceConfig_full returns a test configuration with all fields
func testAccLocationResourceConfig_full(siteName, siteSlug, name, slug, description, facility string) string {
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
  name        = %q
  slug        = %q
  site        = netbox_site.test.id
  status      = "active"
  description = %q
  facility    = %q
}
`, siteName, siteSlug, name, slug, description, facility)
}

// testAccLocationResourceConfig_withParent returns a test configuration with parent location
func testAccLocationResourceConfig_withParent(siteName, siteSlug, parentName, parentSlug, childName, childSlug string) string {
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

resource "netbox_location" "parent" {
  name = %q
  slug = %q
  site = netbox_site.test.id
}

resource "netbox_location" "child" {
  name   = %q
  slug   = %q
  site   = netbox_site.test.id
  parent = netbox_location.parent.id
}
`, siteName, siteSlug, parentName, parentSlug, childName, childSlug)
}
