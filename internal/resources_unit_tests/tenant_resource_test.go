package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTenantResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTenantResource()
	if r == nil {
		t.Fatal("Expected non-nil tenant resource")
	}
}

func TestTenantResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTenantResource()
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
		Optional: []string{"group", "description", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestTenantResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTenantResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_tenant")
}

func TestTenantResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTenantResource()
	testutil.ValidateResourceConfigure(t, r)
}
