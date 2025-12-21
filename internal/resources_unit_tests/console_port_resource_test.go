package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConsolePortResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestConsolePortResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"device", "name"},

		Optional: []string{"custom_fields", "description", "label", "mark_connected", "speed", "tags", "type"},

		Computed: []string{"id"},
	})

}

func TestConsolePortResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_console_port")

}

func TestConsolePortResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortResource()

	testutil.ValidateResourceConfigure(t, r)

}
