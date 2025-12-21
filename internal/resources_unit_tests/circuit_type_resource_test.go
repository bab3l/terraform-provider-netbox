package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCircuitTypeResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewCircuitTypeResource()
	if r == nil {
		t.Fatal("Expected non-nil CircuitType resource")
	}
}

func TestCircuitTypeResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewCircuitTypeResource()
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
		Computed: []string{"id"},
		Optional: []string{"description", "color", "tags", "custom_fields"},
	})
}

func TestCircuitTypeResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewCircuitTypeResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_circuit_type")
}

func TestCircuitTypeResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewCircuitTypeResource()
	testutil.ValidateResourceConfigure(t, r)
}
