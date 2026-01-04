package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestInterfaceTemplateResource(t *testing.T) {
	t.Parallel()

	r := resources.NewInterfaceTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestInterfaceTemplateResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewInterfaceTemplateResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "type"},
		Optional: []string{"device_type", "module_type", "label", "enabled", "mgmt_only", "description", "bridge", "poe_mode", "poe_type", "rf_role"},
		Computed: []string{"id"},
	})
}

func TestInterfaceTemplateResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewInterfaceTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_interface_template")
}

func TestInterfaceTemplateResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewInterfaceTemplateResource()
	testutil.ValidateResourceConfigure(t, r)
}
