package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPowerOutletTemplateResource(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerOutletTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil PowerOutletTemplate resource")
	}
}

func TestPowerOutletTemplateResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerOutletTemplateResource()
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
		Required: []string{"name"},
		Optional: []string{"device_type", "module_type", "label", "type", "power_port", "feed_leg", "description"},
		Computed: []string{"id"},
	})
}

func TestPowerOutletTemplateResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerOutletTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_power_outlet_template")
}

func TestPowerOutletTemplateResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewPowerOutletTemplateResource().(*resources.PowerOutletTemplateResource)
	testutil.ValidateResourceConfigure(t, r)
}
