package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPowerFeedResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerFeedResource()

	if r == nil {

		t.Fatal("Expected non-nil PowerFeed resource")

	}

}

func TestPowerFeedResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerFeedResource()

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

		Required: []string{"power_panel", "name"},

		Optional: []string{"rack", "mark_connected", "description", "tenant", "comments", "tags", "custom_fields"},

		Computed: []string{"id", "status", "type", "supply", "phase", "voltage", "amperage", "max_utilization"},
	})

}

func TestPowerFeedResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerFeedResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_power_feed")

}

func TestPowerFeedResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerFeedResource().(*resources.PowerFeedResource)

	testutil.ValidateResourceConfigure(t, r)

}
