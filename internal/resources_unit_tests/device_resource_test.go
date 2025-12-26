package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeviceResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestDeviceResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"site", "device_type", "role"},
		Optional: []string{"name", "status", "description", "comments", "tenant", "platform", "serial", "asset_tag", "rack", "position", "face", "latitude", "longitude", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestDeviceResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_device")
}

func TestDeviceResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	testutil.ValidateResourceConfigure(t, r)
}
