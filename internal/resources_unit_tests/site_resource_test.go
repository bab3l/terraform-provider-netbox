package resources_unit_tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSiteResource(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteResource()
	if r == nil {
		t.Fatal("Expected non-nil Site resource")
	}
}

func TestSiteResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteResource()
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
		Optional: []string{"status", "description", "comments", "facility", "time_zone", "physical_address", "shipping_address", "latitude", "longitude", "tags", "custom_fields"},
		Computed: []string{"id"},
	})

	testutil.ValidateFloat64AttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["latitude"],
		"latitude",
		reflect.TypeOf(validators.LatitudeValidator{}),
	)
	testutil.ValidateFloat64AttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["longitude"],
		"longitude",
		reflect.TypeOf(validators.LongitudeValidator{}),
	)
}

func TestSiteResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_site")
}

func TestSiteResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteResource()
	testutil.ValidateResourceConfigure(t, r)
}
