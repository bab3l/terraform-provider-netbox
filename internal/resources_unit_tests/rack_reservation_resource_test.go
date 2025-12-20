package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRackReservationResource(t *testing.T) {
	t.Parallel()

	r := resources.NewRackReservationResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestRackReservationResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewRackReservationResource()
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
		Required: []string{"rack", "units", "user", "description"},
		Optional: []string{"tenant", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestRackReservationResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewRackReservationResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rack_reservation")
}

func TestRackReservationResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewRackReservationResource()
	testutil.ValidateResourceConfigure(t, r)
}
