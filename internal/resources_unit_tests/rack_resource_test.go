package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRackResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRackResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestRackResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRackResource()

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

		Required: []string{"name", "site"},

		Optional: []string{"id", "location", "tenant", "status", "role", "rack_type", "serial", "asset_tag", "facility_id", "form_factor", "width", "u_height", "starting_unit", "weight", "max_weight", "weight_unit", "desc_units", "outer_width", "outer_depth", "outer_unit", "mounting_depth", "airflow", "description", "comments", "tags", "custom_fields"},

		Computed: []string{},
	})

}

func TestRackResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRackResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rack")

}

func TestRackResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRackResource()

	testutil.ValidateResourceConfigure(t, r)

}
