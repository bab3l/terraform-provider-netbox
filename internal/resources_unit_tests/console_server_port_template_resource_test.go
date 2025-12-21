package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConsoleServerPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestConsoleServerPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

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

func TestConsoleServerPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_console_server_port_template")

}

func TestConsoleServerPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

	testutil.ValidateResourceConfigure(t, r)

}
