package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestModuleTypeResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleTypeResource()

	if r == nil {

		t.Fatal("Expected non-nil ModuleType resource")

	}

}

func TestModuleTypeResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleTypeResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"manufacturer", "model"},

		Optional: []string{"part_number", "airflow", "weight", "weight_unit", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},

		OptionalComputed: []string{},
	})

}

func TestModuleTypeResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleTypeResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_module_type")

}

func TestModuleTypeResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewModuleTypeResource()

	testutil.ValidateResourceConfigure(t, r)

}
