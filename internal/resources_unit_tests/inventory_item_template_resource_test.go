package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestInventoryItemTemplateResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestInventoryItemTemplateResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemTemplateResource()
	schemaRequest := &resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}
	r.Schema(context.Background(), *schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"device_type", "name"},
		Optional: []string{"parent", "label", "role", "manufacturer", "part_id", "description", "component_type", "component_id"},
		Computed: []string{"id"},
	})
}

func TestInventoryItemTemplateResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_inventory_item_template")
}

func TestInventoryItemTemplateResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemTemplateResource()
	testutil.ValidateResourceConfigure(t, r)
}
