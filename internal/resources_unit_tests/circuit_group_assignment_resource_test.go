package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCircuitGroupAssignmentResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupAssignmentResource()

	if r == nil {

		t.Fatal("Expected non-nil CircuitGroupAssignment resource")

	}

}

func TestCircuitGroupAssignmentResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupAssignmentResource()

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

		Required: []string{"group_id", "circuit_id"},

		Optional: []string{"priority", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestCircuitGroupAssignmentResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupAssignmentResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_circuit_group_assignment")

}

func TestCircuitGroupAssignmentResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupAssignmentResource()

	testutil.ValidateResourceConfigure(t, r)

}
