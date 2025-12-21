package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestInventoryItemResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestInventoryItemResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}
}

func TestInventoryItemResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_inventory_item")
}

func TestInventoryItemResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemResource()
	testutil.ValidateResourceConfigure(t, r)
}
