package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestUserDataSource(t *testing.T) {
	d := datasources.NewUserDataSource()
	if d == nil {
		t.Fatal("User data source should not be nil")
	}
}

func TestUserDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewUserDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "username", "first_name", "last_name", "email", "is_staff", "is_active"}
	for _, attr := range expectedAttrs {
		if _, ok := schema.Attributes[attr]; !ok {
			t.Errorf("Schema should have '%s' attribute", attr)
		}
	}

	// Verify that lookup fields are optional
	idAttr := schema.Attributes["id"]
	if !idAttr.IsOptional() {
		t.Error("'id' attribute should be optional for lookup")
	}
	usernameAttr := schema.Attributes["username"]
	if !usernameAttr.IsOptional() {
		t.Error("'username' attribute should be optional for lookup")
	}
}

func TestUserDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewUserDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_user"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
