package resources

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSiteGroupResource(t *testing.T) {
	// Test that the site group resource can be instantiated
	r := NewSiteGroupResource()
	if r == nil {
		t.Fatal("Site group resource should not be nil")
	}
}

func TestSiteGroupResourceSchema(t *testing.T) {
	ctx := context.Background()
	r := NewSiteGroupResource()

	// Test that the resource schema can be retrieved
	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {
		t.Errorf("Schema should not have errors: %v", schemaResp.Diagnostics.Errors())
	}

	// Check that essential attributes are present
	requiredAttributes := []string{"id", "name", "slug"}
	for _, attr := range requiredAttributes {
		if _, exists := schemaResp.Schema.Attributes[attr]; !exists {
			t.Errorf("Site group resource schema should include %s attribute", attr)
		}
	}

	// Check that optional attributes are present
	optionalAttributes := []string{"parent", "description", "tags", "custom_fields"}
	for _, attr := range optionalAttributes {
		if _, exists := schemaResp.Schema.Attributes[attr]; !exists {
			t.Errorf("Site group resource schema should include %s attribute", attr)
		}
	}
}

func TestSiteGroupResourceMetadata(t *testing.T) {
	ctx := context.Background()
	r := NewSiteGroupResource()

	// Test that the resource metadata can be retrieved
	metadataReq := resource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResp := &resource.MetadataResponse{}

	r.Metadata(ctx, metadataReq, metadataResp)

	expectedTypeName := "netbox_site_group"
	if metadataResp.TypeName != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)
	}
}

func TestSiteGroupResourceConfigure(t *testing.T) {
	ctx := context.Background()
	r := NewSiteGroupResource().(*SiteGroupResource) // Cast to access Configure method

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
