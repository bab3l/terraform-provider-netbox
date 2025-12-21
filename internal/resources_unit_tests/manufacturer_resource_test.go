package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestManufacturerResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewManufacturerResource()
	if r == nil {
		t.Fatal("Expected non-nil manufacturer resource")
	}
}

func TestManufacturerResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewManufacturerResource()
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
	})
}

func TestManufacturerResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewManufacturerResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_manufacturer")
}

func TestManufacturerResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewManufacturerResource()
	testutil.ValidateResourceConfigure(t, r)
}
