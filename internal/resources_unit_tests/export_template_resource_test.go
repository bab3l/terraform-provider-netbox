package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestExportTemplateResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewExportTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestExportTemplateResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewExportTemplateResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"object_types", "template_code", "name"},
		Optional: []string{"description", "mime_type", "file_extension", "as_attachment"},
		Computed: []string{"id"},
	})
}

func TestExportTemplateResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewExportTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_export_template")
}

func TestExportTemplateResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewExportTemplateResource()
	testutil.ValidateResourceConfigure(t, r)
}
