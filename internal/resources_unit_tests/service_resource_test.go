package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestServiceResource(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceResource()
	if r == nil {
		t.Fatal("Expected non-nil Service resource")
	}
}

func TestServiceResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceResource()
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
		Required: []string{"name", "protocol", "ports"},
		Optional: []string{"device", "virtual_machine", "ipaddresses", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestServiceResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_service")
}

func TestServiceResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceResource()
	testutil.ValidateResourceConfigure(t, r)
}
