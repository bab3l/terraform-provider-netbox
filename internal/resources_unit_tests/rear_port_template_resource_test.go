package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRearPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestRearPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "type"},

		Optional: []string{"device_type", "module_type", "label", "color", "positions", "description"},

		Computed: []string{"id"},
	})

}

func TestRearPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rear_port_template")

}

func TestRearPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	testutil.ValidateResourceConfigure(t, r)

}
