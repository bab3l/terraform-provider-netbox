package datasources_unit_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestContactRoleDataSource_Metadata(t *testing.T) {
	d := datasources.NewContactRoleDataSource()
	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}
	d.Metadata(nil, req, resp)

	if resp.TypeName != "netbox_contact_role" {
		t.Errorf("expected TypeName 'netbox_contact_role', got '%s'", resp.TypeName)
	}
}

func TestContactRoleDataSource_Schema(t *testing.T) {
	d := datasources.NewContactRoleDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(nil, req, resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("expected schema attributes, got nil")
	}

	lookupAttrs := []string{"id", "name", "slug"}
	computedAttrs := []string{"description"}

	for _, attr := range lookupAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected lookup attribute '%s' in schema", attr)
		}
	}

	for _, attr := range computedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected computed attribute '%s' in schema", attr)
		}
	}
}

func TestContactRoleDataSource_SchemaDescription(t *testing.T) {
	d := datasources.NewContactRoleDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(nil, req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("expected schema to have a description")
	}
}

func TestContactRoleDataSource_Configure(t *testing.T) {
	d := datasources.NewContactRoleDataSource().(*datasources.ContactRoleDataSource)
	req := datasource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &datasource.ConfigureResponse{}
	d.Configure(nil, req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors with nil provider data, got: %v", resp.Diagnostics)
	}
}
