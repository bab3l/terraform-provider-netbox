package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConsolePortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestConsolePortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name"},

		Optional: []string{"description", "device_type", "label", "module_type", "type"},

		Computed: []string{"id"},
	})

}

func TestConsolePortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_console_port_template")

}

func TestConsolePortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

	testutil.ValidateResourceConfigure(t, r)

}
