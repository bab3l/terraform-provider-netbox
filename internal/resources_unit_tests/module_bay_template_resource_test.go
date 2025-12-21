package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestModuleBayTemplateResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil ModuleBayTemplate resource")

	}

}

func TestModuleBayTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	// Either device_type or module_type is required, but both are optional at schema level

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name"},

		Optional: []string{"device_type", "module_type", "label", "position", "description"},

		Computed: []string{"id"},

		OptionalComputed: []string{},
	})

}

func TestModuleBayTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_module_bay_template")

}

func TestModuleBayTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleBayTemplateResource()

	testutil.ValidateResourceConfigure(t, r)

}
