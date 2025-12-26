package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPowerPortResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortResource()
	if r == nil {
		t.Fatal("Expected non-nil PowerPort resource")
	}
}

func TestPowerPortResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortResource()
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
		Optional: []string{"label", "type", "maximum_draw", "allocated_draw", "description", "mark_connected", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestPowerPortResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_power_port")
}

func TestPowerPortResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortResource().(*resources.PowerPortResource)
	testutil.ValidateResourceConfigure(t, r)
}
