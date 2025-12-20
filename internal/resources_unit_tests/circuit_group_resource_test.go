package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCircuitGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil CircuitGroup resource")

	}

}

func TestCircuitGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

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

		Optional: []string{"description", "tenant", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestCircuitGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_circuit_group")

}

func TestCircuitGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

	testutil.ValidateResourceConfigure(t, r)

}
