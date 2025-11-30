package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
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
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
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
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", "Test Site"),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", "test-site"),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
				),
			},
		},
	})
}
