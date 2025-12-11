// Package resources_test provides unit tests for the NetBox Terraform provider resources.
package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestWirelessLinkResourceMetadata(t *testing.T) {
	r := resources.NewWirelessLinkResource()
	req := resource.MetadataRequest{ProviderTypeName: "netbox"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)

	expected := "netbox_wireless_link"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %q, got %q", expected, resp.TypeName)
	}
}

func TestWirelessLinkResourceSchema(t *testing.T) {
	r := resources.NewWirelessLinkResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)

	// Check that response has no diagnostics errors
	if resp.Diagnostics.HasError() {
		t.Errorf("Schema returned errors: %v", resp.Diagnostics)
	}

	// Verify required attributes exist
	requiredAttrs := []string{"interface_a", "interface_b"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected required attribute %q not found in schema", attr)
		}
	}

	// Verify optional attributes exist
	optionalAttrs := []string{
		"id", "ssid", "status", "tenant", "auth_type", "auth_cipher", "auth_psk",
		"distance", "distance_unit", "description", "comments", "tags", "custom_fields",
	}
	for _, attr := range optionalAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected attribute %q not found in schema", attr)
		}
	}
}
