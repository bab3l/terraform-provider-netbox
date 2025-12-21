package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCircuitResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCircuitResource()

	if r == nil {

		t.Fatal("Expected non-nil Circuit resource")

	}

}

func TestCircuitResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCircuitResource()

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

		Required: []string{"cid", "circuit_provider", "type"},

		Optional: []string{"status", "tenant", "install_date", "termination_date", "commit_rate", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestCircuitResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCircuitResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_circuit")

}

func TestCircuitResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCircuitResource()

	testutil.ValidateResourceConfigure(t, r)

}
