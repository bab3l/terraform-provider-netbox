package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConfigContextResource_Metadata(t *testing.T) {
	r := resources.NewConfigContextResource()

	req := resource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "netbox_config_context" {
		t.Errorf("Expected type name 'netbox_config_context', got '%s'", resp.TypeName)
	}
}

func TestConfigContextResource_Schema(t *testing.T) {
	r := resources.NewConfigContextResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	// Verify required attributes exist
	requiredAttrs := []string{"id", "name", "data"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected attribute '%s' to exist in schema", attr)
		}
	}

	// Verify optional attributes exist
	optionalAttrs := []string{
		"description", "weight", "is_active",
		"regions", "site_groups", "sites", "locations",
		"device_types", "roles", "platforms",
		"cluster_types", "cluster_groups", "clusters",
		"tenant_groups", "tenants", "tags",
	}
	for _, attr := range optionalAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected attribute '%s' to exist in schema", attr)
		}
	}
}

func TestConfigContextResource_SchemaDescription(t *testing.T) {
	r := resources.NewConfigContextResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("Expected schema to have a markdown description")
	}
}

func TestConfigContextResource_Configure(t *testing.T) {
	r := resources.NewConfigContextResource()

	// Verify the resource implements the configurable interface
	configurable, ok := r.(resource.ResourceWithConfigure)
	if !ok {
		t.Skip("Resource does not implement ResourceWithConfigure")
	}

	// Test with nil provider data - should not error
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	configurable.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Configure with nil provider data should not error: %v", resp.Diagnostics)
	}
}
