package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConsoleServerPortResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestConsoleServerPortResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

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

func TestConsoleServerPortResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_console_server_port")

}

func TestConsoleServerPortResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

	testutil.ValidateResourceConfigure(t, r)

}
