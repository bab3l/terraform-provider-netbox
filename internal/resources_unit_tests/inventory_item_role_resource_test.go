package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestInventoryItemRoleResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemRoleResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestInventoryItemRoleResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemRoleResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}
}

func TestInventoryItemRoleResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemRoleResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_inventory_item_role")
}

func TestInventoryItemRoleResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInventoryItemRoleResource()
	testutil.ValidateResourceConfigure(t, r)
}
