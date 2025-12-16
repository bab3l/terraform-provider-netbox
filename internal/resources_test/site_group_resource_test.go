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

func TestSiteGroupResource(t *testing.T) {

	r := resources.NewSiteGroupResource()

	if r == nil {

		t.Fatal("Site group resource should not be nil")

	}

}

func TestSiteGroupResourceSchema(t *testing.T) {

	ctx := context.Background()

	r := resources.NewSiteGroupResource()

	schemaReq := fwresource.SchemaRequest{}

	schemaResp := &fwresource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {

		t.Errorf("Schema should not have errors: %v", schemaResp.Diagnostics.Errors())

	}

	requiredAttributes := []string{"id", "name", "slug"}

	for _, attr := range requiredAttributes {

		if _, exists := schemaResp.Schema.Attributes[attr]; !exists {

			t.Errorf("Site group resource schema should include %s attribute", attr)

		}

	}

	optionalAttributes := []string{"parent", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttributes {

		if _, exists := schemaResp.Schema.Attributes[attr]; !exists {

			t.Errorf("Site group resource schema should include %s attribute", attr)

		}

	}

}

func TestSiteGroupResourceMetadata(t *testing.T) {

	ctx := context.Background()

	r := resources.NewSiteGroupResource()

	metadataReq := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResp := &fwresource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_site_group"

	if metadataResp.TypeName != expectedTypeName {

		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)

	}

}

func TestSiteGroupResourceConfigure(t *testing.T) {

	ctx := context.Background()

	r := resources.NewSiteGroupResource().(*resources.SiteGroupResource)

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

func TestAccSiteGroupResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-site-group")

	slug := testutil.RandomSlug("tf-test-sg")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccSiteGroupResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-site-group-full")

	slug := testutil.RandomSlug("tf-test-sg-full")

	description := "Test site group with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_site_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_site_group.test", "description", description),
				),
			},
		},
	})

}

func TestAccSiteGroupResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-site-group-update")

	slug := testutil.RandomSlug("tf-test-sg-upd")

	updatedName := testutil.RandomName("tf-test-site-group-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_site_group.test", "name", name),
				),
			},

			{

				Config: testAccSiteGroupResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_site_group.test", "name", updatedName),
				),
			},
		},
	})

}

// testAccSiteGroupResourceConfig_basic returns a basic test configuration.

func testAccSiteGroupResourceConfig_basic(name, slug string) string {

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































resource "netbox_site_group" "test" {















  name = %q















  slug = %q















}















`, name, slug)

}

// testAccSiteGroupResourceConfig_full returns a test configuration with all fields.

func testAccSiteGroupResourceConfig_full(name, slug, description string) string {

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































resource "netbox_site_group" "test" {















  name        = %q















  slug        = %q















  description = %q















}















`, name, slug, description)

}
