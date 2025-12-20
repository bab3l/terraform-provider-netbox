package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCircuitTerminationResource(t *testing.T) {
	t.Parallel()

	r := resources.NewCircuitTerminationResource()

	if r == nil {
		t.Fatal("Expected non-nil CircuitTermination resource")
	}
}

func TestCircuitTerminationResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewCircuitTerminationResource()

	schemaRequest := resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	validation := testutil.SchemaValidation{
		Required: []string{"circuit", "term_side"},
		Optional: []string{"site", "provider_network", "port_speed", "upstream_speed", "xconnect_id", "pp_info", "description", "mark_connected", "tags", "custom_fields"},
		Computed: []string{"id"},
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, validation)
}

func TestCircuitTerminationResourceMetadata(t *testing.T) {
	t.Parallel()
	testutil.ValidateResourceMetadata(t, resources.NewCircuitTerminationResource(), "netbox", "netbox_circuit_termination")
}

func TestCircuitTerminationResourceConfigure(t *testing.T) {
	t.Parallel()
	testutil.ValidateResourceConfigure(t, resources.NewCircuitTerminationResource())
}
