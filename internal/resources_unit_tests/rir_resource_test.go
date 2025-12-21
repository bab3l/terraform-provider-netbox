package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRIRResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestRIRResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "slug"},

		Optional: []string{"is_private", "description", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestRIRResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rir")

}

func TestRIRResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

	testutil.ValidateResourceConfigure(t, r)

}
