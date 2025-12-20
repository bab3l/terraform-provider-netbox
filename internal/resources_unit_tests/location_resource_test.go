package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestLocationResource(t *testing.T) {
	t.Parallel()
	r := resources.NewLocationResource()
	if r == nil {
		t.Fatal("Expected non-nil Location resource")
	}
}

func TestLocationResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewLocationResource()
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
		Required: []string{"name", "slug", "site"},
		Optional: []string{"parent", "status", "tenant", "facility", "description", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestLocationResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewLocationResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_location")
}

func TestLocationResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewLocationResource()
	testutil.ValidateResourceConfigure(t, r)
}
