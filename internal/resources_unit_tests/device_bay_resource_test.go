package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeviceBayResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestDeviceBayResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"device", "name"},
		Optional: []string{"label", "description", "installed_device", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestDeviceBayResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_device_bay")
}

func TestDeviceBayResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()
	testutil.ValidateResourceConfigure(t, r)
}
