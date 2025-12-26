package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProviderNetworkResource(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderNetworkResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestProviderNetworkResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderNetworkResource()

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

		Required: []string{"circuit_provider", "name"},

		Optional: []string{"service_id", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestProviderNetworkResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderNetworkResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_provider_network")

}

func TestProviderNetworkResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderNetworkResource()

	testutil.ValidateResourceConfigure(t, r)

}
