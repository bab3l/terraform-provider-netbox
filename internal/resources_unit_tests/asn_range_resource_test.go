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

func TestASNRangeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

	if r == nil {

		t.Fatal("Expected non-nil ASNRange resource")

	}

}

func TestASNRangeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

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

		Required: []string{"name", "slug", "rir", "start", "end"},

		Optional: []string{"tenant", "description", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

	testutil.ValidateStringAttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["start"],
		"start",
		reflect.TypeOf(validators.ASNStringValidator{}),
	)
	testutil.ValidateStringAttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["end"],
		"end",
		reflect.TypeOf(validators.ASNStringValidator{}),
	)

}

func TestASNRangeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_asn_range")

}

func TestASNRangeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

	testutil.ValidateResourceConfigure(t, r)

}
