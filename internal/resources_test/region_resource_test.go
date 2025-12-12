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

func TestRegionResource(t *testing.T) {
	r := resources.NewRegionResource()
	if r == nil {
		t.Fatal("Region resource should not be nil")
	}
}

func TestRegionResourceSchema(t *testing.T) {
	ctx := context.Background()
	r := resources.NewRegionResource()

	schemaReq := fwresource.SchemaRequest{}
	schemaResp := &fwresource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Region resource schema should not have errors: %v", schemaResp.Diagnostics.Errors())
	}

	attrs := schemaResp.Schema.Attributes
	requiredAttrs := []string{"id", "name", "slug"}
	for _, attr := range requiredAttrs {
		if _, ok := attrs[attr]; !ok {
			t.Errorf("Region resource schema should include %s attribute", attr)
		}
	}

	optionalAttrs := []string{"parent", "description", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, ok := attrs[attr]; !ok {
			t.Errorf("Region resource schema should include %s attribute", attr)
		}
	}
}

func TestRegionResourceMetadata(t *testing.T) {
	ctx := context.Background()
	r := resources.NewRegionResource()

	metadataReq := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResp := &fwresource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_region"
	if metadataResp.TypeName != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)
	}
}

func TestRegionResourceConfigure(t *testing.T) {
	ctx := context.Background()
	r := resources.NewRegionResource().(*resources.RegionResource)

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

func TestAccRegionResource_basic(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-region")
	slug := testutil.RandomSlug("tf-test-region")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccRegionResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-region-full")
	slug := testutil.RandomSlug("tf-test-region-full")
	description := "Test region with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("netbox_region.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_region.test", "description", description),
				),
			},
		},
	})
}

func TestAccRegionResource_update(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-region-update")
	slug := testutil.RandomSlug("tf-test-region-upd")
	updatedName := testutil.RandomName("tf-test-region-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", name),
				),
			},
			{
				Config: testAccRegionResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccRegionResource_withParent(t *testing.T) {
	// Generate unique names
	parentName := testutil.RandomName("tf-test-region-parent")
	parentSlug := testutil.RandomSlug("tf-test-region-prnt")
	childName := testutil.RandomName("tf-test-region-child")
	childSlug := testutil.RandomSlug("tf-test-region-chld")

	// Register cleanup (child first, then parent due to dependency)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(childSlug)
	cleanup.RegisterRegionCleanup(parentSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.parent", "id"),
					resource.TestCheckResourceAttr("netbox_region.parent", "name", parentName),
					resource.TestCheckResourceAttrSet("netbox_region.child", "id"),
					resource.TestCheckResourceAttr("netbox_region.child", "name", childName),
					resource.TestCheckResourceAttrPair("netbox_region.child", "parent", "netbox_region.parent", "id"),
				),
			},
		},
	})
}

// testAccRegionResourceConfig_basic returns a basic test configuration
func testAccRegionResourceConfig_basic(name, slug string) string {
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

resource "netbox_region" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

// testAccRegionResourceConfig_full returns a test configuration with all fields
func testAccRegionResourceConfig_full(name, slug, description string) string {
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

resource "netbox_region" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

// testAccRegionResourceConfig_withParent returns a test configuration with parent region
func testAccRegionResourceConfig_withParent(parentName, parentSlug, childName, childSlug string) string {
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

resource "netbox_region" "parent" {
  name = %q
  slug = %q
}

resource "netbox_region" "child" {
  name   = %q
  slug   = %q
  parent = netbox_region.parent.id
}
`, parentName, parentSlug, childName, childSlug)
}
