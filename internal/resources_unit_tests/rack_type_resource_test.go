package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRackTypeResource(t *testing.T) {
	t.Parallel()

	r := resources.NewRackTypeResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestRackTypeResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewRackTypeResource()
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
		Required: []string{"manufacturer", "model", "slug"},
		Optional: []string{"description", "form_factor", "width", "u_height", "starting_unit", "desc_units", "outer_width", "outer_depth", "outer_unit", "weight", "max_weight", "weight_unit", "mounting_depth", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestRackTypeResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewRackTypeResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rack_type")
}

func TestRackTypeResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewRackTypeResource()
	testutil.ValidateResourceConfigure(t, r)
}
