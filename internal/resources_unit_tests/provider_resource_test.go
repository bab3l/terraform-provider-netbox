package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProviderResource(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestProviderResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

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

		Optional: []string{"description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestProviderResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_provider")

}

func TestProviderResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

	testutil.ValidateResourceConfigure(t, r)

}
