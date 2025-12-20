package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPowerOutletResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource()

	if r == nil {

		t.Fatal("Expected non-nil PowerOutlet resource")

	}

}

func TestPowerOutletResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource()

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

		Required: []string{"device", "name"},

		Optional: []string{"label", "type", "power_port", "feed_leg", "description", "mark_connected", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestPowerOutletResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_power_outlet")

}

func TestPowerOutletResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource().(*resources.PowerOutletResource)

	testutil.ValidateResourceConfigure(t, r)

}
