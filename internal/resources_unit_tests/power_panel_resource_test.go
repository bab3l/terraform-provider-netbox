package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPowerPanelResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewPowerPanelResource()
	if r == nil {
		t.Fatal("Expected non-nil PowerPanel resource")
	}
}

func TestPowerPanelResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewPowerPanelResource()
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
		Required: []string{"site", "name"},
		Optional: []string{"location", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestPowerPanelResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewPowerPanelResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_power_panel")
}

func TestPowerPanelResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewPowerPanelResource().(*resources.PowerPanelResource)
	testutil.ValidateResourceConfigure(t, r)
}
