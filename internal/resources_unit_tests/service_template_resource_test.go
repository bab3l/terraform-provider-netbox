package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestServiceTemplateResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewServiceTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil Service Template resource")
	}
}

func TestServiceTemplateResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewServiceTemplateResource()
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
		Required: []string{"name", "ports"},
		Optional: []string{"protocol", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestServiceTemplateResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewServiceTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_service_template")
}

func TestServiceTemplateResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewServiceTemplateResource()
	testutil.ValidateResourceConfigure(t, r)
}
