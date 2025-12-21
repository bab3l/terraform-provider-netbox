package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestContactAssignmentResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactAssignmentResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestContactAssignmentResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactAssignmentResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"object_type", "object_id", "contact_id"},

		Optional: []string{"role_id", "priority", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestContactAssignmentResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactAssignmentResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_contact_assignment")

}

func TestContactAssignmentResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactAssignmentResource()

	testutil.ValidateResourceConfigure(t, r)

}
