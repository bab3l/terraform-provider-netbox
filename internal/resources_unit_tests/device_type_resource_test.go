package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeviceTypeResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewDeviceTypeResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestDeviceTypeResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewDeviceTypeResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"manufacturer", "model", "slug"},
		Optional: []string{"description", "comments", "part_number", "u_height", "is_full_depth", "subdevice_role", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestDeviceTypeResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewDeviceTypeResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_device_type")
}

func TestDeviceTypeResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewDeviceTypeResource()
	testutil.ValidateResourceConfigure(t, r)
}
