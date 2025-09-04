package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSiteResource(t *testing.T) {
	// Test that the site resource can be instantiated
	r := resources.NewSiteGroupResource()
	if r == nil {
		t.Fatal("Site resource should not be nil")
	}
}

func TestSiteResourceSchema(t *testing.T) {
	ctx := context.Background()
	r := resources.NewSiteResource()

	// Test that the resource schema can be retrieved
	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Site resource schema should not have errors: %v", schemaResp.Diagnostics.Errors())
	}

	// Verify essential attributes exist
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

	// Test that the resource metadata can be retrieved
	metadataReq := resource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResp := &resource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_site"
	if metadataResp.TypeName != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)
	}
}

func TestSiteResourceConfigure(t *testing.T) {
	ctx := context.Background()
	r := resources.NewSiteResource().(*resources.SiteResource) // Cast to access Configure method

	// Test with nil provider data (should not error)
	configureReq := resource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResp := &resource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if configureResp.Diagnostics.HasError() {
		t.Error("Configure should not error with nil provider data")
	}

	// Test with correct provider data type
	client := &netbox.APIClient{}
	configureReq.ProviderData = client
	configureResp = &resource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if configureResp.Diagnostics.HasError() {
		t.Errorf("Configure should not error with correct provider data: %v", configureResp.Diagnostics.Errors())
	}

	// Test with incorrect provider data type
	configureReq.ProviderData = "invalid"
	configureResp = &resource.ConfigureResponse{}

	r.Configure(ctx, configureReq, configureResp)

	if !configureResp.Diagnostics.HasError() {
		t.Error("Configure should error with incorrect provider data type")
	}
}
