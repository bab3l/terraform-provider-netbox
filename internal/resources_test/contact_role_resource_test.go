package resources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestContactRoleResource_Metadata(t *testing.T) {
	r := resources.NewContactRoleResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &resource.MetadataResponse{}
	r.Metadata(nil, req, resp)

	if resp.TypeName != "netbox_contact_role" {
		t.Errorf("expected TypeName 'netbox_contact_role', got '%s'", resp.TypeName)
	}
}

func TestContactRoleResource_Schema(t *testing.T) {
	r := resources.NewContactRoleResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(nil, req, resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("expected schema attributes, got nil")
	}

	requiredAttrs := []string{"name", "slug"}
	optionalAttrs := []string{"description", "tags", "custom_fields"}
	computedAttrs := []string{"id"}

	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected required attribute '%s' in schema", attr)
		}
	}

	for _, attr := range optionalAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected optional attribute '%s' in schema", attr)
		}
	}

	for _, attr := range computedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected computed attribute '%s' in schema", attr)
		}
	}
}

func TestContactRoleResource_SchemaDescription(t *testing.T) {
	r := resources.NewContactRoleResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(nil, req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("expected schema to have a description")
	}
}

func TestContactRoleResource_Configure(t *testing.T) {
	r := resources.NewContactRoleResource().(*resources.ContactRoleResource)
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}
	r.Configure(nil, req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors with nil provider data, got: %v", resp.Diagnostics)
	}
}
