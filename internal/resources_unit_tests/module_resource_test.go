package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestModuleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleResource()

	if r == nil {

		t.Fatal("Expected non-nil Module resource")

	}

}

func TestModuleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"device", "module_bay", "module_type"},

		Optional: []string{"serial", "asset_tag", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id", "status"},

		OptionalComputed: []string{},
	})

}

func TestModuleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_module")

}

func TestModuleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleResource()

	testutil.ValidateResourceConfigure(t, r)

}
