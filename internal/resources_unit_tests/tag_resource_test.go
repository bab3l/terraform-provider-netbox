package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTagResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTagResource()
	if r == nil {
		t.Fatal("Expected non-nil Tag resource")
	}
}

func TestTagResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTagResource()
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
		Optional: []string{"color", "description", "object_types"},
		Computed: []string{"id"},
	})
}

func TestTagResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTagResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_tag")
}

func TestTagResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTagResource()
	testutil.ValidateResourceConfigure(t, r)
}
