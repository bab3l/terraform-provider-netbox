package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPowerPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil PowerPortTemplate resource")
	}
}

func TestPowerPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource()
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
		Optional: []string{"device_type", "module_type", "label", "type", "maximum_draw", "allocated_draw", "description"},
		Computed: []string{"id"},
	})
}

func TestPowerPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_power_port_template")
}

func TestPowerPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource().(*resources.PowerPortTemplateResource)
	testutil.ValidateResourceConfigure(t, r)
}
