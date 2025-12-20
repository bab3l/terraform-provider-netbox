package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTenantGroupResource(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	if r == nil {
		t.Fatal("Expected non-nil tenant group resource")
	}
}

func TestTenantGroupResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "slug"},
		Optional: []string{"parent", "description", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestTenantGroupResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_tenant_group")
}

func TestTenantGroupResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	testutil.ValidateResourceConfigure(t, r)
}
