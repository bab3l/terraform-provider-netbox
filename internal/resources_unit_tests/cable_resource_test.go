package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCableResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCableResource()

	if r == nil {

		t.Fatal("Expected non-nil cable resource")

	}

}

func TestCableResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCableResource()

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

		Required: []string{"a_terminations", "b_terminations"},

		Optional: []string{

			"type", "status", "tenant", "label", "color",

			"length", "length_unit", "description", "comments",

			"tags", "custom_fields",
		},

		Computed: []string{"id"},
	})

}

func TestCableResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCableResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_cable")

}

func TestCableResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewCableResource()

	testutil.ValidateResourceConfigure(t, r)

}
