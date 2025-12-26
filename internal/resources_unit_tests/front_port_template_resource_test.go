package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestFrontPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestFrontPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "type", "rear_port"},
		Optional: []string{"device_type", "module_type", "label", "color", "rear_port_position", "description"},
		Computed: []string{"id"},
	})
}

func TestFrontPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_front_port_template")
}

func TestFrontPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()
	testutil.ValidateResourceConfigure(t, r)
}
