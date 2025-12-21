package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestFrontPortResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewFrontPortResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestFrontPortResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewFrontPortResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"device", "name", "type", "rear_port"},
		Optional: []string{"label", "color", "rear_port_position", "description", "mark_connected", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestFrontPortResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewFrontPortResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_front_port")
}

func TestFrontPortResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewFrontPortResource()
	testutil.ValidateResourceConfigure(t, r)
}
