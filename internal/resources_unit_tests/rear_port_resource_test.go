package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRearPortResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestRearPortResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"device", "name", "type"},

		Optional: []string{"label", "color", "positions", "description", "mark_connected", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestRearPortResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rear_port")

}

func TestRearPortResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortResource()

	testutil.ValidateResourceConfigure(t, r)

}
