package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestASNResource(t *testing.T) {

	t.Parallel()

	r := resources.NewASNResource()

	if r == nil {

		t.Fatal("Expected non-nil ASN resource")

	}

}

func TestASNResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewASNResource()

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

		Required: []string{"asn"},

		Optional: []string{"rir", "tenant", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestASNResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewASNResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_asn")

}

func TestASNResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewASNResource()

	testutil.ValidateResourceConfigure(t, r)

}
