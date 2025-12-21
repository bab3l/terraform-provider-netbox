package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestModuleBayResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayResource()

	if r == nil {

		t.Fatal("Expected non-nil ModuleBay resource")

	}

}

func TestModuleBayResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"device", "name"},

		Optional: []string{"label", "position", "description", "tags", "custom_fields"},

		Computed: []string{"id"},

		OptionalComputed: []string{},
	})

}

func TestModuleBayResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_module_bay")

}

func TestModuleBayResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayResource()

	testutil.ValidateResourceConfigure(t, r)

}
